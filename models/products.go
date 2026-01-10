package models

// ProductsInput входные параметры для инструмента products
type ProductsInput struct {
	// Action действие: search, get, create, update, delete, get_by_entity, link, unlink, update_quantity
	Action string `json:"action" jsonschema_description:"Действие: search, get, create, update, delete, get_by_entity, link, unlink, update_quantity"`

	// ProductID ID товара (для get, update)
	ProductID int `json:"product_id,omitempty" jsonschema_description:"ID товара (для get, update)"`

	// Filter параметры поиска
	Filter *ProductFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска (для search)"`

	// Data данные для создания/обновления
	Data *ProductData `json:"data,omitempty" jsonschema_description:"Данные товара (для create, update)"`

	// IDs массив ID для удаления
	IDs []int `json:"ids,omitempty" jsonschema_description:"Массив ID товаров (для delete)"`

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
	Query string `json:"query,omitempty" jsonschema_description:"Поисковый запрос"`
	Limit int    `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Page  int    `json:"page,omitempty" jsonschema_description:"Номер страницы"`
}

// ProductData данные товара
type ProductData struct {
	Name string `json:"name" jsonschema_description:"Название товара"`
	SKU  string `json:"sku,omitempty" jsonschema_description:"Артикул"`
}
