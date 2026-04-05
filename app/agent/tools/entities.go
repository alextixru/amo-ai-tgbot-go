package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/entities"
	toolmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// EntitiesTool — нативный ADK tool для работы с основными сущностями amoCRM (Shadow Tool паттерн).
type EntitiesTool struct {
	service entities.Service
}

// NewEntitiesTool создаёт новый экземпляр EntitiesTool.
func NewEntitiesTool(service entities.Service) *EntitiesTool {
	return &EntitiesTool{service: service}
}

// Name реализует tool.Tool.
func (t *EntitiesTool) Name() string {
	return "entities"
}

// Description реализует tool.Tool.
func (t *EntitiesTool) Description() string {
	return "CRUD для сделок, контактов, компаний amoCRM. Actions: search, get, create, update, sync, link, unlink. Вызови с entity_type + action чтобы получить схему параметров."
}

// IsLongRunning реализует tool.Tool.
func (t *EntitiesTool) IsLongRunning() bool {
	return false
}

// ProcessRequest реализует toolinternal.RequestProcessor — регистрирует Declaration в LLM request.
func (t *EntitiesTool) ProcessRequest(_ tool.Context, req *model.LLMRequest) error {
	return packToolDeclaration(req, t)
}

// Declaration реализует toolinternal.FunctionTool (duck typing).
func (t *EntitiesTool) Declaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"entity_type": {
					Type:        genai.TypeString,
					Description: "Тип сущности: leads, contacts, companies",
					Enum:        []string{"leads", "contacts", "companies"},
				},
				"action": {
					Type:        genai.TypeString,
					Description: "Действие: search, get, create, update, sync, link, unlink, get_chats, link_chats",
					Enum:        []string{"search", "get", "create", "update", "sync", "link", "unlink", "get_chats", "link_chats"},
				},
			},
			Required: []string{"entity_type", "action"},
		},
	}
}

// Run реализует toolinternal.FunctionTool (duck typing).
func (t *EntitiesTool) Run(ctx tool.Context, args any) (map[string]any, error) {
	log.Printf("[entities] handler called, args type=%T value=%v", args, args)
	m, ok := args.(map[string]any)
	if !ok {
		log.Printf("[entities] args is not map[string]any, got %T", args)
		return nil, fmt.Errorf("entities: неверный формат input")
	}

	entityType, _ := m["entity_type"].(string)
	action, _ := m["action"].(string)

	if entityType == "" {
		return nil, fmt.Errorf("entities: entity_type обязателен (leads, contacts, companies)")
	}
	if action == "" {
		return nil, fmt.Errorf("entities: action обязателен")
	}

	if t.entitiesIsSchemaMode(action, m) {
		resp := t.entitiesBuildSchemaResponse(entityType, action)
		b, _ := json.Marshal(resp)
		log.Printf("[entities] schema mode, action=%s entity_type=%s, response size=%d bytes", action, entityType, len(b))
		return resp, nil
	}

	result, err := t.entitiesExecute(ctx, m, entityType, action)
	if err != nil {
		return nil, err
	}
	return toResultMap(result)
}

// entitiesRequiredFields определяет обязательные поля (помимо entity_type/action) для каждого action.
// Если список пустой — action всегда попадает в Execute mode (например, search).
var entitiesRequiredFields = map[string][]string{
	"search":     {},
	"get":        {"id"},
	"create":     {"data_or_data_list"}, // специальный маркер: data ИЛИ data_list
	"update":     {"data_or_data_list"}, // специальный маркер: data ИЛИ data_list (или batch через data_list без id)
	"sync":       {"data"},
	"link":       {"id", "link_to"},
	"unlink":     {"id", "link_to"},
	"get_chats":  {"id"},
	"link_chats": {"chat_links"},
}

// entitiesSchemaFields описывает поля для schema response по каждому action.
type fieldDesc struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required,omitempty"`
}

// entitiesIsSchemaMode определяет, нужно ли вернуть схему.
// Schema mode: если для action есть обязательные поля и хотя бы одно из них отсутствует.
func (t *EntitiesTool) entitiesIsSchemaMode(action string, m map[string]any) bool {
	required, known := entitiesRequiredFields[action]
	if !known {
		// Неизвестный action — отдаём схему чтобы описать доступные actions
		return true
	}
	if len(required) == 0 {
		// search — всегда execute
		return false
	}
	for _, field := range required {
		if field == "data_or_data_list" {
			// Достаточно одного из двух
			hasData := m["data"] != nil
			hasList := m["data_list"] != nil
			if !hasData && !hasList {
				return true
			}
			continue
		}
		if m[field] == nil {
			return true
		}
	}
	return false
}

// entitiesBuildSchemaResponse формирует schema response с полной схемой полей и справочными данными.
func (t *EntitiesTool) entitiesBuildSchemaResponse(entityType, action string) map[string]any {
	svc := t.service

	// Справочные данные
	availableValues := map[string]any{
		"pipelines":    svc.PipelineNames(),
		"statuses":     svc.StatusesByPipeline(),
		"users":        svc.UserNames(),
		"loss_reasons": svc.LossReasonNames(),
		"custom_field_codes": map[string]any{
			"leads":     svc.CustomFieldCodes("leads"),
			"contacts":  svc.CustomFieldCodes("contacts"),
			"companies": svc.CustomFieldCodes("companies"),
		},
	}

	required, optional, description, example := entitiesSchemaForAction(entityType, action)

	return map[string]any{
		"schema":           true,
		"tool":             "entities",
		"action":           action,
		"entity_type":      entityType,
		"description":      description,
		"required_fields":  required,
		"optional_fields":  optional,
		"available_values": availableValues,
		"example":          example,
	}
}

// entitiesSchemaForAction возвращает описание полей, description и пример для action.
func entitiesSchemaForAction(entityType, action string) (required, optional map[string]any, description string, example map[string]any) {
	// Общие опциональные поля для data объекта
	dataFields := map[string]any{
		"name":                  map[string]any{"type": "string", "description": "Название"},
		"responsible_user_name": map[string]any{"type": "string", "description": "Имя ответственного (из available_values.users)"},
		"custom_fields_values":  map[string]any{"type": "object", "description": "Кастомные поля: {\"FIELD_CODE\": \"значение\"} или {\"PHONE\": [{\"value\": \"+7...\", \"enum_code\": \"WORK\"}]}"},
		"tags":                  map[string]any{"type": "array", "description": "Теги: [{\"name\": \"тег\"}]"},
	}
	if entityType == "leads" {
		dataFields["price"] = map[string]any{"type": "integer", "description": "Бюджет сделки"}
		dataFields["pipeline_name"] = map[string]any{"type": "string", "description": "Название воронки (из available_values.pipelines)"}
		dataFields["status_name"] = map[string]any{"type": "string", "description": "Название статуса (из available_values.statuses[pipeline_name])"}
		dataFields["loss_reason_name"] = map[string]any{"type": "string", "description": "Причина отказа (из available_values.loss_reasons)"}
		dataFields["source_name"] = map[string]any{"type": "string", "description": "Название источника"}
		dataFields["embedded_contacts"] = map[string]any{"type": "array", "description": "ID контактов для привязки"}
		dataFields["embedded_companies"] = map[string]any{"type": "array", "description": "ID компаний для привязки"}
	}
	if entityType == "contacts" {
		dataFields["first_name"] = map[string]any{"type": "string", "description": "Имя контакта"}
		dataFields["last_name"] = map[string]any{"type": "string", "description": "Фамилия контакта"}
		dataFields["embedded_companies"] = map[string]any{"type": "array", "description": "ID компаний для привязки"}
	}

	filterFields := map[string]any{
		"query":                  map[string]any{"type": "string", "description": "Поисковый запрос"},
		"limit":                  map[string]any{"type": "integer", "description": "Лимит результатов (макс 250, по умолчанию 50)"},
		"page":                   map[string]any{"type": "integer", "description": "Номер страницы"},
		"ids":                    map[string]any{"type": "array", "description": "Фильтр по ID"},
		"responsible_user_names": map[string]any{"type": "array", "description": "Имена ответственных (из available_values.users)"},
		"created_at_from":        map[string]any{"type": "string", "description": "Дата создания от (ISO-8601, например 2024-01-15T10:00:00Z)"},
		"created_at_to":          map[string]any{"type": "string", "description": "Дата создания до (ISO-8601)"},
		"updated_at_from":        map[string]any{"type": "string", "description": "Дата обновления от (ISO-8601)"},
		"updated_at_to":          map[string]any{"type": "string", "description": "Дата обновления до (ISO-8601)"},
		"custom_fields_values":   map[string]any{"type": "array", "description": "Фильтр по кастомным полям: [{\"field_code\": \"PHONE\", \"values\": [\"+7...\"]}]"},
	}
	if entityType == "leads" {
		filterFields["pipeline_names"] = map[string]any{"type": "array", "description": "Воронки (из available_values.pipelines)"}
		filterFields["statuses"] = map[string]any{"type": "array", "description": "Статусы: [{\"pipeline_name\": \"...\", \"status_name\": \"...\"}]"}
		filterFields["price_from"] = map[string]any{"type": "integer", "description": "Бюджет от"}
		filterFields["price_to"] = map[string]any{"type": "integer", "description": "Бюджет до"}
		filterFields["closed_at_from"] = map[string]any{"type": "string", "description": "Дата закрытия от (ISO-8601)"}
		filterFields["closed_at_to"] = map[string]any{"type": "string", "description": "Дата закрытия до (ISO-8601)"}
	}
	if entityType == "contacts" || entityType == "companies" {
		filterFields["names"] = map[string]any{"type": "array", "description": "Фильтр по названию/имени"}
	}

	linkToField := map[string]any{"type": "object", "description": "Цель связывания: {\"type\": \"leads|contacts|companies\", \"id\": 123}"}

	switch action {
	case "search":
		description = fmt.Sprintf("Поиск %s. Все параметры опциональны — без фильтров возвращает последние записи.", entityType)
		required = map[string]any{}
		optional = map[string]any{
			"filter": map[string]any{"type": "object", "description": "Параметры фильтрации", "fields": filterFields},
			"with":   map[string]any{"type": "array", "description": "Включить связанные данные: leads, contacts, companies, catalog_elements, loss_reason, source"},
		}
		example = map[string]any{
			"entity_type": entityType,
			"action":      "search",
			"filter": map[string]any{
				"query": "пример",
				"limit": 10,
			},
		}

	case "get":
		description = fmt.Sprintf("Получить %s по ID.", entityType)
		required = map[string]any{
			"id": map[string]any{"type": "integer", "description": "ID сущности"},
		}
		optional = map[string]any{
			"with": map[string]any{"type": "array", "description": "Включить связанные данные: leads, contacts, companies, loss_reason, source"},
		}
		example = map[string]any{
			"entity_type": entityType,
			"action":      "get",
			"id":          12345,
		}

	case "create":
		description = fmt.Sprintf("Создать одну или несколько %s.", entityType)
		required = map[string]any{
			"data": map[string]any{"type": "object", "description": "Данные для создания одной записи (или используй data_list для batch)", "fields": dataFields},
		}
		optional = map[string]any{
			"data_list": map[string]any{"type": "array", "description": "Массив данных для batch-создания (вместо data)", "item_fields": dataFields},
		}
		exampleData := map[string]any{"name": "Пример"}
		if entityType == "leads" {
			exampleData["pipeline_name"] = "Основная воронка"
			exampleData["status_name"] = "Новая заявка"
		}
		if entityType == "contacts" {
			exampleData["first_name"] = "Иван"
			exampleData["last_name"] = "Петров"
		}
		example = map[string]any{
			"entity_type": entityType,
			"action":      "create",
			"data":        exampleData,
		}

	case "update":
		description = fmt.Sprintf("Обновить %s по ID или batch через data_list.", entityType)
		required = map[string]any{
			"id":   map[string]any{"type": "integer", "description": "ID сущности (для обновления одной записи)"},
			"data": map[string]any{"type": "object", "description": "Данные для обновления (или используй data_list для batch)", "fields": dataFields},
		}
		optional = map[string]any{
			"data_list": map[string]any{"type": "array", "description": "Batch-обновление: каждый элемент должен содержать id", "item_fields": dataFields},
		}
		example = map[string]any{
			"entity_type": entityType,
			"action":      "update",
			"id":          12345,
			"data":        map[string]any{"name": "Новое название"},
		}

	case "sync":
		description = fmt.Sprintf("Создать или обновить %s (upsert по полям).", entityType)
		required = map[string]any{
			"data": map[string]any{"type": "object", "description": "Данные для sync", "fields": dataFields},
		}
		optional = map[string]any{
			"id": map[string]any{"type": "integer", "description": "ID существующей записи (если известен)"},
		}
		example = map[string]any{
			"entity_type": entityType,
			"action":      "sync",
			"data":        map[string]any{"name": "Пример sync"},
		}

	case "link":
		description = fmt.Sprintf("Связать %s с другой сущностью.", entityType)
		required = map[string]any{
			"id":      map[string]any{"type": "integer", "description": fmt.Sprintf("ID %s", entityType)},
			"link_to": linkToField,
		}
		optional = map[string]any{}
		example = map[string]any{
			"entity_type": entityType,
			"action":      "link",
			"id":          12345,
			"link_to":     map[string]any{"type": "contacts", "id": 67890},
		}

	case "unlink":
		description = fmt.Sprintf("Отвязать %s от другой сущности.", entityType)
		required = map[string]any{
			"id":      map[string]any{"type": "integer", "description": fmt.Sprintf("ID %s", entityType)},
			"link_to": linkToField,
		}
		optional = map[string]any{}
		example = map[string]any{
			"entity_type": entityType,
			"action":      "unlink",
			"id":          12345,
			"link_to":     map[string]any{"type": "contacts", "id": 67890},
		}

	case "get_chats":
		description = "Получить чаты привязанные к контакту (только contacts)."
		required = map[string]any{
			"id": map[string]any{"type": "integer", "description": "ID контакта"},
		}
		optional = map[string]any{}
		example = map[string]any{
			"entity_type": "contacts",
			"action":      "get_chats",
			"id":          12345,
		}

	case "link_chats":
		description = "Привязать чаты к контакту (только contacts)."
		required = map[string]any{
			"chat_links": map[string]any{"type": "array", "description": "Массив ссылок на чаты: [{\"chat_id\": \"...\", \"contact_id\": 123}]"},
		}
		optional = map[string]any{}
		example = map[string]any{
			"entity_type": "contacts",
			"action":      "link_chats",
			"chat_links":  []map[string]any{{"chat_id": "chat-uuid", "contact_id": 12345}},
		}

	default:
		description = fmt.Sprintf("Неизвестный action: %s. Доступные: search, get, create, update, sync, link, unlink, get_chats, link_chats.", action)
		required = map[string]any{}
		optional = map[string]any{}
		example = map[string]any{"entity_type": entityType, "action": "search"}
	}

	return required, optional, description, example
}

// entitiesExecute десериализует input в EntitiesInput и выполняет действие.
func (t *EntitiesTool) entitiesExecute(ctx context.Context, m map[string]any, entityType, action string) (any, error) {
	// JSON roundtrip: map[string]any → EntitiesInput
	raw, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("entities: marshal input: %w", err)
	}
	var input toolmodels.EntitiesInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, fmt.Errorf("entities: unmarshal input: %w", err)
	}
	return t.handleEntities(ctx, input)
}

func (t *EntitiesTool) handleEntities(ctx context.Context, input toolmodels.EntitiesInput) (any, error) {
	switch input.EntityType {
	case "leads":
		return t.handleLeads(ctx, input)
	case "contacts":
		return t.handleContacts(ctx, input)
	case "companies":
		return t.handleCompanies(ctx, input)
	default:
		return nil, fmt.Errorf("unknown entity_type: %s (expected: leads, contacts, companies)", input.EntityType)
	}
}

// ============ LEADS ============

func (t *EntitiesTool) handleLeads(ctx context.Context, input toolmodels.EntitiesInput) (any, error) {
	switch input.Action {
	case "search":
		return t.service.SearchLeads(ctx, input.Filter, input.With)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return t.service.GetLead(ctx, input.ID, input.With)
	case "create":
		if len(input.DataList) > 0 {
			return t.service.CreateLeads(ctx, input.DataList)
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data or data_list is required for action 'create'")
		}
		return t.service.CreateLead(ctx, input.Data)
	case "update":
		if len(input.DataList) > 0 {
			return t.service.UpdateLeads(ctx, input.DataList)
		}
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data or data_list is required for action 'update'")
		}
		return t.service.UpdateLead(ctx, input.ID, input.Data)
	case "sync":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'sync'")
		}
		return t.service.SyncLead(ctx, input.ID, input.Data)
	case "link":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'link'")
		}
		return t.service.LinkLead(ctx, input.ID, input.LinkTo)
	case "unlink":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'unlink'")
		}
		return t.service.UnlinkLead(ctx, input.ID, input.LinkTo)
	default:
		return nil, fmt.Errorf("unknown action for leads: %s", input.Action)
	}
}

// ============ CONTACTS ============

func (t *EntitiesTool) handleContacts(ctx context.Context, input toolmodels.EntitiesInput) (any, error) {
	switch input.Action {
	case "search":
		return t.service.SearchContacts(ctx, input.Filter, input.With)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return t.service.GetContact(ctx, input.ID, input.With)
	case "create":
		if len(input.DataList) > 0 {
			return t.service.CreateContacts(ctx, input.DataList)
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data or data_list is required for action 'create'")
		}
		return t.service.CreateContact(ctx, input.Data)
	case "update":
		if len(input.DataList) > 0 {
			return t.service.UpdateContacts(ctx, input.DataList)
		}
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data or data_list is required for action 'update'")
		}
		return t.service.UpdateContact(ctx, input.ID, input.Data)
	case "sync":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'sync'")
		}
		return t.service.SyncContact(ctx, input.ID, input.Data)
	case "get_chats":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get_chats'")
		}
		return t.service.GetContactChats(ctx, input.ID)
	case "link_chats":
		if len(input.ChatLinks) == 0 {
			return nil, fmt.Errorf("chat_links is required for action 'link_chats'")
		}
		return t.service.LinkContactChats(ctx, input.ChatLinks)
	case "link":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'link'")
		}
		return t.service.LinkContact(ctx, input.ID, input.LinkTo)
	case "unlink":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'unlink'")
		}
		return t.service.UnlinkContact(ctx, input.ID, input.LinkTo)
	default:
		return nil, fmt.Errorf("unknown action for contacts: %s", input.Action)
	}
}

// ============ COMPANIES ============

func (t *EntitiesTool) handleCompanies(ctx context.Context, input toolmodels.EntitiesInput) (any, error) {
	switch input.Action {
	case "search":
		return t.service.SearchCompanies(ctx, input.Filter, input.With)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return t.service.GetCompany(ctx, input.ID, input.With)
	case "create":
		if len(input.DataList) > 0 {
			return t.service.CreateCompanies(ctx, input.DataList)
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data or data_list is required for action 'create'")
		}
		return t.service.CreateCompany(ctx, input.Data)
	case "update":
		if len(input.DataList) > 0 {
			return t.service.UpdateCompanies(ctx, input.DataList)
		}
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data or data_list is required for action 'update'")
		}
		return t.service.UpdateCompany(ctx, input.ID, input.Data)
	case "sync":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'sync'")
		}
		return t.service.SyncCompany(ctx, input.ID, input.Data)
	case "link":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'link'")
		}
		return t.service.LinkCompany(ctx, input.ID, input.LinkTo)
	case "unlink":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'unlink'")
		}
		return t.service.UnlinkCompany(ctx, input.ID, input.LinkTo)
	default:
		return nil, fmt.Errorf("unknown action for companies: %s", input.Action)
	}
}

