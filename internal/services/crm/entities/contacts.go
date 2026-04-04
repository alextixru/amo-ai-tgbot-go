package entities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) SearchContacts(ctx context.Context, filter *gkitmodels.EntitiesFilter, with []string) ([]*models.Contact, error) {
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

	contacts, _, err := s.sdk.Contacts().Get(ctx, f)
	return contacts, err
}

func (s *service) GetContact(ctx context.Context, id int, with []string) (*models.Contact, error) {
	f := filters.NewContactsFilter()
	f.SetIDs([]int{id})
	if len(with) > 0 {
		f.With = with
	}
	contacts, _, err := s.sdk.Contacts().Get(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(contacts) > 0 {
		return contacts[0], nil
	}
	return nil, nil
}

func (s *service) CreateContact(ctx context.Context, data *gkitmodels.EntityData) ([]*models.Contact, error) {
	contact := s.mapToContact(data)
	contacts, _, err := s.sdk.Contacts().Create(ctx, []*models.Contact{contact})
	return contacts, err
}

func (s *service) CreateContacts(ctx context.Context, dataList []gkitmodels.EntityData) ([]*models.Contact, error) {
	contacts := make([]*models.Contact, len(dataList))
	for i, data := range dataList {
		contacts[i] = s.mapToContact(&data)
	}
	result, _, err := s.sdk.Contacts().Create(ctx, contacts)
	return result, err
}

func (s *service) UpdateContact(ctx context.Context, id int, data *gkitmodels.EntityData) ([]*models.Contact, error) {
	contact := s.mapToContact(data)
	contact.ID = id
	contacts, _, err := s.sdk.Contacts().Update(ctx, []*models.Contact{contact})
	return contacts, err
}

func (s *service) UpdateContacts(ctx context.Context, dataList []gkitmodels.EntityData) ([]*models.Contact, error) {
	contacts := make([]*models.Contact, len(dataList))
	for i, data := range dataList {
		contacts[i] = s.mapToContact(&data)
		if contacts[i].ID == 0 && data.ID > 0 {
			contacts[i].ID = data.ID
		}
	}
	result, _, err := s.sdk.Contacts().Update(ctx, contacts)
	return result, err
}

func (s *service) SyncContact(ctx context.Context, id int, data *gkitmodels.EntityData) (*models.Contact, error) {
	contact := s.mapToContact(data)
	if id > 0 {
		contact.ID = id
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

func (s *service) mapToContact(data *gkitmodels.EntityData) *models.Contact {
	contact := &models.Contact{
		Name: data.Name,
	}
	if data.ID > 0 {
		contact.ID = data.ID
	}
	if data.ResponsibleUserID > 0 {
		contact.ResponsibleUserID = data.ResponsibleUserID
	}

	contact.CustomFieldsValues = mapCustomFieldsValues(data.CustomFieldsValues)

	if len(data.Tags) > 0 || len(data.EmbeddedCompanies) > 0 {
		contact.Embedded = &models.ContactEmbedded{}
		if len(data.Tags) > 0 {
			contact.Embedded.Tags = mapTags(data.Tags)
		}
		// Примечание: SDK ContactEmbedded может иметь другие поля, зависит от v4
	}

	return contact
}
