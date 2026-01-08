package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// CatalogsInput входные параметры для инструмента catalogs
type CatalogsInput struct {
	// Action действие: list, get, create, update, list_elements, get_element, create_element, update_element
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, list_elements, get_element, create_element, update_element"`

	// CatalogID ID каталога (для get, update, list_elements, get_element, create_element, update_element)
	CatalogID int `json:"catalog_id,omitempty" jsonschema_description:"ID каталога"`

	// ElementID ID элемента (для get_element, update_element)
	ElementID int `json:"element_id,omitempty" jsonschema_description:"ID элемента каталога"`

	// Filter параметры поиска
	Filter *CatalogFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска"`

	// Data данные каталога (для create, update)
	Data *CatalogData `json:"data,omitempty" jsonschema_description:"Данные каталога"`

	// ElementData данные элемента (для create_element, update_element)
	ElementData *CatalogElementData `json:"element_data,omitempty" jsonschema_description:"Данные элемента каталога"`
}

// CatalogFilter фильтры поиска
type CatalogFilter struct {
	Page  int    `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Limit int    `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Query string `json:"query,omitempty" jsonschema_description:"Поисковый запрос (только для элементов)"`
}

// CatalogData данные каталога
type CatalogData struct {
	Name            string `json:"name" jsonschema_description:"Название каталога"`
	Type            string `json:"type,omitempty" jsonschema_description:"Тип: regular, invoices, products"`
	CanAddElements  bool   `json:"can_add_elements,omitempty" jsonschema_description:"Можно ли добавлять элементы"`
	CanShowInCards  bool   `json:"can_show_in_cards,omitempty" jsonschema_description:"Показывать в карточках"`
	CanLinkMultiple bool   `json:"can_link_multiple,omitempty" jsonschema_description:"Разрешить множественную привязку"`
}

// CatalogElementData данные элемента каталога
type CatalogElementData struct {
	Name               string         `json:"name" jsonschema_description:"Название элемента"`
	CustomFieldsValues map[string]any `json:"custom_fields_values,omitempty" jsonschema_description:"Значения кастомных полей"`
}

// registerCatalogsTool регистрирует инструмент для работы с каталогами
func (r *Registry) registerCatalogsTool() {
	r.addTool(genkit.DefineTool[CatalogsInput, any](
		r.g,
		"catalogs",
		"Работа с каталогами и их элементами. "+
			"Каталоги: list (список), get (по ID), create (создать), update (обновить). "+
			"Элементы: list_elements, get_element, create_element, update_element. "+
			"Для работы с элементами требуется catalog_id.",
		func(ctx *ai.ToolContext, input CatalogsInput) (any, error) {
			return r.handleCatalogs(ctx.Context, input)
		},
	))
}

func (r *Registry) handleCatalogs(ctx context.Context, input CatalogsInput) (any, error) {
	switch input.Action {
	case "list":
		return r.listCatalogs(ctx, input.Filter)
	case "get":
		if input.CatalogID == 0 {
			return nil, fmt.Errorf("catalog_id is required for action 'get'")
		}
		return r.sdk.Catalogs().GetOne(ctx, input.CatalogID)
	case "create":
		if input.Data == nil || input.Data.Name == "" {
			return nil, fmt.Errorf("data.name is required for action 'create'")
		}
		return r.createCatalog(ctx, input.Data)
	case "update":
		if input.CatalogID == 0 {
			return nil, fmt.Errorf("catalog_id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		return r.updateCatalog(ctx, input.CatalogID, input.Data)
	case "list_elements":
		if input.CatalogID == 0 {
			return nil, fmt.Errorf("catalog_id is required for action 'list_elements'")
		}
		return r.listCatalogElements(ctx, input.CatalogID, input.Filter)
	case "get_element":
		if input.CatalogID == 0 {
			return nil, fmt.Errorf("catalog_id is required for action 'get_element'")
		}
		if input.ElementID == 0 {
			return nil, fmt.Errorf("element_id is required for action 'get_element'")
		}
		return r.sdk.CatalogElements(input.CatalogID).GetOne(ctx, input.ElementID, nil)
	case "create_element":
		if input.CatalogID == 0 {
			return nil, fmt.Errorf("catalog_id is required for action 'create_element'")
		}
		if input.ElementData == nil || input.ElementData.Name == "" {
			return nil, fmt.Errorf("element_data.name is required for action 'create_element'")
		}
		return r.createCatalogElement(ctx, input.CatalogID, input.ElementData)
	case "update_element":
		if input.CatalogID == 0 {
			return nil, fmt.Errorf("catalog_id is required for action 'update_element'")
		}
		if input.ElementID == 0 {
			return nil, fmt.Errorf("element_id is required for action 'update_element'")
		}
		if input.ElementData == nil {
			return nil, fmt.Errorf("element_data is required for action 'update_element'")
		}
		return r.updateCatalogElement(ctx, input.CatalogID, input.ElementID, input.ElementData)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

// ============ CATALOGS ============

func (r *Registry) listCatalogs(ctx context.Context, filter *CatalogFilter) ([]*models.Catalog, error) {
	// CatalogsFilter не поддерживает SetLimit/SetPage
	catalogs, _, err := r.sdk.Catalogs().Get(ctx, nil)
	return catalogs, err
}

func (r *Registry) createCatalog(ctx context.Context, data *CatalogData) ([]*models.Catalog, error) {
	catalog := &models.Catalog{
		Name:            data.Name,
		CanAddElements:  data.CanAddElements,
		CanShowInCards:  data.CanShowInCards,
		CanLinkMultiple: data.CanLinkMultiple,
	}
	catalogs, _, err := r.sdk.Catalogs().Create(ctx, []*models.Catalog{catalog})
	return catalogs, err
}

func (r *Registry) updateCatalog(ctx context.Context, id int, data *CatalogData) ([]*models.Catalog, error) {
	catalog := &models.Catalog{
		ID: id,
	}
	if data.Name != "" {
		catalog.Name = data.Name
	}
	catalog.CanAddElements = data.CanAddElements
	catalog.CanShowInCards = data.CanShowInCards
	catalog.CanLinkMultiple = data.CanLinkMultiple

	catalogs, _, err := r.sdk.Catalogs().Update(ctx, []*models.Catalog{catalog})
	return catalogs, err
}

// ============ CATALOG ELEMENTS ============

func (r *Registry) listCatalogElements(ctx context.Context, catalogID int, filter *CatalogFilter) ([]*models.CatalogElement, error) {
	// CatalogElementsService наследует BaseEntityIdService, Get принимает url.Values
	elements, _, err := r.sdk.CatalogElements(catalogID).Get(ctx, nil)
	return elements, err
}

func (r *Registry) createCatalogElement(ctx context.Context, catalogID int, data *CatalogElementData) ([]*models.CatalogElement, error) {
	element := &models.CatalogElement{
		Name: data.Name,
	}
	elements, _, err := r.sdk.CatalogElements(catalogID).Create(ctx, []*models.CatalogElement{element})
	return elements, err
}

func (r *Registry) updateCatalogElement(ctx context.Context, catalogID int, elementID int, data *CatalogElementData) ([]*models.CatalogElement, error) {
	element := &models.CatalogElement{
		ID: elementID,
	}
	if data.Name != "" {
		element.Name = data.Name
	}
	elements, _, err := r.sdk.CatalogElements(catalogID).Update(ctx, []*models.CatalogElement{element})
	return elements, err
}
