package tools

import (
	"encoding/json"
	"fmt"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
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
					var sdkFilter *filters.CustomersFilter
					var with []string
					if input.Filter != nil {
						sdkFilter = &filters.CustomersFilter{}
						if input.Filter.Page > 0 {
							sdkFilter.SetPage(input.Filter.Page)
						}
						if input.Filter.Limit > 0 {
							sdkFilter.SetLimit(input.Filter.Limit)
						}
						if input.Filter.Query != "" {
							sdkFilter.SetQuery(input.Filter.Query)
						}
						if len(input.Filter.ResponsibleUserIDs) > 0 {
							sdkFilter.SetResponsibleUserIDs(input.Filter.ResponsibleUserIDs)
						}
						if len(input.Filter.IDs) > 0 {
							sdkFilter.SetIDs(input.Filter.IDs)
						}
						if len(input.Filter.StatusIDs) > 0 {
							sdkFilter.SetStatusIDs(input.Filter.StatusIDs)
						}
						if len(input.Filter.Names) > 0 {
							sdkFilter.SetNames(input.Filter.Names)
						}
						// NextDate filter is NOT supported by the current SDK version.
						// Supported ranges: CreatedAt, UpdatedAt, ClosestTaskAt.
						with = input.Filter.With
					}
					return r.customersService.ListCustomers(ctx, sdkFilter, with)
				case "get":
					if input.ID == 0 {
						return nil, fmt.Errorf("id is required")
					}
					var with []string
					if input.Filter != nil {
						with = input.Filter.With
					}
					return r.customersService.GetCustomer(ctx, input.ID, with)
				case "create":
					if input.Batch != nil && len(input.Batch) > 0 {
						customers := make([]amomodels.Customer, 0, len(input.Batch))
						for _, d := range input.Batch {
							customers = append(customers, mapCustomerData(d, 0))
						}
						return r.customersService.CreateCustomers(ctx, customers)
					}
					if input.Data == nil {
						return nil, fmt.Errorf("data or batch is required")
					}
					return r.customersService.CreateCustomers(ctx, []amomodels.Customer{mapCustomerData(input.Data, 0)})
				case "update":
					if input.Batch != nil && len(input.Batch) > 0 {
						customers := make([]amomodels.Customer, 0, len(input.Batch))
						for _, d := range input.Batch {
							customers = append(customers, mapCustomerData(d, input.ID))
						}
						return r.customersService.UpdateCustomers(ctx, customers)
					}
					if input.ID == 0 || input.Data == nil {
						return nil, fmt.Errorf("id and data are required")
					}
					return r.customersService.UpdateCustomers(ctx, []amomodels.Customer{mapCustomerData(input.Data, input.ID)})
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

// mapCustomerData конвертирует CustomerData в SDK Customer
func mapCustomerData(d *gkitmodels.CustomerData, id int) amomodels.Customer {
	customer := amomodels.Customer{
		Name:      d.Name,
		NextPrice: d.NextPrice,
		NextDate:  d.NextDate,
		StatusID:  d.StatusID,
	}
	customer.ID = id
	customer.ResponsibleUserID = d.ResponsibleUserID

	if d.CustomFieldsValues != nil {
		customer.CustomFieldsValues = mapCustomerCustomFieldsValues(d.CustomFieldsValues)
	}

	if len(d.TagsToAdd) > 0 {
		for _, tag := range d.TagsToAdd {
			customer.TagsToAdd = append(customer.TagsToAdd, amomodels.Tag{Name: tag})
		}
	}
	if len(d.TagsToDelete) > 0 {
		for _, tag := range d.TagsToDelete {
			customer.TagsToDelete = append(customer.TagsToDelete, amomodels.Tag{Name: tag})
		}
	}

	return customer
}

// mapCustomerCustomFieldsValues конвертирует map[string]any в []CustomFieldValue
func mapCustomerCustomFieldsValues(cfv map[string]any) []amomodels.CustomFieldValue {
	// Пробуем десериализовать через JSON для гибкости (как в complex_create.go)
	data, err := json.Marshal(cfv)
	if err != nil {
		return nil
	}

	var result []amomodels.CustomFieldValue
	if err := json.Unmarshal(data, &result); err != nil {
		// Fallback: пробуем как map field_id -> values
		for fieldID, values := range cfv {
			cfValue := amomodels.CustomFieldValue{
				FieldCode: fieldID,
			}
			// Пытаемся распарсить values как массив
			if valArr, ok := values.([]any); ok {
				for _, v := range valArr {
					if valMap, ok := v.(map[string]any); ok {
						item := amomodels.FieldValueElement{}
						if val, ok := valMap["value"]; ok {
							item.Value = val
						}
						if enum, ok := valMap["enum_code"].(string); ok {
							item.EnumCode = enum
						}
						cfValue.Values = append(cfValue.Values, item)
					}
				}
			}
			if len(cfValue.Values) > 0 {
				result = append(result, cfValue)
			}
		}
	}

	return result
}
