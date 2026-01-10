package admin_integrations

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

// Service определяет бизнес-логику для работы с интеграциями, виджетами и вебхуками.
type Service interface {
	// Webhooks
	ListWebhooks(ctx context.Context) ([]models.Webhook, error)
	SubscribeWebhook(ctx context.Context, destination string, settings []string) (*models.Webhook, error)
	UnsubscribeWebhook(ctx context.Context, destination string, settings []string) error

	// Widgets
	ListWidgets(ctx context.Context) ([]*models.Widget, error)
	GetWidget(ctx context.Context, code string) (*models.Widget, error)
	CreateWidgets(ctx context.Context, widgets []*models.Widget) ([]*models.Widget, error)
	UpdateWidgets(ctx context.Context, widgets []*models.Widget) ([]*models.Widget, error)
	InstallWidget(ctx context.Context, code string) (*models.Widget, error)
	UninstallWidget(ctx context.Context, code string) error

	// Website Buttons
	ListWebsiteButtons(ctx context.Context) ([]*models.WebsiteButton, error)
	GetWebsiteButton(ctx context.Context, id int) (*models.WebsiteButton, error)
	CreateWebsiteButton(ctx context.Context, req *models.WebsiteButtonCreateRequest) (*models.WebsiteButtonCreateResponse, error)
	UpdateWebsiteButton(ctx context.Context, req *models.WebsiteButtonUpdateRequest) (*models.WebsiteButton, error)

	// Chat Templates
	ListChatTemplates(ctx context.Context) ([]*models.ChatTemplate, error)
	DeleteChatTemplate(ctx context.Context, id int) error
	SendChatTemplateOnReview(ctx context.Context, id int) ([]models.ChatTemplateReview, error)
	UpdateChatTemplateReviewStatus(ctx context.Context, templateID, reviewID int, status string) (*models.ChatTemplateReview, error)

	// Short Links
	ListShortLinks(ctx context.Context) ([]models.ShortLink, error)
	CreateShortLink(ctx context.Context, url string) (models.ShortLink, error)
	DeleteShortLink(ctx context.Context, id int) error
}

type service struct {
	sdk *amocrm.SDK
}

// NewService создает новый экземпляр сервиса интеграций.
func NewService(sdk *amocrm.SDK) Service {
	return &service{
		sdk: sdk,
	}
}
