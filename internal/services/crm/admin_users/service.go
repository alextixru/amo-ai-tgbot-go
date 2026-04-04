package admin_users

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// Service определяет бизнес-логику для работы с пользователями и ролями.
type Service interface {
	// Users
	ListUsers(ctx context.Context, filter *gkitmodels.AdminUsersFilter) ([]*amomodels.User, error)
	GetUser(ctx context.Context, id int, with []string) (*amomodels.User, error)
	CreateUsers(ctx context.Context, users []*amomodels.User) ([]*amomodels.User, error)

	// Roles
	ListRoles(ctx context.Context, filter *gkitmodels.AdminUsersFilter) ([]*amomodels.Role, error)
	GetRole(ctx context.Context, id int, with []string) (*amomodels.Role, error)
	CreateRoles(ctx context.Context, roles []*amomodels.Role) ([]*amomodels.Role, error)
	UpdateRoles(ctx context.Context, roles []*amomodels.Role) ([]*amomodels.Role, error)
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
