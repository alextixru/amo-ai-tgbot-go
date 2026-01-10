package tools

import (
	"encoding/json"
	"fmt"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func (r *Registry) RegisterAdminUsersTool() {
	r.addTool(genkit.DefineTool[gkitmodels.AdminUsersInput, any](
		r.g,
		"admin_users",
		"Work with users and roles",
		func(ctx *ai.ToolContext, input gkitmodels.AdminUsersInput) (any, error) {
			switch input.Layer {
			case "users":
				return r.handleUsers(ctx, input)
			case "roles":
				return r.handleRoles(ctx, input)
			default:
				return nil, fmt.Errorf("unknown layer: %s", input.Layer)
			}
		},
	))
}

func (r *Registry) handleUsers(ctx *ai.ToolContext, input gkitmodels.AdminUsersInput) (any, error) {
	switch input.Action {
	case "list", "search":
		return r.adminUsersService.ListUsers(ctx)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for get user")
		}
		return r.adminUsersService.GetUser(ctx, input.ID)
	case "create":
		var users []*amomodels.User
		data, _ := json.Marshal(input.Data["users"])
		if err := json.Unmarshal(data, &users); err != nil {
			return nil, fmt.Errorf("failed to parse users: %w", err)
		}
		return r.adminUsersService.CreateUsers(ctx, users)
	case "update", "delete":
		return nil, fmt.Errorf("action %s is not supported for users by amoCRM API", input.Action)
	default:
		return nil, fmt.Errorf("unknown action for users: %s", input.Action)
	}
}

func (r *Registry) handleRoles(ctx *ai.ToolContext, input gkitmodels.AdminUsersInput) (any, error) {
	switch input.Action {
	case "list", "search":
		return r.adminUsersService.ListRoles(ctx)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for get role")
		}
		return r.adminUsersService.GetRole(ctx, input.ID)
	case "create":
		var roles []*amomodels.Role
		data, _ := json.Marshal(input.Data["roles"])
		if err := json.Unmarshal(data, &roles); err != nil {
			return nil, fmt.Errorf("failed to parse roles: %w", err)
		}
		return r.adminUsersService.CreateRoles(ctx, roles)
	case "update":
		var roles []*amomodels.Role
		data, _ := json.Marshal(input.Data["roles"])
		if err := json.Unmarshal(data, &roles); err != nil {
			return nil, fmt.Errorf("failed to parse roles: %w", err)
		}
		return r.adminUsersService.UpdateRoles(ctx, roles)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for delete role")
		}
		return nil, r.adminUsersService.DeleteRole(ctx, input.ID)
	default:
		return nil, fmt.Errorf("unknown action for roles: %s", input.Action)
	}
}
