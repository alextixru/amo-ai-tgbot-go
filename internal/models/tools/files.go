package tools

// FilesInput входные параметры для инструмента files
type FilesInput struct {
	// Action действие: list, get, upload, update, delete
	Action string `json:"action" jsonschema_description:"Действие: list, get, upload, update, delete"`

	// UUID идентификатор файла (для get, update, delete одного файла)
	UUID string `json:"uuid,omitempty" jsonschema_description:"UUID файла (для get, update, delete)"`

	// UUIDs массив идентификаторов файлов (для batch delete; если указан UUID, он тоже учитывается)
	UUIDs []string `json:"uuids,omitempty" jsonschema_description:"Список UUID файлов для массового удаления"`

	// Filter параметры поиска (для list)
	Filter *FileFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска файлов"`

	// UploadParams параметры загрузки (для upload)
	UploadParams *FileUploadParams `json:"upload_params,omitempty" jsonschema_description:"Параметры загрузки файла"`

	// UpdateData параметры обновления (для update)
	UpdateData *FileUpdateData `json:"update_data,omitempty" jsonschema_description:"Параметры обновления файла (переименование)"`
}

// FileFilter фильтры поиска файлов
type FileFilter struct {
	// Page номер страницы (начиная с 1)
	Page int `json:"page,omitempty" jsonschema_description:"Номер страницы (начиная с 1)"`

	// Limit лимит результатов на странице
	Limit int `json:"limit,omitempty" jsonschema_description:"Лимит результатов на странице"`

	// UUIDs фильтр по конкретным UUID файлов
	UUIDs []string `json:"uuids,omitempty" jsonschema_description:"Фильтр по UUID файлов"`

	// Name поиск по имени файла (точное или частичное совпадение)
	Name string `json:"name,omitempty" jsonschema_description:"Поиск по имени файла"`

	// Term полнотекстовый поиск по содержимому и имени
	Term string `json:"term,omitempty" jsonschema_description:"Полнотекстовый поиск"`

	// Extensions фильтр по расширениям файлов (например: pdf, xlsx, jpg)
	Extensions []string `json:"extensions,omitempty" jsonschema_description:"Фильтр по расширениям файлов (pdf, xlsx, jpg и т.д.)"`

	// Deleted включить удалённые файлы в результаты
	Deleted bool `json:"deleted,omitempty" jsonschema_description:"Включить удалённые файлы"`

	// DateFrom начало диапазона дат в формате RFC3339 (например: 2024-01-01T00:00:00Z)
	DateFrom string `json:"date_from,omitempty" jsonschema_description:"Начало диапазона дат (RFC3339, например 2024-01-01T00:00:00Z)"`

	// DateTo конец диапазона дат в формате RFC3339
	DateTo string `json:"date_to,omitempty" jsonschema_description:"Конец диапазона дат (RFC3339)"`

	// DatePreset пресет периода: today, yesterday, week, month
	DatePreset string `json:"date_preset,omitempty" jsonschema_description:"Пресет периода: today, yesterday, week, month"`

	// SizeFrom минимальный размер файла в байтах
	SizeFrom int `json:"size_from,omitempty" jsonschema_description:"Минимальный размер файла в байтах"`

	// SizeTo максимальный размер файла в байтах
	SizeTo int `json:"size_to,omitempty" jsonschema_description:"Максимальный размер файла в байтах"`
}

// FileUploadParams параметры загрузки файла
type FileUploadParams struct {
	// LocalPath путь к локальному файлу на сервере
	LocalPath string `json:"local_path,omitempty" jsonschema_description:"Путь к локальному файлу на сервере"`

	// FileName имя файла (переопределить автоматически определённое из пути)
	FileName string `json:"file_name,omitempty" jsonschema_description:"Имя файла (если нужно переопределить)"`

	// WithPreview создать превью файла (для изображений)
	WithPreview bool `json:"with_preview,omitempty" jsonschema_description:"Создать превью (для изображений)"`

	// FileUUID UUID существующего файла — загрузить как новую версию
	FileUUID string `json:"file_uuid,omitempty" jsonschema_description:"UUID существующего файла для загрузки новой версии"`
}

// FileUpdateData параметры обновления файла
type FileUpdateData struct {
	// UUID идентификатор файла для обновления
	UUID string `json:"uuid" jsonschema_description:"UUID файла для обновления"`

	// Name новое имя файла
	Name string `json:"name" jsonschema_description:"Новое имя файла"`
}
