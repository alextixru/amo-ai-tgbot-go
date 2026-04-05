package admin_pipelines

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go"
	toolmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// Service определяет бизнес-логику для работы с воронками и статусами.
type Service interface {
	// Pipelines

	// ListPipelines возвращает список воронок.
	// withStatuses=true — включить статусы каждой воронки в один запрос (with=statuses).
	ListPipelines(ctx context.Context, withStatuses bool) (*toolmodels.ListPipelinesOutput, error)

	// GetPipeline возвращает воронку по ID или имени.
	// Если id==0, выполняет lookup по name через ListPipelines.
	// withStatuses=true — включить статусы воронки (WithRelations("statuses")).
	GetPipeline(ctx context.Context, id int, name string, withStatuses bool) (*toolmodels.PipelineOutput, error)

	// CreatePipelines создаёт одну или несколько воронок.
	CreatePipelines(ctx context.Context, data []toolmodels.PipelineData) ([]*toolmodels.PipelineOutput, error)

	// UpdatePipeline обновляет воронку.
	// Если id==0, выполняет lookup по name через ListPipelines.
	UpdatePipeline(ctx context.Context, id int, name string, data toolmodels.PipelineData) (*toolmodels.PipelineOutput, error)

	// DeletePipeline удаляет воронку.
	// Если id==0, выполняет lookup по name через ListPipelines.
	DeletePipeline(ctx context.Context, id int, name string) error

	// Statuses

	// ListStatuses возвращает статусы воронки.
	// Если pipelineID==0, выполняет lookup по pipelineName.
	ListStatuses(ctx context.Context, pipelineID int, pipelineName string) ([]*toolmodels.StatusOutput, error)

	// GetStatus возвращает статус воронки по ID или имени.
	// Если pipelineID==0, выполняет lookup по pipelineName.
	// Если statusID==0, выполняет lookup по statusName.
	GetStatus(ctx context.Context, pipelineID int, pipelineName string, statusID int, statusName string) (*toolmodels.StatusOutput, error)

	// CreateStatus создаёт один статус в воронке.
	// Если pipelineID==0, выполняет lookup по pipelineName.
	CreateStatus(ctx context.Context, pipelineID int, pipelineName string, data toolmodels.StatusData) (*toolmodels.StatusOutput, error)

	// CreateStatuses создаёт несколько статусов в воронке (батч).
	// Если pipelineID==0, выполняет lookup по pipelineName.
	CreateStatuses(ctx context.Context, pipelineID int, pipelineName string, data []toolmodels.StatusData) ([]*toolmodels.StatusOutput, error)

	// UpdateStatus обновляет статус воронки.
	// Если pipelineID==0, выполняет lookup по pipelineName.
	// Если statusID==0, выполняет lookup по statusName.
	UpdateStatus(ctx context.Context, pipelineID int, pipelineName string, statusID int, statusName string, data toolmodels.StatusData) (*toolmodels.StatusOutput, error)

	// DeleteStatus удаляет статус воронки.
	// Если pipelineID==0, выполняет lookup по pipelineName.
	// Если statusID==0, выполняет lookup по statusName.
	DeleteStatus(ctx context.Context, pipelineID int, pipelineName string, statusID int, statusName string) error
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

// resolvePipelineID возвращает ID воронки по имени (если id == 0).
// Возвращает ошибку с подсказкой о доступных воронках если имя не найдено.
func (s *service) resolvePipelineID(ctx context.Context, id int, name string) (int, error) {
	if id != 0 {
		return id, nil
	}
	if name == "" {
		return 0, fmt.Errorf("pipeline_id или pipeline_name обязателен")
	}

	result, err := s.ListPipelines(ctx, false)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить список воронок: %w", err)
	}

	var available []string
	for _, p := range result.Pipelines {
		if p.Name == name {
			return p.ID, nil
		}
		available = append(available, p.Name)
	}

	return 0, fmt.Errorf("воронка %q не найдена. Доступные: %v", name, available)
}

// resolveStatusID возвращает ID статуса по имени (если id == 0).
// Возвращает ошибку с подсказкой о доступных статусах если имя не найдено.
func (s *service) resolveStatusID(ctx context.Context, pipelineID int, id int, name string) (int, error) {
	if id != 0 {
		return id, nil
	}
	if name == "" {
		return 0, fmt.Errorf("status_id или status_name обязателен")
	}

	statuses, err := s.ListStatuses(ctx, pipelineID, "")
	if err != nil {
		return 0, fmt.Errorf("не удалось получить список статусов: %w", err)
	}

	var available []string
	for _, st := range statuses {
		if st.Name == name {
			return st.ID, nil
		}
		available = append(available, st.Name)
	}

	return 0, fmt.Errorf("статус %q не найден. Доступные: %v", name, available)
}
