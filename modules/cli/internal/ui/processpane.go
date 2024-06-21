package ui

import (
	"io"
	"strings"
	"sync/atomic"
	"time"

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

	Title string
}

// NewProcessPane creates a new ProcessPane with a textView and processrunner.ProcessRunner
func NewProcessPane(tApp *tview.Application, pr processrunner.ProcessRunner) *ProcessPane {
	tv := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			tApp.Draw()
		})
	tv.SetBorder(true).
		SetBorderColor(tcell.ColorGray).
		SetTitle(pr.GetTitle())

	ansiWriter := tview.ANSIWriter(tv)

	return &ProcessPane{
		tApp:           tApp,
		textView:       tv,
		ansiWriter:     ansiWriter,
		pr:             pr,
		Title:          pr.GetTitle(),
		TickerInterval: 250,
	}
}

// StartScan starts scanning the stdout of the process and writes to the textView
func (pp *ProcessPane) StartScan() {
	go func() {
		// initialize a ticker for periodic updates
		ticker := time.NewTicker(pp.TickerInterval * time.Millisecond) // adjust the duration as needed
		defer ticker.Stop()

		for range ticker.C {
			currentOutput := pp.pr.GetOutputAndClearBuf() // get the current full output

			// new, unprocessed data.
			pp.tApp.QueueUpdateDraw(func() {
				_, err := pp.ansiWriter.Write([]byte(currentOutput))
				if err != nil {
					log.WithError(err).Error("Error writing to textView")
				}
				atomic.AddInt64(&pp.lineCount, int64(strings.Count(currentOutput, "\n")))
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

// GetTextView returns the textView associated with the ProcessPane.
func (pp *ProcessPane) GetTextView() *tview.TextView {
	return pp.textView
}

// Highlight highlights or unhighlights the ProcessPane's textView.
func (pp *ProcessPane) Highlight(highlight bool) {
	if highlight {
		title := "[black:blue]" + pp.Title + "[::-]"
		pp.textView.SetBorderColor(tcell.ColorBlue).SetTitle(title)
	} else {
		pp.textView.SetBorderColor(tcell.ColorGray).SetTitle(pp.Title)
	}
}

// GetLineCount returns the line count of the ProcessPane's textView.
func (pp *ProcessPane) GetLineCount() int64 {
	return pp.lineCount
}
