package tools

// CatalogsInput входные параметры для инструмента catalogs
type CatalogsInput struct {
	// Action действие: list, get, create, update, delete, list_elements, get_element, create_element, update_element, delete_element, link_element, unlink_element
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, delete (каталоги), list_elements, get_element, create_element, update_element, delete_element, link_element, unlink_element (элементы)"`

	// CatalogName название каталога (для get, update, delete, list_elements, get_element, create_element, update_element, delete_element, link_element, unlink_element)
	CatalogName string `json:"catalog_name,omitempty" jsonschema_description:"Название каталога (например: 'Товары', 'Услуги', 'Счета')"`

	// ElementID ID элемента (для get_element, update_element, delete_element, link_element, unlink_element)
	ElementID int `json:"element_id,omitempty" jsonschema_description:"ID элемента каталога"`

	// With дополнительные данные: invoice_link, supplier_field_values
	With []string `json:"with,omitempty" jsonschema_description:"Дополнительные данные для get_element: invoice_link, supplier_field_values"`

	// Filter параметры поиска
	Filter *CatalogFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска"`

	// Data данные каталога (для create, update)
	Data *CatalogData `json:"data,omitempty" jsonschema_description:"Данные каталога"`

	// ElementData данные элемента (для create_element, update_element)
	ElementData *CatalogElementData `json:"element_data,omitempty" jsonschema_description:"Данные элемента каталога"`

	// LinkData данные для связи элемента (для link_element, unlink_element)
	LinkData *ElementLinkData `json:"link_data,omitempty" jsonschema_description:"Данные связи элемента с сущностью"`
}

// CatalogFilter фильтры поиска элементов и каталогов
type CatalogFilter struct {
	Page  int    `json:"page,omitempty" jsonschema_description:"Номер страницы (начиная с 1)"`
	Limit int    `json:"limit,omitempty" jsonschema_description:"Лимит результатов (по умолчанию 50, максимум 250)"`
	Query string `json:"query,omitempty" jsonschema_description:"Поисковый запрос по названию элемента каталога"`
	IDs   []int  `json:"ids,omitempty" jsonschema_description:"Фильтр по массиву ID элементов"`
	// Type фильтр по типу каталога для action=list: regular, invoices, products
	Type string `json:"type,omitempty" jsonschema_description:"Тип каталога для фильтрации списка: regular, invoices, products"`
}

// ElementLinkData данные для связи элемента каталога с сущностью
type ElementLinkData struct {
	EntityType string         `json:"entity_type" jsonschema_description:"Тип сущности: leads, contacts, companies, customers"`
	EntityID   int            `json:"entity_id" jsonschema_description:"ID сущности"`
	Metadata   map[string]any `json:"metadata,omitempty" jsonschema_description:"Метаданные связи: quantity (float64), price_id (int)"`
}

// CatalogData данные каталога
type CatalogData struct {
	Name            string `json:"name" jsonschema_description:"Название каталога"`
	Type            string `json:"type,omitempty" jsonschema_description:"Тип: regular, invoices, products"`
	CanAddElements  bool   `json:"can_add_elements,omitempty" jsonschema_description:"Можно ли добавлять элементы"`
	CanShowInCards  bool   `json:"can_show_in_cards,omitempty" jsonschema_description:"Показывать в карточках"`
	CanLinkMultiple bool   `json:"can_link_multiple,omitempty" jsonschema_description:"Разрешить множественную привязку"`
}

// CatalogFieldValue значение кастомного поля элемента каталога
type CatalogFieldValue struct {
	FieldCode string `json:"field_code" jsonschema_description:"Код кастомного поля"`
	Value     string `json:"value" jsonschema_description:"Значение поля"`
}

// CatalogElementData данные элемента каталога
type CatalogElementData struct {
	Name               string              `json:"name" jsonschema_description:"Название элемента"`
	CustomFieldsValues []CatalogFieldValue `json:"custom_fields_values,omitempty" jsonschema_description:"Значения кастомных полей: [{field_code, value}]"`
}
