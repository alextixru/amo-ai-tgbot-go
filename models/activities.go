package models

// ActivitiesInput входные параметры для инструмента activities
type ActivitiesInput struct {
	// Parent родительская сущность (ОПЦИОНАЛЬНО)
	Parent *ParentEntity `json:"parent,omitempty" jsonschema_description:"Родительская сущность. Опционально: {type: 'leads'|'contacts'|'companies', id: number}"`

	// Layer тип активности: tasks, notes, calls, events, files, links, tags, subscriptions, talks
	Layer string `json:"layer" jsonschema_description:"Тип: tasks, notes, calls, events, files, links, tags, subscriptions, talks"`

	// Action действие
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, complete, link, unlink, subscribe, unsubscribe, close"`

	// ID идентификатор элемента (для get, update)
	ID int `json:"id,omitempty" jsonschema_description:"ID элемента (для get, update)"`

	// Data данные для создания/обновления
	Data *ActivityData `json:"data,omitempty" jsonschema_description:"Данные (для create, update)"`

	// Filter фильтры для поиска (для action=list)
	Filter *ActivityFilter `json:"filter,omitempty" jsonschema_description:"Фильтры для поиска (только для list)"`

	// ResultText текст результата (для tasks.complete)
	ResultText string `json:"result_text,omitempty" jsonschema_description:"Текст результата (для tasks.complete)"`

	// UserIDs ID пользователей (для subscribe)
	UserIDs []int `json:"user_ids,omitempty" jsonschema_description:"ID пользователей (для subscribe)"`

	// UserID ID пользователя (для unsubscribe)
	UserID int `json:"user_id,omitempty" jsonschema_description:"ID пользователя (для unsubscribe)"`

	// FileUUIDs UUID файлов (для files.link)
	FileUUIDs []string `json:"file_uuids,omitempty" jsonschema_description:"UUID файлов (для files.link)"`

	// FileUUID UUID файла (для files.unlink)
	FileUUID string `json:"file_uuid,omitempty" jsonschema_description:"UUID файла (для files.unlink)"`

	// TalkID ID чата (для talks.close)
	TalkID string `json:"talk_id,omitempty" jsonschema_description:"ID чата (для talks.close)"`

	// LinkTo цель для связывания (для links)
	LinkTo *LinkTarget `json:"link_to,omitempty" jsonschema_description:"Цель связывания (для links.link/unlink)"`
}

// ActivityFilter критерии поиска активностей
type ActivityFilter struct {
	Limit             int    `json:"limit,omitempty" jsonschema_description:"Лимит записей (до 50)"`
	ResponsibleUserID []int  `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственных"`
	IsCompleted       *bool  `json:"is_completed,omitempty" jsonschema_description:"Статус завершения (true/false)"`
	TaskTypeID        []int  `json:"task_type_id,omitempty" jsonschema_description:"ID типов задач"`
	DateRange         string `json:"date_range,omitempty" jsonschema_description:"Диапазон дат: 'today' (сегодня), 'tomorrow' (завтра), 'overdue' (просроченные), 'future' (будущие)"`
	Query             string `json:"query,omitempty" jsonschema_description:"Поисковый запрос по тексту"`
}

// ActivityData данные для создания/обновления активностей
type ActivityData struct {
	// Task fields
	Text              string `json:"text,omitempty" jsonschema_description:"Текст задачи/примечания"`
	CompleteTillAt    int64  `json:"complete_till_at,omitempty" jsonschema_description:"Срок задачи (unix timestamp)"`
	CompleteTill      int64  `json:"complete_till,omitempty" jsonschema_description:"Алиас для complete_till_at"`
	TaskTypeID        int    `json:"task_type_id,omitempty" jsonschema_description:"ID типа задачи"`
	TaskType          int    `json:"task_type,omitempty" jsonschema_description:"Алиас для task_type_id"`
	ResponsibleUserID int    `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственного"`

	// Note fields
	NoteType string `json:"note_type,omitempty" jsonschema_description:"Тип примечания: common, call_in, call_out, etc."`

	// Call fields
	Direction  string `json:"direction,omitempty" jsonschema_description:"Направление звонка: inbound, outbound"`
	Duration   int    `json:"duration,omitempty" jsonschema_description:"Длительность звонка (секунды)"`
	Source     string `json:"source,omitempty" jsonschema_description:"Источник звонка"`
	Phone      string `json:"phone,omitempty" jsonschema_description:"Номер телефона"`
	CallResult string `json:"call_result,omitempty" jsonschema_description:"Результат звонка"`
	CallStatus int    `json:"call_status,omitempty" jsonschema_description:"Статус звонка: 1-успех, и т.д."`
	UniqID     string `json:"uniq,omitempty" jsonschema_description:"Уникальный ID звонка"`
	Link       string `json:"link,omitempty" jsonschema_description:"Ссылка на запись звонка"`

	// Tag fields
	TagName string `json:"tag_name,omitempty" jsonschema_description:"Название тега"`
	TagID   int    `json:"tag_id,omitempty" jsonschema_description:"ID тега"`
}
