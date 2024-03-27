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
	// Refresh the view to call Render again
	RefreshView(p Props)
	// Autoscroll toggle
	ToggleAutoscroll()
	// WordWrap toggle
	ToggleWordWrap()
	// Borderless toggle
	ToggleBorderless()
	// reset the border to true
	ResetBorderless()
}

type SharedState struct {
	isAutoScroll bool
	isWordWrap   bool
	isBorderless bool
}

func defaultSharedState() *SharedState {
	return &SharedState{
		isAutoScroll: true,
		isWordWrap:   false,
		isBorderless: false,
	}
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

	// app level store for ui settings state
	ss *SharedState
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
	// set the shared state defaults for the app
	a.ss = defaultSharedState()

	// create the views
	mainView := NewMainView(a.Application, a.processRunners, a)
	fullscreenView := NewFullscreenView(a.Application, nil, a)

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

func (a *App) RefreshView(p Props) {
	a.Application.SetRoot(a.view.Render(p), true)
}

// SetAutoscroll sets the autoscroll state.
func (a *App) ToggleAutoscroll() {
	a.ss.isAutoScroll = !a.ss.isAutoScroll
}

// SetWordWrap sets the word wrap state.
func (a *App) ToggleWordWrap() {
	a.ss.isWordWrap = !a.ss.isWordWrap
}

// SetBorderless sets the borderless state.
func (a *App) ToggleBorderless() {
	a.ss.isBorderless = !a.ss.isBorderless
}

// SetBorderless sets the borderless state.
func (a *App) ResetBorderless() {
	a.ss.isBorderless = false
}
