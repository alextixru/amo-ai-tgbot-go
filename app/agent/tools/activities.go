package tools

import (
	"encoding/json"
	"fmt"

	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/activities"
	models "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// ActivitiesTool реализует нативный ADK FunctionTool интерфейс для работы с активностями amoCRM.
// Shadow Tool паттерн: минимальная схема видна LLM, полная схема возвращается при первом вызове.
type ActivitiesTool struct {
	service activities.Service
}

// NewActivitiesTool создаёт новый ActivitiesTool с указанным сервисом.
func NewActivitiesTool(service activities.Service) *ActivitiesTool {
	return &ActivitiesTool{service: service}
}

// Name implements tool.Tool.
func (t *ActivitiesTool) Name() string {
	return "activities"
}

// Description implements tool.Tool.
func (t *ActivitiesTool) Description() string {
	return "Активности привязанные к сущностям amoCRM. " +
		"Layers: tasks, notes, calls, events, links, tags, subscriptions, talks. " +
		"Вызови с layer + action чтобы получить схему параметров."
}

// IsLongRunning implements tool.Tool.
func (t *ActivitiesTool) IsLongRunning() bool {
	return false
}

// Declaration implements toolinternal.FunctionTool (duck typing).
func (t *ActivitiesTool) Declaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"layer": {
					Type:        genai.TypeString,
					Description: "Слой активностей: tasks, notes, calls, events, files, links, tags, subscriptions, talks",
					Enum:        []string{"tasks", "notes", "calls", "events", "files", "links", "tags", "subscriptions", "talks"},
				},
				"action": {
					Type:        genai.TypeString,
					Description: "Действие: list, get, create, update, complete, link, unlink, delete, subscribe, unsubscribe, close",
				},
			},
			Required: []string{"layer", "action"},
		},
	}
}

// Run implements toolinternal.FunctionTool (duck typing).
func (t *ActivitiesTool) Run(ctx tool.Context, args any) (map[string]any, error) {
	m, ok := args.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("activities: неверный формат input")
	}
	result, err := t.handleActivitiesShadow(ctx, m)
	if err != nil {
		return nil, err
	}
	return toResultMap(result)
}

// activitiesSchemas содержит полные схемы параметров для каждой комбинации layer+action.
// Возвращается LLM при первом вызове без обязательных полей (Shadow Tool — schema mode).
var activitiesSchemas = map[string]map[string]any{
	// ─── TASKS ───────────────────────────────────────────────────────────────────
	"tasks:list": {
		"schema":          true,
		"tool":            "activities",
		"layer":           "tasks",
		"action":          "list",
		"description":     "Получить список задач. parent необязателен — если не указан, возвращаются задачи всех сущностей.",
		"required_fields": map[string]any{},
		"optional_fields": map[string]any{
			"parent": map[string]any{
				"type":        "object",
				"description": "Родительская сущность {type: leads|contacts|companies, id: number}",
			},
			"filter": map[string]any{
				"type":        "object",
				"description": "Фильтры задач",
				"fields": map[string]any{
					"limit":                  map[string]any{"type": "integer", "description": "Лимит (до 50)"},
					"page":                   map[string]any{"type": "integer", "description": "Страница"},
					"order":                  map[string]any{"type": "string", "description": "Сортировка: complete_till, created_at"},
					"order_dir":              map[string]any{"type": "string", "description": "Направление: asc, desc"},
					"ids":                    map[string]any{"type": "array[integer]", "description": "ID конкретных задач"},
					"responsible_user_names": map[string]any{"type": "array[string]", "description": "Имена ответственных"},
					"created_by_names":       map[string]any{"type": "array[string]", "description": "Имена создателей"},
					"is_completed":           map[string]any{"type": "boolean", "description": "Статус завершения"},
					"task_type":              map[string]any{"type": "string", "description": "Тип: follow_up, meeting"},
					"date_range":             map[string]any{"type": "string", "description": "today, tomorrow, overdue, this_week, next_week"},
					"query":                  map[string]any{"type": "string", "description": "Поисковый запрос"},
					"updated_at":             map[string]any{"type": "integer", "description": "Фильтр по дате изменения от (Unix timestamp)"},
					"updated_at_to":          map[string]any{"type": "integer", "description": "Фильтр по дате изменения до (Unix timestamp)"},
				},
			},
			"with": map[string]any{"type": "array[string]", "description": "Связанные данные: leads, contacts"},
		},
		"example": map[string]any{
			"layer":  "tasks",
			"action": "list",
			"parent": map[string]any{"type": "leads", "id": 12345},
			"filter": map[string]any{"is_completed": false, "date_range": "today"},
		},
	},
	"tasks:get": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "tasks",
		"action":      "get",
		"description": "Получить задачу по ID.",
		"required_fields": map[string]any{
			"id": map[string]any{"type": "integer", "description": "ID задачи"},
		},
		"optional_fields": map[string]any{
			"with": map[string]any{"type": "array[string]", "description": "Связанные данные: leads, contacts"},
		},
		"example": map[string]any{
			"layer":  "tasks",
			"action": "get",
			"id":     42,
		},
	},
	"tasks:create": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "tasks",
		"action":      "create",
		"description": "Создать задачу (одну или пакет) для сущности amoCRM.",
		"required_fields": map[string]any{
			"parent": map[string]any{
				"type":        "object",
				"description": "Родительская сущность — обязательна",
				"fields": map[string]any{
					"type": map[string]any{"type": "string", "description": "leads | contacts | companies"},
					"id":   map[string]any{"type": "integer", "description": "ID сущности"},
				},
			},
		},
		"optional_fields": map[string]any{
			"task_data": map[string]any{
				"type":        "object",
				"description": "Данные одной задачи",
				"fields": map[string]any{
					"text":                  map[string]any{"type": "string", "description": "Текст задачи (обязательно внутри task_data)"},
					"responsible_user_name": map[string]any{"type": "string", "description": "Имя ответственного"},
					"task_type":             map[string]any{"type": "string", "description": "follow_up | meeting"},
					"deadline":              map[string]any{"type": "string", "description": "today, tomorrow, 'in 2 hours', 'in 3 days', '2024-01-15', '2024-01-15T14:00'"},
				},
			},
			"tasks_data": map[string]any{
				"type":        "array[object]",
				"description": "Массив задач для пакетного создания (вместо task_data). Каждый элемент — TaskData.",
			},
		},
		"example": map[string]any{
			"layer":  "tasks",
			"action": "create",
			"parent": map[string]any{"type": "leads", "id": 12345},
			"task_data": map[string]any{
				"text":                  "Позвонить клиенту",
				"responsible_user_name": "Иван Петров",
				"task_type":             "follow_up",
				"deadline":              "tomorrow",
			},
		},
	},
	"tasks:update": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "tasks",
		"action":      "update",
		"description": "Обновить задачу по ID.",
		"required_fields": map[string]any{
			"id":        map[string]any{"type": "integer", "description": "ID задачи"},
			"task_data": map[string]any{"type": "object", "description": "Данные для обновления (text, responsible_user_name, task_type, deadline)"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":     "tasks",
			"action":    "update",
			"id":        42,
			"task_data": map[string]any{"deadline": "in 2 hours"},
		},
	},
	"tasks:complete": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "tasks",
		"action":      "complete",
		"description": "Завершить задачу по ID с текстом результата.",
		"required_fields": map[string]any{
			"id": map[string]any{"type": "integer", "description": "ID задачи"},
		},
		"optional_fields": map[string]any{
			"result_text": map[string]any{"type": "string", "description": "Текст результата выполнения"},
		},
		"example": map[string]any{
			"layer":       "tasks",
			"action":      "complete",
			"id":          42,
			"result_text": "Клиент перезвонил, договорились о встрече",
		},
	},

	// ─── NOTES ───────────────────────────────────────────────────────────────────
	"notes:list": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "notes",
		"action":      "list",
		"description": "Получить список примечаний для сущности.",
		"required_fields": map[string]any{
			"parent": map[string]any{
				"type":        "object",
				"description": "Родительская сущность {type: leads|contacts|companies, id: number}",
			},
		},
		"optional_fields": map[string]any{
			"notes_filter": map[string]any{
				"type":        "object",
				"description": "Фильтры примечаний",
				"fields": map[string]any{
					"limit":      map[string]any{"type": "integer", "description": "Лимит (до 50)"},
					"page":       map[string]any{"type": "integer", "description": "Страница"},
					"ids":        map[string]any{"type": "array[integer]", "description": "ID конкретных примечаний"},
					"note_types": map[string]any{"type": "array[string]", "description": "Типы: common, call_in, call_out, service_message"},
					"updated_at": map[string]any{"type": "integer", "description": "Фильтр по дате изменения (timestamp)"},
				},
			},
			"with": map[string]any{"type": "array[string]", "description": "Связанные данные"},
		},
		"example": map[string]any{
			"layer":  "notes",
			"action": "list",
			"parent": map[string]any{"type": "leads", "id": 12345},
		},
	},
	"notes:get": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "notes",
		"action":      "get",
		"description": "Получить примечание по ID.",
		"required_fields": map[string]any{
			"parent": map[string]any{"type": "object", "description": "{type: leads|contacts|companies, id: number}"},
			"id":     map[string]any{"type": "integer", "description": "ID примечания"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":  "notes",
			"action": "get",
			"parent": map[string]any{"type": "leads", "id": 12345},
			"id":     99,
		},
	},
	"notes:create": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "notes",
		"action":      "create",
		"description": "Создать примечание (одно или пакет) для сущности.",
		"required_fields": map[string]any{
			"parent": map[string]any{
				"type":        "object",
				"description": "{type: leads|contacts|companies, id: number}",
			},
		},
		"optional_fields": map[string]any{
			"note_data": map[string]any{
				"type":        "object",
				"description": "Данные одного примечания",
				"fields": map[string]any{
					"text":      map[string]any{"type": "string", "description": "Текст примечания (обязательно внутри note_data)"},
					"note_type": map[string]any{"type": "string", "description": "common | call_in | call_out | service_message (по умолчанию common)"},
				},
			},
			"notes_data": map[string]any{
				"type":        "array[object]",
				"description": "Массив примечаний для пакетного создания",
			},
		},
		"example": map[string]any{
			"layer":  "notes",
			"action": "create",
			"parent": map[string]any{"type": "leads", "id": 12345},
			"note_data": map[string]any{
				"text":      "Клиент заинтересован в продукте",
				"note_type": "common",
			},
		},
	},
	"notes:update": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "notes",
		"action":      "update",
		"description": "Обновить примечание по ID.",
		"required_fields": map[string]any{
			"parent":    map[string]any{"type": "object", "description": "{type: leads|contacts|companies, id: number}"},
			"id":        map[string]any{"type": "integer", "description": "ID примечания"},
			"note_data": map[string]any{"type": "object", "description": "Данные для обновления (text, note_type)"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":     "notes",
			"action":    "update",
			"parent":    map[string]any{"type": "leads", "id": 12345},
			"id":        99,
			"note_data": map[string]any{"text": "Обновлённый текст примечания"},
		},
	},

	// ─── CALLS ───────────────────────────────────────────────────────────────────
	"calls:create": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "calls",
		"action":      "create",
		"description": "Создать запись о звонке для сущности amoCRM. Единственный поддерживаемый action для calls.",
		"required_fields": map[string]any{
			"parent": map[string]any{
				"type":        "object",
				"description": "{type: leads|contacts|companies, id: number}",
			},
			"call_data": map[string]any{
				"type":        "object",
				"description": "Данные звонка",
				"fields": map[string]any{
					"direction": map[string]any{"type": "string", "description": "inbound | outbound"},
					"duration":  map[string]any{"type": "integer", "description": "Длительность в секундах"},
					"phone":     map[string]any{"type": "string", "description": "Номер телефона"},
				},
			},
		},
		"optional_fields": map[string]any{
			"call_data": map[string]any{
				"type":        "object",
				"description": "Дополнительные поля звонка",
				"fields": map[string]any{
					"source":      map[string]any{"type": "string", "description": "Источник звонка"},
					"call_result": map[string]any{"type": "string", "description": "Результат звонка"},
					"call_status": map[string]any{"type": "integer", "description": "1=оставить_сообщение, 2=перезвонить, 3=недоступен, 4=занято, 5=неверный_номер, 6=нет_ответа, 7=успешный_звонок"},
					"unique_id":   map[string]any{"type": "string", "description": "Уникальный ID звонка"},
					"record_url":  map[string]any{"type": "string", "description": "Ссылка на запись звонка"},
				},
			},
		},
		"example": map[string]any{
			"layer":  "calls",
			"action": "create",
			"parent": map[string]any{"type": "contacts", "id": 9876},
			"call_data": map[string]any{
				"direction":   "outbound",
				"duration":    120,
				"phone":       "+79001234567",
				"call_status": 7,
				"call_result": "Договорились о встрече",
			},
		},
	},

	// ─── EVENTS ──────────────────────────────────────────────────────────────────
	"events:list": {
		"schema":          true,
		"tool":            "activities",
		"layer":           "events",
		"action":          "list",
		"description":     "Получить список событий (история изменений). parent необязателен.",
		"required_fields": map[string]any{},
		"optional_fields": map[string]any{
			"parent": map[string]any{
				"type":        "object",
				"description": "{type: leads|contacts|companies, id: number}",
			},
			"events_filter": map[string]any{
				"type":        "object",
				"description": "Фильтры событий",
				"fields": map[string]any{
					"limit":            map[string]any{"type": "integer", "description": "Лимит (до 100)"},
					"page":             map[string]any{"type": "integer", "description": "Страница"},
					"types":            map[string]any{"type": "array[string]", "description": "Типы: lead_added, lead_status_changed, contact_added и др."},
					"created_by_names": map[string]any{"type": "array[string]", "description": "Имена создателей событий"},
					"with":             map[string]any{"type": "array[string]", "description": "contact_name, lead_name, company_name, note и др."},
				},
			},
		},
		"example": map[string]any{
			"layer":  "events",
			"action": "list",
			"parent": map[string]any{"type": "leads", "id": 12345},
			"events_filter": map[string]any{
				"types": []string{"lead_status_changed", "lead_note_added"},
				"limit": 20,
			},
		},
	},
	"events:get": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "events",
		"action":      "get",
		"description": "Получить событие по ID.",
		"required_fields": map[string]any{
			"id": map[string]any{"type": "integer", "description": "ID события"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":  "events",
			"action": "get",
			"id":     777,
		},
	},

	// ─── FILES ───────────────────────────────────────────────────────────────────
	"files:list": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "files",
		"action":      "list",
		"description": "Получить список файлов прикреплённых к сущности.",
		"required_fields": map[string]any{
			"parent": map[string]any{
				"type":        "object",
				"description": "{type: leads|contacts|companies, id: number}",
			},
		},
		"optional_fields": map[string]any{
			"files_filter": map[string]any{
				"type":        "object",
				"description": "Фильтры файлов",
				"fields": map[string]any{
					"limit":      map[string]any{"type": "integer", "description": "Лимит (до 50)"},
					"page":       map[string]any{"type": "integer", "description": "Страница"},
					"extensions": map[string]any{"type": "array[string]", "description": "Расширения: pdf, docx, xlsx"},
					"term":       map[string]any{"type": "string", "description": "Поиск по имени файла"},
					"uuid":       map[string]any{"type": "string", "description": "UUID конкретного файла"},
				},
			},
		},
		"example": map[string]any{
			"layer":  "files",
			"action": "list",
			"parent": map[string]any{"type": "leads", "id": 12345},
		},
	},
	"files:link": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "files",
		"action":      "link",
		"description": "Прикрепить файлы (по UUID из Drive) к сущности.",
		"required_fields": map[string]any{
			"parent":     map[string]any{"type": "object", "description": "{type: leads|contacts|companies, id: number}"},
			"file_uuids": map[string]any{"type": "array[string]", "description": "UUID файлов для прикрепления"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":      "files",
			"action":     "link",
			"parent":     map[string]any{"type": "leads", "id": 12345},
			"file_uuids": []string{"uuid-1", "uuid-2"},
		},
	},
	"files:unlink": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "files",
		"action":      "unlink",
		"description": "Открепить файл (по UUID) от сущности.",
		"required_fields": map[string]any{
			"parent":    map[string]any{"type": "object", "description": "{type: leads|contacts|companies, id: number}"},
			"file_uuid": map[string]any{"type": "string", "description": "UUID файла для открепления"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":     "files",
			"action":    "unlink",
			"parent":    map[string]any{"type": "leads", "id": 12345},
			"file_uuid": "550e8400-e29b-41d4-a716-446655440000",
		},
	},

	// ─── LINKS ───────────────────────────────────────────────────────────────────
	"links:list": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "links",
		"action":      "list",
		"description": "Получить список связанных сущностей для parent.",
		"required_fields": map[string]any{
			"parent": map[string]any{
				"type":        "object",
				"description": "{type: leads|contacts|companies, id: number}",
			},
		},
		"optional_fields": map[string]any{
			"links_filter": map[string]any{
				"type":        "object",
				"description": "Фильтры связей",
				"fields": map[string]any{
					"to_entity_type": map[string]any{"type": "string", "description": "Тип целевой сущности: leads, contacts, companies, catalog_elements"},
					"to_entity_id":   map[string]any{"type": "integer", "description": "ID конкретной целевой сущности"},
				},
			},
		},
		"example": map[string]any{
			"layer":  "links",
			"action": "list",
			"parent": map[string]any{"type": "leads", "id": 12345},
		},
	},
	"links:link": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "links",
		"action":      "link",
		"description": "Связать сущности. Поддерживает одиночное (link_to) и пакетное (links_to) связывание.",
		"required_fields": map[string]any{
			"parent": map[string]any{
				"type":        "object",
				"description": "{type: leads|contacts|companies, id: number}",
			},
		},
		"optional_fields": map[string]any{
			"link_to": map[string]any{
				"type":        "object",
				"description": "Одна цель связывания {type: leads|contacts|companies, id: number}",
			},
			"links_to": map[string]any{
				"type":        "array[object]",
				"description": "Несколько целей связывания (вместо link_to)",
			},
		},
		"example": map[string]any{
			"layer":   "links",
			"action":  "link",
			"parent":  map[string]any{"type": "leads", "id": 12345},
			"link_to": map[string]any{"type": "contacts", "id": 9876},
		},
	},
	"links:unlink": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "links",
		"action":      "unlink",
		"description": "Разорвать связь между сущностями.",
		"required_fields": map[string]any{
			"parent":  map[string]any{"type": "object", "description": "{type: leads|contacts|companies, id: number}"},
			"link_to": map[string]any{"type": "object", "description": "{type: leads|contacts|companies, id: number}"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":   "links",
			"action":  "unlink",
			"parent":  map[string]any{"type": "leads", "id": 12345},
			"link_to": map[string]any{"type": "contacts", "id": 9876},
		},
	},

	// ─── TAGS ────────────────────────────────────────────────────────────────────
	"tags:list": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "tags",
		"action":      "list",
		"description": "Получить список тегов для типа сущности.",
		"required_fields": map[string]any{
			"parent": map[string]any{
				"type":        "object",
				"description": "Нужен только type для указания типа сущности {type: leads|contacts|companies}. id не обязателен.",
			},
		},
		"optional_fields": map[string]any{
			"tags_filter": map[string]any{
				"type":        "object",
				"description": "Фильтры тегов",
				"fields": map[string]any{
					"limit": map[string]any{"type": "integer", "description": "Лимит (до 50)"},
					"page":  map[string]any{"type": "integer", "description": "Страница"},
					"query": map[string]any{"type": "string", "description": "Поиск по названию (частичное совпадение)"},
					"name":  map[string]any{"type": "string", "description": "Фильтр по точному названию"},
					"ids":   map[string]any{"type": "array[integer]", "description": "ID конкретных тегов"},
				},
			},
		},
		"example": map[string]any{
			"layer":  "tags",
			"action": "list",
			"parent": map[string]any{"type": "leads"},
		},
	},
	"tags:create": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "tags",
		"action":      "create",
		"description": "Создать тег (один или пакет) для типа сущности.",
		"required_fields": map[string]any{
			"parent": map[string]any{
				"type":        "object",
				"description": "{type: leads|contacts|companies} — нужен только type",
			},
		},
		"optional_fields": map[string]any{
			"tag_name": map[string]any{
				"type":        "string",
				"description": "Название одного тега для создания",
			},
			"tag_names": map[string]any{
				"type":        "array[string]",
				"description": "Список названий тегов для пакетного создания",
			},
		},
		"example": map[string]any{
			"layer":    "tags",
			"action":   "create",
			"parent":   map[string]any{"type": "leads"},
			"tag_name": "VIP",
		},
	},
	"tags:delete": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "tags",
		"action":      "delete",
		"description": "Удалить тег по ID или имени для типа сущности.",
		"required_fields": map[string]any{
			"parent": map[string]any{
				"type":        "object",
				"description": "{type: leads|contacts|companies} — нужен только type",
			},
		},
		"optional_fields": map[string]any{
			"tag_id":   map[string]any{"type": "integer", "description": "ID тега для удаления"},
			"tag_name": map[string]any{"type": "string", "description": "Название тега для удаления (альтернатива tag_id)"},
		},
		"example": map[string]any{
			"layer":    "tags",
			"action":   "delete",
			"parent":   map[string]any{"type": "leads"},
			"tag_name": "VIP",
		},
	},

	// ─── SUBSCRIPTIONS ───────────────────────────────────────────────────────────
	"subscriptions:list": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "subscriptions",
		"action":      "list",
		"description": "Получить список подписчиков на уведомления для сущности.",
		"required_fields": map[string]any{
			"parent": map[string]any{
				"type":        "object",
				"description": "{type: leads|contacts|companies, id: number}",
			},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":  "subscriptions",
			"action": "list",
			"parent": map[string]any{"type": "leads", "id": 12345},
		},
	},
	"subscriptions:subscribe": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "subscriptions",
		"action":      "subscribe",
		"description": "Подписать пользователей на уведомления о сущности.",
		"required_fields": map[string]any{
			"parent":     map[string]any{"type": "object", "description": "{type: leads|contacts|companies, id: number}"},
			"user_names": map[string]any{"type": "array[string]", "description": "Имена пользователей для подписки"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":      "subscriptions",
			"action":     "subscribe",
			"parent":     map[string]any{"type": "leads", "id": 12345},
			"user_names": []string{"Иван Петров", "Мария Сидорова"},
		},
	},
	"subscriptions:unsubscribe": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "subscriptions",
		"action":      "unsubscribe",
		"description": "Отписать пользователя от уведомлений о сущности.",
		"required_fields": map[string]any{
			"parent":    map[string]any{"type": "object", "description": "{type: leads|contacts|companies, id: number}"},
			"user_name": map[string]any{"type": "string", "description": "Имя пользователя для отписки"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":     "subscriptions",
			"action":    "unsubscribe",
			"parent":    map[string]any{"type": "leads", "id": 12345},
			"user_name": "Иван Петров",
		},
	},

	// ─── TALKS ───────────────────────────────────────────────────────────────────
	"talks:get": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "talks",
		"action":      "get",
		"description": "Получить информацию о чате (беседе) по talk_id.",
		"required_fields": map[string]any{
			"talk_id": map[string]any{"type": "string", "description": "ID чата"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":   "talks",
			"action":  "get",
			"talk_id": "abc123",
		},
	},
	"talks:close": {
		"schema":      true,
		"tool":        "activities",
		"layer":       "talks",
		"action":      "close",
		"description": "Закрыть активную беседу (чат) по talk_id.",
		"required_fields": map[string]any{
			"talk_id": map[string]any{"type": "string", "description": "ID чата"},
		},
		"optional_fields": map[string]any{
			"force_close": map[string]any{"type": "boolean", "description": "Принудительное закрытие (true/false)"},
		},
		"example": map[string]any{
			"layer":       "talks",
			"action":      "close",
			"talk_id":     "abc123",
			"force_close": false,
		},
	},
}

// activitiesSchemaKey формирует ключ layer:action для поиска в activitiesSchemas.
func activitiesSchemaKey(layer, action string) string {
	return layer + ":" + action
}

// isActivitiesSchemaMode определяет режим работы (schema/execute) по наличию обязательных полей.
// Возвращает true (schema mode) если обязательные поля для layer+action отсутствуют.
func isActivitiesSchemaMode(layer, action string, m map[string]any) bool {
	parentRaw := m["parent"]
	parentMap, _ := parentRaw.(map[string]any)
	parentType, _ := parentMap["type"].(string)
	parentIDRaw := parentMap["id"]
	parentID := 0
	switch v := parentIDRaw.(type) {
	case float64:
		parentID = int(v)
	case int:
		parentID = v
	}

	hasParentType := parentType != ""
	hasParentID := parentID != 0
	hasParent := hasParentType && hasParentID

	id := 0
	switch v := m["id"].(type) {
	case float64:
		id = int(v)
	case int:
		id = v
	}
	hasID := id != 0

	talkID, _ := m["talk_id"].(string)

	switch layer {
	case "tasks":
		switch action {
		case "list":
			// нет обязательных → всегда execute
			return false
		case "get":
			return !hasID
		case "create":
			return !hasParent
		case "update":
			return !hasID || m["task_data"] == nil
		case "complete":
			return !hasID
		}

	case "notes":
		switch action {
		case "list":
			return !hasParent
		case "get":
			return !hasParent || !hasID
		case "create":
			return !hasParent
		case "update":
			return !hasParent || !hasID || m["note_data"] == nil
		}

	case "calls":
		if action == "create" {
			return !hasParent || m["call_data"] == nil
		}

	case "events":
		switch action {
		case "list":
			// нет обязательных → всегда execute
			return false
		case "get":
			return !hasID
		}

	case "files":
		switch action {
		case "list":
			return !hasParent
		case "link":
			fileUUIDs, _ := m["file_uuids"].([]any)
			return !hasParent || len(fileUUIDs) == 0
		case "unlink":
			fileUUID, _ := m["file_uuid"].(string)
			return !hasParent || fileUUID == ""
		}

	case "links":
		switch action {
		case "list":
			return !hasParent
		case "link":
			return !hasParent
		case "unlink":
			return !hasParent || m["link_to"] == nil
		}

	case "tags":
		switch action {
		case "list":
			return !hasParentType
		case "create":
			tagName, _ := m["tag_name"].(string)
			tagNames, _ := m["tag_names"].([]any)
			return !hasParentType || (tagName == "" && len(tagNames) == 0)
		case "delete":
			tagName, _ := m["tag_name"].(string)
			tagID := 0
			switch v := m["tag_id"].(type) {
			case float64:
				tagID = int(v)
			case int:
				tagID = v
			}
			return !hasParentType || (tagName == "" && tagID == 0)
		}

	case "subscriptions":
		switch action {
		case "list":
			return !hasParent
		case "subscribe":
			userNames, _ := m["user_names"].([]any)
			return !hasParent || len(userNames) == 0
		case "unsubscribe":
			userName, _ := m["user_name"].(string)
			return !hasParent || userName == ""
		}

	case "talks":
		return talkID == ""
	}

	// неизвестный layer/action → schema mode
	return true
}

// activitiesSchemaResponse строит schema response для данного layer+action.
// Добавляет available_values из сервиса.
func (t *ActivitiesTool) activitiesSchemaResponse(layer, action string) map[string]any {
	key := activitiesSchemaKey(layer, action)
	schema, ok := activitiesSchemas[key]
	if !ok {
		return map[string]any{
			"schema": true,
			"tool":   "activities",
			"error":  fmt.Sprintf("unknown layer+action: %s+%s", layer, action),
			"available_layers_and_actions": map[string]any{
				"tasks":         []string{"list", "get", "create", "update", "complete"},
				"notes":         []string{"list", "get", "create", "update"},
				"calls":         []string{"create"},
				"events":        []string{"list", "get"},
				"files":         []string{"list", "link", "unlink"},
				"links":         []string{"list", "link", "unlink"},
				"tags":          []string{"list", "create", "delete"},
				"subscriptions": []string{"list", "subscribe", "unsubscribe"},
				"talks":         []string{"get", "close"},
			},
		}
	}

	// Добавляем available_values для layers где они релевантны
	result := make(map[string]any, len(schema)+1)
	for k, v := range schema {
		result[k] = v
	}

	switch layer {
	case "tasks", "subscriptions":
		result["available_values"] = map[string]any{
			"users":        t.service.UserNames(),
			"entity_types": []string{"leads", "contacts", "companies"},
		}
	case "notes", "calls", "files", "links":
		result["available_values"] = map[string]any{
			"entity_types": []string{"leads", "contacts", "companies"},
		}
	case "events":
		result["available_values"] = map[string]any{
			"entity_types": []string{"leads", "contacts", "companies"},
			"event_types":  []string{"lead_added", "lead_status_changed", "lead_deleted", "contact_added", "contact_deleted", "company_added", "company_deleted", "lead_note_added", "contact_note_added", "task_added", "task_completed", "incoming_chat_message"},
		}
	case "tags":
		result["available_values"] = map[string]any{
			"entity_types": []string{"leads", "contacts", "companies"},
		}
	case "talks":
		result["available_values"] = map[string]any{}
	}

	return result
}

// handleActivitiesShadow реализует Shadow Tool логику для activities.
func (t *ActivitiesTool) handleActivitiesShadow(ctx tool.Context, m map[string]any) (any, error) {
	layer, _ := m["layer"].(string)
	action, _ := m["action"].(string)

	// Нет layer или action — возвращаем верхнеуровневую схему
	if layer == "" || action == "" {
		return map[string]any{
			"schema": true,
			"tool":   "activities",
			"error":  "layer and action are required",
			"available_layers_and_actions": map[string]any{
				"tasks":         []string{"list", "get", "create", "update", "complete"},
				"notes":         []string{"list", "get", "create", "update"},
				"calls":         []string{"create"},
				"events":        []string{"list", "get"},
				"files":         []string{"list", "link", "unlink"},
				"links":         []string{"list", "link", "unlink"},
				"tags":          []string{"list", "create", "delete"},
				"subscriptions": []string{"list", "subscribe", "unsubscribe"},
				"talks":         []string{"get", "close"},
			},
			"hint": "Укажи layer + action чтобы получить схему параметров",
		}, nil
	}

	// Определяем режим: schema или execute
	if isActivitiesSchemaMode(layer, action, m) {
		return t.activitiesSchemaResponse(layer, action), nil
	}

	// Execute mode: JSON roundtrip map → ActivitiesInput → существующий handler
	b, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("activities: marshal input: %w", err)
	}

	var input models.ActivitiesInput
	if err := json.Unmarshal(b, &input); err != nil {
		return nil, fmt.Errorf("activities: unmarshal input: %w", err)
	}

	return t.handleActivities(ctx, input)
}

func (t *ActivitiesTool) handleActivities(ctx tool.Context, input models.ActivitiesInput) (any, error) {
	// Валидация parent: теперь опционально для некоторых действий или при наличии фильтра
	// Для большинства действий (кроме list) parent всё еще обязателен
	if input.Action != "list" && (input.Parent == nil || input.Parent.Type == "" || input.Parent.ID == 0) {
		if input.Action == "create" || input.Action == "update" || input.Action == "complete" || input.Action == "link" || input.Action == "unlink" || input.Action == "subscribe" || input.Action == "unsubscribe" {
			if input.Action == "create" && (input.Parent == nil || input.Parent.ID == 0) {
				return nil, fmt.Errorf("parent.id is required for create action")
			}
		}
	}

	switch input.Layer {
	case "tasks":
		return t.handleTasks(ctx, input)
	case "notes":
		return t.handleNotes(ctx, input)
	case "calls":
		return t.handleCalls(ctx, input)
	case "events":
		return t.handleEvents(ctx, input)
	case "files":
		return t.handleFiles(ctx, input)
	case "links":
		return t.handleLinks(ctx, input)
	case "tags":
		return t.handleTags(ctx, input)
	case "subscriptions":
		return t.handleSubscriptions(ctx, input)
	case "talks":
		return t.handleTalks(ctx, input)
	default:
		return nil, fmt.Errorf("unknown layer: %s", input.Layer)
	}
}

func (t *ActivitiesTool) handleTasks(ctx tool.Context, input models.ActivitiesInput) (any, error) {
	switch input.Action {
	case "list":
		return t.service.ListTasks(ctx, input.Parent, input.Filter, input.With)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required")
		}
		return t.service.GetTask(ctx, input.ID, input.With)
	case "create":
		if input.Parent == nil {
			return nil, fmt.Errorf("parent is required for create")
		}
		// Batch create
		if len(input.TasksData) > 0 {
			return t.service.CreateTasks(ctx, *input.Parent, input.TasksData)
		}
		// Single create
		if input.TaskData == nil {
			return nil, fmt.Errorf("task_data or tasks_data is required")
		}
		return t.service.CreateTask(ctx, *input.Parent, input.TaskData)
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required")
		}
		if input.TaskData == nil {
			return nil, fmt.Errorf("task_data is required")
		}
		return t.service.UpdateTask(ctx, input.ID, input.TaskData)
	case "complete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required")
		}
		return t.service.CompleteTask(ctx, input.ID, input.ResultText)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (t *ActivitiesTool) handleNotes(ctx tool.Context, input models.ActivitiesInput) (any, error) {
	if input.Parent == nil {
		return nil, fmt.Errorf("parent is required for notes")
	}
	switch input.Action {
	case "list":
		return t.service.ListNotes(ctx, *input.Parent, input.NotesFilter, input.With)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required")
		}
		return t.service.GetNote(ctx, input.Parent.Type, input.ID)
	case "create":
		// Batch create
		if len(input.NotesData) > 0 {
			return t.service.CreateNotes(ctx, *input.Parent, input.NotesData)
		}
		// Single create
		if input.NoteData == nil {
			return nil, fmt.Errorf("note_data or notes_data is required")
		}
		return t.service.CreateNote(ctx, *input.Parent, input.NoteData)
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required")
		}
		if input.NoteData == nil {
			return nil, fmt.Errorf("note_data is required")
		}
		return t.service.UpdateNote(ctx, input.Parent.Type, input.ID, input.NoteData)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (t *ActivitiesTool) handleCalls(ctx tool.Context, input models.ActivitiesInput) (any, error) {
	if input.Action != "create" {
		return nil, fmt.Errorf("calls only supports 'create' action")
	}
	if input.Parent == nil {
		return nil, fmt.Errorf("parent is required for calls")
	}
	if input.CallData == nil {
		return nil, fmt.Errorf("call_data is required")
	}
	return t.service.CreateCall(ctx, *input.Parent, input.CallData)
}

func (t *ActivitiesTool) handleEvents(ctx tool.Context, input models.ActivitiesInput) (any, error) {
	switch input.Action {
	case "list":
		return t.service.ListEvents(ctx, input.Parent, input.EventsFilter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required")
		}
		return t.service.GetEvent(ctx, input.ID)
	default:
		return nil, fmt.Errorf("events only supports 'list' and 'get' actions")
	}
}

func (t *ActivitiesTool) handleFiles(ctx tool.Context, input models.ActivitiesInput) (any, error) {
	if input.Parent == nil {
		return nil, fmt.Errorf("parent is required for files")
	}
	switch input.Action {
	case "list":
		return t.service.ListFiles(ctx, *input.Parent, input.FilesFilter)
	case "link":
		if len(input.FileUUIDs) == 0 {
			return nil, fmt.Errorf("file_uuids is required")
		}
		return t.service.LinkFiles(ctx, *input.Parent, input.FileUUIDs)
	case "unlink":
		if input.FileUUID == "" {
			return nil, fmt.Errorf("file_uuid is required")
		}
		return nil, t.service.UnlinkFile(ctx, *input.Parent, input.FileUUID)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (t *ActivitiesTool) handleLinks(ctx tool.Context, input models.ActivitiesInput) (any, error) {
	if input.Parent == nil {
		return nil, fmt.Errorf("parent is required for links")
	}
	switch input.Action {
	case "list":
		return t.service.ListLinks(ctx, *input.Parent, input.LinksFilter)
	case "link":
		// Batch link
		if len(input.LinksTo) > 0 {
			return t.service.LinkEntities(ctx, *input.Parent, input.LinksTo)
		}
		// Single link
		if input.LinkTo == nil {
			return nil, fmt.Errorf("link_to or links_to is required")
		}
		return t.service.LinkEntity(ctx, *input.Parent, input.LinkTo)
	case "unlink":
		if input.LinkTo == nil {
			return nil, fmt.Errorf("link_to is required")
		}
		return nil, t.service.UnlinkEntity(ctx, *input.Parent, input.LinkTo)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (t *ActivitiesTool) handleTags(ctx tool.Context, input models.ActivitiesInput) (any, error) {
	entityType := ""
	if input.Parent != nil {
		entityType = input.Parent.Type
	}
	if entityType == "" {
		return nil, fmt.Errorf("parent.type is required for tags")
	}

	switch input.Action {
	case "list":
		return t.service.ListTags(ctx, entityType, input.TagsFilter)
	case "create":
		// Batch create
		if len(input.TagNames) > 0 {
			return t.service.CreateTags(ctx, entityType, input.TagNames)
		}
		// Single create
		if input.TagName == "" {
			return nil, fmt.Errorf("tag_name or tag_names is required")
		}
		return t.service.CreateTag(ctx, entityType, input.TagName)
	case "delete":
		if input.TagID == 0 && input.TagName != "" {
			return nil, t.service.DeleteTagByName(ctx, entityType, input.TagName)
		}
		if input.TagID == 0 {
			return nil, fmt.Errorf("tag_id or tag_name is required for delete")
		}
		return nil, t.service.DeleteTag(ctx, entityType, input.TagID)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (t *ActivitiesTool) handleSubscriptions(ctx tool.Context, input models.ActivitiesInput) (any, error) {
	if input.Parent == nil {
		return nil, fmt.Errorf("parent is required for subscriptions")
	}
	switch input.Action {
	case "list":
		return t.service.ListSubscriptions(ctx, *input.Parent)
	case "subscribe":
		if len(input.UserNames) == 0 {
			return nil, fmt.Errorf("user_names is required")
		}
		return t.service.Subscribe(ctx, *input.Parent, input.UserNames)
	case "unsubscribe":
		if input.UserName == "" {
			return nil, fmt.Errorf("user_name is required")
		}
		return nil, t.service.Unsubscribe(ctx, *input.Parent, input.UserName)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (t *ActivitiesTool) handleTalks(ctx tool.Context, input models.ActivitiesInput) (any, error) {
	switch input.Action {
	case "get":
		if input.TalkID == "" {
			return nil, fmt.Errorf("talk_id is required")
		}
		return t.service.GetTalk(ctx, input.TalkID)
	case "close":
		if input.TalkID == "" {
			return nil, fmt.Errorf("talk_id is required")
		}
		return nil, t.service.CloseTalk(ctx, input.TalkID, input.ForceClose)
	default:
		return nil, fmt.Errorf("unknown action: %s (talks supports 'get' and 'close')", input.Action)
	}
}
