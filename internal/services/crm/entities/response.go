package entities

// EntityResult read-friendly модель для LLM.
// Содержит человекочитаемые поля вместо числовых ID.
// Технические поля (account_id, _links) убраны.
type EntityResult struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`

	// Для leads
	Price        int    `json:"price,omitempty"`
	PipelineName string `json:"pipeline_name,omitempty"`
	StatusName   string `json:"status_name,omitempty"`
	LossReason   string `json:"loss_reason,omitempty"`
	SourceName   string `json:"source_name,omitempty"`
	ClosedAt     string `json:"closed_at,omitempty"`

	// Для contacts
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`

	// Общие
	ResponsibleUserName string `json:"responsible_user_name,omitempty"`
	CreatedByName       string `json:"created_by_name,omitempty"`
	UpdatedByName       string `json:"updated_by_name,omitempty"`

	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`

	Tags               []string           `json:"tags,omitempty"`
	CustomFieldsValues []CustomFieldEntry `json:"custom_fields_values,omitempty"`

	// Связанные сущности
	Contacts  []EntityRef `json:"contacts,omitempty"`
	Companies []EntityRef `json:"companies,omitempty"`
	Leads     []EntityRef `json:"leads,omitempty"`
}

// CustomFieldEntry значение кастомного поля в читаемом формате
type CustomFieldEntry struct {
	FieldCode string `json:"field_code,omitempty"`
	FieldName string `json:"field_name,omitempty"`
	Values    []any  `json:"values"`
}

// EntityRef ссылка на связанную сущность
type EntityRef struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// SearchResult результат поиска с метаданными пагинации
type SearchResult struct {
	Items   []*EntityResult `json:"items"`
	HasMore bool            `json:"has_more"`
	Page    int             `json:"page"`
}

// LinkResult читаемый результат операций link/unlink
type LinkResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
