package tools

import (
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// LinkToolInput входные параметры для инструмента crm_manage_links
type LinkToolInput struct {
	// Action действие: link, unlink
	Action string `json:"action"`

	// FromEntity сущность, К КОТОРОЙ привязываем
	FromEntity EntityLinkInput `json:"from_entity"`

	// ToEntity сущность, КОТОРУЮ привязываем
	ToEntity EntityLinkInput `json:"to_entity"`
}

// EntityLinkInput параметры сущности для связывания
type EntityLinkInput struct {
	// Type тип сущности: leads, contacts, companies, customers
	Type string `json:"type"`
	// ID идентификатор сущности
	ID int `json:"id"`
}

// registerLinksTool регистрирует инструмент для управления связями
func (r *Registry) registerLinksTool() {
	r.addTool(genkit.DefineTool[LinkToolInput, any](
		r.g,
		"crm_manage_links",
		"Управление связями (Links). Привязывайте контакты к сделкам, компании к контактам и т.д.",
		func(ctx *ai.ToolContext, input LinkToolInput) (any, error) {
			if input.FromEntity.ID == 0 || input.FromEntity.Type == "" {
				return nil, fmt.Errorf("from_entity (id, type) is required")
			}
			if input.ToEntity.ID == 0 || input.ToEntity.Type == "" {
				return nil, fmt.Errorf("to_entity (id, type) is required")
			}

			link := models.EntityLink{
				ToEntityID:   input.ToEntity.ID,
				ToEntityType: input.ToEntity.Type,
			}

			switch input.Action {
			case "link":
				links, err := r.sdk.Links().Link(ctx.Context, input.FromEntity.Type, input.FromEntity.ID, []models.EntityLink{link})
				if err != nil {
					return nil, err
				}
				if len(links) == 0 {
					return nil, fmt.Errorf("failed to link entities: empty response")
				}
				return "Entities linked successfully", nil

			case "unlink":
				err := r.sdk.Links().Unlink(ctx.Context, input.FromEntity.Type, input.FromEntity.ID, []models.EntityLink{link})
				if err != nil {
					return nil, err
				}
				return "Entities unlinked successfully", nil

			default:
				return nil, fmt.Errorf("unknown action: %s", input.Action)
			}
		},
	))
}
