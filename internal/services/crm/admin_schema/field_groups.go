package admin_schema

import (
	"context"
	"net/url"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListFieldGroups(ctx context.Context, entityType string, filter url.Values) (*PagedResult[models.CustomFieldGroup], error) {
	// with=fields критически важен: без него Fields в каждой группе будет пустым
	if filter == nil {
		filter = url.Values{}
	}
	filter.Set("with", "fields")
	groups, meta, err := s.sdk.CustomFieldGroups(entityType).Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	return newPagedResult(groups, meta), nil
}

func (s *service) GetFieldGroup(ctx context.Context, entityType string, id string) (*models.CustomFieldGroup, error) {
	params := url.Values{}
	params.Set("with", "fields")
	res, err := s.sdk.CustomFieldGroups(entityType).GetOne(ctx, id, params)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *service) CreateFieldGroups(ctx context.Context, entityType string, groups []models.CustomFieldGroup) (*PagedResult[models.CustomFieldGroup], error) {
	res, meta, err := s.sdk.CustomFieldGroups(entityType).Create(ctx, groups)
	if err != nil {
		return nil, err
	}
	return newPagedResult(res, meta), nil
}

func (s *service) UpdateFieldGroups(ctx context.Context, entityType string, groups []models.CustomFieldGroup) (*PagedResult[models.CustomFieldGroup], error) {
	res, meta, err := s.sdk.CustomFieldGroups(entityType).Update(ctx, groups)
	if err != nil {
		return nil, err
	}
	return newPagedResult(res, meta), nil
}

func (s *service) DeleteFieldGroup(ctx context.Context, entityType string, id string) (*DeleteResult, error) {
	if err := s.sdk.CustomFieldGroups(entityType).Delete(ctx, id); err != nil {
		return nil, err
	}
	return &DeleteResult{Success: true, DeletedID: id}, nil
}
