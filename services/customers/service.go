package customers

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

// Service определяет бизнес-логику для работы с покупателями, бонусами, статусами, транзакциями и сегментами.
type Service interface {
	// Customers
	ListCustomers(ctx context.Context) ([]models.Customer, error)
	GetCustomer(ctx context.Context, id int) (models.Customer, error)
	CreateCustomers(ctx context.Context, customers []models.Customer) ([]models.Customer, error)
	UpdateCustomers(ctx context.Context, customers []models.Customer) ([]models.Customer, error)
	DeleteCustomer(ctx context.Context, id int) error
	LinkCustomer(ctx context.Context, customerID int, links []models.EntityLink) ([]models.EntityLink, error)

	// Bonus Points
	GetBonusPoints(ctx context.Context, customerID int) (*models.BonusPoints, error)
	EarnBonusPoints(ctx context.Context, customerID int, points int) (int, error)
	RedeemBonusPoints(ctx context.Context, customerID int, points int) (int, error)

	// Statuses
	ListCustomerStatuses(ctx context.Context, page, limit int) ([]models.Status, error)
	GetCustomerStatus(ctx context.Context, id int) (models.Status, error)
	CreateCustomerStatuses(ctx context.Context, statuses []models.Status) ([]models.Status, error)
	UpdateCustomerStatuses(ctx context.Context, statuses []models.Status) ([]models.Status, error)
	DeleteCustomerStatus(ctx context.Context, id int) error

	// Transactions
	ListTransactions(ctx context.Context, customerID int) ([]models.Transaction, error)
	CreateTransactions(ctx context.Context, customerID int, transactions []models.Transaction, accrueBonus bool) ([]models.Transaction, error)
	DeleteTransaction(ctx context.Context, customerID int, transactionID int) error

	// Segments
	ListSegments(ctx context.Context) ([]*models.CustomerSegment, error)
	GetSegment(ctx context.Context, id int) (*models.CustomerSegment, error)
	CreateSegments(ctx context.Context, segments []*models.CustomerSegment) ([]*models.CustomerSegment, error)
	DeleteSegment(ctx context.Context, id int) error
}

type service struct {
	sdk *amocrm.SDK
}

// NewService создает новый экземпляр сервиса покупателей.
func NewService(sdk *amocrm.SDK) Service {
	return &service{
		sdk: sdk,
	}
}
