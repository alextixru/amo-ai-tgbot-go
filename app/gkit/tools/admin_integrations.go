package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// AdminIntegrationsInput входные параметры для инструмента admin_integrations
type AdminIntegrationsInput struct {
	// Layer слой: webhooks | widgets | website_buttons | chat_templates | short_links
	Layer string `json:"layer" jsonschema_description:"Слой: webhooks, widgets, website_buttons, chat_templates, short_links"`

	// Action действие
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, delete, subscribe, unsubscribe, install, uninstall"`

	// ID идентификатор (для get, update, delete)
	ID int `json:"id,omitempty" jsonschema_description:"ID элемента"`

	// Code код виджета (для widgets)
	Code string `json:"code,omitempty" jsonschema_description:"Код виджета (для widgets: get, install, uninstall)"`

	// Filter фильтры для list
	Filter *IntegrationsFilter `json:"filter,omitempty" jsonschema_description:"Фильтры для поиска"`

	// Data данные для create/update
	Data map[string]any `json:"data,omitempty" jsonschema_description:"Данные для создания/обновления"`
}

// IntegrationsFilter фильтры для admin_integrations
type IntegrationsFilter struct {
	Limit int `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Page  int `json:"page,omitempty" jsonschema_description:"Номер страницы"`
}

// registerAdminIntegrationsTool регистрирует инструмент для управления интеграциями
func (r *Registry) registerAdminIntegrationsTool() {
	r.addTool(genkit.DefineTool[AdminIntegrationsInput, any](
		r.g,
		"admin_integrations",
		"Управление интеграциями amoCRM. "+
			"Layers: webhooks, widgets, website_buttons, chat_templates, short_links. "+
			"Webhooks: list, subscribe, unsubscribe. "+
			"Widgets: list, get, install, uninstall (используют code). "+
			"WebsiteButtons: list, get (read-only). "+
			"ChatTemplates: list, get, create, update, delete. "+
			"ShortLinks: list, create, delete.",
		func(ctx *ai.ToolContext, input AdminIntegrationsInput) (any, error) {
			return r.handleAdminIntegrations(ctx.Context, input)
		},
	))
}

func (r *Registry) handleAdminIntegrations(ctx context.Context, input AdminIntegrationsInput) (any, error) {
	switch input.Layer {
	case "webhooks":
		return r.handleWebhooks(ctx, input)
	case "widgets":
		return r.handleWidgets(ctx, input)
	case "website_buttons":
		return r.handleWebsiteButtons(ctx, input)
	case "chat_templates":
		return r.handleChatTemplates(ctx, input)
	case "short_links":
		return r.handleShortLinks(ctx, input)
	default:
		return nil, fmt.Errorf("unknown layer: %s", input.Layer)
	}
}

// ============================================================================
// Webhooks
// ============================================================================

func (r *Registry) handleWebhooks(ctx context.Context, input AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "list":
		return r.listWebhooks(ctx, input.Filter)
	case "subscribe":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'subscribe'")
		}
		return r.subscribeWebhook(ctx, input.Data)
	case "unsubscribe":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'unsubscribe'")
		}
		return r.unsubscribeWebhook(ctx, input.Data)
	default:
		return nil, fmt.Errorf("unknown action: %s for webhooks (available: list, subscribe, unsubscribe)", input.Action)
	}
}

func (r *Registry) listWebhooks(ctx context.Context, filter *IntegrationsFilter) ([]models.Webhook, error) {
	f := &services.WebhooksFilter{
		Limit: 50,
		Page:  1,
	}
	if filter != nil {
		if filter.Limit > 0 {
			f.Limit = filter.Limit
		}
		if filter.Page > 0 {
			f.Page = filter.Page
		}
	}
	return r.sdk.Webhooks().Get(ctx, f)
}

func (r *Registry) subscribeWebhook(ctx context.Context, data map[string]any) (*models.Webhook, error) {
	webhook := &models.Webhook{}

	if destination, ok := data["destination"].(string); ok {
		webhook.Destination = destination
	}
	if settings, ok := data["settings"].([]any); ok {
		for _, s := range settings {
			if str, ok := s.(string); ok {
				webhook.Settings = append(webhook.Settings, str)
			}
		}
	}

	return r.sdk.Webhooks().Subscribe(ctx, webhook)
}

func (r *Registry) unsubscribeWebhook(ctx context.Context, data map[string]any) (any, error) {
	webhook := &models.Webhook{}

	if destination, ok := data["destination"].(string); ok {
		webhook.Destination = destination
	}

	err := r.sdk.Webhooks().Unsubscribe(ctx, webhook)
	if err != nil {
		return nil, err
	}
	return map[string]any{"success": true, "unsubscribed": webhook.Destination}, nil
}

// ============================================================================
// Widgets
// ============================================================================

func (r *Registry) handleWidgets(ctx context.Context, input AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "list":
		return r.listWidgets(ctx, input.Filter)
	case "get":
		if input.Code == "" {
			return nil, fmt.Errorf("code is required for action 'get'")
		}
		return r.sdk.Widgets().GetOne(ctx, input.Code)
	case "install":
		if input.Code == "" {
			return nil, fmt.Errorf("code is required for action 'install'")
		}
		return r.sdk.Widgets().Install(ctx, input.Code)
	case "uninstall":
		if input.Code == "" {
			return nil, fmt.Errorf("code is required for action 'uninstall'")
		}
		err := r.sdk.Widgets().Uninstall(ctx, input.Code)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true, "uninstalled_code": input.Code}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s for widgets (available: list, get, install, uninstall)", input.Action)
	}
}

func (r *Registry) listWidgets(ctx context.Context, filter *IntegrationsFilter) ([]models.Widget, error) {
	f := &services.WidgetsFilter{
		Limit: 50,
		Page:  1,
	}
	if filter != nil {
		if filter.Limit > 0 {
			f.Limit = filter.Limit
		}
		if filter.Page > 0 {
			f.Page = filter.Page
		}
	}
	return r.sdk.Widgets().Get(ctx, f)
}

// ============================================================================
// Website Buttons (read-only)
// ============================================================================

func (r *Registry) handleWebsiteButtons(ctx context.Context, input AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "list":
		return r.listWebsiteButtons(ctx, input.Filter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id (source_id) is required for action 'get'")
		}
		return r.sdk.WebsiteButtons().GetOne(ctx, input.ID, nil)
	default:
		return nil, fmt.Errorf("unknown action: %s for website_buttons (available: list, get)", input.Action)
	}
}

func (r *Registry) listWebsiteButtons(ctx context.Context, filter *IntegrationsFilter) ([]models.WebsiteButton, error) {
	f := &services.WebsiteButtonsFilter{
		Limit: 50,
		Page:  1,
	}
	if filter != nil {
		if filter.Limit > 0 {
			f.Limit = filter.Limit
		}
		if filter.Page > 0 {
			f.Page = filter.Page
		}
	}
	return r.sdk.WebsiteButtons().Get(ctx, f, nil)
}

// ============================================================================
// Chat Templates
// ============================================================================

func (r *Registry) handleChatTemplates(ctx context.Context, input AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "list":
		return r.listChatTemplates(ctx, input.Filter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.sdk.ChatTemplates().GetOne(ctx, input.ID, nil)
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		return r.createChatTemplate(ctx, input.Data)
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		return r.updateChatTemplate(ctx, input.ID, input.Data)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'delete'")
		}
		err := r.sdk.ChatTemplates().Delete(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true, "deleted_id": input.ID}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s for chat_templates", input.Action)
	}
}

func (r *Registry) listChatTemplates(ctx context.Context, filter *IntegrationsFilter) ([]services.ChatTemplate, error) {
	f := &services.ChatTemplatesFilter{
		Limit: 50,
		Page:  1,
	}
	if filter != nil {
		if filter.Limit > 0 {
			f.Limit = filter.Limit
		}
		if filter.Page > 0 {
			f.Page = filter.Page
		}
	}
	return r.sdk.ChatTemplates().Get(ctx, f)
}

func (r *Registry) createChatTemplate(ctx context.Context, data map[string]any) ([]services.ChatTemplate, error) {
	template := services.ChatTemplate{}

	if name, ok := data["name"].(string); ok {
		template.Name = name
	}
	if content, ok := data["content"].(string); ok {
		template.Content = content
	}

	return r.sdk.ChatTemplates().Create(ctx, []services.ChatTemplate{template})
}

func (r *Registry) updateChatTemplate(ctx context.Context, id int, data map[string]any) ([]services.ChatTemplate, error) {
	template := services.ChatTemplate{ID: id}

	if name, ok := data["name"].(string); ok {
		template.Name = name
	}
	if content, ok := data["content"].(string); ok {
		template.Content = content
	}

	return r.sdk.ChatTemplates().Update(ctx, []services.ChatTemplate{template})
}

// ============================================================================
// Short Links
// ============================================================================

func (r *Registry) handleShortLinks(ctx context.Context, input AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "list":
		return r.listShortLinks(ctx, input.Filter)
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		return r.createShortLink(ctx, input.Data)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'delete'")
		}
		err := r.sdk.ShortLinks().Delete(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true, "deleted_id": input.ID}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s for short_links (available: list, create, delete)", input.Action)
	}
}

func (r *Registry) listShortLinks(ctx context.Context, filter *IntegrationsFilter) ([]models.ShortLink, error) {
	f := &services.ShortLinksFilter{
		Limit: 50,
		Page:  1,
	}
	if filter != nil {
		if filter.Limit > 0 {
			f.Limit = filter.Limit
		}
		if filter.Page > 0 {
			f.Page = filter.Page
		}
	}
	return r.sdk.ShortLinks().Get(ctx, f)
}

func (r *Registry) createShortLink(ctx context.Context, data map[string]any) ([]models.ShortLink, error) {
	link := models.ShortLink{}

	if url, ok := data["url"].(string); ok {
		link.URL = url
	}
	if entityType, ok := data["entity_type"].(string); ok {
		link.EntityType = entityType
	}
	if entityID, ok := data["entity_id"].(float64); ok {
		link.EntityID = int(entityID)
	}

	return r.sdk.ShortLinks().Create(ctx, []models.ShortLink{link})
}
