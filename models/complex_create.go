package models

// ComplexCreateInput входные параметры для создания сделки с контактами/компанией
type ComplexCreateInput struct {
	// Lead данные сделки
	Lead LeadData `json:"lead" jsonschema_description:"Данные сделки (обязательно)"`

	// Contacts контакты для привязки
	Contacts []ContactData `json:"contacts,omitempty" jsonschema_description:"Контакты для создания и привязки"`

	// Company компания для привязки
	Company *CompanyData `json:"company,omitempty" jsonschema_description:"Компания для создания и привязки"`
}

// LeadData данные сделки
type LeadData struct {
	Name              string `json:"name" jsonschema_description:"Название сделки"`
	Price             int    `json:"price,omitempty" jsonschema_description:"Бюджет"`
	PipelineID        int    `json:"pipeline_id,omitempty" jsonschema_description:"ID воронки"`
	StatusID          int    `json:"status_id,omitempty" jsonschema_description:"ID статуса"`
	ResponsibleUserID int    `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственного"`
}

// ContactData данные контакта
type ContactData struct {
	Name   string `json:"name" jsonschema_description:"Имя контакта"`
	Phone  string `json:"phone,omitempty" jsonschema_description:"Телефон"`
	Email  string `json:"email,omitempty" jsonschema_description:"Email"`
	IsMain bool   `json:"is_main,omitempty" jsonschema_description:"Основной контакт"`
}

// CompanyData данные компании
type CompanyData struct {
	Name string `json:"name" jsonschema_description:"Название компании"`
}
