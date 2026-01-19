package oauth

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

var (
	defaultManager *UserAccountManager
	once           sync.Once
)

// DefaultUserAccountManager returns the singleton instance of UserAccountManager
func DefaultUserAccountManager(cachePath string) *UserAccountManager {
	once.Do(func() {
		defaultManager = NewUserAccountManager(cachePath)
	})
	return defaultManager
}

// UserInfo содержит информацию о пользователе от Google OAuth API.
type UserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// UserAccounts хранит активный и исторические аккаунты.
type UserAccounts struct {
	Active string   `json:"active"`
	Old    []string `json:"old"`
}

// UserAccountManager управляет кешированием Google-аккаунтов.
type UserAccountManager struct {
	cachePath string
	mu        sync.RWMutex
}

// NewUserAccountManager создает новый экземпляр UserAccountManager.
func NewUserAccountManager(cachePath string) *UserAccountManager {
	return &UserAccountManager{
		cachePath: cachePath,
	}
}

// CacheGoogleAccount сохраняет email как активный аккаунт.
// Если уже был активный аккаунт, он перемещается в список старых.
func (uam *UserAccountManager) CacheGoogleAccount(email string) error {
	uam.mu.Lock()
	defer uam.mu.Unlock()

	accounts, err := uam.loadAccounts()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if accounts.Active != "" && accounts.Active != email {
		if !containsString(accounts.Old, accounts.Active) {
			accounts.Old = append(accounts.Old, accounts.Active)
		}
	}

	// Удалить новый email из списка старых, если он там был
	newOld := make([]string, 0, len(accounts.Old))
	for _, old := range accounts.Old {
		if old != email {
			newOld = append(newOld, old)
		}
	}
	accounts.Old = newOld
	accounts.Active = email

	return uam.saveAccounts(accounts)
}

// GetCachedGoogleAccount возвращает email текущего активного аккаунта.
func (uam *UserAccountManager) GetCachedGoogleAccount() (string, error) {
	uam.mu.RLock()
	defer uam.mu.RUnlock()

	accounts, err := uam.loadAccounts()
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return accounts.Active, nil
}

// GetLifetimeGoogleAccounts возвращает количество уникальных аккаунтов за все время.
func (uam *UserAccountManager) GetLifetimeGoogleAccounts() (int, error) {
	uam.mu.RLock()
	defer uam.mu.RUnlock()

	accounts, err := uam.loadAccounts()
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	unique := make(map[string]struct{})
	if accounts.Active != "" {
		unique[accounts.Active] = struct{}{}
	}
	for _, email := range accounts.Old {
		unique[email] = struct{}{}
	}

	return len(unique), nil
}

// ClearCachedGoogleAccount очищает активный аккаунт (перемещая его в старые).
func (uam *UserAccountManager) ClearCachedGoogleAccount() error {
	uam.mu.Lock()
	defer uam.mu.Unlock()

	accounts, err := uam.loadAccounts()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if accounts.Active != "" {
		if !containsString(accounts.Old, accounts.Active) {
			accounts.Old = append(accounts.Old, accounts.Active)
		}
		accounts.Active = ""
	}

	return uam.saveAccounts(accounts)
}

func (uam *UserAccountManager) loadAccounts() (UserAccounts, error) {
	data, err := os.ReadFile(uam.cachePath)
	if err != nil {
		return UserAccounts{}, err
	}

	// Validation for empty file
	if len(bytes.TrimSpace(data)) == 0 {
		return UserAccounts{}, nil
	}

	var accounts UserAccounts
	if err := json.Unmarshal(data, &accounts); err != nil {
		// Graceful fallback for invalid JSON
		return UserAccounts{}, nil
	}

	// Ensure Old slice is initialized
	if accounts.Old == nil {
		accounts.Old = make([]string, 0)
	}

	return accounts, nil
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (uam *UserAccountManager) saveAccounts(accounts UserAccounts) error {
	dir := filepath.Dir(uam.cachePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(uam.cachePath, data, 0600)
}
