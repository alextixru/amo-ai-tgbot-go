package tools

import (
	"encoding/json"
	"fmt"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/admin_users"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// AdminUsersTool реализует нативный ADK FunctionTool интерфейс для управления
// пользователями и ролями amoCRM.
// Shadow Tool паттерн: минимальная схема видна LLM, полная схема возвращается при первом вызове.
type AdminUsersTool struct {
	service admin_users.Service
}

// NewAdminUsersTool создаёт новый AdminUsersTool с указанным сервисом.
func NewAdminUsersTool(service admin_users.Service) *AdminUsersTool {
	return &AdminUsersTool{service: service}
}

// Name implements tool.Tool.
func (t *AdminUsersTool) Name() string {
	return "admin_users"
}

// Description implements tool.Tool.
func (t *AdminUsersTool) Description() string {
	return "Управление пользователями и ролями amoCRM. " +
		"Layers: users (пользователи), roles (роли). " +
		"Actions для users: list, search, get, create. " +
		"Actions для roles: list, search, get, create, update, delete. " +
		"Вызови с layer + action чтобы получить схему параметров."
}

// IsLongRunning implements tool.Tool.
func (t *AdminUsersTool) IsLongRunning() bool {
	return false
}

// Declaration implements toolinternal.FunctionTool (duck typing).
func (t *AdminUsersTool) Declaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"layer": {
					Type:        genai.TypeString,
					Description: "Слой: users (пользователи) или roles (роли)",
					Enum:        []string{"users", "roles"},
				},
				"action": {
					Type:        genai.TypeString,
					Description: "Действие: list, search, get, create, update, delete",
				},
			},
			Required: []string{"layer", "action"},
		},
	}
}

// Run implements toolinternal.FunctionTool (duck typing).
func (t *AdminUsersTool) Run(ctx tool.Context, args any) (map[string]any, error) {
	raw, ok := args.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("admin_users: invalid input type %T", args)
	}

	layer, _ := raw["layer"].(string)
	action, _ := raw["action"].(string)

	if layer == "" {
		result := map[string]any{
			"schema":      true,
			"tool":        "admin_users",
			"error":       "layer is required",
			"layers":      []string{"users", "roles"},
			"description": "Укажите layer (users или roles) и action чтобы получить схему параметров.",
		}
		return result, nil
	}
	if action == "" {
		result := map[string]any{
			"schema":      true,
			"tool":        "admin_users",
			"layer":       layer,
			"error":       "action is required",
			"description": "Укажите action чтобы получить схему параметров.",
		}
		return result, nil
	}

	key := layer + "." + action

	// Проверяем наличие обязательных полей для данного layer.action
	required, known := adminUsersRequiredFields[key]
	if !known {
		return nil, fmt.Errorf("неизвестная комбинация layer=%q action=%q. Вызови без обязательных полей чтобы получить схему.", layer, action)
	}

	// Определяем режим: schema или execute
	isSchemaMode := false
	for _, field := range required {
		if _, ok := raw[field]; !ok {
			isSchemaMode = true
			break
		}
	}

	if isSchemaMode {
		if schema, ok := adminUsersSchema[key]; ok {
			return schema, nil
		}
		return map[string]any{
			"schema": true,
			"tool":   "admin_users",
			"layer":  layer,
			"action": action,
			"error":  "schema not found for this combination",
		}, nil
	}

	// Execute mode: json-roundtrip map → AdminUsersInput
	parsedInput, err := adminUsersMapToInput(raw)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора параметров: %w", err)
	}

	var result any
	switch layer {
	case "users":
		result, err = t.handleUsers(ctx, parsedInput)
	case "roles":
		result, err = t.handleRoles(ctx, parsedInput)
	default:
		return nil, fmt.Errorf("unknown layer: %s", layer)
	}
	if err != nil {
		return nil, err
	}
	return toResultMap(result)
}

// adminUsersSchema содержит полную схему параметров для каждого layer+action.
// Возвращается LLM при первом вызове без обязательных полей (schema mode).
var adminUsersSchema = map[string]map[string]any{
	"users.list": {
		"schema":      true,
		"tool":        "admin_users",
		"layer":       "users",
		"action":      "list",
		"description": "Получить список пользователей amoCRM с пагинацией и фильтрацией.",
		"required_fields": map[string]any{
			"layer":  map[string]any{"type": "string", "value": "users"},
			"action": map[string]any{"type": "string", "value": "list"},
		},
		"optional_fields": map[string]any{
			"filter": map[string]any{
				"type":        "object",
				"description": "Фильтры для поиска",
				"properties": map[string]any{
					"limit": map[string]any{"type": "integer", "description": "Лимит результатов на странице (по умолчанию 50)"},
					"page":  map[string]any{"type": "integer", "description": "Номер страницы (начиная с 1)"},
					"name":  map[string]any{"type": "string", "description": "Фильтр по имени пользователя (client-side)"},
					"email": map[string]any{"type": "string", "description": "Фильтр по email пользователя (client-side)"},
					"order": map[string]any{"type": "object", "description": "Сортировка: {\"created_at\": \"desc\"}"},
				},
			},
		},
		"example": map[string]any{
			"layer":  "users",
			"action": "list",
			"filter": map[string]any{"limit": 20, "page": 1},
		},
	},
	"users.search": {
		"schema":      true,
		"tool":        "admin_users",
		"layer":       "users",
		"action":      "search",
		"description": "Поиск пользователей amoCRM (алиас для list с фильтрацией).",
		"required_fields": map[string]any{
			"layer":  map[string]any{"type": "string", "value": "users"},
			"action": map[string]any{"type": "string", "value": "search"},
		},
		"optional_fields": map[string]any{
			"filter": map[string]any{
				"type":        "object",
				"description": "Фильтры для поиска",
				"properties": map[string]any{
					"limit": map[string]any{"type": "integer", "description": "Лимит результатов"},
					"page":  map[string]any{"type": "integer", "description": "Номер страницы"},
					"name":  map[string]any{"type": "string", "description": "Фильтр по имени"},
					"email": map[string]any{"type": "string", "description": "Фильтр по email"},
				},
			},
		},
		"example": map[string]any{
			"layer":  "users",
			"action": "search",
			"filter": map[string]any{"name": "Иван"},
		},
	},
	"users.get": {
		"schema":      true,
		"tool":        "admin_users",
		"layer":       "users",
		"action":      "get",
		"description": "Получить пользователя по ID. Включает роль, группу, телефон и другие детали.",
		"required_fields": map[string]any{
			"layer":  map[string]any{"type": "string", "value": "users"},
			"action": map[string]any{"type": "string", "value": "get"},
			"id":     map[string]any{"type": "integer", "description": "ID пользователя"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":  "users",
			"action": "get",
			"id":     12345,
		},
	},
	"users.create": {
		"schema":      true,
		"tool":        "admin_users",
		"layer":       "users",
		"action":      "create",
		"description": "Создать одного или нескольких пользователей amoCRM.",
		"required_fields": map[string]any{
			"layer":  map[string]any{"type": "string", "value": "users"},
			"action": map[string]any{"type": "string", "value": "create"},
			"users": map[string]any{
				"type":        "array",
				"description": "Список пользователей для создания",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"name":     map[string]any{"type": "string", "description": "Имя пользователя (обязательное)"},
						"email":    map[string]any{"type": "string", "description": "Email пользователя (обязательное)"},
						"password": map[string]any{"type": "string", "description": "Пароль пользователя"},
						"lang":     map[string]any{"type": "string", "description": "Язык интерфейса: ru, en, es"},
					},
					"required": []string{"name", "email", "password"},
				},
			},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":  "users",
			"action": "create",
			"users": []any{
				map[string]any{"name": "Иван Петров", "email": "ivan@example.com", "password": "SecurePass123", "lang": "ru"},
			},
		},
		"notes": "add_to_group не реализован (API amoCRM не поддерживает назначение пользователей в группы через REST API v4). update и delete не поддерживаются amoCRM API.",
	},
	"roles.list": {
		"schema":      true,
		"tool":        "admin_users",
		"layer":       "roles",
		"action":      "list",
		"description": "Получить список ролей amoCRM с пагинацией и фильтрацией.",
		"required_fields": map[string]any{
			"layer":  map[string]any{"type": "string", "value": "roles"},
			"action": map[string]any{"type": "string", "value": "list"},
		},
		"optional_fields": map[string]any{
			"filter": map[string]any{
				"type":        "object",
				"description": "Фильтры для поиска ролей",
				"properties": map[string]any{
					"limit": map[string]any{"type": "integer", "description": "Лимит результатов (по умолчанию 50)"},
					"page":  map[string]any{"type": "integer", "description": "Номер страницы (начиная с 1)"},
					"name":  map[string]any{"type": "string", "description": "Фильтр по названию роли (client-side)"},
				},
			},
		},
		"example": map[string]any{
			"layer":  "roles",
			"action": "list",
		},
	},
	"roles.search": {
		"schema":      true,
		"tool":        "admin_users",
		"layer":       "roles",
		"action":      "search",
		"description": "Поиск ролей amoCRM (алиас для list с фильтрацией).",
		"required_fields": map[string]any{
			"layer":  map[string]any{"type": "string", "value": "roles"},
			"action": map[string]any{"type": "string", "value": "search"},
		},
		"optional_fields": map[string]any{
			"filter": map[string]any{
				"type":        "object",
				"description": "Фильтры для поиска",
				"properties": map[string]any{
					"name": map[string]any{"type": "string", "description": "Фильтр по названию роли"},
				},
			},
		},
		"example": map[string]any{
			"layer":  "roles",
			"action": "search",
			"filter": map[string]any{"name": "Менеджер"},
		},
	},
	"roles.get": {
		"schema":      true,
		"tool":        "admin_users",
		"layer":       "roles",
		"action":      "get",
		"description": "Получить роль по ID. Включает список пользователей с этой ролью.",
		"required_fields": map[string]any{
			"layer":  map[string]any{"type": "string", "value": "roles"},
			"action": map[string]any{"type": "string", "value": "get"},
			"id":     map[string]any{"type": "integer", "description": "ID роли"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":  "roles",
			"action": "get",
			"id":     42,
		},
	},
	"roles.create": {
		"schema":      true,
		"tool":        "admin_users",
		"layer":       "roles",
		"action":      "create",
		"description": "Создать одну или несколько ролей amoCRM.",
		"required_fields": map[string]any{
			"layer":  map[string]any{"type": "string", "value": "roles"},
			"action": map[string]any{"type": "string", "value": "create"},
			"roles": map[string]any{
				"type":        "array",
				"description": "Список ролей для создания",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"name": map[string]any{"type": "string", "description": "Название роли (обязательное)"},
					},
					"required": []string{"name"},
				},
			},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":  "roles",
			"action": "create",
			"roles":  []any{map[string]any{"name": "Старший менеджер"}},
		},
	},
	"roles.update": {
		"schema":      true,
		"tool":        "admin_users",
		"layer":       "roles",
		"action":      "update",
		"description": "Обновить одну или несколько ролей amoCRM. Требует ID каждой роли.",
		"required_fields": map[string]any{
			"layer":  map[string]any{"type": "string", "value": "roles"},
			"action": map[string]any{"type": "string", "value": "update"},
			"roles": map[string]any{
				"type":        "array",
				"description": "Список ролей для обновления",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id":   map[string]any{"type": "integer", "description": "ID роли (обязательное для обновления)"},
						"name": map[string]any{"type": "string", "description": "Новое название роли"},
					},
					"required": []string{"id", "name"},
				},
			},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":  "roles",
			"action": "update",
			"roles":  []any{map[string]any{"id": 42, "name": "Ведущий менеджер"}},
		},
	},
	"roles.delete": {
		"schema":      true,
		"tool":        "admin_users",
		"layer":       "roles",
		"action":      "delete",
		"description": "Удалить роль по ID.",
		"required_fields": map[string]any{
			"layer":  map[string]any{"type": "string", "value": "roles"},
			"action": map[string]any{"type": "string", "value": "delete"},
			"id":     map[string]any{"type": "integer", "description": "ID роли для удаления"},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"layer":  "roles",
			"action": "delete",
			"id":     42,
		},
	},
}

// adminUsersRequiredFields перечисляет обязательные поля для каждого layer.action (помимо layer и action).
// Если хотя бы одно из них отсутствует — режим Schema.
var adminUsersRequiredFields = map[string][]string{
	"users.list":   {},
	"users.search": {},
	"users.get":    {"id"},
	"users.create": {"users"},
	"roles.list":   {},
	"roles.search": {},
	"roles.get":    {"id"},
	"roles.create": {"roles"},
	"roles.update": {"roles"},
	"roles.delete": {"id"},
}

// adminUsersMapToInput конвертирует map[string]any → AdminUsersInput через json roundtrip.
func adminUsersMapToInput(raw map[string]any) (gkitmodels.AdminUsersInput, error) {
	data, err := json.Marshal(raw)
	if err != nil {
		return gkitmodels.AdminUsersInput{}, fmt.Errorf("marshal: %w", err)
	}
	var input gkitmodels.AdminUsersInput
	if err := json.Unmarshal(data, &input); err != nil {
		return gkitmodels.AdminUsersInput{}, fmt.Errorf("unmarshal: %w", err)
	}
	return input, nil
}

func (t *AdminUsersTool) handleUsers(ctx tool.Context, input gkitmodels.AdminUsersInput) (any, error) {
	switch input.Action {
	case "list", "search":
		return t.service.ListUsers(ctx, input.Filter)

	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for get user")
		}
		return t.service.GetUser(ctx, input.ID)

	case "create":
		if len(input.Users) == 0 {
			return nil, fmt.Errorf("users list is empty: provide at least one user in 'users' field")
		}
		sdkUsers := make([]*amomodels.User, 0, len(input.Users))
		for _, u := range input.Users {
			sdkUsers = append(sdkUsers, &amomodels.User{
				Name:     u.Name,
				Email:    u.Email,
				Password: u.Password,
				Lang:     u.Lang,
			})
		}
		return t.service.CreateUsers(ctx, sdkUsers)

	case "add_to_group":
		return nil, fmt.Errorf("add_to_group: not implemented — amoCRM API does not support assigning users to groups via REST API v4")

	case "update", "delete":
		return nil, fmt.Errorf("action %s is not supported for users by amoCRM API", input.Action)

	default:
		return nil, fmt.Errorf("unknown action for users: %s", input.Action)
	}
}

func (t *AdminUsersTool) handleRoles(ctx tool.Context, input gkitmodels.AdminUsersInput) (any, error) {
	switch input.Action {
	case "list", "search":
		return t.service.ListRoles(ctx, input.Filter)

	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for get role")
		}
		return t.service.GetRole(ctx, input.ID)

	case "create":
		if len(input.Roles) == 0 {
			return nil, fmt.Errorf("roles list is empty: provide at least one role in 'roles' field")
		}
		sdkRoles := make([]*amomodels.Role, 0, len(input.Roles))
		for _, ro := range input.Roles {
			sdkRoles = append(sdkRoles, &amomodels.Role{
				Name: ro.Name,
			})
		}
		return t.service.CreateRoles(ctx, sdkRoles)

	case "update":
		if len(input.Roles) == 0 {
			return nil, fmt.Errorf("roles list is empty: provide at least one role in 'roles' field")
		}
		sdkRoles := make([]*amomodels.Role, 0, len(input.Roles))
		for _, ro := range input.Roles {
			if ro.ID == 0 {
				return nil, fmt.Errorf("role id is required for update (role name: %q)", ro.Name)
			}
			sdkRoles = append(sdkRoles, &amomodels.Role{
				ID:   ro.ID,
				Name: ro.Name,
			})
		}
		return t.service.UpdateRoles(ctx, sdkRoles)

	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for delete role")
		}
		result, err := t.service.DeleteRole(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, fmt.Errorf("unknown action for roles: %s", input.Action)
	}
}
