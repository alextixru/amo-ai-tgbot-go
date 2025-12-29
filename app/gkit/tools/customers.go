package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// CustomersInput входные параметры для инструмента customers
type CustomersInput struct {
	// Layer слой: customers, bonus_points, statuses, transactions, segments
	Layer string `json:"layer" jsonschema_description:"Слой: customers, bonus_points, statuses, transactions, segments"`

	// Action действие (зависит от layer)
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, delete, link, earn_points, redeem_points, etc."`

	// CustomerID ID покупателя (для большинства операций)
	CustomerID int `json:"customer_id,omitempty" jsonschema_description:"ID покупателя"`

	// ID идентификатор объекта (для get, delete)
	ID int `json:"id,omitempty" jsonschema_description:"ID объекта (статус, транзакция, сегмент)"`

	// Filter параметры поиска
	Filter *CustomerFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска"`

	// Data данные покупателя (для create, update)
	Data *CustomerData `json:"data,omitempty" jsonschema_description:"Данные покупателя"`

	// Points количество баллов (для earn_points, redeem_points)
	Points int `json:"points,omitempty" jsonschema_description:"Количество бонусных баллов"`

	// TransactionData данные транзакции
	TransactionData *CustomerTransactionData `json:"transaction_data,omitempty" jsonschema_description:"Данные транзакции"`

	// LinkData данные для привязки
	LinkData *CustomerLinkData `json:"link_data,omitempty" jsonschema_description:"Данные для привязки сущностей"`
}

// CustomerFilter фильтры поиска покупателей
type CustomerFilter struct {
	Page               int    `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Limit              int    `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Query              string `json:"query,omitempty" jsonschema_description:"Поисковый запрос"`
	ResponsibleUserIDs []int  `json:"responsible_user_ids,omitempty" jsonschema_description:"ID ответственных"`
}

// CustomerData данные покупателя
type CustomerData struct {
	Name              string `json:"name" jsonschema_description:"Имя покупателя"`
	ResponsibleUserID int    `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственного"`
	NextDate          int64  `json:"next_date,omitempty" jsonschema_description:"Дата следующей покупки (Unix timestamp)"`
	NextPrice         int    `json:"next_price,omitempty" jsonschema_description:"Ожидаемая сумма"`
	StatusID          int    `json:"status_id,omitempty" jsonschema_description:"ID статуса"`
}

// CustomerTransactionData данные транзакции
type CustomerTransactionData struct {
	Price       int    `json:"price" jsonschema_description:"Сумма транзакции"`
	Comment     string `json:"comment,omitempty" jsonschema_description:"Комментарий"`
	AccrueBonus bool   `json:"accrue_bonus,omitempty" jsonschema_description:"Начислить бонусные баллы"`
}

// CustomerLinkData данные для привязки
type CustomerLinkData struct {
	EntityType string `json:"entity_type" jsonschema_description:"Тип сущности: contacts, companies"`
	EntityID   int    `json:"entity_id" jsonschema_description:"ID сущности"`
}

// registerCustomersTool регистрирует инструмент для работы с покупателями
func (r *Registry) registerCustomersTool() {
	r.addTool(genkit.DefineTool[CustomersInput, any](
		r.g,
		"customers",
		"Работа с покупателями (retention). "+
			"Layers: customers (CRUD+link), bonus_points (earn/redeem), "+
			"statuses (list), transactions (create/list/delete), segments (list). "+
			"Требует customer_id для большинства операций.",
		func(ctx *ai.ToolContext, input CustomersInput) (any, error) {
			return r.handleCustomers(ctx.Context, input)
		},
	))
}

func (r *Registry) handleCustomers(ctx context.Context, input CustomersInput) (any, error) {
	switch input.Layer {
	case "customers":
		return r.handleCustomersLayer(ctx, input)
	case "bonus_points":
		return r.handleBonusPointsLayer(ctx, input)
	case "statuses":
		return r.handleStatusesLayer(ctx, input)
	case "transactions":
		return r.handleTransactionsLayer(ctx, input)
	case "segments":
		return r.handleSegmentsLayer(ctx, input)
	default:
		return nil, fmt.Errorf("unknown layer: %s", input.Layer)
	}
}

// ============ CUSTOMERS LAYER ============

func (r *Registry) handleCustomersLayer(ctx context.Context, input CustomersInput) (any, error) {
	switch input.Action {
	case "list":
		return r.listCustomers(ctx, input.Filter)
	case "get":
		if input.CustomerID == 0 {
			return nil, fmt.Errorf("customer_id is required for action 'get'")
		}
		return r.sdk.Customers().GetOne(ctx, input.CustomerID, nil)
	case "create":
		if input.Data == nil || input.Data.Name == "" {
			return nil, fmt.Errorf("data.name is required for action 'create'")
		}
		return r.createCustomer(ctx, input.Data)
	case "update":
		if input.CustomerID == 0 {
			return nil, fmt.Errorf("customer_id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		return r.updateCustomer(ctx, input.CustomerID, input.Data)
	case "delete":
		if input.CustomerID == 0 {
			return nil, fmt.Errorf("customer_id is required for action 'delete'")
		}
		err := r.sdk.Customers().Delete(ctx, input.CustomerID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true, "deleted_id": input.CustomerID}, nil
	case "link":
		if input.CustomerID == 0 {
			return nil, fmt.Errorf("customer_id is required for action 'link'")
		}
		if input.LinkData == nil {
			return nil, fmt.Errorf("link_data is required for action 'link'")
		}
		return r.linkCustomer(ctx, input.CustomerID, input.LinkData)
	default:
		return nil, fmt.Errorf("unknown action for customers layer: %s", input.Action)
	}
}

func (r *Registry) listCustomers(ctx context.Context, filter *CustomerFilter) ([]models.Customer, error) {
	f := &services.CustomersFilter{
		Limit: 50,
		Page:  1,
	}
	if filter != nil {
		if filter.Limit > 0 {
			f.Limit = filter.Limit
		}
		if filter.Page > 0 {
			f.Page = filter.Page
		}
		if filter.Query != "" {
			f.Query = filter.Query
		}
		if len(filter.ResponsibleUserIDs) > 0 {
			f.ResponsibleUserIDs = filter.ResponsibleUserIDs
		}
	}
	return r.sdk.Customers().Get(ctx, f)
}

func (r *Registry) createCustomer(ctx context.Context, data *CustomerData) ([]models.Customer, error) {
	customer := models.Customer{
		Name: data.Name,
	}
	if data.ResponsibleUserID > 0 {
		customer.ResponsibleUserID = data.ResponsibleUserID
	}
	if data.NextDate > 0 {
		customer.NextDate = data.NextDate
	}
	if data.NextPrice > 0 {
		customer.NextPrice = data.NextPrice
	}
	if data.StatusID > 0 {
		customer.StatusID = data.StatusID
	}
	return r.sdk.Customers().Create(ctx, []models.Customer{customer})
}

func (r *Registry) updateCustomer(ctx context.Context, id int, data *CustomerData) ([]models.Customer, error) {
	customer := models.Customer{}
	customer.ID = id
	if data.Name != "" {
		customer.Name = data.Name
	}
	if data.ResponsibleUserID > 0 {
		customer.ResponsibleUserID = data.ResponsibleUserID
	}
	if data.NextDate > 0 {
		customer.NextDate = data.NextDate
	}
	if data.NextPrice > 0 {
		customer.NextPrice = data.NextPrice
	}
	if data.StatusID > 0 {
		customer.StatusID = data.StatusID
	}
	return r.sdk.Customers().Update(ctx, []models.Customer{customer})
}

func (r *Registry) linkCustomer(ctx context.Context, customerID int, data *CustomerLinkData) ([]models.EntityLink, error) {
	link := models.EntityLink{
		ToEntityType: data.EntityType,
		ToEntityID:   data.EntityID,
	}
	return r.sdk.Customers().Link(ctx, customerID, []models.EntityLink{link})
}

// ============ BONUS POINTS LAYER ============

func (r *Registry) handleBonusPointsLayer(ctx context.Context, input CustomersInput) (any, error) {
	if input.CustomerID == 0 {
		return nil, fmt.Errorf("customer_id is required for bonus_points layer")
	}

	switch input.Action {
	case "get_points":
		return r.sdk.CustomerBonusPoints(input.CustomerID).Get(ctx)
	case "earn_points":
		if input.Points <= 0 {
			return nil, fmt.Errorf("points must be positive for action 'earn_points'")
		}
		newBalance, err := r.sdk.CustomerBonusPoints(input.CustomerID).EarnPoints(ctx, input.Points)
		if err != nil {
			return nil, err
		}
		return map[string]any{"bonus_points": newBalance}, nil
	case "redeem_points":
		if input.Points <= 0 {
			return nil, fmt.Errorf("points must be positive for action 'redeem_points'")
		}
		newBalance, err := r.sdk.CustomerBonusPoints(input.CustomerID).RedeemPoints(ctx, input.Points)
		if err != nil {
			return nil, err
		}
		return map[string]any{"bonus_points": newBalance}, nil
	default:
		return nil, fmt.Errorf("unknown action for bonus_points layer: %s", input.Action)
	}
}

// ============ STATUSES LAYER ============

func (r *Registry) handleStatusesLayer(ctx context.Context, input CustomersInput) (any, error) {
	switch input.Action {
	case "list_statuses":
		page := 1
		limit := 50
		if input.Filter != nil {
			if input.Filter.Page > 0 {
				page = input.Filter.Page
			}
			if input.Filter.Limit > 0 {
				limit = input.Filter.Limit
			}
		}
		return r.sdk.CustomerStatuses().Get(ctx, page, limit)
	default:
		return nil, fmt.Errorf("unknown action for statuses layer: %s", input.Action)
	}
}

// ============ TRANSACTIONS LAYER ============

func (r *Registry) handleTransactionsLayer(ctx context.Context, input CustomersInput) (any, error) {
	if input.CustomerID == 0 {
		return nil, fmt.Errorf("customer_id is required for transactions layer")
	}

	switch input.Action {
	case "list_transactions":
		f := &services.TransactionsFilter{
			Limit: 50,
			Page:  1,
		}
		if input.Filter != nil {
			if input.Filter.Limit > 0 {
				f.Limit = input.Filter.Limit
			}
			if input.Filter.Page > 0 {
				f.Page = input.Filter.Page
			}
		}
		return r.sdk.CustomerTransactions(input.CustomerID).Get(ctx, f)
	case "create_transaction":
		if input.TransactionData == nil || input.TransactionData.Price == 0 {
			return nil, fmt.Errorf("transaction_data.price is required for action 'create_transaction'")
		}
		transaction := services.Transaction{
			Price:   input.TransactionData.Price,
			Comment: input.TransactionData.Comment,
		}
		accrueBonus := true
		if input.TransactionData.AccrueBonus {
			accrueBonus = input.TransactionData.AccrueBonus
		}
		return r.sdk.CustomerTransactions(input.CustomerID).Create(ctx, []services.Transaction{transaction}, accrueBonus)
	case "delete_transaction":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'delete_transaction'")
		}
		err := r.sdk.CustomerTransactions(input.CustomerID).Delete(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true, "deleted_id": input.ID}, nil
	default:
		return nil, fmt.Errorf("unknown action for transactions layer: %s", input.Action)
	}
}

// ============ SEGMENTS LAYER ============

func (r *Registry) handleSegmentsLayer(ctx context.Context, input CustomersInput) (any, error) {
	switch input.Action {
	case "list_segments":
		f := &services.SegmentsFilter{
			Limit: 50,
			Page:  1,
		}
		if input.Filter != nil {
			if input.Filter.Limit > 0 {
				f.Limit = input.Filter.Limit
			}
			if input.Filter.Page > 0 {
				f.Page = input.Filter.Page
			}
		}
		return r.sdk.Segments().Get(ctx, f)
	case "get_segment":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get_segment'")
		}
		return r.sdk.Segments().GetOne(ctx, input.ID)
	default:
		return nil, fmt.Errorf("unknown action for segments layer: %s", input.Action)
	}
}
