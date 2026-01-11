package entities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models"
)

// Service определяет бизнес-логику для работы с основными сущностями amoCRM.
type Service interface {
	// Leads
	SearchLeads(ctx context.Context, filter *gkitmodels.EntitiesFilter, with []string) ([]*models.Lead, error)
	GetLead(ctx context.Context, id int, with []string) (*models.Lead, error)
	CreateLead(ctx context.Context, data *gkitmodels.EntityData) (*models.Lead, error)
	CreateLeads(ctx context.Context, dataList []gkitmodels.EntityData) ([]*models.Lead, error)
	UpdateLead(ctx context.Context, id int, data *gkitmodels.EntityData) (*models.Lead, error)
	UpdateLeads(ctx context.Context, dataList []gkitmodels.EntityData) ([]*models.Lead, error)
	SyncLead(ctx context.Context, id int, data *gkitmodels.EntityData) (*models.Lead, error)
	LinkLead(ctx context.Context, id int, target *gkitmodels.LinkTarget) ([]models.EntityLink, error)
	UnlinkLead(ctx context.Context, id int, target *gkitmodels.LinkTarget) error

	// Contacts
	SearchContacts(ctx context.Context, filter *gkitmodels.EntitiesFilter, with []string) ([]*models.Contact, error)
	GetContact(ctx context.Context, id int, with []string) (*models.Contact, error)
	CreateContact(ctx context.Context, data *gkitmodels.EntityData) ([]*models.Contact, error)
	CreateContacts(ctx context.Context, dataList []gkitmodels.EntityData) ([]*models.Contact, error)
	UpdateContact(ctx context.Context, id int, data *gkitmodels.EntityData) ([]*models.Contact, error)
	UpdateContacts(ctx context.Context, dataList []gkitmodels.EntityData) ([]*models.Contact, error)
	SyncContact(ctx context.Context, id int, data *gkitmodels.EntityData) (*models.Contact, error)
	LinkContact(ctx context.Context, id int, target *gkitmodels.LinkTarget) error
	UnlinkContact(ctx context.Context, id int, target *gkitmodels.LinkTarget) error
	GetContactChats(ctx context.Context, id int) ([]models.ChatLink, error)
	LinkContactChats(ctx context.Context, links []models.ChatLink) ([]models.ChatLink, error)

	// Companies
	SearchCompanies(ctx context.Context, filter *gkitmodels.EntitiesFilter, with []string) ([]*models.Company, error)
	GetCompany(ctx context.Context, id int, with []string) (*models.Company, error)
	CreateCompany(ctx context.Context, data *gkitmodels.EntityData) ([]*models.Company, error)
	CreateCompanies(ctx context.Context, dataList []gkitmodels.EntityData) ([]*models.Company, error)
	UpdateCompany(ctx context.Context, id int, data *gkitmodels.EntityData) ([]*models.Company, error)
	UpdateCompanies(ctx context.Context, dataList []gkitmodels.EntityData) ([]*models.Company, error)
	SyncCompany(ctx context.Context, id int, data *gkitmodels.EntityData) (*models.Company, error)
	LinkCompany(ctx context.Context, id int, target *gkitmodels.LinkTarget) error
	UnlinkCompany(ctx context.Context, id int, target *gkitmodels.LinkTarget) error
}

type service struct {
	sdk *amocrm.SDK
}

// New создает новый экземпляр сервиса сущностей.
func New(sdk *amocrm.SDK) Service {
	return &service{
		sdk: sdk,
	}
}
