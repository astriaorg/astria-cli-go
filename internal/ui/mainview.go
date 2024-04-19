package ui

import (
	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

const (
	MainTitle = "Astria Dev"
)

// MainView represents the initial view when the app is started.
// It shows all the process panes in a vertical layout.
type MainView struct {
	tApp         *tview.Application
	processPanes []*ProcessPane
	s            *StateStore

	selectedPaneIdx int
}

// NewMainView creates a new MainView with the given tview.Application and ProcessPanes.
func NewMainView(tApp *tview.Application, processrunners []processrunner.ProcessRunner, s *StateStore) *MainView {
	// create process panes for the runners
	var processPanes []*ProcessPane
	for _, pr := range processrunners {
		pp := NewProcessPane(tApp, pr)
		processPanes = append(processPanes, pp)
		// start scanning the stdout of the panes
		pp.StartScan()
		// set the defaults
		pp.SetIsWordWrap(s.GetIsWordWrap())
		pp.SetIsAutoScroll(s.GetIsAutoscroll())
		pp.SetIsBorderless(s.GetIsBorderless())
	}

	return &MainView{
		tApp:            tApp,
		processPanes:    processPanes,
		s:               s,
		selectedPaneIdx: 0,
	}
}

// Build the help text legend at the bottom of the main screen with dynamically
// changing setting status
func (mv *MainView) getHelpInfo() string {
	output := " "
	output += "(q)uit | "
	output += "(r)estart selected | "
	output += appendStatus("(a)utoscroll", mv.s.GetIsAutoscroll()) + " | "
	output += appendStatus("(w)rap lines", mv.s.GetIsWordWrap()) + " | "
	output += "(up/down) select pane | "
	output += "(enter) fullscreen selected pane"
	return output
}

// Render returns the tview.Flex that represents the MainView.
func (mv *MainView) Render(_ Props) *tview.Flex {
	innerFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	for _, pp := range mv.processPanes {
		// propagate the shared state to the process panes
		pp.SetIsAutoScroll(mv.s.GetIsAutoscroll())
		pp.SetIsBorderless(mv.s.GetIsBorderless())
		pp.SetIsWordWrap(mv.s.GetIsWordWrap())
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
					mv.s.ToggleAutoscroll()
					for _, pp := range mv.processPanes {
						pp.SetIsAutoScroll(mv.s.GetIsAutoscroll())
					}
				case 'e':
					mv.s.SetPreviousView("main", nil)
					mv.s.SetIsBorderless(false)
					a.SetView("environment", nil)
				case 'q':
					a.Exit()
					return nil
				// hotkey for restarting process
				case 'r':
					// get selected ProcessPane, restart its process, and start scanning again
					selectedPP := mv.processPanes[mv.selectedPaneIdx]
					err := selectedPP.pr.Restart()
					if err != nil {
						log.WithError(err).Error("error restarting process")
					}
					selectedPP.StartScan()
				case 'w':
					mv.s.ToggleWordWrap()
					for _, pp := range mv.processPanes {
						pp.SetIsWordWrap(mv.s.GetIsWordWrap())
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
		default:
			// do nothing. intentionally left blank
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
