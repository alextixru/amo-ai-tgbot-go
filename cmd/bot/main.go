package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"

	"github.com/tihn/amo-ai-tgbot-go/app/gkit"
	tgHandler "github.com/tihn/amo-ai-tgbot-go/app/telegram"
	"github.com/tihn/amo-ai-tgbot-go/config"
	"github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/crm"
	"github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/genkit"
	appCRM "github.com/tihn/amo-ai-tgbot-go/internal/services/crm"
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

	// Genkit client
	genkitClient, err := genkit.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to init Genkit client: %v", err)
	}

	// CRM client
	crmClient, err := crm.New(cfg)
	if err != nil {
		log.Fatalf("Failed to init CRM client: %v", err)
	}

	// === Application ===

	// CRM service (business logic)
	crmService := appCRM.NewService(crmClient)

	// AI agent (needs SDK for tools)
	agent := gkit.NewAgent(genkitClient, crmClient.SDK())

	// Telegram service (business logic)
	telegramSvc := telegram.NewService(agent, crmService)

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

	if cfg.AIProvider == "gemini-cli" {
		log.Print("Bot started with Gemini CLI provider")
	} else {
		log.Printf("Bot started with Ollama model: %s", cfg.OllamaModel)
	}
	b.Start(ctx)
}
