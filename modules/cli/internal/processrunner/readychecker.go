package processrunner

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// ReadyChecker is a struct used within the ProcessRunner to check if the
// process being run has completed all its startup steps.
type ReadyChecker struct {
	// callBackName is the name of the callback function and is used for logging purposes.
	callBackName string
	// callback is the anonymous function that will be called to check if all
	// startup requirements for the process have been completed. The function
	// should return true if all startup checks are complete, and false if any
	// startup checks have not completed.
	callback      func() bool
	retryCount    int
	retryInterval time.Duration
	// haltIfFailed is a flag that determines if the process should halt the app or
	// continue if all retries of the callback complete without success.
	haltIfFailed bool
}

// ReadyCheckerOpts is a struct used to pass options into NewReadyChecker.
type ReadyCheckerOpts struct {
	CallBackName  string
	Callback      func() bool
	RetryCount    int
	RetryInterval time.Duration
	HaltIfFailed  bool
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

// waitUntilReady calls the ReadyChecker.callback function N number of times,
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
			log.Debug(fmt.Sprintf("ReadyChecker callback to '%s' completed successfully.", r.callBackName))
			return nil
		}
		log.Debug(fmt.Sprintf("ReadyChecker callback to '%s': attempt %d, failed to complete. Retrying...", r.callBackName, i+1))
		time.Sleep(r.retryInterval)
	}
	complete := r.callback()
	if !complete && r.haltIfFailed {
		err := fmt.Errorf("ReadyChecker callback to '%s' failed to complete after %d retries. Halting", r.callBackName, r.retryCount)
		panic(err)
	}
	return nil
}
