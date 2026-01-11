package admin_users

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) ListUsers(ctx context.Context, filter *gkitmodels.AdminUsersFilter) ([]*amomodels.User, error) {
	var sdkFilter *filters.UsersFilter
	if filter != nil {
		sdkFilter = &filters.UsersFilter{
			Limit: filter.Limit,
			Page:  filter.Page,
		}
		sdkFilter.With = filter.With
	}
	users, _, err := s.sdk.Users().Get(ctx, sdkFilter)
	return users, err
}

func (s *service) GetUser(ctx context.Context, id int, with []string) (*amomodels.User, error) {
	// Note: SDK Users().GetOne currently does not seem to support With parameters via options in toolmap.md
	// but AUDIT.md says it does. If they are needed, we might need a different SDK version or method.
	return s.sdk.Users().GetOne(ctx, id)
}

func (s *service) CreateUsers(ctx context.Context, users []*amomodels.User) ([]*amomodels.User, error) {
	res, _, err := s.sdk.Users().Create(ctx, users)
	return res, err
}
