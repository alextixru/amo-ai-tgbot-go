package oauth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

// GetAuthURL генерирует URL для авторизации, state и verifier (для PKCE).
// Клиент должен сохранить state и verifier для последующего вызова ExchangeCode.
func GetAuthURL() (authURL, state, verifier string, err error) {
	config := GetOAuthConfig()
	// Для headless/manual flow используем специальный redirect URI
	config.RedirectURL = NoBrowserRedirectURI

	state = generateState()
	verifier = oauth2.GenerateVerifier()

	authURL = config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
	)

	return authURL, state, verifier, nil
}

// ExchangeCode обменивает код авторизации на токен и сохраняет его.
// state должен совпадать со state из GetAuthURL для защиты от CSRF.
func ExchangeCode(ctx context.Context, code, state, expectedState, verifier string) error {
	if state != expectedState {
		return errors.New("state mismatch: possible CSRF attack")
	}

	config := GetOAuthConfig()
	config.RedirectURL = NoBrowserRedirectURI

	// Exchange code using PKCE verifier
	token, err := config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return fmt.Errorf("failed to exchange code for token: %w", err)
	}

	credsPath := GetDefaultCredsPath()
	if err := SaveToken(credsPath, token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	// Fetch and cache user info
	uam := DefaultUserAccountManager(GetGoogleAccountsPath())
	httpClient := oauth2.NewClient(ctx, config.TokenSource(ctx, token))
	if err := fetchAndCacheUserInfo(ctx, httpClient, uam); err != nil {
		// Log warning but don't fail the complete flow?
		// For headless, maybe we should return error if we can't identify the user.
		// But usually auth is successful even if userinfo fails.
		fmt.Printf("Warning: failed to fetch user info: %v\n", err)
	}

	return nil
}

// IsAuthenticated проверяет наличие валидного токена с buffer'ом 5 минут.
func IsAuthenticated() bool {
	credsPath := GetDefaultCredsPath()
	token, err := LoadToken(credsPath)
	if err != nil {
		return false
	}
	// Check if token is valid and not expiring within 5 minutes
	if token.AccessToken == "" {
		return false
	}
	if token.Expiry.IsZero() {
		return true // No expiry set, assume valid
	}
	return token.Expiry.After(time.Now().Add(5 * time.Minute))
}
