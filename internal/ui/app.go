package ui

import (
	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	BottomLegendText = " (q)uit | (w)ord wrap | up/down select pane | enter fullscreen"
	MainTitle        = "Astria Dev"
)

// App contains the tview.Application and other necessary fields to manage the ui
type App struct {
	// Application is the tview application
	*tview.Application

	// flex is the top level flex view
	flex *tview.Flex

	// processPanes holds ProcessPanes, one for each process.
	processPanes []*ProcessPane
	// selectedPane is the currently selected ProcessPane
	selectedPane *ProcessPane
}

// NewApp creates a new tview.Application with the necessary components
func NewApp(processrunners []*processrunner.ProcessRunner) *App {
	tviewApp := tview.NewApplication()
	var processPanes []*ProcessPane
	var selectedPane *ProcessPane

	// create ProcessPane for each process and add to innerFlex
	innerFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	for index, pr := range processrunners {
		pp := NewProcessPane(tviewApp, pr)
		processPanes = append(processPanes, pp)

		shouldFocus := false
		if index == 0 {
			// select and focus the first pane
			shouldFocus = true
			selectedPane = pp
		}
		innerFlex.AddItem(pp.textView, 0, 1, shouldFocus)
	}

	// create main flex view and add help text and innerFlex
	mainWindowHelpInfo := tview.NewTextView().SetText(BottomLegendText)
	flex := tview.NewFlex()
	flex.AddItem(innerFlex, 0, 4, false).
		SetDirection(tview.FlexRow).
		AddItem(mainWindowHelpInfo, 1, 0, false)
	flex.SetTitle(MainTitle).SetBorder(true)

	return &App{
		Application:  tviewApp,
		flex:         flex,
		processPanes: processPanes,
		selectedPane: selectedPane,
	}
}

// Start starts the tui application
func (a *App) Start() {
	// start scanning stdout for each process
	for _, pr := range a.processPanes {
		pr.StartScan()
	}

	// keyboard shortcuts
	a.Application.SetInputCapture(a.KeyboardMainView)

	// set the ui root primitive and run the tview application
	a.Application.SetRoot(a.flex, true)
	if err := a.Application.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	// stop each process
	for _, pp := range a.processPanes {
		pp.pr.Stop()
	}
	// stop the tview application
	a.Application.Stop()
}

func (a *App) KeyboardMainView(evt *tcell.EventKey) *tcell.EventKey {
	switch evt.Key() {
	case tcell.KeyCtrlC:
		a.Stop()
	case tcell.KeyRune:
		switch evt.Rune() {
		case 'a':
			// TODO - autoscroll
		case 'q':
			a.Stop()
		case 'w':
			for _, pp := range a.processPanes {
				pp.ToggleIsWordWrapped()
			}
		}
	case tcell.KeyUp:

	case tcell.KeyDown:
		// TODO - select pane
	case tcell.KeyEnter:
		// TODO - enter fullscreen
	}
	return evt
}
