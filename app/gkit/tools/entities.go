package tools

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	models "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// RegisterEntitiesTool регистрирует инструмент для работы с основными сущностями
func (r *Registry) RegisterEntitiesTool() {
	r.addTool(genkit.DefineTool[models.EntitiesInput, any](
		r.g,
		"entities",
		"Работа с основными сущностями amoCRM: сделки (leads), контакты (contacts), компании (companies). "+
			"Поддерживает: search (поиск), get (получение по ID), create (создание), update (обновление), "+
			"sync (создание или обновление), delete (удаление, только leads), link (связывание), unlink (отвязывание), "+
			"get_chats (получение чатов, только contacts).",
		func(ctx *ai.ToolContext, input models.EntitiesInput) (any, error) {
			return r.handleEntities(ctx.Context, input)
		},
	))
}

func (r *Registry) handleEntities(ctx context.Context, input models.EntitiesInput) (any, error) {
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

func (r *Registry) handleLeads(ctx context.Context, input models.EntitiesInput) (any, error) {
	switch input.Action {
	case "search":
		return r.entitiesService.SearchLeads(ctx, input.Filter, input.With)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.entitiesService.GetLead(ctx, input.ID, input.With)
	case "create":
		if len(input.DataList) > 0 {
			return r.entitiesService.CreateLeads(ctx, input.DataList)
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data or data_list is required for action 'create'")
		}
		return r.entitiesService.CreateLead(ctx, input.Data)
	case "update":
		if len(input.DataList) > 0 {
			return r.entitiesService.UpdateLeads(ctx, input.DataList)
		}
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data or data_list is required for action 'update'")
		}
		return r.entitiesService.UpdateLead(ctx, input.ID, input.Data)
	case "sync":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'sync'")
		}
		return r.entitiesService.SyncLead(ctx, input.ID, input.Data)
	case "link":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'link'")
		}
		return r.entitiesService.LinkLead(ctx, input.ID, input.LinkTo)
	case "unlink":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'unlink'")
		}
		return nil, r.entitiesService.UnlinkLead(ctx, input.ID, input.LinkTo)
	default:
		return nil, fmt.Errorf("unknown action for leads: %s (note: delete not supported)", input.Action)
	}
}

// ============ CONTACTS ============

func (r *Registry) handleContacts(ctx context.Context, input models.EntitiesInput) (any, error) {
	switch input.Action {
	case "search":
		return r.entitiesService.SearchContacts(ctx, input.Filter, input.With)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.entitiesService.GetContact(ctx, input.ID, input.With)
	case "create":
		if len(input.DataList) > 0 {
			return r.entitiesService.CreateContacts(ctx, input.DataList)
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data or data_list is required for action 'create'")
		}
		return r.entitiesService.CreateContact(ctx, input.Data)
	case "update":
		if len(input.DataList) > 0 {
			return r.entitiesService.UpdateContacts(ctx, input.DataList)
		}
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data or data_list is required for action 'update'")
		}
		return r.entitiesService.UpdateContact(ctx, input.ID, input.Data)
	case "sync":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'sync'")
		}
		return r.entitiesService.SyncContact(ctx, input.ID, input.Data)
	case "get_chats":
		// ... existing get_chats handling ...
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get_chats'")
		}
		return r.entitiesService.GetContactChats(ctx, input.ID)
	case "link_chats":
		if len(input.ChatLinks) == 0 {
			return nil, fmt.Errorf("chat_links is required for action 'link_chats'")
		}
		return r.entitiesService.LinkContactChats(ctx, input.ChatLinks)
	case "link":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'link'")
		}
		return nil, r.entitiesService.LinkContact(ctx, input.ID, input.LinkTo)
	case "unlink":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'unlink'")
		}
		return nil, r.entitiesService.UnlinkContact(ctx, input.ID, input.LinkTo)
	default:
		return nil, fmt.Errorf("unknown action for contacts: %s (note: delete not supported)", input.Action)
	}
}

// ============ COMPANIES ============

func (r *Registry) handleCompanies(ctx context.Context, input models.EntitiesInput) (any, error) {
	switch input.Action {
	case "search":
		return r.entitiesService.SearchCompanies(ctx, input.Filter, input.With)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.entitiesService.GetCompany(ctx, input.ID, input.With)
	case "create":
		if len(input.DataList) > 0 {
			return r.entitiesService.CreateCompanies(ctx, input.DataList)
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data or data_list is required for action 'create'")
		}
		return r.entitiesService.CreateCompany(ctx, input.Data)
	case "update":
		if len(input.DataList) > 0 {
			return r.entitiesService.UpdateCompanies(ctx, input.DataList)
		}
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data or data_list is required for action 'update'")
		}
		return r.entitiesService.UpdateCompany(ctx, input.ID, input.Data)
	case "sync":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'sync'")
		}
		return r.entitiesService.SyncCompany(ctx, input.ID, input.Data)
	case "link":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'link'")
		}
		return nil, r.entitiesService.LinkCompany(ctx, input.ID, input.LinkTo)
	case "unlink":
		if input.ID == 0 || input.LinkTo == nil {
			return nil, fmt.Errorf("id and link_to are required for action 'unlink'")
		}
		return nil, r.entitiesService.UnlinkCompany(ctx, input.ID, input.LinkTo)
	default:
		return nil, fmt.Errorf("unknown action for companies: %s (note: delete not supported)", input.Action)
	}
}
