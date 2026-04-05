// Package llm provides LLM model adapters for the AI agent.
package llm

import (
	genaiopenai "github.com/achetronic/adk-utils-go/genai/openai"
	"google.golang.org/adk/model"

	"github.com/tihn/amo-ai-tgbot-go/config"
)

// NewProvider creates an ADK-compatible LLM model from application config.
// Currently supports Ollama via OpenAI-compatible API.
func NewProvider(cfg *config.Config) model.LLM {
	return genaiopenai.New(genaiopenai.Config{
		BaseURL:   cfg.OllamaURL + "/v1",
		ModelName: cfg.OllamaModel,
		APIKey:    "ollama", // Ollama doesn't require a key, but the field is mandatory
	})
}
