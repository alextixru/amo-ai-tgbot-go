package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models"
)

func (r *Registry) RegisterProductsTool() {
	r.addTool(genkit.DefineTool[gkitmodels.ProductsInput, any](
		r.g,
		"products",
		"Работа с товарами (элементами каталога 'products'). "+
			"Поддерживает: search (поиск), get (получение), create (создание), update (обновление), delete (удаление), "+
			"get_by_entity (товары в сделке/контакте), link (привзяка к сущности), unlink (отвязка), "+
			"update_quantity (обновление количества).",
		func(ctx *ai.ToolContext, input gkitmodels.ProductsInput) (any, error) {
			return r.handleProducts(ctx.Context, input)
		},
	))
}

func (r *Registry) handleProducts(ctx context.Context, input gkitmodels.ProductsInput) (any, error) {
	switch input.Action {
	case "search":
		return r.productsService.SearchProducts(ctx, input.Filter)
	case "get":
		if input.ProductID == 0 {
			return nil, fmt.Errorf("product_id is required for action 'get'")
		}
		var with []string
		if input.Filter != nil {
			with = input.Filter.With
		}
		return r.productsService.GetProduct(ctx, input.ProductID, with)
	case "create":
		if input.Data == nil && len(input.Items) == 0 {
			return nil, fmt.Errorf("data or items is required for action 'create'")
		}
		var elements []*models.CatalogElement

		// Поддержка одиночного элемента или батча
		var dataToMarshal any
		if len(input.Items) > 0 {
			dataToMarshal = input.Items
		} else {
			dataToMarshal = []*gkitmodels.ProductData{input.Data}
		}

		elementData, _ := json.Marshal(dataToMarshal)
		if err := json.Unmarshal(elementData, &elements); err != nil {
			return nil, fmt.Errorf("failed to parse elements data: %w", err)
		}
		return r.productsService.CreateProducts(ctx, elements)
	case "update":
		if (input.ProductID == 0 && input.Data != nil) || (input.Data == nil && len(input.Items) == 0) {
			return nil, fmt.Errorf("product_id and data (or items for batch) are required for action 'update'")
		}

		var elements []*models.CatalogElement
		if len(input.Items) > 0 {
			elementData, _ := json.Marshal(input.Items)
			if err := json.Unmarshal(elementData, &elements); err != nil {
				return nil, fmt.Errorf("failed to parse elements data: %w", err)
			}
		} else {
			elementData, _ := json.Marshal(input.Data)
			var element *models.CatalogElement
			if err := json.Unmarshal(elementData, &element); err != nil {
				return nil, fmt.Errorf("failed to parse element data: %w", err)
			}
			element.ID = input.ProductID
			elements = append(elements, element)
		}
		return r.productsService.UpdateProducts(ctx, elements)
	case "delete":
		if len(input.IDs) == 0 {
			return nil, fmt.Errorf("ids array is required for action 'delete'")
		}
		return nil, r.productsService.DeleteProducts(ctx, input.IDs)
	case "get_by_entity":
		if input.Entity == nil {
			return nil, fmt.Errorf("entity (type, id) is required for action 'get_by_entity'")
		}
		return r.productsService.GetProductsByEntity(ctx, input.Entity.Type, input.Entity.ID)
	case "link":
		if input.Entity == nil || input.Product == nil {
			return nil, fmt.Errorf("entity and product are required for action 'link'")
		}
		return nil, r.productsService.LinkProduct(ctx, input.Entity.Type, input.Entity.ID, input.Product.ID, input.Product.Quantity, input.Product.PriceID)
	case "unlink":
		if input.Entity == nil || input.ProductID == 0 {
			return nil, fmt.Errorf("entity and product_id are required for action 'unlink'")
		}
		return nil, r.productsService.UnlinkProduct(ctx, input.Entity.Type, input.Entity.ID, input.ProductID)
	case "update_quantity":
		if input.Entity == nil || input.Product == nil {
			return nil, fmt.Errorf("entity and product (id, quantity, price_id) are required for action 'update_quantity'")
		}
		// В amoCRM обновление количества — это просто повторный Link (Link в v4 обновляет metadata если связь уже есть)
		return nil, r.productsService.LinkProduct(ctx, input.Entity.Type, input.Entity.ID, input.Product.ID, input.Product.Quantity, input.Product.PriceID)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}
