package products

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// hardcoded popular currency codes by ID (amoCRM standard)
var currencyCodes = map[int]string{
	1:  "RUB",
	2:  "USD",
	3:  "EUR",
	4:  "GBP",
	5:  "UAH",
	6:  "KZT",
	7:  "BYR",
	8:  "CNY",
	9:  "TRY",
	10: "AED",
}

func (s *service) findProductsCatalogID(ctx context.Context) (int, error) {
	var findErr error
	s.catalogIDOnce.Do(func() {
		f := filters.NewCatalogsFilter().SetType("products")
		catalogs, _, err := s.sdk.Catalogs().Get(ctx, f)
		if err != nil {
			findErr = fmt.Errorf("failed to find products catalog: %w", err)
			return
		}
		if len(catalogs) == 0 {
			findErr = fmt.Errorf("products catalog not found")
			return
		}
		s.productsCatalogID = catalogs[0].ID
	})
	return s.productsCatalogID, findErr
}

// findFirstPriceFieldID загружает кастомные поля каталога и возвращает ID первого поля типа price
func (s *service) findFirstPriceFieldID(ctx context.Context, catalogID int) (int, error) {
	entityType := fmt.Sprintf("catalogs/%d", catalogID)
	fields, _, err := s.sdk.CustomFields().Get(ctx, entityType, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get catalog custom fields: %w", err)
	}
	for _, f := range fields {
		if f.Type == "price" {
			return f.ID, nil
		}
	}
	return 0, fmt.Errorf("no price field found in products catalog")
}

// resolveCurrencyCode резолвит числовой ID валюты в код (например 1 → "RUB")
func resolveCurrencyCode(id int) string {
	if id == 0 {
		return ""
	}
	if code, ok := currencyCodes[id]; ok {
		return code
	}
	return fmt.Sprintf("[unknown_currency:%d]", id)
}

// convertProductFields конвертирует []ProductFieldInput → []models.CustomFieldValue через field_code
func convertProductFields(fields []gkitmodels.ProductFieldInput) []models.CustomFieldValue {
	if len(fields) == 0 {
		return nil
	}
	result := make([]models.CustomFieldValue, 0, len(fields))
	for _, f := range fields {
		cfv := models.CustomFieldValue{
			FieldCode: f.FieldCode,
			Values: []models.FieldValueElement{
				{Value: f.Value},
			},
		}
		result = append(result, cfv)
	}
	return result
}

// productDataToElement конвертирует ProductData в *models.CatalogElement
func productDataToElement(d gkitmodels.ProductData) *models.CatalogElement {
	el := &models.CatalogElement{
		Name:               d.Name,
		CustomFieldsValues: convertProductFields(d.Fields),
	}
	if d.ID != 0 {
		el.ID = d.ID
	}
	return el
}

func (s *service) SearchProducts(ctx context.Context, filter *gkitmodels.ProductFilter, with []string) (*ProductSearchResult, error) {
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
		} else {
			f.SetLimit(50)
		}
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		} else {
			f.SetPage(1)
		}
	} else {
		f.SetLimit(50)
		f.SetPage(1)
	}

	// Прокидываем with в параметры запроса через BaseFilter.SetWith
	if len(with) > 0 {
		f.SetWith(with...)
	}

	elements, meta, err := s.sdk.CatalogElements(catalogID).Get(ctx, f)
	if err != nil {
		return nil, err
	}

	result := &ProductSearchResult{
		Items: elements,
	}
	if meta != nil {
		result.Page = meta.Page
		result.HasMore = meta.HasMore
	}
	return result, nil
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

func (s *service) CreateProducts(ctx context.Context, items []gkitmodels.ProductData) ([]*models.CatalogElement, error) {
	catalogID, err := s.findProductsCatalogID(ctx)
	if err != nil {
		return nil, err
	}
	elements := make([]*models.CatalogElement, 0, len(items))
	for _, d := range items {
		elements = append(elements, productDataToElement(d))
	}
	created, _, err := s.sdk.CatalogElements(catalogID).Create(ctx, elements)
	return created, err
}

func (s *service) UpdateProducts(ctx context.Context, items []gkitmodels.ProductData) ([]*models.CatalogElement, error) {
	catalogID, err := s.findProductsCatalogID(ctx)
	if err != nil {
		return nil, err
	}
	elements := make([]*models.CatalogElement, 0, len(items))
	for _, d := range items {
		elements = append(elements, productDataToElement(d))
	}
	updated, _, err := s.sdk.CatalogElements(catalogID).Update(ctx, elements)
	return updated, err
}

func (s *service) DeleteProducts(ctx context.Context, ids []int) (*OperationResult, error) {
	catalogID, err := s.findProductsCatalogID(ctx)
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		if err := s.sdk.CatalogElements(catalogID).Delete(ctx, id); err != nil {
			return nil, fmt.Errorf("failed to delete product %d: %w", id, err)
		}
	}
	return &OperationResult{OK: true, Action: "delete"}, nil
}

func (s *service) GetProductsByEntity(ctx context.Context, entityType string, entityID int) ([]ProductWithLink, error) {
	links, err := s.sdk.Links().Get(ctx, entityType, entityID, nil)
	if err != nil {
		return nil, err
	}

	catalogID, err := s.findProductsCatalogID(ctx)
	if err != nil {
		return nil, err
	}

	var result []ProductWithLink
	for _, link := range links {
		if link.ToEntityType != "catalog_elements" {
			continue
		}
		if cid, ok := link.Metadata["catalog_id"].(float64); !ok || int(cid) != catalogID {
			continue
		}

		product, err := s.GetProduct(ctx, link.ToEntityID, []string{"supplier_field_values"})
		if err != nil {
			// Если не удалось загрузить товар — добавляем с пустым Product, не ломаем ответ
			result = append(result, ProductWithLink{Link: *link, Product: nil})
			continue
		}

		// Резолвим числовые ID в читаемые значения
		currCode := resolveCurrencyCode(product.CurrencyID)

		result = append(result, ProductWithLink{Link: *link, Product: product, CurrencyCode: currCode})
	}
	return result, nil
}

func (s *service) LinkProduct(ctx context.Context, entityType string, entityID int, productID int, quantity int, priceID int) (*OperationResult, error) {
	catalogID, err := s.findProductsCatalogID(ctx)
	if err != nil {
		return nil, err
	}

	// Если price_id не передан — берём первое ценовое поле каталога
	if priceID == 0 {
		priceID, err = s.findFirstPriceFieldID(ctx, catalogID)
		if err != nil {
			// Не критично — продолжаем без price_id
			priceID = 0
		}
	}

	link := models.NewCatalogElementLink(entityType, entityID, catalogID, productID, float64(quantity), priceID)
	_, err = s.sdk.Links().Link(ctx, entityType, entityID, []*models.EntityLink{link})
	if err != nil {
		return nil, err
	}
	return &OperationResult{OK: true, Action: "link", ProductID: productID, EntityID: entityID}, nil
}

func (s *service) UnlinkProduct(ctx context.Context, entityType string, entityID int, productID int) (*OperationResult, error) {
	link := &models.EntityLink{
		ToEntityType: "catalog_elements",
		ToEntityID:   productID,
	}
	if err := s.sdk.Links().Unlink(ctx, entityType, entityID, []*models.EntityLink{link}); err != nil {
		return nil, err
	}
	return &OperationResult{OK: true, Action: "unlink", ProductID: productID, EntityID: entityID}, nil
}
