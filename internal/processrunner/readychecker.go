package processrunner

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// ReadyChecker is a struct used within the ProcessRunner to check if the
// process being run has completed all its startup steps.
type ReadyChecker struct {
	callBackName  string
	callback      func() bool
	startDelay    time.Duration
	retryCount    int
	retryInterval time.Duration
	haltIfFailed  bool
}

// ReadyCheckerOpts is a struct used to pass options into NewReadyChecker.
type ReadyCheckerOpts struct {
	// CallBackName is the name of the callback function and is used for logging purposes.
	CallBackName string
	// Callback is the function that will be called to check if the process is ready.
	Callback      func() bool
	RetryCount    int
	RetryInterval time.Duration
	// HaltIfFailed is a flag that determines if the process should halt the app or
	// continue if all retries of the callback complete without success.
	HaltIfFailed bool
}

// NewReadyChecker creates a new ReadyChecker.
func NewReadyChecker(opts ReadyCheckerOpts) ReadyChecker {
	return ReadyChecker{
		callBackName:  opts.CallBackName,
		callback:      opts.Callback,
		retryCount:    opts.RetryCount,
		retryInterval: opts.RetryInterval,
		haltIfFailed:  opts.HaltIfFailed,
	}
}

// WaitUntilReady calls the ReadyChecker.callback function N number of times,
// waiting M amount of time between retries, where N = ReadyChecker.retryCount
// and M = ReadyChecker.retryInterval.
// If the callback returns true, the function returns nil.
// If ReadyChecker.haltIfFailed is false, the function will return nil after all
// retries have been completed without success.
// If ReadyChecker.haltIfFailed is true, the function will panic if the callback
// does not succeed after all retries.
func (r *ReadyChecker) waitUntilReady() error {
	for i := 0; i < r.retryCount-1; i++ {
		complete := r.callback()
		if complete {
			log.Info(fmt.Sprintf("ReadyChecker callback to '%s' completed successfully.", r.callBackName))
			return nil
		}
		log.Debug(fmt.Sprintf("ReadyChecker callback to '%s' run %d failed to complete. Retrying...", r.callBackName, i+1))
		time.Sleep(r.retryInterval)
	}
	complete := r.callback()
	if !complete && r.haltIfFailed {
		err := fmt.Errorf("ReadyChecker callback to '%s' failed to complete after %d retries. Halting.", r.callBackName, r.retryCount)
		panic(err)
	}
	return nil
}
