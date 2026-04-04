package admin_integrations

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
)

// Service определяет бизнес-логику для работы с интеграциями, виджетами и вебхуками.
type Service interface {
	// Webhooks
	ListWebhooks(ctx context.Context, filter *filters.WebhooksFilter) ([]models.Webhook, error)
	SubscribeWebhook(ctx context.Context, destination string, settings []string) (*models.Webhook, error)
	UnsubscribeWebhook(ctx context.Context, destination string, settings []string) error

	// Widgets
	ListWidgets(ctx context.Context, filter *filters.WidgetsFilter) ([]*models.Widget, error)
	GetWidget(ctx context.Context, code string) (*models.Widget, error)
	InstallWidget(ctx context.Context, code string, settings map[string]any) (*models.Widget, error)
	UninstallWidget(ctx context.Context, code string) error

	// Website Buttons
	ListWebsiteButtons(ctx context.Context, filter *services.WebsiteButtonsFilter, with []string) ([]*models.WebsiteButton, error)
	GetWebsiteButton(ctx context.Context, id int, with []string) (*models.WebsiteButton, error)
	CreateWebsiteButton(ctx context.Context, req *models.WebsiteButtonCreateRequest) (*models.WebsiteButtonCreateResponse, error)
	UpdateWebsiteButton(ctx context.Context, req *models.WebsiteButtonUpdateRequest) (*models.WebsiteButton, error)
	AddOnlineChat(ctx context.Context, sourceID int) error

	// Chat Templates
	ListChatTemplates(ctx context.Context, filter *filters.TemplatesFilter) ([]*models.ChatTemplate, error)
	DeleteChatTemplate(ctx context.Context, id int) error
	DeleteChatTemplates(ctx context.Context, ids []int) error
	SendChatTemplateOnReview(ctx context.Context, id int) ([]models.ChatTemplateReview, error)
	UpdateChatTemplateReviewStatus(ctx context.Context, templateID, reviewID int, status string) (*models.ChatTemplateReview, error)

	// Short Links
	ListShortLinks(ctx context.Context, filter *filters.ShortLinksFilter) ([]models.ShortLink, error)
	CreateShortLink(ctx context.Context, url string) (models.ShortLink, error)
	CreateShortLinks(ctx context.Context, urls []string) ([]models.ShortLink, error)
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
