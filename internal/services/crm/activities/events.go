package activities

import (
	"context"
	"strings"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// defaultEventsWith — поля, которые запрашиваем у API по умолчанию,
// чтобы LLM получала имена контактов, сделок и компаний.
var defaultEventsWith = []string{"contact_name", "lead_name", "company_name"}

func (s *service) convertEvent(e *models.Event) *EventOutput {
	if e == nil {
		return nil
	}
	return &EventOutput{
		ID:            e.ID,
		Type:          e.Type,
		EntityID:      e.EntityID,
		EntityType:    e.EntityType,
		CreatedByName: s.resolveUserID(e.CreatedBy),
		CreatedAt:     toISO(e.CreatedAt),
		ValueBefore:   e.ValueBefore,
		ValueAfter:    e.ValueAfter,
	}
}

func (s *service) ListEvents(ctx context.Context, parent *gkitmodels.ParentEntity, filter *gkitmodels.EventsFilter) (*EventsListOutput, error) {
	f := filters.NewEventsFilter()
	if filter != nil {
		if filter.Limit > 0 {
			f.SetLimit(filter.Limit)
		} else {
			f.SetLimit(50)
		}
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		}
		if len(filter.Types) > 0 {
			f.SetTypes(filter.Types)
		}
		// Резолвинг created_by_names → IDs
		if len(filter.CreatedByNames) > 0 {
			ids, err := s.resolveUserNames(filter.CreatedByNames)
			if err != nil {
				return nil, err
			}
			f.SetCreatedBy(ids)
		}
	} else {
		f.SetLimit(50).SetPage(1)
	}

	if parent != nil {
		f.SetEntity([]string{parent.Type})
		f.SetEntityIDs([]int{parent.ID})
	}

	// Добавляем with в query params напрямую
	params := f.ToQueryParams()
	withValues := defaultEventsWith
	if filter != nil && len(filter.With) > 0 {
		withValues = filter.With
	}
	params.Set("with", strings.Join(withValues, ","))

	events, meta, err := s.sdk.Events().Get(ctx, params)
	if err != nil {
		return nil, err
	}

	out := &EventsListOutput{
		Events: make([]*EventOutput, 0, len(events)),
	}
	for _, e := range events {
		out.Events = append(out.Events, s.convertEvent(e))
	}
	if meta != nil {
		out.PageMeta = PageMeta{HasMore: meta.HasMore, Total: meta.TotalItems}
	}
	return out, nil
}

func (s *service) GetEvent(ctx context.Context, id int) (*EventOutput, error) {
	e, err := s.sdk.Events().GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.convertEvent(e), nil
}
