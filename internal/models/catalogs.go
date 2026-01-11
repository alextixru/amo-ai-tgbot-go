package models

// CatalogsInput входные параметры для инструмента catalogs
type CatalogsInput struct {
	// Action действие: list, get, create, update, list_elements, get_element, create_element, update_element, link_element, unlink_element
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, list_elements, get_element, create_element, update_element, link_element, unlink_element"`

	// CatalogID ID каталога (для get, update, list_elements, get_element, create_element, update_element, link_element, unlink_element)
	CatalogID int `json:"catalog_id,omitempty" jsonschema_description:"ID каталога"`

	// ElementID ID элемента (для get_element, update_element, link_element, unlink_element)
	ElementID int `json:"element_id,omitempty" jsonschema_description:"ID элемента каталога"`

	// Filter параметры поиска
	Filter *CatalogFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска"`

	// Data данные каталога (для create, update)
	Data *CatalogData `json:"data,omitempty" jsonschema_description:"Данные каталога"`

	// ElementData данные элемента (для create_element, update_element)
	ElementData *CatalogElementData `json:"element_data,omitempty" jsonschema_description:"Данные элемента каталога"`

	// LinkData данные для связи элемента (для link_element, unlink_element)
	LinkData *ElementLinkData `json:"link_data,omitempty" jsonschema_description:"Данные связи элемента с сущностью"`
}

// CatalogFilter фильтры поиска
type CatalogFilter struct {
	Page  int      `json:"page,omitempty" jsonschema_description:"Номер страницы (начиная с 1)"`
	Limit int      `json:"limit,omitempty" jsonschema_description:"Лимит результатов (по умолчанию 50, максимум 250)"`
	Query string   `json:"query,omitempty" jsonschema_description:"Поисковый запрос по названию элемента каталога"`
	IDs   []int    `json:"ids,omitempty" jsonschema_description:"Фильтр по массиву ID элементов"`
	With  []string `json:"with,omitempty" jsonschema_description:"Дополнительные данные для get_element: invoice_link, supplier_field_values"`
}

// ElementLinkData данные для связи элемента каталога с сущностью
type ElementLinkData struct {
	EntityType string         `json:"entity_type" jsonschema_description:"Тип сущности: leads, contacts, companies, customers"`
	EntityID   int            `json:"entity_id" jsonschema_description:"ID сущности"`
	Metadata   map[string]any `json:"metadata,omitempty" jsonschema_description:"Метаданные связи (quantity, price_id и др.)"`
}

// CatalogData данные каталога
type CatalogData struct {
	Name            string `json:"name" jsonschema_description:"Название каталога"`
	Type            string `json:"type,omitempty" jsonschema_description:"Тип: regular, invoices, products"`
	CanAddElements  bool   `json:"can_add_elements,omitempty" jsonschema_description:"Можно ли добавлять элементы"`
	CanShowInCards  bool   `json:"can_show_in_cards,omitempty" jsonschema_description:"Показывать в карточках"`
	CanLinkMultiple bool   `json:"can_link_multiple,omitempty" jsonschema_description:"Разрешить множественную привязку"`
}

// CatalogElementData данные элемента каталога
type CatalogElementData struct {
	Name               string         `json:"name" jsonschema_description:"Название элемента"`
	CustomFieldsValues map[string]any `json:"custom_fields_values,omitempty" jsonschema_description:"Значения кастомных полей"`
}
