package tools

import (
	"fmt"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func (r *Registry) RegisterCustomersTool() {
	r.addTool(genkit.DefineTool[gkitmodels.CustomersInput, any](
		r.g,
		"customers",
		"Work with customers (Retention)",
		func(ctx *ai.ToolContext, input gkitmodels.CustomersInput) (any, error) {
			switch input.Layer {
			case "customers":
				switch input.Action {
				case "list":
					return r.customersService.ListCustomers(ctx)
				case "get":
					if input.ID == 0 {
						return nil, fmt.Errorf("id is required")
					}
					return r.customersService.GetCustomer(ctx, input.ID)
				case "create":
					if input.Data == nil {
						return nil, fmt.Errorf("data is required")
					}
					customer := amomodels.Customer{
						Name:      input.Data.Name,
						NextPrice: input.Data.NextPrice,
						NextDate:  input.Data.NextDate,
						StatusID:  input.Data.StatusID,
					}
					customer.ResponsibleUserID = input.Data.ResponsibleUserID
					return r.customersService.CreateCustomers(ctx, []amomodels.Customer{customer})
				case "update":
					if input.ID == 0 || input.Data == nil {
						return nil, fmt.Errorf("id and data are required")
					}
					customer := amomodels.Customer{
						Name:      input.Data.Name,
						NextPrice: input.Data.NextPrice,
						NextDate:  input.Data.NextDate,
						StatusID:  input.Data.StatusID,
					}
					customer.ID = input.ID
					customer.ResponsibleUserID = input.Data.ResponsibleUserID
					return r.customersService.UpdateCustomers(ctx, []amomodels.Customer{customer})
				case "delete":
					if input.ID == 0 {
						return nil, fmt.Errorf("id is required")
					}
					return nil, r.customersService.DeleteCustomer(ctx, input.ID)
				case "link":
					if input.CustomerID == 0 || input.LinkData == nil {
						return nil, fmt.Errorf("customer_id and link_data are required")
					}
					link := amomodels.EntityLink{
						ToEntityType: input.LinkData.EntityType,
						ToEntityID:   input.LinkData.EntityID,
					}
					return r.customersService.LinkCustomer(ctx, input.CustomerID, []amomodels.EntityLink{link})
				default:
					return nil, fmt.Errorf("unknown action for customers: %s", input.Action)
				}

			case "bonus_points":
				if input.CustomerID == 0 {
					return nil, fmt.Errorf("customer_id is required")
				}
				switch input.Action {
				case "get":
					return r.customersService.GetBonusPoints(ctx, input.CustomerID)
				case "earn_points":
					return r.customersService.EarnBonusPoints(ctx, input.CustomerID, input.Points)
				case "redeem_points":
					return r.customersService.RedeemBonusPoints(ctx, input.CustomerID, input.Points)
				default:
					return nil, fmt.Errorf("unknown action for bonus_points: %s", input.Action)
				}

			case "statuses":
				switch input.Action {
				case "list":
					page, limit := 0, 0
					if input.Filter != nil {
						page, limit = input.Filter.Page, input.Filter.Limit
					}
					return r.customersService.ListCustomerStatuses(ctx, page, limit)
				case "get":
					if input.ID == 0 {
						return nil, fmt.Errorf("id is required")
					}
					return r.customersService.GetCustomerStatus(ctx, input.ID)
				case "create":
					if input.Data == nil {
						return nil, fmt.Errorf("data is required")
					}
					status := amomodels.Status{
						Name: input.Data.Name,
					}
					return r.customersService.CreateCustomerStatuses(ctx, []amomodels.Status{status})
				case "update":
					if input.ID == 0 || input.Data == nil {
						return nil, fmt.Errorf("id and data are required")
					}
					status := amomodels.Status{
						Name: input.Data.Name,
					}
					status.ID = input.ID
					return r.customersService.UpdateCustomerStatuses(ctx, []amomodels.Status{status})
				case "delete":
					if input.ID == 0 {
						return nil, fmt.Errorf("id is required")
					}
					return nil, r.customersService.DeleteCustomerStatus(ctx, input.ID)
				default:
					return nil, fmt.Errorf("unknown action for statuses: %s", input.Action)
				}

			case "transactions":
				if input.CustomerID == 0 {
					return nil, fmt.Errorf("customer_id is required")
				}
				switch input.Action {
				case "list":
					return r.customersService.ListTransactions(ctx, input.CustomerID)
				case "create":
					if input.TransactionData == nil {
						return nil, fmt.Errorf("transaction_data is required")
					}
					tx := amomodels.Transaction{
						Price:   input.TransactionData.Price,
						Comment: input.TransactionData.Comment,
					}
					return r.customersService.CreateTransactions(ctx, input.CustomerID, []amomodels.Transaction{tx}, input.TransactionData.AccrueBonus)
				case "delete":
					if input.ID == 0 {
						return nil, fmt.Errorf("id is required")
					}
					return nil, r.customersService.DeleteTransaction(ctx, input.CustomerID, input.ID)
				default:
					return nil, fmt.Errorf("unknown action for transactions: %s", input.Action)
				}

			case "segments":
				switch input.Action {
				case "list":
					return r.customersService.ListSegments(ctx)
				case "get":
					if input.ID == 0 {
						return nil, fmt.Errorf("id is required")
					}
					return r.customersService.GetSegment(ctx, input.ID)
				case "create":
					if input.Data == nil {
						return nil, fmt.Errorf("data is required")
					}
					segment := &amomodels.CustomerSegment{
						Name: input.Data.Name,
					}
					return r.customersService.CreateSegments(ctx, []*amomodels.CustomerSegment{segment})
				case "delete":
					if input.ID == 0 {
						return nil, fmt.Errorf("id is required")
					}
					return nil, r.customersService.DeleteSegment(ctx, input.ID)
				default:
					return nil, fmt.Errorf("unknown action for segments: %s", input.Action)
				}

			default:
				return nil, fmt.Errorf("unknown layer: %s", input.Layer)
			}
		},
	))
}
