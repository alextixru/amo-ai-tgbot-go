package crm

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/oauth"

	"github.com/tihn/amo-ai-tgbot-go/internal/config"
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

// GetAccountInfo returns account information as formatted string
func (c *Client) GetAccountInfo(ctx context.Context) (string, error) {
	account, err := c.sdk.Account().GetCurrent(ctx, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"🏢 Аккаунт: %s\n"+
			"🆔 ID: %d\n"+
			"🌐 Subdomain: %s",
		account.Name,
		account.ID,
		account.Subdomain,
	), nil
}

// GetPipelines returns pipelines with statuses as formatted string
func (c *Client) GetPipelines(ctx context.Context) (string, error) {
	pipelines, err := c.sdk.Pipelines().Get(ctx)
	if err != nil {
		return "", err
	}

	if len(pipelines) == 0 {
		return "📭 Воронок нет", nil
	}

	var result string
	for _, p := range pipelines {
		result += fmt.Sprintf("📊 %s (ID: %d)\n", p.Name, p.ID)

		// Get statuses for this pipeline
		statuses, err := c.sdk.Pipelines().GetStatuses(ctx, p.ID)
		if err != nil {
			result += fmt.Sprintf("   ⚠️ Ошибка загрузки статусов: %v\n", err)
			continue
		}

		for i, s := range statuses {
			prefix := "├─"
			if i == len(statuses)-1 {
				prefix = "└─"
			}
			result += fmt.Sprintf("   %s %s (ID: %d)\n", prefix, s.Name, s.ID)
		}
		result += "\n"
	}

	return result, nil
}
