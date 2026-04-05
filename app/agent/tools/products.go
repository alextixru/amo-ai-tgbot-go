package tools

import (
	"encoding/json"
	"fmt"

	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/crm/products"
)

// ProductsTool реализует нативный ADK Tool интерфейс для управления товарами в amoCRM.
type ProductsTool struct {
	service products.Service
}

// NewProductsTool создаёт новый ProductsTool с переданным сервисом товаров.
func NewProductsTool(service products.Service) *ProductsTool {
	return &ProductsTool{service: service}
}

// Name implements tool.Tool.
func (t *ProductsTool) Name() string {
	return "products"
}

// Description implements tool.Tool.
func (t *ProductsTool) Description() string {
	return "Управление товарами в каталоге products amoCRM: поиск, создание, обновление, удаление, " +
		"получение/привязка/отвязка товаров к сделкам и контактам. " +
		"Actions: search, get, create, update, delete, get_by_entity, link, unlink, update_quantity. " +
		"Вызови с action чтобы получить схему параметров и доступные значения."
}

// IsLongRunning implements tool.Tool.
func (t *ProductsTool) IsLongRunning() bool {
	return false
}

// Declaration implements the ADK runnableTool interface — описание функции для LLM.
func (t *ProductsTool) Declaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        t.Name(),
		Description: t.Description(),
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"action": {
					Type:        genai.TypeString,
					Description: "Действие: search, get, create, update, delete, get_by_entity, link, unlink, update_quantity",
				},
				"product_id": {
					Type:        genai.TypeInteger,
					Description: "ID товара (для get, update, unlink)",
				},
				"with": {
					Type:        genai.TypeArray,
					Description: "Дополнительные данные: invoice_link, supplier_field_values",
					Items:       &genai.Schema{Type: genai.TypeString},
				},
				"filter": {
					Type:        genai.TypeObject,
					Description: "Фильтры поиска (для search): query, limit, page, ids",
				},
				"data": {
					Type:        genai.TypeObject,
					Description: "Данные товара для create/update (одиночный): name, fields [{field_code, value}]",
				},
				"items": {
					Type:        genai.TypeArray,
					Description: "Массив товаров для batch create/update",
					Items:       &genai.Schema{Type: genai.TypeObject},
				},
				"ids": {
					Type:        genai.TypeArray,
					Description: "Массив ID товаров (для delete или фильтрации в search)",
					Items:       &genai.Schema{Type: genai.TypeInteger},
				},
				"entity": {
					Type:        genai.TypeObject,
					Description: "Сущность для работы со связями: {type: 'leads'|'contacts'|'companies', id: int}",
				},
				"product": {
					Type:        genai.TypeObject,
					Description: "Данные привязки товара к сущности: {id: int, quantity: int, price_id: int}",
				},
			},
			Required: []string{"action"},
		},
	}
}

// Run implements the ADK runnableTool interface — точка входа при вызове инструмента LLM.
func (t *ProductsTool) Run(ctx tool.Context, args any) (map[string]any, error) {
	raw, ok := args.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("products: expected map[string]any input, got %T", args)
	}
	result, err := t.handleProductsRaw(ctx, raw)
	if err != nil {
		return nil, err
	}
	return toResultMap(result)
}

// productsSchemas полная схема параметров для каждого action инструмента products.
// Возвращается LLM при первом вызове (schema mode) вместо выполнения действия.
var productsSchemas = map[string]map[string]any{
	"search": {
		"description":     "Поиск товаров в каталоге products по запросу, ID или без фильтров.",
		"required_fields": map[string]any{},
		"optional_fields": map[string]any{
			"filter": map[string]any{
				"type":        "object",
				"description": "Фильтры поиска",
				"fields": map[string]any{
					"query": "string — поисковый запрос по названию",
					"limit": "int — лимит (по умолчанию 50, макс 250)",
					"page":  "int — номер страницы",
					"ids":   "[]int — фильтр по массиву ID товаров",
				},
			},
			"with": "[]string — дополнительные данные: invoice_link, supplier_field_values",
		},
		"example": map[string]any{
			"action": "search",
			"filter": map[string]any{"query": "виджет", "limit": 20},
		},
	},
	"get": {
		"description": "Получение товара по ID.",
		"required_fields": map[string]any{
			"product_id": "int — ID товара",
		},
		"optional_fields": map[string]any{
			"with": "[]string — дополнительные данные: invoice_link, supplier_field_values",
		},
		"example": map[string]any{
			"action":     "get",
			"product_id": 12345,
		},
	},
	"create": {
		"description": "Создание одного или нескольких товаров в каталоге products.",
		"required_fields": map[string]any{
			"data OR items": "object/array — данные товара(ов). Используй data для одного, items для нескольких.",
		},
		"optional_fields": map[string]any{},
		"data_fields": map[string]any{
			"name":   "string — название товара (обязательно)",
			"fields": "[]object — кастомные поля: [{field_code: 'SKU', value: 'ABC'}, {field_code: 'PRICE', value: '1000'}]",
		},
		"example": map[string]any{
			"action": "create",
			"data": map[string]any{
				"name":   "Виджет Pro",
				"fields": []map[string]any{{"field_code": "SKU", "value": "WP-001"}, {"field_code": "PRICE", "value": "1500"}},
			},
		},
	},
	"update": {
		"description": "Обновление одного или нескольких товаров.",
		"required_fields": map[string]any{
			"data OR items": "object/array — данные товара(ов) с ID. Используй data для одного (с product_id или data.id), items для batch.",
		},
		"optional_fields": map[string]any{
			"product_id": "int — ID товара (если используется data без id внутри)",
		},
		"data_fields": map[string]any{
			"id":     "int — ID товара (для batch update через items)",
			"name":   "string — новое название",
			"fields": "[]object — кастомные поля: [{field_code, value}]",
		},
		"example": map[string]any{
			"action":     "update",
			"product_id": 12345,
			"data": map[string]any{
				"name":   "Виджет Pro v2",
				"fields": []map[string]any{{"field_code": "PRICE", "value": "2000"}},
			},
		},
	},
	"delete": {
		"description": "Удаление товаров по массиву ID.",
		"required_fields": map[string]any{
			"ids": "[]int — массив ID товаров для удаления",
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action": "delete",
			"ids":    []int{12345, 67890},
		},
	},
	"get_by_entity": {
		"description": "Получение товаров, привязанных к сущности (сделке, контакту или компании).",
		"required_fields": map[string]any{
			"entity": "object — сущность: {type: 'leads'|'contacts'|'companies', id: int}",
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action": "get_by_entity",
			"entity": map[string]any{"type": "leads", "id": 98765},
		},
	},
	"link": {
		"description": "Привязка товара к сущности (сделке, контакту или компании).",
		"required_fields": map[string]any{
			"entity":  "object — сущность: {type: 'leads'|'contacts'|'companies', id: int}",
			"product": "object — товар: {id: int, quantity: int (необяз.), price_id: int (необяз.)}",
		},
		"optional_fields": map[string]any{},
		"notes":           "Если price_id не указан, используется первое ценовое поле каталога.",
		"example": map[string]any{
			"action":  "link",
			"entity":  map[string]any{"type": "leads", "id": 98765},
			"product": map[string]any{"id": 12345, "quantity": 3},
		},
	},
	"unlink": {
		"description": "Отвязка товара от сущности.",
		"required_fields": map[string]any{
			"entity":     "object — сущность: {type: 'leads'|'contacts'|'companies', id: int}",
			"product_id": "int — ID товара для отвязки",
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action":     "unlink",
			"entity":     map[string]any{"type": "leads", "id": 98765},
			"product_id": 12345,
		},
	},
	"update_quantity": {
		"description": "Обновление количества и/или цены товара в привязке к сущности.",
		"required_fields": map[string]any{
			"entity":  "object — сущность: {type: 'leads'|'contacts'|'companies', id: int}",
			"product": "object — товар: {id: int, quantity: int, price_id: int (необяз.)}",
		},
		"optional_fields": map[string]any{},
		"notes":           "В amoCRM реализовано через повторный Link — API обновляет metadata если связь уже существует.",
		"example": map[string]any{
			"action":  "update_quantity",
			"entity":  map[string]any{"type": "leads", "id": 98765},
			"product": map[string]any{"id": 12345, "quantity": 5},
		},
	},
}

// handleProductsRaw — точка входа shadow handler.
// Получает raw map[string]any, определяет режим (schema/execute) и действует соответственно.
func (t *ProductsTool) handleProductsRaw(ctx tool.Context, raw map[string]any) (any, error) {
	action, _ := raw["action"].(string)
	if action == "" {
		return map[string]any{
			"schema":            true,
			"tool":              "products",
			"error":             "action is required",
			"available_actions": productsActionList(),
		}, nil
	}

	if isProductsSchemaMode(raw, action) {
		return t.productsSchemaResponse(ctx, action)
	}

	return t.handleProducts(ctx, raw)
}

// isProductsSchemaMode определяет, нужно ли вернуть схему (true) или выполнить действие (false).
// Возвращает true если отсутствуют обязательные поля для данного action.
func isProductsSchemaMode(raw map[string]any, action string) bool {
	switch action {
	case "search":
		return false // нет обязательных полей — всегда execute
	case "get":
		return !hasIntField(raw, "product_id")
	case "create", "update":
		_, hasData := raw["data"]
		_, hasItems := raw["items"]
		return !hasData && !hasItems
	case "delete":
		return !hasArrayField(raw, "ids")
	case "get_by_entity":
		return !hasObjectField(raw, "entity")
	case "link", "update_quantity":
		return !hasObjectField(raw, "entity") || !hasObjectField(raw, "product")
	case "unlink":
		return !hasObjectField(raw, "entity") || !hasIntField(raw, "product_id")
	default:
		return true // неизвестный action → возвращаем схему с ошибкой
	}
}

// productsSchemaResponse формирует schema response для заданного action.
func (t *ProductsTool) productsSchemaResponse(ctx tool.Context, action string) (any, error) {
	schema, ok := productsSchemas[action]
	if !ok {
		return map[string]any{
			"schema":            true,
			"tool":              "products",
			"error":             fmt.Sprintf("unknown action: %q", action),
			"available_actions": productsActionList(),
		}, nil
	}

	resp := map[string]any{
		"schema":      true,
		"tool":        "products",
		"action":      action,
		"description": schema["description"],
	}
	if rf, ok := schema["required_fields"]; ok {
		resp["required_fields"] = rf
	}
	if of, ok := schema["optional_fields"]; ok {
		resp["optional_fields"] = of
	}
	if df, ok := schema["data_fields"]; ok {
		resp["data_fields"] = df
	}
	if notes, ok := schema["notes"]; ok {
		resp["notes"] = notes
	}
	if ex, ok := schema["example"]; ok {
		resp["example"] = ex
	}

	// available_values: загружаем названия товаров (до 20 штук) для контекста LLM
	resp["available_values"] = t.productsAvailableValues(ctx)

	return resp, nil
}

// productsAvailableValues возвращает справочные данные для schema response.
func (t *ProductsTool) productsAvailableValues(ctx tool.Context) map[string]any {
	vals := map[string]any{}

	result, err := t.service.SearchProducts(ctx, &gkitmodels.ProductFilter{Limit: 20, Page: 1}, nil)
	if err == nil && result != nil {
		names := make([]string, 0, len(result.Items))
		for _, item := range result.Items {
			if item != nil && item.Name != "" {
				names = append(names, item.Name)
			}
		}
		vals["product_names"] = names
		if result.HasMore {
			vals["product_names_note"] = "показаны первые 20; используй action=search с filter.query для уточнения"
		}
	}

	return vals
}

// handleProducts выполняет действие с товарами (Execute mode).
// Конвертирует raw map[string]any в ProductsInput через JSON roundtrip.
func (t *ProductsTool) handleProducts(ctx tool.Context, raw map[string]any) (any, error) {
	b, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("products: marshal input: %w", err)
	}
	var input gkitmodels.ProductsInput
	if err := json.Unmarshal(b, &input); err != nil {
		return nil, fmt.Errorf("products: unmarshal input: %w", err)
	}
	return t.executeProducts(ctx, input)
}

// executeProducts выполняет действие с уже десериализованным ProductsInput.
func (t *ProductsTool) executeProducts(ctx tool.Context, input gkitmodels.ProductsInput) (any, error) {
	switch input.Action {
	case "search":
		return t.service.SearchProducts(ctx, input.Filter, input.With)
	case "get":
		if input.ProductID == 0 {
			return nil, fmt.Errorf("product_id is required for action 'get'")
		}
		return t.service.GetProduct(ctx, input.ProductID, input.With)
	case "create":
		if input.Data == nil && len(input.Items) == 0 {
			return nil, fmt.Errorf("data or items is required for action 'create'")
		}
		var items []gkitmodels.ProductData
		if len(input.Items) > 0 {
			items = input.Items
		} else {
			items = []gkitmodels.ProductData{*input.Data}
		}
		return t.service.CreateProducts(ctx, items)
	case "update":
		if input.Data == nil && len(input.Items) == 0 {
			return nil, fmt.Errorf("data or items (for batch) are required for action 'update'")
		}
		var items []gkitmodels.ProductData
		if len(input.Items) > 0 {
			items = input.Items
		} else {
			data := *input.Data
			if data.ID == 0 {
				data.ID = input.ProductID
			}
			items = []gkitmodels.ProductData{data}
		}
		return t.service.UpdateProducts(ctx, items)
	case "delete":
		if len(input.IDs) == 0 {
			return nil, fmt.Errorf("ids array is required for action 'delete'")
		}
		return t.service.DeleteProducts(ctx, input.IDs)
	case "get_by_entity":
		if input.Entity == nil {
			return nil, fmt.Errorf("entity (type, id) is required for action 'get_by_entity'")
		}
		return t.service.GetProductsByEntity(ctx, input.Entity.Type, input.Entity.ID)
	case "link":
		if input.Entity == nil || input.Product == nil {
			return nil, fmt.Errorf("entity and product are required for action 'link'")
		}
		return t.service.LinkProduct(ctx, input.Entity.Type, input.Entity.ID, input.Product.ID, input.Product.Quantity, input.Product.PriceID)
	case "unlink":
		if input.Entity == nil || input.ProductID == 0 {
			return nil, fmt.Errorf("entity and product_id are required for action 'unlink'")
		}
		return t.service.UnlinkProduct(ctx, input.Entity.Type, input.Entity.ID, input.ProductID)
	case "update_quantity":
		if input.Entity == nil || input.Product == nil {
			return nil, fmt.Errorf("entity and product (id, quantity, price_id) are required for action 'update_quantity'")
		}
		// В amoCRM обновление количества — это повторный Link (v4 обновляет metadata если связь уже есть)
		return t.service.LinkProduct(ctx, input.Entity.Type, input.Entity.ID, input.Product.ID, input.Product.Quantity, input.Product.PriceID)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

// productsActionList возвращает список доступных actions для error messages.
func productsActionList() []string {
	actions := make([]string, 0, len(productsSchemas))
	for a := range productsSchemas {
		actions = append(actions, a)
	}
	return actions
}

// --- helpers ---

func hasIntField(raw map[string]any, key string) bool {
	v, ok := raw[key]
	if !ok {
		return false
	}
	switch n := v.(type) {
	case float64:
		return n != 0
	case int:
		return n != 0
	case int64:
		return n != 0
	}
	return false
}

func hasArrayField(raw map[string]any, key string) bool {
	v, ok := raw[key]
	if !ok {
		return false
	}
	arr, ok := v.([]any)
	return ok && len(arr) > 0
}

func hasObjectField(raw map[string]any, key string) bool {
	v, ok := raw[key]
	if !ok {
		return false
	}
	_, ok = v.(map[string]any)
	return ok
}
