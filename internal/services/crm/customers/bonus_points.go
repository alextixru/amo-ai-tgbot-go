package customers

import (
	"context"
)

func (s *service) GetBonusPoints(ctx context.Context, customerID int) (*BonusPointsInfo, error) {
	bp, err := s.sdk.CustomerBonusPoints(customerID).Get(ctx)
	if err != nil {
		return nil, err
	}
	if bp == nil {
		return nil, nil
	}
	return &BonusPointsInfo{
		BonusPoints: bp.BonusPoints,
	}, nil
}

func (s *service) EarnBonusPoints(ctx context.Context, customerID int, points int) (*BonusPointsResult, error) {
	balance, err := s.sdk.CustomerBonusPoints(customerID).EarnPoints(ctx, points)
	if err != nil {
		return nil, err
	}
	return &BonusPointsResult{
		Balance:   balance,
		Operation: "earn",
		Points:    points,
	}, nil
}

func (s *service) RedeemBonusPoints(ctx context.Context, customerID int, points int) (*BonusPointsResult, error) {
	balance, err := s.sdk.CustomerBonusPoints(customerID).RedeemPoints(ctx, points)
	if err != nil {
		return nil, err
	}
	return &BonusPointsResult{
		Balance:   balance,
		Operation: "redeem",
		Points:    points,
	}, nil
}
