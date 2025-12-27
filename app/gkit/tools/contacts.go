package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// ContactToolInput входные параметры для инструмента crm_contacts
type ContactToolInput struct {
	// Action действие: search, get, create, update
	Action string `json:"action"`

	// Filter фильтр для поиска (используется при action=search)
	Filter *ContactFilterInput `json:"filter,omitempty"`

	// ID идентификатор контакта (используется при action=get, update)
	ID int `json:"id,omitempty"`

	// Data данные для создания или обновления (используется при action=create, update)
	Data *ContactDataInput `json:"data,omitempty"`
}

// ContactFilterInput параметры фильтрации контактов
type ContactFilterInput struct {
	// Query поисковый запрос (имя, телефон, email)
	Query string `json:"query,omitempty"`
	// ResponsibleUserID ID ответственного
	ResponsibleUserID int `json:"responsible_user_id,omitempty"`
	// Limit количество элементов (по умолчанию 50)
	Limit int `json:"limit,omitempty"`
}

// ContactDataInput данные контакты
type ContactDataInput struct {
	// Name имя контакта (обязательно при создании)
	Name string `json:"name,omitempty"`
	// FirstName имя (если нужно разделить)
	FirstName string `json:"first_name,omitempty"`
	// LastName фамилия
	LastName string `json:"last_name,omitempty"`
	// ResponsibleUserID ID ответственного
	ResponsibleUserID int `json:"responsible_user_id,omitempty"`
}

// registerContactsTools регистрирует инструменты для работы с контактами
func (r *Registry) registerContactsTools() {
	r.addTool(genkit.DefineTool[ContactToolInput, any](
		r.g,
		"crm_contacts",
		"Управление контактами (Contacts). Позволяет искать, получать, создавать и обновлять контакты в amoCRM.",
		func(ctx *ai.ToolContext, input ContactToolInput) (any, error) {
			switch input.Action {
			case "search":
				return r.searchContacts(ctx.Context, input.Filter)
			case "get":
				if input.ID == 0 {
					return nil, fmt.Errorf("id is required for 'get' action")
				}
				return r.getContact(ctx.Context, input.ID)
			case "create":
				if input.Data == nil {
					return nil, fmt.Errorf("data is required for 'create' action")
				}
				return r.createContact(ctx.Context, input.Data)
			case "update":
				if input.ID == 0 {
					return nil, fmt.Errorf("id is required for 'update' action")
				}
				if input.Data == nil {
					return nil, fmt.Errorf("data is required for 'update' action")
				}
				return r.updateContact(ctx.Context, input.ID, input.Data)
			default:
				return nil, fmt.Errorf("unknown action: %s", input.Action)
			}
		},
	))
}

// searchContacts выполняет поиск контактов
func (r *Registry) searchContacts(ctx context.Context, filter *ContactFilterInput) ([]models.Contact, error) {
	sdkFilter := &services.ContactsFilter{
		Limit: 50,
	}

	if filter != nil {
		if filter.Limit > 0 {
			sdkFilter.Limit = filter.Limit
		}
		if filter.Query != "" {
			sdkFilter.Query = filter.Query
		}
		if filter.ResponsibleUserID > 0 {
			sdkFilter.FilterByResponsibleUserID = []int{filter.ResponsibleUserID}
		}
	}

	return r.sdk.Contacts().Get(ctx, sdkFilter)
}

// getContact получает контакт по ID
func (r *Registry) getContact(ctx context.Context, id int) (*models.Contact, error) {
	// Добавляем with=leads, чтобы видеть сделки
	return r.sdk.Contacts().GetOne(ctx, id, []string{"leads"})
}

// createContact создает новый контакт
func (r *Registry) createContact(ctx context.Context, data *ContactDataInput) (*models.Contact, error) {
	contact := models.Contact{
		Name:      data.Name,
		FirstName: data.FirstName,
		LastName:  data.LastName,
	}

	if data.ResponsibleUserID > 0 {
		contact.ResponsibleUserID = data.ResponsibleUserID
	}

	createdContacts, err := r.sdk.Contacts().Create(ctx, []models.Contact{contact})
	if err != nil {
		return nil, err
	}

	if len(createdContacts) == 0 {
		return nil, fmt.Errorf("failed to create contact: empty response")
	}

	return &createdContacts[0], nil
}

// updateContact обновляет существующий контакт
func (r *Registry) updateContact(ctx context.Context, id int, data *ContactDataInput) (*models.Contact, error) {
	contact := models.Contact{
		BaseModel: models.BaseModel{ID: id},
	}

	if data.Name != "" {
		contact.Name = data.Name
	}
	if data.FirstName != "" {
		contact.FirstName = data.FirstName
	}
	if data.LastName != "" {
		contact.LastName = data.LastName
	}
	if data.ResponsibleUserID != 0 {
		contact.ResponsibleUserID = data.ResponsibleUserID
	}

	updatedContacts, err := r.sdk.Contacts().Update(ctx, []models.Contact{contact})
	if err != nil {
		return nil, err
	}

	if len(updatedContacts) == 0 {
		return nil, fmt.Errorf("failed to update contact: empty response")
	}

	return &updatedContacts[0], nil
}
