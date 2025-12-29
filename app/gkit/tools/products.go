package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
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
			"Поддерживает: search (поиск), get (получение по ID), create (создание), update (обновление), delete (удаление).",
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
		if input.ProductID == 0 {
			return nil, fmt.Errorf("product_id is required for action 'get'")
		}
		return r.sdk.Products().GetOne(ctx, input.ProductID)
	case "create":
		if input.Data == nil || input.Data.Name == "" {
			return nil, fmt.Errorf("data.name is required for action 'create'")
		}
		return r.createProduct(ctx, input.Data)
	case "update":
		if input.ProductID == 0 {
			return nil, fmt.Errorf("product_id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		return r.updateProduct(ctx, input.ProductID, input.Data)
	case "delete":
		if len(input.IDs) == 0 {
			return nil, fmt.Errorf("ids is required for action 'delete'")
		}
		return nil, r.sdk.Products().Delete(ctx, input.IDs)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (r *Registry) searchProducts(ctx context.Context, filter *ProductFilter) ([]models.CatalogElement, error) {
	f := &services.ProductsFilter{
		Limit: 50,
		Page:  1,
	}
	if filter != nil {
		if filter.Query != "" {
			f.Query = filter.Query
		}
		if filter.Limit > 0 {
			f.Limit = filter.Limit
		}
		if filter.Page > 0 {
			f.Page = filter.Page
		}
	}
	return r.sdk.Products().Get(ctx, f)
}

func (r *Registry) createProduct(ctx context.Context, data *ProductData) ([]models.CatalogElement, error) {
	product := models.CatalogElement{
		Name: data.Name,
	}
	// TODO: добавить поддержку custom_fields для SKU, price и т.д.
	return r.sdk.Products().Create(ctx, []models.CatalogElement{product})
}

func (r *Registry) updateProduct(ctx context.Context, id int, data *ProductData) ([]models.CatalogElement, error) {
	product := models.CatalogElement{
		ID: id,
	}
	if data.Name != "" {
		product.Name = data.Name
	}
	return r.sdk.Products().Update(ctx, []models.CatalogElement{product})
}
