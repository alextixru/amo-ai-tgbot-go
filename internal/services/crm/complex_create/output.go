package complex_create

// ComplexCreateResult — обогащённый ответ для LLM после комплексного создания.
// Содержит имена вместо числовых ID.
type ComplexCreateResult struct {
	// Lead созданная сделка
	Lead CreatedLeadView `json:"lead"`

	// Contacts созданные контакты (если были)
	Contacts []CreatedContactView `json:"contacts,omitempty"`

	// Company созданная компания (если была)
	Company *CreatedCompanyView `json:"company,omitempty"`

	// Merged был ли лид объединён с существующим
	Merged bool `json:"merged,omitempty"`
}

// CreatedLeadView — представление созданной сделки для LLM.
type CreatedLeadView struct {
	ID                  int    `json:"id"`
	Name                string `json:"name,omitempty"`
	Price               int    `json:"price,omitempty"`
	PipelineName        string `json:"pipeline_name,omitempty"`
	StatusName          string `json:"status_name,omitempty"`
	ResponsibleUserName string `json:"responsible_user_name,omitempty"`
	CreatedAt           string `json:"created_at,omitempty"`
}

// CreatedContactView — представление созданного контакта для LLM.
type CreatedContactView struct {
	ID                  int    `json:"id"`
	Name                string `json:"name,omitempty"`
	ResponsibleUserName string `json:"responsible_user_name,omitempty"`
}

// CreatedCompanyView — представление созданной компании для LLM.
type CreatedCompanyView struct {
	ID                  int    `json:"id"`
	Name                string `json:"name,omitempty"`
	ResponsibleUserName string `json:"responsible_user_name,omitempty"`
}
