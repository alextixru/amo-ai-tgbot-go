package customers

import (
	"context"
	"net/url"
	"strconv"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListCustomerStatuses(ctx context.Context, page, limit int) ([]models.Status, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	statuses, _, err := s.sdk.CustomerStatuses().Get(ctx, params)
	return statuses, err
}

func (s *service) GetCustomerStatus(ctx context.Context, id int) (models.Status, error) {
	return s.sdk.CustomerStatuses().GetOne(ctx, id)
}

func (s *service) CreateCustomerStatuses(ctx context.Context, statuses []models.Status) ([]models.Status, error) {
	res, _, err := s.sdk.CustomerStatuses().Create(ctx, statuses)
	return res, err
}

func (s *service) UpdateCustomerStatuses(ctx context.Context, statuses []models.Status) ([]models.Status, error) {
	res, _, err := s.sdk.CustomerStatuses().Update(ctx, statuses)
	return res, err
}

func (s *service) DeleteCustomerStatus(ctx context.Context, id int) error {
	return s.sdk.CustomerStatuses().Delete(ctx, id)
}
