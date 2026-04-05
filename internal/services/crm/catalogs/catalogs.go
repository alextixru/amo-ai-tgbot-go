package catalogs

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/constants"
	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) ListCatalogs(ctx context.Context, filter *gkitmodels.CatalogFilter) (*CatalogListResult, error) {
	var f *filters.CatalogsFilter
	if filter != nil && filter.Type != "" {
		f = filters.NewCatalogsFilter()
		f.SetType(filter.Type)
	}

	catalogs, meta, err := s.sdk.Catalogs().Get(ctx, f)
	if err != nil {
		return nil, err
	}

	items := make([]*CatalogItem, 0, len(catalogs))
	for _, c := range catalogs {
		items = append(items, s.normalizeCatalog(c))
	}

	result := &CatalogListResult{
		Items: items,
	}
	if meta != nil {
		result.Total = meta.TotalItems
		result.Page = meta.Page
		result.HasMore = meta.HasMore
	}

	return result, nil
}

func (s *service) GetCatalog(ctx context.Context, name string) (*CatalogItem, error) {
	id, err := s.resolveCatalogName(name)
	if err != nil {
		return nil, err
	}
	c, err := s.sdk.Catalogs().GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.normalizeCatalog(c), nil
}

func (s *service) CreateCatalog(ctx context.Context, data *gkitmodels.CatalogData) (*CatalogItem, error) {
	catalog := mapCatalogDataToModel(data)
	res, _, err := s.sdk.Catalogs().Create(ctx, []*models.Catalog{catalog})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	created := res[0]
	// Обновляем внутренние мапы
	s.catalogsByName[created.Name] = created.ID
	s.catalogsByID[created.ID] = created.Name
	return s.normalizeCatalog(created), nil
}

func (s *service) UpdateCatalog(ctx context.Context, name string, data *gkitmodels.CatalogData) (*CatalogItem, error) {
	id, err := s.resolveCatalogName(name)
	if err != nil {
		return nil, err
	}
	catalog := mapCatalogDataToModel(data)
	catalog.ID = id // BUG2 fix: явно проставляем ID
	res, _, err := s.sdk.Catalogs().Update(ctx, []*models.Catalog{catalog})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	updated := res[0]
	// Обновляем внутренние мапы если имя изменилось
	if updated.Name != name {
		delete(s.catalogsByName, name)
	}
	s.catalogsByName[updated.Name] = updated.ID
	s.catalogsByID[updated.ID] = updated.Name
	return s.normalizeCatalog(updated), nil
}

func (s *service) DeleteCatalog(ctx context.Context, name string) error {
	id, err := s.resolveCatalogName(name)
	if err != nil {
		return err
	}
	if err := s.sdk.Catalogs().Delete(ctx, id); err != nil {
		return err
	}
	// Убираем из внутренних мап
	delete(s.catalogsByName, name)
	delete(s.catalogsByID, id)
	return nil
}

// mapCatalogDataToModel преобразует CatalogData из tool input в SDK-модель.
// Исправляет BUG1: вместо json roundtrip — явный маппинг полей.
func mapCatalogDataToModel(data *gkitmodels.CatalogData) *models.Catalog {
	if data == nil {
		return &models.Catalog{}
	}
	return &models.Catalog{
		Name:            data.Name,
		Type:            constants.CatalogType(data.Type),
		CanAddElements:  data.CanAddElements,
		CanShowInCards:  data.CanShowInCards,
		CanLinkMultiple: data.CanLinkMultiple,
	}
}
