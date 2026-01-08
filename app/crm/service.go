package crm

import (
	"context"
	"fmt"

	infraCRM "github.com/tihn/amo-ai-tgbot-go/infrastructure/crm"
)

// Service provides CRM business logic
type Service struct {
	client *infraCRM.Client
}

// NewService creates a new CRM service
func NewService(client *infraCRM.Client) *Service {
	return &Service{client: client}
}

// Client returns the underlying CRM client
func (s *Service) Client() *infraCRM.Client {
	return s.client
}

// Healthcheck checks API connectivity
func (s *Service) Healthcheck(ctx context.Context) error {
	return s.client.Healthcheck(ctx)
}

// GetAccountInfo returns account information as formatted string
func (s *Service) GetAccountInfo(ctx context.Context) (string, error) {
	account, err := s.client.SDK().Account().GetCurrent(ctx, nil)
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
func (s *Service) GetPipelines(ctx context.Context) (string, error) {
	pipelines, _, err := s.client.SDK().Pipelines().Get(ctx, nil)
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
		statuses, _, err := s.client.SDK().Statuses(p.ID).Get(ctx, nil)
		if err != nil {
			result += fmt.Sprintf("   ⚠️ Ошибка загрузки статусов: %v\n", err)
			continue
		}

		for i, st := range statuses {
			prefix := "├─"
			if i == len(statuses)-1 {
				prefix = "└─"
			}
			result += fmt.Sprintf("   %s %s (ID: %d)\n", prefix, st.Name, st.ID)
		}
		result += "\n"
	}

	return result, nil
}
