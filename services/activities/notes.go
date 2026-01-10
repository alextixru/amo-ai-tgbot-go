package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) ListNotes(ctx context.Context, parent gkitmodels.ParentEntity) ([]*models.Note, error) {
	f := filters.NewNotesFilter()
	f.SetLimit(50)
	f.SetPage(1)
	notes, _, err := s.sdk.Notes().GetByParent(ctx, parent.Type, parent.ID, f)
	return notes, err
}

func (s *service) GetNote(ctx context.Context, entityType string, id int) (*models.Note, error) {
	return s.sdk.Notes().GetOne(ctx, entityType, id, nil)
}

func (s *service) CreateNote(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.ActivityData) (*models.Note, error) {
	note := &models.Note{
		EntityID: parent.ID,
		Params: &models.NoteParams{
			Text: data.Text,
		},
	}
	if data.NoteType != "" {
		note.NoteType = models.NoteType(data.NoteType)
	} else {
		note.NoteType = models.NoteTypeCommon
	}
	notes, _, err := s.sdk.Notes().Create(ctx, parent.Type, []*models.Note{note})
	if err != nil {
		return nil, err
	}
	if len(notes) > 0 {
		return notes[0], nil
	}
	return nil, nil
}

func (s *service) UpdateNote(ctx context.Context, entityType string, id int, data *gkitmodels.ActivityData) (*models.Note, error) {
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
