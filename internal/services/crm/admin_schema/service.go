package admin_schema

import (
	"context"
	"net/url"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
)

// PagedResult обёртка с результатами и метаданными пагинации
type PagedResult[T any] struct {
	Items      []T              `json:"items"`
	TotalItems int              `json:"total_items,omitempty"`
	Page       int              `json:"page,omitempty"`
	HasMore    bool             `json:"has_more"`
	Meta       *services.PageMeta `json:"-"`
}

// DeleteResult результат операции удаления
type DeleteResult struct {
	Success   bool   `json:"success"`
	DeletedID any    `json:"deleted_id"`
}

// Service определяет бизнес-логику для работы со структурой данных CRM.
type Service interface {
	// Custom Fields
	ListCustomFields(ctx context.Context, entityType string, filter *filters.CustomFieldsFilter) (*PagedResult[*models.CustomField], error)
	GetCustomField(ctx context.Context, entityType string, id int) (*models.CustomField, error)
	CreateCustomFields(ctx context.Context, entityType string, fields []*models.CustomField) (*PagedResult[*models.CustomField], error)
	UpdateCustomFields(ctx context.Context, entityType string, fields []*models.CustomField) (*PagedResult[*models.CustomField], error)
	DeleteCustomField(ctx context.Context, entityType string, id int) (*DeleteResult, error)

	// Field Groups
	ListFieldGroups(ctx context.Context, entityType string, filter url.Values) (*PagedResult[models.CustomFieldGroup], error)
	GetFieldGroup(ctx context.Context, entityType string, id string) (*models.CustomFieldGroup, error)
	CreateFieldGroups(ctx context.Context, entityType string, groups []models.CustomFieldGroup) (*PagedResult[models.CustomFieldGroup], error)
	UpdateFieldGroups(ctx context.Context, entityType string, groups []models.CustomFieldGroup) (*PagedResult[models.CustomFieldGroup], error)
	DeleteFieldGroup(ctx context.Context, entityType string, id string) (*DeleteResult, error)

	// Loss Reasons
	ListLossReasons(ctx context.Context, filter url.Values) (*PagedResult[*models.LossReason], error)
	GetLossReason(ctx context.Context, id int) (*models.LossReason, error)
	CreateLossReasons(ctx context.Context, reasons []*models.LossReason) (*PagedResult[*models.LossReason], error)
	DeleteLossReason(ctx context.Context, id int) (*DeleteResult, error)

	// Sources
	ListSources(ctx context.Context, filter *filters.SourcesFilter) (*PagedResult[*models.Source], error)
	GetSource(ctx context.Context, id int) (*models.Source, error)
	CreateSources(ctx context.Context, sources []*models.Source) (*PagedResult[*models.Source], error)
	UpdateSources(ctx context.Context, sources []*models.Source) (*PagedResult[*models.Source], error)
	DeleteSource(ctx context.Context, id int) (*DeleteResult, error)
}

type service struct {
	sdk *amocrm.SDK
}

// NewService создает новый экземпляр сервиса схемы.
func NewService(sdk *amocrm.SDK) Service {
	return &service{
		sdk: sdk,
	}
}

// newPagedResult создаёт PagedResult из списка и PageMeta SDK
func newPagedResult[T any](items []T, meta *services.PageMeta) *PagedResult[T] {
	r := &PagedResult[T]{
		Items: items,
		Meta:  meta,
	}
	if meta != nil {
		r.TotalItems = meta.TotalItems
		r.Page = meta.Page
		r.HasMore = meta.HasMore
	}
	if r.Items == nil {
		r.Items = []T{}
	}
	return r
}
