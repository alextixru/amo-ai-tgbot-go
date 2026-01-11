package telegram

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	tgsvc "github.com/tihn/amo-ai-tgbot-go/internal/services/telegram"
)

// Handler processes Telegram messages
type Handler struct {
	svc   *tgsvc.Service
	debug bool
}

// NewHandler creates a new Handler with Telegram service
func NewHandler(svc *tgsvc.Service, debug bool) *Handler {
	return &Handler{
		svc:   svc,
		debug: debug,
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
		response = h.svc.HandleStart()
	case text == "/status" || text == "/healthcheck":
		response = h.svc.HandleHealthcheck(ctx)
	case text == "/account":
		response = h.svc.HandleAccount(ctx)
	case text == "/pipelines":
		response = h.svc.HandlePipelines(ctx)
	case text != "" && text[0] == '/':
		response = "❓ Неизвестная команда. Используй /start для списка команд."
	default:
		h.debugLog("🤖 Processing with AI...")
		response, err = h.svc.ProcessAI(ctx, telegramUserID, chatID, text)
		if err != nil {
			log.Printf("AI error: %v", err)
			response = fmt.Sprintf("❌ Ошибка AI: %v", err)
		} else {
			h.debugLog("🤖 AI response received")
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
