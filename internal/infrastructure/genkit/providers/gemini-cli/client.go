package geminicli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/genai"
)

const (
	BaseURL    = "https://cloudcode-pa.googleapis.com"
	APIVersion = "v1internal"
	UserAgent  = "Genkit-CLI-Go/1.0"
)

type Client struct {
	httpClient *http.Client
	projectID  string
	sessionID  string
}

type ClientOptions struct {
	CredsPath string
	NoBrowser bool
}

type userAgentTransport struct {
	base http.RoundTripper
}

func (t *userAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", UserAgent)
	return t.base.RoundTrip(r)
}

func NewClient(ctx context.Context, opts ClientOptions) (*Client, error) {
	ts, err := GetTokenSource(ctx, opts.CredsPath, opts.NoBrowser)
	if err != nil {
		return nil, fmt.Errorf("failed to get token source: %v", err)
	}

	httpClient := oauth2.NewClient(ctx, ts)
	httpClient.Transport = &userAgentTransport{base: httpClient.Transport}

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
			Platform:   getPlatform(),
			PluginType: "GEMINI",
		},
	}

	var resp LoadCodeAssistResponse
	err := c.postWithRetry(ctx, "loadCodeAssist", reqBody, &resp)
	if err != nil {
		// Check if user is affected by VPC Service Controls
		var apiErr *APIError
		if errors.As(err, &apiErr) && isVPCSCAffectedUser(apiErr) {
			// VPC-SC users are assumed to be on standard tier
			// We can proceed without explicit project ID
			return nil
		}
		return err
	}

	if resp.CloudaicompanionProject != "" {
		c.projectID = resp.CloudaicompanionProject
	} else {
		return fmt.Errorf("no projectId received from loadCodeAssist")
	}

	return nil
}

func getPlatform() string {
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return "DARWIN_ARM64"
		}
		return "DARWIN_AMD64"
	case "linux":
		if runtime.GOARCH == "arm64" {
			return "LINUX_ARM64"
		}
		return "LINUX_AMD64"
	case "windows":
		return "WINDOWS_AMD64"
	default:
		return "PLATFORM_UNSPECIFIED"
	}
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
	SessionID         string                 `json:"session_id,omitempty"`
	CachedContent     string                 `json:"cachedContent,omitempty"`
}

// GenerationConfig is a subset of SDK config that matches internal API
type GenerationConfig struct {
	Temperature        *float32              `json:"temperature,omitempty"`
	TopP               *float32              `json:"topP,omitempty"`
	TopK               *float32              `json:"topK,omitempty"`
	CandidateCount     int32                 `json:"candidateCount,omitempty"`
	MaxOutputTokens    int32                 `json:"maxOutputTokens,omitempty"`
	StopSequences      []string              `json:"stopSequences,omitempty"`
	ResponseMIMEType   string                `json:"responseMimeType,omitempty"`
	ThinkingConfig     *genai.ThinkingConfig `json:"thinkingConfig,omitempty"`
	PresencePenalty    *float32              `json:"presencePenalty,omitempty"`
	FrequencyPenalty   *float32              `json:"frequencyPenalty,omitempty"`
	Seed               *int32                `json:"seed,omitempty"`
	ResponseModalities []string              `json:"responseModalities,omitempty"`
	ResponseLogprobs   bool                  `json:"responseLogprobs,omitempty"`
	Logprobs           *int32                `json:"logprobs,omitempty"`
}

// CAResponse wraps the response from Code Assist API
type CAResponse struct {
	Response *genai.GenerateContentResponse `json:"response"`
	TraceID  string                         `json:"traceId,omitempty"`
}

// APIError represents a structured error from Code Assist API
type APIError struct {
	StatusCode int
	Message    string
	Details    []ErrorDetail `json:"details,omitempty"`
	RawBody    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.RawBody)
}

func (e *APIError) IsRetryable() bool {
	return e.StatusCode == 429 || e.StatusCode >= 500
}

// isVPCSCAffectedUser проверяет, является ли ошибка результатом VPC Service Controls
func isVPCSCAffectedUser(err *APIError) bool {
	if err == nil {
		return false
	}
	for _, detail := range err.Details {
		if detail.Reason == "SECURITY_POLICY_VIOLATED" {
			return true
		}
	}
	return false
}

type ErrorDetail struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Domain  string `json:"domain,omitempty"`
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
	err := c.postWithRetry(ctx, "generateContent", caReq, &caResp)
	if err != nil {
		return nil, err
	}

	return caResp.Response, nil
}

// RetrieveUserQuota получает информацию о квотах пользователя
func (c *Client) RetrieveUserQuota(ctx context.Context) (*RetrieveUserQuotaResponse, error) {
	req := &RetrieveUserQuotaRequest{
		Project:   c.projectID,
		UserAgent: UserAgent,
	}
	var resp RetrieveUserQuotaResponse
	err := c.postWithRetry(ctx, "retrieveUserQuota", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CountTokens подсчитывает токены для заданного контента
func (c *Client) CountTokens(ctx context.Context, model string, contents []*genai.Content) (int32, error) {
	// Strip thoughts for CountTokens API compatibility (matches TS implementation)
	contents = StripThoughtsFromContents(contents)

	req := &CountTokensRequest{
		Model:    model,
		Project:  c.projectID,
		Contents: contents,
	}
	var resp CountTokensResponse
	err := c.postWithRetry(ctx, "countTokens", req, &resp)
	if err != nil {
		return 0, err
	}
	return resp.TotalTokens, nil
}

// StripThoughtsFromContents removes Thought parts from contents for CountTokens API compatibility.
func StripThoughtsFromContents(contents []*genai.Content) []*genai.Content {
	result := make([]*genai.Content, 0, len(contents))
	for _, c := range contents {
		if c == nil {
			continue
		}
		var filteredParts []*genai.Part
		for _, p := range c.Parts {
			if p == nil {
				continue
			}
			if p.Thought {
				continue // Skip pure thought parts
			}
			// Clear signature from non-thought parts (CountTokens doesn't like it)
			if len(p.ThoughtSignature) > 0 {
				newPart := *p
				newPart.ThoughtSignature = nil
				filteredParts = append(filteredParts, &newPart)
			} else {
				filteredParts = append(filteredParts, p)
			}
		}
		if len(filteredParts) > 0 {
			result = append(result, &genai.Content{
				Role:  c.Role,
				Parts: filteredParts,
			})
		}
	}
	return result
}

// GenerateStream implements streaming for Code Assist API using SSE
func (c *Client) GenerateStream(ctx context.Context, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (<-chan *genai.GenerateContentResponse, <-chan error) {
	respChan := make(chan *genai.GenerateContentResponse)
	errChan := make(chan error, 1)

	contents = EnsureThoughtSignatures(contents)

	innerReq := &VertexGenerateRequest{
		Contents: contents,
	}

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

	caReq := &CARequest{
		Model:        model,
		Project:      c.projectID,
		UserPromptID: fmt.Sprintf("genkit-%s", c.projectID),
		Request:      innerReq,
	}

	go func() {
		defer close(respChan)
		defer close(errChan)

		jsonBody, err := json.Marshal(caReq)
		if err != nil {
			errChan <- err
			return
		}

		url := fmt.Sprintf("%s/%s:streamGenerateContent?alt=sse", BaseURL, APIVersion)
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			errChan <- err
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errChan <- &APIError{StatusCode: resp.StatusCode, RawBody: string(body)}
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		var bufferedLines []string

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data: ") {
				bufferedLines = append(bufferedLines, strings.TrimPrefix(line, "data: "))
			} else if line == "" {
				if len(bufferedLines) == 0 {
					continue
				}
				var chunk CAResponse
				if err := json.Unmarshal([]byte(strings.Join(bufferedLines, "\n")), &chunk); err != nil {
					errChan <- fmt.Errorf("failed to decode SSE chunk: %v", err)
					return
				}
				if chunk.Response != nil {
					respChan <- chunk.Response
				}
				bufferedLines = nil
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- err
		}
	}()

	return respChan, errChan
}

func (c *Client) postWithRetry(ctx context.Context, method string, body any, out any) error {
	maxRetries := 3
	baseDelay := 1 * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := c.postRaw(ctx, method, body, out)
		if err == nil {
			return nil
		}

		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.IsRetryable() && attempt < maxRetries {
			delay := baseDelay * time.Duration(1<<attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
				continue
			}
		}
		return err
	}
	return errors.New("max retries exceeded")
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
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			RawBody:    string(body),
		}
		// Try to parse structured error
		var errData struct {
			Error struct {
				Message string        `json:"message"`
				Details []ErrorDetail `json:"details"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &errData); err == nil {
			apiErr.Message = errData.Error.Message
			apiErr.Details = errData.Error.Details
		}
		return apiErr
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return err
		}
	}

	return nil
}
