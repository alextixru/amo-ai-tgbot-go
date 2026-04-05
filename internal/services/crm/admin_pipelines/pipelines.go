package admin_pipelines

import (
	"context"
	"net/url"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	sdkservices "github.com/alextixru/amocrm-sdk-go/core/services"
	toolmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) ListPipelines(ctx context.Context, withStatuses bool) (*toolmodels.ListPipelinesOutput, error) {
	params := url.Values{}
	if withStatuses {
		params.Set("with", "statuses")
	}

	pipelines, meta, err := s.sdk.Pipelines().Get(ctx, params)
	if err != nil {
		return nil, err
	}

	out := make([]*toolmodels.PipelineOutput, 0, len(pipelines))
	for _, p := range pipelines {
		out = append(out, pipelineToOutput(p))
	}

	return &toolmodels.ListPipelinesOutput{
		Pipelines: out,
		PageMeta:  meta,
	}, nil
}

func (s *service) GetPipeline(ctx context.Context, id int, name string, withStatuses bool) (*toolmodels.PipelineOutput, error) {
	resolvedID, err := s.resolvePipelineID(ctx, id, name)
	if err != nil {
		return nil, err
	}

	var p *amomodels.Pipeline
	if withStatuses {
		p, err = s.sdk.Pipelines().GetOne(ctx, resolvedID, sdkservices.WithRelations("statuses"))
	} else {
		p, err = s.sdk.Pipelines().GetOne(ctx, resolvedID)
	}
	if err != nil {
		return nil, err
	}

	return pipelineToOutput(p), nil
}

func (s *service) CreatePipelines(ctx context.Context, data []toolmodels.PipelineData) ([]*toolmodels.PipelineOutput, error) {
	pipelines := make([]*amomodels.Pipeline, 0, len(data))
	for _, d := range data {
		pipelines = append(pipelines, pipelineDataToModel(d))
	}

	created, _, err := s.sdk.Pipelines().Create(ctx, pipelines)
	if err != nil {
		return nil, err
	}

	out := make([]*toolmodels.PipelineOutput, 0, len(created))
	for _, p := range created {
		out = append(out, pipelineToOutput(p))
	}
	return out, nil
}

func (s *service) UpdatePipeline(ctx context.Context, id int, name string, data toolmodels.PipelineData) (*toolmodels.PipelineOutput, error) {
	resolvedID, err := s.resolvePipelineID(ctx, id, name)
	if err != nil {
		return nil, err
	}

	p := pipelineDataToModel(data)
	p.ID = resolvedID

	updated, err := s.sdk.Pipelines().UpdateOne(ctx, p)
	if err != nil {
		return nil, err
	}

	return pipelineToOutput(updated), nil
}

func (s *service) DeletePipeline(ctx context.Context, id int, name string) error {
	resolvedID, err := s.resolvePipelineID(ctx, id, name)
	if err != nil {
		return err
	}
	return s.sdk.Pipelines().DeleteOne(ctx, resolvedID)
}

// --- helpers ---

// pipelineToOutput конвертирует SDK-модель воронки в LLM-friendly вывод.
// Убирает account_id и _links, добавляет статусы если они были загружены.
func pipelineToOutput(p *amomodels.Pipeline) *toolmodels.PipelineOutput {
	out := &toolmodels.PipelineOutput{
		ID:           p.ID,
		Name:         p.Name,
		Sort:         p.Sort,
		IsMain:       p.IsMain,
		IsUnsortedOn: p.IsUnsortedOn,
		IsArchive:    p.IsArchive,
	}

	if p.Embedded != nil && len(p.Embedded.Statuses) > 0 {
		out.Statuses = make([]*toolmodels.StatusOutput, 0, len(p.Embedded.Statuses))
		for i := range p.Embedded.Statuses {
			out.Statuses = append(out.Statuses, statusToOutput(&p.Embedded.Statuses[i]))
		}
	}

	return out
}

// pipelineDataToModel конвертирует типизированные данные от LLM в SDK-модель.
func pipelineDataToModel(d toolmodels.PipelineData) *amomodels.Pipeline {
	return &amomodels.Pipeline{
		Name:         d.Name,
		Sort:         d.Sort,
		IsMain:       d.IsMain,
		IsUnsortedOn: d.IsUnsortedOn,
	}
}
