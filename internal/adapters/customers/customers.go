package customers

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListCustomers(ctx context.Context, filter *filters.CustomersFilter, with []string) ([]models.Customer, error) {
	if filter == nil {
		filter = &filters.CustomersFilter{}
	}
	if len(with) > 0 {
		filter.With = with
	}
	customers, _, err := s.sdk.Customers().Get(ctx, filter)
	return customers, err
}

func (s *service) GetCustomer(ctx context.Context, id int, with []string) (models.Customer, error) {
	filter := &filters.CustomersFilter{}
	filter.IDs = []int{id}
	if len(with) > 0 {
		filter.With = with
	}
	customers, _, err := s.sdk.Customers().Get(ctx, filter)
	if err != nil {
		return models.Customer{}, err
	}
	if len(customers) > 0 {
		return customers[0], nil
	}
	return models.Customer{}, nil
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
