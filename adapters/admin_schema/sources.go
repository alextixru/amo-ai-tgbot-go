package admin_schema

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListSources(ctx context.Context, filter *filters.SourcesFilter) ([]*models.Source, error) {
	sources, _, err := s.sdk.Sources().Get(ctx, filter)
	return sources, err
}

func (s *service) GetSource(ctx context.Context, id int) (*models.Source, error) {
	return s.sdk.Sources().GetOne(ctx, id)
}

func (s *service) CreateSources(ctx context.Context, sources []*models.Source) ([]*models.Source, error) {
	res, _, err := s.sdk.Sources().Create(ctx, sources)
	return res, err
}

func (s *service) UpdateSources(ctx context.Context, sources []*models.Source) ([]*models.Source, error) {
	res, _, err := s.sdk.Sources().Update(ctx, sources)
	return res, err
}

func (s *service) DeleteSource(ctx context.Context, id int) error {
	return s.sdk.Sources().Delete(ctx, id)
}
