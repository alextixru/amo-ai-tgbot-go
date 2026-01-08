package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// ProductsInput входные параметры для инструмента products
type ProductsInput struct {
	// Action действие: search, get, create, update, delete
	Action string `json:"action" jsonschema_description:"Действие: search, get, create, update, delete"`

	// ProductID ID товара (для get, update)
	ProductID int `json:"product_id,omitempty" jsonschema_description:"ID товара (для get, update)"`

	// Filter параметры поиска
	Filter *ProductFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска (для search)"`

	// Data данные для создания/обновления
	Data *ProductData `json:"data,omitempty" jsonschema_description:"Данные товара (для create, update)"`

	// IDs массив ID для удаления
	IDs []int `json:"ids,omitempty" jsonschema_description:"Массив ID товаров (для delete)"`
}

// ProductFilter фильтры поиска товаров
type ProductFilter struct {
	Query string `json:"query,omitempty" jsonschema_description:"Поисковый запрос"`
	Limit int    `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Page  int    `json:"page,omitempty" jsonschema_description:"Номер страницы"`
}

// ProductData данные товара
type ProductData struct {
	Name string `json:"name" jsonschema_description:"Название товара"`
	SKU  string `json:"sku,omitempty" jsonschema_description:"Артикул"`
}

// registerProductsTool регистрирует инструмент для работы с товарами
func (r *Registry) registerProductsTool() {
	r.addTool(genkit.DefineTool[ProductsInput, any](
		r.g,
		"products",
		"Работа с товарами amoCRM (каталог ID=1). "+
			"⚠️ ВНИМАНИЕ: ProductsService.Get/Create/Update возвращают ErrNotAvailableForAction согласно API. "+
			"Для работы с товарами используйте 'catalogs' tool с catalog_id товарного каталога вместо этого tool. "+
			"Поддерживает: search (поиск, может не работать), get (получение по ID), create/update (недоступны), delete (удаление).",
		func(ctx *ai.ToolContext, input ProductsInput) (any, error) {
			return r.handleProducts(ctx.Context, input)
		},
	))
}

func (r *Registry) handleProducts(ctx context.Context, input ProductsInput) (any, error) {
	switch input.Action {
	case "search":
		return r.searchProducts(ctx, input.Filter)
	case "get":
		return nil, fmt.Errorf("ProductsService.GetOne is not available. Please use 'catalogs' tool with your products catalog ID")
	case "create":
		return nil, fmt.Errorf("ProductsService.Create is not available. Please use 'catalogs' tool with your products catalog ID")
	case "update":
		return nil, fmt.Errorf("ProductsService.Update is not available. Please use 'catalogs' tool with your products catalog ID")
	case "delete":
		return nil, fmt.Errorf("ProductsService.Delete is not available. Please use 'catalogs' tool with your products catalog ID")
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (r *Registry) searchProducts(ctx context.Context, filter *ProductFilter) ([]*models.CatalogElement, error) {
	return nil, fmt.Errorf("ProductsService.Get is not available. Please use 'catalogs' tool with your products catalog ID")
}

func (r *Registry) createProduct(ctx context.Context, data *ProductData) ([]*models.CatalogElement, error) {
	return nil, fmt.Errorf("ProductsService.Create is not available. Please use 'catalogs' tool with your products catalog ID")
}

func (r *Registry) updateProduct(ctx context.Context, id int, data *ProductData) ([]*models.CatalogElement, error) {
	return nil, fmt.Errorf("ProductsService.Update is not available. Please use 'catalogs' tool with your products catalog ID")
}
