package geminicli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core/api"
	"github.com/firebase/genkit/go/genkit"
	"google.golang.org/genai"
)

const (
	ProviderName = "gemini-cli"
)

// CodeAssist is a Genkit plugin for interacting with the Google Code Assist API.
type CodeAssist struct {
	CredsPath string

	client  *Client
	mu      sync.Mutex
	initted bool
}

// Name returns the name of the plugin.
func (ca *CodeAssist) Name() string {
	return ProviderName
}

// Init initializes the plugin.
func (ca *CodeAssist) Init(ctx context.Context) []api.Action {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	if ca.initted {
		panic("plugin already initialized")
	}

	client, err := NewClient(ctx, ClientOptions{
		CredsPath: ca.CredsPath,
		NoBrowser: os.Getenv("NO_BROWSER") == "true",
	})
	if err != nil {
		panic(fmt.Errorf("CodeAssist.Init: %w", err))
	}
	ca.client = client
	ca.initted = true

	return []api.Action{}
}

// DefineModel defines a model with the given name and registers it in Genkit.
func (ca *CodeAssist) DefineModel(gk *genkit.Genkit, name string, opts *ai.ModelOptions) (ai.Model, error) {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	if !ca.initted {
		return nil, errors.New("CodeAssist plugin not initialized")
	}

	if opts == nil {
		modelOpts, ok := SupportedModels[name]
		if !ok {
			return nil, fmt.Errorf("unknown model %q", name)
		}
		opts = &modelOpts
	}

	return ca.defineModelInternal(gk, name, opts), nil
}

// defineModelInternal creates and registers a model in Genkit registry.
func (ca *CodeAssist) defineModelInternal(gk *genkit.Genkit, name string, opts *ai.ModelOptions) ai.Model {
	client := ca.client
	return genkit.DefineModel(gk, api.NewName(ProviderName, name), opts, func(ctx context.Context, req *ai.ModelRequest, cb ai.ModelStreamCallback) (*ai.ModelResponse, error) {
		// Convert Genkit request to SDK types
		gcc, err := toGeminiRequest(req)
		if err != nil {
			return nil, err
		}

		// Enable thoughts and configuration based on model version to match CLI behavior.
		if gcc.ThinkingConfig == nil {
			if name == Gemini3ProPreview || name == Gemini3FlashPreview {
				gcc.ThinkingConfig = &genai.ThinkingConfig{
					IncludeThoughts: true,
					ThinkingLevel:   genai.ThinkingLevelHigh,
				}
			} else if strings.HasPrefix(name, "gemini-2.5-") {
				gcc.ThinkingConfig = &genai.ThinkingConfig{
					IncludeThoughts: true,
					ThinkingBudget:  genai.Ptr(int32(8192)),
				}
			}
		}

		// Convert messages to genai contents
		var genaiContents []*genai.Content
		for _, m := range req.Messages {
			if m.Role == ai.RoleSystem {
				continue
			}
			parts, err := toGeminiParts(m.Content)
			if err != nil {
				return nil, err
			}
			genaiContents = append(genaiContents, &genai.Content{
				Role:  string(m.Role),
				Parts: parts,
			})
		}

		// Non-streaming mode
		if cb == nil {
			resp, err := client.Generate(ctx, name, genaiContents, gcc)
			if err != nil {
				return nil, err
			}
			// Record telemetry (best effort)
			_ = client.RecordConversationOffered(ctx, name)
			r, err := translateResponse(resp)
			if err != nil {
				return nil, err
			}
			r.Request = req
			return r, nil
		}

		// Streaming mode
		respChan, errChan := client.GenerateStream(ctx, name, genaiContents, gcc)

		var genaiParts []*genai.Part
		var chunks []*ai.Part
		var lastResp *genai.GenerateContentResponse

		for {
			select {
			case resp, ok := <-respChan:
				if !ok {
					// Channel closed, build final response
					if lastResp == nil {
						return nil, fmt.Errorf("no response received from stream")
					}

					// Create merged candidate with accumulated parts (matching official plugin)
					finishReason := genai.FinishReasonUnspecified
					if len(lastResp.Candidates) > 0 {
						finishReason = lastResp.Candidates[0].FinishReason
					}

					merged := &genai.GenerateContentResponse{
						Candidates: []*genai.Candidate{
							{
								FinishReason: finishReason,
								Content: &genai.Content{
									Role:  string(ai.RoleModel),
									Parts: genaiParts,
								},
							},
						},
						UsageMetadata:  lastResp.UsageMetadata,
						PromptFeedback: lastResp.PromptFeedback,
					}

					// Record telemetry (best effort)
					_ = client.RecordConversationOffered(ctx, name)

					r, err := translateResponse(merged)
					if err != nil {
						return nil, err
					}
					r.Message.Content = chunks
					r.Request = req
					return r, nil
				}

				lastResp = resp

				// Translate chunk and call callback
				for _, c := range resp.Candidates {
					tc, err := translateCandidate(c)
					if err != nil {
						return nil, err
					}

					err = cb(ctx, &ai.ModelResponseChunk{
						Content: tc.Message.Content,
						Role:    ai.RoleModel,
					})
					if err != nil {
						return nil, err
					}

					if c.Content != nil && len(c.Content.Parts) > 0 {
						genaiParts = append(genaiParts, c.Content.Parts...)
					}
					chunks = append(chunks, tc.Message.Content...)
				}

			case err := <-errChan:
				if err != nil {
					return nil, err
				}

			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	})
}

// DefineAllModels defines all supported models.
func (ca *CodeAssist) DefineAllModels(gk *genkit.Genkit) ([]ai.Model, error) {
	var models []ai.Model
	for name, opts := range SupportedModels {
		m, err := ca.DefineModel(gk, name, &opts)
		if err != nil {
			return nil, err
		}
		models = append(models, m)
	}
	return models, nil
}

// GetQuota возвращает текущие квоты пользователя
func (ca *CodeAssist) GetQuota(ctx context.Context) (*RetrieveUserQuotaResponse, error) {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	if !ca.initted {
		return nil, errors.New("CodeAssist plugin not initialized")
	}
	return ca.client.RetrieveUserQuota(ctx)
}
