package complex_create

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
)

func (s *service) CreateComplex(ctx context.Context, lead *models.Lead) (*services.ComplexLeadResult, error) {
	return s.sdk.Leads().AddOneComplex(ctx, lead)
}

func (s *service) CreateComplexBatch(ctx context.Context, leads []*models.Lead) ([]services.ComplexLeadResult, error) {
	return s.sdk.Leads().AddComplex(ctx, leads)
}
