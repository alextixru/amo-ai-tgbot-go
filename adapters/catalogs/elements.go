package catalogs

import (
	"context"
	"net/url"
	"strings"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) ListElements(ctx context.Context, catalogID int, filter *gkitmodels.CatalogFilter) ([]*models.CatalogElement, error) {
	var f *filters.CatalogElementsFilter
	if filter != nil {
		f = filters.NewCatalogElementsFilter()
		if filter.Query != "" {
			f.SetQuery(filter.Query)
		}
		if len(filter.IDs) > 0 {
			f.SetIDs(filter.IDs)
		}
		if filter.Limit > 0 {
			f.SetLimit(filter.Limit)
		}
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		}
	}
	elements, _, err := s.sdk.CatalogElements(catalogID).Get(ctx, f)
	return elements, err
}

func (s *service) GetElement(ctx context.Context, catalogID, elementID int, with []string) (*models.CatalogElement, error) {
	var params url.Values
	if len(with) > 0 {
		params = url.Values{}
		params.Set("with", strings.Join(with, ","))
	}
	return s.sdk.CatalogElements(catalogID).GetOne(ctx, elementID, params)
}

func (s *service) CreateElements(ctx context.Context, catalogID int, elements []*models.CatalogElement) ([]*models.CatalogElement, error) {
	res, _, err := s.sdk.CatalogElements(catalogID).Create(ctx, elements)
	return res, err
}

func (s *service) UpdateElements(ctx context.Context, catalogID int, elements []*models.CatalogElement) ([]*models.CatalogElement, error) {
	res, _, err := s.sdk.CatalogElements(catalogID).Update(ctx, elements)
	return res, err
}

func (s *service) LinkElement(ctx context.Context, catalogID, elementID int, entityType string, entityID int, metadata map[string]interface{}) error {
	return s.sdk.CatalogElements(catalogID).Link(ctx, elementID, entityType, entityID, metadata)
}

func (s *service) UnlinkElement(ctx context.Context, catalogID, elementID int, entityType string, entityID int) error {
	return s.sdk.CatalogElements(catalogID).Unlink(ctx, elementID, entityType, entityID)
}
