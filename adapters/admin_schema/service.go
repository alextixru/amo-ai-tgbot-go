package admin_schema

import (
	"context"
	"net/url"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

// Service определяет бизнес-логику для работы со структурой данных CRM.
type Service interface {
	// Custom Fields
	ListCustomFields(ctx context.Context, entityType string, filter *filters.CustomFieldsFilter) ([]*models.CustomField, error)
	GetCustomField(ctx context.Context, entityType string, id int) (*models.CustomField, error)
	CreateCustomFields(ctx context.Context, entityType string, fields []*models.CustomField) ([]*models.CustomField, error)
	UpdateCustomFields(ctx context.Context, entityType string, fields []*models.CustomField) ([]*models.CustomField, error)
	DeleteCustomField(ctx context.Context, entityType string, id int) error

	// Field Groups
	ListFieldGroups(ctx context.Context, entityType string, filter url.Values) ([]models.CustomFieldGroup, error)
	GetFieldGroup(ctx context.Context, entityType string, id string) (*models.CustomFieldGroup, error)
	CreateFieldGroups(ctx context.Context, entityType string, groups []models.CustomFieldGroup) ([]models.CustomFieldGroup, error)
	UpdateFieldGroups(ctx context.Context, entityType string, groups []models.CustomFieldGroup) ([]models.CustomFieldGroup, error)
	DeleteFieldGroup(ctx context.Context, entityType string, id string) error

	// Loss Reasons
	ListLossReasons(ctx context.Context, filter url.Values) ([]*models.LossReason, error)
	GetLossReason(ctx context.Context, id int) (*models.LossReason, error)
	CreateLossReasons(ctx context.Context, reasons []*models.LossReason) ([]*models.LossReason, error)
	DeleteLossReason(ctx context.Context, id int) error

	// Sources
	ListSources(ctx context.Context, filter *filters.SourcesFilter) ([]*models.Source, error)
	GetSource(ctx context.Context, id int) (*models.Source, error)
	CreateSources(ctx context.Context, sources []*models.Source) ([]*models.Source, error)
	UpdateSources(ctx context.Context, sources []*models.Source) ([]*models.Source, error)
	DeleteSource(ctx context.Context, id int) error
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
