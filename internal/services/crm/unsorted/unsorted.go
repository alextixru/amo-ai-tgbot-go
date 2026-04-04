package unsorted

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) ListUnsorted(ctx context.Context, filter *gkitmodels.UnsortedFilter) ([]*models.Unsorted, error) {
	f := filters.NewUnsortedFilter()
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
		if len(filter.PipelineID) > 0 {
			f.SetPipelineID(filter.PipelineID[0])
		}
	}
	items, _, err := s.sdk.Unsorted().Get(ctx, f)
	return items, err
}

func (s *service) GetUnsorted(ctx context.Context, uid string) (*models.Unsorted, error) {
	f := filters.NewUnsortedFilter().SetUIDs([]string{uid})
	items, _, err := s.sdk.Unsorted().Get(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("unsorted record %s not found", uid)
	}
	return items[0], nil
}

func (s *service) CreateUnsorted(ctx context.Context, category string, items []*models.Unsorted) ([]*models.Unsorted, error) {
	result, _, err := s.sdk.Unsorted().Create(ctx, category, items)
	return result, err
}

func (s *service) AcceptUnsorted(ctx context.Context, uid string, userID int, statusID int) (*models.UnsortedAcceptResult, error) {
	params := map[string]interface{}{}
	if userID > 0 {
		params["user_id"] = userID
	}
	if statusID > 0 {
		params["status_id"] = statusID
	}
	return s.sdk.Unsorted().Accept(ctx, uid, params)
}

func (s *service) DeclineUnsorted(ctx context.Context, uid string, userID int) (*models.UnsortedDeclineResult, error) {
	params := map[string]interface{}{}
	if userID > 0 {
		params["user_id"] = userID
	}
	return s.sdk.Unsorted().Decline(ctx, uid, params)
}

func (s *service) LinkUnsorted(ctx context.Context, uid string, leadID int) (*models.UnsortedLinkResult, error) {
	params := map[string]interface{}{
		"lead_id": leadID,
	}
	return s.sdk.Unsorted().Link(ctx, uid, params)
}

func (s *service) SummaryUnsorted(ctx context.Context, filter *gkitmodels.UnsortedFilter) (*models.UnsortedSummary, error) {
	f := filters.NewUnsortedSummaryFilter()
	if filter != nil {
		if len(filter.PipelineID) > 0 {
			f.SetPipelineID(filter.PipelineID[0])
		}
	}
	return s.sdk.Unsorted().Summary(ctx, f)
}
