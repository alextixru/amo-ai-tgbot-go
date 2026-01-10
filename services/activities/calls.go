package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

// ============ CALLS ============

func (s *service) CreateCall(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.CallData) (*models.Call, error) {
	call := models.Call{
		EntityID:   parent.ID,
		EntityType: parent.Type,
		Duration:   data.Duration,
		Source:     data.Source,
		Phone:      data.Phone,
		CallResult: data.CallResult,
		CallStatus: models.CallStatus(data.CallStatus),
	}
	if data.Direction != "" {
		call.Direction = models.CallDirection(data.Direction)
	}
	if data.UniqueID != "" {
		call.Uniq = data.UniqueID
	}
	if data.RecordURL != "" {
		call.Link = data.RecordURL
	}
	return s.sdk.Calls().CreateOne(ctx, &call)
}
