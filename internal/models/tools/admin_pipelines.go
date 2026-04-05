package tools

import sdkservices "github.com/alextixru/amocrm-sdk-go/core/services"

// StatusColors допустимые hex-цвета для статусов воронки
// Полный список значений, принимаемых API amoCRM
var StatusColors = []string{
	"#fffeb2", "#fffd7f", "#fff000", "#ffeab2", "#ffdc7f",
	"#ffce5a", "#ffdbdb", "#ffc8c8", "#ff8f92", "#d6eaff",
	"#c1e0ff", "#98cbff", "#ebffb1", "#deff81", "#87f2c0",
	"#f9deff", "#f3beff", "#ccc8f9", "#eb93ff", "#f2f3f4",
	"#e6e8ea",
}

// PipelineData данные для создания/обновления воронки
type PipelineData struct {
	// Name название воронки
	Name string `json:"name,omitempty" jsonschema_description:"Название воронки"`

	// Sort порядок сортировки
	Sort int `json:"sort,omitempty" jsonschema_description:"Порядок сортировки"`

	// IsMain является ли основной воронкой
	IsMain bool `json:"is_main,omitempty" jsonschema_description:"Является ли основной воронкой"`

	// IsUnsortedOn включено ли неразобранное
	IsUnsortedOn bool `json:"is_unsorted_on,omitempty" jsonschema_description:"Включить неразобранное"`
}

// StatusData данные для создания/обновления статуса
type StatusData struct {
	// Name название статуса
	Name string `json:"name,omitempty" jsonschema_description:"Название статуса"`

	// Sort порядок сортировки
	Sort int `json:"sort,omitempty" jsonschema_description:"Порядок сортировки"`

	// Color цвет статуса (hex). Допустимые значения: #fffeb2 #fffd7f #fff000 #ffeab2 #ffdc7f #ffce5a #ffdbdb #ffc8c8 #ff8f92 #d6eaff #c1e0ff #98cbff #ebffb1 #deff81 #87f2c0 #f9deff #f3beff #ccc8f9 #eb93ff #f2f3f4 #e6e8ea
	Color string `json:"color,omitempty" jsonschema_description:"Цвет статуса в hex. Допустимые: #fffeb2 #fffd7f #fff000 #ffeab2 #ffdc7f #ffce5a #ffdbdb #ffc8c8 #ff8f92 #d6eaff #c1e0ff #98cbff #ebffb1 #deff81 #87f2c0 #f9deff #f3beff #ccc8f9 #eb93ff #f2f3f4 #e6e8ea"`

	// Type тип статуса: "regular" (обычный), "won" (выигран), "lost" (проигран)
	Type string `json:"type,omitempty" jsonschema_description:"Тип статуса: regular (обычный), won (выигран), lost (проигран)"`
}

// AdminPipelinesInput входные параметры для инструмента admin_pipelines
type AdminPipelinesInput struct {
	// Action действие
	// Воронки: list, get, create, update, delete
	// Статусы: list_statuses, get_status, create_status, update_status, delete_status
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, delete (воронки); list_statuses, get_status, create_status, update_status, delete_status (статусы)"`

	// PipelineID идентификатор воронки (числовой)
	PipelineID int `json:"pipeline_id,omitempty" jsonschema_description:"ID воронки. Альтернатива pipeline_name."`

	// PipelineName имя воронки для поиска по имени (если pipeline_id == 0)
	PipelineName string `json:"pipeline_name,omitempty" jsonschema_description:"Название воронки. Используется для поиска вместо pipeline_id."`

	// StatusID идентификатор статуса (числовой)
	StatusID int `json:"status_id,omitempty" jsonschema_description:"ID статуса. Альтернатива status_name."`

	// StatusName имя статуса для поиска по имени (если status_id == 0)
	StatusName string `json:"status_name,omitempty" jsonschema_description:"Название статуса. Используется для поиска вместо status_id."`

	// WithStatuses включить статусы воронки в ответ (за один запрос)
	WithStatuses bool `json:"with_statuses,omitempty" jsonschema_description:"Включить статусы воронки в ответ (для list и get). Позволяет получить всё за один вызов."`

	// Pipeline данные для создания/обновления одной воронки
	Pipeline *PipelineData `json:"pipeline,omitempty" jsonschema_description:"Данные воронки для create/update"`

	// Status данные для создания/обновления одного статуса
	Status *StatusData `json:"status,omitempty" jsonschema_description:"Данные статуса для create_status/update_status"`

	// Items батч-данные для создания нескольких воронок или статусов
	// Для create: массив PipelineData; для create_status: массив StatusData
	Items []map[string]any `json:"items,omitempty" jsonschema_description:"Батч-режим (data.items): массив PipelineData (для create) или StatusData (для create_status)"`
}

// --- Output types ---

// PipelineOutput результат операции с воронкой, адаптированный для LLM
// Не содержит account_id и _links
type PipelineOutput struct {
	ID           int             `json:"id"`
	Name         string          `json:"name,omitempty"`
	Sort         int             `json:"sort,omitempty"`
	IsMain       bool            `json:"is_main,omitempty"`
	IsUnsortedOn bool            `json:"is_unsorted_on,omitempty"`
	IsArchive    bool            `json:"is_archive,omitempty"`
	Statuses     []*StatusOutput `json:"statuses,omitempty"`
}

// StatusOutput результат операции со статусом, адаптированный для LLM
// Не содержит account_id и _links
type StatusOutput struct {
	ID         int    `json:"id"`
	Name       string `json:"name,omitempty"`
	Sort       int    `json:"sort,omitempty"`
	Color      string `json:"color,omitempty"`
	PipelineID int    `json:"pipeline_id,omitempty"`
	IsEditable bool   `json:"is_editable,omitempty"`

	// Семантические метки типа статуса
	TypeLabel string `json:"type_label,omitempty"` // "regular", "won", "lost"
	IsWon     bool   `json:"is_won,omitempty"`
	IsLost    bool   `json:"is_lost,omitempty"`
	IsClosed  bool   `json:"is_closed,omitempty"`
}

// ListPipelinesOutput результат ListPipelines с метаданными пагинации
type ListPipelinesOutput struct {
	Pipelines []*PipelineOutput      `json:"pipelines"`
	PageMeta  *sdkservices.PageMeta  `json:"page_meta,omitempty"`
}
