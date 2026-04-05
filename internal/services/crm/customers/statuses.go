package customers

import (
	"context"
	"net/url"
	"strconv"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListCustomerStatuses(ctx context.Context, page, limit int) ([]StatusOutput, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	statuses, _, err := s.sdk.CustomerStatuses().Get(ctx, params)
	if err != nil {
		return nil, err
	}
	out := make([]StatusOutput, 0, len(statuses))
	for _, st := range statuses {
		out = append(out, statusToOutput(st))
	}
	return out, nil
}

func (s *service) GetCustomerStatus(ctx context.Context, id int) (*StatusOutput, error) {
	st, err := s.sdk.CustomerStatuses().GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	o := statusToOutput(st)
	return &o, nil
}

func (s *service) CreateCustomerStatuses(ctx context.Context, names []string) ([]StatusOutput, error) {
	statuses := make([]models.Status, 0, len(names))
	for _, name := range names {
		statuses = append(statuses, models.Status{Name: name})
	}
	res, _, err := s.sdk.CustomerStatuses().Create(ctx, statuses)
	if err != nil {
		return nil, err
	}
	out := make([]StatusOutput, 0, len(res))
	for _, st := range res {
		out = append(out, statusToOutput(st))
	}
	return out, nil
}

func (s *service) UpdateCustomerStatus(ctx context.Context, id int, name string) (*StatusOutput, error) {
	st := models.Status{Name: name}
	st.ID = id
	res, _, err := s.sdk.CustomerStatuses().Update(ctx, []models.Status{st})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	o := statusToOutput(res[0])
	return &o, nil
}

func (s *service) DeleteCustomerStatus(ctx context.Context, id int) error {
	return s.sdk.CustomerStatuses().Delete(ctx, id)
}

func statusToOutput(st models.Status) StatusOutput {
	return StatusOutput{
		ID:    st.ID,
		Name:  st.Name,
		Color: st.Color,
		Sort:  st.Sort,
	}
}
