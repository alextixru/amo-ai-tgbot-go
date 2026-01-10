package admin_users

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

// Service определяет бизнес-логику для работы с пользователями и ролями.
type Service interface {
	// Users
	ListUsers(ctx context.Context) ([]*models.User, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	CreateUsers(ctx context.Context, users []*models.User) ([]*models.User, error)

	// Roles
	ListRoles(ctx context.Context) ([]*models.Role, error)
	GetRole(ctx context.Context, id int) (*models.Role, error)
	CreateRoles(ctx context.Context, roles []*models.Role) ([]*models.Role, error)
	UpdateRoles(ctx context.Context, roles []*models.Role) ([]*models.Role, error)
	DeleteRole(ctx context.Context, id int) error
}

type service struct {
	sdk *amocrm.SDK
}

// NewService создает новый экземпляр сервиса пользователей и ролей.
func NewService(sdk *amocrm.SDK) Service {
	return &service{
		sdk: sdk,
	}
}
