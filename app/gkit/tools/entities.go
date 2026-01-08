package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// EntitiesInput входные параметры для инструмента entities
type EntitiesInput struct {
	// EntityType тип сущности: leads, contacts, companies
	EntityType string `json:"entity_type" jsonschema_description:"Тип сущности: leads, contacts, companies"`

	// Action действие: search, get, create, update, sync, delete, link, unlink, get_chats, link_chats
	Action string `json:"action" jsonschema_description:"Действие: search, get, create, update, sync, delete, link, unlink, get_chats, link_chats"`

	// ID идентификатор сущности (для get, update, delete, link, unlink)
	ID int `json:"id,omitempty" jsonschema_description:"ID сущности (для get, update, delete, link, unlink)"`

	// Filter параметры поиска (для search)
	Filter *EntitiesFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска (для action=search)"`

	// Data данные для создания/обновления
	Data *EntityData `json:"data,omitempty" jsonschema_description:"Данные сущности (для create, update)"`

	// LinkTo цель для связывания
	LinkTo *LinkTarget `json:"link_to,omitempty" jsonschema_description:"Цель связывания (для link, unlink)"`
}

// EntitiesFilter фильтры поиска сущностей
type EntitiesFilter struct {
	Query             string `json:"query,omitempty" jsonschema_description:"Поисковый запрос"`
	Limit             int    `json:"limit,omitempty" jsonschema_description:"Лимит результатов (макс 250)"`
	Page              int    `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	ResponsibleUserID []int  `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственных"`
	PipelineID        []int  `json:"pipeline_id,omitempty" jsonschema_description:"ID воронок (только для leads)"`
	StatusID          []int  `json:"status_id,omitempty" jsonschema_description:"ID статусов (только для leads)"`
}

// EntityData данные сущности для create/update
type EntityData struct {
	Name               string         `json:"name,omitempty" jsonschema_description:"Название"`
	Price              int            `json:"price,omitempty" jsonschema_description:"Бюджет (только для leads)"`
	StatusID           int            `json:"status_id,omitempty" jsonschema_description:"ID статуса (только для leads)"`
	PipelineID         int            `json:"pipeline_id,omitempty" jsonschema_description:"ID воронки (только для leads)"`
	ResponsibleUserID  int            `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственного"`
	CustomFieldsValues map[string]any `json:"custom_fields_values,omitempty" jsonschema_description:"Значения кастомных полей"`
}

// LinkTarget цель для связывания сущностей
type LinkTarget struct {
	Type string `json:"type" jsonschema_description:"Тип целевой сущности: leads, contacts, companies"`
	ID   int    `json:"id" jsonschema_description:"ID целевой сущности"`
}

// registerEntitiesTool регистрирует инструмент для работы с основными сущностями
func (r *Registry) registerEntitiesTool() {
	r.addTool(genkit.DefineTool[EntitiesInput, any](
		r.g,
		"entities",
		"Работа с основными сущностями amoCRM: сделки (leads), контакты (contacts), компании (companies). "+
			"Поддерживает: search (поиск), get (получение по ID), create (создание), update (обновление), "+
			"sync (создание или обновление), delete (удаление, только leads), link (связывание), unlink (отвязывание), "+
			"get_chats (получение чатов, только contacts), link_chats (привязка чатов, только contacts).",
		func(ctx *ai.ToolContext, input EntitiesInput) (any, error) {
			return r.handleEntities(ctx.Context, input)
		},
	))
}

func (r *Registry) handleEntities(ctx context.Context, input EntitiesInput) (any, error) {
	// Валидация entity_type
	switch input.EntityType {
	case "leads":
		return r.handleLeads(ctx, input)
	case "contacts":
		return r.handleContacts(ctx, input)
	case "companies":
		return r.handleCompanies(ctx, input)
	default:
		return nil, fmt.Errorf("unknown entity_type: %s (expected: leads, contacts, companies)", input.EntityType)
	}
}

// ============ LEADS ============

func (r *Registry) handleLeads(ctx context.Context, input EntitiesInput) (any, error) {
	switch input.Action {
	case "search":
		return r.searchLeads(ctx, input.Filter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.sdk.Leads().GetOne(ctx, input.ID)
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		return r.createLead(ctx, input.Data)
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		return r.updateLead(ctx, input.ID, input.Data)
	case "sync":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'sync'")
		}
		lead := &models.Lead{
			Name:       input.Data.Name,
			Price:      input.Data.Price,
			StatusID:   input.Data.StatusID,
			PipelineID: input.Data.PipelineID,
		}
		if input.ID > 0 {
			lead.ID = input.ID
		}
		if input.Data.ResponsibleUserID > 0 {
			lead.ResponsibleUserID = input.Data.ResponsibleUserID
		}
		return r.sdk.Leads().SyncOne(ctx, lead, []string{"contacts", "companies"})
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'delete'")
		}
		return nil, r.sdk.Leads().Delete(ctx, input.ID)
	case "link":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'link'")
		}
		return r.linkLead(ctx, input.ID, input.LinkTo)
	case "unlink":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'unlink'")
		}
		return nil, r.unlinkLead(ctx, input.ID, input.LinkTo)
	default:
		return nil, fmt.Errorf("unknown action for leads: %s", input.Action)
	}
}

func (r *Registry) searchLeads(ctx context.Context, filter *EntitiesFilter) ([]*models.Lead, error) {
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
		// StatusIDs не поддерживается напрямую — требуется SetStatuses с LeadStatusFilter
	}
	leads, _, err := r.sdk.Leads().Get(ctx, f)
	return leads, err
}

func (r *Registry) createLead(ctx context.Context, data *EntityData) (*models.Lead, error) {
	lead := &models.Lead{
		Name:       data.Name,
		Price:      data.Price,
		StatusID:   data.StatusID,
		PipelineID: data.PipelineID,
	}
	if data.ResponsibleUserID > 0 {
		lead.ResponsibleUserID = data.ResponsibleUserID
	}
	return r.sdk.Leads().CreateOne(ctx, lead)
}

func (r *Registry) updateLead(ctx context.Context, id int, data *EntityData) (*models.Lead, error) {
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
	return r.sdk.Leads().UpdateOne(ctx, lead)
}

func (r *Registry) linkLead(ctx context.Context, leadID int, target *LinkTarget) ([]models.EntityLink, error) {
	link := models.EntityLink{
		ToEntityID:   target.ID,
		ToEntityType: target.Type,
	}
	return r.sdk.Leads().Link(ctx, leadID, []models.EntityLink{link})
}

func (r *Registry) unlinkLead(ctx context.Context, leadID int, target *LinkTarget) error {
	link := models.EntityLink{
		ToEntityID:   target.ID,
		ToEntityType: target.Type,
	}
	return r.sdk.Leads().Unlink(ctx, leadID, []models.EntityLink{link})
}

// ============ CONTACTS ============

func (r *Registry) handleContacts(ctx context.Context, input EntitiesInput) (any, error) {
	switch input.Action {
	case "search":
		return r.searchContacts(ctx, input.Filter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.sdk.Contacts().GetOne(ctx, input.ID)
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		return r.createContact(ctx, input.Data)
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		return r.updateContact(ctx, input.ID, input.Data)
	case "sync":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'sync'")
		}
		contact := &models.Contact{
			Name: input.Data.Name,
		}
		if input.ID > 0 {
			contact.ID = input.ID
		}
		if input.Data.ResponsibleUserID > 0 {
			contact.ResponsibleUserID = input.Data.ResponsibleUserID
		}
		return r.sdk.Contacts().SyncOne(ctx, contact, []string{"leads", "companies"})
	case "get_chats":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get_chats'")
		}
		return r.sdk.Contacts().GetChats(ctx, input.ID)
	case "link_chats":
		// Требуется расширить EntitiesInput для поддержки ChatLinks
		return nil, fmt.Errorf("link_chats requires ChatLinks field in input (not yet implemented)")
	case "link":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'link'")
		}
		return nil, r.sdk.Contacts().Link(ctx, input.ID, input.LinkTo.Type, input.LinkTo.ID, nil)
	case "unlink":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'unlink'")
		}
		return nil, r.sdk.Contacts().Unlink(ctx, input.ID, input.LinkTo.Type, input.LinkTo.ID)
	default:
		return nil, fmt.Errorf("unknown action for contacts: %s (note: delete not supported)", input.Action)
	}
}

func (r *Registry) searchContacts(ctx context.Context, filter *EntitiesFilter) ([]*models.Contact, error) {
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
	contacts, _, err := r.sdk.Contacts().Get(ctx, f)
	return contacts, err
}

func (r *Registry) createContact(ctx context.Context, data *EntityData) ([]*models.Contact, error) {
	contact := &models.Contact{
		Name: data.Name,
	}
	if data.ResponsibleUserID > 0 {
		contact.ResponsibleUserID = data.ResponsibleUserID
	}
	contacts, _, err := r.sdk.Contacts().Create(ctx, []*models.Contact{contact})
	return contacts, err
}

func (r *Registry) updateContact(ctx context.Context, id int, data *EntityData) ([]*models.Contact, error) {
	contact := &models.Contact{}
	contact.ID = id
	if data.Name != "" {
		contact.Name = data.Name
	}
	if data.ResponsibleUserID > 0 {
		contact.ResponsibleUserID = data.ResponsibleUserID
	}
	contacts, _, err := r.sdk.Contacts().Update(ctx, []*models.Contact{contact})
	return contacts, err
}

// ============ COMPANIES ============

func (r *Registry) handleCompanies(ctx context.Context, input EntitiesInput) (any, error) {
	switch input.Action {
	case "search":
		return r.searchCompanies(ctx, input.Filter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.sdk.Companies().GetOne(ctx, input.ID)
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		return r.createCompany(ctx, input.Data)
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		return r.updateCompany(ctx, input.ID, input.Data)
	case "sync":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'sync'")
		}
		company := &models.Company{
			Name: input.Data.Name,
		}
		if input.ID > 0 {
			company.ID = input.ID
		}
		if input.Data.ResponsibleUserID > 0 {
			company.ResponsibleUserID = input.Data.ResponsibleUserID
		}
		return r.sdk.Companies().SyncOne(ctx, company, []string{"leads", "contacts"})
	case "link":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'link'")
		}
		return nil, r.sdk.Companies().Link(ctx, input.ID, input.LinkTo.Type, input.LinkTo.ID, nil)
	case "unlink":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'unlink'")
		}
		return nil, r.sdk.Companies().Unlink(ctx, input.ID, input.LinkTo.Type, input.LinkTo.ID)
	default:
		return nil, fmt.Errorf("unknown action for companies: %s (note: delete not supported)", input.Action)
	}
}

func (r *Registry) searchCompanies(ctx context.Context, filter *EntitiesFilter) ([]*models.Company, error) {
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
	companies, _, err := r.sdk.Companies().Get(ctx, f)
	return companies, err
}

func (r *Registry) createCompany(ctx context.Context, data *EntityData) ([]*models.Company, error) {
	company := &models.Company{
		Name: data.Name,
	}
	if data.ResponsibleUserID > 0 {
		company.ResponsibleUserID = data.ResponsibleUserID
	}
	companies, _, err := r.sdk.Companies().Create(ctx, []*models.Company{company})
	return companies, err
}

func (r *Registry) updateCompany(ctx context.Context, id int, data *EntityData) ([]*models.Company, error) {
	company := &models.Company{}
	company.ID = id
	if data.Name != "" {
		company.Name = data.Name
	}
	if data.ResponsibleUserID > 0 {
		company.ResponsibleUserID = data.ResponsibleUserID
	}
	companies, _, err := r.sdk.Companies().Update(ctx, []*models.Company{company})
	return companies, err
}
