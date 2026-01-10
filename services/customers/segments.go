package customers

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListSegments(ctx context.Context) ([]*models.CustomerSegment, error) {
	res, _, err := s.sdk.Segments().Get(ctx, nil)
	return res, err
}

func (s *service) GetSegment(ctx context.Context, id int) (*models.CustomerSegment, error) {
	return s.sdk.Segments().GetOne(ctx, id)
}

func (s *service) CreateSegments(ctx context.Context, segments []*models.CustomerSegment) ([]*models.CustomerSegment, error) {
	res, _, err := s.sdk.Segments().Create(ctx, segments)
	return res, err
}

func (s *service) DeleteSegment(ctx context.Context, id int) error {
	return s.sdk.Segments().Delete(ctx, id)
}
