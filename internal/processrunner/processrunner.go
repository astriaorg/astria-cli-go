package processrunner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"
)

// ProcessRunner is an interface that represents a process to be run.
type ProcessRunner interface {
	Start(depStarted <-chan bool) error
	Wait() error
	Stop()
	GetDidStart() <-chan bool
	GetStdout() io.ReadCloser
	GetStderr() io.ReadCloser
	GetTitle() string
	GetScanner() *bufio.Scanner
}

// ProcessRunner is a struct that represents a process to be run.
type processRunner struct {
	// cmd is the exec.Cmd to be run
	cmd *exec.Cmd
	// Title is the title of the process
	title string

	didStart chan bool
	stdout   io.ReadCloser
	stderr   io.ReadCloser
	ctx      context.Context
	cancel   context.CancelFunc
	scanner  *bufio.Scanner
}

type NewProcessRunnerOpts struct {
	Title   string
	BinPath string
	Env     []string
	Args    []string
}

// NewProcessRunner creates a new ProcessRunner.
// It creates a new exec.Cmd with the given binPath and args, and sets the environment.
func NewProcessRunner(ctx context.Context, opts NewProcessRunnerOpts) ProcessRunner {
	ctx, cancel := context.WithCancel(ctx)

	cmd := exec.Command(opts.BinPath, opts.Args...)
	cmd.Env = opts.Env
	return &processRunner{
		cmd:      cmd,
		title:    opts.Title,
		didStart: make(chan bool),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start starts the process and returns the ProcessRunner and an error.
// It takes a channel that's closed when the dependency starts.
// This allows us to control the order of process startup.
func (pr *processRunner) Start(depStarted <-chan bool) error {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		select {
		// wait for the dependency to start
		case <-depStarted:
		case <-pr.ctx.Done():
			return
		}

		stdout, err := pr.cmd.StdoutPipe()
		if err != nil {
			fmt.Println("error obtaining stdout:", err)
			return
		}
		pr.stdout = stdout

		stderr, err := pr.cmd.StderrPipe()
		if err != nil {
			fmt.Println("error obtaining stderr:", err)
			return
		}
		pr.stderr = stderr

		err = pr.cmd.Start()
		if err != nil {
			fmt.Println("error starting process:", err)
			return
		}

		pr.scanner = bufio.NewScanner(pr.stdout)

		// signal that this process started
		close(pr.didStart)
	}()

	wg.Wait()

	if pr.ctx.Err() != nil {
		// the context was cancelled, return the context's error
		return pr.ctx.Err()
	}

	return nil
}

// Wait waits for the process to finish.
func (pr *processRunner) Wait() error {
	return pr.cmd.Wait()
}

// Stop stops the process.
func (pr *processRunner) Stop() {
	// send SIGINT to the process
	if err := pr.cmd.Process.Signal(syscall.SIGINT); err != nil {
		fmt.Println("Error sending SIGINT:", err)
	}
	// this will terminate the process if it's running
	pr.cancel()
}

// GetDidStart returns a channel that's closed when the process starts.
func (pr *processRunner) GetDidStart() <-chan bool {
	return pr.didStart
}

// GetStdout provides a reader for the process's stdout.
func (pr *processRunner) GetStdout() io.ReadCloser {
	return pr.stdout
}

// GetStderr provides a reader for the process's stderr.
func (pr *processRunner) GetStderr() io.ReadCloser {
	return pr.stderr
}

// GetTitle returns the title of the process.
func (pr *processRunner) GetTitle() string {
	return pr.title
}

// GetScanner returns a scanner for the process's stdout.
func (pr *processRunner) GetScanner() *bufio.Scanner {
	return pr.scanner
}
