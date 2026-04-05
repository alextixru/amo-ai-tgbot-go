package tools

// AdminIntegrationsInput входные параметры для инструмента admin_integrations
type AdminIntegrationsInput struct {
	// Layer слой: webhooks | widgets | website_buttons | chat_templates | short_links
	Layer string `json:"layer" jsonschema_description:"Слой: webhooks, widgets, website_buttons, chat_templates, short_links"`

	// Action действие
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, delete, subscribe, unsubscribe, install, uninstall, delete_many, send_review, update_review, add_chat"`

	// ID идентификатор (для get, update, delete)
	ID int `json:"id,omitempty" jsonschema_description:"ID элемента (source_id для website_buttons)"`

	// Code код виджета (для widgets)
	Code string `json:"code,omitempty" jsonschema_description:"Код виджета (для widgets: get, install, uninstall)"`

	// Filter фильтры для list
	Filter *IntegrationsFilter `json:"filter,omitempty" jsonschema_description:"Фильтры для поиска"`

	// With параметры обогащения ответа (не фильтрация)
	// Для website_buttons: 'scripts' — включить скрипты, 'deleted' — включить удалённые кнопки
	// Для webhooks, widgets: не используется
	With []string `json:"with,omitempty" jsonschema_description:"Параметры обогащения: для website_buttons — 'scripts' (скрипты), 'deleted' (удалённые кнопки)"`

	// Destination URL для webhook subscribe/unsubscribe
	Destination string `json:"destination,omitempty" jsonschema_description:"URL вебхука (для webhooks: subscribe, unsubscribe)"`

	// EventTypes список событий для webhook
	// Доступные события: add_lead, update_lead, delete_lead, restore_lead, status_lead,
	// add_contact, update_contact, delete_contact, restore_contact,
	// add_company, update_company, delete_company, restore_company,
	// add_task, update_task, delete_task, complete_task,
	// add_note, update_note, delete_note
	EventTypes []string `json:"event_types,omitempty" jsonschema_description:"События вебхука. Доступные: add_lead, update_lead, delete_lead, restore_lead, status_lead, add_contact, update_contact, delete_contact, restore_contact, add_company, update_company, delete_company, restore_company, add_task, update_task, delete_task, complete_task, add_note, update_note, delete_note"`

	// WebsiteButton данные для создания/обновления кнопки сайта
	WebsiteButton *WebsiteButtonData `json:"website_button,omitempty" jsonschema_description:"Данные кнопки на сайте (для website_buttons: create, update)"`

	// ChatTemplate данные для создания/обновления шаблона чата
	ChatTemplate *ChatTemplateData `json:"chat_template,omitempty" jsonschema_description:"Данные шаблона чата (для chat_templates: create, update)"`

	// ReviewID ID ревью (для chat_templates: update_review)
	ReviewID int `json:"review_id,omitempty" jsonschema_description:"ID ревью шаблона (для chat_templates: update_review — берётся из ответа send_review)"`

	// ReviewStatus новый статус ревью (для chat_templates: update_review)
	ReviewStatus string `json:"review_status,omitempty" jsonschema_description:"Новый статус ревью (для chat_templates: update_review)"`

	// IDs идентификаторы для батч-операций (delete_many)
	IDs []int `json:"ids,omitempty" jsonschema_description:"Массив ID элементов (для chat_templates: delete_many)"`

	// URL ссылка для создания короткой ссылки (для short_links: create)
	URL string `json:"url,omitempty" jsonschema_description:"URL для создания короткой ссылки (для short_links: create)"`

	// URLs массив URL для батч-создания (short_links)
	URLs []string `json:"urls,omitempty" jsonschema_description:"Массив URL (для short_links: create — батч)"`

	// EntityID ID сущности для привязки короткой ссылки (для short_links: create)
	EntityID int `json:"entity_id,omitempty" jsonschema_description:"ID сущности для привязки короткой ссылки (для short_links: create)"`

	// EntityType тип сущности для привязки короткой ссылки (для short_links: create)
	// Доступные значения: leads, contacts, companies, customers
	EntityType string `json:"entity_type,omitempty" jsonschema_description:"Тип сущности: leads, contacts, companies, customers (для short_links: create)"`

	// Settings настройки виджета
	Settings map[string]any `json:"settings,omitempty" jsonschema_description:"Настройки виджета (для widgets: install)"`
}

// IntegrationsFilter фильтры для admin_integrations
type IntegrationsFilter struct {
	Limit        int      `json:"limit,omitempty" jsonschema_description:"Лимит результатов (по умолчанию 50)"`
	Page         int      `json:"page,omitempty" jsonschema_description:"Номер страницы (по умолчанию 1)"`
	Destination  string   `json:"destination,omitempty" jsonschema_description:"URL вебхука для фильтрации (для webhooks: list)"`
	ExternalIDs  []string `json:"external_ids,omitempty" jsonschema_description:"Внешние ID для фильтрации (для chat_templates: list)"`
	TemplateType string   `json:"template_type,omitempty" jsonschema_description:"Тип шаблона: amocrm или waba (для chat_templates: list)"`
}

// WebsiteButtonData типизированные данные для создания/обновления кнопки сайта
type WebsiteButtonData struct {
	// Name название кнопки
	Name string `json:"name,omitempty" jsonschema_description:"Название кнопки"`

	// PipelineID ID воронки (числовой, получить через admin_pipelines: list)
	PipelineID *int `json:"pipeline_id,omitempty" jsonschema_description:"ID воронки (числовой). Получить список воронок через admin_pipelines: list"`

	// TrustedWebsites список доверенных сайтов
	TrustedWebsites []string `json:"trusted_websites,omitempty" jsonschema_description:"Список доверенных сайтов (например: ['example.com'])"`

	// IsDuplicationControlEnabled включить контроль дублей
	IsDuplicationControlEnabled *bool `json:"is_duplication_control_enabled,omitempty" jsonschema_description:"Включить контроль дублей"`
}

// ChatTemplateData типизированные данные для создания/обновления шаблона чата
type ChatTemplateData struct {
	// Name название шаблона
	Name string `json:"name,omitempty" jsonschema_description:"Название шаблона"`

	// Content содержимое шаблона
	Content string `json:"content,omitempty" jsonschema_description:"Текст шаблона"`

	// ExternalID внешний ID шаблона
	ExternalID string `json:"external_id,omitempty" jsonschema_description:"Внешний ID шаблона"`

	// Type тип шаблона: amocrm или waba
	Type string `json:"type,omitempty" jsonschema_description:"Тип шаблона: amocrm или waba"`

	// IsEditable можно ли редактировать
	IsEditable bool `json:"is_editable,omitempty" jsonschema_description:"Можно ли редактировать шаблон"`

	// WabaHeader заголовок WABA-шаблона
	WabaHeader string `json:"waba_header,omitempty" jsonschema_description:"Заголовок WABA-шаблона"`

	// WabaFooter подвал WABA-шаблона
	WabaFooter string `json:"waba_footer,omitempty" jsonschema_description:"Подвал WABA-шаблона"`

	// WabaCategory категория WABA: UTILITY, AUTHENTICATION, MARKETING
	WabaCategory string `json:"waba_category,omitempty" jsonschema_description:"Категория WABA: UTILITY, AUTHENTICATION, MARKETING"`

	// WabaLanguage язык WABA (например: ru, en)
	WabaLanguage string `json:"waba_language,omitempty" jsonschema_description:"Язык WABA (например: ru, en)"`
}
