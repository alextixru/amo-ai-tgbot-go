package admin_integrations

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListWidgets(ctx context.Context) ([]*models.Widget, error) {
	widgets, _, err := s.sdk.Widgets().Get(ctx, nil)
	return widgets, err
}

func (s *service) GetWidget(ctx context.Context, code string) (*models.Widget, error) {
	return s.sdk.Widgets().GetByCode(ctx, code)
}

func (s *service) CreateWidgets(ctx context.Context, widgets []*models.Widget) ([]*models.Widget, error) {
	res, _, err := s.sdk.Widgets().Add(ctx, widgets)
	return res, err
}

func (s *service) UpdateWidgets(ctx context.Context, widgets []*models.Widget) ([]*models.Widget, error) {
	res, _, err := s.sdk.Widgets().Update(ctx, widgets)
	return res, err
}

func (s *service) InstallWidget(ctx context.Context, code string) (*models.Widget, error) {
	return s.sdk.Widgets().InstallByCode(ctx, code)
}

func (s *service) UninstallWidget(ctx context.Context, code string) error {
	return s.sdk.Widgets().Uninstall(ctx, code)
}
