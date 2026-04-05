package agent

import (
	"context"
	"fmt"

	amocrm "github.com/alextixru/amocrm-sdk-go"

	"github.com/tihn/amo-ai-tgbot-go/app/agent/tools"
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
	"github.com/tihn/amo-ai-tgbot-go/internal/services/session"
)

// Agent handles AI processing. Currently a stub — will be replaced by ADK Runner.
type Agent struct {
	registry *tools.Registry
	store    session.Store
}

// NewAgent creates a new AI agent with registered tools.
// CRM services are initialized and tools are registered, ready for ADK Runner.
func NewAgent(ctx context.Context, sdk *amocrm.SDK) (*Agent, error) {
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

	// Регистрируем все tools через framework-agnostic registry
	registry := tools.NewRegistry(
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

	return &Agent{
		registry: registry,
		store:    store,
	}, nil
}

// Process processes a user message.
// TODO: Replace stub with ADK Runner integration.
func (a *Agent) Process(ctx context.Context, sessionID, message string) (string, error) {
	return "🔧 AI-агент в процессе миграции на ADK. Скоро вернусь!", nil
}

// Tools returns the registered tool definitions (for ADK adapter).
func (a *Agent) Tools() []tools.ToolDefinition {
	return a.registry.AllTools()
}

// Store returns the session store (for ADK adapter).
func (a *Agent) Store() session.Store {
	return a.store
}
