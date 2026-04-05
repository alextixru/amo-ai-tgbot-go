package customers

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
)

func (s *service) ListTransactions(ctx context.Context, customerID int, page, limit int) (*TransactionsListOutput, error) {
	f := &services.TransactionsFilter{}
	if page > 0 {
		f.Page = page
	}
	if limit > 0 {
		f.Limit = limit
	}
	res, meta, err := s.sdk.CustomerTransactions(customerID).Get(ctx, f)
	if err != nil {
		return nil, err
	}
	out := &TransactionsListOutput{
		Transactions: make([]TransactionOutput, 0, len(res)),
	}
	for _, tx := range res {
		out.Transactions = append(out.Transactions, transactionToOutput(tx))
	}
	if meta != nil {
		out.HasMore = meta.HasMore
	}
	return out, nil
}

func (s *service) CreateTransactions(ctx context.Context, customerID int, price int, comment string, accrueBonus bool) (*TransactionsListOutput, error) {
	txs := []models.Transaction{
		{
			Price:   price,
			Comment: comment,
		},
	}
	res, err := s.sdk.CustomerTransactions(customerID).Create(ctx, txs, accrueBonus)
	if err != nil {
		return nil, err
	}
	out := &TransactionsListOutput{
		Transactions: make([]TransactionOutput, 0, len(res)),
	}
	for _, tx := range res {
		out.Transactions = append(out.Transactions, transactionToOutput(tx))
	}
	return out, nil
}

func (s *service) DeleteTransaction(ctx context.Context, customerID int, transactionID int) error {
	return s.sdk.CustomerTransactions(customerID).Delete(ctx, transactionID)
}

// transactionToOutput конвертирует models.Transaction в TransactionOutput.
func transactionToOutput(tx models.Transaction) TransactionOutput {
	return TransactionOutput{
		ID:        tx.ID,
		Price:     tx.Price,
		Comment:   tx.Comment,
		CreatedAt: toISO(int64(tx.CreatedAt)),
	}
}
