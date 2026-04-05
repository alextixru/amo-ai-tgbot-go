package files

import (
	"context"
	"fmt"
	"time"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) ListFiles(ctx context.Context, filter *gkitmodels.FileFilter) (*FileListResult, error) {
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
		if filter.Name != "" {
			f.SetName(filter.Name)
		}
		if filter.Term != "" {
			f.SetTerm(filter.Term)
		}
		if len(filter.Extensions) > 0 {
			f.SetExtensions(filter.Extensions)
		}
		if filter.Deleted {
			f.SetDeleted(true)
		}
		if filter.DatePreset != "" {
			f.SetDatePreset(filter.DatePreset)
		} else if filter.DateFrom != "" || filter.DateTo != "" {
			from, to, err := parseDateRange(filter.DateFrom, filter.DateTo)
			if err != nil {
				return nil, fmt.Errorf("invalid date range: %w", err)
			}
			f.SetDate(from, to, "")
		}
		if filter.SizeFrom > 0 || filter.SizeTo > 0 {
			var from, to *int
			if filter.SizeFrom > 0 {
				v := filter.SizeFrom
				from = &v
			}
			if filter.SizeTo > 0 {
				v := filter.SizeTo
				to = &v
			}
			f.SetSize(from, to, nil)
		}
	}

	files, meta, err := s.sdk.Files().Get(ctx, f)
	if err != nil {
		return nil, err
	}

	result := &FileListResult{
		Items: files,
	}
	if meta != nil {
		result.Total = meta.TotalItems
		result.HasMore = meta.HasMore
	}

	return result, nil
}

func (s *service) GetFile(ctx context.Context, uuid string, withDeleted bool) (*models.File, error) {
	if !withDeleted {
		return s.sdk.Files().GetOneByUUID(ctx, uuid)
	}

	// Для получения удалённого файла используем ListFiles с фильтром по UUID и deleted=true
	f := filters.NewFilesFilter()
	f.SetUUID([]string{uuid})
	f.SetDeleted(true)

	files, _, err := s.sdk.Files().Get(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("file %s not found (including deleted)", uuid)
	}
	return files[0], nil
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

func (s *service) DeleteFiles(ctx context.Context, uuids []string) error {
	var files []*models.File
	for _, uuid := range uuids {
		files = append(files, &models.File{UUID: uuid})
	}
	return s.sdk.Files().Delete(ctx, files)
}

// parseDateRange разбирает строки RFC3339 в Unix-timestamp для SDK-фильтра.
// Возвращает nil для пустых значений.
func parseDateRange(from, to string) (*int, *int, error) {
	var fromTS, toTS *int

	if from != "" {
		t, err := time.Parse(time.RFC3339, from)
		if err != nil {
			return nil, nil, fmt.Errorf("date_from %q is not valid RFC3339: %w", from, err)
		}
		v := int(t.Unix())
		fromTS = &v
	}

	if to != "" {
		t, err := time.Parse(time.RFC3339, to)
		if err != nil {
			return nil, nil, fmt.Errorf("date_to %q is not valid RFC3339: %w", to, err)
		}
		v := int(t.Unix())
		toTS = &v
	}

	return fromTS, toTS, nil
}
