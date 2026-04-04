package admin_schema

import (
	"context"
	"net/url"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListLossReasons(ctx context.Context, filter url.Values) ([]*models.LossReason, error) {
	reasons, _, err := s.sdk.LossReasons().Get(ctx, filter)
	return reasons, err
}

func (s *service) GetLossReason(ctx context.Context, id int) (*models.LossReason, error) {
	return s.sdk.LossReasons().GetOne(ctx, id)
}

func (s *service) CreateLossReasons(ctx context.Context, reasons []*models.LossReason) ([]*models.LossReason, error) {
	res, _, err := s.sdk.LossReasons().Create(ctx, reasons)
	return res, err
}

func (s *service) DeleteLossReason(ctx context.Context, id int) error {
	return s.sdk.LossReasons().DeleteOne(ctx, id)
}
