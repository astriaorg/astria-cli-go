package ui

import (
	"bufio"

	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/rivo/tview"
)

const (
	BottomLegendText = " (q)uit | (w)ord wrap | up/down select pane | enter fullscreen"
	MainTitle        = "Astria Dev"
)

// App is a struct for managing the ui.
type App struct {
	*tview.Application
	flexView     *tview.Flex
	ProcessPanes []*ProcessPane
	SelectedPane *ProcessPane
}

// NewApp creates a new tview.Application with the necessary components
func NewApp(processrunners []*processrunner.ProcessRunner) *App {
	tApp := tview.NewApplication()
	var processPanes []*ProcessPane
	var selectedPane *ProcessPane

	// create ProcessPane for each process and add to innerFlex
	innerFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	for index, pr := range processrunners {
		pp := NewProcessPane(tApp, pr)
		processPanes = append(processPanes, pp)

		shouldFocus := false
		if index == 0 {
			shouldFocus = true
			selectedPane = pp
		}
		innerFlex.AddItem(pp.textView, 0, 1, shouldFocus)
	}

	// create main flex view and add help text and innerFlex
	mainWindowHelpInfo := tview.NewTextView().SetText(BottomLegendText)
	flexView := tview.NewFlex()
	flexView.AddItem(innerFlex, 0, 4, false).
		SetDirection(tview.FlexRow).
		AddItem(mainWindowHelpInfo, 1, 0, false)
	flexView.SetTitle(MainTitle).SetBorder(true)

	return &App{
		Application:  tApp,
		flexView:     flexView,
		ProcessPanes: processPanes,
		SelectedPane: selectedPane,
	}
}

// Start starts the application
func (a *App) Start() {
	// start scanning stdout for each process
	for _, pr := range a.ProcessPanes {
		pr.StartScan()
	}
	// set the ui root primitive and run the tview application
	a.Application.SetRoot(a.flexView, true)
	if err := a.Application.Run(); err != nil {
		panic(err)
	}
}

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
	tv.SetTitle(pr.Title).SetBorder(true)

	return &ProcessPane{
		tApp:     tApp,
		textView: tv,
		pr:       pr,
		title:    pr.Title,
	}
}

// StartScan starts scanning the stdout of the process and writes to the textView
func (pp *ProcessPane) StartScan() {
	// ansi writer
	ansiWriter := tview.ANSIWriter(pp.textView)

	// new scanner to scan stdout
	stdoutScanner := bufio.NewScanner(pp.pr.Stdout)

	// scan stdout and write using ansiWriter
	go func() {
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
		if err := pp.pr.Cmd.Wait(); err != nil {
			panic(err)
		}
	}()
}
