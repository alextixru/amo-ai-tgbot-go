package activities

import (
	"context"
	"fmt"
	"strconv"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
)

// convertTalk конвертирует models.Talk в TalkOutput
func (s *service) convertTalk(t *models.Talk) *TalkOutput {
	if t == nil {
		return nil
	}
	return &TalkOutput{
		TalkID:     t.TalkID,
		ChatID:     t.ChatID,
		ContactID:  t.ContactID,
		EntityID:   t.EntityID,
		EntityType: t.EntityType,
		Origin:     t.Origin,
		IsInWork:   t.IsInWork,
		IsRead:     t.IsRead,
		Rate:       t.Rate,
		CreatedAt:  toISO(t.CreatedAt),
		UpdatedAt:  toISO(t.UpdatedAt),
	}
}

// GetTalk получает беседу по ID
func (s *service) GetTalk(ctx context.Context, talkID string) (*TalkOutput, error) {
	id, err := strconv.Atoi(talkID)
	if err != nil {
		return nil, fmt.Errorf("invalid talk ID (must be integer): %w", err)
	}
	talk, err := s.sdk.Talks().GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.convertTalk(talk), nil
}

// CloseTalk закрывает беседу
func (s *service) CloseTalk(ctx context.Context, talkID string, forceClose bool) error {
	opts := &services.TalkCloseOptions{ForceClose: forceClose}
	return s.sdk.Talks().Close(ctx, talkID, opts)
}
