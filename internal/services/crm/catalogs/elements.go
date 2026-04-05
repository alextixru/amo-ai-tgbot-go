package catalogs

import (
	"context"
	"net/url"
	"strings"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) ListElements(ctx context.Context, catalogName string, filter *gkitmodels.CatalogFilter) (*ElementListResult, error) {
	id, err := s.resolveCatalogName(catalogName)
	if err != nil {
		return nil, err
	}

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

	elements, meta, err := s.sdk.CatalogElements(id).Get(ctx, f)
	if err != nil {
		return nil, err
	}

	items := make([]*ElementItem, 0, len(elements))
	for _, e := range elements {
		items = append(items, s.normalizeElement(e))
	}

	result := &ElementListResult{
		Items: items,
	}
	if meta != nil {
		result.Total = meta.TotalItems
		result.Page = meta.Page
		result.HasMore = meta.HasMore
	}

	return result, nil
}

func (s *service) GetElement(ctx context.Context, catalogName string, elementID int, with []string) (*ElementItem, error) {
	id, err := s.resolveCatalogName(catalogName)
	if err != nil {
		return nil, err
	}

	// Всегда запрашиваем все доступные with-параметры
	defaultWith := []string{"invoice_link", "supplier_field_values"}
	merged := mergeWith(defaultWith, with)

	params := url.Values{}
	params.Set("with", strings.Join(merged, ","))

	e, err := s.sdk.CatalogElements(id).GetOne(ctx, elementID, params)
	if err != nil {
		return nil, err
	}
	return s.normalizeElement(e), nil
}

func (s *service) CreateElement(ctx context.Context, catalogName string, data *gkitmodels.CatalogElementData) (*ElementItem, error) {
	id, err := s.resolveCatalogName(catalogName)
	if err != nil {
		return nil, err
	}

	element := mapElementDataToModel(data)
	res, _, err := s.sdk.CatalogElements(id).Create(ctx, []*models.CatalogElement{element})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return s.normalizeElement(res[0]), nil
}

func (s *service) UpdateElement(ctx context.Context, catalogName string, elementID int, data *gkitmodels.CatalogElementData) (*ElementItem, error) {
	id, err := s.resolveCatalogName(catalogName)
	if err != nil {
		return nil, err
	}

	element := mapElementDataToModel(data)
	element.ID = elementID // BUG2 fix для элементов: явно проставляем ID
	res, _, err := s.sdk.CatalogElements(id).Update(ctx, []*models.CatalogElement{element})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return s.normalizeElement(res[0]), nil
}

func (s *service) DeleteElement(ctx context.Context, catalogName string, elementID int) error {
	id, err := s.resolveCatalogName(catalogName)
	if err != nil {
		return err
	}
	return s.sdk.CatalogElements(id).Delete(ctx, elementID)
}

func (s *service) LinkElement(ctx context.Context, catalogName string, elementID int, entityType string, entityID int, metadata map[string]interface{}) error {
	id, err := s.resolveCatalogName(catalogName)
	if err != nil {
		return err
	}
	return s.sdk.CatalogElements(id).Link(ctx, elementID, entityType, entityID, metadata)
}

func (s *service) UnlinkElement(ctx context.Context, catalogName string, elementID int, entityType string, entityID int) error {
	id, err := s.resolveCatalogName(catalogName)
	if err != nil {
		return err
	}
	return s.sdk.CatalogElements(id).Unlink(ctx, elementID, entityType, entityID)
}

// mapElementDataToModel преобразует CatalogElementData из tool input в SDK-модель.
// Исправляет BUG1 для элементов: явный маппинг вместо json roundtrip.
func mapElementDataToModel(data *gkitmodels.CatalogElementData) *models.CatalogElement {
	if data == nil {
		return &models.CatalogElement{}
	}

	element := &models.CatalogElement{
		Name: data.Name,
	}

	// Маппинг кастомных полей: field_code → CustomFieldValue
	if len(data.CustomFieldsValues) > 0 {
		cfv := make([]models.CustomFieldValue, 0, len(data.CustomFieldsValues))
		for _, f := range data.CustomFieldsValues {
			cfv = append(cfv, models.CustomFieldValue{
				FieldCode: f.FieldCode,
				Values: []models.FieldValueElement{
					{Value: f.Value},
				},
			})
		}
		element.CustomFieldsValues = cfv
	}

	return element
}

// mergeWith объединяет два среза строк без дублей (base — приоритет).
func mergeWith(base, extra []string) []string {
	seen := make(map[string]struct{}, len(base))
	result := make([]string, 0, len(base)+len(extra))
	for _, v := range base {
		seen[v] = struct{}{}
		result = append(result, v)
	}
	for _, v := range extra {
		if _, ok := seen[v]; !ok {
			result = append(result, v)
		}
	}
	return result
}
