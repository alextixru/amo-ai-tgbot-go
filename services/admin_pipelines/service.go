package admin_pipelines

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

// Service определяет бизнес-логику для работы с воронками и статусами.
type Service interface {
	// Pipelines
	ListPipelines(ctx context.Context) ([]*models.Pipeline, error)
	GetPipeline(ctx context.Context, id int) (*models.Pipeline, error)
	CreatePipelines(ctx context.Context, pipelines []*models.Pipeline) ([]*models.Pipeline, error)
	UpdatePipeline(ctx context.Context, pipeline *models.Pipeline) (*models.Pipeline, error)
	DeletePipeline(ctx context.Context, id int) error

	// Statuses
	ListStatuses(ctx context.Context, pipelineID int) ([]*models.Status, error)
	GetStatus(ctx context.Context, pipelineID, statusID int) (*models.Status, error)
	CreateStatus(ctx context.Context, pipelineID int, status *models.Status) (*models.Status, error)
	UpdateStatus(ctx context.Context, pipelineID int, status *models.Status) (*models.Status, error)
	DeleteStatus(ctx context.Context, pipelineID, statusID int) error
}

type service struct {
	sdk *amocrm.SDK
}

// New создает новый экземпляр сервиса воронок.
func New(sdk *amocrm.SDK) Service {
	return &service{
		sdk: sdk,
	}
}
