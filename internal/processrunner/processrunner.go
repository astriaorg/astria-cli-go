package processrunner

import (
	"fmt"
	"io"
	"os/exec"
)

// ProcessRunner is a struct that represents a process to be run.
type ProcessRunner struct {
	*exec.Cmd
	binPath     string
	environment []string
	args        []string
	IsRunning   chan bool
	Title       string
	Stdout      io.ReadCloser
	Stderr      io.ReadCloser
}

// NewProcessRunner creates a new ProcessRunner with a title, binPath, and environment.
func NewProcessRunner(title string, binPath string, environment []string, args []string) *ProcessRunner {
	command := exec.Command(binPath, args...)
	command.Env = environment
	return &ProcessRunner{
		Cmd:         command,
		binPath:     binPath,
		environment: environment,
		args:        args,
		IsRunning:   make(chan bool),
		Title:       title,
	}
}

// Start starts the process and returns the ProcessRunner and an error.
func (pr *ProcessRunner) Start(depIsRunning chan bool) (*ProcessRunner, error) {
	stdout, err := pr.Cmd.StdoutPipe()
	if err != nil {
		return pr, err
	}
	pr.Stdout = stdout

	stderr, err := pr.Cmd.StderrPipe()
	if err != nil {
		return pr, err
	}
	pr.Stderr = stderr

	// start the process
	go func() {
		// wait for the dependency to start
		<-depIsRunning
		err = pr.Cmd.Start()
		if err != nil {
			fmt.Println("error starting process:", err)
			panic(err)
		}

		// signal that this process is running so
		pr.IsRunning <- true
	}()

	return pr, nil
}
