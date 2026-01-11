package models

// AdminUsersInput входные параметры для инструмента admin_users
type AdminUsersInput struct {
	// Layer слой: users | roles
	Layer string `json:"layer" jsonschema_description:"Слой: users (пользователи), roles (роли)"`

	// Action действие
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update (только roles), delete (только roles)"`

	// ID идентификатор пользователя или роли
	ID int `json:"id,omitempty" jsonschema_description:"ID пользователя или роли"`

	// Filter фильтры для list
	Filter *AdminUsersFilter `json:"filter,omitempty" jsonschema_description:"Фильтры для поиска"`

	// Data данные для create/update
	Data map[string]any `json:"data,omitempty" jsonschema_description:"Данные для создания/обновления"`
}

// AdminUsersFilter фильтры для admin_users
type AdminUsersFilter struct {
	Limit int      `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Page  int      `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	With  []string `json:"with,omitempty" jsonschema_description:"Связанные данные. Для users: role, uuid, group, amojo_id, user_rank, phone_number. Для roles: users."`
}
