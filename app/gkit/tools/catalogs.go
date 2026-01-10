package tools

import (
	"encoding/json"
	"fmt"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func (r *Registry) RegisterCatalogsTool() {
	r.addTool(genkit.DefineTool[gkitmodels.CatalogsInput, any](
		r.g,
		"catalogs",
		"Work with catalogs and their elements",
		func(ctx *ai.ToolContext, input gkitmodels.CatalogsInput) (any, error) {
			switch input.Action {
			// Catalogs
			case "list":
				return r.catalogsService.ListCatalogs(ctx)
			case "get":
				if input.CatalogID == 0 {
					return nil, fmt.Errorf("catalog_id is required")
				}
				return r.catalogsService.GetCatalog(ctx, input.CatalogID)
			case "create":
				var catalogs []*amomodels.Catalog
				data, _ := json.Marshal(input.Data) // Assuming input.Data maps to Catalog models
				if err := json.Unmarshal(data, &catalogs); err != nil {
					return nil, fmt.Errorf("failed to parse catalogs data: %w", err)
				}
				return r.catalogsService.CreateCatalogs(ctx, catalogs)
			case "update":
				var catalogs []*amomodels.Catalog
				data, _ := json.Marshal(input.Data)
				if err := json.Unmarshal(data, &catalogs); err != nil {
					return nil, fmt.Errorf("failed to parse catalogs data: %w", err)
				}
				return r.catalogsService.UpdateCatalogs(ctx, catalogs)

			// Elements
			case "list_elements":
				if input.CatalogID == 0 {
					return nil, fmt.Errorf("catalog_id is required")
				}
				return r.catalogsService.ListElements(ctx, input.CatalogID)
			case "get_element":
				if input.CatalogID == 0 || input.ElementID == 0 {
					return nil, fmt.Errorf("catalog_id and element_id are required")
				}
				return r.catalogsService.GetElement(ctx, input.CatalogID, input.ElementID)
			case "create_element":
				if input.CatalogID == 0 {
					return nil, fmt.Errorf("catalog_id is required")
				}
				var elements []*amomodels.CatalogElement
				data, _ := json.Marshal(input.ElementData)
				if err := json.Unmarshal(data, &elements); err != nil {
					return nil, fmt.Errorf("failed to parse elements data: %w", err)
				}
				return r.catalogsService.CreateElements(ctx, input.CatalogID, elements)
			case "update_element":
				if input.CatalogID == 0 {
					return nil, fmt.Errorf("catalog_id is required")
				}
				var elements []*amomodels.CatalogElement
				data, _ := json.Marshal(input.ElementData)
				if err := json.Unmarshal(data, &elements); err != nil {
					return nil, fmt.Errorf("failed to parse elements data: %w", err)
				}
				return r.catalogsService.UpdateElements(ctx, input.CatalogID, elements)

			default:
				return nil, fmt.Errorf("unknown action: %s", input.Action)
			}
		},
	))
}
