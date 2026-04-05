package catalogs

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// CatalogItem нормализованное представление каталога для LLM
type CatalogItem struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Type            string `json:"type,omitempty"`
	CanAddElements  bool   `json:"can_add_elements,omitempty"`
	CanShowInCards  bool   `json:"can_show_in_cards,omitempty"`
	CanLinkMultiple bool   `json:"can_link_multiple,omitempty"`
	CanBeDeleted    bool   `json:"can_be_deleted,omitempty"`
	Sort            int    `json:"sort,omitempty"`
	CreatedBy       string `json:"created_by,omitempty"`
	UpdatedBy       string `json:"updated_by,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}

// CatalogListResult результат списка каталогов с пагинацией
type CatalogListResult struct {
	Items   []*CatalogItem `json:"items"`
	Total   int            `json:"total"`
	Page    int            `json:"page"`
	HasMore bool           `json:"has_more"`
}

// ElementItem нормализованное представление элемента каталога для LLM
type ElementItem struct {
	ID                 int                        `json:"id"`
	Name               string                     `json:"name"`
	CatalogName        string                     `json:"catalog_name,omitempty"`
	CatalogID          int                        `json:"catalog_id,omitempty"`
	CurrencyCode       string                     `json:"currency_code,omitempty"`
	IsDeleted          bool                       `json:"is_deleted,omitempty"`
	CustomFieldsValues []models.CustomFieldValue  `json:"custom_fields_values,omitempty"`
	Quantity           float64                    `json:"quantity,omitempty"`
	InvoiceLink        string                     `json:"invoice_link,omitempty"`
	CreatedBy          string                     `json:"created_by,omitempty"`
	UpdatedBy          string                     `json:"updated_by,omitempty"`
	CreatedAt          string                     `json:"created_at,omitempty"`
	UpdatedAt          string                     `json:"updated_at,omitempty"`
}

// ElementListResult результат списка элементов с пагинацией
type ElementListResult struct {
	Items   []*ElementItem `json:"items"`
	Total   int            `json:"total"`
	Page    int            `json:"page"`
	HasMore bool           `json:"has_more"`
}

// Service определяет бизнес-логику для работы с каталогами и их элементами.
type Service interface {
	// CatalogNames возвращает список доступных каталогов для подсказок LLM
	CatalogNames() []string

	// Catalogs
	ListCatalogs(ctx context.Context, filter *gkitmodels.CatalogFilter) (*CatalogListResult, error)
	GetCatalog(ctx context.Context, name string) (*CatalogItem, error)
	CreateCatalog(ctx context.Context, data *gkitmodels.CatalogData) (*CatalogItem, error)
	UpdateCatalog(ctx context.Context, name string, data *gkitmodels.CatalogData) (*CatalogItem, error)
	DeleteCatalog(ctx context.Context, name string) error

	// Catalog Elements
	ListElements(ctx context.Context, catalogName string, filter *gkitmodels.CatalogFilter) (*ElementListResult, error)
	GetElement(ctx context.Context, catalogName string, elementID int, with []string) (*ElementItem, error)
	CreateElement(ctx context.Context, catalogName string, data *gkitmodels.CatalogElementData) (*ElementItem, error)
	UpdateElement(ctx context.Context, catalogName string, elementID int, data *gkitmodels.CatalogElementData) (*ElementItem, error)
	DeleteElement(ctx context.Context, catalogName string, elementID int) error
	LinkElement(ctx context.Context, catalogName string, elementID int, entityType string, entityID int, metadata map[string]interface{}) error
	UnlinkElement(ctx context.Context, catalogName string, elementID int, entityType string, entityID int) error
}

type service struct {
	sdk            *amocrm.SDK
	catalogsByName map[string]int    // имя → ID
	catalogsByID   map[int]string    // ID → имя
}

// New создает новый экземпляр сервиса каталогов и загружает справочник каталогов из SDK.
func New(ctx context.Context, sdk *amocrm.SDK) (Service, error) {
	s := &service{
		sdk:            sdk,
		catalogsByName: make(map[string]int),
		catalogsByID:   make(map[int]string),
	}
	if err := s.loadCatalogs(ctx); err != nil {
		return nil, fmt.Errorf("catalogs: init: %w", err)
	}
	return s, nil
}

// NewService создает сервис без загрузки каталогов (для обратной совместимости).
// Предпочтительно использовать New(ctx, sdk).
func NewService(sdk *amocrm.SDK) Service {
	return &service{
		sdk:            sdk,
		catalogsByName: make(map[string]int),
		catalogsByID:   make(map[int]string),
	}
}

// loadCatalogs загружает все каталоги из SDK и строит внутренние мапы.
func (s *service) loadCatalogs(ctx context.Context) error {
	catalogs, _, err := s.sdk.Catalogs().Get(ctx, nil)
	if err != nil {
		return err
	}
	for _, c := range catalogs {
		if c == nil {
			continue
		}
		s.catalogsByName[c.Name] = c.ID
		s.catalogsByID[c.ID] = c.Name
	}
	return nil
}

// CatalogNames возвращает отсортированный список доступных имён каталогов.
func (s *service) CatalogNames() []string {
	names := make([]string, 0, len(s.catalogsByName))
	for name := range s.catalogsByName {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// resolveCatalogName резолвит имя каталога в ID. Возвращает ошибку с подсказкой если не найдено.
func (s *service) resolveCatalogName(name string) (int, error) {
	id, ok := s.catalogsByName[name]
	if !ok {
		available := strings.Join(s.CatalogNames(), ", ")
		if available == "" {
			available = "(каталоги не загружены)"
		}
		return 0, fmt.Errorf("каталог %q не найден. Доступные: %s", name, available)
	}
	return id, nil
}

// resolveCatalogID резолвит ID каталога в имя. Возвращает "[unknown:ID]" если не найдено.
func (s *service) resolveCatalogID(id int) string {
	if name, ok := s.catalogsByID[id]; ok {
		return name
	}
	return fmt.Sprintf("[unknown:%d]", id)
}

// normalizeCatalog преобразует SDK-модель Catalog в CatalogItem для LLM.
func (s *service) normalizeCatalog(c *models.Catalog) *CatalogItem {
	if c == nil {
		return nil
	}
	item := &CatalogItem{
		ID:              c.ID,
		Name:            c.Name,
		Type:            string(c.Type),
		CanAddElements:  c.CanAddElements,
		CanShowInCards:  c.CanShowInCards,
		CanLinkMultiple: c.CanLinkMultiple,
		CanBeDeleted:    c.CanBeDeleted,
		Sort:            c.Sort,
		CreatedBy:       fmt.Sprintf("[unknown:%d]", c.CreatedBy),
		UpdatedBy:       fmt.Sprintf("[unknown:%d]", c.UpdatedBy),
	}
	if c.CreatedBy == 0 {
		item.CreatedBy = ""
	}
	if c.UpdatedBy == 0 {
		item.UpdatedBy = ""
	}
	if c.CreatedAt != 0 {
		item.CreatedAt = time.Unix(c.CreatedAt, 0).UTC().Format(time.RFC3339)
	}
	if c.UpdatedAt != 0 {
		item.UpdatedAt = time.Unix(c.UpdatedAt, 0).UTC().Format(time.RFC3339)
	}
	return item
}

// normalizeElement преобразует SDK-модель CatalogElement в ElementItem для LLM.
func (s *service) normalizeElement(e *models.CatalogElement) *ElementItem {
	if e == nil {
		return nil
	}
	item := &ElementItem{
		ID:                 e.ID,
		Name:               e.Name,
		CatalogID:          e.CatalogID,
		CatalogName:        s.resolveCatalogID(e.CatalogID),
		IsDeleted:          e.IsDeleted,
		CustomFieldsValues: e.CustomFieldsValues,
		Quantity:           e.Quantity,
		InvoiceLink:        e.InvoiceLink,
	}
	if e.CurrencyID != 0 {
		item.CurrencyCode = fmt.Sprintf("[unknown:%d]", e.CurrencyID)
	}
	if e.CreatedBy != 0 {
		item.CreatedBy = fmt.Sprintf("[unknown:%d]", e.CreatedBy)
	}
	if e.UpdatedBy != 0 {
		item.UpdatedBy = fmt.Sprintf("[unknown:%d]", e.UpdatedBy)
	}
	if e.CreatedAt != 0 {
		item.CreatedAt = time.Unix(e.CreatedAt, 0).UTC().Format(time.RFC3339)
	}
	if e.UpdatedAt != 0 {
		item.UpdatedAt = time.Unix(e.UpdatedAt, 0).UTC().Format(time.RFC3339)
	}
	return item
}
