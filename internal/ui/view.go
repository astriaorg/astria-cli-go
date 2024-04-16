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

// Props is an empty interface for passing data to the view.
type Props interface{}

type View interface {
	// Render returns the tview.Flex that represents the view.
	// The Props argument is used to pass data to the view.
	Render(p Props) *tview.Flex
	// GetKeyboard is a callback for defining keyboard shortcuts
	GetKeyboard(a AppController) func(evt *tcell.EventKey) *tcell.EventKey
}

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

// Append the settings status to the end of the input string
func appendStatus(text string, status bool) string {
	if status {
		return text + ": [black:white]ON [-:-]"
	} else {
		return text + ": [white:darkslategray]off[-:-]"
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

// FullscreenView represents the fullscreen view when a pane is selected.
type FullscreenView struct {
	tApp        *tview.Application
	processPane *ProcessPane
	s           *StateStore
}

// NewFullscreenView creates a new FullscreenView with the given tview.Application and ProcessPane.
func NewFullscreenView(tApp *tview.Application, processPane *ProcessPane, s *StateStore) *FullscreenView {
	return &FullscreenView{
		tApp:        tApp,
		processPane: processPane,
		s:           s,
	}
}

// Build the help text legend at the bottom of the fullscreen view with dynamically
// changing setting status
func (fv *FullscreenView) getHelpInfo() string {
	output := " "
	output += "(q/esc) back | "
	output += "(r)estart | "
	output += appendStatus("(a)utoscroll", fv.s.GetIsAutoscroll()) + " | "
	output += appendStatus("(w)rap lines", fv.s.GetIsWordWrap()) + " | "
	output += appendStatus("(b)orderless", fv.s.GetIsBorderless()) + " | "
	output += "(0/1) jump to head/tail" + " | "
	output += "(up/down or mousewheel) scroll if autoscroll is off"
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
		fv.s.ResetBorderless()
		fv.processPane.SetIsBorderless(fv.s.GetIsBorderless())
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
				// hotkey for autoscroll
				case 'a':
					fv.s.ToggleAutoscroll()
					fv.processPane.SetIsAutoScroll(fv.s.GetIsAutoscroll())
				// hotkey for borderless
				case 'b':
					fv.s.ToggleBorderless()
					fv.processPane.SetIsBorderless(fv.s.GetIsBorderless())
				// hotkey for quitting fullscreen
				case 'q':
					backToMain()
					return nil
				// hotkey for restarting process
				case 'r':
					err := fv.processPane.pr.Restart()
					if err != nil {
						log.WithError(err).Error("error restarting process")
					}
					fv.processPane.StartScan()
				// hotkey for word wrap
				case 'w':
					fv.s.ToggleWordWrap()
					fv.processPane.SetIsWordWrap(fv.s.GetIsWordWrap())
				// hotkey for jumping to the head of the logs
				case '0':
					fv.s.DisableAutoscroll()
					fv.processPane.textView.ScrollToBeginning()
				// hotkey for jumping to the tail of the logs
				case '1':
					fv.s.DisableAutoscroll()
					fv.processPane.textView.ScrollTo(int(fv.processPane.GetLineCount()), 0)
				}
				// needed to call the Render method again to refresh the help info
				a.RefreshView(fv.processPane)
				return nil
			}
		case tcell.KeyEscape:
			backToMain()
			return nil
		// only listen to up and down keys if autoscroll is off
		case tcell.KeyUp:
			if !fv.s.GetIsAutoscroll() {
				row, _ := fv.processPane.textView.GetScrollOffset()
				fv.processPane.textView.ScrollTo(row-1, 0)
				return nil
			}
		case tcell.KeyDown:
			if !fv.s.GetIsAutoscroll() {
				row, _ := fv.processPane.textView.GetScrollOffset()
				fv.processPane.textView.ScrollTo(row+1, 0)
				return nil
			}
		default:
			// do nothing. intentionally left blank
		}
		return evt
	}
}
