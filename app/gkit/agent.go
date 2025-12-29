package gkit

import (
	"context"

	amocrm "github.com/alextixru/amocrm-sdk-go"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/tihn/amo-ai-tgbot-go/app/gkit/flows"
	"github.com/tihn/amo-ai-tgbot-go/app/gkit/tools"
	genkitClient "github.com/tihn/amo-ai-tgbot-go/infrastructure/genkit"
)

// Agent handles AI processing with Genkit flows
type Agent struct {
	g        *genkit.Genkit
	model    ai.Model
	chatFlow func(context.Context, flows.ChatInput) (flows.ChatOutput, error)
}

// NewAgent creates a new AI agent with registered flows and tools
func NewAgent(client *genkitClient.Client, sdk *amocrm.SDK) *Agent {
	g := client.G
	model := client.Model

	// Регистрируем все tools (видны в Genkit UI)
	registry := tools.NewRegistry(g, sdk).RegisterAll()

	// Регистрируем Chat Flow с tools
	chatRunner := flows.DefineChatFlow(g, model, registry.AllTools())

	return &Agent{
		g:        g,
		model:    model,
		chatFlow: chatRunner,
	}
}

// Process processes a user message using the chat flow with user context
func (a *Agent) Process(ctx context.Context, message string, userContext map[string]any) (string, error) {
	output, err := a.chatFlow(ctx, flows.ChatInput{
		Message:     message,
		UserContext: userContext,
	})
	if err != nil {
		return "", err
	}
	return output.Response, nil
}
