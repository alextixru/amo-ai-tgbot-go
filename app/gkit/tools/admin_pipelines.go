package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// AdminPipelinesInput входные параметры для инструмента admin_pipelines
type AdminPipelinesInput struct {
	// Action действие: list, get, create, update, delete, list_statuses, get_status, create_status, update_status, delete_status
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, delete, list_statuses, get_status, create_status, update_status, delete_status"`

	// PipelineID идентификатор воронки
	PipelineID int `json:"pipeline_id,omitempty" jsonschema_description:"ID воронки (для большинства операций)"`

	// StatusID идентификатор статуса
	StatusID int `json:"status_id,omitempty" jsonschema_description:"ID статуса (для операций со статусами)"`

	// Data данные для create/update
	Data map[string]any `json:"data,omitempty" jsonschema_description:"Данные для создания/обновления"`
}

// registerAdminPipelinesTool регистрирует инструмент для управления воронками
func (r *Registry) registerAdminPipelinesTool() {
	r.addTool(genkit.DefineTool[AdminPipelinesInput, any](
		r.g,
		"admin_pipelines",
		"Управление воронками и статусами amoCRM. "+
			"Pipeline actions: list (все воронки), get (по ID), create, update, delete. "+
			"Status actions: list_statuses, get_status, create_status, update_status, delete_status. "+
			"Для статусов требуется pipeline_id. Для update_status и delete_status также status_id.",
		func(ctx *ai.ToolContext, input AdminPipelinesInput) (any, error) {
			return r.handleAdminPipelines(ctx.Context, input)
		},
	))
}

func (r *Registry) handleAdminPipelines(ctx context.Context, input AdminPipelinesInput) (any, error) {
	switch input.Action {
	// Pipeline actions
	case "list":
		pipelines, _, err := r.sdk.Pipelines().Get(ctx, nil)
		return pipelines, err
	case "get":
		if input.PipelineID == 0 {
			return nil, fmt.Errorf("pipeline_id is required for action 'get'")
		}
		return r.sdk.Pipelines().GetOne(ctx, input.PipelineID)
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		return r.createPipeline(ctx, input.Data)
	case "update":
		if input.PipelineID == 0 {
			return nil, fmt.Errorf("pipeline_id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		return r.updatePipeline(ctx, input.PipelineID, input.Data)
	case "delete":
		if input.PipelineID == 0 {
			return nil, fmt.Errorf("pipeline_id is required for action 'delete'")
		}
		err := r.sdk.Pipelines().Delete(ctx, []int{input.PipelineID})
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true, "deleted_pipeline_id": input.PipelineID}, nil

	// Status actions (using StatusesService)
	case "list_statuses", "get_statuses":
		if input.PipelineID == 0 {
			return nil, fmt.Errorf("pipeline_id is required for action 'list_statuses'")
		}
		statuses, _, err := r.sdk.Statuses(input.PipelineID).Get(ctx, nil)
		return statuses, err
	case "get_status":
		if input.PipelineID == 0 {
			return nil, fmt.Errorf("pipeline_id is required for action 'get_status'")
		}
		if input.StatusID == 0 {
			return nil, fmt.Errorf("status_id is required for action 'get_status'")
		}
		return r.sdk.Statuses(input.PipelineID).GetOne(ctx, input.StatusID, nil)
	case "create_status":
		if input.PipelineID == 0 {
			return nil, fmt.Errorf("pipeline_id is required for action 'create_status'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create_status'")
		}
		return r.createStatus(ctx, input.PipelineID, input.Data)
	case "update_status":
		if input.PipelineID == 0 {
			return nil, fmt.Errorf("pipeline_id is required for action 'update_status'")
		}
		if input.StatusID == 0 {
			return nil, fmt.Errorf("status_id is required for action 'update_status'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update_status'")
		}
		return r.updateStatus(ctx, input.PipelineID, input.StatusID, input.Data)
	case "delete_status":
		if input.PipelineID == 0 {
			return nil, fmt.Errorf("pipeline_id is required for action 'delete_status'")
		}
		if input.StatusID == 0 {
			return nil, fmt.Errorf("status_id is required for action 'delete_status'")
		}
		err := r.sdk.Statuses(input.PipelineID).DeleteOne(ctx, input.StatusID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true, "deleted_status_id": input.StatusID, "pipeline_id": input.PipelineID}, nil

	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

// ============================================================================
// Pipeline helpers
// ============================================================================

func (r *Registry) createPipeline(ctx context.Context, data map[string]any) ([]*models.Pipeline, error) {
	pipeline := models.Pipeline{}

	if name, ok := data["name"].(string); ok {
		pipeline.Name = name
	}
	if sort, ok := data["sort"].(float64); ok {
		pipeline.Sort = int(sort)
	}
	if isMain, ok := data["is_main"].(bool); ok {
		pipeline.IsMain = isMain
	}
	if isUnsortedOn, ok := data["is_unsorted_on"].(bool); ok {
		pipeline.IsUnsortedOn = isUnsortedOn
	}

	pipelines, _, err := r.sdk.Pipelines().Create(ctx, []*models.Pipeline{&pipeline})
	return pipelines, err
}

func (r *Registry) updatePipeline(ctx context.Context, id int, data map[string]any) ([]*models.Pipeline, error) {
	pipeline := models.Pipeline{ID: id}

	if name, ok := data["name"].(string); ok {
		pipeline.Name = name
	}
	if sort, ok := data["sort"].(float64); ok {
		pipeline.Sort = int(sort)
	}
	if isMain, ok := data["is_main"].(bool); ok {
		pipeline.IsMain = isMain
	}
	if isUnsortedOn, ok := data["is_unsorted_on"].(bool); ok {
		pipeline.IsUnsortedOn = isUnsortedOn
	}

	return r.sdk.Pipelines().Update(ctx, []*models.Pipeline{&pipeline})
}

// ============================================================================
// Status helpers
// ============================================================================

func (r *Registry) createStatus(ctx context.Context, pipelineID int, data map[string]any) (*models.Status, error) {
	status := &models.Status{}

	if name, ok := data["name"].(string); ok {
		status.Name = name
	}
	if sort, ok := data["sort"].(float64); ok {
		status.Sort = int(sort)
	}
	if color, ok := data["color"].(string); ok {
		status.Color = color
	}

	statuses, _, err := r.sdk.Statuses(pipelineID).Create(ctx, []*models.Status{status})
	if err != nil {
		return nil, err
	}
	if len(statuses) == 0 {
		return nil, fmt.Errorf("no status returned from create")
	}
	return statuses[0], nil
}

func (r *Registry) updateStatus(ctx context.Context, pipelineID int, statusID int, data map[string]any) (*models.Status, error) {
	status := &models.Status{ID: statusID}

	if name, ok := data["name"].(string); ok {
		status.Name = name
	}
	if sort, ok := data["sort"].(float64); ok {
		status.Sort = int(sort)
	}
	if color, ok := data["color"].(string); ok {
		status.Color = color
	}

	return r.sdk.Statuses(pipelineID).UpdateOne(ctx, status)
}
