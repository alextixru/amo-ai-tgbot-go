package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// NoteToolInput входные параметры для инструмента crm_notes
type NoteToolInput struct {
	// Action действие: list, create
	Action string `json:"action"`

	// Filter фильтр для поиска (используется при action=list)
	Filter *NoteFilterInput `json:"filter,omitempty"`

	// EntityType тип сущности (leads, contacts, companies, customers) для notes
	EntityType string `json:"entity_type,omitempty"`

	// Data данные для создания (используется при action=create)
	Data *NoteDataInput `json:"data,omitempty"`
}

// NoteFilterInput параметры фильтрации примечаний
type NoteFilterInput struct {
	// NoteType тип примечания (common, call_in, etc.)
	NoteType string `json:"note_type,omitempty"`
	// EntityID ID сущности (фильтр по привязке)
	EntityID int `json:"entity_id,omitempty"`
	// Limit количество элементов (по умолчанию 50)
	Limit int `json:"limit,omitempty"`
}

// NoteDataInput данные примечания
type NoteDataInput struct {
	// NoteType тип примечания (common)
	NoteType string `json:"note_type,omitempty"`
	// Text текст примечания (для common)
	Text string `json:"text,omitempty"`
	// EntityID ID сущности (к которой привязываем)
	EntityID int `json:"entity_id,omitempty"`
}

// registerNotesTools регистрирует инструменты для работы с примечаниями
func (r *Registry) registerNotesTools() {
	r.addTool(genkit.DefineTool[NoteToolInput, any](
		r.g,
		"crm_notes",
		"Управление примечаниями (Notes). Позволяет создавать и просматривать примечания к сущностям.",
		func(ctx *ai.ToolContext, input NoteToolInput) (any, error) {
			if input.EntityType == "" {
				return nil, fmt.Errorf("entity_type is required for notes (leads, contacts, companies, customers)")
			}

			switch input.Action {
			case "list":
				return r.listNotes(ctx.Context, input.EntityType, input.Filter)
			case "create":
				if input.Data == nil {
					return nil, fmt.Errorf("data is required for 'create' action")
				}
				// ID сущности должен быть в Data
				if input.Data.EntityID == 0 {
					return nil, fmt.Errorf("data.entity_id is required for 'create' action")
				}
				return r.createNote(ctx.Context, input.EntityType, input.Data)
			default:
				return nil, fmt.Errorf("unknown action: %s", input.Action)
			}
		},
	))
}

// listNotes выполняет поиск примечаний
func (r *Registry) listNotes(ctx context.Context, entityType string, filter *NoteFilterInput) ([]models.Note, error) {
	sdkFilter := &services.NotesFilter{
		Limit: 50,
	}

	if filter != nil {
		if filter.Limit > 0 {
			sdkFilter.Limit = filter.Limit
		}
		if filter.NoteType != "" {
			sdkFilter.FilterByNoteType = []string{filter.NoteType}
		}

		// Если передан EntityID, используем GetByParent
		if filter.EntityID > 0 {
			return r.sdk.Notes().GetByParent(ctx, entityType, filter.EntityID, sdkFilter)
		}
	}

	// Иначе используем общий Get
	return r.sdk.Notes().Get(ctx, entityType, sdkFilter)
}

// createNote создает новое примечание
func (r *Registry) createNote(ctx context.Context, entityType string, data *NoteDataInput) (*models.Note, error) {
	note := models.Note{
		EntityID: data.EntityID,
	}

	// Устанавливаем тип
	if data.NoteType != "" {
		note.NoteType = models.NoteType(data.NoteType)
	} else {
		note.NoteType = models.NoteTypeCommon
	}

	// Для common типа текст передается в Params
	if note.NoteType == models.NoteTypeCommon {
		note.Params = &models.NoteParams{
			Text: data.Text,
		}
	}

	createdNotes, err := r.sdk.Notes().Create(ctx, entityType, []models.Note{note})
	if err != nil {
		return nil, err
	}

	if len(createdNotes) == 0 {
		return nil, fmt.Errorf("failed to create note: empty response")
	}

	return &createdNotes[0], nil
}
