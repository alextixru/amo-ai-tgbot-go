package activities

import (
	"context"
)

func (s *service) CloseTalk(ctx context.Context, talkID string, forceClose bool) error {
	// TODO: Find correct model for TalkCloseOptions if force_close is needed.
	// AUDIT.md mentions TalkCloseOptions exists in SDK, but it's not found in current models package.
	return s.sdk.Talks().Close(ctx, talkID, nil)
}
