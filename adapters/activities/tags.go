package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) ListTags(ctx context.Context, entityType string, filter *gkitmodels.TagsFilter) ([]*models.Tag, error) {
	f := filters.NewTagsFilter()
	if filter != nil {
		if filter.Limit > 0 {
			f.SetLimit(filter.Limit)
		} else {
			f.SetLimit(50)
		}
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		}
		if filter.Query != "" {
			f.SetQuery(filter.Query)
		}
		if filter.Name != "" {
			f.SetName(filter.Name)
		}
		if len(filter.IDs) > 0 {
			f.SetIDs(filter.IDs)
		}
	} else {
		f.SetLimit(50).SetPage(1)
	}
	tags, _, err := s.sdk.Tags().Get(ctx, entityType, f)
	return tags, err
}

func (s *service) CreateTag(ctx context.Context, entityType string, name string) (*models.Tag, error) {
	tags, err := s.CreateTags(ctx, entityType, []string{name})
	if err != nil {
		return nil, err
	}
	if len(tags) > 0 {
		return tags[0], nil
	}
	return nil, nil
}

func (s *service) CreateTags(ctx context.Context, entityType string, names []string) ([]*models.Tag, error) {
	items := make([]*models.Tag, len(names))
	for i, name := range names {
		items[i] = &models.Tag{
			Name: name,
		}
	}
	tags, _, err := s.sdk.Tags().Create(ctx, entityType, items)
	return tags, err
}

func (s *service) DeleteTag(ctx context.Context, entityType string, tagID int) error {
	tag := &models.Tag{ID: tagID}
	return s.sdk.Tags().Delete(ctx, entityType, []*models.Tag{tag})
}
