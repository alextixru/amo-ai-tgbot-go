package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) ListFiles(ctx context.Context, parent gkitmodels.ParentEntity) ([]models.FileLink, error) {
	svc := services.NewEntityFilesService(s.sdk.Client(), parent.Type, parent.ID)
	files, _, err := svc.Get(ctx, 1, 50)
	return files, err
}

func (s *service) LinkFiles(ctx context.Context, parent gkitmodels.ParentEntity, fileUUIDs []string) ([]models.FileLink, error) {
	svc := services.NewEntityFilesService(s.sdk.Client(), parent.Type, parent.ID)
	links, _, err := svc.Link(ctx, fileUUIDs)
	return links, err
}

func (s *service) UnlinkFile(ctx context.Context, parent gkitmodels.ParentEntity, fileUUID string) error {
	svc := services.NewEntityFilesService(s.sdk.Client(), parent.Type, parent.ID)
	return svc.Unlink(ctx, fileUUID)
}
