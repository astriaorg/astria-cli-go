package processrunner

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessRunner(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := NewProcessRunnerOpts{
		Title:   "Test Echo",
		BinPath: "/bin/echo",
		Args:    []string{"hello, world"},
	}

	pr := NewProcessRunner(ctx, opts)
	assert.NotNil(t, pr, "ProcessRunner should not be nil")

	depStarted := make(chan bool)
	close(depStarted)
	err := pr.Start(ctx, depStarted)
	assert.Nil(t, err, "Process should start without error")

	// wait for the process to signal it has started
	select {
	case <-pr.GetDidStart():
		// expected path
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for process to start")
	}

	// give some time for the process to complete and write its output
	time.Sleep(1 * time.Second)

	pr.Stop()
	output := pr.GetOutput()
	expectedOutput := "hello, world\nprocess exited cleanly"
	assert.Contains(t, output, expectedOutput, "Output should contain the expected text")
}
