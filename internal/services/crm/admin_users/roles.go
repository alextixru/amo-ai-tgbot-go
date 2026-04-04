package admin_users

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) ListRoles(ctx context.Context, filter *gkitmodels.AdminUsersFilter) ([]*amomodels.Role, error) {
	params := url.Values{}
	if filter != nil {
		if filter.Limit > 0 {
			params.Set("limit", strconv.Itoa(filter.Limit))
		}
		if filter.Page > 0 {
			params.Set("page", strconv.Itoa(filter.Page))
		}
		if len(filter.With) > 0 {
			params.Set("with", strings.Join(filter.With, ","))
		}
	}
	roles, _, err := s.sdk.Roles().Get(ctx, params)
	return roles, err
}

func (s *service) GetRole(ctx context.Context, id int, with []string) (*amomodels.Role, error) {
	// Note: SDK Roles().GetOne supports ...GetOneOption, but we don't have the correct option for With fields here.
	return s.sdk.Roles().GetOne(ctx, id)
}

func (s *service) CreateRoles(ctx context.Context, roles []*amomodels.Role) ([]*amomodels.Role, error) {
	res, _, err := s.sdk.Roles().Create(ctx, roles)
	return res, err
}

func (s *service) UpdateRoles(ctx context.Context, roles []*amomodels.Role) ([]*amomodels.Role, error) {
	res, _, err := s.sdk.Roles().Update(ctx, roles)
	return res, err
}

func (s *service) DeleteRole(ctx context.Context, id int) error {
	return s.sdk.Roles().Delete(ctx, id)
}
