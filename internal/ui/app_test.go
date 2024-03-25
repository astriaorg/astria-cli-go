package ui

import (
	"testing"

	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/astria/astria-cli-go/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	mockRunner1 := new(testutils.MockProcessRunner)
	mockRunner1.On("GetTitle").Return("Process 1")
	mockRunner2 := new(testutils.MockProcessRunner)
	mockRunner2.On("GetTitle").Return("Process 2")
	runners := []processrunner.ProcessRunner{mockRunner1, mockRunner2}
	app := NewApp(runners)

	assert.Equal(t, 2, len(app.processPanes), "expected 2 process panes")
	assert.NotEmptyf(t, app.Application, "expected app.Application to be not empty")
	assert.Equal(t, MainTitle, app.flex.GetTitle(), "expected title to be MainTitle")
	assert.True(t, app.isAutoScroll, "expected isAutoScroll to be true")
	assert.False(t, app.isWordWrap, "expected isWordWrap to be false")
}

func TestToggleAutoScrollSyncsWithPanes(t *testing.T) {
	mockRunner1 := new(testutils.MockProcessRunner)
	mockRunner1.On("GetTitle").Return("Process 1")
	mockRunner2 := new(testutils.MockProcessRunner)
	mockRunner2.On("GetTitle").Return("Process 2")
	runners := []processrunner.ProcessRunner{mockRunner1, mockRunner2}
	app := NewApp(runners)

	initialState := app.isAutoScroll
	app.toggleAutoScroll()

	assert.NotEqual(t, initialState, app.isAutoScroll, "Expected isAutoScroll to toggle, but it did not")
	for _, pp := range app.processPanes {
		assert.Equal(t, app.isAutoScroll, pp.isAutoScroll, "Expected ProcessPane isAutoScroll to match app isAutoScroll")
	}
}

func TestToggleWordWrapSyncsWithPanes(t *testing.T) {
	mockRunner1 := new(testutils.MockProcessRunner)
	mockRunner1.On("GetTitle").Return("Process 1")
	mockRunner2 := new(testutils.MockProcessRunner)
	mockRunner2.On("GetTitle").Return("Process 2")
	runners := []processrunner.ProcessRunner{mockRunner1, mockRunner2}
	app := NewApp(runners)

	initialState := app.isWordWrap
	app.toggleWordWrap()

	assert.NotEqual(t, initialState, app.isWordWrap, "Expected isWordWrap to toggle, but it did not")
	for _, pp := range app.processPanes {
		assert.Equal(t, app.isWordWrap, pp.isWordWrap, "Expected ProcessPane isWordWrap to match app isWordWrap")
	}
}

//func TestKeyPressToggles(t *testing.T) {
//	mockRunner1 := new(testutils.MockProcessRunner)
//	mockRunner1.On("GetTitle").Return("Process 1")
//	mockRunner1.On("GetStdout").Return(mock.AnythingOfType("*io.ReadCloser"))
//	mockRunner1.On("GetStderr").Return(mock.AnythingOfType("*io.ReadCloser"))
//	mockRunner2 := new(testutils.MockProcessRunner)
//	mockRunner2.On("GetTitle").Return("Process 2")
//	mockRunner2.On("GetStdout").Return(mock.AnythingOfType("*io.ReadCloser"))
//	mockRunner2.On("GetStderr").Return(mock.AnythingOfType("*io.ReadCloser"))
//	runners := []processrunner.ProcessRunner{mockRunner1, mockRunner2}
//	app := NewApp(runners)
//	// create a SimulationScreen so we can inject keys
//	simScreen := tcell.NewSimulationScreen("UTF-8")
//	simScreen.Init()
//	simScreen.SetSize(512, 512)
//	app.Application.SetScreen(simScreen)
//	app.Start()
//
//	initialAutoscroll := app.isAutoScroll
//	fmt.Println("initialAutoscroll", app.isAutoScroll)
//	simScreen.InjectKey(tcell.KeyRune, 'a', tcell.ModNone)
//	fmt.Println("after", app.isAutoScroll)
//
//	assert.NotEqualf(t, initialAutoscroll, app.isAutoScroll, "Expected isAutoScroll to toggle, but it did not")
//	// check that the panes are in sync
//	for _, pp := range app.processPanes {
//		assert.Equal(t, app.isAutoScroll, pp.isAutoScroll, "Expected ProcessPane isAutoScroll to match app isAutoScroll")
//	}
//}
