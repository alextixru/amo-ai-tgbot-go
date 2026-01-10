package catalogs

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListCatalogs(ctx context.Context) ([]*models.Catalog, error) {
	catalogs, _, err := s.sdk.Catalogs().Get(ctx, nil)
	return catalogs, err
}

func (s *service) GetCatalog(ctx context.Context, id int) (*models.Catalog, error) {
	return s.sdk.Catalogs().GetOne(ctx, id)
}

func (s *service) CreateCatalogs(ctx context.Context, catalogs []*models.Catalog) ([]*models.Catalog, error) {
	res, _, err := s.sdk.Catalogs().Create(ctx, catalogs)
	return res, err
}

func (s *service) UpdateCatalogs(ctx context.Context, catalogs []*models.Catalog) ([]*models.Catalog, error) {
	res, _, err := s.sdk.Catalogs().Update(ctx, catalogs)
	return res, err
}
