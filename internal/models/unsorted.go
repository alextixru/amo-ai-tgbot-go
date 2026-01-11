package models

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
	SourceUID  string         `json:"source_uid,omitempty" jsonschema_description:"Уникальный идентификатор источника"`
	SourceName string         `json:"source_name,omitempty" jsonschema_description:"Название источника"`
	PipelineID int            `json:"pipeline_id,omitempty" jsonschema_description:"ID воронки для создания сделки"`
	CreatedAt  int            `json:"created_at,omitempty" jsonschema_description:"Unix timestamp создания"`
	Data       map[string]any `json:"data,omitempty" jsonschema_description:"Дополнительные данные (custom_fields_values, tags и т.д.)"`
}

// UnsortedFilter фильтры поиска неразобранного
type UnsortedFilter struct {
	Page       int      `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Limit      int      `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Category   []string `json:"category,omitempty" jsonschema_description:"Категории: sip, mail, forms, chats"`
	PipelineID []int    `json:"pipeline_id,omitempty" jsonschema_description:"ID воронок"`
}

// UnsortedAcceptParams параметры принятия неразобранного
type UnsortedAcceptParams struct {
	UserID   int `json:"user_id,omitempty" jsonschema_description:"ID ответственного пользователя"`
	StatusID int `json:"status_id,omitempty" jsonschema_description:"ID статуса для создаваемой сделки"`
}

// UnsortedLinkData данные для привязки неразобранного
type UnsortedLinkData struct {
	LeadID int `json:"lead_id" jsonschema_description:"ID существующей сделки для привязки"`
}
