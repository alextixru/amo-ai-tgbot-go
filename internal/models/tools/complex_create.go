package tools

// ComplexCreateInput входные параметры для создания сделки с контактами/компанией
type ComplexCreateInput struct {
	// Lead данные сделки
	Lead LeadData `json:"lead" jsonschema_description:"Данные сделки (обязательно)"`

	// Contacts контакты для привязки
	Contacts []ContactData `json:"contacts,omitempty" jsonschema_description:"Контакты для создания и привязки"`

	// Company компания для привязки
	Company *CompanyData `json:"company,omitempty" jsonschema_description:"Компания для создания и привязки"`
}

// ComplexCreateBatchInput входные параметры для батч-создания сделок
type ComplexCreateBatchInput struct {
	Items []ComplexCreateInput `json:"items" jsonschema_description:"Массив сделок для создания (до 50 штук)"`
}

// LeadData данные сделки
type LeadData struct {
	Name               string         `json:"name" jsonschema_description:"Название сделки"`
	Price              int            `json:"price,omitempty" jsonschema_description:"Бюджет"`
	PipelineID         int            `json:"pipeline_id,omitempty" jsonschema_description:"ID воронки"`
	StatusID           int            `json:"status_id,omitempty" jsonschema_description:"ID статуса"`
	ResponsibleUserID  int            `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственного"`
	CustomFieldsValues map[string]any `json:"custom_fields_values,omitempty" jsonschema_description:"Кастомные поля сделки. Формат: {field_id: [{value: значение}]}"`
	Tags               []string       `json:"tags,omitempty" jsonschema_description:"Теги сделки (названия)"`
}

// ContactData данные контакта
type ContactData struct {
	Name               string         `json:"name" jsonschema_description:"Имя контакта (полное ФИО)"`
	FirstName          string         `json:"first_name,omitempty" jsonschema_description:"Имя (отдельно)"`
	LastName           string         `json:"last_name,omitempty" jsonschema_description:"Фамилия (отдельно)"`
	Phone              string         `json:"phone,omitempty" jsonschema_description:"Телефон (будет добавлен как кастомное поле PHONE)"`
	Email              string         `json:"email,omitempty" jsonschema_description:"Email (будет добавлен как кастомное поле EMAIL)"`
	IsMain             bool           `json:"is_main,omitempty" jsonschema_description:"Основной контакт"`
	ResponsibleUserID  int            `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственного"`
	CustomFieldsValues map[string]any `json:"custom_fields_values,omitempty" jsonschema_description:"Прочие кастомные поля контакта"`
}

// CompanyData данные компании
type CompanyData struct {
	Name               string         `json:"name" jsonschema_description:"Название компании"`
	ResponsibleUserID  int            `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственного"`
	CustomFieldsValues map[string]any `json:"custom_fields_values,omitempty" jsonschema_description:"Кастомные поля компании"`
}
