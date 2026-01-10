package admin_integrations

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListWebsiteButtons(ctx context.Context) ([]*models.WebsiteButton, error) {
	buttons, _, err := s.sdk.WebsiteButtons().Get(ctx, nil, nil)
	return buttons, err
}

func (s *service) GetWebsiteButton(ctx context.Context, id int) (*models.WebsiteButton, error) {
	return s.sdk.WebsiteButtons().GetOne(ctx, id, nil)
}

func (s *service) CreateWebsiteButton(ctx context.Context, req *models.WebsiteButtonCreateRequest) (*models.WebsiteButtonCreateResponse, error) {
	return s.sdk.WebsiteButtons().CreateAsync(ctx, req)
}

func (s *service) UpdateWebsiteButton(ctx context.Context, req *models.WebsiteButtonUpdateRequest) (*models.WebsiteButton, error) {
	return s.sdk.WebsiteButtons().UpdateAsync(ctx, req)
}
