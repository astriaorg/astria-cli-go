package ui

import (
	"testing"

	"github.com/astriaorg/astria-cli-go/modules/cli/internal/processrunner"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	mockRunner1 := new(testutils.MockProcessRunner)
	mockRunner1.On("GetTitle").Return("Process 1")
	mockRunner2 := new(testutils.MockProcessRunner)
	mockRunner2.On("GetTitle").Return("Process 2")
	runners := []processrunner.ProcessRunner{mockRunner1, mockRunner2}
	app := NewApp(runners)

	assert.NotNil(t, app)
}
