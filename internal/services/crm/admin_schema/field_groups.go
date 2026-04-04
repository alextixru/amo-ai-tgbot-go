package admin_schema

import (
	"context"
	"net/url"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListFieldGroups(ctx context.Context, entityType string, filter url.Values) ([]models.CustomFieldGroup, error) {
	groups, _, err := s.sdk.CustomFieldGroups(entityType).Get(ctx, filter)
	return groups, err
}

func (s *service) GetFieldGroup(ctx context.Context, entityType string, id string) (*models.CustomFieldGroup, error) {
	res, err := s.sdk.CustomFieldGroups(entityType).GetOne(ctx, id, url.Values{})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *service) CreateFieldGroups(ctx context.Context, entityType string, groups []models.CustomFieldGroup) ([]models.CustomFieldGroup, error) {
	res, _, err := s.sdk.CustomFieldGroups(entityType).Create(ctx, groups)
	return res, err
}

func (s *service) UpdateFieldGroups(ctx context.Context, entityType string, groups []models.CustomFieldGroup) ([]models.CustomFieldGroup, error) {
	res, _, err := s.sdk.CustomFieldGroups(entityType).Update(ctx, groups)
	return res, err
}

func (s *service) DeleteFieldGroup(ctx context.Context, entityType string, id string) error {
	return s.sdk.CustomFieldGroups(entityType).Delete(ctx, id)
}
