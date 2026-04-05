package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func convertFileLinks(links []models.FileLink) *FilesListOutput {
	out := &FilesListOutput{
		Files: make([]FileOutput, 0, len(links)),
	}
	for _, l := range links {
		out.Files = append(out.Files, FileOutput{UUID: l.FileUUID})
	}
	return out
}

func (s *service) ListFiles(ctx context.Context, parent gkitmodels.ParentEntity, filter *gkitmodels.FilesFilter) (*FilesListOutput, error) {
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
	}
	files, _, err := svc.Get(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return convertFileLinks(files), nil
}

func (s *service) LinkFiles(ctx context.Context, parent gkitmodels.ParentEntity, fileUUIDs []string) (*FilesListOutput, error) {
	svc := services.NewEntityFilesService(s.sdk.Client(), parent.Type, parent.ID)
	links, _, err := svc.Link(ctx, fileUUIDs)
	if err != nil {
		return nil, err
	}
	return convertFileLinks(links), nil
}

func (s *service) UnlinkFile(ctx context.Context, parent gkitmodels.ParentEntity, fileUUID string) error {
	svc := services.NewEntityFilesService(s.sdk.Client(), parent.Type, parent.ID)
	return svc.Unlink(ctx, fileUUID)
}
