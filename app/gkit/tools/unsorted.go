package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// UnsortedInput входные параметры для инструмента unsorted
type UnsortedInput struct {
	// Action действие: list, get, accept, decline, link, summary
	Action string `json:"action" jsonschema_description:"Действие: list, get, accept, decline, link, summary"`

	// UID идентификатор неразобранного (для get, accept, decline, link)
	UID string `json:"uid,omitempty" jsonschema_description:"UID записи неразобранного"`

	// Filter параметры поиска (для list, summary)
	Filter *UnsortedFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска"`

	// AcceptParams параметры принятия (для accept)
	AcceptParams *UnsortedAcceptParams `json:"accept_params,omitempty" jsonschema_description:"Параметры принятия заявки"`

	// LinkData данные привязки (для link)
	LinkData *UnsortedLinkData `json:"link_data,omitempty" jsonschema_description:"Данные для привязки к сделке"`
}

// UnsortedFilter фильтры поиска неразобранного
type UnsortedFilter struct {
	Page       int      `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Limit      int      `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Category   []string `json:"category,omitempty" jsonschema_description:"Категории: sip, mail, forms, chats"`
	PipelineID []int    `json:"pipeline_id,omitempty" jsonschema_description:"ID воронок"`
}

// UnsortedAcceptParams параметры принятия неразобранного
type UnsortedAcceptParams struct {
	UserID   int `json:"user_id,omitempty" jsonschema_description:"ID ответственного пользователя"`
	StatusID int `json:"status_id,omitempty" jsonschema_description:"ID статуса для создаваемой сделки"`
}

// UnsortedLinkData данные для привязки неразобранного
type UnsortedLinkData struct {
	LeadID int `json:"lead_id" jsonschema_description:"ID существующей сделки для привязки"`
}

// registerUnsortedTool регистрирует инструмент для работы с неразобранным
func (r *Registry) registerUnsortedTool() {
	r.addTool(genkit.DefineTool[UnsortedInput, any](
		r.g,
		"unsorted",
		"Работа с неразобранным (входящие заявки). "+
			"Actions: list (список), get (по UID), accept (принять → создать сделку), "+
			"decline (отклонить), link (привязать к сделке), summary (статистика). "+
			"Категории: sip (звонки), mail (почта), forms (формы), chats (чаты).",
		func(ctx *ai.ToolContext, input UnsortedInput) (any, error) {
			return r.handleUnsorted(ctx.Context, input)
		},
	))
}

func (r *Registry) handleUnsorted(ctx context.Context, input UnsortedInput) (any, error) {
	switch input.Action {
	case "list", "search":
		return r.listUnsorted(ctx, input.Filter)
	case "get":
		if input.UID == "" {
			return nil, fmt.Errorf("uid is required for action 'get'")
		}
		items, _, err := r.sdk.Unsorted().Get(ctx, filters.NewUnsortedFilter().SetUIDs([]string{input.UID}))
		if err != nil {
			return nil, err
		}
		if len(items) == 0 {
			return nil, fmt.Errorf("unsorted item with uid %s not found", input.UID)
		}
		return items[0], nil
	case "accept":
		if input.UID == "" {
			return nil, fmt.Errorf("uid is required for action 'accept'")
		}
		return r.acceptUnsorted(ctx, input.UID, input.AcceptParams)
	case "decline":
		if input.UID == "" {
			return nil, fmt.Errorf("uid is required for action 'decline'")
		}
		return r.sdk.Unsorted().Decline(ctx, input.UID, nil)
	case "link":
		if input.UID == "" {
			return nil, fmt.Errorf("uid is required for action 'link'")
		}
		if input.LinkData == nil || input.LinkData.LeadID == 0 {
			return nil, fmt.Errorf("link_data.lead_id is required for action 'link'")
		}
		return r.linkUnsorted(ctx, input.UID, input.LinkData)
	case "summary":
		return r.unsortedSummary(ctx, input.Filter)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (r *Registry) listUnsorted(ctx context.Context, filter *UnsortedFilter) ([]*models.Unsorted, error) {
	f := filters.NewUnsortedFilter()
	f.SetLimit(50)
	f.SetPage(1)
	if filter != nil {
		if filter.Limit > 0 {
			f.SetLimit(filter.Limit)
		}
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		}
		if len(filter.Category) > 0 {
			f.SetCategory(filter.Category)
		}
		if len(filter.PipelineID) > 0 {
			f.SetPipelineID(filter.PipelineID[0])
		}
	}
	items, _, err := r.sdk.Unsorted().Get(ctx, f)
	return items, err
}

func (r *Registry) acceptUnsorted(ctx context.Context, uid string, params *UnsortedAcceptParams) (*models.UnsortedAcceptResult, error) {
	acceptParams := make(map[string]interface{})
	if params != nil {
		if params.UserID > 0 {
			acceptParams["user_id"] = params.UserID
		}
		if params.StatusID > 0 {
			acceptParams["status_id"] = params.StatusID
		}
	}
	return r.sdk.Unsorted().Accept(ctx, uid, acceptParams)
}

func (r *Registry) linkUnsorted(ctx context.Context, uid string, data *UnsortedLinkData) (*models.UnsortedLinkResult, error) {
	linkData := map[string]interface{}{
		"link": map[string]interface{}{
			"entity_id":   data.LeadID,
			"entity_type": "leads",
		},
	}
	return r.sdk.Unsorted().Link(ctx, uid, linkData)
}

func (r *Registry) unsortedSummary(ctx context.Context, filter *UnsortedFilter) (*models.UnsortedSummary, error) {
	f := filters.NewUnsortedSummaryFilter()
	if filter != nil {
		if len(filter.PipelineID) > 0 {
			f.SetPipelineID(filter.PipelineID[0])
		}
	}
	return r.sdk.Unsorted().Summary(ctx, f)
}
