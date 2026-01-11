package models

// FilesInput входные параметры для инструмента files
type FilesInput struct {
	// Action действие: list, get, upload, update, delete
	Action string `json:"action" jsonschema_description:"Действие: list, get, upload, update, delete"`

	// UUID идентификатор файла (для get, delete)
	UUID string `json:"uuid,omitempty" jsonschema_description:"UUID файла"`

	// UUIDs массив идентификаторов файлов (для batch delete)
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
	Page  int      `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Limit int      `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	UUIDs []string `json:"uuids,omitempty" jsonschema_description:"Фильтр по UUID файлов"`
	With  []string `json:"with,omitempty" jsonschema_description:"Включить дополнительные данные: deleted (удалённые), unbilled (неоплаченные)"`
}

// FileUploadParams параметры загрузки файла
type FileUploadParams struct {
	LocalPath   string `json:"local_path,omitempty" jsonschema_description:"Путь к локальному файлу"`
	FileName    string `json:"file_name,omitempty" jsonschema_description:"Имя файла (переопределить)"`
	WithPreview bool   `json:"with_preview,omitempty" jsonschema_description:"Создать превью"`
	FileUUID    string `json:"file_uuid,omitempty" jsonschema_description:"UUID существующего файла для загрузки новой версии"`
}

// FileUpdateData параметры обновления файла
type FileUpdateData struct {
	UUID string `json:"uuid" jsonschema_description:"UUID файла для обновления"`
	Name string `json:"name" jsonschema_description:"Новое имя файла"`
}
