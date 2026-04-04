package admin_pipelines

import (
	"context"
	"fmt"
	"net/url"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListStatuses(ctx context.Context, pipelineID int) ([]*models.Status, error) {
	statuses, _, err := s.sdk.Statuses(pipelineID).Get(ctx, url.Values{})
	return statuses, err
}

func (s *service) GetStatus(ctx context.Context, pipelineID, statusID int) (*models.Status, error) {
	return s.sdk.Statuses(pipelineID).GetOne(ctx, statusID, url.Values{})
}

func (s *service) CreateStatus(ctx context.Context, pipelineID int, status *models.Status) (*models.Status, error) {
	res, _, err := s.sdk.Statuses(pipelineID).Create(ctx, []*models.Status{status})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no status returned after creation")
	}
	return res[0], nil
}

func (s *service) CreateStatuses(ctx context.Context, pipelineID int, statuses []*models.Status) ([]*models.Status, error) {
	res, _, err := s.sdk.Statuses(pipelineID).Create(ctx, statuses)
	return res, err
}

func (s *service) UpdateStatus(ctx context.Context, pipelineID int, status *models.Status) (*models.Status, error) {
	return s.sdk.Statuses(pipelineID).UpdateOne(ctx, status)
}

func (s *service) DeleteStatus(ctx context.Context, pipelineID, statusID int) error {
	return s.sdk.Statuses(pipelineID).DeleteOne(ctx, statusID)
}
