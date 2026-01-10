package tools

import (
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
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

type Registry struct {
	g *genkit.Genkit

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

	tools []ai.ToolRef
}

func NewRegistry(
	g *genkit.Genkit,
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
		g:                        g,
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
		tools:                    make([]ai.ToolRef, 0),
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

// AllTools возвращает все зарегистрированные инструменты
func (r *Registry) AllTools() []ai.ToolRef {
	return r.tools
}

// addTool добавляет инструмент в список
func (r *Registry) addTool(tool ai.Tool) {
	r.tools = append(r.tools, tool)
}
