package unsorted

import (
	"context"
	"fmt"
	"strings"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) ListUnsorted(ctx context.Context, filter *gkitmodels.UnsortedFilter) (*UnsortedListOutput, error) {
	f := filters.NewUnsortedFilter()
	f.SetWith("leads", "contacts", "companies")

	if filter != nil {
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		}
		if filter.Limit > 0 {
			f.SetLimit(filter.Limit)
		}
		if len(filter.Category) > 0 {
			f.SetCategory(filter.Category)
		}
		if filter.PipelineName != "" {
			id, err := s.resolvePipelineName(filter.PipelineName)
			if err != nil {
				return nil, err
			}
			if id != 0 {
				f.SetPipelineID(id)
			}
		}
		if filter.Order != "" {
			parts := strings.SplitN(filter.Order, " ", 2)
			if len(parts) == 2 {
				f.SetOrder(parts[0], parts[1])
			}
		}
	}

	items, meta, err := s.sdk.Unsorted().Get(ctx, f)
	if err != nil {
		return nil, err
	}

	out := make([]*UnsortedOutput, 0, len(items))
	for _, item := range items {
		out = append(out, s.unsortedToOutput(item))
	}
	return &UnsortedListOutput{Items: out, PageMeta: meta}, nil
}

func (s *service) GetUnsorted(ctx context.Context, uid string) (*UnsortedOutput, error) {
	f := filters.NewUnsortedFilter().SetUIDs([]string{uid})
	f.SetWith("leads", "contacts", "companies")

	items, _, err := s.sdk.Unsorted().Get(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("unsorted record %s not found", uid)
	}
	return s.unsortedToOutput(items[0]), nil
}

func (s *service) CreateUnsorted(ctx context.Context, category string, items []gkitmodels.UnsortedCreateItem) ([]*UnsortedOutput, error) {
	sdkItems := make([]*models.Unsorted, 0, len(items))
	for _, item := range items {
		u := &models.Unsorted{
			SourceUID:  item.SourceUID,
			SourceName: item.SourceName,
			CreatedAt:  rfc3339ToUnix(item.CreatedAt),
		}
		if item.PipelineName != "" {
			id, err := s.resolvePipelineName(item.PipelineName)
			if err != nil {
				return nil, err
			}
			u.PipelineID = id
		}
		sdkItems = append(sdkItems, u)
	}

	result, _, err := s.sdk.Unsorted().Create(ctx, category, sdkItems)
	if err != nil {
		return nil, err
	}

	out := make([]*UnsortedOutput, 0, len(result))
	for _, item := range result {
		out = append(out, s.unsortedToOutput(item))
	}
	return out, nil
}

func (s *service) AcceptUnsorted(ctx context.Context, uid string, params *gkitmodels.UnsortedAcceptParams) (*UnsortedActionResult, error) {
	apiParams := map[string]interface{}{}

	if params != nil {
		if params.UserName != "" {
			userID, err := s.resolveUserName(params.UserName)
			if err != nil {
				return nil, err
			}
			if userID > 0 {
				apiParams["user_id"] = userID
			}
		}
		if params.StatusName != "" {
			pipelineID, err := s.resolvePipelineName(params.PipelineName)
			if err != nil {
				return nil, err
			}
			statusID, err := s.resolveStatusName(pipelineID, params.StatusName)
			if err != nil {
				return nil, err
			}
			if statusID > 0 {
				apiParams["status_id"] = statusID
			}
		}
	}

	result, err := s.sdk.Unsorted().Accept(ctx, uid, apiParams)
	if err != nil {
		return nil, err
	}
	return &UnsortedActionResult{UID: result.UID, Success: result.Result}, nil
}

func (s *service) DeclineUnsorted(ctx context.Context, uid string, params *gkitmodels.UnsortedDeclineParams) (*UnsortedActionResult, error) {
	apiParams := map[string]interface{}{}

	if params != nil && params.UserName != "" {
		userID, err := s.resolveUserName(params.UserName)
		if err != nil {
			return nil, err
		}
		if userID > 0 {
			apiParams["user_id"] = userID
		}
	}

	result, err := s.sdk.Unsorted().Decline(ctx, uid, apiParams)
	if err != nil {
		return nil, err
	}
	return &UnsortedActionResult{UID: result.UID, Success: result.Result}, nil
}

func (s *service) LinkUnsorted(ctx context.Context, uid string, leadID int) (*UnsortedActionResult, error) {
	params := map[string]interface{}{
		"lead_id": leadID,
	}
	result, err := s.sdk.Unsorted().Link(ctx, uid, params)
	if err != nil {
		return nil, err
	}
	return &UnsortedActionResult{UID: result.UID, Success: result.Result}, nil
}

func (s *service) SummaryUnsorted(ctx context.Context, filter *gkitmodels.UnsortedFilter) (*models.UnsortedSummary, error) {
	f := filters.NewUnsortedSummaryFilter()
	if filter != nil {
		if filter.PipelineName != "" {
			id, err := s.resolvePipelineName(filter.PipelineName)
			if err != nil {
				return nil, err
			}
			if id != 0 {
				f.SetPipelineID(id)
			}
		}
		if filter.CreatedAtFrom != "" || filter.CreatedAtTo != "" {
			var fromPtr, toPtr *int
			if filter.CreatedAtFrom != "" {
				v := int(rfc3339ToUnix(filter.CreatedAtFrom))
				fromPtr = &v
			}
			if filter.CreatedAtTo != "" {
				v := int(rfc3339ToUnix(filter.CreatedAtTo))
				toPtr = &v
			}
			f.SetCreatedAt(fromPtr, toPtr)
		}
	}
	return s.sdk.Unsorted().Summary(ctx, f)
}
