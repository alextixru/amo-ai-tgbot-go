package tools

import (
	"context"

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
)

// ToolDefinition describes a single tool in a framework-agnostic way.
// ADK adapter will convert these into ADK tool format.
type ToolDefinition struct {
	Name        string
	Description string
	InputSchema map[string]any                                    // JSON Schema for the tool input
	Handler     func(ctx context.Context, input any) (any, error) // Business logic handler
}

// Registry holds CRM services and produces framework-agnostic tool definitions.
type Registry struct {
	entitiesService          entities.Service
	activitiesService        activities.Service
	complexCreateService     complex_create.Service
	productsService          products.Service
	catalogsService          catalogs.Service
	filesService             files.Service
	unsortedService          unsorted.Service
	customersService         customers.Service
	adminSchemaService       admin_schema.Service
	adminPipelinesService    admin_pipelines.Service
	adminUsersService        admin_users.Service
	adminIntegrationsService admin_integrations.Service

	tools []ToolDefinition
}

func NewRegistry(
	entitiesService entities.Service,
	activitiesService activities.Service,
	complexCreateService complex_create.Service,
	productsService products.Service,
	catalogsService catalogs.Service,
	filesService files.Service,
	unsortedService unsorted.Service,
	customersService customers.Service,
	adminSchemaService admin_schema.Service,
	adminPipelinesService admin_pipelines.Service,
	adminUsersService admin_users.Service,
	adminIntegrationsService admin_integrations.Service,
) *Registry {
	return &Registry{
		entitiesService:          entitiesService,
		activitiesService:        activitiesService,
		complexCreateService:     complexCreateService,
		productsService:          productsService,
		catalogsService:          catalogsService,
		filesService:             filesService,
		unsortedService:          unsortedService,
		customersService:         customersService,
		adminSchemaService:       adminSchemaService,
		adminPipelinesService:    adminPipelinesService,
		adminUsersService:        adminUsersService,
		adminIntegrationsService: adminIntegrationsService,
		tools:                    make([]ToolDefinition, 0),
	}
}

func (r *Registry) RegisterAll() {
	r.RegisterEntitiesTool()
	r.RegisterActivitiesTool()
	r.RegisterComplexCreateTool()
	r.RegisterProductsTool()
	r.RegisterCatalogsTool()
	r.RegisterFilesTool()
	r.RegisterUnsortedTool()
	r.RegisterCustomersTool()
	r.RegisterAdminSchemaTool()
	r.RegisterAdminPipelinesTool()
	r.RegisterAdminUsersTool()
	r.RegisterAdminIntegrationsTool()
}

// AllTools returns all registered tool definitions.
func (r *Registry) AllTools() []ToolDefinition {
	return r.tools
}

// addTool adds a tool definition to the registry.
func (r *Registry) addTool(def ToolDefinition) {
	r.tools = append(r.tools, def)
}
