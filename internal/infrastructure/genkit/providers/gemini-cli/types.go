package geminicli

import "google.golang.org/genai"

// LoadCodeAssist types for project discovery
type LoadCodeAssistRequest struct {
	CloudaicompanionProject string         `json:"cloudaicompanionProject,omitempty"`
	Metadata                ClientMetadata `json:"metadata"`
}

type ClientMetadata struct {
	IDEType     string `json:"ideType,omitempty"`
	Platform    string `json:"platform,omitempty"`
	PluginType  string `json:"pluginType,omitempty"`
	DuetProject string `json:"duetProject,omitempty"`
}

type LoadCodeAssistResponse struct {
	CurrentTier             *UserTier `json:"currentTier,omitempty"`
	CloudaicompanionProject string    `json:"cloudaicompanionProject,omitempty"`
}

type UserTier struct {
	ID string `json:"id"`
}

// RetrieveUserQuotaRequest — запрос квот пользователя
type RetrieveUserQuotaRequest struct {
	Project   string `json:"project"`
	UserAgent string `json:"userAgent,omitempty"`
}

// BucketInfo — информация о квоте
type BucketInfo struct {
	RemainingAmount   string  `json:"remainingAmount,omitempty"`
	RemainingFraction float64 `json:"remainingFraction,omitempty"`
	ResetTime         string  `json:"resetTime,omitempty"`
	TokenType         string  `json:"tokenType,omitempty"`
	ModelID           string  `json:"modelId,omitempty"`
}

// RetrieveUserQuotaResponse — ответ с квотами
type RetrieveUserQuotaResponse struct {
	Buckets []BucketInfo `json:"buckets,omitempty"`
}

// CountTokensRequest — запрос для подсчёта токенов
type CountTokensRequest struct {
	Model    string           `json:"model"`
	Project  string           `json:"project,omitempty"`
	Contents []*genai.Content `json:"contents"`
}

// CountTokensResponse — ответ с количеством токенов
type CountTokensResponse struct {
	TotalTokens int32 `json:"totalTokens"`
}
