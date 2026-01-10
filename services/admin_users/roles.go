package admin_users

import (
	"context"
	"net/url"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListRoles(ctx context.Context) ([]*models.Role, error) {
	roles, _, err := s.sdk.Roles().Get(ctx, url.Values{})
	return roles, err
}

func (s *service) GetRole(ctx context.Context, id int) (*models.Role, error) {
	return s.sdk.Roles().GetOne(ctx, id)
}

func (s *service) CreateRoles(ctx context.Context, roles []*models.Role) ([]*models.Role, error) {
	res, _, err := s.sdk.Roles().Create(ctx, roles)
	return res, err
}

func (s *service) UpdateRoles(ctx context.Context, roles []*models.Role) ([]*models.Role, error) {
	res, _, err := s.sdk.Roles().Update(ctx, roles)
	return res, err
}

func (s *service) DeleteRole(ctx context.Context, id int) error {
	return s.sdk.Roles().Delete(ctx, id)
}
