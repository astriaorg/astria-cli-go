package ui

import (
	"bufio"

	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/rivo/tview"
)

// ProcessPane is a struct containing a tview.TextView and processrunner.ProcessRunner
type ProcessPane struct {
	tApp     *tview.Application
	textView *tview.TextView
	pr       *processrunner.ProcessRunner
	title    string
}

// NewProcessPane creates a new ProcessPane with a textView and processrunner.ProcessRunner
func NewProcessPane(tApp *tview.Application, pr *processrunner.ProcessRunner) *ProcessPane {
	tv := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			tApp.Draw()
		})
	tv.SetTitle(pr.GetTitle()).SetBorder(true)

	return &ProcessPane{
		tApp:     tApp,
		textView: tv,
		pr:       pr,
		title:    pr.GetTitle(),
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
