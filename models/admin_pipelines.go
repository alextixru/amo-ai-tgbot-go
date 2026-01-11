package models

// AdminPipelinesInput входные параметры для инструмента admin_pipelines
type AdminPipelinesInput struct {
	// Action действие: list, get, create, update, delete, list_statuses, get_status, create_status, update_status, delete_status
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, delete, list_statuses, get_status, create_status, update_status, delete_status"`

	// PipelineID идентификатор воронки
	PipelineID int `json:"pipeline_id,omitempty" jsonschema_description:"ID воронки (для большинства операций)"`

	// StatusID идентификатор статуса
	StatusID int `json:"status_id,omitempty" jsonschema_description:"ID статуса (для операций со статусами)"`

	// Data данные для create/update
	Data map[string]any `json:"data,omitempty" jsonschema_description:"Данные для создания/обновления. Для pipeline: {name: string, sort: int, is_main: bool, is_unsorted_on: bool}. Для status: {name: string, color: string (например '#fffeb2'), sort: int, type: int}. Для батч-создания статусов: {statuses: [{name, color, sort, type}, ...]}"`
}
