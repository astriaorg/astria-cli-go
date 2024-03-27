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

	selectedPaneIdx int
}

// NewMainView creates a new MainView with the given tview.Application and ProcessPanes.
func NewMainView(tApp *tview.Application, processrunners []processrunner.ProcessRunner) *MainView {
	// create process panes for the runners
	var processPanes []*ProcessPane
	for _, pr := range processrunners {
		pp := NewProcessPane(tApp, pr)
		// start scanning the stdout of the panes
		processPanes = append(processPanes, pp)
		// start scanning the stdout of the panes
		pp.StartScan()
		// set the defaults
		pp.SetIsWordWrap(false)
		pp.SetIsAutoScroll(true)
		pp.SetIsBorderless(false)
	}

	return &MainView{
		tApp:            tApp,
		processPanes:    processPanes,
		selectedPaneIdx: 0,
	}
}

// Render returns the tview.Flex that represents the MainView.
func (mv *MainView) Render(_ Props) *tview.Flex {
	innerFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	for _, pp := range mv.processPanes {
		innerFlex.AddItem(pp.GetTextView(), 0, 1, true)
	}

	mainWindowHelpInfo := tview.NewTextView().SetText(MainLegendText)
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
			switch evt.Rune() {
			case 'a':
				for _, pp := range mv.processPanes {
					pp.SetIsAutoScroll(!pp.isAutoScroll)
				}
			case 'q':
				a.Exit()
			case 'w':
				for _, pp := range mv.processPanes {
					pp.SetIsWordWrap(!pp.isWordWrap)
				}
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
}

// NewFullscreenView creates a new FullscreenView with the given tview.Application and ProcessPane.
func NewFullscreenView(tApp *tview.Application, processPane *ProcessPane) *FullscreenView {
	return &FullscreenView{
		tApp:        tApp,
		processPane: processPane,
	}
}

// Render returns the tview.Flex that represents the FullscreenView.
func (fv *FullscreenView) Render(p Props) *tview.Flex {
	fv.processPane = p.(*ProcessPane)
	// build tview text views and flex
	help := tview.NewTextView().
		SetDynamicColors(true).
		SetText(FullscreenLegendText).
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
		fv.processPane.SetIsBorderless(false)
		a.SetView("main", nil)
	}
	return func(evt *tcell.EventKey) *tcell.EventKey {
		switch evt.Key() {
		case tcell.KeyCtrlC:
			a.Exit()
		case tcell.KeyRune:
			switch evt.Rune() {
			case 'a':
				fv.processPane.SetIsAutoScroll(!fv.processPane.isAutoScroll)
			case 'b':
				fv.processPane.SetIsBorderless(!fv.processPane.isBorderless)
			case 'q':
				backToMain()
			case 'w':
				fv.processPane.SetIsWordWrap(!fv.processPane.isWordWrap)
			}
		case tcell.KeyEscape:
			backToMain()
		}
		return evt
	}
}
