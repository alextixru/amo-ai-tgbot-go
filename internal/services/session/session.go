// Package session provides chat history storage for multi-turn conversations.
package session

import (
	"sync"

	"github.com/firebase/genkit/go/ai"
)

// MaxHistoryMessages limits the number of messages stored per session.
// Prevents memory overflow in long conversations.
const MaxHistoryMessages = 20

// Store defines the interface for chat history persistence.
type Store interface {
	// Load retrieves chat history for a session.
	Load(sessionID string) []*ai.Message

	// Save stores chat history for a session.
	Save(sessionID string, history []*ai.Message)

	// Clear removes chat history for a session.
	Clear(sessionID string)
}

// MemoryStore is an in-memory implementation of Store.
// Note: History is lost on application restart.
// For persistence, migrate to Redis or database implementation.
type MemoryStore struct {
	mu    sync.RWMutex
	store map[string][]*ai.Message
}

// NewMemoryStore creates a new in-memory session store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[string][]*ai.Message),
	}
}

// Load retrieves chat history for a session.
func (m *MemoryStore) Load(sessionID string) []*ai.Message {
	m.mu.RLock()
	defer m.mu.RUnlock()

	history, ok := m.store[sessionID]
	if !ok {
		return nil
	}

	// Return a copy to avoid race conditions
	result := make([]*ai.Message, len(history))
	copy(result, history)
	return result
}

// Save stores chat history for a session.
// Trims history to MaxHistoryMessages if exceeded.
func (m *MemoryStore) Save(sessionID string, history []*ai.Message) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Trim to max messages, keeping most recent
	if len(history) > MaxHistoryMessages {
		history = history[len(history)-MaxHistoryMessages:]
	}

	// Store a copy
	stored := make([]*ai.Message, len(history))
	copy(stored, history)
	m.store[sessionID] = stored
}

// Clear removes chat history for a session.
func (m *MemoryStore) Clear(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, sessionID)
}
