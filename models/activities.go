package models

// ActivitiesInput входные параметры для инструмента activities
type ActivitiesInput struct {
	// Parent родительская сущность (ОПЦИОНАЛЬНО)
	Parent *ParentEntity `json:"parent,omitempty" jsonschema_description:"Родительская сущность {type: leads|contacts|companies, id: number}"`

	// Layer тип активности
	Layer string `json:"layer" jsonschema:"enum=tasks,enum=notes,enum=calls,enum=events,enum=files,enum=links,enum=tags,enum=subscriptions,enum=talks" jsonschema_description:"Слой активности"`

	// Action действие
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, complete, link, unlink, subscribe, unsubscribe, close — зависит от layer"`

	// ID идентификатор элемента (для get, update, complete)
	ID int `json:"id,omitempty" jsonschema_description:"ID элемента (для get, update, complete)"`

	// Типизированные Data по layer (только одна заполняется)
	TaskData *TaskData `json:"task_data,omitempty" jsonschema_description:"Данные задачи (layer=tasks)"`
	NoteData *NoteData `json:"note_data,omitempty" jsonschema_description:"Данные примечания (layer=notes)"`
	CallData *CallData `json:"call_data,omitempty" jsonschema_description:"Данные звонка (layer=calls)"`

	// Filter фильтры для поиска (только для layer=tasks, action=list)
	Filter *TasksFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска (layer=tasks, action=list)"`

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

	// TagName название тега (для tags.create)
	TagName string `json:"tag_name,omitempty" jsonschema_description:"Название тега (для tags.create)"`

	// TagID ID тега (для tags.delete)
	TagID int `json:"tag_id,omitempty" jsonschema_description:"ID тега (для tags.delete)"`
}

// TasksFilter критерии поиска задач
type TasksFilter struct {
	Limit             int    `json:"limit,omitempty" jsonschema_description:"Лимит записей (до 50)"`
	ResponsibleUserID []int  `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственных"`
	IsCompleted       *bool  `json:"is_completed,omitempty" jsonschema_description:"Статус завершения (true/false)"`
	TaskTypeID        int    `json:"task_type_id,omitempty" jsonschema_description:"ID типа задачи"`
	DateRange         string `json:"date_range,omitempty" jsonschema:"enum=today,enum=tomorrow,enum=overdue,enum=this_week,enum=next_week" jsonschema_description:"Диапазон дат"`
	Query             string `json:"query,omitempty" jsonschema_description:"Поисковый запрос по тексту"`
}

// TaskData данные для создания/обновления задачи
type TaskData struct {
	Text              string `json:"text" jsonschema_description:"Текст задачи (обязательно)"`
	ResponsibleUserID int    `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственного"`
	TaskTypeID        int    `json:"task_type_id,omitempty" jsonschema_description:"ID типа задачи (1=звонок, 2=встреча, 3=письмо)"`
	Deadline          string `json:"deadline,omitempty" jsonschema_description:"Срок выполнения: 'today', 'tomorrow', 'in 2 hours', 'in 3 days', '2024-01-15', '2024-01-15T14:00'"`
}

// NoteData данные для создания/обновления примечания
type NoteData struct {
	Text     string `json:"text" jsonschema_description:"Текст примечания (обязательно)"`
	NoteType string `json:"note_type,omitempty" jsonschema:"enum=common,enum=call_in,enum=call_out,enum=service_message" jsonschema_description:"Тип примечания (по умолчанию common)"`
}

// CallData данные для создания звонка
type CallData struct {
	Direction  string `json:"direction" jsonschema:"enum=inbound,enum=outbound" jsonschema_description:"Направление звонка: inbound=входящий, outbound=исходящий"`
	Duration   int    `json:"duration" jsonschema_description:"Длительность звонка (секунды)"`
	Source     string `json:"source,omitempty" jsonschema_description:"Источник звонка"`
	Phone      string `json:"phone" jsonschema_description:"Номер телефона"`
	CallResult string `json:"call_result,omitempty" jsonschema_description:"Результат звонка"`
	CallStatus int    `json:"call_status,omitempty" jsonschema_description:"Статус: 1=успех, 2=занято, 3=нет ответа, 4=не удалось, 5=голосовая почта, 6=неправильный номер"`
	UniqueID   string `json:"unique_id,omitempty" jsonschema_description:"Уникальный ID звонка"`
	RecordURL  string `json:"record_url,omitempty" jsonschema_description:"Ссылка на запись звонка"`
}
