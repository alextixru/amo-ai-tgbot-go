package tools

// AdminIntegrationsInput входные параметры для инструмента admin_integrations
type AdminIntegrationsInput struct {
	// Layer слой: webhooks | widgets | website_buttons | chat_templates | short_links
	Layer string `json:"layer" jsonschema_description:"Слой: webhooks, widgets, website_buttons, chat_templates, short_links"`

	// Action действие
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, delete, subscribe, unsubscribe, install, uninstall"`

	// ID идентификатор (для get, update, delete)
	ID int `json:"id,omitempty" jsonschema_description:"ID элемента"`

	// Code код виджета (для widgets)
	Code string `json:"code,omitempty" jsonschema_description:"Код виджета (для widgets: get, install, uninstall)"`

	// Filter фильтры для list
	Filter *IntegrationsFilter `json:"filter,omitempty" jsonschema_description:"Фильтры для поиска"`

	// Data данные для create/update
	Data map[string]any `json:"data,omitempty" jsonschema_description:"Данные для создания/обновления"`

	// IDs идентификаторы для батч-операций (delete_many)
	IDs []int `json:"ids,omitempty" jsonschema_description:"Массив ID элементов (для chat_templates: delete_many)"`

	// URLs массив URL для батч-создания (short_links)
	URLs []string `json:"urls,omitempty" jsonschema_description:"Массив URL (для short_links: create)"`

	// Settings настройки виджета
	Settings map[string]any `json:"settings,omitempty" jsonschema_description:"Настройки виджета (для widgets: install)"`
}

// IntegrationsFilter фильтры для admin_integrations
type IntegrationsFilter struct {
	Limit       int      `json:"limit,omitempty" jsonschema_description:"Лимит результатов (по умолчанию 50)"`
	Page        int      `json:"page,omitempty" jsonschema_description:"Номер страницы (по умолчанию 1)"`
	Destination string   `json:"destination,omitempty" jsonschema_description:"URL вебхука для фильтрации (для webhooks)"`
	ExternalIDs []string `json:"external_ids,omitempty" jsonschema_description:"Внешние ID для фильтрации (для chat_templates)"`
	With        []string `json:"with,omitempty" jsonschema_description:"Связанные данные: scripts (для website_buttons)"`
}
