package models

// LinkTarget цель для связывания сущностей
// Используется в entities и activities
type LinkTarget struct {
	Type string `json:"type" jsonschema_description:"Тип целевой сущности: leads, contacts, companies"`
	ID   int    `json:"id" jsonschema_description:"ID целевой сущности"`
}

// ParentEntity родительская сущность
// Используется в activities для указания родителя активности
type ParentEntity struct {
	Type string `json:"type" jsonschema_description:"Тип: leads, contacts, companies"`
	ID   int    `json:"id" jsonschema_description:"ID сущности"`
}
