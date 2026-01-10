package tools

import (
	"encoding/json"
	"fmt"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func (r *Registry) RegisterAdminPipelinesTool() {
	r.addTool(genkit.DefineTool[gkitmodels.AdminPipelinesInput, any](
		r.g,
		"admin_pipelines",
		"Work with pipelines and statuses",
		func(ctx *ai.ToolContext, input gkitmodels.AdminPipelinesInput) (any, error) {
			switch input.Action {
			// Pipelines
			case "list", "search":
				return r.adminPipelinesService.ListPipelines(ctx)
			case "get":
				if input.PipelineID == 0 {
					return nil, fmt.Errorf("pipeline_id is required")
				}
				return r.adminPipelinesService.GetPipeline(ctx, input.PipelineID)
			case "create":
				var pipelines []*amomodels.Pipeline
				data, _ := json.Marshal(input.Data["pipelines"])
				if err := json.Unmarshal(data, &pipelines); err != nil {
					return nil, fmt.Errorf("failed to parse pipelines: %w", err)
				}
				return r.adminPipelinesService.CreatePipelines(ctx, pipelines)
			case "update":
				var p amomodels.Pipeline
				data, _ := json.Marshal(input.Data)
				if err := json.Unmarshal(data, &p); err != nil {
					return nil, fmt.Errorf("failed to parse pipeline data: %w", err)
				}
				if p.ID == 0 {
					p.ID = input.PipelineID
				}
				if p.ID == 0 {
					return nil, fmt.Errorf("pipeline id is required for update")
				}
				return r.adminPipelinesService.UpdatePipeline(ctx, &p)
			case "delete":
				if input.PipelineID == 0 {
					return nil, fmt.Errorf("pipeline_id is required")
				}
				return nil, r.adminPipelinesService.DeletePipeline(ctx, input.PipelineID)

			// Statuses
			case "list_statuses", "get_statuses":
				if input.PipelineID == 0 {
					return nil, fmt.Errorf("pipeline_id is required")
				}
				return r.adminPipelinesService.ListStatuses(ctx, input.PipelineID)
			case "get_status":
				if input.PipelineID == 0 || input.StatusID == 0 {
					return nil, fmt.Errorf("pipeline_id and status_id are required")
				}
				return r.adminPipelinesService.GetStatus(ctx, input.PipelineID, input.StatusID)
			case "create_status":
				if input.PipelineID == 0 {
					return nil, fmt.Errorf("pipeline_id is required")
				}
				var s amomodels.Status
				data, _ := json.Marshal(input.Data)
				if err := json.Unmarshal(data, &s); err != nil {
					return nil, fmt.Errorf("failed to parse status data: %w", err)
				}
				return r.adminPipelinesService.CreateStatus(ctx, input.PipelineID, &s)
			case "update_status":
				if input.PipelineID == 0 {
					return nil, fmt.Errorf("pipeline_id is required")
				}
				var s amomodels.Status
				data, _ := json.Marshal(input.Data)
				if err := json.Unmarshal(data, &s); err != nil {
					return nil, fmt.Errorf("failed to parse status data: %w", err)
				}
				if s.ID == 0 {
					s.ID = input.StatusID
				}
				if s.ID == 0 {
					return nil, fmt.Errorf("status id is required for update")
				}
				return r.adminPipelinesService.UpdateStatus(ctx, input.PipelineID, &s)
			case "delete_status":
				if input.PipelineID == 0 || input.StatusID == 0 {
					return nil, fmt.Errorf("pipeline_id and status_id are required")
				}
				return nil, r.adminPipelinesService.DeleteStatus(ctx, input.PipelineID, input.StatusID)

			default:
				return nil, fmt.Errorf("unknown action: %s", input.Action)
			}
		},
	))
}
