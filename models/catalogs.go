package models

// CatalogsInput входные параметры для инструмента catalogs
type CatalogsInput struct {
	// Action действие: list, get, create, update, list_elements, get_element, create_element, update_element
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, list_elements, get_element, create_element, update_element"`

	// CatalogID ID каталога (для get, update, list_elements, get_element, create_element, update_element)
	CatalogID int `json:"catalog_id,omitempty" jsonschema_description:"ID каталога"`

	// ElementID ID элемента (для get_element, update_element)
	ElementID int `json:"element_id,omitempty" jsonschema_description:"ID элемента каталога"`

	// Filter параметры поиска
	Filter *CatalogFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска"`

	// Data данные каталога (для create, update)
	Data *CatalogData `json:"data,omitempty" jsonschema_description:"Данные каталога"`

	// ElementData данные элемента (для create_element, update_element)
	ElementData *CatalogElementData `json:"element_data,omitempty" jsonschema_description:"Данные элемента каталога"`
}

// CatalogFilter фильтры поиска
type CatalogFilter struct {
	Page  int    `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Limit int    `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Query string `json:"query,omitempty" jsonschema_description:"Поисковый запрос (только для элементов)"`
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
