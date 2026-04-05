package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
	admin_integrations "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/admin_integrations"
)

// AdminIntegrationsTool реализует нативный ADK FunctionTool интерфейс для управления
// интеграциями amoCRM: вебхуки, виджеты, кнопки сайта, шаблоны чата, короткие ссылки.
type AdminIntegrationsTool struct {
	service admin_integrations.Service
}

// NewAdminIntegrationsTool создаёт новый AdminIntegrationsTool с указанным сервисом.
func NewAdminIntegrationsTool(service admin_integrations.Service) *AdminIntegrationsTool {
	return &AdminIntegrationsTool{service: service}
}

// Name implements tool.Tool.
func (t *AdminIntegrationsTool) Name() string {
	return "admin_integrations"
}

// Description implements tool.Tool.
func (t *AdminIntegrationsTool) Description() string {
	return "Управление вебхуками, виджетами, кнопками сайта, шаблонами чата и короткими ссылками amoCRM. " +
		"Layers: webhooks, widgets, website_buttons, chat_templates, short_links. " +
		"Вызови с {\"layer\": \"<layer>\", \"action\": \"<action>\"} чтобы получить схему параметров."
}

// IsLongRunning implements tool.Tool.
func (t *AdminIntegrationsTool) IsLongRunning() bool {
	return false
}

// ProcessRequest реализует toolinternal.RequestProcessor — регистрирует Declaration в LLM request.
func (t *AdminIntegrationsTool) ProcessRequest(_ tool.Context, req *model.LLMRequest) error {
	return packToolDeclaration(req, t)
}

// Declaration implements toolinternal.FunctionTool (duck typing).
func (t *AdminIntegrationsTool) Declaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"layer": {
					Type:        genai.TypeString,
					Description: "Слой: webhooks, widgets, website_buttons, chat_templates, short_links",
					Enum:        []string{"webhooks", "widgets", "website_buttons", "chat_templates", "short_links"},
				},
				"action": {
					Type:        genai.TypeString,
					Description: "Действие для выбранного layer (list, get, create, update, delete, subscribe, unsubscribe, install, uninstall, add_chat, send_review, update_review, delete_many)",
				},
			},
			Required: []string{"layer", "action"},
		},
	}
}

// Run implements toolinternal.FunctionTool (duck typing).
func (t *AdminIntegrationsTool) Run(ctx tool.Context, args any) (map[string]any, error) {
	m, ok := args.(map[string]any)
	if !ok {
		b, err := json.Marshal(args)
		if err != nil {
			return toResultMap(adminIntegrationsSchema("", ""))
		}
		if err := json.Unmarshal(b, &m); err != nil {
			return toResultMap(adminIntegrationsSchema("", ""))
		}
	}

	layer, _ := m["layer"].(string)
	action, _ := m["action"].(string)

	if layer == "" || action == "" {
		return toResultMap(adminIntegrationsSchema(layer, action))
	}

	// Schema mode: обязательные поля для данного action отсутствуют
	if adminIntegrationsIsSchemaMode(layer, action, m) {
		return toResultMap(adminIntegrationsSchema(layer, action))
	}

	// Execute mode: JSON roundtrip map → AdminIntegrationsInput → существующий handler
	b, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("admin_integrations: marshal input: %w", err)
	}
	var input gkitmodels.AdminIntegrationsInput
	if err := json.Unmarshal(b, &input); err != nil {
		return nil, fmt.Errorf("admin_integrations: unmarshal input: %w", err)
	}

	var result any
	switch input.Layer {
	case "webhooks":
		result, err = t.handleWebhooks(ctx, input)
	case "widgets":
		result, err = t.handleWidgets(ctx, input)
	case "website_buttons":
		result, err = t.handleWebsiteButtons(ctx, input)
	case "chat_templates":
		result, err = t.handleChatTemplates(ctx, input)
	case "short_links":
		result, err = t.handleShortLinks(ctx, input)
	default:
		return nil, fmt.Errorf("unknown layer: %s", input.Layer)
	}
	if err != nil {
		return nil, err
	}
	return toResultMap(result)
}

// adminIntegrationsSchema — полная схема полей, возвращаемая в schema mode.
// Не содержит dynamic available_values — admin tools работают с именами и ID напрямую.
func adminIntegrationsSchema(layer, action string) map[string]any {
	base := map[string]any{
		"schema": true,
		"tool":   "admin_integrations",
		"layer":  layer,
		"action": action,
	}

	switch layer {
	case "webhooks":
		return adminIntegrationsWebhooksSchema(base, action)
	case "widgets":
		return adminIntegrationsWidgetsSchema(base, action)
	case "website_buttons":
		return adminIntegrationsWebsiteButtonsSchema(base, action)
	case "chat_templates":
		return adminIntegrationsChatTemplatesSchema(base, action)
	case "short_links":
		return adminIntegrationsShortLinksSchema(base, action)
	default:
		base["description"] = "Неизвестный layer. Доступные: webhooks, widgets, website_buttons, chat_templates, short_links."
		base["layers"] = []string{"webhooks", "widgets", "website_buttons", "chat_templates", "short_links"}
		return base
	}
}

func adminIntegrationsWebhooksSchema(base map[string]any, action string) map[string]any {
	switch action {
	case "search", "list":
		base["description"] = "Получить список вебхуков."
		base["optional_fields"] = map[string]any{
			"filter": map[string]any{
				"type":        "object",
				"description": "Фильтр",
				"fields": map[string]any{
					"destination": map[string]any{"type": "string", "description": "Фильтр по URL вебхука"},
				},
			},
		}
		base["example"] = map[string]any{"layer": "webhooks", "action": "list"}
	case "subscribe":
		base["description"] = "Подписаться на вебхук."
		base["required_fields"] = map[string]any{
			"destination": map[string]any{"type": "string", "description": "URL вебхука"},
		}
		base["optional_fields"] = map[string]any{
			"event_types": map[string]any{
				"type":        "array",
				"items":       "string",
				"description": "События: add_lead, update_lead, delete_lead, restore_lead, status_lead, add_contact, update_contact, delete_contact, restore_contact, add_company, update_company, delete_company, restore_company, add_task, update_task, delete_task, complete_task, add_note, update_note, delete_note",
			},
		}
		base["example"] = map[string]any{
			"layer":       "webhooks",
			"action":      "subscribe",
			"destination": "https://example.com/webhook",
			"event_types": []string{"add_lead", "update_lead"},
		}
	case "unsubscribe":
		base["description"] = "Отписаться от вебхука."
		base["required_fields"] = map[string]any{
			"destination": map[string]any{"type": "string", "description": "URL вебхука"},
		}
		base["optional_fields"] = map[string]any{
			"event_types": map[string]any{"type": "array", "items": "string", "description": "События для отписки (пусто = все)"},
		}
		base["example"] = map[string]any{
			"layer":       "webhooks",
			"action":      "unsubscribe",
			"destination": "https://example.com/webhook",
		}
	default:
		base["description"] = "Доступные actions для webhooks: list, subscribe, unsubscribe."
		base["actions"] = []string{"list", "subscribe", "unsubscribe"}
	}
	return base
}

func adminIntegrationsWidgetsSchema(base map[string]any, action string) map[string]any {
	switch action {
	case "search", "list":
		base["description"] = "Получить список виджетов."
		base["optional_fields"] = map[string]any{
			"filter": map[string]any{
				"type":   "object",
				"fields": map[string]any{"limit": map[string]any{"type": "integer"}, "page": map[string]any{"type": "integer"}},
			},
		}
		base["example"] = map[string]any{"layer": "widgets", "action": "list"}
	case "get":
		base["description"] = "Получить виджет по коду."
		base["required_fields"] = map[string]any{
			"code": map[string]any{"type": "string", "description": "Код виджета"},
		}
		base["example"] = map[string]any{"layer": "widgets", "action": "get", "code": "widget_code"}
	case "install":
		base["description"] = "Установить виджет."
		base["required_fields"] = map[string]any{
			"code": map[string]any{"type": "string", "description": "Код виджета"},
		}
		base["optional_fields"] = map[string]any{
			"settings": map[string]any{"type": "object", "description": "Настройки виджета (ключ-значение)"},
		}
		base["example"] = map[string]any{"layer": "widgets", "action": "install", "code": "widget_code"}
	case "uninstall":
		base["description"] = "Удалить виджет."
		base["required_fields"] = map[string]any{
			"code": map[string]any{"type": "string", "description": "Код виджета"},
		}
		base["example"] = map[string]any{"layer": "widgets", "action": "uninstall", "code": "widget_code"}
	default:
		base["description"] = "Доступные actions для widgets: list, get, install, uninstall."
		base["actions"] = []string{"list", "get", "install", "uninstall"}
	}
	return base
}

func adminIntegrationsWebsiteButtonsSchema(base map[string]any, action string) map[string]any {
	switch action {
	case "search", "list":
		base["description"] = "Получить список кнопок на сайте."
		base["optional_fields"] = map[string]any{
			"filter": map[string]any{"type": "object", "fields": map[string]any{"limit": map[string]any{"type": "integer"}, "page": map[string]any{"type": "integer"}}},
			"with":   map[string]any{"type": "array", "items": "string", "description": "Обогащение: 'scripts' (скрипты), 'deleted' (удалённые)"},
		}
		base["example"] = map[string]any{"layer": "website_buttons", "action": "list"}
	case "get":
		base["description"] = "Получить кнопку по source_id."
		base["required_fields"] = map[string]any{
			"id": map[string]any{"type": "integer", "description": "source_id кнопки"},
		}
		base["optional_fields"] = map[string]any{
			"with": map[string]any{"type": "array", "items": "string", "description": "'scripts', 'deleted'"},
		}
		base["example"] = map[string]any{"layer": "website_buttons", "action": "get", "id": 123}
	case "create":
		base["description"] = "Создать кнопку на сайте."
		base["required_fields"] = map[string]any{
			"website_button": map[string]any{
				"type":        "object",
				"description": "Данные кнопки",
				"required_fields": map[string]any{
					"name": map[string]any{"type": "string", "description": "Название кнопки"},
				},
				"optional_fields": map[string]any{
					"pipeline_id":                    map[string]any{"type": "integer", "description": "ID воронки (получить через admin_pipelines: list)"},
					"trusted_websites":               map[string]any{"type": "array", "items": "string", "description": "Доверенные сайты"},
					"is_duplication_control_enabled": map[string]any{"type": "boolean", "description": "Контроль дублей"},
				},
			},
		}
		base["example"] = map[string]any{
			"layer":  "website_buttons",
			"action": "create",
			"website_button": map[string]any{
				"name":             "Кнопка на сайте",
				"trusted_websites": []string{"example.com"},
			},
		}
	case "update":
		base["description"] = "Обновить кнопку на сайте."
		base["required_fields"] = map[string]any{
			"id": map[string]any{"type": "integer", "description": "source_id кнопки"},
		}
		base["optional_fields"] = map[string]any{
			"website_button": map[string]any{
				"type": "object",
				"fields": map[string]any{
					"name":                           map[string]any{"type": "string"},
					"pipeline_id":                    map[string]any{"type": "integer"},
					"trusted_websites":               map[string]any{"type": "array", "items": "string"},
					"is_duplication_control_enabled": map[string]any{"type": "boolean"},
				},
			},
		}
		base["example"] = map[string]any{
			"layer": "website_buttons", "action": "update", "id": 123,
			"website_button": map[string]any{"name": "Новое название"},
		}
	case "add_chat":
		base["description"] = "Добавить онлайн-чат к кнопке."
		base["required_fields"] = map[string]any{
			"id": map[string]any{"type": "integer", "description": "source_id кнопки"},
		}
		base["example"] = map[string]any{"layer": "website_buttons", "action": "add_chat", "id": 123}
	default:
		base["description"] = "Доступные actions для website_buttons: list, get, create, update, add_chat."
		base["actions"] = []string{"list", "get", "create", "update", "add_chat"}
	}
	return base
}

func adminIntegrationsChatTemplatesSchema(base map[string]any, action string) map[string]any {
	switch action {
	case "search", "list":
		base["description"] = "Получить список шаблонов чата."
		base["optional_fields"] = map[string]any{
			"filter": map[string]any{
				"type": "object",
				"fields": map[string]any{
					"limit":         map[string]any{"type": "integer"},
					"page":          map[string]any{"type": "integer"},
					"external_ids":  map[string]any{"type": "array", "items": "string", "description": "Внешние ID"},
					"template_type": map[string]any{"type": "string", "description": "Тип шаблона: amocrm или waba"},
				},
			},
		}
		base["example"] = map[string]any{"layer": "chat_templates", "action": "list"}
	case "create":
		base["description"] = "Создать шаблон чата."
		base["required_fields"] = map[string]any{
			"chat_template": map[string]any{
				"type": "object",
				"required_fields": map[string]any{
					"name":    map[string]any{"type": "string", "description": "Название шаблона"},
					"content": map[string]any{"type": "string", "description": "Текст шаблона"},
					"type":    map[string]any{"type": "string", "description": "Тип: amocrm или waba"},
				},
				"optional_fields": map[string]any{
					"external_id":   map[string]any{"type": "string", "description": "Внешний ID"},
					"is_editable":   map[string]any{"type": "boolean"},
					"waba_header":   map[string]any{"type": "string", "description": "Заголовок WABA"},
					"waba_footer":   map[string]any{"type": "string", "description": "Подвал WABA"},
					"waba_category": map[string]any{"type": "string", "description": "UTILITY, AUTHENTICATION, MARKETING"},
					"waba_language": map[string]any{"type": "string", "description": "Язык (ru, en)"},
				},
			},
		}
		base["example"] = map[string]any{
			"layer": "chat_templates", "action": "create",
			"chat_template": map[string]any{"name": "Приветствие", "content": "Здравствуйте!", "type": "amocrm"},
		}
	case "update":
		base["description"] = "Обновить шаблон чата."
		base["required_fields"] = map[string]any{
			"id":            map[string]any{"type": "integer", "description": "ID шаблона"},
			"chat_template": map[string]any{"type": "object", "description": "Поля для обновления (те же что в create)"},
		}
		base["example"] = map[string]any{
			"layer": "chat_templates", "action": "update", "id": 456,
			"chat_template": map[string]any{"name": "Новое название"},
		}
	case "delete":
		base["description"] = "Удалить шаблон чата."
		base["required_fields"] = map[string]any{
			"id": map[string]any{"type": "integer", "description": "ID шаблона"},
		}
		base["example"] = map[string]any{"layer": "chat_templates", "action": "delete", "id": 456}
	case "delete_many":
		base["description"] = "Удалить несколько шаблонов чата."
		base["required_fields"] = map[string]any{
			"ids": map[string]any{"type": "array", "items": "integer", "description": "Массив ID шаблонов"},
		}
		base["example"] = map[string]any{"layer": "chat_templates", "action": "delete_many", "ids": []int{1, 2, 3}}
	case "send_review":
		base["description"] = "Отправить шаблон на ревью."
		base["required_fields"] = map[string]any{
			"id": map[string]any{"type": "integer", "description": "ID шаблона"},
		}
		base["example"] = map[string]any{"layer": "chat_templates", "action": "send_review", "id": 456}
	case "update_review":
		base["description"] = "Обновить статус ревью шаблона."
		base["required_fields"] = map[string]any{
			"id":            map[string]any{"type": "integer", "description": "ID шаблона"},
			"review_id":     map[string]any{"type": "integer", "description": "ID ревью (из ответа send_review)"},
			"review_status": map[string]any{"type": "string", "description": "Новый статус ревью"},
		}
		base["example"] = map[string]any{
			"layer": "chat_templates", "action": "update_review",
			"id": 456, "review_id": 789, "review_status": "approved",
		}
	default:
		base["description"] = "Доступные actions для chat_templates: list, create, update, delete, delete_many, send_review, update_review."
		base["actions"] = []string{"list", "create", "update", "delete", "delete_many", "send_review", "update_review"}
	}
	return base
}

func adminIntegrationsShortLinksSchema(base map[string]any, action string) map[string]any {
	switch action {
	case "search", "list":
		base["description"] = "Получить список коротких ссылок."
		base["optional_fields"] = map[string]any{
			"filter": map[string]any{"type": "object", "fields": map[string]any{"limit": map[string]any{"type": "integer"}, "page": map[string]any{"type": "integer"}}},
		}
		base["example"] = map[string]any{"layer": "short_links", "action": "list"}
	case "create":
		base["description"] = "Создать короткую ссылку (одну или батч)."
		base["required_fields"] = map[string]any{
			"url": map[string]any{"type": "string", "description": "URL для сокращения (для одной ссылки)"},
		}
		base["optional_fields"] = map[string]any{
			"urls":        map[string]any{"type": "array", "items": "string", "description": "Массив URL (батч-создание; если указан, url/entity_id/entity_type игнорируются)"},
			"entity_id":   map[string]any{"type": "integer", "description": "ID сущности для привязки"},
			"entity_type": map[string]any{"type": "string", "description": "Тип: leads, contacts, companies, customers"},
		}
		base["example"] = map[string]any{
			"layer": "short_links", "action": "create",
			"url": "https://example.com/very/long/url",
		}
	case "delete":
		base["description"] = "Удалить короткую ссылку."
		base["required_fields"] = map[string]any{
			"id": map[string]any{"type": "integer", "description": "ID ссылки"},
		}
		base["example"] = map[string]any{"layer": "short_links", "action": "delete", "id": 123}
	default:
		base["description"] = "Доступные actions для short_links: list, create, delete."
		base["actions"] = []string{"list", "create", "delete"}
	}
	return base
}

// isSchemaMode определяет, нужно ли вернуть схему вместо выполнения.
// Логика: проверяем наличие обязательных полей для конкретного layer+action.
func adminIntegrationsIsSchemaMode(layer, action string, m map[string]any) bool {
	strVal := func(key string) string {
		v, _ := m[key].(string)
		return v
	}
	intVal := func(key string) int {
		switch v := m[key].(type) {
		case float64:
			return int(v)
		case int:
			return v
		}
		return 0
	}
	hasKey := func(key string) bool {
		v, ok := m[key]
		return ok && v != nil
	}

	switch layer {
	case "webhooks":
		switch action {
		case "subscribe", "unsubscribe":
			return strVal("destination") == ""
		}
	case "widgets":
		switch action {
		case "get", "install", "uninstall":
			return strVal("code") == ""
		}
	case "website_buttons":
		switch action {
		case "get", "update", "add_chat":
			return intVal("id") == 0
		case "create":
			return !hasKey("website_button")
		}
	case "chat_templates":
		switch action {
		case "create":
			return !hasKey("chat_template")
		case "update":
			return intVal("id") == 0 || !hasKey("chat_template")
		case "delete", "send_review":
			return intVal("id") == 0
		case "delete_many":
			ids, _ := m["ids"].([]any)
			return len(ids) == 0
		case "update_review":
			return intVal("id") == 0 || intVal("review_id") == 0 || strVal("review_status") == ""
		}
	case "short_links":
		switch action {
		case "create":
			// Батч-режим через urls не требует url
			if urls, ok := m["urls"].([]any); ok && len(urls) > 0 {
				return false
			}
			return strVal("url") == ""
		case "delete":
			return intVal("id") == 0
		}
	}
	// list/search и неизвестные action — не schema mode (выполняем или отдаём ошибку)
	return false
}

// unixToRFC3339 конвертирует Unix timestamp в строку RFC3339. Возвращает пустую строку если ts == 0.
func unixToRFC3339(ts int64) string {
	if ts == 0 {
		return ""
	}
	return time.Unix(ts, 0).UTC().Format(time.RFC3339)
}

func (t *AdminIntegrationsTool) handleWebhooks(ctx context.Context, input gkitmodels.AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "search", "list":
		var filter *filters.WebhooksFilter
		if input.Filter != nil {
			filter = filters.NewWebhooksFilter()
			if input.Filter.Destination != "" {
				filter.SetDestination(input.Filter.Destination)
			}
		}
		webhooks, err := t.service.ListWebhooks(ctx, filter)
		if err != nil {
			return nil, err
		}
		// Форматируем timestamps, убираем AccountID
		type webhookOut struct {
			ID          int      `json:"id,omitempty"`
			Destination string   `json:"destination,omitempty"`
			Settings    []string `json:"settings,omitempty"`
			CreatedAt   string   `json:"created_at,omitempty"`
			UpdatedAt   string   `json:"updated_at,omitempty"`
			Disabled    bool     `json:"disabled,omitempty"`
		}
		out := make([]webhookOut, 0, len(webhooks))
		for _, w := range webhooks {
			out = append(out, webhookOut{
				ID:          w.ID,
				Destination: w.Destination,
				Settings:    w.Settings,
				CreatedAt:   unixToRFC3339(w.CreatedAt),
				UpdatedAt:   unixToRFC3339(w.UpdatedAt),
				Disabled:    w.Disabled,
			})
		}
		return out, nil
	case "subscribe":
		dest := input.Destination
		if dest == "" {
			return nil, fmt.Errorf("destination is required for subscribe")
		}
		return t.service.SubscribeWebhook(ctx, dest, input.EventTypes)
	case "unsubscribe":
		dest := input.Destination
		if dest == "" {
			return nil, fmt.Errorf("destination is required for unsubscribe")
		}
		return nil, t.service.UnsubscribeWebhook(ctx, dest, input.EventTypes)
	default:
		return nil, fmt.Errorf("unknown action for webhooks: %s", input.Action)
	}
}

func (t *AdminIntegrationsTool) handleWidgets(ctx context.Context, input gkitmodels.AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "search", "list":
		var filter *filters.WidgetsFilter
		if input.Filter != nil {
			filter = filters.NewWidgetsFilter()
			if input.Filter.Limit > 0 {
				filter.SetLimit(input.Filter.Limit)
			}
			if input.Filter.Page > 0 {
				filter.SetPage(input.Filter.Page)
			}
		}
		return t.service.ListWidgets(ctx, filter)
	case "get":
		if input.Code == "" {
			return nil, fmt.Errorf("code is required for get widget")
		}
		return t.service.GetWidget(ctx, input.Code)
	case "install":
		if input.Code == "" {
			return nil, fmt.Errorf("code is required for install widget")
		}
		return t.service.InstallWidget(ctx, input.Code, input.Settings)
	case "uninstall":
		if input.Code == "" {
			return nil, fmt.Errorf("code is required for uninstall widget")
		}
		return nil, t.service.UninstallWidget(ctx, input.Code)
	default:
		return nil, fmt.Errorf("unknown action for widgets: %s", input.Action)
	}
}

func (t *AdminIntegrationsTool) handleWebsiteButtons(ctx context.Context, input gkitmodels.AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "search", "list":
		var filter *services.WebsiteButtonsFilter
		if input.Filter != nil {
			filter = &services.WebsiteButtonsFilter{}
			if input.Filter.Limit > 0 {
				filter.Limit = input.Filter.Limit
			}
			if input.Filter.Page > 0 {
				filter.Page = input.Filter.Page
			}
		}
		buttons, err := t.service.ListWebsiteButtons(ctx, filter, input.With)
		if err != nil {
			return nil, err
		}
		return stripWebsiteButtonsAccountID(buttons), nil
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for get website button")
		}
		button, err := t.service.GetWebsiteButton(ctx, input.ID, input.With)
		if err != nil {
			return nil, err
		}
		if button == nil {
			return nil, nil
		}
		return stripWebsiteButtonAccountID(button), nil
	case "create":
		if input.WebsiteButton == nil {
			return nil, fmt.Errorf("website_button is required for create")
		}
		req := &amomodels.WebsiteButtonCreateRequest{
			Name:            input.WebsiteButton.Name,
			TrustedWebsites: input.WebsiteButton.TrustedWebsites,
		}
		if input.WebsiteButton.PipelineID != nil {
			req.PipelineID = *input.WebsiteButton.PipelineID
		}
		if input.WebsiteButton.IsDuplicationControlEnabled != nil {
			req.IsDuplicationControlEnabled = *input.WebsiteButton.IsDuplicationControlEnabled
		}
		return t.service.CreateWebsiteButton(ctx, req)
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id (source_id) is required for update website button")
		}
		req := &amomodels.WebsiteButtonUpdateRequest{
			SourceID: input.ID, // BUG FIX: заполняем SourceID из input.ID
		}
		if input.WebsiteButton != nil {
			req.Name = input.WebsiteButton.Name
			req.TrustedWebsites = input.WebsiteButton.TrustedWebsites
			req.PipelineID = input.WebsiteButton.PipelineID
			req.IsDuplicationControlEnabled = input.WebsiteButton.IsDuplicationControlEnabled
		}
		button, err := t.service.UpdateWebsiteButton(ctx, req)
		if err != nil {
			return nil, err
		}
		if button == nil {
			return nil, nil
		}
		return stripWebsiteButtonAccountID(button), nil
	case "add_chat":
		if input.ID == 0 {
			return nil, fmt.Errorf("id (source_id) is required for add_chat")
		}
		return nil, t.service.AddOnlineChat(ctx, input.ID)
	default:
		return nil, fmt.Errorf("unknown action for website_buttons: %s", input.Action)
	}
}

// websiteButtonOut — ответная структура без AccountID
type websiteButtonOut struct {
	SourceID                    int    `json:"source_id,omitempty"`
	ButtonID                    *int   `json:"button_id,omitempty"`
	Name                        string `json:"name,omitempty"`
	PipelineID                  *int   `json:"pipeline_id,omitempty"`
	IsDuplicationControlEnabled bool   `json:"is_duplication_control_enabled,omitempty"`
	CreationStatus              string `json:"creation_status,omitempty"`
	Script                      string `json:"script,omitempty"`
	IsDeleted                   bool   `json:"is_deleted,omitempty"`
}

func stripWebsiteButtonAccountID(b *amomodels.WebsiteButton) websiteButtonOut {
	return websiteButtonOut{
		SourceID:                    b.SourceID,
		ButtonID:                    b.ButtonID,
		Name:                        b.Name,
		PipelineID:                  b.PipelineID,
		IsDuplicationControlEnabled: b.IsDuplicationControlEnabled,
		CreationStatus:              b.CreationStatus,
		Script:                      b.Script,
		IsDeleted:                   b.IsDeleted,
	}
}

func stripWebsiteButtonsAccountID(buttons []*amomodels.WebsiteButton) []websiteButtonOut {
	out := make([]websiteButtonOut, 0, len(buttons))
	for _, b := range buttons {
		if b != nil {
			out = append(out, stripWebsiteButtonAccountID(b))
		}
	}
	return out
}

func (t *AdminIntegrationsTool) handleChatTemplates(ctx context.Context, input gkitmodels.AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "search", "list":
		var filter *filters.TemplatesFilter
		if input.Filter != nil {
			filter = filters.NewTemplatesFilter()
			if input.Filter.Limit > 0 {
				filter.SetLimit(input.Filter.Limit)
			}
			if input.Filter.Page > 0 {
				filter.SetPage(input.Filter.Page)
			}
			if len(input.Filter.ExternalIDs) > 0 {
				filter.SetExternalIDs(input.Filter.ExternalIDs)
			}
		}
		templates, err := t.service.ListChatTemplates(ctx, filter)
		if err != nil {
			return nil, err
		}
		// Фильтрация по TemplateType на стороне клиента (API не поддерживает)
		if input.Filter != nil && input.Filter.TemplateType != "" {
			filtered := templates[:0]
			for _, tmpl := range templates {
				if string(tmpl.Type) == input.Filter.TemplateType {
					filtered = append(filtered, tmpl)
				}
			}
			templates = filtered
		}
		return formatChatTemplates(templates), nil
	case "create":
		if input.ChatTemplate == nil {
			return nil, fmt.Errorf("chat_template is required for create")
		}
		tmpl := chatTemplateDataToModel(input.ChatTemplate)
		result, err := t.service.CreateChatTemplate(ctx, tmpl)
		if err != nil {
			return nil, err
		}
		return formatChatTemplate(result), nil
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for update chat template")
		}
		if input.ChatTemplate == nil {
			return nil, fmt.Errorf("chat_template is required for update")
		}
		tmpl := chatTemplateDataToModel(input.ChatTemplate)
		tmpl.ID = input.ID
		result, err := t.service.UpdateChatTemplate(ctx, tmpl)
		if err != nil {
			return nil, err
		}
		return formatChatTemplate(result), nil
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for delete chat template")
		}
		return nil, t.service.DeleteChatTemplate(ctx, input.ID)
	case "delete_many":
		if len(input.IDs) == 0 {
			return nil, fmt.Errorf("ids are required for delete_many")
		}
		return nil, t.service.DeleteChatTemplates(ctx, input.IDs)
	case "send_review":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for send chat template on review")
		}
		reviews, err := t.service.SendChatTemplateOnReview(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return formatChatTemplateReviews(reviews), nil
	case "update_review":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for update chat template review status")
		}
		if input.ReviewID == 0 {
			return nil, fmt.Errorf("review_id is required for update_review (получить из ответа send_review)")
		}
		if input.ReviewStatus == "" {
			return nil, fmt.Errorf("review_status is required for update_review")
		}
		review, err := t.service.UpdateChatTemplateReviewStatus(ctx, input.ID, input.ReviewID, input.ReviewStatus)
		if err != nil {
			return nil, err
		}
		return formatChatTemplateReview(review), nil
	default:
		return nil, fmt.Errorf("unknown action for chat_templates: %s", input.Action)
	}
}

// chatTemplateOut — ответная структура без AccountID, с форматированными timestamps
type chatTemplateOut struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Content    string `json:"content,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
	Type       string `json:"type,omitempty"`
	IsEditable bool   `json:"is_editable,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	WabaHeader string `json:"waba_header,omitempty"`
	WabaFooter string `json:"waba_footer,omitempty"`
}

// chatTemplateReviewOut — ответная структура ревью с форматированным timestamp
type chatTemplateReviewOut struct {
	ID        int    `json:"id,omitempty"`
	Status    string `json:"status,omitempty"`
	Reason    string `json:"reason,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

func formatChatTemplate(tmpl *amomodels.ChatTemplate) chatTemplateOut {
	if tmpl == nil {
		return chatTemplateOut{}
	}
	return chatTemplateOut{
		ID:         tmpl.ID,
		Name:       tmpl.Name,
		Content:    tmpl.Content,
		ExternalID: tmpl.ExternalID,
		Type:       string(tmpl.Type),
		IsEditable: tmpl.IsEditable,
		CreatedAt:  unixToRFC3339(tmpl.CreatedAt),
		UpdatedAt:  unixToRFC3339(tmpl.UpdatedAt),
		WabaHeader: tmpl.WabaHeader,
		WabaFooter: tmpl.WabaFooter,
	}
}

func formatChatTemplates(templates []*amomodels.ChatTemplate) []chatTemplateOut {
	out := make([]chatTemplateOut, 0, len(templates))
	for _, tmpl := range templates {
		out = append(out, formatChatTemplate(tmpl))
	}
	return out
}

func formatChatTemplateReview(rev *amomodels.ChatTemplateReview) chatTemplateReviewOut {
	if rev == nil {
		return chatTemplateReviewOut{}
	}
	return chatTemplateReviewOut{
		ID:        rev.ID,
		Status:    rev.Status,
		Reason:    rev.Reason,
		CreatedAt: unixToRFC3339(rev.CreatedAt),
	}
}

func formatChatTemplateReviews(reviews []amomodels.ChatTemplateReview) []chatTemplateReviewOut {
	out := make([]chatTemplateReviewOut, 0, len(reviews))
	for i := range reviews {
		out = append(out, formatChatTemplateReview(&reviews[i]))
	}
	return out
}

func chatTemplateDataToModel(d *gkitmodels.ChatTemplateData) *amomodels.ChatTemplate {
	return &amomodels.ChatTemplate{
		Name:         d.Name,
		Content:      d.Content,
		ExternalID:   d.ExternalID,
		Type:         amomodels.ChatTemplateType(d.Type),
		IsEditable:   d.IsEditable,
		WabaHeader:   d.WabaHeader,
		WabaFooter:   d.WabaFooter,
		WabaCategory: amomodels.ChatTemplateCategory(d.WabaCategory),
		WabaLanguage: d.WabaLanguage,
	}
}

func (t *AdminIntegrationsTool) handleShortLinks(ctx context.Context, input gkitmodels.AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "search", "list":
		var filter *filters.ShortLinksFilter
		if input.Filter != nil {
			filter = filters.NewShortLinksFilter()
			if input.Filter.Limit > 0 {
				filter.SetLimit(input.Filter.Limit)
			}
			if input.Filter.Page > 0 {
				filter.SetPage(input.Filter.Page)
			}
		}
		return t.service.ListShortLinks(ctx, filter)
	case "create":
		if len(input.URLs) > 0 {
			// Батч-создание: entity_id/entity_type не поддерживается для батча
			links := make([]amomodels.ShortLink, 0, len(input.URLs))
			for _, u := range input.URLs {
				links = append(links, amomodels.ShortLink{URL: u})
			}
			return t.service.CreateShortLinks(ctx, links)
		}
		// Одиночное создание
		u := input.URL
		if u == "" {
			return nil, fmt.Errorf("url is required for create short link")
		}
		link := amomodels.ShortLink{
			URL:        u,
			EntityID:   input.EntityID,
			EntityType: input.EntityType,
		}
		return t.service.CreateShortLink(ctx, link)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for delete short link")
		}
		return nil, t.service.DeleteShortLink(ctx, input.ID)
	default:
		return nil, fmt.Errorf("unknown action for short_links: %s", input.Action)
	}
}
