package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm"
	"github.com/tihn/amo-ai-tgbot-go/app/gkit"
	appctx "github.com/tihn/amo-ai-tgbot-go/internal/services/context"
)

// Handler processes Telegram messages
type Handler struct {
	agent      *gkit.Agent
	crm        *crm.Service
	ctxBuilder *appctx.Builder
	debug      bool
}

// NewHandler creates a new Handler with AI agent and CRM service
func NewHandler(agent *gkit.Agent, crmService *crm.Service, debug bool) *Handler {
	return &Handler{
		agent:      agent,
		crm:        crmService,
		ctxBuilder: appctx.NewBuilder(crmService.Client().SDK()),
		debug:      debug,
	}
}

func (h *Handler) debugLog(format string, v ...any) {
	if h.debug {
		log.Printf(format, v...)
	}
}

// HandleMessage handles incoming text messages
func (h *Handler) HandleMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	text := update.Message.Text
	chatID := update.Message.Chat.ID
	telegramUserID := update.Message.From.ID

	h.debugLog("📨 Received message: %q from chat %d", text, chatID)

	var response string
	var err error

	// Handle commands
	switch {
	case text == "/start":
		response = h.handleStart()
	case text == "/status" || text == "/healthcheck":
		response = h.handleHealthcheck(ctx)
	case text == "/account":
		response = h.handleAccount(ctx)
	case text == "/pipelines":
		response = h.handlePipelines(ctx)
	case strings.HasPrefix(text, "/"):
		response = "❓ Неизвестная команда. Используй /start для списка команд."
	default:
		// Build user context
		userCtx := h.ctxBuilder.MustBuild(ctx, telegramUserID)

		// Process with AI (chatID as sessionID for history)
		sessionID := fmt.Sprintf("tg_%d", chatID)
		h.debugLog("🤖 Sending to Ollama...")
		response, err = h.agent.Process(ctx, sessionID, text, userCtx.ToMap())
		if err != nil {
			log.Printf("AI error: %v", err)
			response = fmt.Sprintf("❌ Ошибка AI: %v", err)
		} else {
			h.debugLog("🤖 Ollama response received")
		}
	}

	h.debugLog("📤 Sending response (%d chars)...", len(response))
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      response,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Printf("❌ SendMessage error: %v", err)
	} else {
		h.debugLog("✅ Response sent")
	}
}

func (h *Handler) handleStart() string {
	return `👋 Привет! Я amoCRM AI бот.

📋 Доступные команды:
• /status — проверить подключение к amoCRM
• /account — информация об аккаунте
• /pipelines — список воронок и статусов

💬 Или просто напиши мне что-нибудь — я отвечу через AI!`
}

func (h *Handler) handleHealthcheck(ctx context.Context) string {
	err := h.crm.Healthcheck(ctx)
	if err != nil {
		return fmt.Sprintf("❌ amoCRM недоступен\n\nОшибка: %v", err)
	}
	return "✅ amoCRM доступен!"
}

func (h *Handler) handleAccount(ctx context.Context) string {
	info, err := h.crm.GetAccountInfo(ctx)
	if err != nil {
		return fmt.Sprintf("❌ Ошибка получения аккаунта\n\n%v", err)
	}
	return info
}

func (h *Handler) handlePipelines(ctx context.Context) string {
	info, err := h.crm.GetPipelines(ctx)
	if err != nil {
		return fmt.Sprintf("❌ Ошибка получения воронок\n\n%v", err)
	}
	return info
}
