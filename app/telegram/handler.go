package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

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
	var keyboard *models.InlineKeyboardMarkup
	var err error

	// Handle commands
	switch {
	case text == "/start":
		response, keyboard = h.svc.HandleStart(telegramUserID)
	case text == "/connect":
		response, keyboard = h.svc.HandleConnect(telegramUserID, chatID)
	case strings.HasPrefix(text, "/auth "):
		code := strings.TrimSpace(strings.TrimPrefix(text, "/auth "))
		response, keyboard = h.svc.HandleAuthCode(ctx, telegramUserID, code)
	case text == "/me":
		response = h.svc.HandleMe(telegramUserID)
	case text == "/disconnect":
		response = h.svc.HandleDisconnect(telegramUserID)
	case text == "/status" || text == "/healthcheck":
		response = h.svc.HandleHealthcheck(ctx)
	case text == "/account":
		response = h.svc.HandleAccount(ctx)
	case text == "/pipelines":
		response = h.svc.HandlePipelines(ctx)
	case text != "" && text[0] == '/':
		response = "❓ Неизвестная команда. Используй /start для списка команд."
	default:
		// Check if user is waiting for auth code
		if h.svc.IsWaitingCode(telegramUserID) {
			h.debugLog("🔐 User is waiting for auth code, processing as code...")
			response, keyboard = h.svc.HandleAuthCode(ctx, telegramUserID, strings.TrimSpace(text))
		} else {
			h.debugLog("🤖 Processing with AI...")
			response, err = h.svc.ProcessAI(ctx, telegramUserID, chatID, text)
			if err != nil {
				log.Printf("AI error: %v", err)
				response = fmt.Sprintf("❌ Ошибка AI: %v", err)
			} else {
				h.debugLog("🤖 AI response received")
			}
		}
	}

	h.sendResponse(ctx, b, chatID, response, keyboard)
}

// HandleCallback handles inline button callbacks
func (h *Handler) HandleCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	messageID := update.CallbackQuery.Message.Message.ID
	telegramUserID := update.CallbackQuery.From.ID
	data := update.CallbackQuery.Data

	h.debugLog("📨 Received callback: %q from user %d", data, telegramUserID)

	var response string
	var keyboard *models.InlineKeyboardMarkup

	switch data {
	case "auth_start":
		response, keyboard = h.svc.ShowAuthWaiting(telegramUserID, chatID)
	case "auth_panel":
		response, keyboard = h.svc.ShowAuthPanel(telegramUserID)
	case "auth_cancel":
		response, keyboard = h.svc.CancelAuth(telegramUserID)
	case "auth_disconnect":
		response, keyboard = h.svc.Disconnect(telegramUserID)
	case "back_main":
		response, keyboard = h.svc.HandleStart(telegramUserID)
	default:
		response = "❓ Неизвестное действие."
	}

	// Answer callback to remove "loading" state
	_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	// Edit the existing message instead of sending a new one
	h.editMessage(ctx, b, chatID, messageID, response, keyboard)
}

func (h *Handler) sendResponse(ctx context.Context, b *bot.Bot, chatID int64, text string, keyboard *models.InlineKeyboardMarkup) {
	h.debugLog("📤 Sending response (%d chars)...", len(text))

	params := &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	}

	if keyboard != nil {
		params.ReplyMarkup = keyboard
	}

	_, err := b.SendMessage(ctx, params)
	if err != nil {
		log.Printf("❌ SendMessage error: %v", err)
	} else {
		h.debugLog("✅ Response sent")
	}
}

func (h *Handler) editMessage(ctx context.Context, b *bot.Bot, chatID int64, messageID int, text string, keyboard *models.InlineKeyboardMarkup) {
	h.debugLog("📝 Editing message %d...", messageID)

	params := &bot.EditMessageTextParams{
		ChatID:    chatID,
		MessageID: messageID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	}

	if keyboard != nil {
		params.ReplyMarkup = keyboard
	}

	_, err := b.EditMessageText(ctx, params)
	if err != nil {
		log.Printf("❌ EditMessageText error: %v", err)
		// Fallback to sending new message
		h.sendResponse(ctx, b, chatID, text, keyboard)
	} else {
		h.debugLog("✅ Message edited")
	}
}
