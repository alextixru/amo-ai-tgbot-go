package customers

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListCustomers(ctx context.Context) ([]models.Customer, error) {
	customers, _, err := s.sdk.Customers().Get(ctx, nil)
	return customers, err
}

func (s *service) GetCustomer(ctx context.Context, id int) (models.Customer, error) {
	return s.sdk.Customers().GetOne(ctx, id)
}

func (s *service) CreateCustomers(ctx context.Context, customers []models.Customer) ([]models.Customer, error) {
	res, _, err := s.sdk.Customers().Create(ctx, customers)
	return res, err
}

func (s *service) UpdateCustomers(ctx context.Context, customers []models.Customer) ([]models.Customer, error) {
	res, _, err := s.sdk.Customers().Update(ctx, customers)
	return res, err
}

func (s *service) DeleteCustomer(ctx context.Context, id int) error {
	return s.sdk.Customers().Delete(ctx, id)
}

func (s *service) LinkCustomer(ctx context.Context, customerID int, links []models.EntityLink) ([]models.EntityLink, error) {
	return s.sdk.Customers().Link(ctx, customerID, links)
}
