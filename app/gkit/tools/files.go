package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (r *Registry) RegisterFilesTool() {
	r.addTool(genkit.DefineTool[gkitmodels.FilesInput, any](
		r.g,
		"files",
		"Работа с файловым хранилищем (amoCRM Drive). "+
			"Поддерживает: list (список), get (получение по UUID), upload (загрузка), delete (удаление).",
		func(ctx *ai.ToolContext, input gkitmodels.FilesInput) (any, error) {
			return r.handleDriveFiles(ctx.Context, input)
		},
	))
}

func (r *Registry) handleDriveFiles(ctx context.Context, input gkitmodels.FilesInput) (any, error) {
	switch input.Action {
	case "list":
		return r.filesService.ListFiles(ctx, input.Filter)
	case "get":
		if input.UUID == "" {
			return nil, fmt.Errorf("uuid is required for action 'get'")
		}
		return r.filesService.GetFile(ctx, input.UUID)
	case "upload":
		if input.UploadParams == nil {
			return nil, fmt.Errorf("upload_params is required for action 'upload'")
		}
		params := services.FileUploadParams{
			LocalPath:   input.UploadParams.LocalPath,
			FileName:    input.UploadParams.FileName,
			WithPreview: input.UploadParams.WithPreview,
		}
		return r.filesService.UploadFile(ctx, params)
	case "delete":
		if input.UUID == "" {
			return nil, fmt.Errorf("uuid is required for action 'delete'")
		}
		return nil, r.filesService.DeleteFile(ctx, input.UUID)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}
