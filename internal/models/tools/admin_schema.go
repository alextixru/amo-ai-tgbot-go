package tools

// AdminSchemaInput входные параметры для инструмента admin_schema
type AdminSchemaInput struct {
	// Layer слой схемы: custom_fields | field_groups | loss_reasons | sources
	Layer string `json:"layer" jsonschema_description:"Слой схемы: custom_fields, field_groups, loss_reasons, sources"`

	// Action действие: list | get | create | update | delete
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, delete. ВАЖНО: update недоступен для loss_reasons (API ограничение)"`

	// EntityType тип сущности (для custom_fields и field_groups): leads | contacts | companies | customers
	EntityType string `json:"entity_type,omitempty" jsonschema_description:"Тип сущности: leads, contacts, companies, customers (для custom_fields и field_groups)"`

	// ID идентификатор элемента (для get, update, delete)
	ID int `json:"id,omitempty" jsonschema_description:"ID элемента (для custom_fields, loss_reasons, sources)"`

	// GroupID идентификатор группы полей (string в API)
	GroupID string `json:"group_id,omitempty" jsonschema_description:"ID группы полей (для field_groups get/update/delete)"`

	// Filter фильтры для list
	Filter *SchemaFilter `json:"filter,omitempty" jsonschema_description:"Фильтры для поиска"`

	// CustomFieldData данные для create/update кастомного поля
	CustomField *CustomFieldData `json:"custom_field,omitempty" jsonschema_description:"Данные кастомного поля для create/update (используется с layer=custom_fields)"`

	// CustomFields список кастомных полей для batch create/update
	CustomFields []CustomFieldData `json:"custom_fields,omitempty" jsonschema_description:"Список кастомных полей для batch create/update (используется с layer=custom_fields)"`

	// FieldGroup данные группы полей для create/update
	FieldGroup *FieldGroupData `json:"field_group,omitempty" jsonschema_description:"Данные группы полей для create/update (используется с layer=field_groups)"`

	// FieldGroups список групп для batch create/update
	FieldGroups []FieldGroupData `json:"field_groups,omitempty" jsonschema_description:"Список групп полей для batch create/update (используется с layer=field_groups)"`

	// LossReason данные причины отказа для create
	LossReason *LossReasonData `json:"loss_reason,omitempty" jsonschema_description:"Данные причины отказа для create (используется с layer=loss_reasons)"`

	// LossReasons список причин отказа для batch create
	LossReasons []LossReasonData `json:"loss_reasons,omitempty" jsonschema_description:"Список причин отказа для batch create (используется с layer=loss_reasons)"`

	// Source данные источника для create/update
	Source *SourceData `json:"source,omitempty" jsonschema_description:"Данные источника для create/update (используется с layer=sources)"`

	// Sources список источников для batch create/update
	Sources []SourceData `json:"sources,omitempty" jsonschema_description:"Список источников для batch create/update (используется с layer=sources)"`
}

// SchemaFilter фильтры для поиска в admin_schema
type SchemaFilter struct {
	Limit int `json:"limit,omitempty" jsonschema_description:"Лимит результатов (по умолчанию 50)"`
	Page  int `json:"page,omitempty" jsonschema_description:"Номер страницы (по умолчанию 1)"`

	// Name фильтр по имени (client-side, частичное совпадение без учёта регистра)
	Name string `json:"name,omitempty" jsonschema_description:"Фильтр по имени (частичное совпадение, без учёта регистра)"`

	// Для custom_fields
	IDs   []int    `json:"ids,omitempty" jsonschema_description:"Фильтр по ID полей (для custom_fields)"`
	Types []string `json:"types,omitempty" jsonschema_description:"Фильтр по типам полей: text, numeric, checkbox, select, multiselect, date, url, textarea, radiobutton, streetaddress, smart_address, birthday, legal_entity, date_time, price, category, items, chained_list, tracking_data, linked_entity, file, payer, supplier, multitext, monetary (для custom_fields)"`

	// Order сортировка для custom_fields: {"created_at": "desc"}, {"updated_at": "asc"}
	Order map[string]string `json:"order,omitempty" jsonschema_description:"Сортировка для custom_fields: ключ — поле (created_at, updated_at, id), значение — asc или desc"`

	// Для sources
	ExternalIDs []string `json:"external_ids,omitempty" jsonschema_description:"Фильтр по external_id источников (для sources)"`
}

// CustomFieldData данные для создания/обновления кастомного поля
type CustomFieldData struct {
	// ID обязателен для update
	ID int `json:"id,omitempty" jsonschema_description:"ID поля (обязателен для update)"`

	// Name название поля
	Name string `json:"name,omitempty" jsonschema_description:"Название поля"`

	// Type тип поля: text, numeric, checkbox, select, multiselect, date, url, textarea, radiobutton, multitext и др.
	Type string `json:"type,omitempty" jsonschema_description:"Тип поля: text, numeric, checkbox, select, multiselect, date, url, textarea, radiobutton, multitext и др."`

	// Code символьный код поля (латиница, уникален в аккаунте)
	Code string `json:"code,omitempty" jsonschema_description:"Символьный код поля (латиница, уникален в аккаунте)"`

	// Sort порядок сортировки
	Sort int `json:"sort,omitempty" jsonschema_description:"Порядок сортировки"`

	// GroupID ID/код группы полей (например 'general', 'statistic' — получить через layer=field_groups)
	GroupID string `json:"group_id,omitempty" jsonschema_description:"ID группы полей. Допустимые значения узнать через action=list, layer=field_groups"`

	// IsAPIOnly только для API (не отображается в интерфейсе)
	IsAPIOnly bool `json:"is_api_only,omitempty" jsonschema_description:"Только для API — поле не отображается в интерфейсе CRM"`

	// IsRequired поле обязательно для заполнения
	IsRequired bool `json:"is_required,omitempty" jsonschema_description:"Поле обязательно для заполнения"`

	// Enums варианты значений (для select, multiselect, radiobutton)
	Enums []CustomFieldEnumData `json:"enums,omitempty" jsonschema_description:"Варианты значений (для select, multiselect, radiobutton)"`
}

// CustomFieldEnumData вариант значения для select/multiselect полей
type CustomFieldEnumData struct {
	// ID обязателен при обновлении существующего варианта
	ID int `json:"id,omitempty" jsonschema_description:"ID варианта (обязателен при обновлении существующего)"`

	// Value отображаемое значение
	Value string `json:"value" jsonschema_description:"Отображаемое значение"`

	// Sort порядок сортировки
	Sort int `json:"sort,omitempty" jsonschema_description:"Порядок сортировки"`

	// Code символьный код варианта
	Code string `json:"code,omitempty" jsonschema_description:"Символьный код варианта"`
}

// FieldGroupData данные для создания/обновления группы полей
type FieldGroupData struct {
	// ID обязателен для update (строковый идентификатор группы)
	ID string `json:"id,omitempty" jsonschema_description:"ID группы (обязателен для update)"`

	// Name название группы
	Name string `json:"name,omitempty" jsonschema_description:"Название группы полей"`

	// Sort порядок сортировки
	Sort int `json:"sort,omitempty" jsonschema_description:"Порядок сортировки"`
}

// LossReasonData данные для создания причины отказа
type LossReasonData struct {
	// Name название причины отказа
	Name string `json:"name" jsonschema_description:"Название причины отказа"`

	// Sort порядок сортировки
	Sort int `json:"sort,omitempty" jsonschema_description:"Порядок сортировки"`
}

// SourceData данные для создания/обновления источника
type SourceData struct {
	// ID обязателен для update
	ID int `json:"id,omitempty" jsonschema_description:"ID источника (обязателен для update)"`

	// Name название источника
	Name string `json:"name,omitempty" jsonschema_description:"Название источника"`

	// ExternalID внешний идентификатор источника
	ExternalID string `json:"external_id,omitempty" jsonschema_description:"Внешний идентификатор источника"`

	// OriginCode код источника
	OriginCode string `json:"origin_code,omitempty" jsonschema_description:"Код источника"`

	// PipelineID ID воронки, к которой привязан источник
	PipelineID int `json:"pipeline_id,omitempty" jsonschema_description:"ID воронки (получить через admin_pipelines)"`

	// Default источник по умолчанию
	Default bool `json:"default,omitempty" jsonschema_description:"Источник по умолчанию"`
}
