package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/oauth2"
)

// UserData stores token along with user info (email)
type UserData struct {
	Token *oauth2.Token `json:"token"`
	Email string        `json:"email,omitempty"`
}

// FileTokenStore stores tokens as JSON files per user.
// Path: {baseDir}/tokens/{telegram_user_id}.json
type FileTokenStore struct {
	baseDir string
	mu      sync.RWMutex
}

// NewFileTokenStore creates a new file-based token store.
func NewFileTokenStore(baseDir string) *FileTokenStore {
	return &FileTokenStore{
		baseDir: baseDir,
	}
}

func (f *FileTokenStore) tokenPath(telegramUserID int64) string {
	return filepath.Join(f.baseDir, "tokens", fmt.Sprintf("%d.json", telegramUserID))
}

// SaveToken stores a token for a Telegram user.
func (f *FileTokenStore) SaveToken(telegramUserID int64, token *oauth2.Token) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.saveUserDataLocked(telegramUserID, token, "")
}

// SaveTokenWithEmail stores a token with email for a Telegram user.
func (f *FileTokenStore) SaveTokenWithEmail(telegramUserID int64, token *oauth2.Token, email string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.saveUserDataLocked(telegramUserID, token, email)
}

func (f *FileTokenStore) saveUserDataLocked(telegramUserID int64, token *oauth2.Token, email string) error {
	path := f.tokenPath(telegramUserID)
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create token directory: %w", err)
	}

	// Try to preserve existing email if not provided
	if email == "" {
		if existing, err := f.loadUserDataLocked(telegramUserID); err == nil && existing != nil {
			email = existing.Email
		}
	}

	userData := &UserData{
		Token: token,
		Email: email,
	}

	data, err := json.MarshalIndent(userData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// LoadToken retrieves a token for a Telegram user.
func (f *FileTokenStore) LoadToken(telegramUserID int64) (*oauth2.Token, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	userData, err := f.loadUserDataLocked(telegramUserID)
	if err != nil {
		return nil, err
	}
	if userData == nil {
		return nil, nil
	}
	return userData.Token, nil
}

// LoadUserEmail retrieves the email for a Telegram user.
func (f *FileTokenStore) LoadUserEmail(telegramUserID int64) (string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	userData, err := f.loadUserDataLocked(telegramUserID)
	if err != nil {
		return "", err
	}
	if userData == nil {
		return "", nil
	}
	return userData.Email, nil
}

func (f *FileTokenStore) loadUserDataLocked(telegramUserID int64) (*UserData, error) {
	path := f.tokenPath(telegramUserID)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	// Try new UserData format first
	var userData UserData
	if err := json.Unmarshal(data, &userData); err == nil && userData.Token != nil {
		return &userData, nil
	}

	// Fallback: try legacy oauth2.Token format
	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &UserData{Token: &token}, nil
}

// DeleteToken removes a token for a Telegram user.
func (f *FileTokenStore) DeleteToken(telegramUserID int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	path := f.tokenPath(telegramUserID)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete token file: %w", err)
	}

	return nil
}

// HasToken returns true if a token exists for the user.
func (f *FileTokenStore) HasToken(telegramUserID int64) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	path := f.tokenPath(telegramUserID)
	_, err := os.Stat(path)
	return err == nil
}
