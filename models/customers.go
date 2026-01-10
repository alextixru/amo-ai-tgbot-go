package models

// CustomersInput входные параметры для инструмента customers
type CustomersInput struct {
	// Layer слой: customers, bonus_points, statuses, transactions, segments
	Layer string `json:"layer" jsonschema_description:"Слой: customers, bonus_points, statuses, transactions, segments"`

	// Action действие (зависит от layer)
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, delete, link, earn_points, redeem_points, etc."`

	// CustomerID ID покупателя (для большинства операций)
	CustomerID int `json:"customer_id,omitempty" jsonschema_description:"ID покупателя"`

	// ID идентификатор объекта (для get, delete)
	ID int `json:"id,omitempty" jsonschema_description:"ID объекта (статус, транзакция, сегмент)"`

	// Filter параметры поиска
	Filter *CustomerFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска"`

	// Data данные покупателя (для create, update)
	Data *CustomerData `json:"data,omitempty" jsonschema_description:"Данные покупателя"`

	// Points количество баллов (для earn_points, redeem_points)
	Points int `json:"points,omitempty" jsonschema_description:"Количество бонусных баллов"`

	// TransactionData данные транзакции
	TransactionData *CustomerTransactionData `json:"transaction_data,omitempty" jsonschema_description:"Данные транзакции"`

	// LinkData данные для привязки
	LinkData *CustomerLinkData `json:"link_data,omitempty" jsonschema_description:"Данные для привязки сущностей"`
}

// CustomerFilter фильтры поиска покупателей
type CustomerFilter struct {
	Page               int    `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Limit              int    `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Query              string `json:"query,omitempty" jsonschema_description:"Поисковый запрос"`
	ResponsibleUserIDs []int  `json:"responsible_user_ids,omitempty" jsonschema_description:"ID ответственных"`
}

// CustomerData данные покупателя
type CustomerData struct {
	Name              string `json:"name" jsonschema_description:"Имя покупателя"`
	ResponsibleUserID int    `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственного"`
	NextDate          int64  `json:"next_date,omitempty" jsonschema_description:"Дата следующей покупки (Unix timestamp)"`
	NextPrice         int    `json:"next_price,omitempty" jsonschema_description:"Ожидаемая сумма"`
	StatusID          int    `json:"status_id,omitempty" jsonschema_description:"ID статуса"`
}

// CustomerTransactionData данные транзакции
type CustomerTransactionData struct {
	Price       int    `json:"price" jsonschema_description:"Сумма транзакции"`
	Comment     string `json:"comment,omitempty" jsonschema_description:"Комментарий"`
	AccrueBonus bool   `json:"accrue_bonus,omitempty" jsonschema_description:"Начислить бонусные баллы"`
}

// CustomerLinkData данные для привязки
type CustomerLinkData struct {
	EntityType string `json:"entity_type" jsonschema_description:"Тип сущности: contacts, companies"`
	EntityID   int    `json:"entity_id" jsonschema_description:"ID сущности"`
}
