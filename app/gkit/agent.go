package gkit

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	genkitClient "github.com/tihn/amo-ai-tgbot-go/infrastructure/genkit"
)

// Agent handles AI processing with Genkit
type Agent struct {
	g     *genkit.Genkit
	model ai.Model
}

// NewAgent creates a new AI agent from Genkit client
func NewAgent(client *genkitClient.Client) *Agent {
	return &Agent{
		g:     client.G,
		model: client.Model,
	}
}

// Process processes a user message and returns AI response
func (a *Agent) Process(ctx context.Context, message string) (string, error) {
	resp, err := genkit.Generate(ctx, a.g,
		ai.WithModel(a.model),
		ai.WithPrompt(message),
	)
	if err != nil {
		return "", err
	}
	return resp.Text(), nil
}
