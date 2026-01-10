package entities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) SearchLeads(ctx context.Context, filter *gkitmodels.EntitiesFilter) ([]*models.Lead, error) {
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
		if len(filter.ResponsibleUserID) > 0 {
			f.SetResponsibleUserIDs(filter.ResponsibleUserID)
		}
		if len(filter.PipelineID) > 0 {
			f.SetPipelineIDs(filter.PipelineID)
		}
	}
	leads, _, err := s.sdk.Leads().Get(ctx, f)
	return leads, err
}

func (s *service) GetLead(ctx context.Context, id int) (*models.Lead, error) {
	return s.sdk.Leads().GetOne(ctx, id)
}

func (s *service) CreateLead(ctx context.Context, data *gkitmodels.EntityData) (*models.Lead, error) {
	lead := &models.Lead{
		Name:       data.Name,
		Price:      data.Price,
		StatusID:   data.StatusID,
		PipelineID: data.PipelineID,
	}
	if data.ResponsibleUserID > 0 {
		lead.ResponsibleUserID = data.ResponsibleUserID
	}
	return s.sdk.Leads().CreateOne(ctx, lead)
}

func (s *service) UpdateLead(ctx context.Context, id int, data *gkitmodels.EntityData) (*models.Lead, error) {
	lead := &models.Lead{}
	lead.ID = id
	if data.Name != "" {
		lead.Name = data.Name
	}
	if data.Price > 0 {
		lead.Price = data.Price
	}
	if data.StatusID > 0 {
		lead.StatusID = data.StatusID
	}
	if data.PipelineID > 0 {
		lead.PipelineID = data.PipelineID
	}
	if data.ResponsibleUserID > 0 {
		lead.ResponsibleUserID = data.ResponsibleUserID
	}
	return s.sdk.Leads().UpdateOne(ctx, lead)
}

func (s *service) SyncLead(ctx context.Context, id int, data *gkitmodels.EntityData) (*models.Lead, error) {
	lead := &models.Lead{
		Name:       data.Name,
		Price:      data.Price,
		StatusID:   data.StatusID,
		PipelineID: data.PipelineID,
	}
	if id > 0 {
		lead.ID = id
	}
	if data.ResponsibleUserID > 0 {
		lead.ResponsibleUserID = data.ResponsibleUserID
	}
	return s.sdk.Leads().SyncOne(ctx, lead, []string{"contacts", "companies"})
}

func (s *service) DeleteLead(ctx context.Context, id int) error {
	return s.sdk.Leads().Delete(ctx, id)
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
