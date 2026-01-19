package geminicli

import (
	"context"
)

// RecordConversationOfferedRequest reflects the request for recording a conversation offered/started event.
type RecordConversationOfferedRequest struct {
	Model   string `json:"model"`
	Project string `json:"project"`
}

// RecordConversationInteractionRequest reflects the request for recording user interaction with a conversation.
type RecordConversationInteractionRequest struct {
	Model   string `json:"model"`
	Project string `json:"project"`
}

// RecordConversationOffered sends a metric event to the Code Assist API indicating a conversation was offered.
func (c *Client) RecordConversationOffered(ctx context.Context, model string) error {
	req := &RecordConversationOfferedRequest{
		Model:   model,
		Project: c.projectID,
	}
	return c.postWithRetry(ctx, "recordConversationOffered", req, nil)
}

// RecordConversationInteraction sends a metric event to the Code Assist API indicating a user interaction.
func (c *Client) RecordConversationInteraction(ctx context.Context, model string) error {
	req := &RecordConversationInteractionRequest{
		Model:   model,
		Project: c.projectID,
	}
	return c.postWithRetry(ctx, "recordConversationInteraction", req, nil)
}

// recordMetrics is a helper for general metrics (matching recordCodeAssistMetrics in TS)
func (c *Client) recordMetrics(ctx context.Context, metrics any) error {
	// For now we only implement the most important ones to match CLI basic behavior.
	return c.postWithRetry(ctx, "recordCodeAssistMetrics", metrics, nil)
}
