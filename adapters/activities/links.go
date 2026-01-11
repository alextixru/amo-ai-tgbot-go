package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) ListLinks(ctx context.Context, parent gkitmodels.ParentEntity, filter *gkitmodels.LinksFilter) ([]*models.EntityLink, error) {
	f := filters.NewLinksFilter()
	if filter != nil {
		if filter.ToEntityType != "" {
			f.SetToEntityType(filter.ToEntityType)
		}
		if filter.ToEntityID != 0 {
			f.SetToEntityID(filter.ToEntityID)
		}
	}
	return s.sdk.Links().Get(ctx, parent.Type, parent.ID, f)
}

func (s *service) LinkEntity(ctx context.Context, parent gkitmodels.ParentEntity, target *gkitmodels.LinkTarget) ([]*models.EntityLink, error) {
	return s.LinkEntities(ctx, parent, []gkitmodels.LinkTarget{*target})
}

func (s *service) LinkEntities(ctx context.Context, parent gkitmodels.ParentEntity, targets []gkitmodels.LinkTarget) ([]*models.EntityLink, error) {
	items := make([]*models.EntityLink, len(targets))
	for i, t := range targets {
		items[i] = &models.EntityLink{
			ToEntityID:   t.ID,
			ToEntityType: t.Type,
		}
	}
	return s.sdk.Links().Link(ctx, parent.Type, parent.ID, items)
}

func (s *service) UnlinkEntity(ctx context.Context, parent gkitmodels.ParentEntity, target *gkitmodels.LinkTarget) error {
	link := &models.EntityLink{
		ToEntityID:   target.ID,
		ToEntityType: target.Type,
	}
	return s.sdk.Links().Unlink(ctx, parent.Type, parent.ID, []*models.EntityLink{link})
}
