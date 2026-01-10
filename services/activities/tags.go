package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListTags(ctx context.Context, entityType string) ([]*models.Tag, error) {
	f := filters.NewTagsFilter()
	f.SetLimit(50)
	f.SetPage(1)
	tags, _, err := s.sdk.Tags().Get(ctx, entityType, f)
	return tags, err
}

func (s *service) CreateTag(ctx context.Context, entityType string, name string) (*models.Tag, error) {
	tag := &models.Tag{
		Name: name,
	}
	tags, _, err := s.sdk.Tags().Create(ctx, entityType, []*models.Tag{tag})
	if err != nil {
		return nil, err
	}
	if len(tags) > 0 {
		return tags[0], nil
	}
	return nil, nil
}

func (s *service) DeleteTag(ctx context.Context, entityType string, tagID int) error {
	tag := &models.Tag{ID: tagID}
	return s.sdk.Tags().Delete(ctx, entityType, []*models.Tag{tag})
}
