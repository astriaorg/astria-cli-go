package ui

import (
	"fmt"

	"github.com/astriaorg/astria-cli-go/modules/cli/internal/processrunner"
	"github.com/rivo/tview"
)

// AppController is an interface for the App to control the views and itself.
type AppController interface {
	// Start starts the app.
	Start(*StateStore)
	// Exit exits the app.
	Exit()
	// SetView sets the current view.
	SetView(view string, p Props)
	// RefreshView resets keyboard input and sets the root tview component,
	// which effectively refreshes the view.
	RefreshView(p Props)
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
func (a *App) Start(stateStore *StateStore) {
	// if a state store wans't passed in, create a default one
	if stateStore == nil {
		stateStore = NewStateStore()
	}

	// create the views
	mainView := NewMainView(a.Application, a.processRunners, stateStore)
	fullscreenView := NewFullscreenView(a.Application, nil, stateStore)

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
		a.Exit()
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
func (a *App) SetView(view string, p Props) {
	a.view = a.viewMap[view]
	a.RefreshView(p)
}

// RefreshView refreshes the view by calling Render again and resetting the
// newly rendered view as the root.
func (a *App) RefreshView(p Props) {
	a.Application.SetInputCapture(nil)
	a.Application.SetMouseCapture(nil)
	a.Application.SetInputCapture(a.view.GetKeyboard(a))
	a.Application.SetRoot(a.view.Render(p), true)
}
