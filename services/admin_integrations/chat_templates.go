package admin_integrations

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListChatTemplates(ctx context.Context) ([]*models.ChatTemplate, error) {
	templates, _, err := s.sdk.ChatTemplates().Get(ctx, nil)
	return templates, err
}

func (s *service) DeleteChatTemplate(ctx context.Context, id int) error {
	return s.sdk.ChatTemplates().Delete(ctx, id)
}

func (s *service) SendChatTemplateOnReview(ctx context.Context, id int) ([]models.ChatTemplateReview, error) {
	return s.sdk.ChatTemplates().SendOnReview(ctx, id)
}

func (s *service) UpdateChatTemplateReviewStatus(ctx context.Context, templateID, reviewID int, status string) (*models.ChatTemplateReview, error) {
	return s.sdk.ChatTemplates().UpdateReviewStatus(ctx, templateID, reviewID, status)
}
