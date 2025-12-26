package ai

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"

	"github.com/tihn/amo-ai-tgbot-go/infrastructure/config"
)

// Agent handles AI processing with Genkit
type Agent struct {
	g     *genkit.Genkit
	model ai.Model
}

// New creates a new AI agent with Ollama backend
func New(ctx context.Context, cfg *config.Config) (*Agent, error) {
	ollamaPlugin := &ollama.Ollama{ServerAddress: cfg.OllamaURL}

	g := genkit.Init(ctx, genkit.WithPlugins(ollamaPlugin))

	// Define the model (required for Ollama plugin)
	model := ollamaPlugin.DefineModel(g, ollama.ModelDefinition{
		Name: cfg.OllamaModel,
		Type: "chat",
	}, nil)

	return &Agent{
		g:     g,
		model: model,
	}, nil
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
