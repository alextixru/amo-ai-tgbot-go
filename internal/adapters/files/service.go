package files

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models"
)

type Service interface {
	ListFiles(ctx context.Context, filter *gkitmodels.FileFilter) ([]*models.File, error)
	GetFile(ctx context.Context, uuid string) (*models.File, error)
	UploadFile(ctx context.Context, params services.FileUploadParams) (*models.File, error)
	UpdateFile(ctx context.Context, uuid, name string) (*models.File, error)
	DeleteFile(ctx context.Context, uuid string) error
	DeleteFiles(ctx context.Context, uuids []string) error
}

type service struct {
	sdk *amocrm.SDK
}

func NewService(sdk *amocrm.SDK) Service {
	return &service{sdk: sdk}
}
