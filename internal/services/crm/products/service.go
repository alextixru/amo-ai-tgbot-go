package products

import (
	"context"
	"sync"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// ProductWithLink обогащённая связь: метаданные связи + детали товара
type ProductWithLink struct {
	Link         models.EntityLink      `json:"link"`
	Product      *models.CatalogElement `json:"product"`
	CurrencyCode string                 `json:"currency_code,omitempty"`
}

// ProductSearchResult результат поиска товаров с метаданными пагинации
type ProductSearchResult struct {
	Items   []*models.CatalogElement `json:"items"`
	Page    int                      `json:"page"`
	HasMore bool                     `json:"has_more"`
}

// OperationResult подтверждение выполненной операции
type OperationResult struct {
	OK        bool   `json:"ok"`
	Action    string `json:"action"`
	ProductID int    `json:"product_id,omitempty"`
	EntityID  int    `json:"entity_id,omitempty"`
}

// Service интерфейс сервиса товаров
type Service interface {
	SearchProducts(ctx context.Context, filter *gkitmodels.ProductFilter, with []string) (*ProductSearchResult, error)
	GetProduct(ctx context.Context, id int, with []string) (*models.CatalogElement, error)
	CreateProducts(ctx context.Context, items []gkitmodels.ProductData) ([]*models.CatalogElement, error)
	UpdateProducts(ctx context.Context, items []gkitmodels.ProductData) ([]*models.CatalogElement, error)
	DeleteProducts(ctx context.Context, ids []int) (*OperationResult, error)

	// Работа со связями
	GetProductsByEntity(ctx context.Context, entityType string, entityID int) ([]ProductWithLink, error)
	LinkProduct(ctx context.Context, entityType string, entityID int, productID int, quantity int, priceID int) (*OperationResult, error)
	UnlinkProduct(ctx context.Context, entityType string, entityID int, productID int) (*OperationResult, error)
}

type service struct {
	sdk               *amocrm.SDK
	catalogIDOnce     sync.Once
	productsCatalogID int
}

func NewService(sdk *amocrm.SDK) Service {
	return &service{sdk: sdk}
}
