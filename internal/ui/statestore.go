package ui

import "sync"

// AppState is a struct containing the state of the application
type AppState struct {
	isAutoScroll bool
	isWordWrap   bool
	isBorderless bool
}

// StateStore is a struct that controls the state of the application
type StateStore struct {
	state AppState
	mutex sync.Mutex
}

// NewStateStore creates a new StateStore
func NewStateStore() *StateStore {
	return &StateStore{
		state: AppState{
			isAutoScroll: true,
			isWordWrap:   false,
			isBorderless: false,
		},
	}
}

// ToggleAutoscroll toggles the autoscroll state
func (s *StateStore) ToggleAutoscroll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.state.isAutoScroll = !s.state.isAutoScroll
}

// GetIsAutoscroll returns the autoscroll state
func (s *StateStore) GetIsAutoscroll() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.state.isAutoScroll
}

// ToggleWordWrap toggles the word wrap state
func (s *StateStore) ToggleWordWrap() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.state.isWordWrap = !s.state.isWordWrap
}

// GetIsWordWrap returns the word wrap state
func (s *StateStore) GetIsWordWrap() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.state.isWordWrap
}

// ToggleBorderless toggles the borderless state
func (s *StateStore) ToggleBorderless() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.state.isBorderless = !s.state.isBorderless
}

// ResetBorderless resets the borderless state to false or "off"
func (s *StateStore) ResetBorderless() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.state.isBorderless = false
}

// GetIsBorderless returns the borderless state
func (s *StateStore) GetIsBorderless() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.state.isBorderless
}
