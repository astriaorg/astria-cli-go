package ui

import (
	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MainLegendText       = " (q)uit | (a)utoscroll | (w)ord wrap | (up/down) select pane | (enter) fullscreen"
	FullscreenLegendText = " (q/esc) back | (a)utoscroll | (w)ord wrap "
	MainTitle            = "Astria Dev"
)

// App contains the tview.Application and other necessary fields to manage the ui
type App struct {
	// Application is the tview application
	*tview.Application

	// flex is the top level flex view
	flex *tview.Flex
	// prevFlex is the previously shown flex view
	prevFlex *tview.Flex

	// processPanes holds ProcessPanes, one for each process.
	processPanes []*ProcessPane
	// selectedPaneIdx is the index of the currently selected ProcessPane
	selectedPaneIdx int

	// ui state
	isAutoScroll bool
	isWordWrap   bool
}

// NewApp creates a new tview.Application with the necessary components
func NewApp(processrunners []*processrunner.ProcessRunner) *App {
	tviewApp := tview.NewApplication()
	var processPanes []*ProcessPane

	// create ProcessPane for each process and adds it's textView to innerFlex
	innerFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	for _, pr := range processrunners {
		pp := NewProcessPane(tviewApp, pr)
		processPanes = append(processPanes, pp)
		innerFlex.AddItem(pp.textView, 0, 1, true)
	}

	// create main flex view and add help text and innerFlex
	mainWindowHelpInfo := tview.NewTextView().SetText(MainLegendText)
	flex := tview.NewFlex()
	flex.AddItem(innerFlex, 0, 4, false).
		SetDirection(tview.FlexRow).
		AddItem(mainWindowHelpInfo, 1, 0, false)
	flex.SetTitle(MainTitle).SetBorder(true)

	return &App{
		Application:     tviewApp,
		flex:            flex,
		processPanes:    processPanes,
		selectedPaneIdx: 0, // select first pane on start
		isAutoScroll:    true,
		isWordWrap:      false,
	}
}

// Start starts the tview application.
func (a *App) Start() {
	// start scanning stdout for each process
	for _, pp := range a.processPanes {
		pp.StartScan()
	}

	// keyboard shortcuts
	a.setInputCapture(a.getKeyboardMainView)

	// set initial process pane local state
	// NOTE - right now, we're keeping all the panes' local ui state in sync with
	//  the top level ui state defined in this struct. This means that toggling
	//  e.g. word wrap happens at the global level and affects all panes, even
	//  when in fullscreen view. This is a design decision that can be changed.
	for _, pp := range a.processPanes {
		pp.SetIsAutoScroll(a.isAutoScroll)
		pp.SetIsWordWrap(a.isWordWrap)
	}

	// redraw after startup to ensure ui state is correct
	a.redraw()

	// run the tview application
	if err := a.Application.Run(); err != nil {
		panic(err)
	}
}

// stop stops the tview application.
func (a *App) stop() {
	// stop each process
	for _, pp := range a.processPanes {
		pp.StopProcess()
	}
	// stop the tview application
	a.Application.Stop()
}

// setIsFullscreen sets the view to fullscreen view or main view depending on the bool.
func (a *App) setIsFullscreen(isFullscreen bool) {
	if isFullscreen {
		// set prevFlex to current flex
		// NOTE - we only have 2 views right now, so this is okay to do as an easy way to go between views,
		//  but creating a View struct or similar would be a good place to start refactoring.
		a.prevFlex = a.flex
		// build tview text views and flex
		help := tview.NewTextView().
			SetDynamicColors(true).
			SetText(FullscreenLegendText).
			SetChangedFunc(func() {
				a.Application.Draw()
			})
		selectedPane := a.processPanes[a.selectedPaneIdx]
		flex := tview.NewFlex().AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(selectedPane.GetTextView(), 0, 1, true).
			AddItem(help, 1, 0, false), 0, 4, false)
		a.setInputCapture(a.getKeyboardFullscreenView)
		// set current flex
		a.flex = flex
	} else {
		a.setInputCapture(a.getKeyboardMainView)
		// we only have 2 views right now, so this means we're going back to main view
		a.flex = a.prevFlex
	}
	a.redraw()
}

// setInputCapture sets the input capture function for the tview.Application.
// It must first set the input capture to nil before setting the new input capture.
func (a *App) setInputCapture(cb func(*tcell.EventKey) *tcell.EventKey) {
	a.Application.SetInputCapture(nil)
	a.Application.SetInputCapture(cb)
}

// getKeyboardMainView handles keyboard input for the main view.
func (a *App) getKeyboardMainView(evt *tcell.EventKey) *tcell.EventKey {
	switch evt.Key() {
	case tcell.KeyCtrlC:
		a.stop()
	case tcell.KeyRune:
		switch evt.Rune() {
		case 'a':
			a.toggleAutoScroll()
		case 'q':
			a.stop()
		case 'w':
			a.toggleWordWrap()
		}
	case tcell.KeyDown:
		// we want the down key to increment the selected index,
		// bc top pane starts at 0
		a.incrementSelectedPaneIdx()
	case tcell.KeyUp:
		a.decrementSelectedPaneIdx()
	case tcell.KeyEnter:
		a.setIsFullscreen(true)
	}
	return evt
}

// getKeyboardFullscreenView handles keyboard input for the fullscreen view.
func (a *App) getKeyboardFullscreenView(evt *tcell.EventKey) *tcell.EventKey {
	switch evt.Key() {
	case tcell.KeyCtrlC:
		// CtrlC should always completely stop the app no matter the view
		a.stop()
	case tcell.KeyRune:
		switch evt.Rune() {
		case 'a':
			a.toggleAutoScroll()
		case 'q':
			// q goes back to main view when on fullscreen view
			a.setIsFullscreen(false)
		case 'w':
			a.toggleWordWrap()
		}
	case tcell.KeyEscape:
		a.setIsFullscreen(false)
	}
	return evt
}

// incrementSelectedPaneIdx increments the selectedPaneIdx.
func (a *App) incrementSelectedPaneIdx() {
	a.selectedPaneIdx = (a.selectedPaneIdx + 1) % len(a.processPanes)
	a.redraw()
}

// decrementSelectedPaneIdx decrements the selectedPaneIdx.
func (a *App) decrementSelectedPaneIdx() {
	paneLen := len(a.processPanes)
	a.selectedPaneIdx = (a.selectedPaneIdx - 1 + paneLen) % paneLen
	a.redraw()
}

// toggleAutoScroll toggles the ui state and updates the ProcessPanes.
func (a *App) toggleAutoScroll() {
	a.isAutoScroll = !a.isAutoScroll
	for _, pp := range a.processPanes {
		pp.SetIsAutoScroll(a.isAutoScroll)
	}
}

// toggleWordWrap toggles the ui state and updates the ProcessPanes.
func (a *App) toggleWordWrap() {
	a.isWordWrap = !a.isWordWrap
	for _, pp := range a.processPanes {
		pp.SetIsWordWrap(a.isWordWrap)
	}
}

// redraw sets the tview's Root primitive and ensures correct visual state.
func (a *App) redraw() {
	// set the root primitive
	a.Application.SetRoot(a.flex, true)

	// ui treatment for the selected pane
	for idx, pp := range a.processPanes {
		if idx == a.selectedPaneIdx {
			pp.Highlight(true)
			a.Application.SetFocus(pp.textView)
		} else {
			pp.Highlight(false)
		}
	}
}
