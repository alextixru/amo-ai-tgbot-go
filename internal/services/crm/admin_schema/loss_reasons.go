package admin_schema

import (
	"context"
	"net/url"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListLossReasons(ctx context.Context, filter url.Values) (*PagedResult[*models.LossReason], error) {
	reasons, meta, err := s.sdk.LossReasons().Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	return newPagedResult(reasons, meta), nil
}

func (s *service) GetLossReason(ctx context.Context, id int) (*models.LossReason, error) {
	return s.sdk.LossReasons().GetOne(ctx, id)
}

func (s *service) CreateLossReasons(ctx context.Context, reasons []*models.LossReason) (*PagedResult[*models.LossReason], error) {
	res, meta, err := s.sdk.LossReasons().Create(ctx, reasons)
	if err != nil {
		return nil, err
	}
	return newPagedResult(res, meta), nil
}

func (s *service) DeleteLossReason(ctx context.Context, id int) (*DeleteResult, error) {
	if err := s.sdk.LossReasons().DeleteOne(ctx, id); err != nil {
		return nil, err
	}
	return &DeleteResult{Success: true, DeletedID: id}, nil
}
