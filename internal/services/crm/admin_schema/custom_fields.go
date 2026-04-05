package admin_schema

import (
	"context"
	"net/url"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListCustomFields(ctx context.Context, entityType string, filter *filters.CustomFieldsFilter) (*PagedResult[*models.CustomField], error) {
	fields, meta, err := s.sdk.CustomFields().Get(ctx, entityType, filter)
	if err != nil {
		return nil, err
	}
	return newPagedResult(fields, meta), nil
}

func (s *service) GetCustomField(ctx context.Context, entityType string, id int) (*models.CustomField, error) {
	return s.sdk.CustomFields().GetOne(ctx, entityType, id, url.Values{})
}

func (s *service) CreateCustomFields(ctx context.Context, entityType string, fields []*models.CustomField) (*PagedResult[*models.CustomField], error) {
	res, meta, err := s.sdk.CustomFields().Create(ctx, entityType, fields)
	if err != nil {
		return nil, err
	}
	return newPagedResult(res, meta), nil
}

func (s *service) UpdateCustomFields(ctx context.Context, entityType string, fields []*models.CustomField) (*PagedResult[*models.CustomField], error) {
	res, meta, err := s.sdk.CustomFields().Update(ctx, entityType, fields)
	if err != nil {
		return nil, err
	}
	return newPagedResult(res, meta), nil
}

func (s *service) DeleteCustomField(ctx context.Context, entityType string, id int) (*DeleteResult, error) {
	if err := s.sdk.CustomFields().Delete(ctx, entityType, id); err != nil {
		return nil, err
	}
	return &DeleteResult{Success: true, DeletedID: id}, nil
}
