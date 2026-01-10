package models

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
}

// IntegrationsFilter фильтры для admin_integrations
type IntegrationsFilter struct {
	Limit int `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Page  int `json:"page,omitempty" jsonschema_description:"Номер страницы"`
}
