package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// CustomerToolInput входные параметры для инструмента crm_customers
type CustomerToolInput struct {
	// Action действие: search, get, create, update
	Action string `json:"action"`

	// Filter фильтр для поиска (используется при action=search)
	Filter *CustomerFilterInput `json:"filter,omitempty"`

	// ID идентификатор покупателя (используется при action=get, update)
	ID int `json:"id,omitempty"`

	// Data данные для создания или обновления (используется при action=create, update)
	Data *CustomerDataInput `json:"data,omitempty"`
}

// CustomerFilterInput параметры фильтрации покупателей
type CustomerFilterInput struct {
	// Query поисковый запрос
	Query string `json:"query,omitempty"`
	// ResponsibleUserID ID ответственного
	ResponsibleUserID int `json:"responsible_user_id,omitempty"`
	// Limit количество элементов (по умолчанию 50)
	Limit int `json:"limit,omitempty"`
}

// CustomerDataInput данные покупателя
type CustomerDataInput struct {
	// Name имя покупателя (обязательно при создании)
	Name string `json:"name,omitempty"`
	// NextPrice ожидаемая сумма следующей покупки
	NextPrice int `json:"next_price,omitempty"`
	// Periodicity периодичность покупок (в днях)
	Periodicity int `json:"periodicity,omitempty"`
	// ResponsibleUserID ID ответственного
	ResponsibleUserID int `json:"responsible_user_id,omitempty"`
}

// registerCustomersTools регистрирует инструменты для работы с покупателями
func (r *Registry) registerCustomersTools() {
	r.addTool(genkit.DefineTool[CustomerToolInput, any](
		r.g,
		"crm_customers",
		"Управление покупателями (Customers). Позволяет искать, получать, создавать и обновлять покупателей в amoCRM.",
		func(ctx *ai.ToolContext, input CustomerToolInput) (any, error) {
			switch input.Action {
			case "search":
				return r.searchCustomers(ctx.Context, input.Filter)
			case "get":
				if input.ID == 0 {
					return nil, fmt.Errorf("id is required for 'get' action")
				}
				return r.getCustomer(ctx.Context, input.ID)
			case "create":
				if input.Data == nil {
					return nil, fmt.Errorf("data is required for 'create' action")
				}
				return r.createCustomer(ctx.Context, input.Data)
			case "update":
				if input.ID == 0 {
					return nil, fmt.Errorf("id is required for 'update' action")
				}
				if input.Data == nil {
					return nil, fmt.Errorf("data is required for 'update' action")
				}
				return r.updateCustomer(ctx.Context, input.ID, input.Data)
			default:
				return nil, fmt.Errorf("unknown action: %s", input.Action)
			}
		},
	))
}

// searchCustomers выполняет поиск покупателей
func (r *Registry) searchCustomers(ctx context.Context, filter *CustomerFilterInput) ([]models.Customer, error) {
	sdkFilter := &services.CustomersFilter{
		Limit: 50,
	}

	if filter != nil {
		if filter.Limit > 0 {
			sdkFilter.Limit = filter.Limit
		}
		if filter.Query != "" {
			sdkFilter.Query = filter.Query
		}
		if filter.ResponsibleUserID > 0 {
			sdkFilter.ResponsibleUserIDs = []int{filter.ResponsibleUserID}
		}
	}

	return r.sdk.Customers().Get(ctx, sdkFilter)
}

// getCustomer получает покупателя по ID
func (r *Registry) getCustomer(ctx context.Context, id int) (*models.Customer, error) {
	return r.sdk.Customers().GetOne(ctx, id, []string{"contacts", "companies"})
}

// createCustomer создает нового покупателя
func (r *Registry) createCustomer(ctx context.Context, data *CustomerDataInput) (*models.Customer, error) {
	customer := models.Customer{
		Name:        data.Name,
		NextPrice:   data.NextPrice,
		Periodicity: data.Periodicity,
	}

	if data.ResponsibleUserID > 0 {
		customer.ResponsibleUserID = data.ResponsibleUserID
	}

	createdCustomers, err := r.sdk.Customers().Create(ctx, []models.Customer{customer})
	if err != nil {
		return nil, err
	}

	if len(createdCustomers) == 0 {
		return nil, fmt.Errorf("failed to create customer: empty response")
	}

	return &createdCustomers[0], nil
}

// updateCustomer обновляет существующего покупателя
func (r *Registry) updateCustomer(ctx context.Context, id int, data *CustomerDataInput) (*models.Customer, error) {
	customer := models.Customer{
		BaseModel: models.BaseModel{ID: id},
	}

	if data.Name != "" {
		customer.Name = data.Name
	}
	if data.NextPrice != 0 {
		customer.NextPrice = data.NextPrice
	}
	if data.Periodicity != 0 {
		customer.Periodicity = data.Periodicity
	}
	if data.ResponsibleUserID != 0 {
		customer.ResponsibleUserID = data.ResponsibleUserID
	}

	updatedCustomers, err := r.sdk.Customers().Update(ctx, []models.Customer{customer})
	if err != nil {
		return nil, err
	}

	if len(updatedCustomers) == 0 {
		return nil, fmt.Errorf("failed to update customer: empty response")
	}

	return &updatedCustomers[0], nil
}
