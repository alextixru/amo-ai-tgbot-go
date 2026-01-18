package geminicli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/genai"
)

const (
	BaseURL    = "https://cloudcode-pa.googleapis.com"
	APIVersion = "v1internal"
)

type Client struct {
	httpClient *http.Client
	projectID  string
}

func NewClient(ctx context.Context, credsPath string) (*Client, error) {
	ts, err := GetTokenSource(ctx, credsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get token source: %v", err)
	}

	httpClient := oauth2.NewClient(ctx, ts)
	c := &Client{httpClient: httpClient}

	// Discovery and onboarding
	if err := c.setupUser(ctx); err != nil {
		return nil, fmt.Errorf("failed to setup user: %v", err)
	}

	return c, nil
}

func (c *Client) setupUser(ctx context.Context) error {
	// Try to get projectId via loadCodeAssist
	reqBody := LoadCodeAssistRequest{
		Metadata: ClientMetadata{
			IDEType:    "GEMINI_CLI",
			Platform:   "DARWIN_ARM64",
			PluginType: "GEMINI",
		},
	}

	var resp LoadCodeAssistResponse
	err := c.postRaw(ctx, "loadCodeAssist", reqBody, &resp)
	if err != nil {
		return err
	}

	if resp.CloudaicompanionProject != "" {
		c.projectID = resp.CloudaicompanionProject
	} else {
		return fmt.Errorf("no projectId received from loadCodeAssist")
	}

	return nil
}

// CARequest wraps the request in the format expected by Code Assist API
type CARequest struct {
	Model        string                 `json:"model"`
	Project      string                 `json:"project,omitempty"`
	UserPromptID string                 `json:"user_prompt_id,omitempty"`
	Request      *VertexGenerateRequest `json:"request"`
}

// VertexGenerateRequest matches the inner request structure of Code Assist API
type VertexGenerateRequest struct {
	Contents          []*genai.Content       `json:"contents"`
	SystemInstruction *genai.Content         `json:"systemInstruction,omitempty"`
	Tools             []*genai.Tool          `json:"tools,omitempty"`
	ToolConfig        *genai.ToolConfig      `json:"toolConfig,omitempty"`
	SafetySettings    []*genai.SafetySetting `json:"safetySettings,omitempty"`
	GenerationConfig  *GenerationConfig      `json:"generationConfig,omitempty"`
}

// GenerationConfig is a subset of SDK config that matches internal API
type GenerationConfig struct {
	Temperature      *float32              `json:"temperature,omitempty"`
	TopP             *float32              `json:"topP,omitempty"`
	TopK             *float32              `json:"topK,omitempty"`
	CandidateCount   int32                 `json:"candidateCount,omitempty"`
	MaxOutputTokens  int32                 `json:"maxOutputTokens,omitempty"`
	StopSequences    []string              `json:"stopSequences,omitempty"`
	ResponseMIMEType string                `json:"responseMimeType,omitempty"`
	ThinkingConfig   *genai.ThinkingConfig `json:"thinkingConfig,omitempty"`
}

// CAResponse wraps the response from Code Assist API
type CAResponse struct {
	Response *genai.GenerateContentResponse `json:"response"`
	TraceID  string                         `json:"traceId,omitempty"`
}

func (c *Client) Generate(ctx context.Context, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	// Add thought signatures for Gemini 3 Pro compatibility
	contents = EnsureThoughtSignatures(contents)

	// Build the inner request
	innerReq := &VertexGenerateRequest{
		Contents: contents,
	}

	// Apply config fields to request
	if config != nil {
		innerReq.GenerationConfig = &GenerationConfig{
			Temperature:      config.Temperature,
			TopP:             config.TopP,
			TopK:             config.TopK,
			CandidateCount:   config.CandidateCount,
			MaxOutputTokens:  config.MaxOutputTokens,
			StopSequences:    config.StopSequences,
			ResponseMIMEType: config.ResponseMIMEType,
			ThinkingConfig:   config.ThinkingConfig,
		}
		innerReq.SystemInstruction = config.SystemInstruction
		innerReq.Tools = config.Tools
		innerReq.ToolConfig = config.ToolConfig
		innerReq.SafetySettings = config.SafetySettings
	}

	// Wrap in Code Assist format
	caReq := &CARequest{
		Model:        model,
		Project:      c.projectID,
		UserPromptID: fmt.Sprintf("genkit-%s", c.projectID),
		Request:      innerReq,
	}

	var caResp CAResponse
	err := c.postRaw(ctx, "generateContent", caReq, &caResp)
	if err != nil {
		return nil, err
	}

	return caResp.Response, nil
}

// postRaw is used for all Code Assist API methods
func (c *Client) postRaw(ctx context.Context, method string, body any, out any) error {
	url := fmt.Sprintf("%s/%s:%s", BaseURL, APIVersion, method)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return err
		}
	}

	return nil
}
