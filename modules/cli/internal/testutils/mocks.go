package testutils

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockProcessRunner is a mock type for the ProcessRunner interface
type MockProcessRunner struct {
	mock.Mock
}

func (m *MockProcessRunner) Restart() error {
	args := m.Called()
	return args.Error(0)
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

func (m *MockProcessRunner) GetOutputAndClearBuf() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockProcessRunner) GetEnvironmentPath() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockProcessRunner) GetEnvironment() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockProcessRunner) GetBinPath() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockProcessRunner) GetInfo() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockProcessRunner) CanWriteToLog() bool {
	return false
}

func (m *MockProcessRunner) WriteToLog(data string) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockProcessRunner) GetStartMinimized() bool {
	return false
}

func (m *MockProcessRunner) GetHighlightColor() string {
	return "blue"
}

func (m *MockProcessRunner) GetBorderColor() string {
	return "gray"
}
