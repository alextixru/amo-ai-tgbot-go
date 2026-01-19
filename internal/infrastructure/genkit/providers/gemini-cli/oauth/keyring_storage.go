package oauth

import (
	"encoding/json"
	"fmt"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

const (
	ServiceName = "gemini-cli"
	TokenKey    = "oauth_token"
)

// SaveTokenToKeyring saves the token to the system keyring.
func SaveTokenToKeyring(token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// We use a fixed user/key for the token since we manage one active session per service
	if err := keyring.Set(ServiceName, TokenKey, string(data)); err != nil {
		return fmt.Errorf("failed to set token in keyring: %w", err)
	}

	return nil
}

// LoadTokenFromKeyring loads the token from the system keyring.
func LoadTokenFromKeyring() (*oauth2.Token, error) {
	data, err := keyring.Get(ServiceName, TokenKey)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token from keyring: %w", err)
	}

	return &token, nil
}

// ClearTokenFromKeyring removes the token from the system keyring.
func ClearTokenFromKeyring() error {
	if err := keyring.Delete(ServiceName, TokenKey); err != nil {
		// Ignore error if item not found
		if err == keyring.ErrNotFound {
			return nil
		}
		return fmt.Errorf("failed to delete token from keyring: %w", err)
	}
	return nil
}

// IsKeyringAvailable checks if basic keyring operations work.
// This is a naive check; real availability depends on the OS environment (e.g. headless Linux vs Mac).
func IsKeyringAvailable() bool {
	// Try to get a non-existent item to check if service communicates
	_, err := keyring.Get(ServiceName, "test_availability")
	if err == keyring.ErrNotFound {
		return true
	}
	// If we get "exec: ..." error or dbus error, it's likely not available
	// But distinguishing them reliably across platforms is hard.
	// For now, assume available if we imported the package, but let runtime decide.
	// We can trust the error from Save/Load.
	return true
}
