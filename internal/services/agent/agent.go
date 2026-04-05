// Package agent defines the interface for AI agent implementations.
package agent

import "context"

// Processor is the interface for processing user messages through an AI agent.
// Telegram service depends on this interface, not a concrete implementation.
type Processor interface {
	Process(ctx context.Context, sessionID, message string) (string, error)
}
