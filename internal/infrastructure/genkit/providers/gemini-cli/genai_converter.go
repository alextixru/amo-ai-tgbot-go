package geminicli

import (
	"encoding/json"

	"github.com/firebase/genkit/go/ai"
	"google.golang.org/genai"
)

// ToGenaiContents converts Genkit messages to Genai contents
func ToGenaiContents(msgs []*ai.Message) []*genai.Content {
	var contents []*genai.Content
	for _, m := range msgs {
		if m.Role == ai.RoleSystem {
			continue // System instructions handled separately in config or by SDK if needed
		}
		contents = append(contents, &genai.Content{
			Role:  string(m.Role),
			Parts: ToGenaiParts(m.Content),
		})
	}
	return contents
}

// ToGenaiSystemInstruction extracts system instruction from messages
func ToGenaiSystemInstruction(msgs []*ai.Message) *genai.Content {
	for _, m := range msgs {
		if m.Role == ai.RoleSystem {
			return &genai.Content{
				Role:  "system",
				Parts: ToGenaiParts(m.Content),
			}
		}
	}
	return nil
}

// ToGenaiParts converts Genkit parts to Genai parts
func ToGenaiParts(parts []*ai.Part) []*genai.Part {
	var res []*genai.Part
	for _, p := range parts {
		if p.IsText() {
			res = append(res, &genai.Part{Text: p.Text})
		} else if p.IsMedia() {
			res = append(res, &genai.Part{
				InlineData: &genai.Blob{
					MIMEType: p.ContentType,
					Data:     []byte(p.Text), // Genkit stores data as string in Text for blobs
				},
			})
		} else if p.IsToolRequest() {
			args, _ := p.ToolRequest.Input.(map[string]any)
			res = append(res, &genai.Part{
				FunctionCall: &genai.FunctionCall{
					Name: p.ToolRequest.Name,
					Args: args,
				},
			})
		} else if p.IsToolResponse() {
			var resp map[string]any
			if m, ok := p.ToolResponse.Output.(map[string]any); ok {
				resp = m
			} else {
				// If output is not a map (e.g. an array), wrap it
				resp = map[string]any{"output": p.ToolResponse.Output}
			}
			res = append(res, &genai.Part{
				FunctionResponse: &genai.FunctionResponse{
					Name:     p.ToolResponse.Name,
					Response: resp,
				},
			})
		}
	}
	return res
}

// ToGenaiConfig converts Genkit config to Genai config
func ToGenaiConfig(reqConfig any, tools []*ai.ToolDefinition) *genai.GenerateContentConfig {
	config := &genai.GenerateContentConfig{}

	if reqConfig != nil {
		// Genkit uses map[string]any or a struct for config
		// We try to extract common fields
		data, _ := json.Marshal(reqConfig)
		var m map[string]any
		json.Unmarshal(data, &m)

		if v, ok := m["temperature"].(float64); ok {
			f := float32(v)
			config.Temperature = &f
		}
		if v, ok := m["topP"].(float64); ok {
			f := float32(v)
			config.TopP = &f
		}
		if v, ok := m["topK"].(float64); ok {
			f := float32(v)
			config.TopK = &f
		}
		if v, ok := m["maxOutputTokens"].(float64); ok {
			config.MaxOutputTokens = int32(v)
		}
		if v, ok := m["stopSequences"].([]any); ok {
			var stops []string
			for _, s := range v {
				if str, ok := s.(string); ok {
					stops = append(stops, str)
				}
			}
			config.StopSequences = stops
		}

		// Support for thinkingConfig if passed in ExtraOptions or similar
		if tc, ok := m["thinkingConfig"].(map[string]any); ok {
			config.ThinkingConfig = &genai.ThinkingConfig{}
			if budget, ok := tc["thinkingBudget"].(float64); ok {
				b := int32(budget)
				config.ThinkingConfig.ThinkingBudget = &b
			}
		}
	}

	if len(tools) > 0 {
		var genaiTools []*genai.Tool
		for _, t := range tools {
			genaiTools = append(genaiTools, &genai.Tool{
				FunctionDeclarations: []*genai.FunctionDeclaration{
					{
						Name:        t.Name,
						Description: t.Description,
						Parameters:  convertToGenaiSchema(t.InputSchema),
					},
				},
			})
		}
		config.Tools = genaiTools
	}

	return config
}

func convertToGenaiSchema(schema map[string]any) *genai.Schema {
	if schema == nil {
		return nil
	}
	data, _ := json.Marshal(schema)
	var s genai.Schema
	json.Unmarshal(data, &s)
	return &s
}

// FromGenaiResponse converts Genai response back to Genkit response
func FromGenaiResponse(resp *genai.GenerateContentResponse, req *ai.ModelRequest) *ai.ModelResponse {
	if resp == nil || len(resp.Candidates) == 0 {
		return &ai.ModelResponse{
			Message: &ai.Message{
				Role:    ai.RoleModel,
				Content: []*ai.Part{},
			},
			Request: req,
		}
	}

	cand := resp.Candidates[0]

	// Check if Content is nil
	if cand.Content == nil {
		return &ai.ModelResponse{
			Message: &ai.Message{
				Role:    ai.RoleModel,
				Content: []*ai.Part{},
			},
			Request: req,
		}
	}

	var usage *ai.GenerationUsage
	if resp.UsageMetadata != nil {
		usage = &ai.GenerationUsage{
			InputTokens:  int(resp.UsageMetadata.PromptTokenCount),
			OutputTokens: int(resp.UsageMetadata.CandidatesTokenCount),
			TotalTokens:  int(resp.UsageMetadata.TotalTokenCount),
		}
	}

	return &ai.ModelResponse{
		Message: &ai.Message{
			Role:    ai.Role(cand.Content.Role),
			Content: FromGenaiParts(cand.Content.Parts),
		},
		Usage:   usage,
		Request: req,
	}
}

// FromGenaiParts converts Genai parts to Genkit parts
func FromGenaiParts(parts []*genai.Part) []*ai.Part {
	var res []*ai.Part
	for _, p := range parts {
		if p.Text != "" {
			res = append(res, ai.NewTextPart(p.Text))
		} else if p.InlineData != nil {
			res = append(res, ai.NewMediaPart(p.InlineData.MIMEType, string(p.InlineData.Data)))
		} else if p.FunctionCall != nil {
			res = append(res, ai.NewToolRequestPart(&ai.ToolRequest{
				Name:  p.FunctionCall.Name,
				Input: p.FunctionCall.Args,
			}))
		} else if p.FunctionResponse != nil {
			res = append(res, ai.NewToolResponsePart(&ai.ToolResponse{
				Name:   p.FunctionResponse.Name,
				Output: p.FunctionResponse.Response,
			}))
		}
	}
	return res
}

// SyntheticThoughtSignature is used for Gemini 3 Pro models that require thought signatures
const SyntheticThoughtSignature = "skip_thought_signature_validator"

// EnsureThoughtSignatures adds synthetic thought signatures to function calls
// in the active loop (since last user text message) for Gemini 3 Pro compatibility.
func EnsureThoughtSignatures(contents []*genai.Content) []*genai.Content {
	// Find the start of the active loop - last user turn with text (not function response)
	activeLoopStartIndex := -1
	for i := len(contents) - 1; i >= 0; i-- {
		content := contents[i]
		if content.Role == "user" {
			// Check if this user message has text (not just function response)
			for _, part := range content.Parts {
				if part.Text != "" {
					activeLoopStartIndex = i
					break
				}
			}
			if activeLoopStartIndex != -1 {
				break
			}
		}
	}

	if activeLoopStartIndex == -1 {
		return contents
	}

	// Shallow copy the array to avoid modifying original
	newContents := make([]*genai.Content, len(contents))
	copy(newContents, contents)

	// Iterate through model messages in the active loop
	for i := activeLoopStartIndex; i < len(newContents); i++ {
		content := newContents[i]
		if content.Role != "model" || len(content.Parts) == 0 {
			continue
		}

		// Find the first function call in this model turn
		for j, part := range content.Parts {
			if part.FunctionCall != nil {
				// Check if it already has a thought signature
				if len(part.ThoughtSignature) == 0 {
					// Create new parts slice and replace the part with one that has signature
					newParts := make([]*genai.Part, len(content.Parts))
					copy(newParts, content.Parts)
					newParts[j] = &genai.Part{
						FunctionCall:     part.FunctionCall,
						ThoughtSignature: []byte(SyntheticThoughtSignature),
					}
					newContents[i] = &genai.Content{
						Role:  content.Role,
						Parts: newParts,
					}
				}
				break // Only consider the first function call per model turn
			}
		}
	}

	return newContents
}
