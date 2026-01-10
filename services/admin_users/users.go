package admin_users

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListUsers(ctx context.Context) ([]*models.User, error) {
	users, _, err := s.sdk.Users().Get(ctx, nil)
	return users, err
}

func (s *service) GetUser(ctx context.Context, id int) (*models.User, error) {
	return s.sdk.Users().GetOne(ctx, id)
}

func (s *service) CreateUsers(ctx context.Context, users []*models.User) ([]*models.User, error) {
	res, _, err := s.sdk.Users().Create(ctx, users)
	return res, err
}
