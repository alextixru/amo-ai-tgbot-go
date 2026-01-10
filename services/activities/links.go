package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) ListLinks(ctx context.Context, parent gkitmodels.ParentEntity) ([]*models.EntityLink, error) {
	return s.sdk.Links().Get(ctx, parent.Type, parent.ID, nil)
}

func (s *service) LinkEntity(ctx context.Context, parent gkitmodels.ParentEntity, target *gkitmodels.LinkTarget) ([]*models.EntityLink, error) {
	link := &models.EntityLink{
		ToEntityID:   target.ID,
		ToEntityType: target.Type,
	}
	return s.sdk.Links().Link(ctx, parent.Type, parent.ID, []*models.EntityLink{link})
}

func (s *service) UnlinkEntity(ctx context.Context, parent gkitmodels.ParentEntity, target *gkitmodels.LinkTarget) error {
	link := &models.EntityLink{
		ToEntityID:   target.ID,
		ToEntityType: target.Type,
	}
	return s.sdk.Links().Unlink(ctx, parent.Type, parent.ID, []*models.EntityLink{link})
}
