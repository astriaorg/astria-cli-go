package processrunner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// ProcessRunner is an interface that represents a process to be run.
type ProcessRunner interface {
	Start(depStarted <-chan bool) error
	Stop()
	GetDidStart() <-chan bool
	GetTitle() string
	GetOutput() string
	SetExitStatusString(status string)
	GetExitStatusString() string
}

// ProcessRunner is a struct that represents a process to be run.
type processRunner struct {
	// cmd is the exec.Cmd to be run
	cmd *exec.Cmd
	// Title is the title of the process
	title string

	exitStatusString     string
	exitStatusStringLock *sync.Mutex

	didStart  chan bool
	stdout    io.ReadCloser
	stderr    io.ReadCloser
	ctx       context.Context
	cancel    context.CancelFunc
	outputBuf *bytes.Buffer
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

	cmd := exec.CommandContext(ctx, opts.BinPath, opts.Args...)
	cmd.Env = opts.Env
	return &processRunner{
		cmd:                  cmd,
		title:                opts.Title,
		didStart:             make(chan bool),
		outputBuf:            new(bytes.Buffer),
		ctx:                  ctx,
		cancel:               cancel,
		exitStatusStringLock: &sync.Mutex{},
	}
}

// Start starts the process and returns the ProcessRunner and an error.
// It takes a channel that's closed when the dependency starts.
// This allows us to control the order of process startup.
func (pr *processRunner) Start(depStarted <-chan bool) error {
	select {
	case <-depStarted:
	// continue if the dependency has started.
	case <-pr.ctx.Done():
		log.Info("Context cancelled before starting process", pr.title)
		return pr.ctx.Err()
	}

	// get stdout and stderr
	stdout, err := pr.cmd.StdoutPipe()
	if err != nil {
		log.WithError(err).Errorf("Error obtaining stdout for process %s", pr.title)
		return err
	}
	pr.stdout = stdout

	stderr, err := pr.cmd.StderrPipe()
	if err != nil {
		log.WithError(err).Errorf("Error obtaining stderr for process %s", pr.title)
		return err
	}
	pr.stderr = stderr

	// multiwriter to write both stdout and stderr to the same buffer
	mw := io.MultiWriter(pr.outputBuf)
	go io.Copy(mw, stdout)
	go io.Copy(mw, stderr)

	// actually start the process
	if err := pr.cmd.Start(); err != nil {
		log.WithError(err).Errorf("error starting process %s", pr.title)
		return err
	}

	// signal that this process has started.
	close(pr.didStart)

	// asynchronously monitor process
	go func() {
		err = pr.cmd.Wait()
		if err != nil {
			err = fmt.Errorf("process exited with error: %w", err)
			log.Error(err)
			pr.SetExitStatusString(err.Error())
		} else {
			s := fmt.Sprintf("process exited cleanly")
			log.Infof(s)
			pr.SetExitStatusString(s)
		}
	}()

	return nil
}

func (pr *processRunner) SetExitStatusString(status string) {
	if status != "" {
		pr.exitStatusStringLock.Lock()
		defer pr.exitStatusStringLock.Unlock()
		pr.exitStatusString = status
	}
	pr.cancel()
}

// Stop stops the process.
func (pr *processRunner) Stop() {
	// send SIGINT to the process
	if err := pr.cmd.Process.Signal(syscall.SIGINT); err != nil {
		log.WithError(err).Errorf("Error sending SIGINT for process %s", pr.title)
	}
}

// GetDidStart returns a channel that's closed when the process starts.
func (pr *processRunner) GetDidStart() <-chan bool {
	return pr.didStart
}

// GetTitle returns the title of the process.
func (pr *processRunner) GetTitle() string {
	return pr.title
}

// GetExitStatusString returns the exit status string of the process.
func (pr *processRunner) GetExitStatusString() string {
	pr.exitStatusStringLock.Lock()
	defer pr.exitStatusStringLock.Unlock()
	return pr.exitStatusString
}

// GetOutput returns the combined stdout and stderr output of the process.
func (pr *processRunner) GetOutput() string {
	return pr.outputBuf.String()
}
