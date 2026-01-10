package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) ListSubscriptions(ctx context.Context, parent gkitmodels.ParentEntity) ([]models.Subscription, error) {
	svc := services.NewEntitySubscriptionsService(s.sdk.Client(), parent.Type, parent.ID)
	return svc.Get(ctx, 1, 50)
}

func (s *service) Subscribe(ctx context.Context, parent gkitmodels.ParentEntity, userIDs []int) ([]models.Subscription, error) {
	svc := services.NewEntitySubscriptionsService(s.sdk.Client(), parent.Type, parent.ID)
	return svc.Subscribe(ctx, userIDs)
}

func (s *service) Unsubscribe(ctx context.Context, parent gkitmodels.ParentEntity, userID int) error {
	svc := services.NewEntitySubscriptionsService(s.sdk.Client(), parent.Type, parent.ID)
	return svc.Unsubscribe(ctx, userID)
}
