package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) ListEvents(ctx context.Context, parent gkitmodels.ParentEntity) ([]*models.Event, error) {
	f := filters.NewEventsFilter()
	f.SetLimit(50)
	f.SetPage(1)
	f.SetEntity([]string{parent.Type})
	f.SetEntityIDs([]int{parent.ID})
	events, _, err := s.sdk.Events().Get(ctx, f.ToQueryParams())
	return events, err
}

func (s *service) GetEvent(ctx context.Context, id int) (*models.Event, error) {
	return s.sdk.Events().GetOne(ctx, id)
}
