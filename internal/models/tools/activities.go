package tools

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

	// Типизированные Data по layer (для одиночных операций)
	TaskData *TaskData `json:"task_data,omitempty" jsonschema_description:"Данные задачи (layer=tasks). Используй tasks_data для батч-создания"`
	NoteData *NoteData `json:"note_data,omitempty" jsonschema_description:"Данные примечания (layer=notes). Используй notes_data для батч-создания"`
	CallData *CallData `json:"call_data,omitempty" jsonschema_description:"Данные звонка (layer=calls)"`

	// Массивы для батч-операций (используются если Action=create/link)
	TasksData []TaskData   `json:"tasks_data,omitempty" jsonschema_description:"Массив данных задач для батч-создания (layer=tasks)"`
	NotesData []NoteData   `json:"notes_data,omitempty" jsonschema_description:"Массив данных примечаний для батч-создания (layer=notes)"`
	TagNames  []string     `json:"tag_names,omitempty" jsonschema_description:"Список названий тегов для батч-создания (layer=tags)"`
	LinksTo   []LinkTarget `json:"links_to,omitempty" jsonschema_description:"Цели для батч-связывания (layer=links)"`

	// Фильтры поиска (layer=... , action=list)
	Filter       *TasksFilter  `json:"filter,omitempty" jsonschema_description:"Фильтры для задач (layer=tasks)"`
	EventsFilter *EventsFilter `json:"events_filter,omitempty" jsonschema_description:"Фильтры для событий (layer=events)"`
	NotesFilter  *NotesFilter  `json:"notes_filter,omitempty" jsonschema_description:"Фильтры для примечаний (layer=notes)"`
	TagsFilter   *TagsFilter   `json:"tags_filter,omitempty" jsonschema_description:"Фильтры для тегов (layer=tags)"`
	FilesFilter  *FilesFilter  `json:"files_filter,omitempty" jsonschema_description:"Фильтры для файлов (layer=files)"`
	LinksFilter  *LinksFilter  `json:"links_filter,omitempty" jsonschema_description:"Фильтры для связей (layer=links)"`

	// Специфические параметры действий
	ResultText string   `json:"result_text,omitempty" jsonschema_description:"Текст результата (для tasks.complete)"`
	ForceClose bool     `json:"force_close,omitempty" jsonschema_description:"Принудительное закрытие беседы (для talks.close)"`
	With       []string `json:"with,omitempty" jsonschema_description:"Связанные данные для задач (например: leads, contacts)"`

	// Параметры пользователей (для subscriptions)
	UserNames []string `json:"user_names,omitempty" jsonschema_description:"Имена пользователей (для subscribe)"`
	UserName  string   `json:"user_name,omitempty" jsonschema_description:"Имя пользователя (для unsubscribe)"`

	// Параметры файлов (для files.link/unlink)
	FileUUIDs []string `json:"file_uuids,omitempty" jsonschema_description:"UUID файлов (для files.link)"`
	FileUUID  string   `json:"file_uuid,omitempty" jsonschema_description:"UUID файла (для files.unlink)"`

	// Параметры чатов (для talks.close/get)
	TalkID string `json:"talk_id,omitempty" jsonschema_description:"ID чата (для talks.close и talks.get)"`

	// Одиночная цель для связывания (совместимость)
	LinkTo *LinkTarget `json:"link_to,omitempty" jsonschema_description:"Цель связывания (для links.link/unlink)"`

	// Название тега (для create/delete)
	TagName string `json:"tag_name,omitempty" jsonschema_description:"Название тега (для tags.create и tags.delete — альтернатива tag_id)"`
	TagID   int    `json:"tag_id,omitempty" jsonschema_description:"ID тега (для tags.delete)"`
}

// TasksFilter критерии поиска задач
type TasksFilter struct {
	Limit                  int      `json:"limit,omitempty" jsonschema_description:"Лимит записей (до 50)"`
	Page                   int      `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Order                  string   `json:"order,omitempty" jsonschema:"enum=complete_till,enum=created_at" jsonschema_description:"Поле сортировки (по умолчанию complete_till)"`
	OrderDir               string   `json:"order_dir,omitempty" jsonschema:"enum=asc,enum=desc" jsonschema_description:"Направление сортировки (по умолчанию asc)"`
	IDs                    []int    `json:"ids,omitempty" jsonschema_description:"ID конкретных задач"`
	ResponsibleUserNames   []string `json:"responsible_user_names,omitempty" jsonschema_description:"Имена ответственных пользователей"`
	CreatedByNames         []string `json:"created_by_names,omitempty" jsonschema_description:"Имена создателей задач"`
	IsCompleted            *bool    `json:"is_completed,omitempty" jsonschema_description:"Статус завершения (true/false)"`
	TaskType               string   `json:"task_type,omitempty" jsonschema:"enum=follow_up,enum=meeting" jsonschema_description:"Тип задачи: follow_up=Связаться, meeting=Встреча"`
	DateRange              string   `json:"date_range,omitempty" jsonschema:"enum=today,enum=tomorrow,enum=overdue,enum=this_week,enum=next_week" jsonschema_description:"Диапазон дат (клиентская фильтрация)"`
	Query                  string   `json:"query,omitempty" jsonschema_description:"Поисковый запрос (влияет только на лимиты в задачах)"`
	UpdatedAt              *int64   `json:"updated_at,omitempty" jsonschema_description:"Фильтр по дате изменения — от (Unix timestamp)"`
	UpdatedAtTo            *int64   `json:"updated_at_to,omitempty" jsonschema_description:"Фильтр по дате изменения — до (Unix timestamp)"`
}

// EventsFilter критерии фильтрации событий
type EventsFilter struct {
	Limit          int      `json:"limit,omitempty" jsonschema_description:"Лимит записей (до 100)"`
	Page           int      `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Types          []string `json:"types,omitempty" jsonschema_description:"Типы событий: lead_added, lead_status_changed, contact_added, etc."`
	CreatedByNames []string `json:"created_by_names,omitempty" jsonschema_description:"Имена создателей событий"`
	With           []string `json:"with,omitempty" jsonschema_description:"Дополнительные данные: contact_name, lead_name, company_name, catalog_element_name, customer_name, catalog_name, note"`
}

// NotesFilter критерии фильтрации примечаний
type NotesFilter struct {
	Limit     int      `json:"limit,omitempty" jsonschema_description:"Лимит записей (до 50)"`
	Page      int      `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	IDs       []int    `json:"ids,omitempty" jsonschema_description:"ID конкретных примечаний"`
	NoteTypes []string `json:"note_types,omitempty" jsonschema:"enum=common,enum=call_in,enum=call_out,enum=service_message" jsonschema_description:"Типы примечаний"`
	UpdatedAt *int64   `json:"updated_at,omitempty" jsonschema_description:"Фильтр по дате изменения (timestamp)"`
}

// TagsFilter критерии фильтрации тегов
type TagsFilter struct {
	Limit int    `json:"limit,omitempty" jsonschema_description:"Лимит записей (до 50)"`
	Page  int    `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Query string `json:"query,omitempty" jsonschema_description:"Поиск по названию (частичное совпадение)"`
	Name  string `json:"name,omitempty" jsonschema_description:"Фильтр по точному названию тега"`
	IDs   []int  `json:"ids,omitempty" jsonschema_description:"ID конкретных тегов"`
}

// FilesFilter критерии фильтрации файлов
type FilesFilter struct {
	Limit      int      `json:"limit,omitempty" jsonschema_description:"Лимит записей (до 50)"`
	Page       int      `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Extensions []string `json:"extensions,omitempty" jsonschema_description:"Расширения файлов: pdf, docx, xlsx"`
	Term       string   `json:"term,omitempty" jsonschema_description:"Поисковый термин (имя файла)"`
	UUID       string   `json:"uuid,omitempty" jsonschema_description:"UUID конкретного файла"`
}

// LinksFilter критерии фильтрации связей
type LinksFilter struct {
	ToEntityType string `json:"to_entity_type,omitempty" jsonschema_description:"Тип связанной сущности: leads, contacts, companies, catalog_elements"`
	ToEntityID   int    `json:"to_entity_id,omitempty" jsonschema_description:"ID конкретной связанной сущности"`
}

// TaskData данные для создания/обновления задачи
type TaskData struct {
	Text                string `json:"text" jsonschema_description:"Текст задачи (обязательно)"`
	ResponsibleUserName string `json:"responsible_user_name,omitempty" jsonschema_description:"Имя ответственного пользователя"`
	TaskType            string `json:"task_type,omitempty" jsonschema:"enum=follow_up,enum=meeting" jsonschema_description:"Тип задачи: follow_up=Связаться, meeting=Встреча"`
	Deadline            string `json:"deadline,omitempty" jsonschema_description:"Срок выполнения: 'today', 'tomorrow', 'in 2 hours', 'in 3 days', '2024-01-15', '2024-01-15T14:00'"`
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
	CallStatus int    `json:"call_status,omitempty" jsonschema_description:"Статус: 1=оставить_сообщение, 2=перезвонить, 3=недоступен, 4=занято, 5=неверный_номер, 6=нет_ответа, 7=успешный_звонок"`
	UniqueID   string `json:"unique_id,omitempty" jsonschema_description:"Уникальный ID звонка"`
	RecordURL  string `json:"record_url,omitempty" jsonschema_description:"Ссылка на запись звонка"`
}
