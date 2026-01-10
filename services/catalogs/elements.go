package catalogs

import (
	"context"
	"net/url"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListElements(ctx context.Context, catalogID int) ([]*models.CatalogElement, error) {
	elements, _, err := s.sdk.CatalogElements(catalogID).Get(ctx, nil)
	return elements, err
}

func (s *service) GetElement(ctx context.Context, catalogID, elementID int) (*models.CatalogElement, error) {
	return s.sdk.CatalogElements(catalogID).GetOne(ctx, elementID, url.Values{})
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
