package models

// FilesInput входные параметры для инструмента files
type FilesInput struct {
	// Action действие: list, get, delete, upload
	Action string `json:"action" jsonschema_description:"Действие: list, get, delete, upload"`

	// UUID идентификатор файла (для get, delete)
	UUID string `json:"uuid,omitempty" jsonschema_description:"UUID файла"`

	// Filter параметры поиска (для list)
	Filter *FileFilter `json:"filter,omitempty" jsonschema_description:"Фильтры поиска файлов"`

	// UploadParams параметры загрузки (для upload)
	UploadParams *FileUploadParams `json:"upload_params,omitempty" jsonschema_description:"Параметры загрузки файла"`
}

// FileFilter фильтры поиска файлов
type FileFilter struct {
	Page  int      `json:"page,omitempty" jsonschema_description:"Номер страницы"`
	Limit int      `json:"limit,omitempty" jsonschema_description:"Лимит результатов"`
	UUIDs []string `json:"uuids,omitempty" jsonschema_description:"Фильтр по UUID файлов"`
}

// FileUploadParams параметры загрузки файла
type FileUploadParams struct {
	LocalPath   string `json:"local_path,omitempty" jsonschema_description:"Путь к локальному файлу"`
	FileName    string `json:"file_name,omitempty" jsonschema_description:"Имя файла (переопределить)"`
	WithPreview bool   `json:"with_preview,omitempty" jsonschema_description:"Создать превью"`
}
