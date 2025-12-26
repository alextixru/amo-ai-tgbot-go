package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"

	appCRM "github.com/tihn/amo-ai-tgbot-go/app/crm"
	"github.com/tihn/amo-ai-tgbot-go/app/gkit"
	tgHandler "github.com/tihn/amo-ai-tgbot-go/app/telegram"
	"github.com/tihn/amo-ai-tgbot-go/infrastructure/config"
	"github.com/tihn/amo-ai-tgbot-go/infrastructure/crm"
	"github.com/tihn/amo-ai-tgbot-go/infrastructure/genkit"
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

	// AI agent
	agent := gkit.NewAgent(genkitClient)

	// CRM service (business logic)
	crmService := appCRM.NewService(crmClient)

	// Telegram handler
	handler := tgHandler.NewHandler(agent, crmService, cfg.Debug)

	// === Start Bot ===

	opts := []bot.Option{
		bot.WithDefaultHandler(handler.HandleMessage),
		bot.WithSkipGetMe(),
	}

	b, err := bot.New(cfg.TelegramToken, opts...)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bot started with Ollama model: %s", cfg.OllamaModel)
	b.Start(ctx)
}
