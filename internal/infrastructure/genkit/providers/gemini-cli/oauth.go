package geminicli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	ClientID     = "681255809395-oo8ft2oprdrnp9e3aqf6av3hmdib135j.apps.googleusercontent.com"
	ClientSecret = "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"
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
		RedirectURL:  "http://localhost:0/oauth2callback", // Will be updated with actual port
	}
}

func GetTokenSource(ctx context.Context, credsPath string) (oauth2.TokenSource, error) {
	config := GetOAuthConfig()

	// 1. Try to load cached token
	token, err := LoadToken(credsPath)
	if err == nil {
		return config.TokenSource(ctx, token), nil
	}

	// 2. No token or invalid, run OAuth flow
	token, err = RunOAuthFlow(ctx, config)
	if err != nil {
		return nil, err
	}

	// 3. Save new token
	if err := SaveToken(credsPath, token); err != nil {
		fmt.Printf("Warning: failed to save token: %v\n", err)
	}

	return config.TokenSource(ctx, token), nil
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

func RunOAuthFlow(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	// Find available port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("failed to start listener: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	config.RedirectURL = fmt.Sprintf("http://localhost:%d/oauth2callback", port)

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("\n--- Gemini CLI Authentication Required ---\n")
	fmt.Printf("Please open this URL in your browser to authorize Gemini CLI:\n\n%s\n\n", authURL)

	codeChan := make(chan string)
	errChan := make(chan error)

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/oauth2callback" {
				http.NotFound(w, r)
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
		token, err := config.Exchange(ctx, code)
		if err != nil {
			return nil, fmt.Errorf("failed to exchange code for token: %v", err)
		}
		_ = server.Shutdown(ctx)
		return token, nil
	case err := <-errChan:
		_ = server.Shutdown(ctx)
		return nil, err
	case <-ctx.Done():
		_ = server.Shutdown(ctx)
		return nil, ctx.Err()
	}
}
