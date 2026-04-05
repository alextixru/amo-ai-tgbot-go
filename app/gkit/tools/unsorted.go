package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// unsortedRequiredFields определяет обязательные поля для каждого action.
// Если все они отсутствуют или пусты — handler возвращает схему.
var unsortedRequiredFields = map[string][]string{
	"list":    {},         // нет обязательных — всегда execute
	"summary": {},         // нет обязательных — всегда execute
	"get":     {"uid"},
	"accept":  {"uid"},
	"decline": {"uid"},
	"link":    {"uid", "lead_id"},
	"create":  {"category", "items"},
}

// unsortedSchemas содержит полные описания полей для каждого action.
var unsortedSchemas = map[string]map[string]any{
	"list": {
		"description": "Список неразобранных заявок с пагинацией. Возвращает items + page_meta.",
		"required_fields": map[string]any{},
		"optional_fields": map[string]any{
			"filter": map[string]any{
				"type":        "object",
				"description": "Фильтры поиска",
				"fields": map[string]any{
					"page":            "number — номер страницы (по умолчанию 1)",
					"limit":           "number — лимит результатов",
					"category":        "[]string — категории: sip, mail, forms, chats",
					"pipeline_name":   "string — имя воронки для фильтрации",
					"created_at_from": "string — начало диапазона (RFC3339, напр. 2024-01-01T00:00:00Z)",
					"created_at_to":   "string — конец диапазона (RFC3339, напр. 2024-01-31T23:59:59Z)",
					"order":           "string — сортировка: 'created_at asc' или 'created_at desc'",
				},
			},
		},
		"example": map[string]any{
			"action": "list",
			"filter": map[string]any{
				"limit":         10,
				"pipeline_name": "Основная воронка",
				"order":         "created_at desc",
			},
		},
	},
	"get": {
		"description": "Получить одну запись неразобранного по UID. Возвращает полный объект с вложенными сделками, контактами, компаниями.",
		"required_fields": map[string]any{
			"uid": "string — уникальный идентификатор записи неразобранного",
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action": "get",
			"uid":    "abc123def456",
		},
	},
	"accept": {
		"description": "Принять заявку из Неразобранного — создаёт сделку в воронке. Возвращает {uid, success}.",
		"required_fields": map[string]any{
			"uid": "string — UID записи неразобранного",
		},
		"optional_fields": map[string]any{
			"accept_params": map[string]any{
				"type":        "object",
				"description": "Параметры принятия (все опциональны)",
				"fields": map[string]any{
					"user_name":     "string — имя ответственного пользователя",
					"pipeline_name": "string — имя воронки для создаваемой сделки",
					"status_name":   "string — имя статуса для создаваемой сделки (требует pipeline_name)",
				},
			},
		},
		"example": map[string]any{
			"action": "accept",
			"uid":    "abc123def456",
			"accept_params": map[string]any{
				"user_name":     "Иван Петров",
				"pipeline_name": "Основная воронка",
				"status_name":   "В работе",
			},
		},
	},
	"decline": {
		"description": "Отклонить заявку из Неразобранного. Возвращает {uid, success}.",
		"required_fields": map[string]any{
			"uid": "string — UID записи неразобранного",
		},
		"optional_fields": map[string]any{
			"decline_params": map[string]any{
				"type":        "object",
				"description": "Параметры отклонения (все опциональны)",
				"fields": map[string]any{
					"user_name": "string — имя пользователя, выполняющего отклонение",
				},
			},
		},
		"example": map[string]any{
			"action": "decline",
			"uid":    "abc123def456",
			"decline_params": map[string]any{
				"user_name": "Иван Петров",
			},
		},
	},
	"link": {
		"description": "Привязать запись Неразобранного к существующей сделке. Возвращает {uid, success}.",
		"required_fields": map[string]any{
			"uid":     "string — UID записи неразобранного",
			"lead_id": "number — ID существующей сделки для привязки",
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action":  "link",
			"uid":     "abc123def456",
			"lead_id": 12345,
		},
	},
	"summary": {
		"description": "Статистика по Неразобранному (количество по категориям). Возвращает агрегированные данные.",
		"required_fields": map[string]any{},
		"optional_fields": map[string]any{
			"filter": map[string]any{
				"type":        "object",
				"description": "Фильтры для статистики (все опциональны)",
				"fields": map[string]any{
					"pipeline_name":   "string — имя воронки",
					"created_at_from": "string — начало периода (RFC3339)",
					"created_at_to":   "string — конец периода (RFC3339)",
				},
			},
		},
		"example": map[string]any{
			"action": "summary",
			"filter": map[string]any{
				"pipeline_name":   "Основная воронка",
				"created_at_from": "2024-01-01T00:00:00Z",
				"created_at_to":   "2024-01-31T23:59:59Z",
			},
		},
	},
	"create": {
		"description": "Создать одну или несколько заявок в Неразобранном (батч). Возвращает массив созданных объектов.",
		"required_fields": map[string]any{
			"category": "string — категория источника: sip, forms, chats",
			"items": map[string]any{
				"type":        "array",
				"description": "Массив создаваемых заявок (минимум одна)",
				"item_fields": map[string]any{
					"source_uid":    "string — уникальный идентификатор источника (опционально)",
					"source_name":   "string — название источника (опционально)",
					"pipeline_name": "string — имя воронки (опционально)",
					"created_at":    "string — дата создания RFC3339 (опционально)",
					"data":          "object — дополнительные данные (опционально)",
				},
			},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action":   "create",
			"category": "forms",
			"items": []map[string]any{
				{
					"source_name":   "Сайт",
					"pipeline_name": "Основная воронка",
				},
			},
		},
	},
}

func (r *Registry) RegisterUnsortedTool() {
	r.addTool(genkit.DefineTool[map[string]any, any](
		r.g,
		"unsorted",
		"Работа с входящими заявками amoCRM (Неразобранное). "+
			"Actions: list (список), get (по UID), accept (принять → создаёт сделку), "+
			"decline (отклонить), link (привязать к сделке), summary (статистика), create (создать заявку). "+
			"Вызови с action чтобы получить схему параметров и доступные значения.",
		func(ctx *ai.ToolContext, input map[string]any) (any, error) {
			return r.handleUnsortedShadow(ctx.Context, input)
		},
	))
}

// handleUnsortedShadow реализует Shadow Tool паттерн:
// Schema mode: action присутствует, обязательные поля action'а отсутствуют → возвращает схему + available_values.
// Execute mode: все обязательные поля присутствуют → выполняет действие.
func (r *Registry) handleUnsortedShadow(ctx context.Context, input map[string]any) (any, error) {
	action, _ := input["action"].(string)
	if action == "" {
		return map[string]any{
			"schema":           true,
			"tool":             "unsorted",
			"error":            "поле action обязательно",
			"available_actions": []string{"list", "get", "accept", "decline", "link", "summary", "create"},
		}, nil
	}

	required, known := unsortedRequiredFields[action]
	if !known {
		return nil, fmt.Errorf("unknown action: %s. Доступные: list, get, accept, decline, link, summary, create", action)
	}

	// Определяем режим: schema если хотя бы одно обязательное поле отсутствует.
	if isSchemaMode(input, required) {
		return r.unsortedSchemaResponse(action), nil
	}

	// Execute mode: json roundtrip map → UnsortedInput.
	return r.executeUnsorted(ctx, input, action)
}

// isSchemaMode возвращает true если хотя бы одно обязательное поле отсутствует или пустое.
func isSchemaMode(input map[string]any, required []string) bool {
	for _, field := range required {
		v, ok := input[field]
		if !ok || v == nil || v == "" || v == 0.0 {
			return true
		}
		// lead_id приходит как float64 из JSON
		if f, ok := v.(float64); ok && f == 0 {
			return true
		}
	}
	return false
}

// unsortedSchemaResponse строит ответ со схемой и available_values из сервиса.
func (r *Registry) unsortedSchemaResponse(action string) map[string]any {
	schema, ok := unsortedSchemas[action]
	if !ok {
		schema = map[string]any{"description": "Схема для action " + action + " не найдена"}
	}

	resp := map[string]any{
		"schema": true,
		"tool":   "unsorted",
		"action": action,
	}
	for k, v := range schema {
		resp[k] = v
	}

	resp["available_values"] = map[string]any{
		"pipelines": r.unsortedService.PipelineNames(),
		"statuses":  r.unsortedService.StatusNames(),
		"users":     r.unsortedService.UserNames(),
	}

	return resp
}

// executeUnsorted выполняет действие: json roundtrip map → UnsortedInput → handleUnsorted.
func (r *Registry) executeUnsorted(ctx context.Context, input map[string]any, action string) (any, error) {
	// Для "link" поле lead_id находится на верхнем уровне, но UnsortedInput ожидает link_data.lead_id.
	// Нормализуем: если есть lead_id на верхнем уровне — оборачиваем в link_data.
	if action == "link" {
		if _, hasLinkData := input["link_data"]; !hasLinkData {
			if leadID, ok := input["lead_id"]; ok {
				input["link_data"] = map[string]any{"lead_id": leadID}
			}
		}
	}

	// Для "create" поля category и items находятся на верхнем уровне,
	// но UnsortedInput ожидает create_data.{category, items}.
	if action == "create" {
		if _, hasCreateData := input["create_data"]; !hasCreateData {
			createData := map[string]any{}
			if cat, ok := input["category"]; ok {
				createData["category"] = cat
			}
			if items, ok := input["items"]; ok {
				createData["items"] = items
			}
			input["create_data"] = createData
		}
	}

	data, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("unsorted: marshal input: %w", err)
	}
	var typed gkitmodels.UnsortedInput
	if err := json.Unmarshal(data, &typed); err != nil {
		return nil, fmt.Errorf("unsorted: unmarshal input: %w", err)
	}
	return r.handleUnsorted(ctx, typed)
}

func (r *Registry) handleUnsorted(ctx context.Context, input gkitmodels.UnsortedInput) (any, error) {
	switch input.Action {
	case "list":
		return r.unsortedService.ListUnsorted(ctx, input.Filter)
	case "get":
		if input.UID == "" {
			return nil, fmt.Errorf("uid is required for action 'get'")
		}
		return r.unsortedService.GetUnsorted(ctx, input.UID)
	case "create":
		if input.CreateData == nil || input.CreateData.Category == "" || len(input.CreateData.Items) == 0 {
			return nil, fmt.Errorf("create_data with category and items is required for action 'create'")
		}
		return r.unsortedService.CreateUnsorted(ctx, input.CreateData.Category, input.CreateData.Items)
	case "accept":
		if input.UID == "" {
			return nil, fmt.Errorf("uid is required for action 'accept'")
		}
		return r.unsortedService.AcceptUnsorted(ctx, input.UID, input.AcceptParams)
	case "decline":
		if input.UID == "" {
			return nil, fmt.Errorf("uid is required for action 'decline'")
		}
		return r.unsortedService.DeclineUnsorted(ctx, input.UID, input.DeclineParams)
	case "link":
		if input.UID == "" || input.LinkData == nil || input.LinkData.LeadID == 0 {
			return nil, fmt.Errorf("uid and link_data.lead_id are required for action 'link'")
		}
		return r.unsortedService.LinkUnsorted(ctx, input.UID, input.LinkData.LeadID)
	case "summary":
		return r.unsortedService.SummaryUnsorted(ctx, input.Filter)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}
