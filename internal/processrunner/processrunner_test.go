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

func TestProcessRunner_StartError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := NewProcessRunnerOpts{
		Title:   "Nonexistent Command",
		BinPath: "/path/to/nonexistent",
	}

	pr := NewProcessRunner(ctx, opts)
	depStarted := make(chan bool)
	close(depStarted)

	err := pr.Start(ctx, depStarted)
	assert.NotNil(t, err, "Expected an error for a nonexistent command")
}

func TestProcessRunner_ImmediateExitWithError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// exits with an error code immediately
	opts := NewProcessRunnerOpts{
		Title:   "Fail Command",
		BinPath: "bash",
		Args:    []string{"-c", "exit 1"},
	}

	pr := NewProcessRunner(ctx, opts)
	depStarted := make(chan bool)
	close(depStarted)

	err := pr.Start(ctx, depStarted)
	assert.Nil(t, err, "Process should start without error even if it exits with an error")

	time.Sleep(1 * time.Second)

	output := pr.GetOutput()
	assert.Contains(t, output, "process exited with error", "Output should contain the error exit status")
}

func TestProcessRunner_LongRunningProcess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// sleep command to simulate a long-running operation
	opts := NewProcessRunnerOpts{
		Title:   "Sleep Command",
		BinPath: "sleep",
		Args:    []string{"1"}, // sleep for a second
	}

	pr := NewProcessRunner(ctx, opts)
	depStarted := make(chan bool)
	close(depStarted)

	err := pr.Start(ctx, depStarted)
	assert.Nil(t, err, "Process should start without error")

	// wait longer than the sleep command duration
	<-time.After(2 * time.Second)

	output := pr.GetOutput()
	assert.Equal(t, "process exited cleanly", output, "Expected clean exit after sleep")
}
