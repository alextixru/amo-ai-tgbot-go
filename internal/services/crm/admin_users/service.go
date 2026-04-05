package admin_users

import (
	"context"
	"sync"

	"github.com/alextixru/amocrm-sdk-go"
	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// UserView — DTO пользователя с расшифрованными именами вместо числовых ID.
// Возвращается LLM вместо сырой SDK-модели.
type UserView struct {
	ID          int              `json:"id"`
	Name        string           `json:"name,omitempty"`
	Email       string           `json:"email,omitempty"`
	Lang        string           `json:"lang,omitempty"`
	AmojoID     string           `json:"amojo_id,omitempty"`
	PhoneNumber string           `json:"phone_number,omitempty"`
	Rank        amomodels.UserRank `json:"user_rank,omitempty"`
	IsAdmin     bool             `json:"is_admin,omitempty"`
	IsActive    bool             `json:"is_active,omitempty"`
	// RoleID числовой ID роли
	RoleID   int    `json:"role_id,omitempty"`
	// RoleName расшифрованное название роли (может быть пустым если роль не загружена)
	RoleName string `json:"role_name,omitempty"`
	// GroupID числовой ID группы по умолчанию
	GroupID   int    `json:"group_id,omitempty"`
	// GroupName расшифрованное название группы (заполняется из embedded groups)
	GroupName string `json:"group_name,omitempty"`
	// Groups все группы пользователя (заполняется при with=group)
	Groups []amomodels.UserGroup `json:"groups,omitempty"`
	// Roles встроенные роли (заполняется при with=role)
	Roles []*amomodels.Role `json:"roles,omitempty"`
}

// PagedUsersResult — пагинированный список пользователей.
type PagedUsersResult struct {
	Items   []*UserView `json:"items"`
	Page    int         `json:"page"`
	Total   int         `json:"total"`
	HasNext bool        `json:"has_next"`
}

// PagedRolesResult — пагинированный список ролей.
type PagedRolesResult struct {
	Items   []*amomodels.Role `json:"items"`
	Page    int               `json:"page"`
	Total   int               `json:"total"`
	HasNext bool              `json:"has_next"`
}

// DeleteResult — результат удаления сущности.
type DeleteResult struct {
	Success   bool `json:"success"`
	DeletedID int  `json:"deleted_id"`
}

// Service определяет бизнес-логику для работы с пользователями и ролями.
type Service interface {
	// Users

	// ListUsers возвращает пагинированный список пользователей с расшифрованными ID.
	// Автоматически загружает with=role,group. Поддерживает client-side фильтрацию по Name и Email.
	ListUsers(ctx context.Context, filter *gkitmodels.AdminUsersFilter) (*PagedUsersResult, error)

	// GetUser возвращает пользователя по ID с расшифрованными ролью и группой.
	// Автоматически запрашивает with=role,group,uuid,phone_number,user_rank.
	GetUser(ctx context.Context, id int) (*UserView, error)

	// CreateUsers создаёт новых пользователей.
	CreateUsers(ctx context.Context, users []*amomodels.User) ([]*amomodels.User, error)

	// Roles

	// ListRoles возвращает пагинированный список ролей.
	// Автоматически загружает with=users. Поддерживает client-side фильтрацию по Name.
	ListRoles(ctx context.Context, filter *gkitmodels.AdminUsersFilter) (*PagedRolesResult, error)

	// GetRole возвращает роль по ID.
	// Автоматически запрашивает with=users.
	GetRole(ctx context.Context, id int) (*amomodels.Role, error)

	// CreateRoles создаёт новые роли.
	CreateRoles(ctx context.Context, roles []*amomodels.Role) ([]*amomodels.Role, error)

	// UpdateRoles обновляет роли.
	UpdateRoles(ctx context.Context, roles []*amomodels.Role) ([]*amomodels.Role, error)

	// DeleteRole удаляет роль по ID. Возвращает подтверждение удаления.
	DeleteRole(ctx context.Context, id int) (*DeleteResult, error)
}

type service struct {
	sdk *amocrm.SDK

	// rolesCache lazy-cache: id роли → название роли.
	// Заполняется при первом вызове ListUsers или GetUser.
	// Сброс не нужен для MVP (рестарт = перезагрузка).
	rolesMu    sync.RWMutex
	rolesCache map[int]string
}

// NewService создает новый экземпляр сервиса пользователей и ролей.
func NewService(sdk *amocrm.SDK) Service {
	return &service{
		sdk:        sdk,
		rolesCache: make(map[int]string),
	}
}
