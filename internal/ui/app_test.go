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
}
