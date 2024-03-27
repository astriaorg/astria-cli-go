package testutils

import (
	"bufio"

	"github.com/stretchr/testify/mock"
)

type MockProcessRunner struct {
	mock.Mock
}

func (m *MockProcessRunner) Start(depStarted <-chan bool) error {
	args := m.Called(depStarted)
	return args.Error(0)
}
func (m *MockProcessRunner) Wait() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockProcessRunner) Stop() {
}

func (m *MockProcessRunner) GetDidStart() <-chan bool {
	args := m.Called()
	return args.Get(0).(<-chan bool)
}

func (m *MockProcessRunner) GetTitle() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockProcessRunner) GetScanner() *bufio.Scanner {
	args := m.Called()
	return args.Get(0).(*bufio.Scanner)
}
