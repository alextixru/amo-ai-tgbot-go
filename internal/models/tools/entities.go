package tools

import (
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

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

// StatusPair пара воронка+статус для фильтрации сделок
type StatusPair struct {
	PipelineName string `json:"pipeline_name" jsonschema_description:"Название воронки"`
	StatusName   string `json:"status_name" jsonschema_description:"Название статуса"`
}

// EntitiesFilter фильтры поиска сущностей
type EntitiesFilter struct {
	// Общие фильтры
	Query string   `json:"query,omitempty" jsonschema_description:"Поисковый запрос"`
	Limit int      `json:"limit,omitempty" jsonschema_description:"Лимит результатов (макс 250, по умолчанию 50)"`
	Page  int      `json:"page,omitempty" jsonschema_description:"Номер страницы (по умолчанию 1)"`
	IDs   []int    `json:"ids,omitempty" jsonschema_description:"Фильтр по ID сущностей"`
	Names []string `json:"names,omitempty" jsonschema_description:"Фильтр по названию (только contacts, companies)"`

	// Ответственные — по именам
	ResponsibleUserNames []string `json:"responsible_user_names,omitempty" jsonschema_description:"Имена ответственных пользователей"`
	CreatedByNames       []string `json:"created_by_names,omitempty" jsonschema_description:"Имена создателей"`
	UpdatedByNames       []string `json:"updated_by_names,omitempty" jsonschema_description:"Имена обновивших"`

	// Диапазоны дат в ISO-8601 (например "2024-01-15T10:00:00Z" или "2024-01-15")
	CreatedAtFrom string `json:"created_at_from,omitempty" jsonschema_description:"Дата создания от (ISO-8601, например 2024-01-15T10:00:00Z)"`
	CreatedAtTo   string `json:"created_at_to,omitempty" jsonschema_description:"Дата создания до (ISO-8601)"`
	UpdatedAtFrom string `json:"updated_at_from,omitempty" jsonschema_description:"Дата обновления от (ISO-8601)"`
	UpdatedAtTo   string `json:"updated_at_to,omitempty" jsonschema_description:"Дата обновления до (ISO-8601)"`
	ClosedAtFrom  string `json:"closed_at_from,omitempty" jsonschema_description:"Дата закрытия от, только leads (ISO-8601)"`
	ClosedAtTo    string `json:"closed_at_to,omitempty" jsonschema_description:"Дата закрытия до, только leads (ISO-8601)"`

	// Фильтры только для leads — по именам
	PipelineNames []string     `json:"pipeline_names,omitempty" jsonschema_description:"Названия воронок (только leads)"`
	Statuses      []StatusPair `json:"statuses,omitempty" jsonschema_description:"Фильтр по статусам — пары {pipeline_name, status_name} (только leads)"`

	PriceFrom int `json:"price_from,omitempty" jsonschema_description:"Бюджет от (только leads)"`
	PriceTo   int `json:"price_to,omitempty" jsonschema_description:"Бюджет до (только leads)"`

	// CustomFieldsValues для поиска по кастомным полям
	CustomFieldsValues []CustomFieldFilter `json:"custom_fields_values,omitempty" jsonschema_description:"Фильтр по кастомным полям"`

	// With параметры
	With []string `json:"with,omitempty" jsonschema_description:"Включить связанные данные: leads,contacts,companies,catalog_elements,loss_reason,source"`
}

// CustomFieldFilter фильтр по кастомному полю
type CustomFieldFilter struct {
	FieldCode string `json:"field_code" jsonschema_description:"Символьный код кастомного поля (например PHONE, EMAIL, UTM_SOURCE)"`
	Values    []string `json:"values" jsonschema_description:"Значения для фильтрации"`
}

// EntityData данные сущности для create/update
type EntityData struct {
	ID    int    `json:"id,omitempty" jsonschema_description:"ID сущности (для batch update)"`
	Name  string `json:"name,omitempty" jsonschema_description:"Название"`
	Price int    `json:"price,omitempty" jsonschema_description:"Бюджет (только для leads)"`

	// Имена вместо числовых ID
	StatusName          string `json:"status_name,omitempty" jsonschema_description:"Название статуса (только для leads, например 'Новая заявка')"`
	PipelineName        string `json:"pipeline_name,omitempty" jsonschema_description:"Название воронки (только для leads, например 'Основная воронка')"`
	ResponsibleUserName string `json:"responsible_user_name,omitempty" jsonschema_description:"Имя ответственного пользователя (например 'Иван Петров')"`

	// Дополнительные поля для leads
	LossReasonName string `json:"loss_reason_name,omitempty" jsonschema_description:"Название причины отказа (для проигранных сделок)"`
	SourceName     string `json:"source_name,omitempty" jsonschema_description:"Название источника сделки"`

	// Имя и фамилия для contacts
	FirstName string `json:"first_name,omitempty" jsonschema_description:"Имя контакта (только contacts)"`
	LastName  string `json:"last_name,omitempty" jsonschema_description:"Фамилия контакта (только contacts)"`

	CustomFieldsValues map[string]any `json:"custom_fields_values,omitempty" jsonschema_description:"Значения кастомных полей. Ключ — code поля (например PHONE), значение — строка или массив {value, enum_code}"`
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
