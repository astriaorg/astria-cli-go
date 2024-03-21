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
		Application:  tviewApp,
		flex:         flex,
		processPanes: processPanes,
		// select first pane on start
		selectedPaneIdx: 0,
	}
}

// Start starts the tview application.
func (a *App) Start() {
	// start scanning stdout for each process
	for _, pr := range a.processPanes {
		pr.StartScan()
	}

	// keyboard shortcuts
	a.setInputCapture(a.getKeyboardMainView)

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
		pp.pr.Stop()
	}
	// stop the tview application
	a.Application.Stop()
}

// setIsFullscreen sets the view to fullscreen view or main view depending on the bool.
func (a *App) setIsFullscreen(isFullscreen bool) {
	if isFullscreen {
		// set prevFlex to current flex
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
			// TODO - autoscroll
		case 'q':
			a.stop()
		case 'w':
			for _, pp := range a.processPanes {
				pp.ToggleIsWordWrapped()
			}
		}
	case tcell.KeyDown:
		// we want the down key to increment the selected index, bc top starts at 0
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
		// CtrlC should always stop the app no matter the view
		a.stop()
	case tcell.KeyRune:
		switch evt.Rune() {
		case 'a':
			// TODO - autoscroll
		case 'q':
			a.setIsFullscreen(false)
		case 'w':
			for _, pp := range a.processPanes {
				pp.ToggleIsWordWrapped()
			}
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

// redraw updates the panes to show visual treatment for selected pane.
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
