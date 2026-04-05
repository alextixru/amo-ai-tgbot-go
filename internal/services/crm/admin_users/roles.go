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

func (s *service) ListRoles(ctx context.Context, filter *gkitmodels.AdminUsersFilter) (*PagedRolesResult, error) {
	sdkFilter := filters.NewUsersFilter()
	// Автоматически подгружаем список пользователей каждой роли
	sdkFilter.SetWith("users")

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

	roles, meta, err := s.sdk.Roles().Get(ctx, sdkFilter.ToQueryParams())
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}

	// Client-side фильтрация по имени (API не поддерживает)
	if filter != nil && filter.Name != "" {
		nameLC := strings.ToLower(filter.Name)
		filtered := roles[:0]
		for _, r := range roles {
			if strings.Contains(strings.ToLower(r.Name), nameLC) {
				filtered = append(filtered, r)
			}
		}
		roles = filtered
	}

	result := &PagedRolesResult{
		Items: roles,
	}
	if meta != nil {
		result.Page = meta.Page
		result.Total = meta.TotalItems
		result.HasNext = meta.HasMore
	}

	return result, nil
}

func (s *service) GetRole(ctx context.Context, id int) (*amomodels.Role, error) {
	// Автоматически подгружаем список пользователей роли
	role, err := s.sdk.Roles().GetOne(ctx, id, sdkservices.WithRelations("users"))
	if err != nil {
		return nil, fmt.Errorf("get role %d: %w", id, err)
	}
	return role, nil
}

func (s *service) CreateRoles(ctx context.Context, roles []*amomodels.Role) ([]*amomodels.Role, error) {
	res, _, err := s.sdk.Roles().Create(ctx, roles)
	if err != nil {
		return nil, fmt.Errorf("create roles: %w", err)
	}
	return res, nil
}

func (s *service) UpdateRoles(ctx context.Context, roles []*amomodels.Role) ([]*amomodels.Role, error) {
	res, _, err := s.sdk.Roles().Update(ctx, roles)
	if err != nil {
		return nil, fmt.Errorf("update roles: %w", err)
	}
	return res, nil
}

func (s *service) DeleteRole(ctx context.Context, id int) (*DeleteResult, error) {
	if err := s.sdk.Roles().Delete(ctx, id); err != nil {
		return nil, fmt.Errorf("delete role %d: %w", id, err)
	}
	return &DeleteResult{
		Success:   true,
		DeletedID: id,
	}, nil
}
