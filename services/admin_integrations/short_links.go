package admin_integrations

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListShortLinks(ctx context.Context) ([]models.ShortLink, error) {
	links, _, err := s.sdk.ShortLinks().Get(ctx, nil)
	return links, err
}

func (s *service) CreateShortLink(ctx context.Context, url string) (models.ShortLink, error) {
	links := []models.ShortLink{{URL: url}}
	res, _, err := s.sdk.ShortLinks().Create(ctx, links)
	if err != nil {
		return models.ShortLink{}, err
	}
	if len(res) == 0 {
		return models.ShortLink{}, nil
	}
	return res[0], nil
}

func (s *service) DeleteShortLink(ctx context.Context, id int) error {
	return s.sdk.ShortLinks().Delete(ctx, id)
}
