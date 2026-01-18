package files

import (
	"context"
	"strings"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) ListFiles(ctx context.Context, filter *gkitmodels.FileFilter) ([]*models.File, error) {
	f := filters.NewFilesFilter()
	if filter != nil {
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		}
		if filter.Limit > 0 {
			f.SetLimit(filter.Limit)
		}
		if len(filter.UUIDs) > 0 {
			f.SetUUID(filter.UUIDs)
		}
		if len(filter.With) > 0 {
			f.SetWith(strings.Join(filter.With, ","))
		}
	}
	files, _, err := s.sdk.Files().Get(ctx, f)
	return files, err
}

func (s *service) GetFile(ctx context.Context, uuid string) (*models.File, error) {
	return s.sdk.Files().GetOneByUUID(ctx, uuid)
}

func (s *service) UploadFile(ctx context.Context, params services.FileUploadParams) (*models.File, error) {
	return s.sdk.Files().UploadOne(ctx, params)
}

func (s *service) UpdateFile(ctx context.Context, uuid, name string) (*models.File, error) {
	file := &models.File{
		UUID: uuid,
		Name: name,
	}
	return s.sdk.Files().UpdateOne(ctx, file)
}

func (s *service) DeleteFile(ctx context.Context, uuid string) error {
	return s.sdk.Files().DeleteOne(ctx, uuid)
}

func (s *service) DeleteFiles(ctx context.Context, uuids []string) error {
	var files []*models.File
	for _, uuid := range uuids {
		files = append(files, &models.File{UUID: uuid})
	}
	return s.sdk.Files().Delete(ctx, files)
}
