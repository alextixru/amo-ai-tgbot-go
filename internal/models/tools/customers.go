package tools

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

	// With связанные сущности для обогащения ответа
	With []string `json:"with,omitempty" jsonschema_description:"Связанные сущности: catalog_elements, contacts, companies, segments"`

	// Filter параметры поиска
	Filter *CustomerFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска"`

	// Data данные покупателя (для create, update)
	Data *CustomerData `json:"data,omitempty" jsonschema_description:"Данные покупателя"`

	// Batch массив покупателей (для батч-создания/обновления)
	Batch []*CustomerData `json:"batch,omitempty" jsonschema_description:"Массив покупателей для батч-создания/обновления"`

	// Points количество баллов (для earn_points, redeem_points)
	Points int `json:"points,omitempty" jsonschema_description:"Количество бонусных баллов"`

	// TransactionData данные транзакции
	TransactionData *CustomerTransactionData `json:"transaction_data,omitempty" jsonschema_description:"Данные транзакции"`

	// LinkData данные для привязки
	LinkData *CustomerLinkData `json:"link_data,omitempty" jsonschema_description:"Данные для привязки сущностей"`
}

// CustomerFilter фильтры поиска покупателей
type CustomerFilter struct {
	Page                 int      `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Limit                int      `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	Query                string   `json:"query,omitempty" jsonschema_description:"Поисковый запрос"`
	ResponsibleUserNames []string `json:"responsible_user_names,omitempty" jsonschema_description:"Имена ответственных пользователей"`
	IDs                  []int    `json:"ids,omitempty" jsonschema_description:"ID покупателей для фильтрации"`
	StatusNames          []string `json:"status_names,omitempty" jsonschema_description:"Названия статусов покупателей"`
	Names                []string `json:"names,omitempty" jsonschema_description:"Имена покупателей для поиска"`
	// NextDateFrom дата следующей покупки от в формате ISO 8601, например '2024-03-01T00:00:00Z'.
	// ВНИМАНИЕ: этот фильтр не поддерживается текущей версией API — будет возвращена ошибка.
	NextDateFrom string `json:"next_date_from,omitempty" jsonschema_description:"Дата следующей покупки от (ISO 8601, например 2024-03-01T00:00:00Z). Не поддерживается API — вернёт ошибку."`
	// NextDateTo дата следующей покупки до в формате ISO 8601.
	// ВНИМАНИЕ: этот фильтр не поддерживается текущей версией API — будет возвращена ошибка.
	NextDateTo string `json:"next_date_to,omitempty" jsonschema_description:"Дата следующей покупки до (ISO 8601, например 2024-03-31T23:59:59Z). Не поддерживается API — вернёт ошибку."`
}

// CustomerData данные покупателя
type CustomerData struct {
	Name                string               `json:"name" jsonschema_description:"Имя покупателя"`
	ResponsibleUserName string               `json:"responsible_user_name,omitempty" jsonschema_description:"Имя ответственного пользователя"`
	// NextDate дата следующей покупки в формате ISO 8601, например '2024-06-01T00:00:00Z'
	NextDate            string               `json:"next_date,omitempty" jsonschema_description:"Дата следующей покупки (ISO 8601, например 2024-06-01T00:00:00Z)"`
	NextPrice           int                  `json:"next_price,omitempty" jsonschema_description:"Ожидаемая сумма следующей покупки"`
	StatusName          string               `json:"status_name,omitempty" jsonschema_description:"Название статуса покупателя"`
	Periodicity         int                  `json:"periodicity,omitempty" jsonschema_description:"Периодичность покупок (в днях)"`
	CustomFieldsValues  []CustomerFieldValue `json:"custom_fields_values,omitempty" jsonschema_description:"Значения кастомных полей"`
	TagsToAdd           []string             `json:"tags_to_add,omitempty" jsonschema_description:"Теги для добавления"`
	TagsToDelete        []string             `json:"tags_to_delete,omitempty" jsonschema_description:"Теги для удаления"`
}

// CustomerFieldValue значение кастомного поля покупателя
type CustomerFieldValue struct {
	// FieldCode код поля (строковый идентификатор)
	FieldCode string `json:"field_code" jsonschema_description:"Код кастомного поля"`
	// Value значение поля
	Value string `json:"value" jsonschema_description:"Значение поля"`
	// EnumCode код варианта для полей типа select/multiselect
	EnumCode string `json:"enum_code,omitempty" jsonschema_description:"Код варианта (для select/multiselect полей)"`
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
