package processrunner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/astria/astria-cli-go/internal/safebuffer"
	log "github.com/sirupsen/logrus"
)

// ProcessRunner is an interface that represents a process to be run.
type ProcessRunner interface {
	Restart() error
	Start(ctx context.Context, depStarted <-chan bool) error
	Stop()
	GetDidStart() <-chan bool
	GetTitle() string
	GetOutputAndClearBuf() string
	GetInfo() string
	GetEnvironment() []string
}

// ProcessRunner is a struct that represents a process to be run.
type processRunner struct {
	// cmd is the exec.Cmd to be run
	cmd *exec.Cmd
	// Title is the title of the process
	title string

	// saving the opts so we can use them for restarts
	opts NewProcessRunnerOpts
	// Env is the environment variables for the process
	env []string
	// NOTE - only saving the context on a struct so that we can use it to restart the process
	ctx context.Context

	didStart  chan bool
	outputBuf *safebuffer.SafeBuffer

	readyChecker *ReadyChecker
}

type NewProcessRunnerOpts struct {
	Title        string
	BinPath      string
	EnvPath      string
	EnvOverrides []string
	Args         []string
	ReadyCheck   *ReadyChecker
}

// NewProcessRunner creates a new ProcessRunner.
// It creates a new exec.Cmd with the given binPath and args, and sets the
// environment. If no envPath is provided, it uses the current environment using
// os.Environ().
func NewProcessRunner(ctx context.Context, opts NewProcessRunnerOpts) ProcessRunner {
	var env []string
	if opts.EnvPath != "" {
		env = GetEnvironment(opts.EnvPath)
	} else {
		env = os.Environ()
	}

	if opts.EnvOverrides != nil {
		env = applyEnvOverrides(env, opts.EnvOverrides)
	}

	// using exec.CommandContext to allow for cancellation from caller
	cmd := exec.CommandContext(ctx, opts.BinPath, opts.Args...)
	cmd.Env = env
	return &processRunner{
		ctx:          ctx,
		cmd:          cmd,
		title:        opts.Title,
		didStart:     make(chan bool),
		outputBuf:    &safebuffer.SafeBuffer{},
		opts:         opts,
		env:          env,
		readyChecker: opts.ReadyCheck,
	}
}

// Restart stops the process and starts it again.
func (pr *processRunner) Restart() error {
	log.Debug(fmt.Sprintf("Stopping process %s", pr.title))
	pr.Stop()

	// NOTE - you have to recreate the exec.Cmd. you can't just call cmd.Start() again.
	cmd := exec.CommandContext(pr.ctx, pr.opts.BinPath, pr.opts.Args...)
	// setting env again
	cmd.Env = pr.env
	pr.cmd = cmd

	// must recreate the didStart channel because it was previously closed
	pr.didStart = make(chan bool)

	// must create a new channel that triggers process start
	shouldStart := make(chan bool)
	close(shouldStart)

	// start the process again
	err := pr.Start(pr.ctx, shouldStart)
	if err != nil {
		return err
	}

	// add a new line for visual separation
	s := fmt.Sprintf("\n[black:white][astria-go] %s process restarted[-:-]\n", pr.title)
	_, err = pr.outputBuf.WriteString(s)
	if err != nil {
		return err
	}

	return nil
}

// Start starts the process and returns the ProcessRunner and an error.
// It takes a channel that's closed when the dependency starts.
// This allows us to control the order of process startup.
func (pr *processRunner) Start(ctx context.Context, depStarted <-chan bool) error {
	log.Debug(fmt.Sprintf("Starting process %s", pr.title))

	select {
	case <-depStarted:
	// continue if the dependency has started.
	case <-ctx.Done():
		log.Info("Context cancelled before starting process", pr.title)
		return ctx.Err()
	}

	// get stdout and stderr
	stdout, err := pr.cmd.StdoutPipe()
	if err != nil {
		log.WithError(err).Errorf("Error obtaining stdout for process %s", pr.title)
		return err
	}
	stderr, err := pr.cmd.StderrPipe()
	if err != nil {
		log.WithError(err).Errorf("Error obtaining stderr for process %s", pr.title)
		return err
	}

	// multiwriter to write both stdout and stderr to the same buffer
	go func() {
		_, err := io.Copy(pr.outputBuf, stdout)
		if err != nil {
			log.WithError(err).Error("Error during io.Copy to stdout")
		}
	}()
	go func() {
		_, err := io.Copy(pr.outputBuf, stderr)
		if err != nil {
			log.WithError(err).Error("Error during io.Copy to stderr")
		}
	}()

	// actually start the process
	if err := pr.cmd.Start(); err != nil {
		log.WithError(err).Errorf("Error starting process %s", pr.title)
		return err
	}

	// run the readiness check if present
	if pr.readyChecker != nil {
		err := pr.readyChecker.waitUntilReady()
		if err != nil {
			log.WithError(err).Errorf("Error when running readiness check for %s", pr.title)
		}
	}

	// signal that this process has started.
	close(pr.didStart)

	// asynchronously monitor process
	go func() {
		err = pr.cmd.Wait()
		if err != nil {
			logErr := fmt.Errorf("%s process exited with error: %w", pr.title, err)
			outputErr := fmt.Errorf("[white:red][astria-go] %s[-:-]", logErr)
			log.Error(logErr)
			_, err := pr.outputBuf.WriteString(outputErr.Error())
			if err != nil {
				return
			}
		} else {
			exitStatusMessage := fmt.Sprintf("%s process exited cleanly", pr.title)
			outputStatusMessage := fmt.Sprintf("[black:white][astria-go] %s[-:-]", exitStatusMessage)
			log.Infof(exitStatusMessage)
			_, err := pr.outputBuf.WriteString(outputStatusMessage)
			if err != nil {
				return
			}
		}
	}()

	return nil
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

// GetOutputAndClearBuf returns the combined stdout and stderr output of the process.
func (pr *processRunner) GetOutputAndClearBuf() string {
	defer pr.outputBuf.Reset()

	o := pr.outputBuf.String()

	return o
}

// GetInfo returns the formatted binary path and environment path of the process.
func (pr *processRunner) GetInfo() string {
	binaryPathTitle := " " + pr.GetTitle() + " binary path:"
	environmentPathTitle := " Environment path:"
	var maxLen int
	if len(binaryPathTitle) > len(environmentPathTitle) {
		maxLen = len(binaryPathTitle)
	} else {
		maxLen = len(environmentPathTitle)
	}
	output := ""
	output += fmt.Sprintf("%-*s", maxLen+1, binaryPathTitle) + pr.opts.BinPath + "\n"
	output += fmt.Sprintf("%-*s", maxLen+1, environmentPathTitle) + pr.opts.EnvPath + "\n"
	return output
}

// GetEnvironment returns the environment variables for the process.
func (pr *processRunner) GetEnvironment() []string {
	return pr.env
}

// applyEnvOverrides updates the environment variable with the overrides environment variables.
func applyEnvOverrides(env []string, overrides []string) []string {
	// create a map of the environment variables
	envMap := make(map[string]string)
	for _, e := range env {
		kv := strings.SplitN(e, "=", 2)
		envMap[kv[0]] = kv[1]
	}

	// apply the overrides
	for _, e := range overrides {
		kv := strings.SplitN(e, "=", 2)
		envMap[kv[0]] = kv[1]
	}

	// convert the map back to a slice
	var newEnv []string
	for k, v := range envMap {
		newEnv = append(newEnv, fmt.Sprintf("%s=%s", k, v))
	}

	return newEnv
}
