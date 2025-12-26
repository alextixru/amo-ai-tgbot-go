package providers

import (
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"

	"github.com/tihn/amo-ai-tgbot-go/infrastructure/config"
)

// InitOllama initializes Ollama provider and returns the model
func InitOllama(g *genkit.Genkit, cfg *config.Config) ai.Model {
	plugin := &ollama.Ollama{ServerAddress: cfg.OllamaURL}

	model := plugin.DefineModel(g, ollama.ModelDefinition{
		Name: cfg.OllamaModel,
		Type: "chat",
	}, nil)

	return model
}

// OllamaPlugin returns Ollama plugin for Genkit initialization
func OllamaPlugin(cfg *config.Config) *ollama.Ollama {
	return &ollama.Ollama{ServerAddress: cfg.OllamaURL}
}
