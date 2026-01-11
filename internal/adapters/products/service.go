package products

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models"
)

type Service interface {
	SearchProducts(ctx context.Context, filter *gkitmodels.ProductFilter) ([]*models.CatalogElement, error)
	GetProduct(ctx context.Context, id int, with []string) (*models.CatalogElement, error)
	CreateProducts(ctx context.Context, elements []*models.CatalogElement) ([]*models.CatalogElement, error)
	UpdateProducts(ctx context.Context, elements []*models.CatalogElement) ([]*models.CatalogElement, error)
	DeleteProducts(ctx context.Context, ids []int) error

	// Работа со связями
	GetProductsByEntity(ctx context.Context, entityType string, entityID int) ([]models.EntityLink, error)
	LinkProduct(ctx context.Context, entityType string, entityID int, productID int, quantity int, priceID int) error
	UnlinkProduct(ctx context.Context, entityType string, entityID int, productID int) error
}

type service struct {
	sdk *amocrm.SDK
}

func NewService(sdk *amocrm.SDK) Service {
	return &service{sdk: sdk}
}
