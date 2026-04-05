package tools

// UnsortedInput входные параметры для инструмента unsorted
type UnsortedInput struct {
	// Action действие: list, get, accept, decline, link, summary, create
	Action string `json:"action" jsonschema_description:"Действие: list, get, accept, decline, link, summary, create"`

	// UID идентификатор неразобранного (для get, accept, decline, link)
	UID string `json:"uid,omitempty" jsonschema_description:"UID записи неразобранного"`

	// Filter параметры поиска (для list, summary)
	Filter *UnsortedFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска"`

	// AcceptParams параметры принятия (для accept)
	AcceptParams *UnsortedAcceptParams `json:"accept_params,omitempty" jsonschema_description:"Параметры принятия заявки"`

	// DeclineParams параметры отклонения (для decline)
	DeclineParams *UnsortedDeclineParams `json:"decline_params,omitempty" jsonschema_description:"Параметры отклонения заявки"`

	// LinkData данные привязки (для link)
	LinkData *UnsortedLinkData `json:"link_data,omitempty" jsonschema_description:"Данные для привязки к сделке"`

	// CreateData данные для создания (для create)
	CreateData *UnsortedCreateData `json:"create_data,omitempty" jsonschema_description:"Данные для создания записи в Неразобранном"`
}

// UnsortedCreateData данные для создания записи в Неразобранном
type UnsortedCreateData struct {
	Category string               `json:"category" jsonschema_description:"Категория источника: sip, forms, chats. Обязательно."`
	Items    []UnsortedCreateItem `json:"items" jsonschema_description:"Массив создаваемых заявок (батч). Минимум одна."`
}

// UnsortedCreateItem данные одной заявки
type UnsortedCreateItem struct {
	SourceUID    string         `json:"source_uid,omitempty" jsonschema_description:"Уникальный идентификатор источника"`
	SourceName   string         `json:"source_name,omitempty" jsonschema_description:"Название источника"`
	PipelineName string         `json:"pipeline_name,omitempty" jsonschema_description:"Название воронки для создания сделки"`
	CreatedAt    string         `json:"created_at,omitempty" jsonschema_description:"Дата создания в формате RFC3339 (например, 2024-01-15T10:30:00Z)"`
	Data         map[string]any `json:"data,omitempty" jsonschema_description:"Дополнительные данные (custom_fields_values, tags и т.д.)"`
}

// UnsortedFilter фильтры поиска неразобранного
type UnsortedFilter struct {
	Page           int      `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Limit          int      `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Category       []string `json:"category,omitempty" jsonschema_description:"Категории: sip, mail, forms, chats"`
	PipelineName   string   `json:"pipeline_name,omitempty" jsonschema_description:"Название воронки для фильтрации"`
	CreatedAtFrom  string   `json:"created_at_from,omitempty" jsonschema_description:"Начало диапазона даты создания (RFC3339, например 2024-01-01T00:00:00Z)"`
	CreatedAtTo    string   `json:"created_at_to,omitempty" jsonschema_description:"Конец диапазона даты создания (RFC3339, например 2024-01-31T23:59:59Z)"`
	Order          string   `json:"order,omitempty" jsonschema_description:"Сортировка по дате создания: 'created_at asc' или 'created_at desc'"`
}

// UnsortedAcceptParams параметры принятия неразобранного
type UnsortedAcceptParams struct {
	UserName     string `json:"user_name,omitempty" jsonschema_description:"Имя ответственного пользователя"`
	PipelineName string `json:"pipeline_name,omitempty" jsonschema_description:"Название воронки для создаваемой сделки"`
	StatusName   string `json:"status_name,omitempty" jsonschema_description:"Название статуса для создаваемой сделки"`
}

// UnsortedDeclineParams параметры отклонения неразобранного
type UnsortedDeclineParams struct {
	UserName string `json:"user_name,omitempty" jsonschema_description:"Имя пользователя, выполняющего отклонение"`
}

// UnsortedLinkData данные для привязки неразобранного
type UnsortedLinkData struct {
	LeadID int `json:"lead_id" jsonschema_description:"ID существующей сделки для привязки"`
}
