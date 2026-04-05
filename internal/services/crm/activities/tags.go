package activities

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func convertTags(tags []*models.Tag) []*TagOutput {
	out := make([]*TagOutput, 0, len(tags))
	for _, t := range tags {
		if t == nil {
			continue
		}
		out = append(out, &TagOutput{ID: t.ID, Name: t.Name})
	}
	return out
}

func (s *service) ListTags(ctx context.Context, entityType string, filter *gkitmodels.TagsFilter) ([]*TagOutput, error) {
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
	if err != nil {
		return nil, err
	}
	return convertTags(tags), nil
}

func (s *service) CreateTag(ctx context.Context, entityType string, name string) (*TagOutput, error) {
	tags, err := s.CreateTags(ctx, entityType, []string{name})
	if err != nil {
		return nil, err
	}
	if len(tags) > 0 {
		return tags[0], nil
	}
	return nil, nil
}

func (s *service) CreateTags(ctx context.Context, entityType string, names []string) ([]*TagOutput, error) {
	items := make([]*models.Tag, len(names))
	for i, name := range names {
		items[i] = &models.Tag{Name: name}
	}
	tags, _, err := s.sdk.Tags().Create(ctx, entityType, items)
	if err != nil {
		return nil, err
	}
	return convertTags(tags), nil
}

func (s *service) DeleteTag(ctx context.Context, entityType string, tagID int) error {
	tag := &models.Tag{ID: tagID}
	return s.sdk.Tags().Delete(ctx, entityType, []*models.Tag{tag})
}

// DeleteTagByName находит тег по имени и удаляет его.
func (s *service) DeleteTagByName(ctx context.Context, entityType string, tagName string) error {
	f := filters.NewTagsFilter()
	f.SetName(tagName).SetLimit(5)
	tags, _, err := s.sdk.Tags().Get(ctx, entityType, f)
	if err != nil {
		return fmt.Errorf("DeleteTagByName: поиск тега '%s': %w", tagName, err)
	}
	// Ищем точное совпадение
	var found *models.Tag
	for _, t := range tags {
		if t != nil && t.Name == tagName {
			found = t
			break
		}
	}
	if found == nil {
		return fmt.Errorf("тег '%s' не найден для сущности '%s'", tagName, entityType)
	}
	return s.sdk.Tags().Delete(ctx, entityType, []*models.Tag{found})
}
