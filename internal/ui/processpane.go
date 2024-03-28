package ui

import (
	"fmt"
	"io"

	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ProcessPane is a struct containing a tview.TextView and processrunner.ProcessRunner
type ProcessPane struct {
	tApp       *tview.Application
	textView   *tview.TextView
	pr         processrunner.ProcessRunner
	ansiWriter io.Writer

	Title     string
	lineCount int
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
		tApp:       tApp,
		textView:   tv,
		ansiWriter: ansiWriter,
		pr:         pr,
		Title:      pr.GetTitle(),
	}
}

// StartScan starts scanning the stdout of the process and writes to the textView
func (pp *ProcessPane) StartScan() {
	// scan stdout and write using ansiWriter
	go func() {
		stdoutScanner := pp.pr.GetScanner()
		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			pp.tApp.QueueUpdateDraw(func() {
				_, err := pp.ansiWriter.Write([]byte(line + "\n"))
				pp.lineCount++
				if err != nil {
					fmt.Println("error writing to textView:", err)
					panic(err)
				}
			})
		}
		if err := stdoutScanner.Err(); err != nil {
			fmt.Println("error reading stdout:", err)
			panic(err)
		}
		// FIXME - do i need to wait??
		if err := pp.pr.Wait(); err != nil {
			fmt.Println("error waiting for process:", err)
			panic(err)
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
	// NOTE - the verbage for isBorderless is opposite of SetBorder
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

func (pp *ProcessPane) GetLineCount() int {
	return pp.lineCount
}
