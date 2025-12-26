package gkit

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/tihn/amo-ai-tgbot-go/app/gkit/flows"
	genkitClient "github.com/tihn/amo-ai-tgbot-go/infrastructure/genkit"
)

// Agent handles AI processing with Genkit flows
type Agent struct {
	g        *genkit.Genkit
	model    ai.Model
	chatFlow func(context.Context, flows.ChatInput) (flows.ChatOutput, error)
}

// NewAgent creates a new AI agent with registered flows
func NewAgent(client *genkitClient.Client) *Agent {
	g := client.G
	model := client.Model

	// Регистрируем Chat Flow (виден в Genkit UI)
	chatRunner := flows.DefineChatFlow(g, model)

	return &Agent{
		g:        g,
		model:    model,
		chatFlow: chatRunner,
	}
}

// Process processes a user message using the chat flow
func (a *Agent) Process(ctx context.Context, message string) (string, error) {
	output, err := a.chatFlow(ctx, flows.ChatInput{Message: message})
	if err != nil {
		return "", err
	}
	return output.Response, nil
}
