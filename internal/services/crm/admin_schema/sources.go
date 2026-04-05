package admin_schema

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListSources(ctx context.Context, filter *filters.SourcesFilter) (*PagedResult[*models.Source], error) {
	sources, meta, err := s.sdk.Sources().Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	return newPagedResult(sources, meta), nil
}

func (s *service) GetSource(ctx context.Context, id int) (*models.Source, error) {
	return s.sdk.Sources().GetOne(ctx, id)
}

func (s *service) CreateSources(ctx context.Context, sources []*models.Source) (*PagedResult[*models.Source], error) {
	res, meta, err := s.sdk.Sources().Create(ctx, sources)
	if err != nil {
		return nil, err
	}
	return newPagedResult(res, meta), nil
}

func (s *service) UpdateSources(ctx context.Context, sources []*models.Source) (*PagedResult[*models.Source], error) {
	res, meta, err := s.sdk.Sources().Update(ctx, sources)
	if err != nil {
		return nil, err
	}
	return newPagedResult(res, meta), nil
}

func (s *service) DeleteSource(ctx context.Context, id int) (*DeleteResult, error) {
	if err := s.sdk.Sources().Delete(ctx, id); err != nil {
		return nil, err
	}
	return &DeleteResult{Success: true, DeletedID: id}, nil
}
