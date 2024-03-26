package ui

import (
	"fmt"

	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/rivo/tview"
)

// App contains the tview.Application and other necessary fields to manage the ui
type App struct {
	// Application is the tview application
	*tview.Application

	// view is the current view
	view View

	// processRunners is a slice of processrunner.ProcessRunner
	processRunners []processrunner.ProcessRunner
}

// NewApp creates a new tview.Application with the necessary components
func NewApp(processrunners []processrunner.ProcessRunner) *App {
	tviewApp := tview.NewApplication()
	return &App{
		Application:    tviewApp,
		processRunners: processrunners,
	}
}

// Start starts the tview application.
func (a *App) Start() {
	// start with main view
	a.SetView("main", nil)

	// run the tview application
	if err := a.Application.Run(); err != nil {
		panic(err)
	}
}

// Exit stops all the process runners and stops the tview application.
func (a *App) Exit() {
	for _, pr := range a.processRunners {
		pr.Stop()
	}
	a.Application.Stop()
}

// SetView sets the view to the specified view.
func (a *App) SetView(view string, selectedPane *ProcessPane) {
	fmt.Println("setting view to", view)
	if view == "main" {
		a.view = NewMainView(a.Application, a.processRunners)
	}
	if view == "fullscreen" {
		a.view = NewFullscreenView(a.Application, selectedPane)
	}
	a.Application.SetInputCapture(nil)
	a.Application.SetInputCapture(a.view.GetKeyboard(*a))
	a.Application.SetRoot(a.view.Render(), true)
}
