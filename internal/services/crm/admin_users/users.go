package admin_users

import (
	"context"
	"fmt"
	"strings"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	sdkservices "github.com/alextixru/amocrm-sdk-go/core/services"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// warmRolesCache загружает все роли в кеш если он ещё пуст.
// Потокобезопасно. Используется для расшифровки role_id → role_name.
func (s *service) warmRolesCache(ctx context.Context) error {
	s.rolesMu.RLock()
	populated := len(s.rolesCache) > 0
	s.rolesMu.RUnlock()
	if populated {
		return nil
	}

	s.rolesMu.Lock()
	defer s.rolesMu.Unlock()
	// Повторная проверка после взятия write-lock (double-checked locking)
	if len(s.rolesCache) > 0 {
		return nil
	}

	roles, _, err := s.sdk.Roles().Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("warm roles cache: %w", err)
	}
	for _, r := range roles {
		if r != nil {
			s.rolesCache[r.ID] = r.Name
		}
	}
	return nil
}

// roleName возвращает название роли по ID из кеша. При отсутствии — пустую строку.
func (s *service) roleName(id int) string {
	if id == 0 {
		return ""
	}
	s.rolesMu.RLock()
	defer s.rolesMu.RUnlock()
	return s.rolesCache[id]
}

// toUserView конвертирует SDK-модель пользователя в DTO с расшифрованными именами.
func (s *service) toUserView(u *amomodels.User) *UserView {
	v := &UserView{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Lang:        u.Lang,
		AmojoID:     u.AmojoID,
		PhoneNumber: u.PhoneNumber,
		Rank:        u.Rank,
	}
	if u.Rights != nil {
		v.IsAdmin = u.Rights.IsAdmin
		v.IsActive = u.Rights.IsActive
		v.RoleID = u.Rights.RoleID
		v.GroupID = u.Rights.GroupID
		v.RoleName = s.roleName(u.Rights.RoleID)
	}
	if u.Embedded != nil {
		v.Roles = u.Embedded.Roles
		v.Groups = u.Embedded.Groups
		// Расшифровываем group_id → group_name из embedded groups
		if v.GroupID != 0 && len(u.Embedded.Groups) > 0 {
			for _, g := range u.Embedded.Groups {
				if g.ID == v.GroupID {
					v.GroupName = g.Name
					break
				}
			}
		}
	}
	return v
}

func (s *service) ListUsers(ctx context.Context, filter *gkitmodels.AdminUsersFilter) (*PagedUsersResult, error) {
	// Прогреваем кеш ролей для расшифровки role_id → role_name
	if err := s.warmRolesCache(ctx); err != nil {
		// Не фатально: продолжаем без расшифровки имён ролей
		_ = err
	}

	sdkFilter := filters.NewUsersFilter()
	// Автоматически подгружаем роль и группу — сервис сам решает что загружать
	sdkFilter.SetWith("role", "group")

	if filter != nil {
		if filter.Limit > 0 {
			sdkFilter.SetLimit(filter.Limit)
		}
		if filter.Page > 0 {
			sdkFilter.SetPage(filter.Page)
		}
		if len(filter.Order) > 0 {
			for field, dir := range filter.Order {
				sdkFilter.SetOrder(field, dir)
			}
		}
	}

	users, meta, err := s.sdk.Users().Get(ctx, sdkFilter)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	// Client-side фильтрация по имени и email (API не поддерживает)
	if filter != nil {
		nameLC := strings.ToLower(filter.Name)
		emailLC := strings.ToLower(filter.Email)
		if nameLC != "" || emailLC != "" {
			filtered := users[:0]
			for _, u := range users {
				if nameLC != "" && !strings.Contains(strings.ToLower(u.Name), nameLC) {
					continue
				}
				if emailLC != "" && !strings.Contains(strings.ToLower(u.Email), emailLC) {
					continue
				}
				filtered = append(filtered, u)
			}
			users = filtered
		}
	}

	views := make([]*UserView, 0, len(users))
	for _, u := range users {
		views = append(views, s.toUserView(u))
	}

	result := &PagedUsersResult{
		Items: views,
	}
	if meta != nil {
		result.Page = meta.Page
		result.Total = meta.TotalItems
		result.HasNext = meta.HasMore
	}

	return result, nil
}

func (s *service) GetUser(ctx context.Context, id int) (*UserView, error) {
	// Прогреваем кеш ролей для расшифровки role_id → role_name
	if err := s.warmRolesCache(ctx); err != nil {
		_ = err
	}

	// Запрашиваем все доступные with-параметры
	user, err := s.sdk.Users().GetOne(ctx, id,
		sdkservices.WithRelations("role", "group", "uuid", "amojo_id", "user_rank", "phone_number"),
	)
	if err != nil {
		return nil, fmt.Errorf("get user %d: %w", id, err)
	}
	return s.toUserView(user), nil
}

func (s *service) CreateUsers(ctx context.Context, users []*amomodels.User) ([]*amomodels.User, error) {
	// Валидация обязательных полей до вызова API
	for i, u := range users {
		if u.Name == "" {
			return nil, fmt.Errorf("users[%d]: name is required", i)
		}
		if u.Email == "" {
			return nil, fmt.Errorf("users[%d]: email is required", i)
		}
	}

	res, _, err := s.sdk.Users().Create(ctx, users)
	if err != nil {
		return nil, fmt.Errorf("create users: %w", err)
	}
	return res, nil
}
