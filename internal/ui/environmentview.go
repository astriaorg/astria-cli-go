package ui

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

// FullscreenView represents the fullscreen view when a pane is selected.
type EnvironmentView struct {
	tApp *tview.Application
	s    *StateStore

	textView *tview.TextView

	previousView string
	lineCount    int64
}

// NewFullscreenView creates a new FullscreenView with the given tview.Application and ProcessPane.
func NewEnvironmentView(tApp *tview.Application, processrunners []processrunner.ProcessRunner, s *StateStore) *EnvironmentView {
	if len(processrunners) == 0 {
		log.Error("no process runners provided to environment view")
		return nil
	}

	// the primary text view for the environment view does not change like it does in
	// other views, so we only need to set it up here instead of in the Render method
	tv := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			tApp.Draw()
		})
	tv.SetBorder(true).
		SetBorderColor(tcell.ColorGray).
		SetTitle(" Environment ")
	ansiWriter := tview.ANSIWriter(tv)

	// formate the binary names and their paths
	lineCount := int64(0)
	output := ""
	longestTitle := 0
	for _, pr := range processrunners {
		if len(pr.GetTitle()) > longestTitle {
			longestTitle = len(pr.GetTitle())
		}
	}
	for _, pr := range processrunners {
		output += fmt.Sprintf("%-*s", longestTitle+2, pr.GetTitle()+":") + pr.GetBinPath() + "\n"
		lineCount++
	}
	output += "\n"
	lineCount++
	ansiWriter.Write([]byte(output))

	// load the text for the environment view
	envPath := processrunners[0].GetEnvironmentPath()
	sourceCode, err := os.ReadFile(envPath)
	if err != nil {
		panic(err)
	}
	content := string(sourceCode)
	lines := strings.Split(content, "\n")

	// Filter out empty lines and comments
	var filteredLines []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "#") {
			filteredLines = append(filteredLines, trimmedLine)
		}
	}
	// remove duplicate lines
	seen := make(map[string]bool)
	var unique []string
	for _, entry := range filteredLines {
		if _, found := seen[entry]; !found {
			seen[entry] = true
			unique = append(unique, entry)
		}
	}
	sort.Strings(unique)
	envForFormatting := strings.Join(unique, "\n")

	// TODO - add configuration for the lexer and style?
	// there isn't a lexer for .env files in chroma, the python lexer happends
	// to work well which is why it is used here
	lexer := lexers.Get("python")
	if lexer == nil {
		lexer = lexers.Fallback
	}
	style := styles.Get("monokai")
	iterator, err := lexer.Tokenise(nil, envForFormatting)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	formatter := formatters.TTY256 // Assuming we're outputting to a terminal that supports 256 colors
	if err := formatter.Format(&buf, style, iterator); err != nil {
		panic(err)
	}

	ansiWriter.Write(buf.Bytes())

	return &EnvironmentView{
		tApp:      tApp,
		textView:  tv,
		s:         s,
		lineCount: lineCount,
	}
}

// Build the help text legend at the bottom of the fullscreen view with dynamically
// changing setting status
func (ev *EnvironmentView) getHelpInfo() string {
	output := " "
	output += "(q/esc/e) back | "
	output += appendStatus("(w)rap lines", ev.s.GetIsWordWrap()) + " | "
	output += appendStatus("(b)orderless", ev.s.GetIsBorderless()) + " | "
	output += "(0/1) jump to head/tail" + " | "
	output += "(up/down or mousewheel) scroll if autoscroll is off"
	return output
}

// Render returns the tview.Flex that represents the FullscreenView.
func (ev *EnvironmentView) Render(_ Props) *tview.Flex {
	// build tview text views and flex
	help := tview.NewTextView().
		SetDynamicColors(true).
		SetText(ev.getHelpInfo()).
		SetChangedFunc(func() {
			ev.tApp.Draw()
		})
	// update the shared state for the evnironment view
	ev.textView.SetBorder(!ev.s.GetIsBorderless())
	flex := tview.NewFlex().AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ev.textView, 0, 1, true).
		AddItem(help, 1, 0, false), 0, 4, false)
	return flex
}

// GetKeyboard is a callback for defining keyboard shortcuts.
func (ev *EnvironmentView) GetKeyboard(a AppController) func(evt *tcell.EventKey) *tcell.EventKey {
	backToPreviousView := func() {
		ev.s.ResetBorderless()
		ev.textView.SetBorder(ev.s.GetIsBorderless())
		a.RefreshView(ev.textView)
		prevView, prevProps := ev.s.GetPreviousView()
		if prevProps != nil {
			a.RefreshView(prevProps)
		}
		a.SetView(prevView, prevProps)
	}
	return func(evt *tcell.EventKey) *tcell.EventKey {
		switch evt.Key() {
		case tcell.KeyCtrlC:
			a.Exit()
			return nil
		case tcell.KeyRune:
			{
				switch evt.Rune() {
				case 'b':
					ev.s.ToggleBorderless()
					ev.textView.SetBorder(ev.s.GetIsBorderless())
				// hotkeys for returning to previous view
				case 'e', 'q':
					backToPreviousView()
					return nil
				// hotkey for word wrap
				case 'w':
					ev.s.ToggleWordWrap()
					ev.textView.SetWrap(ev.s.GetIsWordWrap())
				// hotkey for jumping to the head of the logs
				case '0':
					ev.s.DisableAutoscroll()
					ev.textView.ScrollToBeginning()
				// hotkey for jumping to the tail of the logs
				case '1':
					ev.s.DisableAutoscroll()
					ev.textView.ScrollTo(int(ev.GetLineCount()), 0)
				}
				// needed to call the Render method again to refresh the help info
				a.RefreshView(nil)
				return nil
			}
		case tcell.KeyEscape:
			backToPreviousView()
			return nil
		case tcell.KeyUp:
			row, _ := ev.textView.GetScrollOffset()
			ev.textView.ScrollTo(row-1, 0)
			return nil
		case tcell.KeyDown:
			row, _ := ev.textView.GetScrollOffset()
			ev.textView.ScrollTo(row+1, 0)
			return nil
		default:
			// do nothing. intentionally left blank
		}
		return evt
	}
}

func (ev *EnvironmentView) GetLineCount() int64 {
	return ev.lineCount
}
