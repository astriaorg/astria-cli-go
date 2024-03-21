package ui

import (
	"bufio"

	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ProcessPane is a struct containing a tview.TextView and processrunner.ProcessRunner
type ProcessPane struct {
	tApp     *tview.Application
	textView *tview.TextView
	pr       *processrunner.ProcessRunner
	Title    string

	isAutoScroll bool
	isWordWrap   bool
}

// NewProcessPane creates a new ProcessPane with a textView and processrunner.ProcessRunner
func NewProcessPane(tApp *tview.Application, pr *processrunner.ProcessRunner) *ProcessPane {
	tv := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			tApp.Draw()
		})
	tv.SetBorder(true).
		SetBorderColor(tcell.ColorGray).
		SetTitle(pr.GetTitle())

	return &ProcessPane{
		tApp:     tApp,
		textView: tv,
		pr:       pr,
		Title:    pr.GetTitle(),

		isAutoScroll: true,
		isWordWrap:   false,
	}
}

// StartScan starts scanning the stdout of the process and writes to the textView
func (pp *ProcessPane) StartScan() {
	// scan stdout and write using ansiWriter
	go func() {
		// ansi writer
		ansiWriter := tview.ANSIWriter(pp.textView)

		// new scanner to scan stdout
		stdoutScanner := bufio.NewScanner(pp.pr.GetStdout())
		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			pp.tApp.QueueUpdateDraw(func() {
				_, err := ansiWriter.Write([]byte(line + "\n"))
				if err != nil {
					panic(err)
				}
			})
		}
		if err := stdoutScanner.Err(); err != nil {
			panic(err)
		}
		if err := pp.pr.Wait(); err != nil {
			panic(err)
		}
	}()
}

// ToggleIsWordWrapped sets the word wrap of the textView.
func (pp *ProcessPane) ToggleIsWordWrapped() {
	pp.isWordWrap = !pp.isWordWrap
	// set the textview's word wrap
	pp.textView.SetWrap(pp.isWordWrap)
}

// GetTextView returns the textView associated with the ProcessPane.
func (pp *ProcessPane) GetTextView() *tview.TextView {
	return pp.textView
}

// Highlight highlights or unhighlights the ProcessPane's textView.
func (pp *ProcessPane) Highlight(highlight bool) {
	if highlight {
		title := "[black:lightblue]" + pp.Title + "[::-]"
		pp.textView.SetBorderColor(tcell.ColorLightBlue).SetTitle(title)
	} else {
		pp.textView.SetBorderColor(tcell.ColorGray).SetTitle(pp.Title)
	}
}