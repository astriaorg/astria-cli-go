package ui

import (
	"io"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/config"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/processrunner"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

// ProcessPane is a struct containing a tview.TextView and processrunner.ProcessRunner
type ProcessPane struct {
	tApp           *tview.Application
	textView       *tview.TextView
	lineCount      int64
	pr             processrunner.ProcessRunner
	ansiWriter     io.Writer
	TickerInterval time.Duration

	Title       string
	isMinimized bool

	HighlightColor tcell.Color
	BorderColor    tcell.Color
}

// NewProcessPane creates a new ProcessPane with a textView and processrunner.ProcessRunner
func NewProcessPane(tApp *tview.Application, pr processrunner.ProcessRunner) *ProcessPane {
	tv := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetMaxLines(pr.GetMaxUiLogLines()).
		SetChangedFunc(func() {
			tApp.Draw()
		})
	tv.SetBorder(true).
		SetTitle(pr.GetTitle())

	ansiWriter := tview.ANSIWriter(tv)

	highlightColor := tcell.GetColor(pr.GetHighlightColor())
	if highlightColor == 0 {
		log.Debugf("Highlight color %s could not be parsed, using default %s", pr.GetHighlightColor(), tcell.ColorBlue.Name())
		highlightColor = tcell.ColorBlue
	}

	borderColor := tcell.GetColor(pr.GetBorderColor())
	if borderColor == 0 {
		log.Debugf("Border color %s could not be parsed, using default %s", pr.GetBorderColor(), tcell.ColorGray.Name())
		borderColor = tcell.ColorGray
	}

	return &ProcessPane{
		tApp:           tApp,
		textView:       tv,
		ansiWriter:     ansiWriter,
		pr:             pr,
		Title:          pr.GetTitle(),
		TickerInterval: 250,
		isMinimized:    pr.GetStartMinimized(),
		HighlightColor: highlightColor,
		BorderColor:    borderColor,
	}
}

// StartScan starts scanning the stdout of the process and writes to the textView
func (pp *ProcessPane) StartScan() {
	go func() {
		// initialize a ticker for periodic updates
		ticker := time.NewTicker(pp.TickerInterval * time.Millisecond) // adjust the duration as needed
		defer ticker.Stop()

		filter := pp.pr.GetLogFilter()
		if filter == "" {
			log.Warnf("No log filter found, using default: \"%s\"", config.DefaultLogFilter)
			filter = config.DefaultLogFilter
		}

		re := regexp.MustCompile(pp.pr.GetLogFilter())

		for range ticker.C {
			currentOutput := pp.pr.GetOutputAndClearBuf() // get the current full output

			// get individual lines from the output
			currentOutputLines := strings.Split(currentOutput, "\n")
			var keptLines []string
			// filter out lines that don't match the include filter
			for _, line := range currentOutputLines {
				if re.MatchString(line) {
					keptLines = append(keptLines, line)
				}
			}

			var builder strings.Builder
			for i, s := range keptLines {
				builder.WriteString(s)
				if i < len(keptLines)-1 {
					builder.WriteString("\n")
				}
			}
			result := builder.String()

			// new, unprocessed data.
			pp.tApp.QueueUpdateDraw(func() {
				// write output data to logs if possible
				if pp.pr.CanWriteToLog() {
					err := pp.pr.WriteToLog(result)
					if err != nil {
						log.WithError(err).Error("Error writing to log")
					}
				}
				// write output data to ui element
				_, err := pp.ansiWriter.Write([]byte(result))
				if err != nil {
					log.WithError(err).Error("Error writing to textView")
				}
				atomic.AddInt64(&pp.lineCount, int64(strings.Count(result, "\n")))
			})
		}
	}()
}

// SetIsAutoScroll sets the auto scroll of the textView.
func (pp *ProcessPane) SetIsAutoScroll(isAutoScroll bool) {
	if isAutoScroll {
		pp.textView.ScrollToEnd()
	} else {
		currentOffset, _ := pp.textView.GetScrollOffset()
		pp.textView.ScrollTo(currentOffset, 0)
	}
}

// SetIsWordWrap sets the word wrap of the textView.
func (pp *ProcessPane) SetIsWordWrap(isWordWrap bool) {
	// set the textview's word wrap
	pp.textView.SetWrap(isWordWrap)
}

// SetIsBorderless sets the border of the textView.
func (pp *ProcessPane) SetIsBorderless(isBorderless bool) {
	// NOTE - the verbiage for isBorderless is opposite of SetBorder
	// therefore, when isBorderless is true, we want to set the border to false
	// for the textView, and vice versa
	pp.textView.SetBorder(!isBorderless)
}

// TODO: description
func (pp *ProcessPane) SetMinimized(isMinimized bool) {
	pp.isMinimized = isMinimized
}

// GetTextView returns the textView associated with the ProcessPane.
func (pp *ProcessPane) GetTextView() *tview.TextView {
	return pp.textView
}

// Highlight highlights or unhighlights the ProcessPane's textView.
func (pp *ProcessPane) Highlight(highlight bool) {
	highlightTitleFormat := "[black:" + pp.HighlightColor.Name() + "]"
	if highlight {
		title := highlightTitleFormat + pp.Title + "[::-]"
		pp.textView.SetBorderColor(pp.HighlightColor).SetTitle(title)
	} else {
		pp.textView.SetBorderColor(pp.BorderColor).SetTitle(pp.Title)
	}
}

// GetLineCount returns the line count of the ProcessPane's textView.
func (pp *ProcessPane) GetLineCount() int64 {
	return pp.lineCount
}
