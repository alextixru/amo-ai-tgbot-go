package genkit

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/config"
	"github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/genkit/providers"
	geminicli "github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/genkit/providers/gemini-cli"
)

// Client wraps Genkit instance and model
type Client struct {
	G     *genkit.Genkit
	Model ai.Model
}

// New creates a new Genkit client with configured provider
func New(ctx context.Context, cfg *config.Config) (*Client, error) {
	var model ai.Model
	var g *genkit.Genkit

	switch cfg.AIProvider {
	case "gemini-cli":
		gcli, err := geminicli.New(ctx, cfg.GeminiCLICredsPath)
		if err != nil {
			return nil, err
		}

		g = genkit.Init(ctx,
			genkit.WithPlugins(gcli),
			genkit.WithPromptDir("./app/gkit/prompts"),
		)

		model = gcli.DefineModel(g, "") // Use default model
	default: // ollama
		plugin := providers.OllamaPlugin(cfg)

		g = genkit.Init(ctx,
			genkit.WithPlugins(plugin),
			genkit.WithPromptDir("./app/gkit/prompts"),
		)

		model = providers.InitOllama(g, plugin, cfg)
	}

	return &Client{
		G:     g,
		Model: model,
	}, nil
}
