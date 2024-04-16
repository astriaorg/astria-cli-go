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

func (m *MockProcessRunner) GetOutputAndClearBuf() string {
	args := m.Called()
	return args.String(0)
}
