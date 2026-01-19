package geminicli

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	ClientID     = "681255809395-oo8ft2oprdrnp9e3aqf6av3hmdib135j.apps.googleusercontent.com"
	ClientSecret = "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"
	// OAuthTimeout is the maximum time to wait for OAuth authentication (5 minutes, matching gemini-cli)
	OAuthTimeout = 5 * time.Minute
	// NoBrowserRedirectURI is the official Google redirect URI for manual code copy-paste
	NoBrowserRedirectURI = "https://codeassist.google.com/authcode"
)

var Scopes = []string{
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",
}

func GetOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Scopes:       Scopes,
		Endpoint:     google.Endpoint,
	}
}

func GetTokenSource(ctx context.Context, credsPath string, noBrowser bool) (oauth2.TokenSource, error) {
	config := GetOAuthConfig()

	// 1. Try to load cached token
	token, err := LoadToken(credsPath)
	if err == nil {
		return &PersistingTokenSource{
			base:      config.TokenSource(ctx, token),
			credsPath: credsPath,
			lastToken: token,
		}, nil
	}

	// 2. No token or invalid, run OAuth flow
	if noBrowser {
		token, err = RunOAuthFlowNoBrowser(ctx, config)
	} else {
		token, err = RunOAuthFlow(ctx, config)
	}

	if err != nil {
		return nil, err
	}

	// 3. Save new token
	if err := SaveToken(credsPath, token); err != nil {
		fmt.Printf("Warning: failed to save token: %v\n", err)
	}

	return &PersistingTokenSource{
		base:      config.TokenSource(ctx, token),
		credsPath: credsPath,
		lastToken: token,
	}, nil
}

// PersistingTokenSource wraps an oauth2.TokenSource and saves tokens to disk on refresh.
type PersistingTokenSource struct {
	base      oauth2.TokenSource
	credsPath string
	lastToken *oauth2.Token
	mu        sync.Mutex
}

func (pts *PersistingTokenSource) Token() (*oauth2.Token, error) {
	pts.mu.Lock()
	defer pts.mu.Unlock()

	token, err := pts.base.Token()
	if err != nil {
		return nil, err
	}

	// Check if token was refreshed
	if pts.lastToken == nil || token.AccessToken != pts.lastToken.AccessToken {
		// Only save if either AccessToken changed or Expiry is significantly different
		// (though oauth2.Token.Expiry is usually enough, TS compares Credentials)
		if err := SaveToken(pts.credsPath, token); err != nil {
			fmt.Printf("Warning: failed to persist token: %v\n", err)
		}
		pts.lastToken = token
	}

	return token, nil
}

func LoadToken(path string) (*oauth2.Token, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func SaveToken(path string, token *oauth2.Token) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func RunOAuthFlow(parentCtx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	// Create context with timeout (5 minutes, matching gemini-cli)
	ctx, cancel := context.WithTimeout(parentCtx, OAuthTimeout)
	defer cancel()

	// Find available port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("failed to start listener: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	config.RedirectURL = fmt.Sprintf("http://localhost:%d/oauth2callback", port)

	// PKCE and State
	state := generateState()
	verifier := oauth2.GenerateVerifier()

	authURL := config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
	)

	fmt.Printf("\n--- Gemini CLI Authentication Required ---\n")
	fmt.Printf("Please open this URL in your browser to authorize Gemini CLI:\n\n%s\n\n", authURL)
	fmt.Printf("Waiting for authentication (timeout: %v)...\n", OAuthTimeout)

	codeChan := make(chan string)
	errChan := make(chan error)

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/oauth2callback" {
				http.NotFound(w, r)
				return
			}
			// Verify state to prevent CSRF
			if r.URL.Query().Get("state") != state {
				errMsg := "invalid oauth state"
				http.Error(w, errMsg, http.StatusForbidden)
				errChan <- errors.New(errMsg)
				return
			}

			code := r.URL.Query().Get("code")
			if code == "" {
				errMsg := "no code found in redirect URL"
				fmt.Fprintf(w, "Error: %s", errMsg)
				errChan <- errors.New(errMsg)
				return
			}
			fmt.Fprintf(w, "Authentication successful! You can close this tab.")
			codeChan <- code
		}),
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case code := <-codeChan:
		// Exchange code using PKCE verifier
		token, err := config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
		if err != nil {
			return nil, fmt.Errorf("failed to exchange code for token: %v", err)
		}
		_ = server.Shutdown(ctx)
		return token, nil
	case err := <-errChan:
		_ = server.Shutdown(ctx)
		return nil, err
	case <-ctx.Done():
		_ = server.Shutdown(context.Background())
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("authentication timed out after %v - the browser tab may have gotten stuck", OAuthTimeout)
		}
		return nil, ctx.Err()
	}
}

// RunOAuthFlowNoBrowser — авторизация через ручной ввод кода (для SSH/удаленных серверов)
func RunOAuthFlowNoBrowser(parentCtx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	ctx, cancel := context.WithTimeout(parentCtx, OAuthTimeout)
	defer cancel()

	// Find available port is NOT needed for manual flow, but config.RedirectURL MUST match one of the allowed redirect URIs.
	// Gemini CLI uses "https://codeassist.google.com/authcode" for manual flow.
	config.RedirectURL = NoBrowserRedirectURI

	// PKCE and State
	state := generateState()
	verifier := oauth2.GenerateVerifier()

	authURL := config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
	)

	fmt.Printf("\n--- Gemini CLI Authentication Required (NO_BROWSER mode) ---\n")
	fmt.Printf("Please open this URL in any browser to authorize Gemini CLI:\n\n%s\n\n", authURL)
	fmt.Printf("After authorization, copy the code provided on the page and paste it here.\n")
	fmt.Printf("Waiting for code (timeout: %v)...\n", OAuthTimeout)
	fmt.Printf("Enter authorization code: ")

	// Read code from stdin
	codeChan := make(chan string)
	errChan := make(chan error)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		code, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err
			return
		}
		codeChan <- strings.TrimSpace(code)
	}()

	select {
	case code := <-codeChan:
		if code == "" {
			return nil, errors.New("authorization code is required")
		}
		// Exchange code using PKCE verifier
		token, err := config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
		if err != nil {
			return nil, fmt.Errorf("failed to exchange code for token: %v", err)
		}
		return token, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("authentication timed out after %v", OAuthTimeout)
		}
		return nil, ctx.Err()
	}
}
