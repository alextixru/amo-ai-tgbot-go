package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (r *Registry) RegisterUnsortedTool() {
	r.addTool(genkit.DefineTool[gkitmodels.UnsortedInput, any](
		r.g,
		"unsorted",
		"Работа с входящими заявками (Неразобранное). "+
			"Поддерживает: list (список), get (получение по UID), accept (принять), decline (отклонить), "+
			"link (привязать к существующей сделке), summary (статистика), create (создать заявку).",
		func(ctx *ai.ToolContext, input gkitmodels.UnsortedInput) (any, error) {
			return r.handleUnsorted(ctx.Context, input)
		},
	))
}

func (r *Registry) handleUnsorted(ctx context.Context, input gkitmodels.UnsortedInput) (any, error) {
	switch input.Action {
	case "search", "list":
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
		var items []*models.Unsorted
		for _, item := range input.CreateData.Items {
			u := &models.Unsorted{
				SourceUID:  item.SourceUID,
				SourceName: item.SourceName,
				PipelineID: item.PipelineID,
				CreatedAt:  int64(item.CreatedAt),
			}
			items = append(items, u)
		}
		return r.unsortedService.CreateUnsorted(ctx, input.CreateData.Category, items)
	case "accept":
		if input.UID == "" {
			return nil, fmt.Errorf("uid is required for action 'accept'")
		}
		userID := 0
		statusID := 0
		if input.AcceptParams != nil {
			userID = input.AcceptParams.UserID
			statusID = input.AcceptParams.StatusID
		}
		return r.unsortedService.AcceptUnsorted(ctx, input.UID, userID, statusID)
	case "decline":
		if input.UID == "" {
			return nil, fmt.Errorf("uid is required for action 'decline'")
		}
		userID := 0
		if input.AcceptParams != nil {
			userID = input.AcceptParams.UserID
		}
		return r.unsortedService.DeclineUnsorted(ctx, input.UID, userID)
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
