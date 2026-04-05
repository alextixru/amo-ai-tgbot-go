package tools

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/tool"

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

// CRMToolset implements tool.Toolset — returns all CRM tools for ADK agent.
type CRMToolset struct {
	tools []tool.Tool
}

// NewCRMToolset creates a toolset with all 12 CRM tools.
func NewCRMToolset(
	entitiesSvc entities.Service,
	activitiesSvc activities.Service,
	complexCreateSvc complex_create.Service,
	productsSvc products.Service,
	catalogsSvc catalogs.Service,
	filesSvc files.Service,
	unsortedSvc unsorted.Service,
	customersSvc customers.Service,
	adminSchemaSvc admin_schema.Service,
	adminPipelinesSvc admin_pipelines.Service,
	adminUsersSvc admin_users.Service,
	adminIntegrationsSvc admin_integrations.Service,
) *CRMToolset {
	return &CRMToolset{
		tools: []tool.Tool{
			NewEntitiesTool(entitiesSvc),
			NewActivitiesTool(activitiesSvc),
			NewComplexCreateTool(complexCreateSvc),
			NewProductsTool(productsSvc),
			NewCatalogsTool(catalogsSvc),
			NewFilesTool(filesSvc),
			NewUnsortedTool(unsortedSvc),
			NewCustomersTool(customersSvc),
			NewAdminSchemaTool(adminSchemaSvc),
			NewAdminPipelinesTool(adminPipelinesSvc),
			NewAdminUsersTool(adminUsersSvc),
			NewAdminIntegrationsTool(adminIntegrationsSvc),
		},
	}
}

// Name implements tool.Toolset.
func (ts *CRMToolset) Name() string {
	return "crm_tools"
}

// Tools implements tool.Toolset.
func (ts *CRMToolset) Tools(_ agent.ReadonlyContext) ([]tool.Tool, error) {
	return ts.tools, nil
}
