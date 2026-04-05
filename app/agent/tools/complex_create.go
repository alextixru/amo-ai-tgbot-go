package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/complex_create"
)

// ComplexCreateTool реализует нативный ADK Tool интерфейс для комплексного создания сделок.
type ComplexCreateTool struct {
	service complex_create.Service
}

// NewComplexCreateTool создаёт новый экземпляр ComplexCreateTool.
func NewComplexCreateTool(service complex_create.Service) *ComplexCreateTool {
	return &ComplexCreateTool{service: service}
}

// Name возвращает имя инструмента.
func (t *ComplexCreateTool) Name() string {
	return "complex_create"
}

// Description возвращает описание инструмента.
func (t *ComplexCreateTool) Description() string {
	return "Создать сделку с контактами и компанией в amoCRM. " +
		"Actions: create (одна сделка), create_batch (до 50 сделок). " +
		"Вызови с {\"action\": \"create\"} без других полей чтобы получить схему параметров и доступные значения."
}

// IsLongRunning возвращает false — инструмент не является длительной операцией.
func (t *ComplexCreateTool) IsLongRunning() bool {
	return false
}

// Declaration возвращает FunctionDeclaration для регистрации в genai.
func (t *ComplexCreateTool) Declaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"action": {
					Type:        genai.TypeString,
					Description: "Действие: create (одна сделка) или create_batch (до 50 сделок). Если не указан — возвращается схема параметров.",
				},
				"lead": {
					Type:        genai.TypeObject,
					Description: "Данные сделки (обязательно для action=create).",
					Properties: map[string]*genai.Schema{
						"name":                  {Type: genai.TypeString, Description: "Название сделки"},
						"price":                 {Type: genai.TypeInteger, Description: "Бюджет"},
						"pipeline_name":         {Type: genai.TypeString, Description: "Название воронки"},
						"status_name":           {Type: genai.TypeString, Description: "Название статуса в воронке"},
						"responsible_user_name": {Type: genai.TypeString, Description: "Имя ответственного пользователя"},
						"custom_fields_values":  {Type: genai.TypeObject, Description: "Кастомные поля сделки. Формат: {field_id: [{value: значение}]}"},
						"tags": {
							Type:        genai.TypeArray,
							Description: "Теги сделки (названия)",
							Items:       &genai.Schema{Type: genai.TypeString},
						},
					},
				},
				"contacts": {
					Type:        genai.TypeArray,
					Description: "Контакты для создания и привязки (для action=create).",
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"name":                  {Type: genai.TypeString, Description: "Имя контакта (полное ФИО)"},
							"first_name":            {Type: genai.TypeString, Description: "Имя (отдельно)"},
							"last_name":             {Type: genai.TypeString, Description: "Фамилия (отдельно)"},
							"phone":                 {Type: genai.TypeString, Description: "Телефон (добавляется как кастомное поле PHONE)"},
							"email":                 {Type: genai.TypeString, Description: "Email (добавляется как кастомное поле EMAIL)"},
							"is_main":               {Type: genai.TypeBoolean, Description: "Основной контакт"},
							"responsible_user_name": {Type: genai.TypeString, Description: "Имя ответственного пользователя"},
							"custom_fields_values":  {Type: genai.TypeObject, Description: "Прочие кастомные поля контакта"},
						},
					},
				},
				"company": {
					Type:        genai.TypeObject,
					Description: "Компания для создания и привязки (для action=create).",
					Properties: map[string]*genai.Schema{
						"name":                  {Type: genai.TypeString, Description: "Название компании"},
						"responsible_user_name": {Type: genai.TypeString, Description: "Имя ответственного пользователя"},
						"custom_fields_values":  {Type: genai.TypeObject, Description: "Кастомные поля компании"},
					},
				},
				"items": {
					Type:        genai.TypeArray,
					Description: "Массив сделок для создания (до 50 штук). Каждый элемент: {lead, contacts?, company?}. Используется для action=create_batch.",
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"lead":     {Type: genai.TypeObject, Description: "Данные сделки"},
							"contacts": {Type: genai.TypeArray, Description: "Контакты", Items: &genai.Schema{Type: genai.TypeObject}},
							"company":  {Type: genai.TypeObject, Description: "Компания"},
						},
					},
				},
			},
			Required: []string{"action"},
		},
	}
}

// Run выполняет инструмент. args — map[string]any с параметрами вызова.
func (t *ComplexCreateTool) Run(ctx tool.Context, args any) (map[string]any, error) {
	m, ok := args.(map[string]any)
	if !ok {
		// Попробуем через JSON roundtrip
		b, err := json.Marshal(args)
		if err != nil {
			return t.complexCreateSchema(), nil
		}
		if err := json.Unmarshal(b, &m); err != nil {
			return t.complexCreateSchema(), nil
		}
	}

	action, _ := m["action"].(string)

	switch action {
	case "create":
		return t.handleComplexCreate(ctx, m)
	case "create_batch":
		return t.handleComplexCreateBatch(ctx, m)
	default:
		// Нет action или неизвестный — возвращаем схему
		return t.complexCreateSchema(), nil
	}
}

// complexCreateSchema — полная схема полей, возвращаемая в schema mode.
// Возвращается LLM при первом вызове без обязательных полей.
func (t *ComplexCreateTool) complexCreateSchema() map[string]any {
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
			"pipelines": t.service.PipelineNames(),
			"statuses":  t.service.StatusesByPipeline(),
			"users":     t.service.UserNames(),
		},
	}
}

// handleComplexCreate — execute mode для create.
// Граница: lead.name должен быть непустым.
func (t *ComplexCreateTool) handleComplexCreate(ctx context.Context, m map[string]any) (map[string]any, error) {
	leadRaw, hasLead := m["lead"]
	if !hasLead || leadRaw == nil {
		return t.complexCreateSchema(), nil
	}

	leadMap, ok := leadRaw.(map[string]any)
	if !ok {
		return t.complexCreateSchema(), nil
	}

	name, _ := leadMap["name"].(string)
	if name == "" {
		return t.complexCreateSchema(), nil
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

	result, err := t.service.CreateComplex(ctx, &input)
	if err != nil {
		return nil, err
	}
	return toResultMap(result)
}

// handleComplexCreateBatch — execute mode для create_batch.
// Граница: items должен быть непустым массивом.
func (t *ComplexCreateTool) handleComplexCreateBatch(ctx context.Context, m map[string]any) (map[string]any, error) {
	itemsRaw, hasItems := m["items"]
	if !hasItems || itemsRaw == nil {
		return t.complexCreateSchema(), nil
	}

	itemsSlice, ok := itemsRaw.([]any)
	if !ok || len(itemsSlice) == 0 {
		return t.complexCreateSchema(), nil
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

	result, err := t.service.CreateComplexBatch(ctx, input.Items)
	if err != nil {
		return nil, err
	}
	return toResultMap(result)
}
