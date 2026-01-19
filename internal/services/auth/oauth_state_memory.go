package auth

import (
	"sync"
)

// MemoryStateStore is an in-memory implementation of OAuthStateStore.
// Note: States are lost on application restart.
type MemoryStateStore struct {
	mu          sync.RWMutex
	states      map[string]*PendingAuth
	stateByUser map[int64]string // telegramUserID -> state
}

// NewMemoryStateStore creates a new in-memory state store.
func NewMemoryStateStore() *MemoryStateStore {
	return &MemoryStateStore{
		states:      make(map[string]*PendingAuth),
		stateByUser: make(map[int64]string),
	}
}

// Save stores a pending auth state.
func (m *MemoryStateStore) Save(state string, auth *PendingAuth) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Remove old state for this user if exists
	if oldState, exists := m.stateByUser[auth.TelegramUserID]; exists {
		delete(m.states, oldState)
	}

	m.states[state] = auth
	m.stateByUser[auth.TelegramUserID] = state
	return nil
}

// Get retrieves a pending auth by state parameter.
func (m *MemoryStateStore) Get(state string) (*PendingAuth, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	auth, exists := m.states[state]
	if !exists {
		return nil, nil
	}

	if auth.IsExpired() {
		return nil, nil
	}

	return auth, nil
}

// GetByUserID retrieves a pending auth by Telegram user ID.
func (m *MemoryStateStore) GetByUserID(telegramUserID int64) (*PendingAuth, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state, exists := m.stateByUser[telegramUserID]
	if !exists {
		return nil, nil
	}

	auth, exists := m.states[state]
	if !exists {
		return nil, nil
	}

	if auth.IsExpired() {
		return nil, nil
	}

	return auth, nil
}

// Delete removes a pending auth state.
func (m *MemoryStateStore) Delete(state string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if auth, exists := m.states[state]; exists {
		delete(m.stateByUser, auth.TelegramUserID)
		delete(m.states, state)
	}

	return nil
}

// UpdateFlowState updates the flow state for a user's pending auth.
func (m *MemoryStateStore) UpdateFlowState(telegramUserID int64, flowState FlowState) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	state, exists := m.stateByUser[telegramUserID]
	if !exists {
		return nil
	}

	if auth, exists := m.states[state]; exists {
		auth.FlowState = flowState
	}

	return nil
}

// DeleteByUserID removes a pending auth state by user ID.
func (m *MemoryStateStore) DeleteByUserID(telegramUserID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if state, exists := m.stateByUser[telegramUserID]; exists {
		delete(m.states, state)
		delete(m.stateByUser, telegramUserID)
	}

	return nil
}

// Cleanup removes expired states. Call periodically if needed.
func (m *MemoryStateStore) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for state, auth := range m.states {
		if auth.IsExpired() {
			delete(m.stateByUser, auth.TelegramUserID)
			delete(m.states, state)
		}
	}
}
