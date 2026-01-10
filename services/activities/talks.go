package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) ListTalks(ctx context.Context, parent gkitmodels.ParentEntity) ([]*models.Talk, error) {
	params := map[string]string{
		"limit":       "50",
		"page":        "1",
		"entity_type": parent.Type,
	}
	return s.sdk.Talks().Get(ctx, params)
}

func (s *service) CloseTalk(ctx context.Context, talkID string) error {
	return s.sdk.Talks().Close(ctx, talkID, nil)
}
