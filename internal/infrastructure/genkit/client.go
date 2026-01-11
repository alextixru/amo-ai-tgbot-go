package genkit

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/config"
	"github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/genkit/providers"
)

// Client wraps Genkit instance and model
type Client struct {
	G     *genkit.Genkit
	Model ai.Model
}

// New creates a new Genkit client with configured provider
func New(ctx context.Context, cfg *config.Config) (*Client, error) {
	// Get plugin for configured provider
	plugin := providers.OllamaPlugin(cfg)

	// Initialize Genkit with plugin and prompts directory
	g := genkit.Init(ctx,
		genkit.WithPlugins(plugin),
		genkit.WithPromptDir("./app/gkit/prompts"),
	)

	// Initialize model from provider (using the same plugin instance)
	model := providers.InitOllama(g, plugin, cfg)

	return &Client{
		G:     g,
		Model: model,
	}, nil
}
