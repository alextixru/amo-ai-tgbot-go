package models

// ProductsInput входные параметры для инструмента products
type ProductsInput struct {
	// Action действие: search, get, create, update, delete, get_by_entity, link, unlink, update_quantity
	Action string `json:"action" jsonschema_description:"Действие: search, get, create, update, delete, get_by_entity, link, unlink, update_quantity"`

	// ProductID ID товара (для get, update)
	ProductID int `json:"product_id,omitempty" jsonschema_description:"ID товара (для get, update)"`

	// Filter параметры поиска
	Filter *ProductFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска (для search)"`

	// Data данные для создания/обновления (одиночный элемент)
	Data *ProductData `json:"data,omitempty" jsonschema_description:"Данные товара (для create, update, одиночный)"`

	// Items массив данных товаров (для batch create/update)
	Items []ProductData `json:"items,omitempty" jsonschema_description:"Массив товаров (для batch create/update)"`

	// IDs массив ID для удаления или фильтрации
	IDs []int `json:"ids,omitempty" jsonschema_description:"Массив ID товаров (для delete или фильтрации в search)"`

	// Entity сущность для привязки (для get_by_entity, link, unlink)
	Entity *EntityReference `json:"entity,omitempty" jsonschema_description:"Сущность для работы со связями (leads, contacts, companies)"`

	// Product данные для привязки
	Product *ProductLinkData `json:"product,omitempty" jsonschema_description:"Данные привязки товара к сущности"`
}

// EntityReference ссылка на сущность
type EntityReference struct {
	Type string `json:"type" jsonschema_description:"Тип сущности: leads, contacts, companies"`
	ID   int    `json:"id" jsonschema_description:"ID сущности"`
}

// ProductLinkData данные привязки товара
type ProductLinkData struct {
	ID       int `json:"id" jsonschema_description:"ID товара"`
	Quantity int `json:"quantity,omitempty" jsonschema_description:"Количество"`
	PriceID  int `json:"price_id,omitempty" jsonschema_description:"ID цены"`
}

// ProductFilter фильтры поиска товаров
type ProductFilter struct {
	Query string   `json:"query,omitempty" jsonschema_description:"Поисковый запрос"`
	Limit int      `json:"limit,omitempty" jsonschema_description:"Лимит результатов (по умолчанию 50, макс 250)"`
	Page  int      `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	IDs   []int    `json:"ids,omitempty" jsonschema_description:"Фильтр по массиву ID товаров"`
	With  []string `json:"with,omitempty" jsonschema_description:"Дополнительные данные для get: invoice_link, supplier_field_values"`
}

// ProductData данные товара
type ProductData struct {
	Name               string         `json:"name,omitempty" jsonschema_description:"Название товара"`
	CustomFieldsValues map[string]any `json:"custom_fields_values,omitempty" jsonschema_description:"Значения кастомных полей (SKU, цена, описание и др.) в формате SDK"`
}
