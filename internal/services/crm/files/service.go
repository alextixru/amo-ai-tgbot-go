package files

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// FileListResult результат листинга файлов с метаданными пагинации
type FileListResult struct {
	// Items список файлов
	Items []*models.File `json:"items"`

	// Total общее количество элементов (если API вернул)
	Total int `json:"total"`

	// HasMore есть ли ещё страницы
	HasMore bool `json:"has_more"`
}

// Service интерфейс сервиса работы с файлами
type Service interface {
	// ListFiles возвращает список файлов с поддержкой расширенной фильтрации и пагинации
	ListFiles(ctx context.Context, filter *gkitmodels.FileFilter) (*FileListResult, error)

	// GetFile возвращает файл по UUID; withDeleted=true позволяет получить удалённый файл
	GetFile(ctx context.Context, uuid string, withDeleted bool) (*models.File, error)

	// UploadFile загружает файл в amoCRM Drive
	UploadFile(ctx context.Context, params services.FileUploadParams) (*models.File, error)

	// UpdateFile переименовывает файл
	UpdateFile(ctx context.Context, uuid, name string) (*models.File, error)

	// DeleteFiles удаляет один или несколько файлов по UUID
	DeleteFiles(ctx context.Context, uuids []string) error
}

type service struct {
	sdk *amocrm.SDK
}

func NewService(sdk *amocrm.SDK) Service {
	return &service{sdk: sdk}
}
