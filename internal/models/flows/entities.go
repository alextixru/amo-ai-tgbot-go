package flows

// EntitiesFlowInput упрощённый вход для entities flow
type EntitiesFlowInput struct {
	BaseFlowInput

	// EntityType тип сущности
	EntityType string `json:"entity_type" jsonschema_description:"Тип: leads, contacts, companies"`

	// Action действие (для mode=direct)
	Action string `json:"action,omitempty" jsonschema_description:"Действие: get, search, create, update"`

	// ID идентификатор (для mode=direct)
	ID int `json:"id,omitempty" jsonschema_description:"ID сущности"`

	// Data упрощённые данные (для mode=direct)
	Data map[string]any `json:"data,omitempty" jsonschema_description:"Данные сущности"`

	// Query поисковый запрос (для mode=direct, action=search)
	Query string `json:"query,omitempty" jsonschema_description:"Поисковый запрос"`
}
