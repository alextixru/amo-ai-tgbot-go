package entities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) SearchContacts(ctx context.Context, filter *gkitmodels.EntitiesFilter) ([]*models.Contact, error) {
	f := filters.NewContactsFilter()
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
	contacts, _, err := s.sdk.Contacts().Get(ctx, f)
	return contacts, err
}

func (s *service) GetContact(ctx context.Context, id int) (*models.Contact, error) {
	return s.sdk.Contacts().GetOne(ctx, id)
}

func (s *service) CreateContact(ctx context.Context, data *gkitmodels.EntityData) ([]*models.Contact, error) {
	contact := &models.Contact{
		Name: data.Name,
	}
	if data.ResponsibleUserID > 0 {
		contact.ResponsibleUserID = data.ResponsibleUserID
	}
	contacts, _, err := s.sdk.Contacts().Create(ctx, []*models.Contact{contact})
	return contacts, err
}

func (s *service) UpdateContact(ctx context.Context, id int, data *gkitmodels.EntityData) ([]*models.Contact, error) {
	contact := &models.Contact{}
	contact.ID = id
	if data.Name != "" {
		contact.Name = data.Name
	}
	if data.ResponsibleUserID > 0 {
		contact.ResponsibleUserID = data.ResponsibleUserID
	}
	contacts, _, err := s.sdk.Contacts().Update(ctx, []*models.Contact{contact})
	return contacts, err
}

func (s *service) SyncContact(ctx context.Context, id int, data *gkitmodels.EntityData) (*models.Contact, error) {
	contact := &models.Contact{
		Name: data.Name,
	}
	if id > 0 {
		contact.ID = id
	}
	if data.ResponsibleUserID > 0 {
		contact.ResponsibleUserID = data.ResponsibleUserID
	}
	return s.sdk.Contacts().SyncOne(ctx, contact, []string{"leads", "companies"})
}

func (s *service) GetContactChats(ctx context.Context, id int) ([]models.ChatLink, error) {
	return s.sdk.Contacts().GetChats(ctx, id)
}

func (s *service) LinkContact(ctx context.Context, id int, target *gkitmodels.LinkTarget) error {
	return s.sdk.Contacts().Link(ctx, id, target.Type, target.ID, nil)
}

func (s *service) UnlinkContact(ctx context.Context, id int, target *gkitmodels.LinkTarget) error {
	return s.sdk.Contacts().Unlink(ctx, id, target.Type, target.ID)
}

func (s *service) LinkContactChats(ctx context.Context, links []models.ChatLink) ([]models.ChatLink, error) {
	return s.sdk.Contacts().LinkChats(ctx, links)
}
