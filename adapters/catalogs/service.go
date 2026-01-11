package catalogs

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

// Service определяет бизнес-логику для работы с каталогами и их элементами.
type Service interface {
	// Catalogs
	ListCatalogs(ctx context.Context) ([]*models.Catalog, error)
	GetCatalog(ctx context.Context, id int) (*models.Catalog, error)
	CreateCatalogs(ctx context.Context, catalogs []*models.Catalog) ([]*models.Catalog, error)
	UpdateCatalogs(ctx context.Context, catalogs []*models.Catalog) ([]*models.Catalog, error)

	// Catalog Elements
	ListElements(ctx context.Context, catalogID int, filter *gkitmodels.CatalogFilter) ([]*models.CatalogElement, error)
	GetElement(ctx context.Context, catalogID, elementID int, with []string) (*models.CatalogElement, error)
	CreateElements(ctx context.Context, catalogID int, elements []*models.CatalogElement) ([]*models.CatalogElement, error)
	UpdateElements(ctx context.Context, catalogID int, elements []*models.CatalogElement) ([]*models.CatalogElement, error)
	LinkElement(ctx context.Context, catalogID, elementID int, entityType string, entityID int, metadata map[string]interface{}) error
	UnlinkElement(ctx context.Context, catalogID, elementID int, entityType string, entityID int) error
}

type service struct {
	sdk *amocrm.SDK
}

// NewService создает новый экземпляр сервиса каталогов.
func NewService(sdk *amocrm.SDK) Service {
	return &service{
		sdk: sdk,
	}
}
