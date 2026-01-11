package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models"
)

func (s *service) ListNotes(ctx context.Context, parent gkitmodels.ParentEntity, filter *gkitmodels.NotesFilter, with []string) ([]*models.Note, error) {
	f := filters.NewNotesFilter()
	if filter != nil {
		if filter.Limit > 0 {
			f.SetLimit(filter.Limit)
		} else {
			f.SetLimit(50)
		}
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		}
		if len(filter.IDs) > 0 {
			f.SetIDs(filter.IDs)
		}
		if len(filter.NoteTypes) > 0 {
			f.SetNoteTypes(filter.NoteTypes)
		}
		if filter.UpdatedAt != nil {
			val := int(*filter.UpdatedAt)
			f.SetUpdatedAt(&val, nil)
		}
	} else {
		f.SetLimit(50).SetPage(1)
	}

	if len(with) > 0 {
		f.With = with
	}

	notes, _, err := s.sdk.Notes().GetByParent(ctx, parent.Type, parent.ID, f)
	return notes, err
}

func (s *service) GetNote(ctx context.Context, entityType string, id int) (*models.Note, error) {
	return s.sdk.Notes().GetOne(ctx, entityType, id, nil)
}

func (s *service) CreateNote(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.NoteData) (*models.Note, error) {
	notes, err := s.CreateNotes(ctx, parent, []gkitmodels.NoteData{*data})
	if err != nil {
		return nil, err
	}
	if len(notes) > 0 {
		return notes[0], nil
	}
	return nil, nil
}

func (s *service) CreateNotes(ctx context.Context, parent gkitmodels.ParentEntity, data []gkitmodels.NoteData) ([]*models.Note, error) {
	items := make([]*models.Note, len(data))
	for i, d := range data {
		note := &models.Note{
			EntityID: parent.ID,
			Params: &models.NoteParams{
				Text: d.Text,
			},
		}
		if d.NoteType != "" {
			note.NoteType = models.NoteType(d.NoteType)
		} else {
			note.NoteType = models.NoteTypeCommon
		}
		items[i] = note
	}
	notes, _, err := s.sdk.Notes().Create(ctx, parent.Type, items)
	return notes, err
}

func (s *service) UpdateNote(ctx context.Context, entityType string, id int, data *gkitmodels.NoteData) (*models.Note, error) {
	note := &models.Note{
		BaseModel: models.BaseModel{ID: id},
		Params: &models.NoteParams{
			Text: data.Text,
		},
	}
	notes, _, err := s.sdk.Notes().Update(ctx, entityType, []*models.Note{note})
	if err != nil {
		return nil, err
	}
	if len(notes) > 0 {
		return notes[0], nil
	}
	return nil, nil
}
