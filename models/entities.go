package models

import "github.com/alextixru/amocrm-sdk-go/core/models"

// EntitiesInput входные параметры для инструмента entities
type EntitiesInput struct {
	// EntityType тип сущности: leads, contacts, companies
	EntityType string `json:"entity_type" jsonschema_description:"Тип сущности: leads, contacts, companies"`

	// Action действие: search, get, create, update, sync, link, unlink, get_chats, link_chats
	Action string `json:"action" jsonschema_description:"Действие: search, get, create, update, sync, delete (только если доступно), link, unlink, get_chats (только contacts), link_chats (только contacts)"`

	// ID идентификатор сущности (для get, update, delete, link, unlink)
	ID int `json:"id,omitempty" jsonschema_description:"ID сущности (для get, update, delete, link, unlink)"`

	// Filter параметры поиска (для search)
	Filter *EntitiesFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска (для action=search)"`

	// Data данные для создания/обновления
	Data *EntityData `json:"data,omitempty" jsonschema_description:"Данные сущности (для create, update, sync)"`

	// DataList данные для batch создания/обновления
	DataList []EntityData `json:"data_list,omitempty" jsonschema_description:"Массив данных для batch create/update (вместо data)"`

	// With параметры для включения связанных данных
	With []string `json:"with,omitempty" jsonschema_description:"Связанные данные для get/search: leads,contacts,companies,catalog_elements,loss_reason,source"`

	// LinkTo цель для связывания
	LinkTo *LinkTarget `json:"link_to,omitempty" jsonschema_description:"Цель связывания (для link, unlink)"`

	// ChatLinks ссылки на чаты для link_chats
	ChatLinks []models.ChatLink `json:"chat_links,omitempty" jsonschema_description:"Ссылки на чаты (для link_chats)"`
}

// EntitiesFilter фильтры поиска сущностей
type EntitiesFilter struct {
	// Общие фильтры
	Query             string   `json:"query,omitempty" jsonschema_description:"Поисковый запрос"`
	Limit             int      `json:"limit,omitempty" jsonschema_description:"Лимит результатов (макс 250, по умолчанию 50)"`
	Page              int      `json:"page,omitempty" jsonschema_description:"Номер страницы (по умолчанию 1)"`
	IDs               []int    `json:"ids,omitempty" jsonschema_description:"Фильтр по ID сущностей"`
	Names             []string `json:"names,omitempty" jsonschema_description:"Фильтр по названию (только contacts, companies)"`
	ResponsibleUserID []int    `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственных"`
	CreatedBy         []int    `json:"created_by,omitempty" jsonschema_description:"ID создателей"`
	UpdatedBy         []int    `json:"updated_by,omitempty" jsonschema_description:"ID обновивших"`

	// Диапазоны дат (Unix timestamp)
	CreatedAtFrom int64 `json:"created_at_from,omitempty" jsonschema_description:"Дата создания от (Unix timestamp)"`
	CreatedAtTo   int64 `json:"created_at_to,omitempty" jsonschema_description:"Дата создания до (Unix timestamp)"`
	UpdatedAtFrom int64 `json:"updated_at_from,omitempty" jsonschema_description:"Дата обновления от (Unix timestamp)"`
	UpdatedAtTo   int64 `json:"updated_at_to,omitempty" jsonschema_description:"Дата обновления до (Unix timestamp)"`
	ClosedAtFrom  int64 `json:"closed_at_from,omitempty" jsonschema_description:"Дата закрытия от, только leads (Unix timestamp)"`
	ClosedAtTo    int64 `json:"closed_at_to,omitempty" jsonschema_description:"Дата закрытия до, только leads (Unix timestamp)"`

	// Фильтры только для leads
	PipelineID []int `json:"pipeline_id,omitempty" jsonschema_description:"ID воронок (только leads)"`
	StatusID   []int `json:"status_id,omitempty" jsonschema_description:"ID статусов (только leads)"`
	PriceFrom  int   `json:"price_from,omitempty" jsonschema_description:"Бюджет от (только leads)"`
	PriceTo    int   `json:"price_to,omitempty" jsonschema_description:"Бюджет до (только leads)"`

	// CustomFieldsValues для поиска по кастомным полям
	CustomFieldsValues []CustomFieldFilter `json:"custom_fields_values,omitempty" jsonschema_description:"Фильтр по кастомным полям"`

	// With параметры
	With []string `json:"with,omitempty" jsonschema_description:"Включить связанные данные: leads,contacts,companies,catalog_elements,loss_reason,source"`
}

// CustomFieldFilter фильтр по кастомному полю
type CustomFieldFilter struct {
	FieldID int   `json:"field_id" jsonschema_description:"ID кастомного поля"`
	Values  []any `json:"values" jsonschema_description:"Значения для фильтрации"`
}

// EntityData данные сущности для create/update
type EntityData struct {
	ID                 int            `json:"id,omitempty" jsonschema_description:"ID сущности (для batch update)"`
	Name               string         `json:"name,omitempty" jsonschema_description:"Название"`
	Price              int            `json:"price,omitempty" jsonschema_description:"Бюджет (только для leads)"`
	StatusID           int            `json:"status_id,omitempty" jsonschema_description:"ID статуса (только для leads)"`
	PipelineID         int            `json:"pipeline_id,omitempty" jsonschema_description:"ID воронки (только для leads)"`
	ResponsibleUserID  int            `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственного"`
	CustomFieldsValues map[string]any `json:"custom_fields_values,omitempty" jsonschema_description:"Значения кастомных полей"`
	Tags               []EntityTag    `json:"tags,omitempty" jsonschema_description:"Теги сущности"`

	// Embedded связанные сущности (для create с привязкой)
	EmbeddedContacts  []int `json:"embedded_contacts,omitempty" jsonschema_description:"ID контактов для привязки (leads, companies)"`
	EmbeddedCompanies []int `json:"embedded_companies,omitempty" jsonschema_description:"ID компаний для привязки (leads, contacts)"`
}

// EntityTag тег сущности
type EntityTag struct {
	ID   int    `json:"id,omitempty" jsonschema_description:"ID тега (для существующего)"`
	Name string `json:"name,omitempty" jsonschema_description:"Название тега (для нового)"`
}
