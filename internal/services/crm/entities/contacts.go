package entities

import (
	"context"
	"fmt"

	sdkfilters "github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) SearchContacts(ctx context.Context, filter *gkitmodels.EntitiesFilter, with []string) (*SearchResult, error) {
	f := sdkfilters.NewContactsFilter()
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
			cfMap := buildCustomFieldsFilter(filter.CustomFieldsValues, s.customFieldsContacts)
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

	contacts, meta, err := s.sdk.Contacts().Get(ctx, f)
	if err != nil {
		return nil, err
	}

	items := make([]*EntityResult, 0, len(contacts))
	for _, c := range contacts {
		items = append(items, s.contactToResult(c))
	}

	result := &SearchResult{Items: items, Page: 1}
	if meta != nil {
		result.HasMore = meta.HasMore
		result.Page = meta.Page
	}
	return result, nil
}

func (s *service) GetContact(ctx context.Context, id int, with []string) (*EntityResult, error) {
	f := sdkfilters.NewContactsFilter()
	f.SetIDs([]int{id})
	if len(with) > 0 {
		f.With = with
	}
	contacts, _, err := s.sdk.Contacts().Get(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(contacts) == 0 {
		return nil, nil
	}
	return s.contactToResult(contacts[0]), nil
}

func (s *service) CreateContact(ctx context.Context, data *gkitmodels.EntityData) (*EntityResult, error) {
	contact, err := s.mapToContact(data)
	if err != nil {
		return nil, err
	}
	created, _, err := s.sdk.Contacts().Create(ctx, []*models.Contact{contact})
	if err != nil {
		return nil, err
	}
	if len(created) == 0 {
		return nil, fmt.Errorf("контакт не был создан")
	}
	return s.contactToResult(created[0]), nil
}

func (s *service) CreateContacts(ctx context.Context, dataList []gkitmodels.EntityData) ([]*EntityResult, error) {
	contacts := make([]*models.Contact, 0, len(dataList))
	for i := range dataList {
		c, err := s.mapToContact(&dataList[i])
		if err != nil {
			return nil, fmt.Errorf("item %d: %w", i, err)
		}
		contacts = append(contacts, c)
	}
	created, _, err := s.sdk.Contacts().Create(ctx, contacts)
	if err != nil {
		return nil, err
	}
	results := make([]*EntityResult, 0, len(created))
	for _, c := range created {
		results = append(results, s.contactToResult(c))
	}
	return results, nil
}

func (s *service) UpdateContact(ctx context.Context, id int, data *gkitmodels.EntityData) (*EntityResult, error) {
	contact, err := s.mapToContact(data)
	if err != nil {
		return nil, err
	}
	contact.ID = id
	updated, _, err := s.sdk.Contacts().Update(ctx, []*models.Contact{contact})
	if err != nil {
		return nil, err
	}
	if len(updated) == 0 {
		return nil, fmt.Errorf("контакт не был обновлён")
	}
	return s.contactToResult(updated[0]), nil
}

func (s *service) UpdateContacts(ctx context.Context, dataList []gkitmodels.EntityData) ([]*EntityResult, error) {
	contacts := make([]*models.Contact, 0, len(dataList))
	for i := range dataList {
		c, err := s.mapToContact(&dataList[i])
		if err != nil {
			return nil, fmt.Errorf("item %d: %w", i, err)
		}
		if c.ID == 0 && dataList[i].ID > 0 {
			c.ID = dataList[i].ID
		}
		contacts = append(contacts, c)
	}
	updated, _, err := s.sdk.Contacts().Update(ctx, contacts)
	if err != nil {
		return nil, err
	}
	results := make([]*EntityResult, 0, len(updated))
	for _, c := range updated {
		results = append(results, s.contactToResult(c))
	}
	return results, nil
}

func (s *service) SyncContact(ctx context.Context, id int, data *gkitmodels.EntityData) (*EntityResult, error) {
	contact, err := s.mapToContact(data)
	if err != nil {
		return nil, err
	}
	if id > 0 {
		contact.ID = id
	}
	synced, err := s.sdk.Contacts().SyncOne(ctx, contact, []string{"leads", "companies"})
	if err != nil {
		return nil, err
	}
	return s.contactToResult(synced), nil
}

func (s *service) GetContactChats(ctx context.Context, id int) (any, error) {
	return s.sdk.Contacts().GetChats(ctx, id)
}

func (s *service) LinkContact(ctx context.Context, id int, target *gkitmodels.LinkTarget) (*LinkResult, error) {
	if err := s.sdk.Contacts().Link(ctx, id, target.Type, target.ID, nil); err != nil {
		return nil, err
	}
	return &LinkResult{
		Success: true,
		Message: fmt.Sprintf("Контакт %d связан с %s %d", id, target.Type, target.ID),
	}, nil
}

func (s *service) UnlinkContact(ctx context.Context, id int, target *gkitmodels.LinkTarget) (*LinkResult, error) {
	if err := s.sdk.Contacts().Unlink(ctx, id, target.Type, target.ID); err != nil {
		return nil, err
	}
	return &LinkResult{
		Success: true,
		Message: fmt.Sprintf("Контакт %d отвязан от %s %d", id, target.Type, target.ID),
	}, nil
}

func (s *service) LinkContactChats(ctx context.Context, links any) (any, error) {
	chatLinks, ok := links.([]models.ChatLink)
	if !ok {
		return nil, fmt.Errorf("неверный тип для chat_links")
	}
	return s.sdk.Contacts().LinkChats(ctx, chatLinks)
}

// mapToContact конвертирует EntityData в SDK Contact, резолвя имена в ID.
func (s *service) mapToContact(data *gkitmodels.EntityData) (*models.Contact, error) {
	contact := &models.Contact{
		Name:      data.Name,
		FirstName: data.FirstName,
		LastName:  data.LastName,
	}
	if data.ID > 0 {
		contact.ID = data.ID
	}

	// Ответственный — имя → ID
	if data.ResponsibleUserName != "" {
		id, err := s.resolveUserID(data.ResponsibleUserName)
		if err != nil {
			return nil, err
		}
		contact.ResponsibleUserID = id
	}

	contact.CustomFieldsValues = mapCustomFieldsValues(data.CustomFieldsValues)

	hasEmbedded := len(data.Tags) > 0 || len(data.EmbeddedCompanies) > 0
	if hasEmbedded {
		contact.Embedded = &models.ContactEmbedded{}
		if len(data.Tags) > 0 {
			contact.Embedded.Tags = mapTags(data.Tags)
		}
		// Привязка компаний к контакту через embedded
		if len(data.EmbeddedCompanies) > 0 {
			companies := make([]*models.Company, len(data.EmbeddedCompanies))
			for i, cid := range data.EmbeddedCompanies {
				companies[i] = &models.Company{}
				companies[i].ID = cid
			}
			contact.Embedded.Companies = companies
		}
	}

	return contact, nil
}
