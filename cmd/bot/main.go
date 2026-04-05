package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"

	appagent "github.com/tihn/amo-ai-tgbot-go/app/agent"
	"github.com/tihn/amo-ai-tgbot-go/app/agent/tools"
	tgHandler "github.com/tihn/amo-ai-tgbot-go/app/telegram"
	"github.com/tihn/amo-ai-tgbot-go/config"
	"github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/crm"
	"github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/llm"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/auth"
	crmActivities "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/activities"
	crmAdminIntegrations "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/admin_integrations"
	crmAdminPipelines "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/admin_pipelines"
	crmAdminSchema "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/admin_schema"
	crmAdminUsers "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/admin_users"
	crmCatalogs "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/catalogs"
	crmComplexCreate "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/complex_create"
	crmCustomers "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/customers"
	crmEntities "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/entities"
	crmFiles "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/files"
	crmProducts "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/products"
	crmUnsorted "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/unsorted"
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

	// LLM provider (Ollama via OpenAI-compatible API)
	llmModel := llm.NewProvider(cfg)

	// === CRM Services ===
	sdk := crmClient.SDK()

	entitiesSvc, err := crmEntities.New(ctx, sdk)
	if err != nil {
		log.Fatalf("Failed to init entities service: %v", err)
	}
	activitiesSvc, err := crmActivities.New(ctx, sdk)
	if err != nil {
		log.Fatalf("Failed to init activities service: %v", err)
	}
	complexCreateSvc, err := crmComplexCreate.New(ctx, sdk)
	if err != nil {
		log.Fatalf("Failed to init complex_create service: %v", err)
	}
	catalogsSvc, err := crmCatalogs.New(ctx, sdk)
	if err != nil {
		log.Fatalf("Failed to init catalogs service: %v", err)
	}
	unsortedSvc, err := crmUnsorted.New(ctx, sdk)
	if err != nil {
		log.Fatalf("Failed to init unsorted service: %v", err)
	}
	customersSvc, err := crmCustomers.New(ctx, sdk)
	if err != nil {
		log.Fatalf("Failed to init customers service: %v", err)
	}

	productsSvc := crmProducts.NewService(sdk)
	filesSvc := crmFiles.NewService(sdk)
	adminSchemaSvc := crmAdminSchema.NewService(sdk)
	adminPipelinesSvc := crmAdminPipelines.New(sdk)
	adminUsersSvc := crmAdminUsers.NewService(sdk)
	adminIntegrationsSvc := crmAdminIntegrations.NewService(sdk)

	// CRM Toolset for ADK agent
	crmToolset := tools.NewCRMToolset(
		entitiesSvc, activitiesSvc, complexCreateSvc, productsSvc,
		catalogsSvc, filesSvc, unsortedSvc, customersSvc,
		adminSchemaSvc, adminPipelinesSvc, adminUsersSvc, adminIntegrationsSvc,
	)

	// AI agent with CRM tools
	aiAgent, err := appagent.NewAgent(ctx, llmModel, crmToolset)
	if err != nil {
		log.Fatalf("Failed to init AI agent: %v", err)
	}

	// === ADK Web UI (debug) ===

	adkLauncher := full.NewLauncher()
	launcherConfig := &launcher.Config{
		SessionService: aiAgent.SessionService(),
		AgentLoader:    agent.NewSingleLoader(aiAgent.ADKAgent()),
	}

	go func() {
		// "web api webui" — запускает HTTP сервер с REST API и Web UI на :8080
		if err := adkLauncher.Execute(ctx, launcherConfig, []string{"web", "api", "webui"}); err != nil {
			log.Printf("ADK Web UI error: %v", err)
		}
	}()

	// === Telegram Bot ===

	// Telegram service (business logic)
	telegramSvc := telegram.NewService(aiAgent, crmClient, authService)

	// Telegram handler
	handler := tgHandler.NewHandler(telegramSvc, cfg.Debug)

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

	log.Print("Bot started (AI agent: ADK Runner, Web UI: http://localhost:8080)")
	b.Start(ctx)
}
