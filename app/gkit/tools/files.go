package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// filesSchemas содержит полные схемы параметров для каждого action tool files.
// Возвращается LLM при первом вызове без обязательных полей (Shadow Tool — schema mode).
var filesSchemas = map[string]any{
	"list": map[string]any{
		"schema":      true,
		"tool":        "files",
		"action":      "list",
		"description": "Получить список файлов из amoCRM Drive с поддержкой фильтрации и пагинации.",
		"required_fields": map[string]any{},
		"optional_fields": map[string]any{
			"filter": map[string]any{
				"type":        "object",
				"description": "Параметры фильтрации",
				"fields": map[string]any{
					"page":        map[string]any{"type": "integer", "description": "Номер страницы (начиная с 1)"},
					"limit":       map[string]any{"type": "integer", "description": "Лимит результатов на странице"},
					"uuids":       map[string]any{"type": "array[string]", "description": "Фильтр по UUID файлов"},
					"name":        map[string]any{"type": "string", "description": "Поиск по имени файла"},
					"term":        map[string]any{"type": "string", "description": "Полнотекстовый поиск по содержимому и имени"},
					"extensions":  map[string]any{"type": "array[string]", "description": "Фильтр по расширениям файлов (pdf, xlsx, jpg и т.д.)"},
					"deleted":     map[string]any{"type": "boolean", "description": "Включить удалённые файлы в результаты"},
					"date_from":   map[string]any{"type": "string", "description": "Начало диапазона дат (RFC3339, например 2024-01-01T00:00:00Z)"},
					"date_to":     map[string]any{"type": "string", "description": "Конец диапазона дат (RFC3339)"},
					"date_preset": map[string]any{"type": "string", "description": "Пресет периода: today, yesterday, week, month"},
					"size_from":   map[string]any{"type": "integer", "description": "Минимальный размер файла в байтах"},
					"size_to":     map[string]any{"type": "integer", "description": "Максимальный размер файла в байтах"},
				},
			},
		},
		"example": map[string]any{
			"action": "list",
			"filter": map[string]any{
				"extensions":  []string{"pdf", "xlsx"},
				"date_preset": "week",
				"limit":       20,
			},
		},
	},
	"get": map[string]any{
		"schema":      true,
		"tool":        "files",
		"action":      "get",
		"description": "Получить файл по UUID. Поддерживает получение удалённых файлов через filter.deleted=true.",
		"required_fields": map[string]any{
			"uuid": map[string]any{"type": "string", "description": "UUID файла"},
		},
		"optional_fields": map[string]any{
			"filter": map[string]any{
				"type":        "object",
				"description": "Дополнительные параметры",
				"fields": map[string]any{
					"deleted": map[string]any{"type": "boolean", "description": "true — искать среди удалённых файлов"},
				},
			},
		},
		"example": map[string]any{
			"action": "get",
			"uuid":   "550e8400-e29b-41d4-a716-446655440000",
		},
	},
	"upload": map[string]any{
		"schema":      true,
		"tool":        "files",
		"action":      "upload",
		"description": "Загрузить файл в amoCRM Drive. Файл должен быть доступен по локальному пути на сервере.",
		"required_fields": map[string]any{
			"upload_params": map[string]any{
				"type":        "object",
				"description": "Параметры загрузки",
				"fields": map[string]any{
					"local_path": map[string]any{"type": "string", "description": "Путь к локальному файлу на сервере (обязательно)"},
				},
			},
		},
		"optional_fields": map[string]any{
			"upload_params": map[string]any{
				"type":        "object",
				"description": "Дополнительные параметры загрузки",
				"fields": map[string]any{
					"file_name":    map[string]any{"type": "string", "description": "Имя файла (если нужно переопределить автоматически определённое из пути)"},
					"with_preview": map[string]any{"type": "boolean", "description": "Создать превью (для изображений)"},
					"file_uuid":    map[string]any{"type": "string", "description": "UUID существующего файла для загрузки новой версии"},
				},
			},
		},
		"example": map[string]any{
			"action": "upload",
			"upload_params": map[string]any{
				"local_path":   "/tmp/report.pdf",
				"file_name":    "monthly_report.pdf",
				"with_preview": false,
			},
		},
	},
	"update": map[string]any{
		"schema":      true,
		"tool":        "files",
		"action":      "update",
		"description": "Переименовать файл в amoCRM Drive.",
		"required_fields": map[string]any{
			"update_data": map[string]any{
				"type":        "object",
				"description": "Данные для обновления",
				"fields": map[string]any{
					"uuid": map[string]any{"type": "string", "description": "UUID файла для переименования (обязательно)"},
					"name": map[string]any{"type": "string", "description": "Новое имя файла (обязательно)"},
				},
			},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action": "update",
			"update_data": map[string]any{
				"uuid": "550e8400-e29b-41d4-a716-446655440000",
				"name": "new_filename.pdf",
			},
		},
	},
	"delete": map[string]any{
		"schema":      true,
		"tool":        "files",
		"action":      "delete",
		"description": "Удалить один или несколько файлов из amoCRM Drive по UUID.",
		"required_fields": map[string]any{
			"uuid_or_uuids": map[string]any{
				"type":        "string | array[string]",
				"description": "uuid (строка) — для одного файла; uuids (массив строк) — для массового удаления; можно указать оба одновременно",
			},
		},
		"optional_fields": map[string]any{},
		"example": map[string]any{
			"action": "delete",
			"uuids":  []string{"uuid-1", "uuid-2", "uuid-3"},
		},
	},
}

// RegisterFilesTool регистрирует инструмент для работы с файловым хранилищем amoCRM Drive.
// Shadow Tool паттерн: LLM видит минимальную схему (только action).
// При вызове без обязательных полей — возвращает полную схему.
// При вызове с обязательными полями — выполняет действие.
func (r *Registry) RegisterFilesTool() {
	r.addTool(genkit.DefineTool[map[string]any, any](
		r.g,
		"files",
		"Файловое хранилище amoCRM Drive. "+
			"Actions: list (список файлов), get (по UUID), upload (загрузка), update (переименование), delete (удаление). "+
			"Вызови с action чтобы получить схему параметров.",
		func(ctx *ai.ToolContext, input map[string]any) (any, error) {
			return r.handleDriveFilesShadow(ctx.Context, input)
		},
	))
}

// handleDriveFilesShadow реализует Shadow Tool логику для tool files.
// Schema mode: если обязательные поля для action отсутствуют → возвращает полную схему.
// Execute mode: если обязательные поля присутствуют → выполняет действие.
func (r *Registry) handleDriveFilesShadow(ctx context.Context, input map[string]any) (any, error) {
	action, _ := input["action"].(string)
	if action == "" {
		return map[string]any{
			"schema":           true,
			"tool":             "files",
			"error":            "action is required",
			"available_actions": []string{"list", "get", "upload", "update", "delete"},
			"hint":             "Вызови с action чтобы получить схему параметров для нужного действия",
		}, nil
	}

	// Определяем режим работы по наличию обязательных полей для action
	if isFilesSchemaMode(action, input) {
		schema, ok := filesSchemas[action]
		if !ok {
			return nil, fmt.Errorf("unknown action: %s (available: list, get, upload, update, delete)", action)
		}
		return schema, nil
	}

	// Execute mode: десериализуем map в FilesInput через JSON roundtrip
	filedInput, err := mapToFilesInput(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	return r.handleDriveFiles(ctx, filedInput)
}

// isFilesSchemaMode определяет режим работы по наличию обязательных полей.
// Возвращает true (Schema mode) если обязательные поля для action отсутствуют.
func isFilesSchemaMode(action string, input map[string]any) bool {
	switch action {
	case "list":
		// list не имеет обязательных полей — всегда Execute mode
		return false

	case "get":
		uuid, _ := input["uuid"].(string)
		return uuid == ""

	case "upload":
		uploadParams, ok := input["upload_params"].(map[string]any)
		if !ok {
			return true
		}
		localPath, _ := uploadParams["local_path"].(string)
		return localPath == ""

	case "update":
		updateData, ok := input["update_data"].(map[string]any)
		if !ok {
			return true
		}
		uuid, _ := updateData["uuid"].(string)
		name, _ := updateData["name"].(string)
		return uuid == "" || name == ""

	case "delete":
		uuid, _ := input["uuid"].(string)
		uuids, _ := input["uuids"].([]any)
		return uuid == "" && len(uuids) == 0

	default:
		// Неизвестный action — Execute mode вернёт ошибку
		return false
	}
}

// mapToFilesInput конвертирует map[string]any в FilesInput через JSON roundtrip.
func mapToFilesInput(input map[string]any) (gkitmodels.FilesInput, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return gkitmodels.FilesInput{}, err
	}
	var result gkitmodels.FilesInput
	if err := json.Unmarshal(data, &result); err != nil {
		return gkitmodels.FilesInput{}, err
	}
	return result, nil
}

func (r *Registry) handleDriveFiles(ctx context.Context, input gkitmodels.FilesInput) (any, error) {
	switch input.Action {
	case "list":
		return r.filesService.ListFiles(ctx, input.Filter)

	case "get":
		if input.UUID == "" {
			return nil, fmt.Errorf("uuid is required for action 'get'")
		}
		withDeleted := input.Filter != nil && input.Filter.Deleted
		return r.filesService.GetFile(ctx, input.UUID, withDeleted)

	case "upload":
		if input.UploadParams == nil {
			return nil, fmt.Errorf("upload_params is required for action 'upload'")
		}
		if input.UploadParams.LocalPath == "" {
			return nil, fmt.Errorf("upload_params.local_path is required")
		}
		params := services.FileUploadParams{
			LocalPath:   input.UploadParams.LocalPath,
			FileName:    input.UploadParams.FileName,
			WithPreview: input.UploadParams.WithPreview,
			FileUUID:    input.UploadParams.FileUUID,
		}
		return r.filesService.UploadFile(ctx, params)

	case "update":
		if input.UpdateData == nil {
			return nil, fmt.Errorf("update_data is required for action 'update'")
		}
		if input.UpdateData.UUID == "" || input.UpdateData.Name == "" {
			return nil, fmt.Errorf("uuid and name are required in update_data")
		}
		return r.filesService.UpdateFile(ctx, input.UpdateData.UUID, input.UpdateData.Name)

	case "delete":
		uuids := normalizeDeleteUUIDs(input.UUID, input.UUIDs)
		if len(uuids) == 0 {
			return nil, fmt.Errorf("uuid or uuids is required for action 'delete'")
		}
		return nil, r.filesService.DeleteFiles(ctx, uuids)

	default:
		return nil, fmt.Errorf("unknown action: %s (available: list, get, upload, update, delete)", input.Action)
	}
}

// normalizeDeleteUUIDs объединяет одиночный UUID и массив UUIDs в единый список,
// исключая дубли.
func normalizeDeleteUUIDs(single string, batch []string) []string {
	seen := make(map[string]struct{}, len(batch)+1)
	result := make([]string, 0, len(batch)+1)

	if single != "" {
		seen[single] = struct{}{}
		result = append(result, single)
	}
	for _, u := range batch {
		if u == "" {
			continue
		}
		if _, ok := seen[u]; !ok {
			seen[u] = struct{}{}
			result = append(result, u)
		}
	}
	return result
}
