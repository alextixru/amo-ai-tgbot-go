package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) convertNote(n *models.Note) *NoteOutput {
	if n == nil {
		return nil
	}
	out := &NoteOutput{
		ID:            n.ID,
		EntityID:      n.EntityID,
		NoteType:      string(n.NoteType),
		CreatedByName: s.resolveUserID(n.CreatedBy),
		UpdatedByName: s.resolveUserID(n.UpdatedBy),
		CreatedAt:     toISO(n.CreatedAt),
		UpdatedAt:     toISO(n.UpdatedAt),
	}
	// Текст может быть в Params.Text или в поле Text
	if n.Params != nil && n.Params.Text != "" {
		out.Text = n.Params.Text
	} else if n.Text != "" {
		out.Text = n.Text
	}
	return out
}

func (s *service) ListNotes(ctx context.Context, parent gkitmodels.ParentEntity, filter *gkitmodels.NotesFilter, with []string) ([]*NoteOutput, error) {
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
	if err != nil {
		return nil, err
	}
	out := make([]*NoteOutput, 0, len(notes))
	for _, n := range notes {
		out = append(out, s.convertNote(n))
	}
	return out, nil
}

func (s *service) GetNote(ctx context.Context, entityType string, id int) (*NoteOutput, error) {
	n, err := s.sdk.Notes().GetOne(ctx, entityType, id, nil)
	if err != nil {
		return nil, err
	}
	return s.convertNote(n), nil
}

func (s *service) CreateNote(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.NoteData) (*NoteOutput, error) {
	notes, err := s.CreateNotes(ctx, parent, []gkitmodels.NoteData{*data})
	if err != nil {
		return nil, err
	}
	if len(notes) > 0 {
		return notes[0], nil
	}
	return nil, nil
}

func (s *service) CreateNotes(ctx context.Context, parent gkitmodels.ParentEntity, data []gkitmodels.NoteData) ([]*NoteOutput, error) {
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
	if err != nil {
		return nil, err
	}
	out := make([]*NoteOutput, 0, len(notes))
	for _, n := range notes {
		out = append(out, s.convertNote(n))
	}
	return out, nil
}

func (s *service) UpdateNote(ctx context.Context, entityType string, id int, data *gkitmodels.NoteData) (*NoteOutput, error) {
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
		return s.convertNote(notes[0]), nil
	}
	return nil, nil
}
