package auth

import (
	"golang.org/x/oauth2"
)

// UserTokenStore manages OAuth tokens for Telegram users.
type UserTokenStore interface {
	// SaveToken stores a token for a Telegram user.
	SaveToken(telegramUserID int64, token *oauth2.Token) error

	// LoadToken retrieves a token for a Telegram user.
	// Returns nil if not found.
	LoadToken(telegramUserID int64) (*oauth2.Token, error)

	// DeleteToken removes a token for a Telegram user.
	DeleteToken(telegramUserID int64) error

	// HasToken returns true if a token exists for the user.
	HasToken(telegramUserID int64) bool
}
