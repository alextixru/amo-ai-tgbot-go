package admin_pipelines

import (
	"context"
	"net/url"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListPipelines(ctx context.Context) ([]*models.Pipeline, error) {
	pipelines, _, err := s.sdk.Pipelines().Get(ctx, url.Values{})
	return pipelines, err
}

func (s *service) GetPipeline(ctx context.Context, id int) (*models.Pipeline, error) {
	return s.sdk.Pipelines().GetOne(ctx, id)
}

func (s *service) CreatePipelines(ctx context.Context, pipelines []*models.Pipeline) ([]*models.Pipeline, error) {
	res, _, err := s.sdk.Pipelines().Create(ctx, pipelines)
	return res, err
}

func (s *service) UpdatePipeline(ctx context.Context, pipeline *models.Pipeline) (*models.Pipeline, error) {
	return s.sdk.Pipelines().UpdateOne(ctx, pipeline)
}

func (s *service) DeletePipeline(ctx context.Context, id int) error {
	return s.sdk.Pipelines().DeleteOne(ctx, id)
}
