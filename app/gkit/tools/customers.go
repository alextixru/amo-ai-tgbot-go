package tools

import (
	"encoding/json"
	"fmt"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// customersRequiredFields определяет обязательные поля для каждой комбинации layer+action.
// Пустой список — действие не требует доп. полей (выполняем сразу, schema mode не нужен).
// Если ни одно из перечисленных полей не присутствует в input — возвращаем схему.
var customersRequiredFields = map[string]map[string][]string{
	"customers": {
		"list":   {},
		"get":    {"id"},
		"create": {"data", "batch"},
		"update": {"id"},
		"delete": {"id"},
		"link":   {"customer_id", "link_data"},
	},
	"bonus_points": {
		"get":           {"customer_id"},
		"earn_points":   {"customer_id"},
		"redeem_points": {"customer_id"},
	},
	"statuses": {
		"list":   {},
		"get":    {"id"},
		"create": {"data"},
		"update": {"id"},
		"delete": {"id"},
	},
	"transactions": {
		"list":   {"customer_id"},
		"create": {"customer_id", "transaction_data"},
		"delete": {"customer_id", "id"},
	},
	"segments": {
		"list":   {},
		"get":    {"id"},
		"create": {"data"},
		"delete": {"id"},
	},
}

// customersIsSchemaMode определяет нужно ли вернуть схему или выполнить действие.
// Schema mode: layer или action неизвестны, либо ни одно обязательное поле не присутствует.
func customersIsSchemaMode(layer, action string, input map[string]any) bool {
	layerReqs, ok := customersRequiredFields[layer]
	if !ok {
		return true
	}
	fields, ok := layerReqs[action]
	if !ok {
		return true
	}
	if len(fields) == 0 {
		return false // нет обязательных полей — execute mode
	}
	// Execute mode если хотя бы одно обязательное поле ненулевое
	for _, f := range fields {
		v, exists := input[f]
		if !exists || v == nil {
			continue
		}
		switch tv := v.(type) {
		case string:
			if tv != "" {
				return false
			}
		case float64:
			if tv != 0 {
				return false
			}
		case bool:
			return false
		case map[string]any:
			if len(tv) > 0 {
				return false
			}
		case []any:
			if len(tv) > 0 {
				return false
			}
		default:
			return false
		}
	}
	return true
}

// customersSchemaResponse формирует ответ со схемой для данного layer+action.
// Включает описание полей, примеры и available_values из сервиса.
func customersSchemaResponse(layer, action string, availableValues map[string]any) map[string]any {
	type e = map[string]any

	schemas := map[string]map[string]e{
		"customers": {
			"list": {
				"description": "Список покупателей с фильтрацией",
				"required_fields": e{},
				"optional_fields": e{
					"filter": e{
						"page":                   "int — номер страницы",
						"limit":                  "int — лимит результатов",
						"query":                  "string — поисковый запрос",
						"responsible_user_names": "[]string — имена ответственных",
						"ids":                    "[]int — ID покупателей",
						"status_names":           "[]string — названия статусов покупателей",
						"names":                  "[]string — имена покупателей",
					},
					"with": "[]string — доп. данные: catalog_elements, contacts, companies, segments",
				},
				"example": e{
					"layer": "customers", "action": "list",
					"filter": e{"limit": 20, "status_names": []string{"VIP"}},
				},
			},
			"get": {
				"description": "Получить покупателя по ID",
				"required_fields": e{
					"id": "int — ID покупателя",
				},
				"optional_fields": e{
					"with": "[]string — доп. данные: catalog_elements, contacts, companies, segments",
				},
				"example": e{
					"layer": "customers", "action": "get", "id": 12345,
					"with": []string{"contacts", "segments"},
				},
			},
			"create": {
				"description": "Создать покупателя. Для батч-создания используй batch вместо data.",
				"required_fields": e{
					"data": e{
						"name":                  "string (required) — имя покупателя",
						"responsible_user_name":  "string — имя ответственного (из available_values.users)",
						"next_date":              "string — дата след. покупки (ISO 8601, например 2024-06-01T00:00:00Z)",
						"next_price":             "int — ожидаемая сумма след. покупки",
						"status_name":            "string — статус (из available_values.customer_statuses)",
						"periodicity":            "int — периодичность покупок (дни)",
						"tags_to_add":            "[]string — теги для добавления",
						"custom_fields_values":   "[]object{field_code, value, enum_code} — кастомные поля",
					},
				},
				"optional_fields": e{
					"batch": "[]CustomerData — массив покупателей для батч-создания вместо data",
				},
				"example": e{
					"layer": "customers", "action": "create",
					"data": e{
						"name":                  "Иван Петров",
						"responsible_user_name":  "Мария Сидорова",
						"status_name":            "Новый",
						"next_price":             5000,
					},
				},
			},
			"update": {
				"description": "Обновить покупателя. Для батч-обновления используй batch.",
				"required_fields": e{
					"id": "int — ID покупателя",
				},
				"optional_fields": e{
					"data": e{
						"name":                  "string",
						"responsible_user_name":  "string — из available_values.users",
						"next_date":              "string (ISO 8601)",
						"next_price":             "int",
						"status_name":            "string — из available_values.customer_statuses",
						"periodicity":            "int",
						"tags_to_add":            "[]string",
						"tags_to_delete":         "[]string",
						"custom_fields_values":   "[]object{field_code, value, enum_code}",
					},
					"batch": "[]CustomerData — для батч-обновления",
				},
				"example": e{
					"layer": "customers", "action": "update",
					"id": 12345,
					"data": e{"next_price": 7000, "status_name": "VIP"},
				},
			},
			"delete": {
				"description": "Удалить покупателя",
				"required_fields": e{
					"id": "int — ID покупателя",
				},
				"example": e{
					"layer": "customers", "action": "delete", "id": 12345,
				},
			},
			"link": {
				"description": "Привязать покупателя к контакту или компании",
				"required_fields": e{
					"customer_id": "int — ID покупателя",
					"link_data": e{
						"entity_type": "string — тип сущности: contacts, companies",
						"entity_id":   "int — ID сущности",
					},
				},
				"example": e{
					"layer": "customers", "action": "link",
					"customer_id": 12345,
					"link_data":   e{"entity_type": "contacts", "entity_id": 67890},
				},
			},
		},
		"bonus_points": {
			"get": {
				"description": "Получить баланс бонусных баллов покупателя",
				"required_fields": e{
					"customer_id": "int — ID покупателя",
				},
				"example": e{
					"layer": "bonus_points", "action": "get", "customer_id": 12345,
				},
			},
			"earn_points": {
				"description": "Начислить бонусные баллы покупателю",
				"required_fields": e{
					"customer_id": "int — ID покупателя",
					"points":      "int — количество баллов для начисления",
				},
				"example": e{
					"layer": "bonus_points", "action": "earn_points",
					"customer_id": 12345, "points": 100,
				},
			},
			"redeem_points": {
				"description": "Списать бонусные баллы покупателя",
				"required_fields": e{
					"customer_id": "int — ID покупателя",
					"points":      "int — количество баллов для списания",
				},
				"example": e{
					"layer": "bonus_points", "action": "redeem_points",
					"customer_id": 12345, "points": 50,
				},
			},
		},
		"statuses": {
			"list": {
				"description": "Список статусов покупателей",
				"required_fields": e{},
				"optional_fields": e{
					"filter": e{
						"page":  "int — номер страницы",
						"limit": "int — лимит результатов",
					},
				},
				"example": e{"layer": "statuses", "action": "list"},
			},
			"get": {
				"description": "Получить статус покупателя по ID",
				"required_fields": e{
					"id": "int — ID статуса",
				},
				"example": e{"layer": "statuses", "action": "get", "id": 1},
			},
			"create": {
				"description": "Создать статус покупателя",
				"required_fields": e{
					"data": e{
						"name": "string (required) — название статуса",
					},
				},
				"example": e{
					"layer": "statuses", "action": "create",
					"data": e{"name": "VIP"},
				},
			},
			"update": {
				"description": "Обновить статус покупателя",
				"required_fields": e{
					"id": "int — ID статуса",
					"data": e{
						"name": "string — новое название",
					},
				},
				"example": e{
					"layer": "statuses", "action": "update",
					"id": 1, "data": e{"name": "Супер VIP"},
				},
			},
			"delete": {
				"description": "Удалить статус покупателя",
				"required_fields": e{
					"id": "int — ID статуса",
				},
				"example": e{"layer": "statuses", "action": "delete", "id": 1},
			},
		},
		"transactions": {
			"list": {
				"description": "Список транзакций покупателя",
				"required_fields": e{
					"customer_id": "int — ID покупателя",
				},
				"optional_fields": e{
					"filter": e{
						"page":  "int — номер страницы",
						"limit": "int — лимит результатов",
					},
				},
				"example": e{
					"layer": "transactions", "action": "list", "customer_id": 12345,
				},
			},
			"create": {
				"description": "Создать транзакцию (фиксация покупки) для покупателя",
				"required_fields": e{
					"customer_id": "int — ID покупателя",
					"transaction_data": e{
						"price":        "int (required) — сумма транзакции",
						"comment":      "string — комментарий",
						"accrue_bonus": "bool — начислить бонусные баллы",
					},
				},
				"example": e{
					"layer": "transactions", "action": "create",
					"customer_id": 12345,
					"transaction_data": e{
						"price": 3500, "comment": "Оплата за март", "accrue_bonus": true,
					},
				},
			},
			"delete": {
				"description": "Удалить транзакцию покупателя",
				"required_fields": e{
					"customer_id": "int — ID покупателя",
					"id":          "int — ID транзакции",
				},
				"example": e{
					"layer": "transactions", "action": "delete",
					"customer_id": 12345, "id": 999,
				},
			},
		},
		"segments": {
			"list": {
				"description": "Список сегментов покупателей",
				"required_fields": e{},
				"optional_fields": e{
					"filter": e{
						"page":  "int — номер страницы",
						"limit": "int — лимит результатов",
					},
				},
				"example": e{"layer": "segments", "action": "list"},
			},
			"get": {
				"description": "Получить сегмент по ID",
				"required_fields": e{
					"id": "int — ID сегмента",
				},
				"example": e{"layer": "segments", "action": "get", "id": 5},
			},
			"create": {
				"description": "Создать сегмент покупателей",
				"required_fields": e{
					"data": e{
						"name": "string (required) — название сегмента",
					},
				},
				"example": e{
					"layer": "segments", "action": "create",
					"data": e{"name": "Постоянные клиенты"},
				},
			},
			"delete": {
				"description": "Удалить сегмент покупателей",
				"required_fields": e{
					"id": "int — ID сегмента",
				},
				"example": e{"layer": "segments", "action": "delete", "id": 5},
			},
		},
	}

	availableLayers := map[string]any{
		"customers":    []string{"list", "get", "create", "update", "delete", "link"},
		"bonus_points": []string{"get", "earn_points", "redeem_points"},
		"statuses":     []string{"list", "get", "create", "update", "delete"},
		"transactions": []string{"list", "create", "delete"},
		"segments":     []string{"list", "get", "create", "delete"},
	}

	layerSchemas, ok := schemas[layer]
	if !ok {
		resp := map[string]any{
			"schema":           true,
			"tool":             "customers",
			"error":            fmt.Sprintf("unknown layer: %q", layer),
			"available_layers": availableLayers,
		}
		if len(availableValues) > 0 {
			resp["available_values"] = availableValues
		}
		return resp
	}

	actionSchema, ok := layerSchemas[action]
	if !ok {
		available := make([]string, 0, len(layerSchemas))
		for a := range layerSchemas {
			available = append(available, a)
		}
		resp := map[string]any{
			"schema":            true,
			"tool":              "customers",
			"layer":             layer,
			"error":             fmt.Sprintf("unknown action: %q for layer %q", action, layer),
			"available_actions": available,
		}
		if len(availableValues) > 0 {
			resp["available_values"] = availableValues
		}
		return resp
	}

	resp := map[string]any{
		"schema": true,
		"tool":   "customers",
		"layer":  layer,
		"action": action,
	}
	for k, v := range actionSchema {
		resp[k] = v
	}
	if len(availableValues) > 0 {
		resp["available_values"] = availableValues
	}
	return resp
}

// RegisterCustomersTool регистрирует customers tool по Shadow Tool паттерну.
//
// LLM видит минимальное описание (layer + action).
// Вызов только с layer+action → handler возвращает полную схему полей + available_values.
// Вызов с обязательными полями → выполняет действие через CRM-сервис.
func (r *Registry) RegisterCustomersTool() {
	r.addTool(genkit.DefineTool[map[string]any, any](
		r.g,
		"customers",
		"Управление покупателями (Retention) в amoCRM: покупатели, бонусные баллы, статусы, транзакции, сегменты. "+
			"Layers: customers (CRUD + link), bonus_points (get/earn_points/redeem_points), "+
			"statuses (CRUD), transactions (list/create/delete), segments (CRUD). "+
			"Вызови с layer+action чтобы получить полную схему параметров и доступные значения.",
		func(ctx *ai.ToolContext, rawInput map[string]any) (any, error) {
			layer, _ := rawInput["layer"].(string)
			action, _ := rawInput["action"].(string)

			availableValues := map[string]any{
				"users":             r.customersService.UserNames(),
				"customer_statuses": r.customersService.StatusNames(),
			}

			if layer == "" || action == "" {
				return map[string]any{
					"schema": true,
					"tool":   "customers",
					"error":  "layer и action обязательны",
					"available_layers": map[string]any{
						"customers":    []string{"list", "get", "create", "update", "delete", "link"},
						"bonus_points": []string{"get", "earn_points", "redeem_points"},
						"statuses":     []string{"list", "get", "create", "update", "delete"},
						"transactions": []string{"list", "create", "delete"},
						"segments":     []string{"list", "get", "create", "delete"},
					},
					"hint":             "Укажи layer + action чтобы получить схему параметров",
					"available_values": availableValues,
				}, nil
			}

			if customersIsSchemaMode(layer, action, rawInput) {
				return customersSchemaResponse(layer, action, availableValues), nil
			}

			// Execute mode: десериализуем map в CustomersInput через JSON roundtrip
			rawBytes, err := json.Marshal(rawInput)
			if err != nil {
				return nil, fmt.Errorf("customers: marshal input: %w", err)
			}
			var fullInput gkitmodels.CustomersInput
			if err := json.Unmarshal(rawBytes, &fullInput); err != nil {
				return nil, fmt.Errorf("customers: parse input: %w", err)
			}

			return r.executeCustomers(ctx, &fullInput)
		},
	))
}

// executeCustomers выполняет действие customers tool.
// Вызывается только в Execute mode (все обязательные поля присутствуют).
func (r *Registry) executeCustomers(ctx *ai.ToolContext, input *gkitmodels.CustomersInput) (any, error) {
	switch input.Layer {
	case "customers":
		switch input.Action {
		case "list":
			return r.customersService.ListCustomers(ctx, input.Filter, input.With)
		case "get":
			if input.ID == 0 {
				return nil, fmt.Errorf("id is required")
			}
			return r.customersService.GetCustomer(ctx, input.ID, input.With)
		case "create":
			if len(input.Batch) > 0 {
				return r.customersService.CreateCustomers(ctx, input.Batch)
			}
			if input.Data == nil {
				return nil, fmt.Errorf("data or batch is required")
			}
			return r.customersService.CreateCustomers(ctx, []*gkitmodels.CustomerData{input.Data})
		case "update":
			if input.ID == 0 {
				return nil, fmt.Errorf("id is required")
			}
			if len(input.Batch) > 0 {
				return r.customersService.UpdateCustomers(ctx, input.ID, input.Batch)
			}
			if input.Data == nil {
				return nil, fmt.Errorf("data or batch is required")
			}
			return r.customersService.UpdateCustomers(ctx, input.ID, []*gkitmodels.CustomerData{input.Data})
		case "delete":
			if input.ID == 0 {
				return nil, fmt.Errorf("id is required")
			}
			return nil, r.customersService.DeleteCustomer(ctx, input.ID)
		case "link":
			if input.CustomerID == 0 || input.LinkData == nil {
				return nil, fmt.Errorf("customer_id and link_data are required")
			}
			return nil, r.customersService.LinkCustomer(ctx, input.CustomerID, input.LinkData.EntityType, input.LinkData.EntityID)
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
			if input.Data == nil || input.Data.Name == "" {
				return nil, fmt.Errorf("data.name is required")
			}
			return r.customersService.CreateCustomerStatuses(ctx, []string{input.Data.Name})
		case "update":
			if input.ID == 0 || input.Data == nil {
				return nil, fmt.Errorf("id and data are required")
			}
			return r.customersService.UpdateCustomerStatus(ctx, input.ID, input.Data.Name)
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
			page, limit := 0, 0
			if input.Filter != nil {
				page, limit = input.Filter.Page, input.Filter.Limit
			}
			return r.customersService.ListTransactions(ctx, input.CustomerID, page, limit)
		case "create":
			if input.TransactionData == nil {
				return nil, fmt.Errorf("transaction_data is required")
			}
			return r.customersService.CreateTransactions(ctx, input.CustomerID, input.TransactionData.Price, input.TransactionData.Comment, input.TransactionData.AccrueBonus)
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
			page, limit := 0, 0
			if input.Filter != nil {
				page, limit = input.Filter.Page, input.Filter.Limit
			}
			return r.customersService.ListSegments(ctx, page, limit)
		case "get":
			if input.ID == 0 {
				return nil, fmt.Errorf("id is required")
			}
			return r.customersService.GetSegment(ctx, input.ID)
		case "create":
			if input.Data == nil || input.Data.Name == "" {
				return nil, fmt.Errorf("data.name is required")
			}
			return r.customersService.CreateSegments(ctx, []string{input.Data.Name})
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
}
