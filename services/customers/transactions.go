package customers

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListTransactions(ctx context.Context, customerID int) ([]models.Transaction, error) {
	res, _, err := s.sdk.CustomerTransactions(customerID).Get(ctx, nil)
	return res, err
}

func (s *service) CreateTransactions(ctx context.Context, customerID int, transactions []models.Transaction, accrueBonus bool) ([]models.Transaction, error) {
	return s.sdk.CustomerTransactions(customerID).Create(ctx, transactions, accrueBonus)
}

func (s *service) DeleteTransaction(ctx context.Context, customerID int, transactionID int) error {
	return s.sdk.CustomerTransactions(customerID).Delete(ctx, transactionID)
}
