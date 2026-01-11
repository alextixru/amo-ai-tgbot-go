package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models"
)

func (s *service) ListEvents(ctx context.Context, parent *gkitmodels.ParentEntity, filter *gkitmodels.EventsFilter) ([]*models.Event, error) {
	f := filters.NewEventsFilter()
	if filter != nil {
		if filter.Limit > 0 {
			f.SetLimit(filter.Limit)
		} else {
			f.SetLimit(50)
		}
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		}
		if len(filter.Types) > 0 {
			f.SetTypes(filter.Types)
		}
		if len(filter.CreatedBy) > 0 {
			f.SetCreatedBy(filter.CreatedBy)
		}
	} else {
		f.SetLimit(50).SetPage(1)
	}

	if parent != nil {
		f.SetEntity([]string{parent.Type})
		f.SetEntityIDs([]int{parent.ID})
	}

	events, _, err := s.sdk.Events().Get(ctx, f.ToQueryParams())
	return events, err
}

func (s *service) GetEvent(ctx context.Context, id int) (*models.Event, error) {
	return s.sdk.Events().GetOne(ctx, id)
}
