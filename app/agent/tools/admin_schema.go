package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/admin_schema"
)

// AdminSchemaTool — нативный ADK tool для управления схемой данных CRM.
type AdminSchemaTool struct {
	service admin_schema.Service
}

// NewAdminSchemaTool создаёт новый экземпляр AdminSchemaTool.
func NewAdminSchemaTool(service admin_schema.Service) *AdminSchemaTool {
	return &AdminSchemaTool{service: service}
}

// Name реализует tool.Tool.
func (t *AdminSchemaTool) Name() string {
	return "admin_schema"
}

// Description реализует tool.Tool.
func (t *AdminSchemaTool) Description() string {
	return "Управление схемой данных CRM: кастомные поля, группы полей, причины отказа, источники сделок. " +
		"Layers: custom_fields, field_groups, loss_reasons, sources. " +
		"Вызови с layer + action чтобы получить схему параметров."
}

// IsLongRunning реализует tool.Tool.
func (t *AdminSchemaTool) IsLongRunning() bool {
	return false
}

// Declaration реализует toolinternal.FunctionTool (duck typing).
func (t *AdminSchemaTool) Declaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"layer": {
					Type:        genai.TypeString,
					Description: "Слой схемы: custom_fields, field_groups, loss_reasons, sources",
					Enum:        []string{"custom_fields", "field_groups", "loss_reasons", "sources"},
				},
				"action": {
					Type:        genai.TypeString,
					Description: "Действие: list, get, create, update, delete",
					Enum:        []string{"list", "get", "create", "update", "delete"},
				},
				"entity_type": {
					Type:        genai.TypeString,
					Description: "Тип сущности (для custom_fields и field_groups): leads, contacts, companies, customers",
				},
				"id": {
					Type:        genai.TypeInteger,
					Description: "ID ресурса (для get/delete)",
				},
				"group_id": {
					Type:        genai.TypeString,
					Description: "ID группы полей (строковый, для field_groups get/delete)",
				},
				"filter": {
					Type:        genai.TypeObject,
					Description: "Фильтры для list-запросов",
				},
				"custom_field": {
					Type:        genai.TypeObject,
					Description: "Одиночное кастомное поле для create/update",
				},
				"custom_fields": {
					Type:        genai.TypeArray,
					Description: "Батч кастомных полей для create/update",
					Items:       &genai.Schema{Type: genai.TypeObject},
				},
				"field_group": {
					Type:        genai.TypeObject,
					Description: "Одиночная группа полей для create/update",
				},
				"field_groups": {
					Type:        genai.TypeArray,
					Description: "Батч групп полей для create/update",
					Items:       &genai.Schema{Type: genai.TypeObject},
				},
				"loss_reason": {
					Type:        genai.TypeObject,
					Description: "Одиночная причина отказа для create",
				},
				"loss_reasons": {
					Type:        genai.TypeArray,
					Description: "Батч причин отказа для create",
					Items:       &genai.Schema{Type: genai.TypeObject},
				},
				"source": {
					Type:        genai.TypeObject,
					Description: "Одиночный источник для create/update",
				},
				"sources": {
					Type:        genai.TypeArray,
					Description: "Батч источников для create/update",
					Items:       &genai.Schema{Type: genai.TypeObject},
				},
			},
			Required: []string{"layer", "action"},
		},
	}
}

// Run реализует toolinternal.FunctionTool (duck typing).
func (t *AdminSchemaTool) Run(ctx tool.Context, args any) (map[string]any, error) {
	input, ok := args.(map[string]any)
	if !ok {
		// Попытка через JSON roundtrip
		b, err := json.Marshal(args)
		if err != nil {
			return nil, fmt.Errorf("не удалось сериализовать input: %w", err)
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, fmt.Errorf("не удалось разобрать input как map: %w", err)
		}
		input = m
	}

	layer, _ := input["layer"].(string)
	action, _ := input["action"].(string)

	if layer == "" || action == "" {
		return map[string]any{
			"schema":      true,
			"tool":        "admin_schema",
			"description": "Управление схемой данных CRM. Укажи layer и action для получения полной схемы.",
			"layers": map[string]any{
				"custom_fields": "Кастомные поля сущностей. Actions: list, get, create, update, delete",
				"field_groups":  "Группы кастомных полей. Actions: list, get, create, update, delete",
				"loss_reasons":  "Причины отказа (проигрыша). Actions: list, get, create, delete",
				"sources":       "Источники сделок. Actions: list, get, create, update, delete",
			},
			"usage": "Вызови с layer + action, например: {\"layer\": \"custom_fields\", \"action\": \"list\"}",
		}, nil
	}

	// Schema mode: обязательные поля отсутствуют
	if adminSchemaIsSchemaMode(input, layer, action) {
		layerSchemas, ok := adminSchemaSchemas[layer]
		if !ok {
			return nil, fmt.Errorf("неизвестный layer: %s. Доступные: custom_fields, field_groups, loss_reasons, sources", layer)
		}
		schema, ok := layerSchemas[action]
		if !ok {
			return nil, fmt.Errorf("неизвестный action %q для layer %q", action, layer)
		}
		return toResultMap(schema)
	}

	// Execute mode: JSON roundtrip в AdminSchemaInput
	b, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("не удалось сериализовать input: %w", err)
	}
	var typed gkitmodels.AdminSchemaInput
	if err := json.Unmarshal(b, &typed); err != nil {
		return nil, fmt.Errorf("не удалось разобрать input как AdminSchemaInput: %w", err)
	}

	var result any
	switch typed.Layer {
	case "custom_fields":
		result, err = t.handleCustomFields(ctx, typed)
	case "field_groups":
		result, err = t.handleFieldGroups(ctx, typed)
	case "loss_reasons":
		result, err = t.handleLossReasons(ctx, typed)
	case "sources":
		result, err = t.handleSources(ctx, typed)
	default:
		return nil, fmt.Errorf("unknown layer: %s", typed.Layer)
	}
	if err != nil {
		return nil, err
	}
	return toResultMap(result)
}

// adminSchemaRequiredFields определяет обязательные поля для каждой комбинации layer+action.
// Если все перечисленные поля отсутствуют в input — возвращаем схему.
// Пустой срез означает, что layer+action не требует дополнительных полей (например list).
var adminSchemaRequiredFields = map[string]map[string][]string{
	"custom_fields": {
		"list":   {"entity_type"},
		"get":    {"entity_type", "id"},
		"create": {"entity_type"},
		"update": {"entity_type"},
		"delete": {"entity_type", "id"},
	},
	"field_groups": {
		"list":   {"entity_type"},
		"get":    {"entity_type", "group_id"},
		"create": {"entity_type"},
		"update": {"entity_type"},
		"delete": {"entity_type", "group_id"},
	},
	"loss_reasons": {
		"list":   {},
		"get":    {"id"},
		"create": {},
		"delete": {"id"},
	},
	"sources": {
		"list":   {},
		"get":    {"id"},
		"create": {},
		"update": {},
		"delete": {"id"},
	},
}

// adminSchemaSchemas содержит полные схемы для каждой комбинации layer+action.
var adminSchemaSchemas = map[string]map[string]any{
	"custom_fields": {
		"list": map[string]any{
			"schema":      true,
			"tool":        "admin_schema",
			"layer":       "custom_fields",
			"action":      "list",
			"description": "Список кастомных полей для указанной сущности. Поддерживает фильтрацию по типу, ID и имени (client-side).",
			"required_fields": map[string]any{
				"entity_type": map[string]any{"type": "string", "description": "Тип сущности: leads, contacts, companies, customers"},
			},
			"optional_fields": map[string]any{
				"filter": map[string]any{
					"type": "object",
					"description": "Фильтры",
					"fields": map[string]any{
						"limit":  map[string]any{"type": "integer", "description": "Лимит результатов (по умолчанию 50)"},
						"page":   map[string]any{"type": "integer", "description": "Номер страницы"},
						"name":   map[string]any{"type": "string", "description": "Фильтр по имени (частичное совпадение, client-side)"},
						"ids":    map[string]any{"type": "array", "description": "Фильтр по ID полей"},
						"types":  map[string]any{"type": "array", "description": "Фильтр по типам: text, numeric, checkbox, select, multiselect, date, url, textarea, radiobutton, multitext и др."},
						"order":  map[string]any{"type": "object", "description": `Сортировка: {"created_at": "desc"} или {"updated_at": "asc"}`},
					},
				},
			},
			"example": map[string]any{
				"layer": "custom_fields", "action": "list", "entity_type": "leads",
				"filter": map[string]any{"limit": 50, "name": "Бюджет"},
			},
		},
		"get": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "custom_fields", "action": "get",
			"description": "Получить кастомное поле по ID.",
			"required_fields": map[string]any{
				"entity_type": map[string]any{"type": "string", "description": "Тип сущности: leads, contacts, companies, customers"},
				"id":          map[string]any{"type": "integer", "description": "ID кастомного поля"},
			},
			"example": map[string]any{"layer": "custom_fields", "action": "get", "entity_type": "leads", "id": 12345},
		},
		"create": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "custom_fields", "action": "create",
			"description": "Создать одно или несколько кастомных полей для сущности.",
			"required_fields": map[string]any{
				"entity_type": map[string]any{"type": "string", "description": "Тип сущности: leads, contacts, companies, customers"},
			},
			"optional_fields": map[string]any{
				"custom_field": map[string]any{
					"type": "object", "description": "Одиночное поле для создания",
					"fields": map[string]any{
						"name":        map[string]any{"type": "string", "description": "Название поля"},
						"type":        map[string]any{"type": "string", "description": "Тип: text, numeric, checkbox, select, multiselect, date, url, textarea, radiobutton, multitext и др."},
						"code":        map[string]any{"type": "string", "description": "Символьный код (латиница, уникален в аккаунте)"},
						"sort":        map[string]any{"type": "integer", "description": "Порядок сортировки"},
						"group_id":    map[string]any{"type": "string", "description": "ID группы полей (узнать через layer=field_groups, action=list)"},
						"is_api_only": map[string]any{"type": "boolean", "description": "Только для API, не отображается в CRM"},
						"is_required": map[string]any{"type": "boolean", "description": "Поле обязательно для заполнения"},
						"enums":       map[string]any{"type": "array", "description": `Варианты для select/multiselect/radiobutton: [{"value": "Вариант 1", "sort": 1}]`},
					},
				},
				"custom_fields": map[string]any{"type": "array", "description": "Батч-создание: массив объектов custom_field"},
			},
			"example": map[string]any{
				"layer": "custom_fields", "action": "create", "entity_type": "leads",
				"custom_field": map[string]any{"name": "Бюджет", "type": "numeric", "sort": 100},
			},
		},
		"update": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "custom_fields", "action": "update",
			"description": "Обновить одно или несколько кастомных полей. ID поля обязателен внутри custom_field/custom_fields.",
			"required_fields": map[string]any{
				"entity_type": map[string]any{"type": "string", "description": "Тип сущности: leads, contacts, companies, customers"},
			},
			"optional_fields": map[string]any{
				"custom_field":  map[string]any{"type": "object", "description": "Поле для обновления (поле id внутри объекта обязательно)"},
				"custom_fields": map[string]any{"type": "array", "description": "Батч-обновление"},
			},
			"example": map[string]any{
				"layer": "custom_fields", "action": "update", "entity_type": "leads",
				"custom_field": map[string]any{"id": 12345, "name": "Новый бюджет"},
			},
		},
		"delete": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "custom_fields", "action": "delete",
			"description": "Удалить кастомное поле по ID.",
			"required_fields": map[string]any{
				"entity_type": map[string]any{"type": "string", "description": "Тип сущности: leads, contacts, companies, customers"},
				"id":          map[string]any{"type": "integer", "description": "ID кастомного поля"},
			},
			"example": map[string]any{"layer": "custom_fields", "action": "delete", "entity_type": "leads", "id": 12345},
		},
	},
	"field_groups": {
		"list": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "field_groups", "action": "list",
			"description": "Список групп кастомных полей для указанной сущности.",
			"required_fields": map[string]any{
				"entity_type": map[string]any{"type": "string", "description": "Тип сущности: leads, contacts, companies, customers"},
			},
			"optional_fields": map[string]any{
				"filter": map[string]any{
					"type": "object",
					"fields": map[string]any{
						"limit": map[string]any{"type": "integer", "description": "Лимит результатов"},
						"page":  map[string]any{"type": "integer", "description": "Номер страницы"},
						"name":  map[string]any{"type": "string", "description": "Фильтр по имени (client-side)"},
					},
				},
			},
			"example": map[string]any{"layer": "field_groups", "action": "list", "entity_type": "leads"},
		},
		"get": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "field_groups", "action": "get",
			"description": "Получить группу полей по group_id.",
			"required_fields": map[string]any{
				"entity_type": map[string]any{"type": "string", "description": "Тип сущности"},
				"group_id":    map[string]any{"type": "string", "description": "ID группы полей (строковый идентификатор, например 'general')"},
			},
			"example": map[string]any{"layer": "field_groups", "action": "get", "entity_type": "leads", "group_id": "general"},
		},
		"create": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "field_groups", "action": "create",
			"description": "Создать одну или несколько групп полей.",
			"required_fields": map[string]any{
				"entity_type": map[string]any{"type": "string", "description": "Тип сущности"},
			},
			"optional_fields": map[string]any{
				"field_group":  map[string]any{"type": "object", "description": `Одна группа: {"name": "Группа", "sort": 10}`},
				"field_groups": map[string]any{"type": "array", "description": "Батч-создание"},
			},
			"example": map[string]any{
				"layer": "field_groups", "action": "create", "entity_type": "leads",
				"field_group": map[string]any{"name": "Финансы", "sort": 10},
			},
		},
		"update": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "field_groups", "action": "update",
			"description": "Обновить группу полей. ID группы (строка) обязателен внутри field_group/field_groups.",
			"required_fields": map[string]any{
				"entity_type": map[string]any{"type": "string", "description": "Тип сущности"},
			},
			"optional_fields": map[string]any{
				"field_group":  map[string]any{"type": "object", "description": "Группа для обновления (поле id внутри обязательно)"},
				"field_groups": map[string]any{"type": "array", "description": "Батч-обновление"},
			},
			"example": map[string]any{
				"layer": "field_groups", "action": "update", "entity_type": "leads",
				"field_group": map[string]any{"id": "my-group", "name": "Обновлённая группа"},
			},
		},
		"delete": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "field_groups", "action": "delete",
			"description": "Удалить группу полей по group_id.",
			"required_fields": map[string]any{
				"entity_type": map[string]any{"type": "string", "description": "Тип сущности"},
				"group_id":    map[string]any{"type": "string", "description": "ID группы полей"},
			},
			"example": map[string]any{"layer": "field_groups", "action": "delete", "entity_type": "leads", "group_id": "my-group"},
		},
	},
	"loss_reasons": {
		"list": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "loss_reasons", "action": "list",
			"description": "Список причин отказа (причин проигрыша сделки).",
			"required_fields": map[string]any{},
			"optional_fields": map[string]any{
				"filter": map[string]any{
					"type": "object",
					"fields": map[string]any{
						"limit": map[string]any{"type": "integer", "description": "Лимит"},
						"page":  map[string]any{"type": "integer", "description": "Страница"},
						"name":  map[string]any{"type": "string", "description": "Фильтр по имени (client-side)"},
					},
				},
			},
			"example": map[string]any{"layer": "loss_reasons", "action": "list"},
		},
		"get": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "loss_reasons", "action": "get",
			"description": "Получить причину отказа по ID.",
			"required_fields": map[string]any{
				"id": map[string]any{"type": "integer", "description": "ID причины отказа"},
			},
			"example": map[string]any{"layer": "loss_reasons", "action": "get", "id": 42},
		},
		"create": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "loss_reasons", "action": "create",
			"description": "Создать одну или несколько причин отказа.",
			"required_fields": map[string]any{},
			"optional_fields": map[string]any{
				"loss_reason":  map[string]any{"type": "object", "description": `Одна причина: {"name": "Дорого", "sort": 10}`},
				"loss_reasons": map[string]any{"type": "array", "description": "Батч-создание"},
			},
			"note": "Хотя бы одно из полей loss_reason или loss_reasons обязательно",
			"example": map[string]any{
				"layer": "loss_reasons", "action": "create",
				"loss_reason": map[string]any{"name": "Дорого", "sort": 10},
			},
		},
		"delete": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "loss_reasons", "action": "delete",
			"description": "Удалить причину отказа по ID.",
			"required_fields": map[string]any{
				"id": map[string]any{"type": "integer", "description": "ID причины отказа"},
			},
			"example": map[string]any{"layer": "loss_reasons", "action": "delete", "id": 42},
		},
	},
	"sources": {
		"list": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "sources", "action": "list",
			"description": "Список источников сделок.",
			"required_fields": map[string]any{},
			"optional_fields": map[string]any{
				"filter": map[string]any{
					"type": "object",
					"fields": map[string]any{
						"name":         map[string]any{"type": "string", "description": "Фильтр по имени (client-side)"},
						"external_ids": map[string]any{"type": "array", "description": "Фильтр по external_id"},
					},
				},
			},
			"example": map[string]any{"layer": "sources", "action": "list"},
		},
		"get": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "sources", "action": "get",
			"description": "Получить источник по ID.",
			"required_fields": map[string]any{
				"id": map[string]any{"type": "integer", "description": "ID источника"},
			},
			"example": map[string]any{"layer": "sources", "action": "get", "id": 100},
		},
		"create": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "sources", "action": "create",
			"description": "Создать один или несколько источников.",
			"required_fields": map[string]any{},
			"optional_fields": map[string]any{
				"source": map[string]any{
					"type": "object", "description": "Один источник",
					"fields": map[string]any{
						"name":        map[string]any{"type": "string", "description": "Название источника"},
						"external_id": map[string]any{"type": "string", "description": "Внешний идентификатор"},
						"origin_code": map[string]any{"type": "string", "description": "Код источника"},
						"pipeline_id": map[string]any{"type": "integer", "description": "ID воронки (получить через admin_pipelines)"},
						"default":     map[string]any{"type": "boolean", "description": "Источник по умолчанию"},
					},
				},
				"sources": map[string]any{"type": "array", "description": "Батч-создание"},
			},
			"note": "Хотя бы одно из полей source или sources обязательно",
			"example": map[string]any{
				"layer": "sources", "action": "create",
				"source": map[string]any{"name": "Сайт", "external_id": "website"},
			},
		},
		"update": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "sources", "action": "update",
			"description": "Обновить один или несколько источников. ID источника обязателен внутри source/sources.",
			"required_fields": map[string]any{},
			"optional_fields": map[string]any{
				"source":  map[string]any{"type": "object", "description": "Источник для обновления (поле id внутри обязательно)"},
				"sources": map[string]any{"type": "array", "description": "Батч-обновление"},
			},
			"note": "Хотя бы одно из полей source или sources обязательно",
			"example": map[string]any{
				"layer": "sources", "action": "update",
				"source": map[string]any{"id": 100, "name": "Новый сайт"},
			},
		},
		"delete": map[string]any{
			"schema": true, "tool": "admin_schema", "layer": "sources", "action": "delete",
			"description": "Удалить источник по ID.",
			"required_fields": map[string]any{
				"id": map[string]any{"type": "integer", "description": "ID источника"},
			},
			"example": map[string]any{"layer": "sources", "action": "delete", "id": 100},
		},
	},
}

// adminSchemaIsSchemaMode определяет, является ли вызов запросом схемы.
// Возвращает true если хотя бы одно обязательное поле для данного layer+action отсутствует в input.
func adminSchemaIsSchemaMode(input map[string]any, layer, action string) bool {
	layerActions, ok := adminSchemaRequiredFields[layer]
	if !ok {
		return false
	}
	required, ok := layerActions[action]
	if !ok {
		return false
	}
	for _, field := range required {
		v, exists := input[field]
		if !exists || v == nil || v == "" || v == 0.0 {
			return true
		}
	}
	return false
}

func buildCustomFieldsFilter(f *gkitmodels.SchemaFilter) *filters.CustomFieldsFilter {
	if f == nil {
		return nil
	}
	filter := &filters.CustomFieldsFilter{
		BaseFilter: filters.NewBaseFilter(),
	}
	if f.Limit > 0 {
		filter.Limit = f.Limit
	}
	if f.Page > 0 {
		filter.Page = f.Page
	}
	if len(f.IDs) > 0 {
		filter.IDs = f.IDs
	}
	if len(f.Types) > 0 {
		filter.Types = f.Types
	}
	for field, dir := range f.Order {
		filter.SetOrder(field, dir)
	}
	return filter
}

func buildSourcesFilter(f *gkitmodels.SchemaFilter) *filters.SourcesFilter {
	if f == nil {
		return nil
	}
	filter := &filters.SourcesFilter{}
	if len(f.ExternalIDs) > 0 {
		filter.ExternalIDs = f.ExternalIDs
	}
	return filter
}

func buildBaseFilter(f *gkitmodels.SchemaFilter) url.Values {
	if f == nil {
		return url.Values{}
	}
	v := url.Values{}
	if f.Limit > 0 {
		v.Add("limit", fmt.Sprintf("%d", f.Limit))
	}
	if f.Page > 0 {
		v.Add("page", fmt.Sprintf("%d", f.Page))
	}
	return v
}

// filterByName выполняет client-side фильтрацию по имени (частичное совпадение без учёта регистра).
// Возвращает индексы совпадающих элементов через функцию getName.
func filterByName[T any](items []T, name string, getName func(T) string) []T {
	if name == "" {
		return items
	}
	needle := strings.ToLower(name)
	result := make([]T, 0, len(items))
	for _, item := range items {
		if strings.Contains(strings.ToLower(getName(item)), needle) {
			result = append(result, item)
		}
	}
	return result
}

// customFieldDataToModel конвертирует CustomFieldData в SDK-модель
func customFieldDataToModel(d gkitmodels.CustomFieldData) *amomodels.CustomField {
	f := &amomodels.CustomField{
		ID:         d.ID,
		Name:       d.Name,
		Type:       amomodels.CustomFieldType(d.Type),
		Code:       d.Code,
		Sort:       d.Sort,
		GroupID:    d.GroupID,
		IsAPIOnly:  d.IsAPIOnly,
		IsRequired: d.IsRequired,
	}
	for _, e := range d.Enums {
		f.Enums = append(f.Enums, amomodels.CustomFieldEnum{
			ID:    e.ID,
			Value: e.Value,
			Sort:  e.Sort,
			Code:  e.Code,
		})
	}
	return f
}

// fieldGroupDataToModel конвертирует FieldGroupData в SDK-модель
func fieldGroupDataToModel(d gkitmodels.FieldGroupData) amomodels.CustomFieldGroup {
	return amomodels.CustomFieldGroup{
		ID:   d.ID,
		Name: d.Name,
		Sort: d.Sort,
	}
}

// lossReasonDataToModel конвертирует LossReasonData в SDK-модель
func lossReasonDataToModel(d gkitmodels.LossReasonData) *amomodels.LossReason {
	return &amomodels.LossReason{
		Name: d.Name,
		Sort: d.Sort,
	}
}

// sourceDataToModel конвертирует SourceData в SDK-модель
func sourceDataToModel(d gkitmodels.SourceData) *amomodels.Source {
	return &amomodels.Source{
		ID:         d.ID,
		Name:       d.Name,
		ExternalID: d.ExternalID,
		OriginCode: d.OriginCode,
		PipelineID: d.PipelineID,
		Default:    d.Default,
	}
}

func (t *AdminSchemaTool) handleCustomFields(ctx context.Context, input gkitmodels.AdminSchemaInput) (any, error) {
	if input.EntityType == "" {
		return nil, fmt.Errorf("entity_type is required for custom_fields")
	}
	switch input.Action {
	case "list":
		filter := buildCustomFieldsFilter(input.Filter)
		result, err := t.service.ListCustomFields(ctx, input.EntityType, filter)
		if err != nil {
			return nil, err
		}
		// client-side фильтрация по имени
		if input.Filter != nil && input.Filter.Name != "" {
			result.Items = filterByName(result.Items, input.Filter.Name, func(f *amomodels.CustomField) string {
				return f.Name
			})
		}
		return result, nil

	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for get")
		}
		return t.service.GetCustomField(ctx, input.EntityType, input.ID)

	case "create":
		fields := collectCustomFields(input)
		if len(fields) == 0 {
			return nil, fmt.Errorf("custom_field or custom_fields is required for create")
		}
		return t.service.CreateCustomFields(ctx, input.EntityType, fields)

	case "update":
		fields := collectCustomFields(input)
		if len(fields) == 0 {
			return nil, fmt.Errorf("custom_field or custom_fields is required for update")
		}
		return t.service.UpdateCustomFields(ctx, input.EntityType, fields)

	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for delete")
		}
		return t.service.DeleteCustomField(ctx, input.EntityType, input.ID)

	default:
		return nil, fmt.Errorf("unknown action for custom_fields: %s", input.Action)
	}
}

// collectCustomFields объединяет одиночное и batch поля из input
func collectCustomFields(input gkitmodels.AdminSchemaInput) []*amomodels.CustomField {
	var fields []*amomodels.CustomField
	if input.CustomField != nil {
		fields = append(fields, customFieldDataToModel(*input.CustomField))
	}
	for _, d := range input.CustomFields {
		fields = append(fields, customFieldDataToModel(d))
	}
	return fields
}

func (t *AdminSchemaTool) handleFieldGroups(ctx context.Context, input gkitmodels.AdminSchemaInput) (any, error) {
	if input.EntityType == "" {
		return nil, fmt.Errorf("entity_type is required for field_groups")
	}
	switch input.Action {
	case "list":
		filter := buildBaseFilter(input.Filter)
		result, err := t.service.ListFieldGroups(ctx, input.EntityType, filter)
		if err != nil {
			return nil, err
		}
		// client-side фильтрация по имени
		if input.Filter != nil && input.Filter.Name != "" {
			result.Items = filterByName(result.Items, input.Filter.Name, func(g amomodels.CustomFieldGroup) string {
				return g.Name
			})
		}
		return result, nil

	case "get":
		if input.GroupID == "" {
			return nil, fmt.Errorf("group_id is required for get")
		}
		return t.service.GetFieldGroup(ctx, input.EntityType, input.GroupID)

	case "create":
		groups := collectFieldGroups(input)
		if len(groups) == 0 {
			return nil, fmt.Errorf("field_group or field_groups is required for create")
		}
		return t.service.CreateFieldGroups(ctx, input.EntityType, groups)

	case "update":
		groups := collectFieldGroups(input)
		if len(groups) == 0 {
			return nil, fmt.Errorf("field_group or field_groups is required for update")
		}
		return t.service.UpdateFieldGroups(ctx, input.EntityType, groups)

	case "delete":
		if input.GroupID == "" {
			return nil, fmt.Errorf("group_id is required for delete")
		}
		return t.service.DeleteFieldGroup(ctx, input.EntityType, input.GroupID)

	default:
		return nil, fmt.Errorf("unknown action for field_groups: %s", input.Action)
	}
}

// collectFieldGroups объединяет одиночную и batch группы из input
func collectFieldGroups(input gkitmodels.AdminSchemaInput) []amomodels.CustomFieldGroup {
	var groups []amomodels.CustomFieldGroup
	if input.FieldGroup != nil {
		groups = append(groups, fieldGroupDataToModel(*input.FieldGroup))
	}
	for _, d := range input.FieldGroups {
		groups = append(groups, fieldGroupDataToModel(d))
	}
	return groups
}

func (t *AdminSchemaTool) handleLossReasons(ctx context.Context, input gkitmodels.AdminSchemaInput) (any, error) {
	switch input.Action {
	case "list":
		filter := buildBaseFilter(input.Filter)
		result, err := t.service.ListLossReasons(ctx, filter)
		if err != nil {
			return nil, err
		}
		// client-side фильтрация по имени
		if input.Filter != nil && input.Filter.Name != "" {
			result.Items = filterByName(result.Items, input.Filter.Name, func(lr *amomodels.LossReason) string {
				return lr.Name
			})
		}
		return result, nil

	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for get")
		}
		return t.service.GetLossReason(ctx, input.ID)

	case "create":
		reasons := collectLossReasons(input)
		if len(reasons) == 0 {
			return nil, fmt.Errorf("loss_reason or loss_reasons is required for create")
		}
		return t.service.CreateLossReasons(ctx, reasons)

	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for delete")
		}
		return t.service.DeleteLossReason(ctx, input.ID)

	default:
		return nil, fmt.Errorf("unknown action for loss_reasons: %s", input.Action)
	}
}

// collectLossReasons объединяет одиночную и batch причины из input
func collectLossReasons(input gkitmodels.AdminSchemaInput) []*amomodels.LossReason {
	var reasons []*amomodels.LossReason
	if input.LossReason != nil {
		reasons = append(reasons, lossReasonDataToModel(*input.LossReason))
	}
	for _, d := range input.LossReasons {
		reasons = append(reasons, lossReasonDataToModel(d))
	}
	return reasons
}

func (t *AdminSchemaTool) handleSources(ctx context.Context, input gkitmodels.AdminSchemaInput) (any, error) {
	switch input.Action {
	case "list":
		filter := buildSourcesFilter(input.Filter)
		result, err := t.service.ListSources(ctx, filter)
		if err != nil {
			return nil, err
		}
		// client-side фильтрация по имени
		if input.Filter != nil && input.Filter.Name != "" {
			result.Items = filterByName(result.Items, input.Filter.Name, func(s *amomodels.Source) string {
				return s.Name
			})
		}
		return result, nil

	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for get")
		}
		return t.service.GetSource(ctx, input.ID)

	case "create":
		sources := collectSources(input)
		if len(sources) == 0 {
			return nil, fmt.Errorf("source or sources is required for create")
		}
		return t.service.CreateSources(ctx, sources)

	case "update":
		sources := collectSources(input)
		if len(sources) == 0 {
			return nil, fmt.Errorf("source or sources is required for update")
		}
		return t.service.UpdateSources(ctx, sources)

	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for delete")
		}
		return t.service.DeleteSource(ctx, input.ID)

	default:
		return nil, fmt.Errorf("unknown action for sources: %s", input.Action)
	}
}

// collectSources объединяет одиночный и batch источники из input
func collectSources(input gkitmodels.AdminSchemaInput) []*amomodels.Source {
	var sources []*amomodels.Source
	if input.Source != nil {
		sources = append(sources, sourceDataToModel(*input.Source))
	}
	for _, d := range input.Sources {
		sources = append(sources, sourceDataToModel(d))
	}
	return sources
}
