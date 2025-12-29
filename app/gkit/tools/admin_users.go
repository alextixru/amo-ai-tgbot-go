package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// AdminUsersInput входные параметры для инструмента admin_users
type AdminUsersInput struct {
	// Layer слой: users | roles
	Layer string `json:"layer" jsonschema_description:"Слой: users (пользователи), roles (роли)"`

	// Action действие
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update (только roles), delete (только roles)"`

	// ID идентификатор пользователя или роли
	ID int `json:"id,omitempty" jsonschema_description:"ID пользователя или роли"`

	// Filter фильтры для list
	Filter *AdminUsersFilter `json:"filter,omitempty" jsonschema_description:"Фильтры для поиска"`

	// Data данные для create/update
	Data map[string]any `json:"data,omitempty" jsonschema_description:"Данные для создания/обновления"`
}

// AdminUsersFilter фильтры для admin_users
type AdminUsersFilter struct {
	Limit int `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Page  int `json:"page,omitempty" jsonschema_description:"Номер страницы"`
}

// registerAdminUsersTool регистрирует инструмент для управления пользователями и ролями
func (r *Registry) registerAdminUsersTool() {
	r.addTool(genkit.DefineTool[AdminUsersInput, any](
		r.g,
		"admin_users",
		"Управление пользователями и ролями amoCRM. "+
			"Layers: users (пользователи), roles (роли). "+
			"Users actions: list, get, create. Roles actions: list, get, create, update, delete. "+
			"Note: Users API ограничен — update/delete пользователей недоступны.",
		func(ctx *ai.ToolContext, input AdminUsersInput) (any, error) {
			return r.handleAdminUsers(ctx.Context, input)
		},
	))
}

func (r *Registry) handleAdminUsers(ctx context.Context, input AdminUsersInput) (any, error) {
	switch input.Layer {
	case "users":
		return r.handleUsers(ctx, input)
	case "roles":
		return r.handleRoles(ctx, input)
	default:
		return nil, fmt.Errorf("unknown layer: %s (expected: users, roles)", input.Layer)
	}
}

// ============================================================================
// Users
// ============================================================================

func (r *Registry) handleUsers(ctx context.Context, input AdminUsersInput) (any, error) {
	switch input.Action {
	case "list":
		return r.listUsers(ctx, input.Filter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.sdk.Users().GetOne(ctx, input.ID)
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		return r.createUser(ctx, input.Data)
	default:
		return nil, fmt.Errorf("unknown action: %s for users (available: list, get, create)", input.Action)
	}
}

func (r *Registry) listUsers(ctx context.Context, filter *AdminUsersFilter) ([]models.User, error) {
	f := &services.UsersFilter{
		Limit: 50,
		Page:  1,
	}
	if filter != nil {
		if filter.Limit > 0 {
			f.Limit = filter.Limit
		}
		if filter.Page > 0 {
			f.Page = filter.Page
		}
	}
	return r.sdk.Users().Get(ctx, f)
}

func (r *Registry) createUser(ctx context.Context, data map[string]any) ([]models.User, error) {
	user := models.User{}

	if name, ok := data["name"].(string); ok {
		user.Name = name
	}
	if email, ok := data["email"].(string); ok {
		user.Email = email
	}

	return r.sdk.Users().Create(ctx, []models.User{user})
}

// ============================================================================
// Roles
// ============================================================================

func (r *Registry) handleRoles(ctx context.Context, input AdminUsersInput) (any, error) {
	switch input.Action {
	case "list":
		return r.listRoles(ctx, input.Filter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.sdk.Roles().GetOne(ctx, input.ID, nil)
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		return r.createRole(ctx, input.Data)
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		return r.updateRole(ctx, input.ID, input.Data)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'delete'")
		}
		err := r.sdk.Roles().Delete(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true, "deleted_role_id": input.ID}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s for roles", input.Action)
	}
}

func (r *Registry) listRoles(ctx context.Context, filter *AdminUsersFilter) ([]models.Role, error) {
	f := &services.RoleFilter{
		Limit: 50,
		Page:  1,
	}
	if filter != nil {
		if filter.Limit > 0 {
			f.Limit = filter.Limit
		}
		if filter.Page > 0 {
			f.Page = filter.Page
		}
	}
	return r.sdk.Roles().Get(ctx, f)
}

func (r *Registry) createRole(ctx context.Context, data map[string]any) ([]models.Role, error) {
	role := models.Role{}

	if name, ok := data["name"].(string); ok {
		role.Name = name
	}

	return r.sdk.Roles().Create(ctx, []models.Role{role})
}

func (r *Registry) updateRole(ctx context.Context, id int, data map[string]any) ([]models.Role, error) {
	role := models.Role{ID: id}

	if name, ok := data["name"].(string); ok {
		role.Name = name
	}

	return r.sdk.Roles().Update(ctx, []models.Role{role})
}
