package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

// FullscreenView represents the fullscreen view when a pane is selected.
type EnvironmentView struct {
	tApp          *tview.Application
	configuration []byte
	s             *StateStore
}

// NewFullscreenView creates a new FullscreenView with the given tview.Application and ProcessPane.
func NewEnvironmentView(tApp *tview.Application, configuration []byte, s *StateStore) *EnvironmentView {
	// TODO - add some syntax highlighting to the environment view
	return &EnvironmentView{
		tApp:          tApp,
		configuration: configuration,
		s:             s,
	}
}

// Build the help text legend at the bottom of the fullscreen view with dynamically
// changing setting status
func (ev *EnvironmentView) getHelpInfo() string {
	output := " "
	output += "(q/esc/e) back | "
	output += appendStatus("(w)rap lines", ev.s.GetIsWordWrap()) + " | "
	output += appendStatus("(b)orderless", ev.s.GetIsBorderless()) + " | "
	output += "(up/down or mousewheel) scroll if autoscroll is off"
	return output
}

// Render returns the tview.Flex that represents the FullscreenView.
func (ev *EnvironmentView) Render(p Props) *tview.Flex {
	// fv.processPane = p.(*ProcessPane)
	// build tview text views and flex
	help := tview.NewTextView().
		SetDynamicColors(true).
		SetText(ev.getHelpInfo()).
		SetChangedFunc(func() {
			ev.tApp.Draw()
		})
	flex := tview.NewFlex().AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ev.processPane.GetTextView(), 0, 1, true).
		AddItem(help, 1, 0, false), 0, 4, false)
	return flex
}

// GetKeyboard is a callback for defining keyboard shortcuts.
func (fv *EnvironmentView) GetKeyboard(a AppController) func(evt *tcell.EventKey) *tcell.EventKey {
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
