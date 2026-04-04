package activities

import (
	"context"
	"fmt"
	"strconv"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
)

// GetTalk получает беседу по ID
func (s *service) GetTalk(ctx context.Context, talkID string) (*models.Talk, error) {
	id, err := strconv.Atoi(talkID)
	if err != nil {
		return nil, fmt.Errorf("invalid talk ID (must be integer): %w", err)
	}
	return s.sdk.Talks().GetOne(ctx, id)
}

// CloseTalk закрывает беседу
func (s *service) CloseTalk(ctx context.Context, talkID string, forceClose bool) error {
	opts := &services.TalkCloseOptions{ForceClose: forceClose}
	return s.sdk.Talks().Close(ctx, talkID, opts)
}
