package admin_pipelines

import (
	"context"
	"fmt"
	"net/url"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	toolmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) ListStatuses(ctx context.Context, pipelineID int, pipelineName string) ([]*toolmodels.StatusOutput, error) {
	resolvedID, err := s.resolvePipelineID(ctx, pipelineID, pipelineName)
	if err != nil {
		return nil, err
	}

	statuses, _, err := s.sdk.Statuses(resolvedID).Get(ctx, url.Values{})
	if err != nil {
		return nil, err
	}

	out := make([]*toolmodels.StatusOutput, 0, len(statuses))
	for _, st := range statuses {
		out = append(out, statusToOutput(st))
	}
	return out, nil
}

func (s *service) GetStatus(ctx context.Context, pipelineID int, pipelineName string, statusID int, statusName string) (*toolmodels.StatusOutput, error) {
	resolvedPipelineID, err := s.resolvePipelineID(ctx, pipelineID, pipelineName)
	if err != nil {
		return nil, err
	}

	resolvedStatusID, err := s.resolveStatusID(ctx, resolvedPipelineID, statusID, statusName)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("with", "descriptions")

	st, err := s.sdk.Statuses(resolvedPipelineID).GetOne(ctx, resolvedStatusID, params)
	if err != nil {
		return nil, err
	}

	return statusToOutput(st), nil
}

func (s *service) CreateStatus(ctx context.Context, pipelineID int, pipelineName string, data toolmodels.StatusData) (*toolmodels.StatusOutput, error) {
	resolvedID, err := s.resolvePipelineID(ctx, pipelineID, pipelineName)
	if err != nil {
		return nil, err
	}

	st := statusDataToModel(data)
	res, _, err := s.sdk.Statuses(resolvedID).Create(ctx, []*amomodels.Status{st})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("статус не был возвращён после создания")
	}
	return statusToOutput(res[0]), nil
}

func (s *service) CreateStatuses(ctx context.Context, pipelineID int, pipelineName string, data []toolmodels.StatusData) ([]*toolmodels.StatusOutput, error) {
	resolvedID, err := s.resolvePipelineID(ctx, pipelineID, pipelineName)
	if err != nil {
		return nil, err
	}

	statuses := make([]*amomodels.Status, 0, len(data))
	for _, d := range data {
		statuses = append(statuses, statusDataToModel(d))
	}

	res, _, err := s.sdk.Statuses(resolvedID).Create(ctx, statuses)
	if err != nil {
		return nil, err
	}

	out := make([]*toolmodels.StatusOutput, 0, len(res))
	for _, st := range res {
		out = append(out, statusToOutput(st))
	}
	return out, nil
}

func (s *service) UpdateStatus(ctx context.Context, pipelineID int, pipelineName string, statusID int, statusName string, data toolmodels.StatusData) (*toolmodels.StatusOutput, error) {
	resolvedPipelineID, err := s.resolvePipelineID(ctx, pipelineID, pipelineName)
	if err != nil {
		return nil, err
	}

	resolvedStatusID, err := s.resolveStatusID(ctx, resolvedPipelineID, statusID, statusName)
	if err != nil {
		return nil, err
	}

	st := statusDataToModel(data)
	st.ID = resolvedStatusID

	updated, err := s.sdk.Statuses(resolvedPipelineID).UpdateOne(ctx, st)
	if err != nil {
		return nil, err
	}

	return statusToOutput(updated), nil
}

func (s *service) DeleteStatus(ctx context.Context, pipelineID int, pipelineName string, statusID int, statusName string) error {
	resolvedPipelineID, err := s.resolvePipelineID(ctx, pipelineID, pipelineName)
	if err != nil {
		return err
	}

	resolvedStatusID, err := s.resolveStatusID(ctx, resolvedPipelineID, statusID, statusName)
	if err != nil {
		return err
	}

	return s.sdk.Statuses(resolvedPipelineID).DeleteOne(ctx, resolvedStatusID)
}

// --- helpers ---

// statusTypeLabel возвращает строковую метку типа статуса.
func statusTypeLabel(st *amomodels.Status) string {
	if st.IsWon() {
		return "won"
	}
	if st.IsLost() {
		return "lost"
	}
	return "regular"
}

// statusToOutput конвертирует SDK-модель статуса в LLM-friendly вывод.
// Убирает account_id и _links, добавляет семантические метки типа.
func statusToOutput(st *amomodels.Status) *toolmodels.StatusOutput {
	return &toolmodels.StatusOutput{
		ID:         st.ID,
		Name:       st.Name,
		Sort:       st.Sort,
		Color:      st.Color,
		PipelineID: st.PipelineID,
		IsEditable: st.IsEditable,
		TypeLabel:  statusTypeLabel(st),
		IsWon:      st.IsWon(),
		IsLost:     st.IsLost(),
		IsClosed:   st.IsClosed(),
	}
}

// statusTypeToInt конвертирует строковый тип ("regular"/"won"/"lost") в числовой для SDK.
func statusTypeToInt(t string) int {
	switch t {
	case "won":
		return int(amomodels.StatusTypeWon)
	case "lost":
		return int(amomodels.StatusTypeLost)
	default:
		return int(amomodels.StatusTypeRegular)
	}
}

// statusDataToModel конвертирует типизированные данные от LLM в SDK-модель.
func statusDataToModel(d toolmodels.StatusData) *amomodels.Status {
	st := &amomodels.Status{
		Name:  d.Name,
		Sort:  d.Sort,
		Color: d.Color,
	}
	if d.Type != "" {
		st.Type = statusTypeToInt(d.Type)
	}
	return st
}
