package geminicli

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core/api"
	"github.com/firebase/genkit/go/genkit"
)

const (
	ProviderName = "gemini-cli"
	DefaultModel = "gemini-2.5-flash"
)

type GeminiCLI struct {
	client *Client
}

func New(ctx context.Context, credsPath string) (*GeminiCLI, error) {
	client, err := NewClient(ctx, credsPath)
	if err != nil {
		return nil, err
	}
	return &GeminiCLI{client: client}, nil
}

func (g *GeminiCLI) Name() string {
	return ProviderName
}

func (g *GeminiCLI) Init(ctx context.Context) []api.Action {
	return nil
}

func (g *GeminiCLI) DefineModel(gk *genkit.Genkit, modelName string) ai.Model {
	if modelName == "" {
		modelName = DefaultModel
	}
	// Code Assist API expects model name without "models/" prefix

	return genkit.DefineModel(gk, modelName, &ai.ModelOptions{
		Label: "Gemini CLI (Code Assist)",
		Supports: &ai.ModelSupports{
			Multiturn:  true,
			SystemRole: true,
			Tools:      true,
			Media:      true,
		},
	}, func(ctx context.Context, req *ai.ModelRequest, cb ai.ModelStreamCallback) (*ai.ModelResponse, error) {
		// Convert Genkit request to SDK types
		contents := ToGenaiContents(req.Messages)
		config := ToGenaiConfig(req.Config, req.Tools)

		// Handle system instruction
		config.SystemInstruction = ToGenaiSystemInstruction(req.Messages)

		// Call the refactored Generate method
		resp, err := g.client.Generate(ctx, modelName, contents, config)
		if err != nil {
			return nil, err
		}

		// Check if response is valid before converting
		if resp == nil || len(resp.Candidates) == 0 {
			return nil, fmt.Errorf("empty response from API")
		}

		// Convert SDK response back to Genkit
		return FromGenaiResponse(resp, req), nil
	})
}
