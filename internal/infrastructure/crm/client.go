package crm

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/oauth"

	"github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/config"
)

// Client wraps amoCRM SDK
type Client struct {
	sdk *amocrm.SDK
}

// New creates a new CRM client based on auth mode
func New(cfg *config.Config) (*Client, error) {
	var sdk *amocrm.SDK
	var err error

	if cfg.AmoCRMAuthMode == config.AuthModeOAuth {
		// OAuth mode with auto-refresh
		provider := oauth.NewProvider(oauth.Config{
			ClientID:     cfg.AmoCRMClientID,
			ClientSecret: cfg.AmoCRMClientSecret,
			RedirectURI:  cfg.AmoCRMRedirectURI,
		})
		storage := oauth.NewFileStorage(".amocrm_tokens.json")

		sdk, err = amocrm.NewWithOAuth(provider, storage)
		if err != nil {
			return nil, fmt.Errorf("failed to init OAuth client: %w", err)
		}
	} else {
		// Token mode
		sdk = amocrm.New(cfg.AmoCRMBaseURL, cfg.AmoCRMToken)
	}

	return &Client{sdk: sdk}, nil
}

// SDK returns the underlying SDK for direct access
func (c *Client) SDK() *amocrm.SDK {
	return c.sdk
}

// Healthcheck checks API connectivity
func (c *Client) Healthcheck(ctx context.Context) error {
	_, err := c.sdk.Account().GetCurrent(ctx, nil)
	return err
}
