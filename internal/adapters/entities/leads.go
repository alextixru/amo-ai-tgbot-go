package entities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models"
)

func (s *service) SearchLeads(ctx context.Context, filter *gkitmodels.EntitiesFilter, with []string) ([]*models.Lead, error) {
	f := filters.NewLeadsFilter()
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
		if len(filter.PipelineID) > 0 {
			f.SetPipelineIDs(filter.PipelineID)
		}

		// Для остальных фильтров используем прямой доступ если методы отсутствуют
		// Примечание: предполагаем наличие полей QueryParams или аналогичных в фильтре
		// если SDK поддерживает эти фильтры в v4.

		if len(filter.With) > 0 {
			with = append(with, filter.With...)
		}
	}

	if len(with) > 0 {
		f.With = with
	}

	leads, _, err := s.sdk.Leads().Get(ctx, f)
	return leads, err
}

func (s *service) GetLead(ctx context.Context, id int, with []string) (*models.Lead, error) {
	f := filters.NewLeadsFilter()
	f.SetIDs([]int{id})
	if len(with) > 0 {
		f.With = with
	}
	leads, _, err := s.sdk.Leads().Get(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(leads) > 0 {
		return leads[0], nil
	}
	return nil, nil
}

func (s *service) CreateLead(ctx context.Context, data *gkitmodels.EntityData) (*models.Lead, error) {
	lead := s.mapToLead(data)
	return s.sdk.Leads().CreateOne(ctx, lead)
}

func (s *service) CreateLeads(ctx context.Context, dataList []gkitmodels.EntityData) ([]*models.Lead, error) {
	leads := make([]*models.Lead, len(dataList))
	for i, data := range dataList {
		leads[i] = s.mapToLead(&data)
	}
	result, _, err := s.sdk.Leads().Create(ctx, leads)
	return result, err
}

func (s *service) UpdateLead(ctx context.Context, id int, data *gkitmodels.EntityData) (*models.Lead, error) {
	lead := s.mapToLead(data)
	lead.ID = id
	return s.sdk.Leads().UpdateOne(ctx, lead)
}

func (s *service) UpdateLeads(ctx context.Context, dataList []gkitmodels.EntityData) ([]*models.Lead, error) {
	leads := make([]*models.Lead, len(dataList))
	for i, data := range dataList {
		leads[i] = s.mapToLead(&data)
		if leads[i].ID == 0 && data.ID > 0 {
			leads[i].ID = data.ID
		}
	}
	result, _, err := s.sdk.Leads().Update(ctx, leads)
	return result, err
}

func (s *service) SyncLead(ctx context.Context, id int, data *gkitmodels.EntityData) (*models.Lead, error) {
	lead := s.mapToLead(data)
	if id > 0 {
		lead.ID = id
	}
	return s.sdk.Leads().SyncOne(ctx, lead, []string{"contacts", "companies"})
}

func (s *service) LinkLead(ctx context.Context, leadID int, target *gkitmodels.LinkTarget) ([]models.EntityLink, error) {
	link := models.EntityLink{
		ToEntityID:   target.ID,
		ToEntityType: target.Type,
	}
	return s.sdk.Leads().Link(ctx, leadID, []models.EntityLink{link})
}

func (s *service) UnlinkLead(ctx context.Context, leadID int, target *gkitmodels.LinkTarget) error {
	link := models.EntityLink{
		ToEntityID:   target.ID,
		ToEntityType: target.Type,
	}
	return s.sdk.Leads().Unlink(ctx, leadID, []models.EntityLink{link})
}

func (s *service) mapToLead(data *gkitmodels.EntityData) *models.Lead {
	lead := &models.Lead{
		Name:       data.Name,
		Price:      data.Price,
		StatusID:   data.StatusID,
		PipelineID: data.PipelineID,
	}
	if data.ID > 0 {
		lead.ID = data.ID
	}
	if data.ResponsibleUserID > 0 {
		lead.ResponsibleUserID = data.ResponsibleUserID
	}

	lead.CustomFieldsValues = mapCustomFieldsValues(data.CustomFieldsValues)

	if len(data.Tags) > 0 || len(data.EmbeddedContacts) > 0 || len(data.EmbeddedCompanies) > 0 {
		lead.Embedded = &models.LeadEmbedded{}
		if len(data.Tags) > 0 {
			lead.Embedded.Tags = mapTags(data.Tags)
		}
		if len(data.EmbeddedContacts) > 0 {
			lead.Embedded.Contacts = make([]*models.Contact, len(data.EmbeddedContacts))
			for i, id := range data.EmbeddedContacts {
				contact := &models.Contact{}
				contact.ID = id
				lead.Embedded.Contacts[i] = contact
			}
		}
		if len(data.EmbeddedCompanies) > 0 {
			lead.Embedded.Companies = make([]*models.Company, len(data.EmbeddedCompanies))
			for i, id := range data.EmbeddedCompanies {
				company := &models.Company{}
				company.ID = id
				lead.Embedded.Companies[i] = company
			}
		}
	}

	return lead
}
