package tools

import (
	"encoding/json"
	"fmt"

	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	admin_pipelines "github.com/tihn/amo-ai-tgbot-go/internal/services/crm/admin_pipelines"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// AdminPipelinesTool реализует нативный ADK FunctionTool интерфейс для управления воронками amoCRM.
// Shadow Tool паттерн: минимальная схема видна LLM, полная схема возвращается при первом вызове.
type AdminPipelinesTool struct {
	service admin_pipelines.Service
}

// NewAdminPipelinesTool создаёт новый AdminPipelinesTool с указанным сервисом.
func NewAdminPipelinesTool(service admin_pipelines.Service) *AdminPipelinesTool {
	return &AdminPipelinesTool{service: service}
}

// Name implements tool.Tool.
func (t *AdminPipelinesTool) Name() string {
	return "admin_pipelines"
}

// Description implements tool.Tool.
func (t *AdminPipelinesTool) Description() string {
	return "Управление воронками и статусами amoCRM. " +
		"Actions: search, get, create, update, delete, get_statuses, get_status, create_status, update_status, delete_status. " +
		"Вызови с action чтобы получить схему параметров."
}

// IsLongRunning implements tool.Tool.
func (t *AdminPipelinesTool) IsLongRunning() bool {
	return false
}

// ProcessRequest реализует toolinternal.RequestProcessor — регистрирует Declaration в LLM request.
func (t *AdminPipelinesTool) ProcessRequest(_ tool.Context, req *model.LLMRequest) error {
	return packToolDeclaration(req, t)
}

// Declaration implements toolinternal.FunctionTool (duck typing).
func (t *AdminPipelinesTool) Declaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"action": {
					Type:        genai.TypeString,
					Description: "Действие над воронкой или статусом: search, get, create, update, delete, get_statuses, get_status, create_status, update_status, delete_status",
					Enum:        []string{"search", "get", "create", "update", "delete", "get_statuses", "get_status", "create_status", "update_status", "delete_status"},
				},
			},
			Required: []string{"action"},
		},
	}
}

// Run implements toolinternal.FunctionTool (duck typing).
func (t *AdminPipelinesTool) Run(ctx tool.Context, args any) (map[string]any, error) {
	raw, ok := args.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("admin_pipelines: invalid input type %T", args)
	}

	// Извлекаем action
	actionVal, ok := raw["action"]
	if !ok || actionVal == nil {
		return nil, fmt.Errorf("поле action обязательно. Доступные: search, get, create, update, delete, get_statuses, get_status, create_status, update_status, delete_status")
	}
	action, ok := actionVal.(string)
	if !ok || action == "" {
		return nil, fmt.Errorf("action должен быть строкой")
	}

	// Schema-режим: обязательных полей нет → возвращаем схему
	if adminPipelinesIsSchemaMode(action, raw) {
		schema, err := adminPipelinesSchema(action)
		if err != nil {
			return nil, err
		}
		return toResultMap(schema)
	}

	// Execute-режим: json roundtrip map → AdminPipelinesInput
	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации входных данных: %w", err)
	}
	var inp gkitmodels.AdminPipelinesInput
	if err := json.Unmarshal(data, &inp); err != nil {
		return nil, fmt.Errorf("ошибка разбора входных данных: %w", err)
	}

	var result any

	switch inp.Action {

	// --- Pipelines ---

	case "list", "search":
		res, err := t.service.ListPipelines(ctx, inp.WithStatuses)
		b, _ := json.Marshal(res)
		fmt.Printf("[admin_pipelines] search result: %s\n", string(b))
		result, err = res, err
		if err != nil {
			return nil, err
		}

	case "get":
		res, err := t.service.GetPipeline(ctx, inp.PipelineID, inp.PipelineName, inp.WithStatuses)
		if err != nil {
			return nil, err
		}
		result = res

	case "create":
		if len(inp.Items) > 0 {
			// Батч-режим: items — массив PipelineData
			itemsData, err := json.Marshal(inp.Items)
			if err != nil {
				return nil, fmt.Errorf("не удалось сериализовать items: %w", err)
			}
			var pipelines []gkitmodels.PipelineData
			if err := json.Unmarshal(itemsData, &pipelines); err != nil {
				return nil, fmt.Errorf("не удалось разобрать items как []PipelineData: %w", err)
			}
			res, err := t.service.CreatePipelines(ctx, pipelines)
			if err != nil {
				return nil, err
			}
			result = res
		} else if inp.Pipeline != nil {
			res, err := t.service.CreatePipelines(ctx, []gkitmodels.PipelineData{*inp.Pipeline})
			if err != nil {
				return nil, err
			}
			result = res
		} else {
			return nil, fmt.Errorf("для create укажите pipeline или items")
		}

	case "update":
		if inp.Pipeline == nil {
			return nil, fmt.Errorf("для update укажите pipeline с данными для обновления")
		}
		res, err := t.service.UpdatePipeline(ctx, inp.PipelineID, inp.PipelineName, *inp.Pipeline)
		if err != nil {
			return nil, err
		}
		result = res

	case "delete":
		if err := t.service.DeletePipeline(ctx, inp.PipelineID, inp.PipelineName); err != nil {
			return nil, err
		}
		result = map[string]any{"deleted": true}

	// --- Statuses ---

	case "list_statuses", "get_statuses":
		res, err := t.service.ListStatuses(ctx, inp.PipelineID, inp.PipelineName)
		if err != nil {
			return nil, err
		}
		result = res

	case "get_status":
		res, err := t.service.GetStatus(ctx, inp.PipelineID, inp.PipelineName, inp.StatusID, inp.StatusName)
		if err != nil {
			return nil, err
		}
		result = res

	case "create_status":
		if len(inp.Items) > 0 {
			// Батч-режим: items — массив StatusData
			itemsData, err := json.Marshal(inp.Items)
			if err != nil {
				return nil, fmt.Errorf("не удалось сериализовать items: %w", err)
			}
			var statuses []gkitmodels.StatusData
			if err := json.Unmarshal(itemsData, &statuses); err != nil {
				return nil, fmt.Errorf("не удалось разобрать items как []StatusData: %w", err)
			}
			res, err := t.service.CreateStatuses(ctx, inp.PipelineID, inp.PipelineName, statuses)
			if err != nil {
				return nil, err
			}
			result = res
		} else if inp.Status != nil {
			res, err := t.service.CreateStatus(ctx, inp.PipelineID, inp.PipelineName, *inp.Status)
			if err != nil {
				return nil, err
			}
			result = res
		} else {
			return nil, fmt.Errorf("для create_status укажите status или items")
		}

	case "update_status":
		if inp.Status == nil {
			return nil, fmt.Errorf("для update_status укажите status с данными для обновления")
		}
		res, err := t.service.UpdateStatus(ctx, inp.PipelineID, inp.PipelineName, inp.StatusID, inp.StatusName, *inp.Status)
		if err != nil {
			return nil, err
		}
		result = res

	case "delete_status":
		if err := t.service.DeleteStatus(ctx, inp.PipelineID, inp.PipelineName, inp.StatusID, inp.StatusName); err != nil {
			return nil, err
		}
		result = map[string]any{"deleted": true}

	default:
		return nil, fmt.Errorf("неизвестное действие: %s", inp.Action)
	}

	return toResultMap(result)
}

// adminPipelinesSchema возвращает полную схему параметров для заданного action.
// Используется в Schema-режиме Shadow Tool.
func adminPipelinesSchema(action string) (any, error) {
	type fieldDesc struct {
		Type        string `json:"type"`
		Description string `json:"description"`
		Required    bool   `json:"required"`
	}

	type schemaResponse struct {
		Schema         bool                 `json:"schema"`
		Tool           string               `json:"tool"`
		Action         string               `json:"action"`
		Description    string               `json:"description"`
		RequiredFields map[string]fieldDesc `json:"required_fields"`
		OptionalFields map[string]fieldDesc `json:"optional_fields,omitempty"`
		Notes          []string             `json:"notes,omitempty"`
		Example        map[string]any       `json:"example"`
	}

	base := schemaResponse{
		Schema: true,
		Tool:   "admin_pipelines",
		Action: action,
	}

	switch action {

	case "search", "list":
		base.Description = "Получить список всех воронок. Не требует обязательных параметров."
		base.RequiredFields = map[string]fieldDesc{}
		base.OptionalFields = map[string]fieldDesc{
			"with_statuses": {Type: "boolean", Description: "Включить статусы каждой воронки в ответ (один запрос)"},
		}
		base.Example = map[string]any{
			"action":        "search",
			"with_statuses": true,
		}

	case "get":
		base.Description = "Получить воронку по ID или имени."
		base.RequiredFields = map[string]fieldDesc{
			"pipeline_id OR pipeline_name": {Type: "int | string", Description: "Одно из двух: числовой ID или точное имя воронки", Required: true},
		}
		base.OptionalFields = map[string]fieldDesc{
			"pipeline_id":   {Type: "integer", Description: "Числовой ID воронки"},
			"pipeline_name": {Type: "string", Description: "Имя воронки (используется если pipeline_id не задан)"},
			"with_statuses": {Type: "boolean", Description: "Включить статусы воронки в ответ"},
		}
		base.Notes = []string{"Укажите pipeline_id или pipeline_name — хотя бы одно"}
		base.Example = map[string]any{
			"action":        "get",
			"pipeline_name": "Продажи",
			"with_statuses": true,
		}

	case "create":
		base.Description = "Создать одну или несколько воронок."
		base.RequiredFields = map[string]fieldDesc{
			"pipeline OR items": {Type: "object | array", Description: "Данные воронки (одна) или массив воронок (батч)", Required: true},
		}
		base.OptionalFields = map[string]fieldDesc{
			"pipeline": {Type: "object", Description: "Данные одной воронки: {name, sort, is_main, is_unsorted_on}"},
			"items":    {Type: "array", Description: "Массив PipelineData для батч-создания: [{name, sort, is_main, is_unsorted_on}, ...]"},
		}
		base.Notes = []string{
			"pipeline.name — обязательное поле внутри pipeline",
			"items[].name — обязательное поле каждого элемента",
		}
		base.Example = map[string]any{
			"action": "create",
			"pipeline": map[string]any{
				"name":           "Новая воронка",
				"sort":           100,
				"is_unsorted_on": true,
			},
		}

	case "update":
		base.Description = "Обновить воронку по ID или имени."
		base.RequiredFields = map[string]fieldDesc{
			"pipeline_id OR pipeline_name": {Type: "int | string", Description: "Идентификатор воронки для обновления", Required: true},
			"pipeline":                     {Type: "object", Description: "Данные для обновления: {name, sort, is_main, is_unsorted_on}", Required: true},
		}
		base.OptionalFields = map[string]fieldDesc{
			"pipeline_id":   {Type: "integer", Description: "Числовой ID воронки"},
			"pipeline_name": {Type: "string", Description: "Имя воронки (используется если pipeline_id не задан)"},
		}
		base.Notes = []string{"Укажите pipeline_id или pipeline_name — хотя бы одно"}
		base.Example = map[string]any{
			"action":        "update",
			"pipeline_name": "Старое название",
			"pipeline": map[string]any{
				"name": "Новое название",
				"sort": 200,
			},
		}

	case "delete":
		base.Description = "Удалить воронку по ID или имени."
		base.RequiredFields = map[string]fieldDesc{
			"pipeline_id OR pipeline_name": {Type: "int | string", Description: "Идентификатор воронки для удаления", Required: true},
		}
		base.OptionalFields = map[string]fieldDesc{
			"pipeline_id":   {Type: "integer", Description: "Числовой ID воронки"},
			"pipeline_name": {Type: "string", Description: "Имя воронки"},
		}
		base.Notes = []string{"Укажите pipeline_id или pipeline_name — хотя бы одно"}
		base.Example = map[string]any{
			"action":        "delete",
			"pipeline_name": "Устаревшая воронка",
		}

	case "get_statuses", "list_statuses":
		base.Description = "Получить список статусов воронки."
		base.RequiredFields = map[string]fieldDesc{
			"pipeline_id OR pipeline_name": {Type: "int | string", Description: "Идентификатор воронки", Required: true},
		}
		base.OptionalFields = map[string]fieldDesc{
			"pipeline_id":   {Type: "integer", Description: "Числовой ID воронки"},
			"pipeline_name": {Type: "string", Description: "Имя воронки"},
		}
		base.Notes = []string{"Укажите pipeline_id или pipeline_name — хотя бы одно"}
		base.Example = map[string]any{
			"action":        "get_statuses",
			"pipeline_name": "Продажи",
		}

	case "get_status":
		base.Description = "Получить статус воронки по ID или имени."
		base.RequiredFields = map[string]fieldDesc{
			"pipeline_id OR pipeline_name": {Type: "int | string", Description: "Идентификатор воронки", Required: true},
			"status_id OR status_name":     {Type: "int | string", Description: "Идентификатор статуса", Required: true},
		}
		base.OptionalFields = map[string]fieldDesc{
			"pipeline_id":   {Type: "integer", Description: "Числовой ID воронки"},
			"pipeline_name": {Type: "string", Description: "Имя воронки"},
			"status_id":     {Type: "integer", Description: "Числовой ID статуса"},
			"status_name":   {Type: "string", Description: "Имя статуса"},
		}
		base.Notes = []string{
			"Укажите pipeline_id или pipeline_name — хотя бы одно",
			"Укажите status_id или status_name — хотя бы одно",
		}
		base.Example = map[string]any{
			"action":        "get_status",
			"pipeline_name": "Продажи",
			"status_name":   "В работе",
		}

	case "create_status":
		base.Description = "Создать один или несколько статусов в воронке."
		base.RequiredFields = map[string]fieldDesc{
			"pipeline_id OR pipeline_name": {Type: "int | string", Description: "Идентификатор воронки", Required: true},
			"status OR items":              {Type: "object | array", Description: "Данные статуса или массив статусов", Required: true},
		}
		base.OptionalFields = map[string]fieldDesc{
			"pipeline_id":   {Type: "integer", Description: "Числовой ID воронки"},
			"pipeline_name": {Type: "string", Description: "Имя воронки"},
			"status": {Type: "object", Description: "Данные одного статуса: {name (обязательное), sort, color, type}. " +
				"type: regular | won | lost. " +
				"color (hex): #fffeb2 #fffd7f #fff000 #ffeab2 #ffdc7f #ffce5a #ffdbdb #ffc8c8 #ff8f92 " +
				"#d6eaff #c1e0ff #98cbff #ebffb1 #deff81 #87f2c0 #f9deff #f3beff #ccc8f9 #eb93ff #f2f3f4 #e6e8ea"},
			"items": {Type: "array", Description: "Массив StatusData для батч-создания: [{name, sort, color, type}, ...]"},
		}
		base.Example = map[string]any{
			"action":        "create_status",
			"pipeline_name": "Продажи",
			"status": map[string]any{
				"name":  "Переговоры",
				"color": "#d6eaff",
				"type":  "regular",
			},
		}

	case "update_status":
		base.Description = "Обновить статус воронки."
		base.RequiredFields = map[string]fieldDesc{
			"pipeline_id OR pipeline_name": {Type: "int | string", Description: "Идентификатор воронки", Required: true},
			"status_id OR status_name":     {Type: "int | string", Description: "Идентификатор статуса", Required: true},
			"status":                       {Type: "object", Description: "Данные для обновления статуса: {name, sort, color, type}", Required: true},
		}
		base.OptionalFields = map[string]fieldDesc{
			"pipeline_id":   {Type: "integer", Description: "Числовой ID воронки"},
			"pipeline_name": {Type: "string", Description: "Имя воронки"},
			"status_id":     {Type: "integer", Description: "Числовой ID статуса"},
			"status_name":   {Type: "string", Description: "Имя статуса"},
		}
		base.Notes = []string{
			"Укажите pipeline_id или pipeline_name — хотя бы одно",
			"Укажите status_id или status_name — хотя бы одно",
		}
		base.Example = map[string]any{
			"action":        "update_status",
			"pipeline_name": "Продажи",
			"status_name":   "Переговоры",
			"status": map[string]any{
				"name":  "Активные переговоры",
				"color": "#87f2c0",
			},
		}

	case "delete_status":
		base.Description = "Удалить статус воронки."
		base.RequiredFields = map[string]fieldDesc{
			"pipeline_id OR pipeline_name": {Type: "int | string", Description: "Идентификатор воронки", Required: true},
			"status_id OR status_name":     {Type: "int | string", Description: "Идентификатор статуса", Required: true},
		}
		base.OptionalFields = map[string]fieldDesc{
			"pipeline_id":   {Type: "integer", Description: "Числовой ID воронки"},
			"pipeline_name": {Type: "string", Description: "Имя воронки"},
			"status_id":     {Type: "integer", Description: "Числовой ID статуса"},
			"status_name":   {Type: "string", Description: "Имя статуса"},
		}
		base.Notes = []string{
			"Укажите pipeline_id или pipeline_name — хотя бы одно",
			"Укажите status_id или status_name — хотя бы одно",
		}
		base.Example = map[string]any{
			"action":        "delete_status",
			"pipeline_name": "Продажи",
			"status_name":   "Устаревший статус",
		}

	default:
		return nil, fmt.Errorf("неизвестное действие: %s. Доступные: search, get, create, update, delete, get_statuses, get_status, create_status, update_status, delete_status", action)
	}

	return base, nil
}

// adminPipelinesIsSchemaMode определяет, находится ли вызов в Schema-режиме.
// Возвращает true если пришёл только action без обязательных полей для данного action.
func adminPipelinesIsSchemaMode(action string, raw map[string]any) bool {
	hasPipeline := raw["pipeline_id"] != nil && raw["pipeline_id"] != float64(0) ||
		raw["pipeline_name"] != nil && raw["pipeline_name"] != ""
	hasStatus := raw["status_id"] != nil && raw["status_id"] != float64(0) ||
		raw["status_name"] != nil && raw["status_name"] != ""
	hasPipelineData := raw["pipeline"] != nil
	hasStatusData := raw["status"] != nil
	hasItems := raw["items"] != nil

	switch action {
	case "search", "list":
		// Не требует обязательных полей — всегда execute
		return false

	case "get":
		return !hasPipeline

	case "create":
		return !hasPipelineData && !hasItems

	case "update":
		return !hasPipeline || !hasPipelineData

	case "delete":
		return !hasPipeline

	case "get_statuses", "list_statuses":
		return !hasPipeline

	case "get_status":
		return !hasPipeline || !hasStatus

	case "create_status":
		return !hasPipeline || (!hasStatusData && !hasItems)

	case "update_status":
		return !hasPipeline || !hasStatus || !hasStatusData

	case "delete_status":
		return !hasPipeline || !hasStatus

	default:
		return false
	}
}
