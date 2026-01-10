package entities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) SearchCompanies(ctx context.Context, filter *gkitmodels.EntitiesFilter) ([]*models.Company, error) {
	f := filters.NewCompaniesFilter()
	f.SetLimit(50)
	f.SetPage(1)
	if filter != nil {
		if filter.Query != "" {
			f.SetQuery(filter.Query)
		}
		if filter.Limit > 0 {
			f.SetLimit(filter.Limit)
		}
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		}
		if len(filter.ResponsibleUserID) > 0 {
			f.SetResponsibleUserIDs(filter.ResponsibleUserID)
		}
	}
	companies, _, err := s.sdk.Companies().Get(ctx, f)
	return companies, err
}

func (s *service) GetCompany(ctx context.Context, id int) (*models.Company, error) {
	return s.sdk.Companies().GetOne(ctx, id)
}

func (s *service) CreateCompany(ctx context.Context, data *gkitmodels.EntityData) ([]*models.Company, error) {
	company := &models.Company{
		Name: data.Name,
	}
	if data.ResponsibleUserID > 0 {
		company.ResponsibleUserID = data.ResponsibleUserID
	}
	companies, _, err := s.sdk.Companies().Create(ctx, []*models.Company{company})
	return companies, err
}

func (s *service) UpdateCompany(ctx context.Context, id int, data *gkitmodels.EntityData) ([]*models.Company, error) {
	company := &models.Company{}
	company.ID = id
	if data.Name != "" {
		company.Name = data.Name
	}
	if data.ResponsibleUserID > 0 {
		company.ResponsibleUserID = data.ResponsibleUserID
	}
	companies, _, err := s.sdk.Companies().Update(ctx, []*models.Company{company})
	return companies, err
}

func (s *service) SyncCompany(ctx context.Context, id int, data *gkitmodels.EntityData) (*models.Company, error) {
	company := &models.Company{
		Name: data.Name,
	}
	if id > 0 {
		company.ID = id
	}
	if data.ResponsibleUserID > 0 {
		company.ResponsibleUserID = data.ResponsibleUserID
	}
	return s.sdk.Companies().SyncOne(ctx, company, []string{"leads", "contacts"})
}

func (s *service) LinkCompany(ctx context.Context, id int, target *gkitmodels.LinkTarget) error {
	return s.sdk.Companies().Link(ctx, id, target.Type, target.ID, nil)
}

func (s *service) UnlinkCompany(ctx context.Context, id int, target *gkitmodels.LinkTarget) error {
	return s.sdk.Companies().Unlink(ctx, id, target.Type, target.ID)
}
