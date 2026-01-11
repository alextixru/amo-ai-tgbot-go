package complex_create

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
)

// Service определяет бизнес-логику для комплексного создания сущностей (сделка + контакты/компания).
type Service interface {
	CreateComplex(ctx context.Context, lead *models.Lead) (*services.ComplexLeadResult, error)
	CreateComplexBatch(ctx context.Context, leads []*models.Lead) ([]services.ComplexLeadResult, error)
}

type service struct {
	sdk *amocrm.SDK
}

// NewService создает новый экземпляр сервиса комплексного создания.
func NewService(sdk *amocrm.SDK) Service {
	return &service{
		sdk: sdk,
	}
}
