package gkit

import (
	"context"

	amocrm "github.com/alextixru/amocrm-sdk-go"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/tihn/amo-ai-tgbot-go/app/gkit/flows"
	"github.com/tihn/amo-ai-tgbot-go/app/gkit/session"
	"github.com/tihn/amo-ai-tgbot-go/app/gkit/tools"
	genkitClient "github.com/tihn/amo-ai-tgbot-go/infrastructure/genkit"
	"github.com/tihn/amo-ai-tgbot-go/services/activities"
	"github.com/tihn/amo-ai-tgbot-go/services/admin_integrations"
	"github.com/tihn/amo-ai-tgbot-go/services/admin_pipelines"
	"github.com/tihn/amo-ai-tgbot-go/services/admin_schema"
	"github.com/tihn/amo-ai-tgbot-go/services/admin_users"
	"github.com/tihn/amo-ai-tgbot-go/services/catalogs"
	"github.com/tihn/amo-ai-tgbot-go/services/complex_create"
	"github.com/tihn/amo-ai-tgbot-go/services/customers"
	"github.com/tihn/amo-ai-tgbot-go/services/entities"
	"github.com/tihn/amo-ai-tgbot-go/services/files"
	"github.com/tihn/amo-ai-tgbot-go/services/products"
	"github.com/tihn/amo-ai-tgbot-go/services/unsorted"
)

// Agent handles AI processing with Genkit flows
type Agent struct {
	g        *genkit.Genkit
	model    ai.Model
	chatFlow func(context.Context, flows.ChatInput) (flows.ChatOutput, error)
}

// NewAgent creates a new AI agent with registered flows and tools
func NewAgent(client *genkitClient.Client, sdk *amocrm.SDK) *Agent {
	g := client.G
	model := client.Model

	// Создаём session store для истории диалогов
	store := session.NewMemoryStore()

	// Инициализируем сервисы
	entitiesSvc := entities.New(sdk)
	activitiesSvc := activities.New(sdk)
	complexCreateSvc := complex_create.NewService(sdk)
	productsSvc := products.NewService(sdk)
	catalogsSvc := catalogs.NewService(sdk)
	filesSvc := files.NewService(sdk)
	unsortedSvc := unsorted.NewService(sdk)
	customersSvc := customers.NewService(sdk)
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
	}
}

// Process processes a user message using the chat flow with user context
// sessionID should be unique per conversation (e.g., Telegram chat ID)
func (a *Agent) Process(ctx context.Context, sessionID, message string, userContext map[string]any) (string, error) {
	output, err := a.chatFlow(ctx, flows.ChatInput{
		SessionID:   sessionID,
		Message:     message,
		UserContext: userContext,
	})
	if err != nil {
		return "", err
	}
	return output.Response, nil
}
