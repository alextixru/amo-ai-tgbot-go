package gkit

import (
	"context"
	"fmt"

	amocrm "github.com/alextixru/amocrm-sdk-go"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/tihn/amo-ai-tgbot-go/app/gkit/tools"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/activities"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/admin_integrations"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/admin_pipelines"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/admin_schema"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/admin_users"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/catalogs"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/complex_create"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/customers"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/entities"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/files"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/products"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/unsorted"
	genkitClient "github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/genkit"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/flows"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/session"
)

// Agent handles AI processing with Genkit flows
type Agent struct {
	g        *genkit.Genkit
	model    ai.Model
	chatFlow func(context.Context, flows.ChatInput) (flows.ChatOutput, error)
}

// NewAgent creates a new AI agent with registered flows and tools.
// Returns an error if any service that requires preloading fails to initialise.
func NewAgent(ctx context.Context, client *genkitClient.Client, sdk *amocrm.SDK) (*Agent, error) {
	g := client.G
	model := client.Model

	// Создаём session store для истории диалогов
	store := session.NewMemoryStore()

	// Инициализируем сервисы
	entitiesSvc, err := entities.New(ctx, sdk)
	if err != nil {
		return nil, fmt.Errorf("NewAgent: %w", err)
	}
	activitiesSvc, err := activities.New(ctx, sdk)
	if err != nil {
		return nil, fmt.Errorf("NewAgent: %w", err)
	}

	complexCreateSvc, err := complex_create.New(ctx, sdk)
	if err != nil {
		return nil, fmt.Errorf("NewAgent: %w", err)
	}

	productsSvc := products.NewService(sdk)

	catalogsSvc, err := catalogs.New(ctx, sdk)
	if err != nil {
		return nil, fmt.Errorf("NewAgent: %w", err)
	}

	filesSvc := files.NewService(sdk)
	unsortedSvc, err := unsorted.New(ctx, sdk)
	if err != nil {
		return nil, fmt.Errorf("NewAgent: %w", err)
	}
	customersSvc, err := customers.New(ctx, sdk)
	if err != nil {
		return nil, fmt.Errorf("NewAgent: %w", err)
	}
	adminSchemaSvc := admin_schema.NewService(sdk)
	adminPipelinesSvc := admin_pipelines.New(sdk)
	adminUsersSvc := admin_users.NewService(sdk)
	adminIntegrationsSvc := admin_integrations.NewService(sdk)

	// Регистрируем все tools через новый транспортный слой
	registry := tools.NewRegistry(
		g,
		entitiesSvc,
		activitiesSvc,
		complexCreateSvc,
		productsSvc,
		catalogsSvc,
		filesSvc,
		unsortedSvc,
		customersSvc,
		adminSchemaSvc,
		adminPipelinesSvc,
		adminUsersSvc,
		adminIntegrationsSvc,
	)
	registry.RegisterAll()

	// Регистрируем Chat Flow с tools и session store
	chatRunner := flows.DefineChatFlow(g, model, registry.AllTools(), store)

	return &Agent{
		g:        g,
		model:    model,
		chatFlow: chatRunner,
	}, nil
}

// Process processes a user message using the chat flow with user context
// sessionID should be unique per conversation (e.g., Telegram chat ID)
func (a *Agent) Process(ctx context.Context, sessionID, message string) (string, error) {
	output, err := a.chatFlow(ctx, flows.ChatInput{
		SessionID: sessionID,
		Message:   message,
	})
	if err != nil {
		return "", err
	}
	return output.Response, nil
}
