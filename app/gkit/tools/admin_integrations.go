package tools

import (
	"encoding/json"
	"fmt"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func (r *Registry) RegisterAdminIntegrationsTool() {
	r.addTool(genkit.DefineTool[gkitmodels.AdminIntegrationsInput, any](
		r.g,
		"admin_integrations",
		"Work with webhooks, widgets, website buttons, chat templates and short links",
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
		return r.adminIntegrationsService.ListWebhooks(ctx)
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
		return r.adminIntegrationsService.ListWidgets(ctx)
	case "get":
		if input.Code == "" {
			return nil, fmt.Errorf("code is required for get widget")
		}
		return r.adminIntegrationsService.GetWidget(ctx, input.Code)
	case "create":
		var widgets []*amomodels.Widget
		data, _ := json.Marshal(input.Data["widgets"])
		if err := json.Unmarshal(data, &widgets); err != nil {
			return nil, fmt.Errorf("failed to parse widgets data: %w", err)
		}
		return r.adminIntegrationsService.CreateWidgets(ctx, widgets)
	case "update":
		var widgets []*amomodels.Widget
		data, _ := json.Marshal(input.Data["widgets"])
		if err := json.Unmarshal(data, &widgets); err != nil {
			return nil, fmt.Errorf("failed to parse widgets data: %w", err)
		}
		return r.adminIntegrationsService.UpdateWidgets(ctx, widgets)
	case "install":
		if input.Code == "" {
			return nil, fmt.Errorf("code is required for install widget")
		}
		return r.adminIntegrationsService.InstallWidget(ctx, input.Code)
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
		return r.adminIntegrationsService.ListWebsiteButtons(ctx)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for get website button")
		}
		return r.adminIntegrationsService.GetWebsiteButton(ctx, input.ID)
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
	default:
		return nil, fmt.Errorf("unknown action for website_buttons: %s", input.Action)
	}
}

func (r *Registry) handleChatTemplates(ctx *ai.ToolContext, input gkitmodels.AdminIntegrationsInput) (any, error) {
	switch input.Action {
	case "search", "list":
		return r.adminIntegrationsService.ListChatTemplates(ctx)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for delete chat template")
		}
		return nil, r.adminIntegrationsService.DeleteChatTemplate(ctx, input.ID)
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
		return r.adminIntegrationsService.ListShortLinks(ctx)
	case "create":
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
