package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/genkit/providers/gemini-cli/oauth"
	"golang.org/x/oauth2"
)

// Service handles OAuth authorization for Telegram users.
type Service struct {
	stateStore OAuthStateStore
	tokenStore UserTokenStore
	mu         sync.RWMutex
}

// NewService creates a new auth service.
func NewService(stateStore OAuthStateStore, tokenStore UserTokenStore) *Service {
	return &Service{
		stateStore: stateStore,
		tokenStore: tokenStore,
	}
}

// StartAuth generates an OAuth URL for a Telegram user.
// Returns the authorization URL that the user should open in their browser.
func (s *Service) StartAuth(telegramUserID, chatID int64) (authURL string, err error) {
	// Generate OAuth URL with PKCE
	url, state, verifier, err := oauth.GetAuthURL()
	if err != nil {
		return "", fmt.Errorf("failed to generate auth URL: %w", err)
	}

	// Store pending auth with waiting state
	pending := &PendingAuth{
		TelegramUserID: telegramUserID,
		ChatID:         chatID,
		State:          state,
		Verifier:       verifier,
		FlowState:      FlowStateWaitingCode,
		CreatedAt:      time.Now(),
	}

	if err := s.stateStore.Save(state, pending); err != nil {
		return "", fmt.Errorf("failed to save pending auth: %w", err)
	}

	return url, nil
}

// IsWaitingCode returns true if user is in the middle of auth flow waiting for code.
func (s *Service) IsWaitingCode(telegramUserID int64) bool {
	pending, err := s.stateStore.GetByUserID(telegramUserID)
	if err != nil || pending == nil {
		return false
	}
	return pending.IsWaitingCode()
}

// CancelAuth cancels the pending authorization flow.
func (s *Service) CancelAuth(telegramUserID int64) error {
	pending, err := s.stateStore.GetByUserID(telegramUserID)
	if err != nil {
		return err
	}
	if pending != nil {
		return s.stateStore.Delete(pending.State)
	}
	return nil
}

// CompleteAuth exchanges the authorization code for tokens.
// The user provides the code they received after authorizing in Google.
func (s *Service) CompleteAuth(ctx context.Context, telegramUserID int64, code string) error {
	// Get pending auth for this user
	pending, err := s.stateStore.GetByUserID(telegramUserID)
	if err != nil {
		return fmt.Errorf("failed to get pending auth: %w", err)
	}
	if pending == nil {
		return fmt.Errorf("no pending authorization found. Please start with /connect first")
	}

	// Exchange code for token using the stored verifier
	config := oauth.GetOAuthConfig()
	config.RedirectURL = oauth.NoBrowserRedirectURI

	token, err := config.Exchange(ctx, code, oauth2.VerifierOption(pending.Verifier))
	if err != nil {
		return fmt.Errorf("failed to exchange code: %w", err)
	}

	// Fetch user info to get email
	email := s.fetchUserEmail(ctx, config, token)

	// Save token with email for this user
	if fileStore, ok := s.tokenStore.(*FileTokenStore); ok {
		if err := fileStore.SaveTokenWithEmail(telegramUserID, token, email); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}
	} else {
		if err := s.tokenStore.SaveToken(telegramUserID, token); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}
	}

	// Clean up pending auth
	if err := s.stateStore.Delete(pending.State); err != nil {
		// Non-fatal, log but don't fail
		fmt.Printf("Warning: failed to delete pending auth: %v\n", err)
	}

	return nil
}

// fetchUserEmail fetches the user's email from Google userinfo API.
func (s *Service) fetchUserEmail(ctx context.Context, config *oauth2.Config, token *oauth2.Token) string {
	httpClient := oauth2.NewClient(ctx, config.TokenSource(ctx, token))
	resp, err := httpClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
	}

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return ""
	}
	return userInfo.Email
}

// GetUserEmail returns the stored email for a Telegram user.
func (s *Service) GetUserEmail(telegramUserID int64) string {
	if fileStore, ok := s.tokenStore.(*FileTokenStore); ok {
		email, _ := fileStore.LoadUserEmail(telegramUserID)
		return email
	}
	return ""
}

// GetTokenSource returns an oauth2.TokenSource for a Telegram user.
// Returns nil if the user is not authenticated.
func (s *Service) GetTokenSource(ctx context.Context, telegramUserID int64) (oauth2.TokenSource, error) {
	token, err := s.tokenStore.LoadToken(telegramUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to load token: %w", err)
	}
	if token == nil {
		return nil, nil
	}

	config := oauth.GetOAuthConfig()

	// Create a persisting token source that auto-refreshes and saves
	return &userTokenSource{
		base:           config.TokenSource(ctx, token),
		tokenStore:     s.tokenStore,
		telegramUserID: telegramUserID,
		lastToken:      token,
	}, nil
}

// IsAuthenticated returns true if the user has a valid token.
func (s *Service) IsAuthenticated(telegramUserID int64) bool {
	token, err := s.tokenStore.LoadToken(telegramUserID)
	if err != nil || token == nil {
		return false
	}

	// Check if token is valid (has access token and not expired within 5 min buffer)
	if token.AccessToken == "" {
		return false
	}
	if token.Expiry.IsZero() {
		return true // No expiry, assume valid (will be refreshed if needed)
	}
	return token.Expiry.After(time.Now().Add(5 * time.Minute))
}

// Logout removes tokens for a Telegram user.
func (s *Service) Logout(telegramUserID int64) error {
	return s.tokenStore.DeleteToken(telegramUserID)
}

// GetPendingAuth returns the pending auth for a user if exists.
func (s *Service) GetPendingAuth(telegramUserID int64) (*PendingAuth, error) {
	return s.stateStore.GetByUserID(telegramUserID)
}

// userTokenSource wraps oauth2.TokenSource and persists refreshed tokens.
type userTokenSource struct {
	base           oauth2.TokenSource
	tokenStore     UserTokenStore
	telegramUserID int64
	lastToken      *oauth2.Token
	mu             sync.Mutex
}

func (ts *userTokenSource) Token() (*oauth2.Token, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	token, err := ts.base.Token()
	if err != nil {
		return nil, err
	}

	// Save if token was refreshed
	if ts.lastToken == nil || token.AccessToken != ts.lastToken.AccessToken {
		if err := ts.tokenStore.SaveToken(ts.telegramUserID, token); err != nil {
			fmt.Printf("Warning: failed to persist refreshed token: %v\n", err)
		}
		ts.lastToken = token
	}

	return token, nil
}
