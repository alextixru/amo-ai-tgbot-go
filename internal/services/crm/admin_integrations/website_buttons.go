package admin_integrations

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
)

func (s *service) ListWebsiteButtons(ctx context.Context, filter *services.WebsiteButtonsFilter, with []string) ([]*models.WebsiteButton, error) {
	buttons, _, err := s.sdk.WebsiteButtons().Get(ctx, filter, with)
	return buttons, err
}

func (s *service) GetWebsiteButton(ctx context.Context, id int, with []string) (*models.WebsiteButton, error) {
	return s.sdk.WebsiteButtons().GetOne(ctx, id, with)
}

func (s *service) CreateWebsiteButton(ctx context.Context, req *models.WebsiteButtonCreateRequest) (*models.WebsiteButtonCreateResponse, error) {
	return s.sdk.WebsiteButtons().CreateAsync(ctx, req)
}

func (s *service) UpdateWebsiteButton(ctx context.Context, req *models.WebsiteButtonUpdateRequest) (*models.WebsiteButton, error) {
	return s.sdk.WebsiteButtons().UpdateAsync(ctx, req)
}

func (s *service) AddOnlineChat(ctx context.Context, sourceID int) error {
	return s.sdk.WebsiteButtons().AddOnlineChatAsync(ctx, sourceID)
}
