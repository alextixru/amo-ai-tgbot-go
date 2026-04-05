package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/services"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) ListSubscriptions(ctx context.Context, parent gkitmodels.ParentEntity) (*SubscriptionsListOutput, error) {
	svc := services.NewEntitySubscriptionsService(s.sdk.Client(), parent.Type, parent.ID)
	subs, err := svc.Get(ctx, 1, 50)
	if err != nil {
		return nil, err
	}
	out := &SubscriptionsListOutput{
		Subscriptions: make([]SubscriptionOutput, 0, len(subs)),
	}
	for _, sub := range subs {
		out.Subscriptions = append(out.Subscriptions, SubscriptionOutput{
			SubscriberID:   sub.SubscriberID,
			SubscriberName: s.resolveUserID(sub.SubscriberID),
		})
	}
	return out, nil
}

func (s *service) Subscribe(ctx context.Context, parent gkitmodels.ParentEntity, userNames []string) (*SubscriptionsListOutput, error) {
	ids, err := s.resolveUserNames(userNames)
	if err != nil {
		return nil, err
	}
	svc := services.NewEntitySubscriptionsService(s.sdk.Client(), parent.Type, parent.ID)
	subs, err := svc.Subscribe(ctx, ids)
	if err != nil {
		return nil, err
	}
	out := &SubscriptionsListOutput{
		Subscriptions: make([]SubscriptionOutput, 0, len(subs)),
	}
	for _, sub := range subs {
		out.Subscriptions = append(out.Subscriptions, SubscriptionOutput{
			SubscriberID:   sub.SubscriberID,
			SubscriberName: s.resolveUserID(sub.SubscriberID),
		})
	}
	return out, nil
}

func (s *service) Unsubscribe(ctx context.Context, parent gkitmodels.ParentEntity, userName string) error {
	uid, err := s.resolveUserName(userName)
	if err != nil {
		return err
	}
	svc := services.NewEntitySubscriptionsService(s.sdk.Client(), parent.Type, parent.ID)
	return svc.Unsubscribe(ctx, uid)
}
