package genkit

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/tihn/amo-ai-tgbot-go/config"
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
		plugin := &geminicli.CodeAssist{
			CredsPath: cfg.GeminiCLICredsPath,
		}

		g = genkit.Init(ctx,
			genkit.WithPlugins(plugin),
			genkit.WithPromptDir("./app/gkit/prompts"),
		)

		var err error
		models, err := plugin.DefineAllModels(g)
		if err != nil {
			return nil, err
		}

		// Pick gemini-2.5-flash as default for now
		for _, m := range models {
			if m.Name() == "gemini-cli/gemini-2.5-flash" {
				model = m
				break
			}
		}
		if model == nil && len(models) > 0 {
			model = models[0]
		}
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
