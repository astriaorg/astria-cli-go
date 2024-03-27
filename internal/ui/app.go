package ui

import (
	"fmt"

	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/rivo/tview"
)

// AppController is an interface for the App to control the views and itself.
type AppController interface {
	// Start starts the app.
	Start()
	// Exit exits the app.
	Exit()
	// SetView sets the current view.
	SetView(view string, p Props)
}

// App contains the tview.Application and other necessary fields to manage the ui
type App struct {
	// Application is the tview application
	*tview.Application

	// view is the current view
	view View

	// viewMap is a map of views, so we can switch between them easily
	viewMap map[string]View

	// processRunners is a list of our running processes
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
	// create the views
	mainView := NewMainView(a.Application, a.processRunners)
	fullscreenView := NewFullscreenView(a.Application, nil)

	// set the views
	a.viewMap = map[string]View{
		"main":       mainView,
		"fullscreen": fullscreenView,
	}

	// show "main" view initially
	a.SetView("main", nil)

	// run the tview application
	if err := a.Application.Run(); err != nil {
		fmt.Println("error running tview application:", err)
		panic(err)
	}
}

// Exit stops all the process runners and stops the tview application.
func (a *App) Exit() {
	for _, pr := range a.processRunners {
		// FIXME - is there a cleaner way to stop the process runners?
		pr.Stop()
	}
	a.Application.Stop()
}

// SetView sets the view to the specified view.
func (a *App) SetView(view string, p Props) {
	a.view = a.viewMap[view]
	a.Application.SetInputCapture(nil)
	a.Application.SetInputCapture(a.view.GetKeyboard(a))
	a.Application.SetRoot(a.view.Render(p), true)
}
