package models

// AdminSchemaInput входные параметры для инструмента admin_schema
type AdminSchemaInput struct {
	// Layer слой схемы: custom_fields | field_groups | loss_reasons | sources
	Layer string `json:"layer" jsonschema_description:"Слой схемы: custom_fields, field_groups, loss_reasons, sources"`

	// Action действие: search | get | create | update | delete
	Action string `json:"action" jsonschema_description:"Действие: search, get, create, update, delete"`

	// EntityType тип сущности (для custom_fields и field_groups): leads | contacts | companies | customers
	EntityType string `json:"entity_type,omitempty" jsonschema_description:"Тип сущности: leads, contacts, companies, customers (для custom_fields и field_groups)"`

	// ID идентификатор элемента (для get, update, delete)
	ID int `json:"id,omitempty" jsonschema_description:"ID элемента (для custom_fields, loss_reasons, sources)"`

	// GroupID идентификатор группы полей (string в API)
	GroupID string `json:"group_id,omitempty" jsonschema_description:"ID группы полей (для field_groups)"`

	// Filter фильтры для search
	Filter *SchemaFilter `json:"filter,omitempty" jsonschema_description:"Фильтры для поиска"`

	// Data данные для create/update
	Data map[string]any `json:"data,omitempty" jsonschema_description:"Данные для создания/обновления"`
}

// SchemaFilter фильтры для поиска в admin_schema
type SchemaFilter struct {
	Limit int `json:"limit,omitempty" jsonschema_description:"Лимит результатов (по умолчанию 50)"`
	Page  int `json:"page,omitempty" jsonschema_description:"Номер страницы (по умолчанию 1)"`
}
