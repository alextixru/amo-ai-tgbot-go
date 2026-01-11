package admin_schema

import (
	"context"
	"net/url"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListCustomFields(ctx context.Context, entityType string, filter *filters.CustomFieldsFilter) ([]*models.CustomField, error) {
	fields, _, err := s.sdk.CustomFields().Get(ctx, entityType, filter)
	return fields, err
}

func (s *service) GetCustomField(ctx context.Context, entityType string, id int) (*models.CustomField, error) {
	return s.sdk.CustomFields().GetOne(ctx, entityType, id, url.Values{})
}

func (s *service) CreateCustomFields(ctx context.Context, entityType string, fields []*models.CustomField) ([]*models.CustomField, error) {
	res, _, err := s.sdk.CustomFields().Create(ctx, entityType, fields)
	return res, err
}

func (s *service) UpdateCustomFields(ctx context.Context, entityType string, fields []*models.CustomField) ([]*models.CustomField, error) {
	res, _, err := s.sdk.CustomFields().Update(ctx, entityType, fields)
	return res, err
}

func (s *service) DeleteCustomField(ctx context.Context, entityType string, id int) error {
	return s.sdk.CustomFields().Delete(ctx, entityType, id)
}
