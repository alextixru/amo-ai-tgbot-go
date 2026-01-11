package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) ListFiles(ctx context.Context, parent gkitmodels.ParentEntity, filter *gkitmodels.FilesFilter) ([]models.FileLink, error) {
	svc := services.NewEntityFilesService(s.sdk.Client(), parent.Type, parent.ID)
	limit := 50
	page := 1
	if filter != nil {
		if filter.Limit > 0 {
			limit = filter.Limit
		}
		if filter.Page > 0 {
			page = filter.Page
		}
		// SDK EntityFilesService usually has restricted filtering compared to global FilesService
	}
	files, _, err := svc.Get(ctx, page, limit)
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
