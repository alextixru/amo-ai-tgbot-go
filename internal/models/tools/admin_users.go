package tools

// AdminUsersInput входные параметры для инструмента admin_users
type AdminUsersInput struct {
	// Layer слой: users | roles
	Layer string `json:"layer" jsonschema_description:"Слой: users (пользователи), roles (роли)"`

	// Action действие
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update (только roles), delete (только roles), add_to_group (только users)"`

	// ID идентификатор пользователя или роли
	ID int `json:"id,omitempty" jsonschema_description:"ID пользователя или роли"`

	// UserID ID пользователя (для add_to_group)
	UserID int `json:"user_id,omitempty" jsonschema_description:"ID пользователя (используется в add_to_group)"`

	// GroupID ID группы (для add_to_group)
	GroupID int `json:"group_id,omitempty" jsonschema_description:"ID группы (используется в add_to_group)"`

	// Filter фильтры для list
	Filter *AdminUsersFilter `json:"filter,omitempty" jsonschema_description:"Фильтры для поиска"`

	// Users типизированные данные для создания пользователей (action=create, layer=users)
	Users []UserCreateInput `json:"users,omitempty" jsonschema_description:"Список пользователей для создания"`

	// Roles типизированные данные для создания/обновления ролей (action=create/update, layer=roles)
	Roles []RoleCreateInput `json:"roles,omitempty" jsonschema_description:"Список ролей для создания или обновления"`
}

// AdminUsersFilter фильтры для admin_users
type AdminUsersFilter struct {
	// Limit количество результатов на странице
	Limit int `json:"limit,omitempty" jsonschema_description:"Лимит результатов на странице (по умолчанию 50)"`

	// Page номер страницы
	Page int `json:"page,omitempty" jsonschema_description:"Номер страницы (начиная с 1)"`

	// Name фильтрация по имени (client-side, выполняется на стороне сервиса после получения данных от API)
	Name string `json:"name,omitempty" jsonschema_description:"Фильтр по имени пользователя/роли (client-side: API не поддерживает, сервис фильтрует локально)"`

	// Email фильтрация по email (client-side)
	Email string `json:"email,omitempty" jsonschema_description:"Фильтр по email пользователя (client-side: API не поддерживает, сервис фильтрует локально)"`

	// Order сортировка результатов: ключ — поле (created_at, updated_at, id), значение — asc/desc
	Order map[string]string `json:"order,omitempty" jsonschema_description:"Сортировка: {\"created_at\": \"desc\"} — последние добавленные первыми"`
}

// UserCreateInput данные для создания пользователя
type UserCreateInput struct {
	// Name имя пользователя (обязательное)
	Name string `json:"name" jsonschema_description:"Имя пользователя (обязательное)"`

	// Email email пользователя (обязательное)
	Email string `json:"email" jsonschema_description:"Email пользователя (обязательное)"`

	// Password пароль пользователя (обязательное при создании)
	Password string `json:"password" jsonschema_description:"Пароль пользователя"`

	// Lang язык интерфейса (ru, en, es и т.д.)
	Lang string `json:"lang,omitempty" jsonschema_description:"Язык интерфейса: ru, en, es"`
}

// RoleCreateInput данные для создания или обновления роли
type RoleCreateInput struct {
	// ID идентификатор роли (для обновления)
	ID int `json:"id,omitempty" jsonschema_description:"ID роли (только для обновления)"`

	// Name название роли (обязательное при создании)
	Name string `json:"name" jsonschema_description:"Название роли (обязательное)"`
}
