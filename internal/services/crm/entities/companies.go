package entities

import (
	"context"
	"fmt"

	sdkfilters "github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) SearchCompanies(ctx context.Context, filter *gkitmodels.EntitiesFilter, with []string) (*SearchResult, error) {
	f := sdkfilters.NewCompaniesFilter()
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
		if len(filter.Names) > 0 {
			f.SetNames(filter.Names)
		}

		// Ответственные: имена → ID
		if len(filter.ResponsibleUserNames) > 0 {
			ids, err := s.resolveUserIDs(filter.ResponsibleUserNames)
			if err != nil {
				return nil, err
			}
			f.SetResponsibleUserIDs(ids)
		}

		// created_by: имена → ID
		if len(filter.CreatedByNames) > 0 {
			ids, err := s.resolveUserIDs(filter.CreatedByNames)
			if err != nil {
				return nil, err
			}
			f.SetCreatedBy(ids)
		}

		// updated_by: имена → ID
		if len(filter.UpdatedByNames) > 0 {
			ids, err := s.resolveUserIDs(filter.UpdatedByNames)
			if err != nil {
				return nil, err
			}
			f.SetUpdatedBy(ids)
		}

		// Даты (ISO-8601 → Unix)
		if filter.CreatedAtFrom != "" || filter.CreatedAtTo != "" {
			from := parseISO(filter.CreatedAtFrom)
			to := parseISO(filter.CreatedAtTo)
			f.SetCreatedAt(intPtrOrNil(from), intPtrOrNil(to))
		}
		if filter.UpdatedAtFrom != "" || filter.UpdatedAtTo != "" {
			from := parseISO(filter.UpdatedAtFrom)
			to := parseISO(filter.UpdatedAtTo)
			f.SetUpdatedAt(intPtrOrNil(from), intPtrOrNil(to))
		}

		// Кастомные поля: field_code → field_id
		if len(filter.CustomFieldsValues) > 0 {
			cfMap := buildCustomFieldsFilter(filter.CustomFieldsValues, s.customFieldsCompanies)
			if len(cfMap) > 0 {
				f.SetCustomFieldsValues(cfMap)
			}
		}

		if len(filter.With) > 0 {
			with = append(with, filter.With...)
		}
	}

	if len(with) > 0 {
		f.With = with
	}

	companies, meta, err := s.sdk.Companies().Get(ctx, f)
	if err != nil {
		return nil, err
	}

	items := make([]*EntityResult, 0, len(companies))
	for _, c := range companies {
		items = append(items, s.companyToResult(c))
	}

	result := &SearchResult{Items: items, Page: 1}
	if meta != nil {
		result.HasMore = meta.HasMore
		result.Page = meta.Page
	}
	return result, nil
}

func (s *service) GetCompany(ctx context.Context, id int, with []string) (*EntityResult, error) {
	f := sdkfilters.NewCompaniesFilter()
	f.SetIDs([]int{id})
	if len(with) > 0 {
		f.With = with
	}
	companies, _, err := s.sdk.Companies().Get(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(companies) == 0 {
		return nil, nil
	}
	return s.companyToResult(companies[0]), nil
}

func (s *service) CreateCompany(ctx context.Context, data *gkitmodels.EntityData) (*EntityResult, error) {
	company, err := s.mapToCompany(data)
	if err != nil {
		return nil, err
	}
	created, _, err := s.sdk.Companies().Create(ctx, []*models.Company{company})
	if err != nil {
		return nil, err
	}
	if len(created) == 0 {
		return nil, fmt.Errorf("компания не была создана")
	}
	return s.companyToResult(created[0]), nil
}

func (s *service) CreateCompanies(ctx context.Context, dataList []gkitmodels.EntityData) ([]*EntityResult, error) {
	companies := make([]*models.Company, 0, len(dataList))
	for i := range dataList {
		c, err := s.mapToCompany(&dataList[i])
		if err != nil {
			return nil, fmt.Errorf("item %d: %w", i, err)
		}
		companies = append(companies, c)
	}
	created, _, err := s.sdk.Companies().Create(ctx, companies)
	if err != nil {
		return nil, err
	}
	results := make([]*EntityResult, 0, len(created))
	for _, c := range created {
		results = append(results, s.companyToResult(c))
	}
	return results, nil
}

func (s *service) UpdateCompany(ctx context.Context, id int, data *gkitmodels.EntityData) (*EntityResult, error) {
	company, err := s.mapToCompany(data)
	if err != nil {
		return nil, err
	}
	company.ID = id
	updated, _, err := s.sdk.Companies().Update(ctx, []*models.Company{company})
	if err != nil {
		return nil, err
	}
	if len(updated) == 0 {
		return nil, fmt.Errorf("компания не была обновлена")
	}
	return s.companyToResult(updated[0]), nil
}

func (s *service) UpdateCompanies(ctx context.Context, dataList []gkitmodels.EntityData) ([]*EntityResult, error) {
	companies := make([]*models.Company, 0, len(dataList))
	for i := range dataList {
		c, err := s.mapToCompany(&dataList[i])
		if err != nil {
			return nil, fmt.Errorf("item %d: %w", i, err)
		}
		if c.ID == 0 && dataList[i].ID > 0 {
			c.ID = dataList[i].ID
		}
		companies = append(companies, c)
	}
	updated, _, err := s.sdk.Companies().Update(ctx, companies)
	if err != nil {
		return nil, err
	}
	results := make([]*EntityResult, 0, len(updated))
	for _, c := range updated {
		results = append(results, s.companyToResult(c))
	}
	return results, nil
}

func (s *service) SyncCompany(ctx context.Context, id int, data *gkitmodels.EntityData) (*EntityResult, error) {
	company, err := s.mapToCompany(data)
	if err != nil {
		return nil, err
	}
	if id > 0 {
		company.ID = id
	}
	synced, err := s.sdk.Companies().SyncOne(ctx, company, []string{"leads", "contacts"})
	if err != nil {
		return nil, err
	}
	return s.companyToResult(synced), nil
}

func (s *service) LinkCompany(ctx context.Context, id int, target *gkitmodels.LinkTarget) (*LinkResult, error) {
	if err := s.sdk.Companies().Link(ctx, id, target.Type, target.ID, nil); err != nil {
		return nil, err
	}
	return &LinkResult{
		Success: true,
		Message: fmt.Sprintf("Компания %d связана с %s %d", id, target.Type, target.ID),
	}, nil
}

func (s *service) UnlinkCompany(ctx context.Context, id int, target *gkitmodels.LinkTarget) (*LinkResult, error) {
	if err := s.sdk.Companies().Unlink(ctx, id, target.Type, target.ID); err != nil {
		return nil, err
	}
	return &LinkResult{
		Success: true,
		Message: fmt.Sprintf("Компания %d отвязана от %s %d", id, target.Type, target.ID),
	}, nil
}

// mapToCompany конвертирует EntityData в SDK Company, резолвя имена в ID.
func (s *service) mapToCompany(data *gkitmodels.EntityData) (*models.Company, error) {
	company := &models.Company{
		Name: data.Name,
	}
	if data.ID > 0 {
		company.ID = data.ID
	}

	// Ответственный — имя → ID
	if data.ResponsibleUserName != "" {
		id, err := s.resolveUserID(data.ResponsibleUserName)
		if err != nil {
			return nil, err
		}
		company.ResponsibleUserID = id
	}

	company.CustomFieldsValues = mapCustomFieldsValues(data.CustomFieldsValues)

	if len(data.Tags) > 0 {
		company.Embedded = &models.CompanyEmbedded{}
		company.Embedded.Tags = mapTags(data.Tags)
	}

	return company, nil
}
