package admin_integrations

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListWebhooks(ctx context.Context, filter *filters.WebhooksFilter) ([]models.Webhook, error) {
	webhooks, _, err := s.sdk.Webhooks().Get(ctx, filter)
	return webhooks, err
}

func (s *service) SubscribeWebhook(ctx context.Context, destination string, settings []string) (*models.Webhook, error) {
	webhook := &models.Webhook{
		Destination: destination,
		Settings:    settings,
	}
	return s.sdk.Webhooks().Subscribe(ctx, webhook)
}

func (s *service) UnsubscribeWebhook(ctx context.Context, destination string, settings []string) error {
	webhook := &models.Webhook{
		Destination: destination,
		Settings:    settings,
	}
	return s.sdk.Webhooks().Unsubscribe(ctx, webhook)
}
