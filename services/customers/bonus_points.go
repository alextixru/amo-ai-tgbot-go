package customers

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) GetBonusPoints(ctx context.Context, customerID int) (*models.BonusPoints, error) {
	return s.sdk.CustomerBonusPoints(customerID).Get(ctx)
}

func (s *service) EarnBonusPoints(ctx context.Context, customerID int, points int) (int, error) {
	return s.sdk.CustomerBonusPoints(customerID).EarnPoints(ctx, points)
}

func (s *service) RedeemBonusPoints(ctx context.Context, customerID int, points int) (int, error) {
	return s.sdk.CustomerBonusPoints(customerID).RedeemPoints(ctx, points)
}
