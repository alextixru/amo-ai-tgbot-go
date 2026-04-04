package admin_integrations

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListWidgets(ctx context.Context, filter *filters.WidgetsFilter) ([]*models.Widget, error) {
	widgets, _, err := s.sdk.Widgets().Get(ctx, filter)
	return widgets, err
}

func (s *service) GetWidget(ctx context.Context, code string) (*models.Widget, error) {
	return s.sdk.Widgets().GetByCode(ctx, code)
}

func (s *service) InstallWidget(ctx context.Context, code string, settings map[string]any) (*models.Widget, error) {
	widget := &models.Widget{
		Code:     code,
		Settings: settings,
	}
	return s.sdk.Widgets().Install(ctx, widget)
}

func (s *service) UninstallWidget(ctx context.Context, code string) error {
	return s.sdk.Widgets().Uninstall(ctx, code)
}
