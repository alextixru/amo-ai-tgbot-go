// Package auth provides OAuth authorization service for Telegram users.
package auth

import (
	"time"
)

// FlowState represents the current state of the auth flow for a user.
type FlowState string

const (
	// FlowStateNone - no active auth flow
	FlowStateNone FlowState = ""
	// FlowStateWaitingCode - user has opened auth URL, waiting for code input
	FlowStateWaitingCode FlowState = "waiting_code"
)

// PendingAuth represents an ongoing OAuth authorization flow.
type PendingAuth struct {
	TelegramUserID int64
	ChatID         int64
	State          string
	Verifier       string // PKCE code verifier
	FlowState      FlowState
	CreatedAt      time.Time
}

// IsExpired returns true if the pending auth has expired (default: 5 minutes).
func (p *PendingAuth) IsExpired() bool {
	return time.Since(p.CreatedAt) > 5*time.Minute
}

// IsWaitingCode returns true if the flow is waiting for authorization code.
func (p *PendingAuth) IsWaitingCode() bool {
	return p.FlowState == FlowStateWaitingCode && !p.IsExpired()
}

// OAuthStateStore manages pending OAuth authorization states.
type OAuthStateStore interface {
	// Save stores a pending auth state.
	Save(state string, auth *PendingAuth) error

	// Get retrieves a pending auth by state parameter.
	// Returns nil if not found or expired.
	Get(state string) (*PendingAuth, error)

	// GetByUserID retrieves a pending auth by Telegram user ID.
	// Returns nil if not found or expired.
	GetByUserID(telegramUserID int64) (*PendingAuth, error)

	// Delete removes a pending auth state.
	Delete(state string) error

	// UpdateFlowState updates the flow state for a user's pending auth.
	UpdateFlowState(telegramUserID int64, state FlowState) error
}
