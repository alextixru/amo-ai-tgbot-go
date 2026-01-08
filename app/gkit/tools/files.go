package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

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

// registerFilesTool регистрирует инструмент для работы с файлами
func (r *Registry) registerFilesTool() {
	r.addTool(genkit.DefineTool[FilesInput, any](
		r.g,
		"files",
		"Работа с файловым хранилищем amoCRM Drive. "+
			"Actions: list (список файлов), get (файл по UUID), delete (удалить), upload (загрузить). "+
			"Для привязки файла к сущности используй tool 'activities' с layer='files'.",
		func(ctx *ai.ToolContext, input FilesInput) (any, error) {
			return r.handleFiles(ctx.Context, input)
		},
	))
}

func (r *Registry) handleFiles(ctx context.Context, input FilesInput) (any, error) {
	switch input.Action {
	case "list", "search":
		return r.listFiles(ctx, input.Filter)
	case "get":
		if input.UUID == "" {
			return nil, fmt.Errorf("uuid is required for action 'get'")
		}
		return r.sdk.Files().GetOneByUUID(ctx, input.UUID)
	case "delete":
		if input.UUID == "" {
			return nil, fmt.Errorf("uuid is required for action 'delete'")
		}
		err := r.sdk.Files().DeleteOne(ctx, input.UUID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true, "deleted_uuid": input.UUID}, nil
	case "upload":
		if input.UploadParams == nil || input.UploadParams.LocalPath == "" {
			return nil, fmt.Errorf("upload_params.local_path is required for action 'upload'")
		}
		return r.uploadFile(ctx, input.UploadParams)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (r *Registry) listFiles(ctx context.Context, filter *FileFilter) ([]*models.File, error) {
	f := filters.NewFilesFilter()
	f.SetLimit(50)
	f.SetPage(1)
	if filter != nil {
		if filter.Limit > 0 {
			f.SetLimit(filter.Limit)
		}
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		}
		if len(filter.UUIDs) > 0 {
			f.SetUUID(filter.UUIDs)
		}
	}
	items, _, err := r.sdk.Files().Get(ctx, f)
	return items, err
}

func (r *Registry) uploadFile(ctx context.Context, params *FileUploadParams) (*models.File, error) {
	uploadParams := services.FileUploadParams{
		LocalPath:   params.LocalPath,
		FileName:    params.FileName,
		WithPreview: params.WithPreview,
	}
	return r.sdk.Files().UploadOne(ctx, uploadParams)
}
