package activities

// PageMeta метаданные пагинации
type PageMeta struct {
	HasMore bool `json:"has_more"`
	Total   int  `json:"total,omitempty"`
}

// TaskOutput задача в читаемом формате
type TaskOutput struct {
	ID                  int         `json:"id"`
	Text                string      `json:"text,omitempty"`
	EntityID            int         `json:"entity_id,omitempty"`
	EntityType          string      `json:"entity_type,omitempty"`
	TaskType            string      `json:"task_type,omitempty"`
	IsCompleted         bool        `json:"is_completed"`
	Deadline            string      `json:"deadline,omitempty"`
	ResponsibleUserName string      `json:"responsible_user_name,omitempty"`
	CreatedByName       string      `json:"created_by_name,omitempty"`
	UpdatedByName       string      `json:"updated_by_name,omitempty"`
	CreatedAt           string      `json:"created_at,omitempty"`
	UpdatedAt           string      `json:"updated_at,omitempty"`
	Result              *TaskResult `json:"result,omitempty"`
}

// TaskResult результат выполнения задачи
type TaskResult struct {
	Text string `json:"text,omitempty"`
}

// TasksListOutput список задач с пагинацией
type TasksListOutput struct {
	Tasks    []*TaskOutput `json:"tasks"`
	PageMeta PageMeta      `json:"page_meta"`
}

// NoteOutput примечание в читаемом формате
type NoteOutput struct {
	ID            int    `json:"id"`
	EntityID      int    `json:"entity_id,omitempty"`
	NoteType      string `json:"note_type,omitempty"`
	Text          string `json:"text,omitempty"`
	CreatedByName string `json:"created_by_name,omitempty"`
	UpdatedByName string `json:"updated_by_name,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

// CallOutput звонок в читаемом формате
type CallOutput struct {
	ID                  int    `json:"id,omitempty"`
	Direction           string `json:"direction,omitempty"`
	Duration            int    `json:"duration,omitempty"`
	Phone               string `json:"phone,omitempty"`
	CallResult          string `json:"call_result,omitempty"`
	CallStatus          int    `json:"call_status,omitempty"`
	Source              string `json:"source,omitempty"`
	UniqueID            string `json:"unique_id,omitempty"`
	RecordURL           string `json:"record_url,omitempty"`
	ResponsibleUserName string `json:"responsible_user_name,omitempty"`
	CreatedByName       string `json:"created_by_name,omitempty"`
	CreatedAt           string `json:"created_at,omitempty"`
}

// EventOutput событие в читаемом формате
type EventOutput struct {
	ID            string      `json:"id,omitempty"`
	Type          string      `json:"type,omitempty"`
	EntityID      int         `json:"entity_id,omitempty"`
	EntityType    string      `json:"entity_type,omitempty"`
	CreatedByName string      `json:"created_by_name,omitempty"`
	CreatedAt     string      `json:"created_at,omitempty"`
	ValueBefore   interface{} `json:"value_before,omitempty"`
	ValueAfter    interface{} `json:"value_after,omitempty"`
	// With-данные: заполняются при запросе с with=contact_name,lead_name,company_name
	ContactName string `json:"contact_name,omitempty"`
	LeadName    string `json:"lead_name,omitempty"`
	CompanyName string `json:"company_name,omitempty"`
}

// EventsListOutput список событий с пагинацией
type EventsListOutput struct {
	Events   []*EventOutput `json:"events"`
	PageMeta PageMeta       `json:"page_meta"`
}

// FilesListOutput список файлов с пагинацией
type FilesListOutput struct {
	Files    []FileOutput `json:"files"`
	PageMeta PageMeta     `json:"page_meta"`
}

// FileOutput файл в читаемом формате
type FileOutput struct {
	UUID string `json:"uuid,omitempty"`
}

// LinkOutput связь в читаемом формате
type LinkOutput struct {
	ToEntityID   int    `json:"to_entity_id,omitempty"`
	ToEntityType string `json:"to_entity_type,omitempty"`
}

// TagOutput тег в читаемом формате
type TagOutput struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// SubscriptionOutput подписка в читаемом формате
type SubscriptionOutput struct {
	SubscriberName string `json:"subscriber_name,omitempty"`
	SubscriberID   int    `json:"subscriber_id,omitempty"`
}

// SubscriptionsListOutput список подписок
type SubscriptionsListOutput struct {
	Subscriptions []SubscriptionOutput `json:"subscriptions"`
}

// TalkOutput беседа в читаемом формате
type TalkOutput struct {
	TalkID     int    `json:"talk_id,omitempty"`
	ChatID     string `json:"chat_id,omitempty"`
	ContactID  int    `json:"contact_id,omitempty"`
	EntityID   int    `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	Origin     string `json:"origin,omitempty"`
	IsInWork   bool   `json:"is_in_work,omitempty"`
	IsRead     bool   `json:"is_read,omitempty"`
	Rate       int    `json:"rate,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}
