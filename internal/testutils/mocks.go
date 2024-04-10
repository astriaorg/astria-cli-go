package testutils

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockProcessRunner is a mock type for the ProcessRunner interface
type MockProcessRunner struct {
	mock.Mock
}

func (m *MockProcessRunner) Start(ctx context.Context, depStarted <-chan bool) error {
	args := m.Called(ctx, depStarted)
	return args.Error(0)
}

func (m *MockProcessRunner) Stop() {
	m.Called()
}

func (m *MockProcessRunner) GetDidStart() <-chan bool {
	args := m.Called()
	return args.Get(0).(<-chan bool)
}

func (m *MockProcessRunner) GetTitle() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockProcessRunner) GetOutput() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockProcessRunner) GetLineCount() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockProcessRunner) CanWriteToLog() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockProcessRunner) WriteToLog(data string) error {
	args := m.Called(data)
	return args.Error(0)
}
