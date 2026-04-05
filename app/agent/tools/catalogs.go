package tools

import (
	"encoding/json"
	"fmt"

	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/catalogs"
)

// catalogsRequiredFields определяет обязательные поля для каждого action.
// Если все обязательные поля присутствуют — режим Execute, иначе — режим Schema.
var catalogsRequiredFields = map[string][]string{
	"list":           {},                                   // нет обязательных полей — всегда Execute
	"get":            {"catalog_name"},
	"create":         {"data"},
	"update":         {"catalog_name", "data"},
	"delete":         {"catalog_name"},
	"list_elements":  {"catalog_name"},
	"get_element":    {"catalog_name", "element_id"},
	"create_element": {"catalog_name", "element_data"},
	"update_element": {"catalog_name", "element_id", "element_data"},
	"delete_element": {"catalog_name", "element_id"},
	"link_element":   {"catalog_name", "element_id", "link_data"},
	"unlink_element": {"catalog_name", "element_id", "link_data"},
}

// catalogsSchemas содержит полную схему полей для каждого action.
var catalogsSchemas = map[string]map[string]any{
	"list": {
		"description":     "Возвращает список всех каталогов. Опционально фильтрует по типу.",
		"required_fields": map[string]any{},
		"optional_fields": map[string]any{
			"filter": map[string]any{
				"type":        "object",
				"description": "Параметры фильтрации",
				"properties": map[string]any{
					"type":  map[string]any{"type": "string", "description": "Тип каталога: regular, invoices, products"},
					"page":  map[string]any{"type": "integer", "description": "Номер страницы (начиная с 1)"},
					"limit": map[string]any{"type": "integer", "description": "Лимит результатов (по умолчанию 50, максимум 250)"},
				},
			},
		},
		"example": map[string]any{
			"action": "list",
		},
	},
	"get": {
		"description": "Возвращает один каталог по имени.",
		"required_fields": map[string]any{
			"catalog_name": map[string]any{"type": "string", "description": "Название каталога из available_values.catalog_names"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action":       "get",
			"catalog_name": "Товары",
		},
	},
	"create": {
		"description": "Создаёт новый каталог.",
		"required_fields": map[string]any{
			"data": map[string]any{
				"type":        "object",
				"description": "Данные нового каталога",
				"properties": map[string]any{
					"name": map[string]any{"type": "string", "description": "Название каталога"},
					"type": map[string]any{"type": "string", "description": "Тип: regular, invoices, products"},
				},
			},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action": "create",
			"data":   map[string]any{"name": "Новый каталог", "type": "regular"},
		},
	},
	"update": {
		"description": "Обновляет существующий каталог по имени.",
		"required_fields": map[string]any{
			"catalog_name": map[string]any{"type": "string", "description": "Название каталога из available_values.catalog_names"},
			"data": map[string]any{
				"type":        "object",
				"description": "Новые данные каталога",
				"properties": map[string]any{
					"name":              map[string]any{"type": "string", "description": "Новое название"},
					"can_add_elements":  map[string]any{"type": "boolean", "description": "Разрешить добавление элементов"},
					"can_show_in_cards": map[string]any{"type": "boolean", "description": "Показывать в карточках"},
					"can_link_multiple": map[string]any{"type": "boolean", "description": "Разрешить множественную привязку"},
				},
			},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action":       "update",
			"catalog_name": "Товары",
			"data":         map[string]any{"name": "Новое название"},
		},
	},
	"delete": {
		"description": "Удаляет каталог по имени.",
		"required_fields": map[string]any{
			"catalog_name": map[string]any{"type": "string", "description": "Название каталога из available_values.catalog_names"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action":       "delete",
			"catalog_name": "Товары",
		},
	},
	"list_elements": {
		"description": "Возвращает список элементов каталога.",
		"required_fields": map[string]any{
			"catalog_name": map[string]any{"type": "string", "description": "Название каталога из available_values.catalog_names"},
		},
		"optional_fields": map[string]any{
			"filter": map[string]any{
				"type":        "object",
				"description": "Параметры фильтрации",
				"properties": map[string]any{
					"query": map[string]any{"type": "string", "description": "Поисковый запрос по названию элемента"},
					"ids":   map[string]any{"type": "array", "items": map[string]any{"type": "integer"}, "description": "Фильтр по массиву ID"},
					"page":  map[string]any{"type": "integer", "description": "Номер страницы"},
					"limit": map[string]any{"type": "integer", "description": "Лимит (макс 250)"},
				},
			},
		},
		"example": map[string]any{
			"action":       "list_elements",
			"catalog_name": "Товары",
		},
	},
	"get_element": {
		"description": "Возвращает элемент каталога по ID.",
		"required_fields": map[string]any{
			"catalog_name": map[string]any{"type": "string", "description": "Название каталога из available_values.catalog_names"},
			"element_id":   map[string]any{"type": "integer", "description": "ID элемента каталога"},
		},
		"optional_fields": map[string]any{
			"with": map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Дополнительные данные: invoice_link, supplier_field_values"},
		},
		"example": map[string]any{
			"action":       "get_element",
			"catalog_name": "Товары",
			"element_id":   123,
		},
	},
	"create_element": {
		"description": "Создаёт новый элемент в каталоге.",
		"required_fields": map[string]any{
			"catalog_name": map[string]any{"type": "string", "description": "Название каталога из available_values.catalog_names"},
			"element_data": map[string]any{
				"type":        "object",
				"description": "Данные нового элемента",
				"properties": map[string]any{
					"name":                 map[string]any{"type": "string", "description": "Название элемента"},
					"custom_fields_values": map[string]any{"type": "array", "description": "Кастомные поля: [{field_code, value}]"},
				},
			},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action":       "create_element",
			"catalog_name": "Товары",
			"element_data": map[string]any{"name": "Новый элемент"},
		},
	},
	"update_element": {
		"description": "Обновляет элемент каталога.",
		"required_fields": map[string]any{
			"catalog_name": map[string]any{"type": "string", "description": "Название каталога из available_values.catalog_names"},
			"element_id":   map[string]any{"type": "integer", "description": "ID элемента каталога"},
			"element_data": map[string]any{
				"type":        "object",
				"description": "Новые данные элемента",
				"properties": map[string]any{
					"name":                 map[string]any{"type": "string", "description": "Новое название"},
					"custom_fields_values": map[string]any{"type": "array", "description": "Кастомные поля: [{field_code, value}]"},
				},
			},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action":       "update_element",
			"catalog_name": "Товары",
			"element_id":   123,
			"element_data": map[string]any{"name": "Обновлённый элемент"},
		},
	},
	"delete_element": {
		"description": "Удаляет элемент каталога.",
		"required_fields": map[string]any{
			"catalog_name": map[string]any{"type": "string", "description": "Название каталога из available_values.catalog_names"},
			"element_id":   map[string]any{"type": "integer", "description": "ID элемента каталога"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action":       "delete_element",
			"catalog_name": "Товары",
			"element_id":   123,
		},
	},
	"link_element": {
		"description": "Привязывает элемент каталога к сущности amoCRM.",
		"required_fields": map[string]any{
			"catalog_name": map[string]any{"type": "string", "description": "Название каталога из available_values.catalog_names"},
			"element_id":   map[string]any{"type": "integer", "description": "ID элемента каталога"},
			"link_data": map[string]any{
				"type":        "object",
				"description": "Данные связи",
				"properties": map[string]any{
					"entity_type": map[string]any{"type": "string", "description": "Тип сущности: leads, contacts, companies, customers"},
					"entity_id":   map[string]any{"type": "integer", "description": "ID сущности"},
					"metadata":    map[string]any{"type": "object", "description": "Метаданные: quantity (float64), price_id (int)"},
				},
			},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action":       "link_element",
			"catalog_name": "Товары",
			"element_id":   123,
			"link_data":    map[string]any{"entity_type": "leads", "entity_id": 456},
		},
	},
	"unlink_element": {
		"description": "Отвязывает элемент каталога от сущности amoCRM.",
		"required_fields": map[string]any{
			"catalog_name": map[string]any{"type": "string", "description": "Название каталога из available_values.catalog_names"},
			"element_id":   map[string]any{"type": "integer", "description": "ID элемента каталога"},
			"link_data": map[string]any{
				"type":        "object",
				"description": "Данные связи для отвязки",
				"properties": map[string]any{
					"entity_type": map[string]any{"type": "string", "description": "Тип сущности: leads, contacts, companies, customers"},
					"entity_id":   map[string]any{"type": "integer", "description": "ID сущности"},
				},
			},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action":       "unlink_element",
			"catalog_name": "Товары",
			"element_id":   123,
			"link_data":    map[string]any{"entity_type": "leads", "entity_id": 456},
		},
	},
}

// CatalogsTool — нативный ADK tool для управления каталогами amoCRM.
type CatalogsTool struct {
	service catalogs.Service
}

// NewCatalogsTool создаёт новый CatalogsTool с указанным сервисом.
func NewCatalogsTool(service catalogs.Service) *CatalogsTool {
	return &CatalogsTool{service: service}
}

// Name возвращает имя инструмента.
func (t *CatalogsTool) Name() string { return "catalogs" }

// Description возвращает описание инструмента.
func (t *CatalogsTool) Description() string {
	return "Управление каталогами и их элементами в amoCRM. " +
		"Actions (каталоги): list, get, create, update, delete. " +
		"Actions (элементы): list_elements, get_element, create_element, update_element, delete_element, link_element, unlink_element. " +
		"Вызови только с action чтобы получить полную схему параметров и список доступных каталогов."
}

// IsLongRunning всегда false для синхронных CRM операций.
func (t *CatalogsTool) IsLongRunning() bool { return false }

// ProcessRequest реализует toolinternal.RequestProcessor — регистрирует Declaration в LLM request.
func (t *CatalogsTool) ProcessRequest(_ tool.Context, req *model.LLMRequest) error {
	return packToolDeclaration(req, t)
}

// Declaration возвращает genai.FunctionDeclaration для регистрации в ADK.
func (t *CatalogsTool) Declaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"action": {
					Type:        genai.TypeString,
					Description: "Действие: list, get, create, update, delete (каталоги); list_elements, get_element, create_element, update_element, delete_element, link_element, unlink_element (элементы)",
				},
				"catalog_name": {
					Type:        genai.TypeString,
					Description: "Название каталога (из available_values.catalog_names)",
				},
				"element_id": {
					Type:        genai.TypeInteger,
					Description: "ID элемента каталога",
				},
				"with": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeString,
					},
					Description: "Дополнительные данные для get_element: invoice_link, supplier_field_values",
				},
				"filter": {
					Type:        genai.TypeObject,
					Description: "Параметры фильтрации (page, limit, query, ids, type)",
				},
				"data": {
					Type:        genai.TypeObject,
					Description: "Данные каталога для create/update",
				},
				"element_data": {
					Type:        genai.TypeObject,
					Description: "Данные элемента каталога для create_element/update_element",
				},
				"link_data": {
					Type:        genai.TypeObject,
					Description: "Данные связи элемента с сущностью для link_element/unlink_element",
				},
			},
			Required: []string{"action"},
		},
	}
}

// Run выполняет инструмент. Реализует Shadow Tool паттерн:
// Schema mode: обязательные поля action'а отсутствуют → возвращает схему + catalog_names.
// Execute mode: все обязательные поля присутствуют → выполняет действие.
func (t *CatalogsTool) Run(ctx tool.Context, args any) (map[string]any, error) {
	m, ok := args.(map[string]any)
	if !ok {
		// args может прийти как []byte или строка — попробуем десериализовать
		raw, err := json.Marshal(args)
		if err != nil {
			return nil, fmt.Errorf("catalogs: неверный формат args")
		}
		if err := json.Unmarshal(raw, &m); err != nil {
			return nil, fmt.Errorf("catalogs: неверный формат args")
		}
	}

	result, err := t.handleCatalogsShadow(ctx, m)
	if err != nil {
		return nil, err
	}
	return toResultMap(result)
}

// handleCatalogsShadow реализует Shadow Tool паттерн:
// Schema mode: action присутствует, обязательные поля action'а отсутствуют → возвращает схему + catalog_names.
// Execute mode: все обязательные поля присутствуют → выполняет действие.
func (t *CatalogsTool) handleCatalogsShadow(ctx tool.Context, input map[string]any) (any, error) {
	action, _ := input["action"].(string)
	if action == "" {
		return map[string]any{
			"schema":            true,
			"tool":              "catalogs",
			"error":             "поле action обязательно",
			"available_actions": []string{"list", "get", "create", "update", "delete", "list_elements", "get_element", "create_element", "update_element", "delete_element", "link_element", "unlink_element"},
			"hint":              "Вызови с action чтобы получить полную схему параметров",
		}, nil
	}

	required, known := catalogsRequiredFields[action]
	if !known {
		return nil, fmt.Errorf("unknown action: %s. Доступные: list, get, create, update, delete, list_elements, get_element, create_element, update_element, delete_element, link_element, unlink_element", action)
	}

	// Граница schema/execute: isSchemaMode определён в unsorted.go (пакет tools)
	if isSchemaMode(input, required) {
		return t.catalogsSchemaResponse(action), nil
	}

	// Execute mode: десериализуем map в CatalogsInput через JSON roundtrip
	catalogsInput, err := mapToCatalogsInput(input)
	if err != nil {
		return nil, fmt.Errorf("catalogs: failed to parse input: %w", err)
	}

	return t.executeCatalogs(ctx, catalogsInput)
}

// catalogsSchemaResponse формирует ответ со схемой для данного action + доступные каталоги.
func (t *CatalogsTool) catalogsSchemaResponse(action string) map[string]any {
	schema, ok := catalogsSchemas[action]
	if !ok {
		schema = map[string]any{
			"description":     fmt.Sprintf("Схема для action '%s' не найдена", action),
			"required_fields": map[string]any{},
			"optional_fields": map[string]any{},
		}
	}

	resp := map[string]any{
		"schema": true,
		"tool":   "catalogs",
		"action": action,
		"available_values": map[string]any{
			"catalog_names": t.service.CatalogNames(),
		},
	}
	for k, v := range schema {
		resp[k] = v
	}
	return resp
}

// mapToCatalogsInput десериализует map[string]any в CatalogsInput через JSON roundtrip.
func mapToCatalogsInput(raw map[string]any) (gkitmodels.CatalogsInput, error) {
	data, err := json.Marshal(raw)
	if err != nil {
		return gkitmodels.CatalogsInput{}, err
	}
	var input gkitmodels.CatalogsInput
	if err := json.Unmarshal(data, &input); err != nil {
		return gkitmodels.CatalogsInput{}, err
	}
	return input, nil
}

// executeCatalogs выполняет действие с каталогами/элементами (Execute mode).
func (t *CatalogsTool) executeCatalogs(ctx tool.Context, input gkitmodels.CatalogsInput) (any, error) {
	switch input.Action {

	// Catalogs
	case "list":
		return t.service.ListCatalogs(ctx, input.Filter)

	case "get":
		if input.CatalogName == "" {
			return nil, fmt.Errorf("catalog_name is required. Available: %s", joinNames(t.service.CatalogNames()))
		}
		return t.service.GetCatalog(ctx, input.CatalogName)

	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for create")
		}
		return t.service.CreateCatalog(ctx, input.Data)

	case "update":
		if input.CatalogName == "" {
			return nil, fmt.Errorf("catalog_name is required. Available: %s", joinNames(t.service.CatalogNames()))
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for update")
		}
		return t.service.UpdateCatalog(ctx, input.CatalogName, input.Data)

	case "delete":
		if input.CatalogName == "" {
			return nil, fmt.Errorf("catalog_name is required. Available: %s", joinNames(t.service.CatalogNames()))
		}
		if err := t.service.DeleteCatalog(ctx, input.CatalogName); err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil

	// Elements
	case "list_elements":
		if input.CatalogName == "" {
			return nil, fmt.Errorf("catalog_name is required. Available: %s", joinNames(t.service.CatalogNames()))
		}
		return t.service.ListElements(ctx, input.CatalogName, input.Filter)

	case "get_element":
		if input.CatalogName == "" {
			return nil, fmt.Errorf("catalog_name is required. Available: %s", joinNames(t.service.CatalogNames()))
		}
		if input.ElementID == 0 {
			return nil, fmt.Errorf("element_id is required")
		}
		return t.service.GetElement(ctx, input.CatalogName, input.ElementID, input.With)

	case "create_element":
		if input.CatalogName == "" {
			return nil, fmt.Errorf("catalog_name is required. Available: %s", joinNames(t.service.CatalogNames()))
		}
		if input.ElementData == nil {
			return nil, fmt.Errorf("element_data is required for create_element")
		}
		return t.service.CreateElement(ctx, input.CatalogName, input.ElementData)

	case "update_element":
		if input.CatalogName == "" {
			return nil, fmt.Errorf("catalog_name is required. Available: %s", joinNames(t.service.CatalogNames()))
		}
		if input.ElementID == 0 {
			return nil, fmt.Errorf("element_id is required")
		}
		if input.ElementData == nil {
			return nil, fmt.Errorf("element_data is required for update_element")
		}
		return t.service.UpdateElement(ctx, input.CatalogName, input.ElementID, input.ElementData)

	case "delete_element":
		if input.CatalogName == "" {
			return nil, fmt.Errorf("catalog_name is required. Available: %s", joinNames(t.service.CatalogNames()))
		}
		if input.ElementID == 0 {
			return nil, fmt.Errorf("element_id is required")
		}
		if err := t.service.DeleteElement(ctx, input.CatalogName, input.ElementID); err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil

	// Link/Unlink
	case "link_element":
		if input.CatalogName == "" {
			return nil, fmt.Errorf("catalog_name is required. Available: %s", joinNames(t.service.CatalogNames()))
		}
		if input.ElementID == 0 {
			return nil, fmt.Errorf("element_id is required")
		}
		if input.LinkData == nil {
			return nil, fmt.Errorf("link_data is required")
		}
		err := t.service.LinkElement(
			ctx, input.CatalogName, input.ElementID,
			input.LinkData.EntityType, input.LinkData.EntityID, input.LinkData.Metadata,
		)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil

	case "unlink_element":
		if input.CatalogName == "" {
			return nil, fmt.Errorf("catalog_name is required. Available: %s", joinNames(t.service.CatalogNames()))
		}
		if input.ElementID == 0 {
			return nil, fmt.Errorf("element_id is required")
		}
		if input.LinkData == nil {
			return nil, fmt.Errorf("link_data is required")
		}
		err := t.service.UnlinkElement(
			ctx, input.CatalogName, input.ElementID,
			input.LinkData.EntityType, input.LinkData.EntityID,
		)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil

	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

// joinNames склеивает срез строк через ", " для вывода подсказок.
func joinNames(names []string) string {
	if len(names) == 0 {
		return "(нет доступных)"
	}
	result := ""
	for i, n := range names {
		if i > 0 {
			result += ", "
		}
		result += n
	}
	return result
}

