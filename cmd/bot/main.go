package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"

	"github.com/tihn/amo-ai-tgbot-go/app/agent"
	tgHandler "github.com/tihn/amo-ai-tgbot-go/app/telegram"
	"github.com/tihn/amo-ai-tgbot-go/config"
	"github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/crm"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/auth"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/telegram"
)

func init() {
	_ = godotenv.Load() // Загружаем .env если есть
}

func main() {
	cfg := config.Load()

	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
	}

	// Debug mode: prompt for missing amoCRM credentials
	if err := config.PromptMissingCredentials(cfg); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// === Infrastructure ===

	// CRM client
	crmClient, err := crm.New(cfg)
	if err != nil {
		log.Fatalf("Failed to init CRM client: %v", err)
	}

	// === Auth Service ===

	// Token storage directory (using ~/.gemini for consistency with gemini-cli)
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	tokenDir := filepath.Join(home, ".gemini")

	// Create auth service with in-memory state store and file-based token store
	stateStore := auth.NewMemoryStateStore()
	tokenStore := auth.NewFileTokenStore(tokenDir)
	authService := auth.NewService(stateStore, tokenStore)

	// === Application ===

	// AI agent (needs SDK for tools)
	aiAgent, err := agent.NewAgent(ctx, crmClient.SDK())
	if err != nil {
		log.Fatalf("Failed to init AI agent: %v", err)
	}

	// Telegram service (business logic)
	telegramSvc := telegram.NewService(aiAgent, crmClient, authService)

	// Telegram handler
	handler := tgHandler.NewHandler(telegramSvc, cfg.Debug)

	// === Start Bot ===

	opts := []bot.Option{
		bot.WithDefaultHandler(handler.HandleMessage),
		bot.WithSkipGetMe(),
	}

	b, err := bot.New(cfg.TelegramToken, opts...)
	if err != nil {
		log.Fatal(err)
	}

	// Register callback handler for inline buttons
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "", bot.MatchTypePrefix, handler.HandleCallback)

	log.Print("Bot started (AI agent: stub, pending ADK migration)")
	b.Start(ctx)
}
