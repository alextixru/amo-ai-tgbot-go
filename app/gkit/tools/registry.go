// Package tools provides Genkit tool definitions for amoCRM SDK operations.
// These tools enable LLM to interact with amoCRM through structured tool calling.
package tools

import (
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	amocrm "github.com/alextixru/amocrm-sdk-go"
)

// Registry holds all registered Genkit tools and provides access to SDK.
type Registry struct {
	g     *genkit.Genkit
	sdk   *amocrm.SDK
	tools []ai.Tool
}

// NewRegistry creates a new tool registry.
func NewRegistry(g *genkit.Genkit, sdk *amocrm.SDK) *Registry {
	return &Registry{
		g:     g,
		sdk:   sdk,
		tools: make([]ai.Tool, 0),
	}
}

// RegisterAll registers all CRM tools and returns the registry.
// TODO: Register 12 unified tools (entities, activities, complex_create, etc.)
func (r *Registry) RegisterAll() *Registry {
	r.registerEntitiesTool()
	r.registerActivitiesTool()
	r.registerComplexCreateTool()
	r.registerProductsTool()
	r.registerCatalogsTool()
	r.registerFilesTool()
	r.registerUnsortedTool()
	r.registerCustomersTool()
	r.registerAdminSchemaTool()
	r.registerAdminPipelinesTool()
	r.registerAdminUsersTool()
	r.registerAdminIntegrationsTool()
	return r
}

// AllTools returns all registered tools for use with ai.WithTools().
func (r *Registry) AllTools() []ai.Tool {
	return r.tools
}

// addTool adds a tool to the registry.
func (r *Registry) addTool(tool ai.Tool) {
	r.tools = append(r.tools, tool)
}

// G returns the Genkit instance.
func (r *Registry) G() *genkit.Genkit {
	return r.g
}

// SDK returns the amoCRM SDK instance.
func (r *Registry) SDK() *amocrm.SDK {
	return r.sdk
}
