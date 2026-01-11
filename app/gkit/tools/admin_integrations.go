package tools

import (
	"encoding/json"
	"fmt"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func (r *Registry) RegisterAdminIntegrationsTool() {
	r.addTool(genkit.DefineTool[gkitmodels.AdminIntegrationsInput, any](
		r.g,
		"admin_integrations",
		"Work with webhooks, widgets, website buttons (with scripts), chat templates and short links",
		func(ctx *ai.ToolContext, input gkitmodels.AdminIntegrationsInput) (any, error) {
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
		},
	))
}

func (r *Registry) handleWebhooks(ctx *ai.ToolContext, input gkitmodels.AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "search", "list":
		var filter *filters.WebhooksFilter
		if input.Filter != nil {
			filter = filters.NewWebhooksFilter()
			if input.Filter.Destination != "" {
				filter.SetDestination(input.Filter.Destination)
			}
		}
		return r.adminIntegrationsService.ListWebhooks(ctx, filter)
	case "subscribe":
		dest, _ := input.Data["destination"].(string)
		if dest == "" {
			return nil, fmt.Errorf("destination is required for subscribe")
		}
		var settings []string
		if s, ok := input.Data["settings"].([]any); ok {
			for _, v := range s {
				if str, ok := v.(string); ok {
					settings = append(settings, str)
				}
			}
		}
		return r.adminIntegrationsService.SubscribeWebhook(ctx, dest, settings)
	case "unsubscribe":
		dest, _ := input.Data["destination"].(string)
		if dest == "" {
			return nil, fmt.Errorf("destination is required for unsubscribe")
		}
		var settings []string
		if s, ok := input.Data["settings"].([]any); ok {
			for _, v := range s {
				if str, ok := v.(string); ok {
					settings = append(settings, str)
				}
			}
		}
		return nil, r.adminIntegrationsService.UnsubscribeWebhook(ctx, dest, settings)
	default:
		return nil, fmt.Errorf("unknown action for webhooks: %s", input.Action)
	}
}

func (r *Registry) handleWidgets(ctx *ai.ToolContext, input gkitmodels.AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "search", "list":
		var filter *filters.WidgetsFilter
		if input.Filter != nil {
			filter = filters.NewWidgetsFilter()
			if input.Filter.Limit > 0 {
				filter.SetLimit(input.Filter.Limit)
			}
			if input.Filter.Page > 0 {
				filter.SetPage(input.Filter.Page)
			}
		}
		return r.adminIntegrationsService.ListWidgets(ctx, filter)
	case "get":
		if input.Code == "" {
			return nil, fmt.Errorf("code is required for get widget")
		}
		return r.adminIntegrationsService.GetWidget(ctx, input.Code)
	case "install":
		if input.Code == "" {
			return nil, fmt.Errorf("code is required for install widget")
		}
		settings := input.Settings
		if settings == nil {
			if s, ok := input.Data["settings"].(map[string]any); ok {
				settings = s
			}
		}
		return r.adminIntegrationsService.InstallWidget(ctx, input.Code, settings)
	case "uninstall":
		if input.Code == "" {
			return nil, fmt.Errorf("code is required for uninstall widget")
		}
		return nil, r.adminIntegrationsService.UninstallWidget(ctx, input.Code)
	default:
		return nil, fmt.Errorf("unknown action for widgets: %s", input.Action)
	}
}

func (r *Registry) handleWebsiteButtons(ctx *ai.ToolContext, input gkitmodels.AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "search", "list":
		var filter *services.WebsiteButtonsFilter
		var with []string
		if input.Filter != nil {
			filter = &services.WebsiteButtonsFilter{}
			if input.Filter.Limit > 0 {
				filter.Limit = input.Filter.Limit
			}
			if input.Filter.Page > 0 {
				filter.Page = input.Filter.Page
			}
			with = input.Filter.With
		}
		return r.adminIntegrationsService.ListWebsiteButtons(ctx, filter, with)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for get website button")
		}
		var with []string
		if input.Filter != nil {
			with = input.Filter.With
		}
		return r.adminIntegrationsService.GetWebsiteButton(ctx, input.ID, with)
	case "create":
		var req amomodels.WebsiteButtonCreateRequest
		data, _ := json.Marshal(input.Data)
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, fmt.Errorf("failed to parse website button create request: %w", err)
		}
		return r.adminIntegrationsService.CreateWebsiteButton(ctx, &req)
	case "update":
		var req amomodels.WebsiteButtonUpdateRequest
		data, _ := json.Marshal(input.Data)
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, fmt.Errorf("failed to parse website button update request: %w", err)
		}
		return r.adminIntegrationsService.UpdateWebsiteButton(ctx, &req)
	case "add_chat":
		if input.ID == 0 {
			return nil, fmt.Errorf("id (source_id) is required for add_chat")
		}
		return nil, r.adminIntegrationsService.AddOnlineChat(ctx, input.ID)
	default:
		return nil, fmt.Errorf("unknown action for website_buttons: %s", input.Action)
	}
}

func (r *Registry) handleChatTemplates(ctx *ai.ToolContext, input gkitmodels.AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "search", "list":
		var filter *filters.TemplatesFilter
		if input.Filter != nil {
			filter = filters.NewTemplatesFilter()
			if input.Filter.Limit > 0 {
				filter.SetLimit(input.Filter.Limit)
			}
			if input.Filter.Page > 0 {
				filter.SetPage(input.Filter.Page)
			}
			if len(input.Filter.ExternalIDs) > 0 {
				filter.SetExternalIDs(input.Filter.ExternalIDs)
			}
		}
		return r.adminIntegrationsService.ListChatTemplates(ctx, filter)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for delete chat template")
		}
		return nil, r.adminIntegrationsService.DeleteChatTemplate(ctx, input.ID)
	case "delete_many":
		if len(input.IDs) == 0 {
			return nil, fmt.Errorf("ids are required for delete_many")
		}
		return nil, r.adminIntegrationsService.DeleteChatTemplates(ctx, input.IDs)
	case "send_review":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for send chat template on review")
		}
		return r.adminIntegrationsService.SendChatTemplateOnReview(ctx, input.ID)
	case "update_review":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for update chat template review status")
		}
		reviewID, _ := input.Data["review_id"].(float64)
		status, _ := input.Data["status"].(string)
		if reviewID == 0 || status == "" {
			return nil, fmt.Errorf("review_id and status are required for update_review")
		}
		return r.adminIntegrationsService.UpdateChatTemplateReviewStatus(ctx, input.ID, int(reviewID), status)
	default:
		return nil, fmt.Errorf("unknown action for chat_templates: %s", input.Action)
	}
}

func (r *Registry) handleShortLinks(ctx *ai.ToolContext, input gkitmodels.AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "search", "list":
		var filter *filters.ShortLinksFilter
		if input.Filter != nil {
			filter = filters.NewShortLinksFilter()
			if input.Filter.Limit > 0 {
				filter.SetLimit(input.Filter.Limit)
			}
			if input.Filter.Page > 0 {
				filter.SetPage(input.Filter.Page)
			}
		}
		return r.adminIntegrationsService.ListShortLinks(ctx, filter)
	case "create":
		if len(input.URLs) > 0 {
			return r.adminIntegrationsService.CreateShortLinks(ctx, input.URLs)
		}
		url, _ := input.Data["url"].(string)
		if url == "" {
			return nil, fmt.Errorf("url is required for create short link")
		}
		return r.adminIntegrationsService.CreateShortLink(ctx, url)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for delete short link")
		}
		return nil, r.adminIntegrationsService.DeleteShortLink(ctx, input.ID)
	default:
		return nil, fmt.Errorf("unknown action for short_links: %s", input.Action)
	}
}
