package models

import "github.com/alextixru/amocrm-sdk-go/core/models"

// EntitiesInput входные параметры для инструмента entities
type EntitiesInput struct {
	// EntityType тип сущности: leads, contacts, companies
	EntityType string `json:"entity_type" jsonschema_description:"Тип сущности: leads, contacts, companies"`

	// Action действие: search, get, create, update, sync, delete, link, unlink, get_chats, link_chats
	Action string `json:"action" jsonschema_description:"Действие: search, get, create, update, sync, delete, link, unlink, get_chats, link_chats"`

	// ID идентификатор сущности (для get, update, delete, link, unlink)
	ID int `json:"id,omitempty" jsonschema_description:"ID сущности (для get, update, delete, link, unlink)"`

	// Filter параметры поиска (для search)
	Filter *EntitiesFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска (для action=search)"`

	// Data данные для создания/обновления
	Data *EntityData `json:"data,omitempty" jsonschema_description:"Данные сущности (для create, update)"`

	// LinkTo цель для связывания
	LinkTo *LinkTarget `json:"link_to,omitempty" jsonschema_description:"Цель связывания (для link, unlink)"`

	// ChatLinks ссылки на чаты для link_chats
	ChatLinks []models.ChatLink `json:"chat_links,omitempty" jsonschema_description:"Ссылки на чаты (для link_chats)"`
}

// EntitiesFilter фильтры поиска сущностей
type EntitiesFilter struct {
	Query             string `json:"query,omitempty" jsonschema_description:"Поисковый запрос"`
	Limit             int    `json:"limit,omitempty" jsonschema_description:"Лимит результатов (макс 250)"`
	Page              int    `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	ResponsibleUserID []int  `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственных"`
	PipelineID        []int  `json:"pipeline_id,omitempty" jsonschema_description:"ID воронок (только для leads)"`
	StatusID          []int  `json:"status_id,omitempty" jsonschema_description:"ID статусов (только для leads)"`
}

// EntityData данные сущности для create/update
type EntityData struct {
	Name               string         `json:"name,omitempty" jsonschema_description:"Название"`
	Price              int            `json:"price,omitempty" jsonschema_description:"Бюджет (только для leads)"`
	StatusID           int            `json:"status_id,omitempty" jsonschema_description:"ID статуса (только для leads)"`
	PipelineID         int            `json:"pipeline_id,omitempty" jsonschema_description:"ID воронки (только для leads)"`
	ResponsibleUserID  int            `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственного"`
	CustomFieldsValues map[string]any `json:"custom_fields_values,omitempty" jsonschema_description:"Значения кастомных полей"`
}
