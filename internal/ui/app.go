package ui

import (
	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/rivo/tview"
)

// App contains the tview.Application and other necessary fields to manage the ui
type App struct {
	// Application is the tview application
	*tview.Application

	// view is the current view
	view View

	// NOTE - keeping track of process panes at global level made many things easier for now
	processPanes []*ProcessPane
}

// NewApp creates a new tview.Application with the necessary components
func NewApp(processrunners []processrunner.ProcessRunner) *App {
	tviewApp := tview.NewApplication()

	// create a ProcessPane for each process runner
	var processPanes []*ProcessPane
	for _, pr := range processrunners {
		pp := NewProcessPane(tviewApp, pr)
		processPanes = append(processPanes, pp)
	}

	return &App{
		Application:  tviewApp,
		processPanes: processPanes,
	}
}

// Start starts the tview application.
func (a *App) Start() {
	// show "main" view initially
	a.SetView("main", nil)

	// start scanning the stdout of all the process panes
	for _, pp := range a.processPanes {
		pp.StartScan()
	}

	// run the tview application
	if err := a.Application.Run(); err != nil {
		panic(err)
	}
}

// Exit stops all the process runners and stops the tview application.
func (a *App) Exit() {
	for _, pp := range a.processPanes {
		pp.StopProcess()
	}
	a.Application.Stop()
}

// SetView sets the view to the specified view.
func (a *App) SetView(view string, selectedPane *ProcessPane) {
	if view == "main" {
		a.view = NewMainView(a.Application, a.processPanes)
	}
	if view == "fullscreen" {
		a.view = NewFullscreenView(a.Application, selectedPane)
	}
	a.Application.SetInputCapture(nil)
	a.Application.SetInputCapture(a.view.GetKeyboard(*a))
	a.Application.SetRoot(a.view.Render(), true)
}
