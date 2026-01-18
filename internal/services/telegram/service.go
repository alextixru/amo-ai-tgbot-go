package telegram

import (
	"context"
	"fmt"

	"github.com/tihn/amo-ai-tgbot-go/app/gkit"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm"
)

// Service handles Telegram business logic
type Service struct {
	agent *gkit.Agent
	crm   *crm.Service
}

// NewService creates a new Telegram service
func NewService(agent *gkit.Agent, crmService *crm.Service) *Service {
	return &Service{
		agent: agent,
		crm:   crmService,
	}
}

// HandleStart returns the start message
func (s *Service) HandleStart() string {
	return `👋 Привет! Я amoCRM AI бот.

📋 Доступные команды:
• /status — проверить подключение к amoCRM
• /account — информация об аккаунте
• /pipelines — список воронок и статусов

💬 Или просто напиши мне что-нибудь — я отвечу через AI!`
}

// HandleHealthcheck checks CRM connectivity
func (s *Service) HandleHealthcheck(ctx context.Context) string {
	err := s.crm.Healthcheck(ctx)
	if err != nil {
		return fmt.Sprintf("❌ amoCRM недоступен\n\nОшибка: %v", err)
	}
	return "✅ amoCRM доступен!"
}

// HandleAccount returns account information
func (s *Service) HandleAccount(ctx context.Context) string {
	info, err := s.crm.GetAccountInfo(ctx)
	if err != nil {
		return fmt.Sprintf("❌ Ошибка получения аккаунта\n\n%v", err)
	}
	return info
}

// HandlePipelines returns pipelines information
func (s *Service) HandlePipelines(ctx context.Context) string {
	info, err := s.crm.GetPipelines(ctx)
	if err != nil {
		return fmt.Sprintf("❌ Ошибка получения воронок\n\n%v", err)
	}
	return info
}

// ProcessAI processes a message through the AI agent
func (s *Service) ProcessAI(ctx context.Context, telegramUserID int64, chatID int64, text string) (string, error) {
	// Process with AI (chatID as sessionID for history)
	sessionID := fmt.Sprintf("tg_%d", chatID)
	return s.agent.Process(ctx, sessionID, text)
}
