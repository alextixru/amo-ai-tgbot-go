package unsorted

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models"
)

type Service interface {
	ListUnsorted(ctx context.Context, filter *gkitmodels.UnsortedFilter) ([]*models.Unsorted, error)
	GetUnsorted(ctx context.Context, uid string) (*models.Unsorted, error)
	CreateUnsorted(ctx context.Context, category string, items []*models.Unsorted) ([]*models.Unsorted, error)
	AcceptUnsorted(ctx context.Context, uid string, userID int, statusID int) (*models.UnsortedAcceptResult, error)
	DeclineUnsorted(ctx context.Context, uid string, userID int) (*models.UnsortedDeclineResult, error)
	LinkUnsorted(ctx context.Context, uid string, leadID int) (*models.UnsortedLinkResult, error)
	SummaryUnsorted(ctx context.Context, filter *gkitmodels.UnsortedFilter) (*models.UnsortedSummary, error)
}

type service struct {
	sdk *amocrm.SDK
}

func NewService(sdk *amocrm.SDK) Service {
	return &service{sdk: sdk}
}
