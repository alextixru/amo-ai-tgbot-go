// Copyright 2025 Google LLC
// SPDX-License-Identifier: Apache-2.0

package geminicli

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"google.golang.org/genai"
)

const (
	// Tool name regex
	toolNameRegex = "^[a-zA-Z_][a-zA-Z0-9_.-]{0,63}$"
)

// configFromRequest converts any supported config type to [genai.GenerateContentConfig].
func configFromRequest(input *ai.ModelRequest) (*genai.GenerateContentConfig, error) {
	var result genai.GenerateContentConfig

	switch config := input.Config.(type) {
	case genai.GenerateContentConfig:
		result = config
	case *genai.GenerateContentConfig:
		result = *config
	case map[string]any:
		data, err := json.Marshal(config)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}
	case nil:
		// Empty but valid config
	default:
		return nil, fmt.Errorf("unexpected config type: %T", input.Config)
	}

	return &result, nil
}

// toGeminiRequest translates an [*ai.ModelRequest] to
// *genai.GenerateContentConfig
func toGeminiRequest(input *ai.ModelRequest) (*genai.GenerateContentConfig, error) {
	gcc, err := configFromRequest(input)
	if err != nil {
		return nil, err
	}

	if gcc.CandidateCount == 0 {
		gcc.CandidateCount = 1
	}

	// Genkit primitive fields must be used instead of go-genai fields
	// i.e.: system prompt, tools, cached content, response schema, etc
	if gcc.CandidateCount != 1 {
		return nil, errors.New("multiple candidates is not supported")
	}
	if gcc.SystemInstruction != nil {
		return nil, errors.New("system instruction must be set using Genkit feature: ai.WithSystemPrompt()")
	}
	if gcc.CachedContent != "" {
		return nil, errors.New("cached content must be set using Genkit feature: ai.WithCacheTTL()")
	}
	if gcc.ResponseSchema != nil {
		return nil, errors.New("response schema must be set using Genkit feature: ai.WithTools() or ai.WithOutputType()")
	}
	if gcc.ResponseMIMEType != "" {
		return nil, errors.New("response MIME type must be set using Genkit feature: ai.WithOutputType(), ai.WithOutputSchema(), ai.WithOutputSchemaByName()")
	}
	if gcc.ResponseJsonSchema != nil {
		return nil, errors.New("response JSON schema must be set using Genkit feature: ai.WithOutputSchema()")
	}

	// Set response MIME type and schema based on output format.
	hasOutput := input.Output != nil
	if hasOutput && len(input.Tools) == 0 {
		switch {
		case input.Output.ContentType == "application/json" || input.Output.Format == "json":
			gcc.ResponseMIMEType = "application/json"
		case input.Output.ContentType == "text/enum" || input.Output.Format == "enum":
			gcc.ResponseMIMEType = "text/x.enum"
		}
	}

	if input.Output != nil && input.Output.Constrained && gcc.ResponseMIMEType != "" {
		schema, err := toGeminiSchema(input.Output.Schema, input.Output.Schema)
		if err != nil {
			return nil, err
		}
		gcc.ResponseSchema = schema
	}

	if len(input.Tools) > 0 {
		tools, err := toGeminiTools(input.Tools)
		if err != nil {
			return nil, err
		}
		gcc.Tools = mergeTools(append(gcc.Tools, tools...))

		tc, err := toGeminiToolChoice(input.ToolChoice, input.Tools)
		if err != nil {
			return nil, err
		}
		gcc.ToolConfig = tc
	}

	var systemParts []*genai.Part
	for _, m := range input.Messages {
		if m.Role == ai.RoleSystem {
			parts, err := toGeminiParts(m.Content)
			if err != nil {
				return nil, err
			}
			systemParts = append(systemParts, parts...)
		}
	}

	if len(systemParts) > 0 {
		gcc.SystemInstruction = &genai.Content{
			Parts: systemParts,
			Role:  string(ai.RoleSystem),
		}
	}

	return gcc, nil
}

// toGeminiTools translates a slice of [ai.ToolDefinition] to a slice of [genai.Tool].
func toGeminiTools(inTools []*ai.ToolDefinition) ([]*genai.Tool, error) {
	var outTools []*genai.Tool
	functions := []*genai.FunctionDeclaration{}

	for _, t := range inTools {
		if !validToolName(t.Name) {
			return nil, fmt.Errorf("invalid tool name: %q", t.Name)
		}
		inputSchema, err := toGeminiSchema(t.InputSchema, t.InputSchema)
		if err != nil {
			return nil, err
		}
		fd := &genai.FunctionDeclaration{
			Name:        t.Name,
			Parameters:  inputSchema,
			Description: t.Description,
		}
		functions = append(functions, fd)
	}

	if len(functions) > 0 {
		outTools = append(outTools, &genai.Tool{
			FunctionDeclarations: functions,
		})
	}

	return outTools, nil
}

// toGeminiParts converts a slice of [ai.Part] to a slice of [genai.Part].
func toGeminiParts(parts []*ai.Part) ([]*genai.Part, error) {
	res := make([]*genai.Part, 0, len(parts))
	for _, p := range parts {
		part, err := toGeminiPart(p)
		if err != nil {
			return nil, err
		}
		res = append(res, part)
	}
	return res, nil
}

// toGeminiPart converts a [ai.Part] to a [genai.Part].
func toGeminiPart(p *ai.Part) (*genai.Part, error) {
	var gp *genai.Part
	switch {
	case p.IsReasoning():
		gp = genai.NewPartFromText(p.Text)
		gp.Thought = true
	case p.IsText():
		gp = genai.NewPartFromText(p.Text)
	case p.IsMedia():
		if strings.HasPrefix(p.Text, "data:") {
			idx := strings.Index(p.Text, ",")
			if idx == -1 {
				return nil, errors.New("invalid data URI")
			}
			data, err := base64.StdEncoding.DecodeString(p.Text[idx+1:])
			if err != nil {
				return nil, err
			}
			gp = genai.NewPartFromBytes(data, p.ContentType)
		} else {
			gp = genai.NewPartFromURI(p.Text, p.ContentType)
		}
	case p.IsData():
		// Handle data URIs for binary content
		idx := strings.Index(p.Text, ",")
		if idx == -1 {
			return nil, errors.New("invalid data URI")
		}
		data, err := base64.StdEncoding.DecodeString(p.Text[idx+1:])
		if err != nil {
			return nil, err
		}
		gp = genai.NewPartFromBytes(data, p.ContentType)
	case p.IsToolResponse():
		toolResp := p.ToolResponse
		var output map[string]any
		if m, ok := toolResp.Output.(map[string]any); ok {
			output = m
		} else {
			output = map[string]any{
				"output": toolResp.Output,
			}
		}
		gp = genai.NewPartFromFunctionResponse(toolResp.Name, output)
	case p.IsToolRequest():
		toolReq := p.ToolRequest
		var input map[string]any
		if m, ok := toolReq.Input.(map[string]any); ok {
			input = m
		} else {
			input = map[string]any{
				"input": toolReq.Input,
			}
		}
		gp = genai.NewPartFromFunctionCall(toolReq.Name, input)
	default:
		return nil, fmt.Errorf("unknown part type: %d", p.Kind)
	}

	if p.Metadata != nil {
		if sig, ok := p.Metadata["signature"].([]byte); ok {
			gp.ThoughtSignature = sig
		}
	}

	return gp, nil
}

// mergeTools consolidates all FunctionDeclarations into a single Tool
func mergeTools(ts []*genai.Tool) []*genai.Tool {
	var decls []*genai.FunctionDeclaration
	var out []*genai.Tool

	for _, t := range ts {
		if t == nil {
			continue
		}
		if len(t.FunctionDeclarations) == 0 {
			out = append(out, t)
			continue
		}
		decls = append(decls, t.FunctionDeclarations...)
	}

	if len(decls) > 0 {
		out = append([]*genai.Tool{{FunctionDeclarations: decls}}, out...)
	}
	return out
}

// toGeminiSchema translates a map representing a standard JSON schema to a more
// limited [genai.Schema].
func toGeminiSchema(originalSchema map[string]any, genkitSchema map[string]any) (*genai.Schema, error) {
	if len(genkitSchema) == 0 {
		return nil, nil
	}
	if v, ok := genkitSchema["$ref"]; ok {
		ref, _ := v.(string)
		s, err := resolveRef(originalSchema, ref)
		if err != nil {
			return nil, err
		}
		return toGeminiSchema(originalSchema, s)
	}

	// Handle "anyOf" subschemas by finding the first valid schema definition
	if v, ok := genkitSchema["anyOf"]; ok {
		if anyOfList, isList := v.([]any); isList {
			for _, item := range anyOfList {
				subSchema, ok := item.(map[string]any)
				if !ok {
					continue
				}
				if subSchemaType, hasType := subSchema["type"]; hasType {
					if typeStr, isString := subSchemaType.(string); isString && typeStr != "null" {
						// Copy title and description from parent
						if title, ok := genkitSchema["title"]; ok {
							subSchema["title"] = title
						}
						if description, ok := genkitSchema["description"]; ok {
							subSchema["description"] = description
						}
						return toGeminiSchema(originalSchema, subSchema)
					}
				}
			}
		}
	}

	schema := &genai.Schema{}
	typeVal, ok := genkitSchema["type"]
	if !ok {
		// No type field - this can happen with boolean schemas or unresolved anyOf
		// Return a permissive schema that accepts any value
		schema.Type = genai.TypeString
		schema.Description = "Any value (flexible type)"
		return schema, nil
	}

	typeStr, _ := typeVal.(string)
	switch typeStr {
	case "string":
		schema.Type = genai.TypeString
	case "float64", "number":
		schema.Type = genai.TypeNumber
	case "integer":
		schema.Type = genai.TypeInteger
	case "boolean":
		schema.Type = genai.TypeBoolean
	case "object":
		schema.Type = genai.TypeObject
	case "array":
		schema.Type = genai.TypeArray
	default:
		return nil, fmt.Errorf("schema type %q not allowed", typeStr)
	}

	if v, ok := genkitSchema["required"]; ok {
		schema.Required = castToStringArray(v)
	}
	if v, ok := genkitSchema["description"]; ok {
		schema.Description = v.(string)
	}
	if v, ok := genkitSchema["format"]; ok {
		schema.Format = v.(string)
	}
	if v, ok := genkitSchema["title"]; ok {
		schema.Title = v.(string)
	}
	if v, ok := genkitSchema["propertyOrdering"]; ok {
		schema.PropertyOrdering = castToStringArray(v)
	}
	if v, ok := genkitSchema["minItems"]; ok {
		if i64, ok := castToInt64(v); ok {
			schema.MinItems = genai.Ptr(i64)
		}
	}
	if v, ok := genkitSchema["maxItems"]; ok {
		if i64, ok := castToInt64(v); ok {
			schema.MaxItems = genai.Ptr(i64)
		}
	}
	if v, ok := genkitSchema["minimum"]; ok {
		if f64, ok := castToFloat64(v); ok {
			schema.Minimum = genai.Ptr(f64)
		}
	}
	if v, ok := genkitSchema["maximum"]; ok {
		if f64, ok := castToFloat64(v); ok {
			schema.Maximum = genai.Ptr(f64)
		}
	}
	if v, ok := genkitSchema["enum"]; ok {
		schema.Enum = castToStringArray(v)
	}
	if v, ok := genkitSchema["items"]; ok {
		switch itemsVal := v.(type) {
		case map[string]any:
			items, err := toGeminiSchema(originalSchema, itemsVal)
			if err != nil {
				return nil, err
			}
			schema.Items = items
		case bool:
			// JSON Schema boolean: true = any type allowed, false = no items
			// Skip setting Items - Gemini will accept any array elements
		}
	}
	if val, ok := genkitSchema["properties"]; ok {
		props := map[string]*genai.Schema{}
		for k, v := range val.(map[string]any) {
			switch propVal := v.(type) {
			case map[string]any:
				p, err := toGeminiSchema(originalSchema, propVal)
				if err != nil {
					return nil, err
				}
				props[k] = p
			case bool:
				// JSON Schema boolean property: skip (any value accepted)
				continue
			}
		}
		schema.Properties = props
	}

	return schema, nil
}

func resolveRef(originalSchema map[string]any, ref string) (map[string]any, error) {
	tkns := strings.Split(ref, "/")
	name := tkns[len(tkns)-1]
	if defs, ok := originalSchema["$defs"].(map[string]any); ok {
		if def, ok := defs[name].(map[string]any); ok {
			return def, nil
		}
	}
	return nil, fmt.Errorf("unable to resolve schema reference")
}

func castToStringArray(v any) []string {
	switch a := v.(type) {
	case []string:
		return a
	case []any:
		var out []string
		for _, it := range a {
			if s, ok := it.(string); ok {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

func toGeminiToolChoice(toolChoice ai.ToolChoice, tools []*ai.ToolDefinition) (*genai.ToolConfig, error) {
	var mode genai.FunctionCallingConfigMode
	switch toolChoice {
	case "":
		return nil, nil
	case ai.ToolChoiceAuto:
		mode = genai.FunctionCallingConfigModeAuto
	case ai.ToolChoiceRequired:
		mode = genai.FunctionCallingConfigModeAny
	case ai.ToolChoiceNone:
		mode = genai.FunctionCallingConfigModeNone
	default:
		return nil, fmt.Errorf("tool choice mode %q not supported", toolChoice)
	}

	var toolNames []string
	if mode == genai.FunctionCallingConfigModeAny {
		for _, t := range tools {
			toolNames = append(toolNames, t.Name)
		}
	}
	return &genai.ToolConfig{
		FunctionCallingConfig: &genai.FunctionCallingConfig{
			Mode:                 mode,
			AllowedFunctionNames: toolNames,
		},
	}, nil
}

// castToInt64 converts v to int64 when possible.
func castToInt64(v any) (int64, bool) {
	switch t := v.(type) {
	case int:
		return int64(t), true
	case int64:
		return t, true
	case float64:
		return int64(t), true
	}
	return 0, false
}

// castToFloat64 converts v to float64 when possible.
func castToFloat64(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	}
	return 0, false
}

// translateCandidate translates from a genai.Candidate to an ai.ModelResponse.
func translateCandidate(cand *genai.Candidate) (*ai.ModelResponse, error) {
	m := &ai.ModelResponse{}
	switch cand.FinishReason {
	case genai.FinishReasonStop:
		m.FinishReason = ai.FinishReasonStop
	case genai.FinishReasonMaxTokens:
		m.FinishReason = ai.FinishReasonLength
	case genai.FinishReasonSafety,
		genai.FinishReasonRecitation,
		genai.FinishReasonLanguage,
		genai.FinishReasonBlocklist,
		genai.FinishReasonProhibitedContent,
		genai.FinishReasonSPII,
		genai.FinishReasonImageSafety,
		genai.FinishReasonImageProhibitedContent,
		genai.FinishReasonImageRecitation:
		m.FinishReason = ai.FinishReasonBlocked
	case genai.FinishReasonMalformedFunctionCall,
		genai.FinishReasonUnexpectedToolCall,
		genai.FinishReasonNoImage,
		genai.FinishReasonImageOther,
		genai.FinishReasonOther:
		m.FinishReason = ai.FinishReasonOther
	case "MISSING_THOUGHT_SIGNATURE":
		// Gemini 3 returns this when thought signatures are missing from the request.
		// The SDK may not have this constant yet, so we match on the string value.
		m.FinishReason = ai.FinishReasonOther
	default:
		if cand.FinishReason != "" && cand.FinishReason != genai.FinishReasonUnspecified {
			m.FinishReason = ai.FinishReasonUnknown
		}
	}

	m.FinishMessage = cand.FinishMessage
	if cand.Content == nil {
		// Return empty message with model role for streaming compatibility
		m.Message = &ai.Message{Role: ai.RoleModel}
		return m, nil
	}
	msg := &ai.Message{Role: ai.Role(cand.Content.Role)}
	for _, part := range cand.Content.Parts {
		var p *ai.Part
		if part.Thought {
			p = ai.NewReasoningPart(part.Text, part.ThoughtSignature)
		} else if part.Text != "" {
			p = ai.NewTextPart(part.Text)
		} else if part.InlineData != nil {
			p = ai.NewMediaPart(part.InlineData.MIMEType, "data:"+part.InlineData.MIMEType+";base64,"+base64.StdEncoding.EncodeToString(part.InlineData.Data))
		} else if part.FileData != nil {
			p = ai.NewMediaPart(part.FileData.MIMEType, part.FileData.FileURI)
		} else if part.FunctionCall != nil {
			p = ai.NewToolRequestPart(&ai.ToolRequest{
				Name:  part.FunctionCall.Name,
				Input: part.FunctionCall.Args,
			})
		} else if part.ExecutableCode != nil {
			p = ai.NewCustomPart(map[string]any{
				"executableCode": map[string]any{
					"code":     part.ExecutableCode.Code,
					"language": part.ExecutableCode.Language,
				},
			})
		} else if part.CodeExecutionResult != nil {
			p = ai.NewCustomPart(map[string]any{
				"codeExecutionResult": map[string]any{
					"output":  part.CodeExecutionResult.Output,
					"outcome": part.CodeExecutionResult.Outcome,
				},
			})
		}
		if p != nil {
			if len(part.ThoughtSignature) > 0 {
				if p.Metadata == nil {
					p.Metadata = make(map[string]any)
				}
				p.Metadata["signature"] = part.ThoughtSignature
			}
			msg.Content = append(msg.Content, p)
		}
	}
	m.Message = msg
	return m, nil
}

// translateResponse translates from a genai.GenerateContentResponse to a ai.ModelResponse.
func translateResponse(resp *genai.GenerateContentResponse) (*ai.ModelResponse, error) {
	if len(resp.Candidates) == 0 {
		// Return empty response for streaming compatibility
		return &ai.ModelResponse{
			Message: &ai.Message{Role: ai.RoleModel},
			Usage:   &ai.GenerationUsage{},
		}, nil
	}

	r, err := translateCandidate(resp.Candidates[0])
	if err != nil {
		return nil, err
	}

	if r.Usage == nil {
		r.Usage = &ai.GenerationUsage{}
	}

	if u := resp.UsageMetadata; u != nil {
		r.Usage.InputTokens = int(u.PromptTokenCount)
		r.Usage.OutputTokens = int(u.CandidatesTokenCount)
		r.Usage.TotalTokens = int(u.TotalTokenCount)
	}

	// Handle PromptFeedback for blocked requests
	if resp.PromptFeedback != nil {
		if resp.PromptFeedback.BlockReason != "" {
			r.FinishReason = ai.FinishReasonBlocked
			r.FinishMessage = fmt.Sprintf("Prompt blocked: %s", resp.PromptFeedback.BlockReason)
		}
	}

	return r, nil
}

// validToolName checks whether the provided tool name matches the regex
func validToolName(n string) bool {
	re := regexp.MustCompile(toolNameRegex)
	return re.MatchString(n)
}

// EnsureThoughtSignatures adds synthetic thought signatures to function calls
// in the active loop for Gemini 3 Compatibility
func EnsureThoughtSignatures(contents []*genai.Content) []*genai.Content {
	activeLoopStartIndex := -1
	for i := len(contents) - 1; i >= 0; i-- {
		content := contents[i]
		if content.Role == "user" {
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

	newContents := make([]*genai.Content, len(contents))
	copy(newContents, contents)

	for i := activeLoopStartIndex; i < len(newContents); i++ {
		content := newContents[i]
		if content.Role != "model" || len(content.Parts) == 0 {
			continue
		}

		for j, part := range content.Parts {
			if part.FunctionCall != nil {
				if len(part.ThoughtSignature) == 0 {
					newParts := make([]*genai.Part, len(content.Parts))
					copy(newParts, content.Parts)
					newParts[j] = &genai.Part{
						FunctionCall:     part.FunctionCall,
						ThoughtSignature: []byte("skip_thought_signature_validator"),
					}
					newContents[i] = &genai.Content{
						Role:  content.Role,
						Parts: newParts,
					}
				}
				break
			}
		}
	}

	return newContents
}
