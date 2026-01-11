package products

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

var (
	productsCatalogID int
	catalogIDOnce     sync.Once
)

func (s *service) findProductsCatalogID(ctx context.Context) (int, error) {
	var err error
	catalogIDOnce.Do(func() {
		f := filters.NewCatalogsFilter().SetType("products")
		catalogs, _, findErr := s.sdk.Catalogs().Get(ctx, f)
		if findErr != nil {
			err = fmt.Errorf("failed to find products catalog: %w", findErr)
			return
		}
		if len(catalogs) == 0 {
			err = fmt.Errorf("products catalog not found")
			return
		}
		productsCatalogID = catalogs[0].ID
	})
	return productsCatalogID, err
}

func (s *service) SearchProducts(ctx context.Context, filter *gkitmodels.ProductFilter) ([]*models.CatalogElement, error) {
	catalogID, err := s.findProductsCatalogID(ctx)
	if err != nil {
		return nil, err
	}

	f := filters.NewCatalogElementsFilter()
	if filter != nil {
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
	} else {
		f.SetLimit(50)
		f.SetPage(1)
	}

	elements, _, err := s.sdk.CatalogElements(catalogID).Get(ctx, f)
	return elements, err
}

func (s *service) GetProduct(ctx context.Context, id int, with []string) (*models.CatalogElement, error) {
	catalogID, err := s.findProductsCatalogID(ctx)
	if err != nil {
		return nil, err
	}
	var params url.Values
	if len(with) > 0 {
		params = url.Values{}
		params.Set("with", strings.Join(with, ","))
	}
	return s.sdk.CatalogElements(catalogID).GetOne(ctx, id, params)
}

func (s *service) CreateProducts(ctx context.Context, elements []*models.CatalogElement) ([]*models.CatalogElement, error) {
	catalogID, err := s.findProductsCatalogID(ctx)
	if err != nil {
		return nil, err
	}
	created, _, err := s.sdk.CatalogElements(catalogID).Create(ctx, elements)
	return created, err
}

func (s *service) UpdateProducts(ctx context.Context, elements []*models.CatalogElement) ([]*models.CatalogElement, error) {
	catalogID, err := s.findProductsCatalogID(ctx)
	if err != nil {
		return nil, err
	}
	updated, _, err := s.sdk.CatalogElements(catalogID).Update(ctx, elements)
	return updated, err
}

func (s *service) DeleteProducts(ctx context.Context, ids []int) error {
	catalogID, err := s.findProductsCatalogID(ctx)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if err := s.sdk.CatalogElements(catalogID).Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete product %d: %w", id, err)
		}
	}
	return nil
}

func (s *service) GetProductsByEntity(ctx context.Context, entityType string, entityID int) ([]models.EntityLink, error) {
	links, err := s.sdk.Links().Get(ctx, entityType, entityID, nil)
	if err != nil {
		return nil, err
	}

	catalogID, err := s.findProductsCatalogID(ctx)
	if err != nil {
		return nil, err
	}

	var productLinks []models.EntityLink
	for _, link := range links {
		// В метаданных связи должен быть catalog_id
		if link.ToEntityType == "catalog_elements" {
			if cid, ok := link.Metadata["catalog_id"].(float64); ok && int(cid) == catalogID {
				productLinks = append(productLinks, *link)
			}
		}
	}
	return productLinks, nil
}

func (s *service) LinkProduct(ctx context.Context, entityType string, entityID int, productID int, quantity int, priceID int) error {
	catalogID, err := s.findProductsCatalogID(ctx)
	if err != nil {
		return err
	}

	link := models.NewCatalogElementLink(entityType, entityID, catalogID, productID, float64(quantity), priceID)
	_, err = s.sdk.Links().Link(ctx, entityType, entityID, []*models.EntityLink{link})
	return err
}

func (s *service) UnlinkProduct(ctx context.Context, entityType string, entityID int, productID int) error {
	link := &models.EntityLink{
		ToEntityType: "catalog_elements",
		ToEntityID:   productID,
	}
	return s.sdk.Links().Unlink(ctx, entityType, entityID, []*models.EntityLink{link})
}
