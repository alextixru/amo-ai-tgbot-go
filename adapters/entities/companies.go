package entities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) SearchCompanies(ctx context.Context, filter *gkitmodels.EntitiesFilter, with []string) ([]*models.Company, error) {
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
		if len(filter.IDs) > 0 {
			f.SetIDs(filter.IDs)
		}
		if len(filter.ResponsibleUserID) > 0 {
			f.SetResponsibleUserIDs(filter.ResponsibleUserID)
		}
		if len(filter.Names) > 0 {
			f.SetNames(filter.Names)
		}

		if len(filter.With) > 0 {
			with = append(with, filter.With...)
		}
	}

	if len(with) > 0 {
		f.With = with
	}

	companies, _, err := s.sdk.Companies().Get(ctx, f)
	return companies, err
}

func (s *service) GetCompany(ctx context.Context, id int, with []string) (*models.Company, error) {
	f := filters.NewCompaniesFilter()
	f.SetIDs([]int{id})
	if len(with) > 0 {
		f.With = with
	}
	companies, _, err := s.sdk.Companies().Get(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(companies) > 0 {
		return companies[0], nil
	}
	return nil, nil
}

func (s *service) CreateCompany(ctx context.Context, data *gkitmodels.EntityData) ([]*models.Company, error) {
	company := s.mapToCompany(data)
	companies, _, err := s.sdk.Companies().Create(ctx, []*models.Company{company})
	return companies, err
}

func (s *service) CreateCompanies(ctx context.Context, dataList []gkitmodels.EntityData) ([]*models.Company, error) {
	companies := make([]*models.Company, len(dataList))
	for i, data := range dataList {
		companies[i] = s.mapToCompany(&data)
	}
	result, _, err := s.sdk.Companies().Create(ctx, companies)
	return result, err
}

func (s *service) UpdateCompany(ctx context.Context, id int, data *gkitmodels.EntityData) ([]*models.Company, error) {
	company := s.mapToCompany(data)
	company.ID = id
	companies, _, err := s.sdk.Companies().Update(ctx, []*models.Company{company})
	return companies, err
}

func (s *service) UpdateCompanies(ctx context.Context, dataList []gkitmodels.EntityData) ([]*models.Company, error) {
	companies := make([]*models.Company, len(dataList))
	for i, data := range dataList {
		companies[i] = s.mapToCompany(&data)
		if companies[i].ID == 0 && data.ID > 0 {
			companies[i].ID = data.ID
		}
	}
	result, _, err := s.sdk.Companies().Update(ctx, companies)
	return result, err
}

func (s *service) SyncCompany(ctx context.Context, id int, data *gkitmodels.EntityData) (*models.Company, error) {
	company := s.mapToCompany(data)
	if id > 0 {
		company.ID = id
	}
	return s.sdk.Companies().SyncOne(ctx, company, []string{"leads", "contacts"})
}

func (s *service) LinkCompany(ctx context.Context, id int, target *gkitmodels.LinkTarget) error {
	return s.sdk.Companies().Link(ctx, id, target.Type, target.ID, nil)
}

func (s *service) UnlinkCompany(ctx context.Context, id int, target *gkitmodels.LinkTarget) error {
	return s.sdk.Companies().Unlink(ctx, id, target.Type, target.ID)
}

func (s *service) mapToCompany(data *gkitmodels.EntityData) *models.Company {
	company := &models.Company{
		Name: data.Name,
	}
	if data.ID > 0 {
		company.ID = data.ID
	}
	if data.ResponsibleUserID > 0 {
		company.ResponsibleUserID = data.ResponsibleUserID
	}

	company.CustomFieldsValues = mapCustomFieldsValues(data.CustomFieldsValues)

	if len(data.Tags) > 0 {
		company.Embedded = &models.CompanyEmbedded{}
		company.Embedded.Tags = mapTags(data.Tags)
	}

	return company
}
