package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) convertCall(c *models.Call) *CallOutput {
	if c == nil {
		return nil
	}
	return &CallOutput{
		ID:                  c.ID,
		Direction:           string(c.Direction),
		Duration:            c.Duration,
		Phone:               c.Phone,
		CallResult:          c.CallResult,
		CallStatus:          int(c.CallStatus),
		Source:              c.Source,
		UniqueID:            c.Uniq,
		RecordURL:           c.Link,
		ResponsibleUserName: s.resolveUserID(c.ResponsibleUserID),
		CreatedByName:       s.resolveUserID(c.CreatedBy),
		CreatedAt:           toISO(c.CreatedAt),
	}
}

func (s *service) CreateCall(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.CallData) (*CallOutput, error) {
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
	c, err := s.sdk.Calls().CreateOne(ctx, &call)
	if err != nil {
		return nil, err
	}
	return s.convertCall(c), nil
}
