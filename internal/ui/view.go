package ui

import (
	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MainLegendText       = " (q)uit | (a)utoscroll | (w)rap lines | (up/down) select pane | (enter) fullscreen selected pane"
	FullscreenLegendText = " (q/esc) back | (a)utoscroll | (w)rap lines | (b)orderless"
	MainTitle            = "Astria Dev"
)

// Props is an empty interface for passing data to the view.
type Props interface{}

type View interface {
	// Render returns the tview.Flex that represents the view.
	// The Props argument is used to pass data to the view.
	Render(p Props) *tview.Flex
	// GetKeyboard is a callback for defining keyboard shortcuts
	// FIXME - is there a way to avoid the App reference here?
	GetKeyboard(a AppController) func(evt *tcell.EventKey) *tcell.EventKey
}

// MainView represents the initial view when the app is started.
// It shows all the process panes in a vertical layout.
type MainView struct {
	// FIXME - how can we avoid having to have a reference of the tview.Application here?
	tApp         *tview.Application
	processPanes []*ProcessPane
	a            *App

	selectedPaneIdx int
}

// NewMainView creates a new MainView with the given tview.Application and ProcessPanes.
func NewMainView(tApp *tview.Application, processrunners []processrunner.ProcessRunner, a *App) *MainView {
	// create process panes for the runners
	var processPanes []*ProcessPane
	for _, pr := range processrunners {
		pp := NewProcessPane(tApp, pr, a.ss)
		// start scanning the stdout of the panes
		processPanes = append(processPanes, pp)
		// start scanning the stdout of the panes
		pp.StartScan()
		// set the defaults
		pp.SetIsWordWrap(a.ss.isWordWrap)
		pp.SetIsAutoScroll(a.ss.isAutoScroll)
		pp.SetIsBorderless(a.ss.isBorderless)
	}

	return &MainView{
		tApp:            tApp,
		processPanes:    processPanes,
		a:               a,
		selectedPaneIdx: 0,
	}
}

// Append the settings status to the end of the input string
func appendStatus(text string, status bool) string {
	if status {
		return text + ": [black:white]ON [-:-]"
	} else {
		return text + ": [white:darkslategray]off[-:-]"
	}
}

// Build the help text legened at the bottom of the main screen with dynamically
// changing setting status
func (mv *MainView) getHelpInfo() string {
	output := " "
	output += "(q)uit | "
	output += appendStatus("(a)utoscroll", mv.a.ss.isAutoScroll) + " | "
	output += appendStatus("(w)rap lines", mv.a.ss.isWordWrap) + " | "
	output += "(up/down) select pane | "
	output += "(enter) fullscreen selected pane"
	return output
}

// Render returns the tview.Flex that represents the MainView.
func (mv *MainView) Render(_ Props) *tview.Flex {
	innerFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	for _, pp := range mv.processPanes {
		// propigate the shared state to the process panes
		pp.SetIsAutoScroll(mv.a.ss.isAutoScroll)
		pp.SetIsBorderless(mv.a.ss.isBorderless)
		pp.SetIsWordWrap(mv.a.ss.isWordWrap)
		innerFlex.AddItem(pp.GetTextView(), 0, 1, true)
	}

	mainWindowHelpInfo := tview.NewTextView().SetDynamicColors(true).SetText(mv.getHelpInfo())
	flex := tview.NewFlex()
	flex.AddItem(innerFlex, 0, 4, false).
		SetDirection(tview.FlexRow).
		AddItem(mainWindowHelpInfo, 1, 0, false)
	flex.SetTitle(MainTitle).SetBorder(true)
	mv.redraw()

	return flex
}

// GetKeyboard is a callback for defining keyboard shortcuts.
func (mv *MainView) GetKeyboard(a AppController) func(evt *tcell.EventKey) *tcell.EventKey {
	return func(evt *tcell.EventKey) *tcell.EventKey {
		switch evt.Key() {
		case tcell.KeyCtrlC:
			a.Exit()
		case tcell.KeyRune:
			{
				switch evt.Rune() {
				case 'a':
					a.ToggleAutoscroll()
					for _, pp := range mv.processPanes {
						pp.SetIsAutoScroll(mv.a.ss.isAutoScroll)
					}
				case 'q':
					a.Exit()
					return nil
				case 'w':
					a.ToggleWordWrap()
					for _, pp := range mv.processPanes {
						pp.SetIsWordWrap(mv.a.ss.isWordWrap)
					}
				}
				a.RefreshView(nil)
			}
		case tcell.KeyDown:
			// we want the down key to increment the selected index,
			// bc top pane starts at 0
			mv.incrementSelectedPaneIdx()
		case tcell.KeyUp:
			mv.decrementSelectedPaneIdx()
		case tcell.KeyEnter:
			a.SetView("fullscreen", mv.processPanes[mv.selectedPaneIdx])
		}
		return evt
	}
}

// incrementSelectedPaneIdx increments the selectedPaneIdx.
func (mv *MainView) incrementSelectedPaneIdx() {
	mv.selectedPaneIdx = (mv.selectedPaneIdx + 1) % len(mv.processPanes)
	mv.redraw()
}

// decrementSelectedPaneIdx decrements the selectedPaneIdx.
func (mv *MainView) decrementSelectedPaneIdx() {
	paneLen := len(mv.processPanes)
	mv.selectedPaneIdx = (mv.selectedPaneIdx - 1 + paneLen) % paneLen
	mv.redraw()
}

// redraw ensures the correct visual state of the panes.
func (mv *MainView) redraw() {
	// ui treatment for the selected pane
	for idx, pp := range mv.processPanes {
		if idx == mv.selectedPaneIdx {
			pp.Highlight(true)
			mv.tApp.SetFocus(pp.textView)
		} else {
			pp.Highlight(false)
		}
	}
}

// FullscreenView represents the fullscreen view when a pane is selected.
type FullscreenView struct {
	tApp        *tview.Application
	processPane *ProcessPane
	a           *App
}

// NewFullscreenView creates a new FullscreenView with the given tview.Application and ProcessPane.
func NewFullscreenView(tApp *tview.Application, processPane *ProcessPane, a *App) *FullscreenView {
	return &FullscreenView{
		tApp:        tApp,
		processPane: processPane,
		a:           a,
	}
}

// Build the help text legened at the bottom of the fullscreen view with dynamically
// changing setting status
func (fv *FullscreenView) getHelpInfo() string {
	output := " "
	output += "(q/esc) back | "
	output += appendStatus("(a)utoscroll", fv.a.ss.isAutoScroll) + " | "
	output += appendStatus("(w)rap lines", fv.a.ss.isWordWrap) + " | "
	output += appendStatus("(b)orderless", fv.a.ss.isBorderless)
	return output
}

// Render returns the tview.Flex that represents the FullscreenView.
func (fv *FullscreenView) Render(p Props) *tview.Flex {
	fv.processPane = p.(*ProcessPane)
	// build tview text views and flex
	help := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fv.getHelpInfo()).
		SetChangedFunc(func() {
			fv.tApp.Draw()
		})
	flex := tview.NewFlex().AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(fv.processPane.GetTextView(), 0, 1, true).
		AddItem(help, 1, 0, false), 0, 4, false)
	return flex
}

// GetKeyboard is a callback for defining keyboard shortcuts.
func (fv *FullscreenView) GetKeyboard(a AppController) func(evt *tcell.EventKey) *tcell.EventKey {
	backToMain := func() {
		// reset borderless state before going back to main view
		a.ResetBorderless()
		fv.processPane.SetIsBorderless(fv.a.ss.isBorderless)
		// rerender the process Pane to apply all settings
		a.RefreshView(fv.processPane)
		// change views
		a.SetView("main", nil)
	}
	return func(evt *tcell.EventKey) *tcell.EventKey {
		switch evt.Key() {
		case tcell.KeyCtrlC:
			a.Exit()
			return nil
		case tcell.KeyRune:
			{
				switch evt.Rune() {
				case 'a':
					a.ToggleAutoscroll()
					fv.processPane.SetIsAutoScroll(fv.a.ss.isAutoScroll)

				case 'b':
					a.ToggleBorderless()
					fv.processPane.SetIsBorderless(fv.a.ss.isBorderless)

				case 'q':
					backToMain()
					return nil
				case 'w':
					a.ToggleWordWrap()
					fv.processPane.SetIsWordWrap(fv.a.ss.isWordWrap)

				}
				// needed to call the Render method again to refresh the help info
				a.RefreshView(fv.processPane)
				return nil
			}
		case tcell.KeyEscape:
			backToMain()
			return nil
		}
		return evt
	}
}
