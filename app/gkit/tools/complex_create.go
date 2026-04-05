package tools

import (
	"encoding/json"
	"fmt"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// complexCreateSchema — полная схема полей, возвращаемая в schema mode.
// Возвращается LLM при первом вызове без обязательных полей.
func (r *Registry) complexCreateSchema() map[string]any {
	return map[string]any{
		"schema":      true,
		"tool":        "complex_create",
		"description": "Создаёт сделку вместе с контактами и/или компанией за один запрос.",
		"actions": map[string]any{
			"create": map[string]any{
				"description": "Создать одну сделку с контактами и/или компанией.",
				"required_fields": map[string]any{
					"lead": map[string]any{
						"type":        "object",
						"description": "Данные сделки (обязательно)",
						"required_fields": map[string]any{
							"name": map[string]any{"type": "string", "description": "Название сделки"},
						},
						"optional_fields": map[string]any{
							"price":                 map[string]any{"type": "integer", "description": "Бюджет"},
							"pipeline_name":         map[string]any{"type": "string", "description": "Название воронки"},
							"status_name":           map[string]any{"type": "string", "description": "Название статуса в воронке"},
							"responsible_user_name": map[string]any{"type": "string", "description": "Имя ответственного пользователя"},
							"custom_fields_values":  map[string]any{"type": "object", "description": "Кастомные поля сделки. Формат: {field_id: [{value: значение}]}"},
							"tags":                  map[string]any{"type": "array", "items": "string", "description": "Теги сделки (названия)"},
						},
					},
				},
				"optional_fields": map[string]any{
					"contacts": map[string]any{
						"type":        "array",
						"description": "Контакты для создания и привязки",
						"item_fields": map[string]any{
							"name":                  map[string]any{"type": "string", "description": "Имя контакта (полное ФИО)"},
							"first_name":            map[string]any{"type": "string", "description": "Имя (отдельно)"},
							"last_name":             map[string]any{"type": "string", "description": "Фамилия (отдельно)"},
							"phone":                 map[string]any{"type": "string", "description": "Телефон (добавляется как кастомное поле PHONE)"},
							"email":                 map[string]any{"type": "string", "description": "Email (добавляется как кастомное поле EMAIL)"},
							"is_main":               map[string]any{"type": "boolean", "description": "Основной контакт"},
							"responsible_user_name": map[string]any{"type": "string", "description": "Имя ответственного пользователя"},
							"custom_fields_values":  map[string]any{"type": "object", "description": "Прочие кастомные поля контакта"},
						},
					},
					"company": map[string]any{
						"type":        "object",
						"description": "Компания для создания и привязки",
						"fields": map[string]any{
							"name":                  map[string]any{"type": "string", "description": "Название компании"},
							"responsible_user_name": map[string]any{"type": "string", "description": "Имя ответственного пользователя"},
							"custom_fields_values":  map[string]any{"type": "object", "description": "Кастомные поля компании"},
						},
					},
				},
				"example": map[string]any{
					"action": "create",
					"lead": map[string]any{
						"name":                  "Новая сделка",
						"pipeline_name":         "Продажи",
						"status_name":           "Новая заявка",
						"responsible_user_name": "Иван Петров",
					},
					"contacts": []map[string]any{
						{
							"name":  "Мария Сидорова",
							"phone": "+79001234567",
							"email": "maria@example.com",
						},
					},
				},
			},
			"create_batch": map[string]any{
				"description": "Создать несколько сделок за один запрос (до 50 штук). Каждый элемент имеет те же поля что и create.",
				"required_fields": map[string]any{
					"items": map[string]any{
						"type":        "array",
						"description": "Массив сделок для создания (до 50 штук). Каждый элемент: {lead, contacts?, company?}",
					},
				},
				"example": map[string]any{
					"action": "create_batch",
					"items": []map[string]any{
						{
							"lead":     map[string]any{"name": "Сделка 1", "pipeline_name": "Продажи"},
							"contacts": []map[string]any{{"name": "Контакт 1", "phone": "+79001111111"}},
						},
						{
							"lead": map[string]any{"name": "Сделка 2", "pipeline_name": "VIP"},
						},
					},
				},
			},
		},
		"available_values": map[string]any{
			"pipelines": r.complexCreateService.PipelineNames(),
			"statuses":  r.complexCreateService.StatusesByPipeline(),
			"users":     r.complexCreateService.UserNames(),
		},
	}
}

func (r *Registry) RegisterComplexCreateTool() {
	r.addTool(genkit.DefineTool[any, any](
		r.g,
		"complex_create",
		"Создать сделку с контактами и компанией в amoCRM. "+
			"Actions: create (одна сделка), create_batch (до 50 сделок). "+
			"Вызови с {\"action\": \"create\"} без других полей чтобы получить схему параметров и доступные значения.",
		func(ctx *ai.ToolContext, rawInput any) (any, error) {
			m, ok := rawInput.(map[string]any)
			if !ok {
				// Попробуем через JSON roundtrip
				b, err := json.Marshal(rawInput)
				if err != nil {
					return r.complexCreateSchema(), nil
				}
				if err := json.Unmarshal(b, &m); err != nil {
					return r.complexCreateSchema(), nil
				}
			}

			action, _ := m["action"].(string)

			switch action {
			case "create":
				return r.handleComplexCreate(ctx, m)
			case "create_batch":
				return r.handleComplexCreateBatch(ctx, m)
			default:
				// Нет action или неизвестный — возвращаем схему
				return r.complexCreateSchema(), nil
			}
		},
	))
}

// handleComplexCreate — execute mode для create.
// Граница: lead.name должен быть непустым.
func (r *Registry) handleComplexCreate(ctx *ai.ToolContext, m map[string]any) (any, error) {
	leadRaw, hasLead := m["lead"]
	if !hasLead || leadRaw == nil {
		return r.complexCreateSchema(), nil
	}

	leadMap, ok := leadRaw.(map[string]any)
	if !ok {
		return r.complexCreateSchema(), nil
	}

	name, _ := leadMap["name"].(string)
	if name == "" {
		return r.complexCreateSchema(), nil
	}

	// JSON roundtrip map → ComplexCreateInput
	b, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("complex_create: marshal input: %w", err)
	}

	var input gkitmodels.ComplexCreateInput
	if err := json.Unmarshal(b, &input); err != nil {
		return nil, fmt.Errorf("complex_create: unmarshal input: %w", err)
	}

	return r.complexCreateService.CreateComplex(ctx.Context, &input)
}

// handleComplexCreateBatch — execute mode для create_batch.
// Граница: items должен быть непустым массивом.
func (r *Registry) handleComplexCreateBatch(ctx *ai.ToolContext, m map[string]any) (any, error) {
	itemsRaw, hasItems := m["items"]
	if !hasItems || itemsRaw == nil {
		return r.complexCreateSchema(), nil
	}

	itemsSlice, ok := itemsRaw.([]any)
	if !ok || len(itemsSlice) == 0 {
		return r.complexCreateSchema(), nil
	}

	// JSON roundtrip map → ComplexCreateBatchInput
	b, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("complex_create_batch: marshal input: %w", err)
	}

	var input gkitmodels.ComplexCreateBatchInput
	if err := json.Unmarshal(b, &input); err != nil {
		return nil, fmt.Errorf("complex_create_batch: unmarshal input: %w", err)
	}

	return r.complexCreateService.CreateComplexBatch(ctx.Context, input.Items)
}
