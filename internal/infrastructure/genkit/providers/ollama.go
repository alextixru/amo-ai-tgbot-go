package providers

import (
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"

	"github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/config"
)

// InitOllama initializes Ollama provider and returns the model
func InitOllama(g *genkit.Genkit, plugin *ollama.Ollama, cfg *config.Config) ai.Model {
	model := plugin.DefineModel(g, ollama.ModelDefinition{
		Name: cfg.OllamaModel,
		Type: "chat",
	}, &ai.ModelOptions{
		Supports: &ai.ModelSupports{
			Multiturn:  true,
			SystemRole: true,
			Tools:      true,
		},
	})

	return model
}

// OllamaPlugin returns Ollama plugin for Genkit initialization
func OllamaPlugin(cfg *config.Config) *ollama.Ollama {
	return &ollama.Ollama{ServerAddress: cfg.OllamaURL}
}
