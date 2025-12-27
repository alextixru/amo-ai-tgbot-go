package gkit

// ChatInput — вход для Chat Flow
type ChatInput struct {
	Message     string         `json:"message"`
	UserContext map[string]any `json:"user_context,omitempty"`
}

// ChatOutput — выход для Chat Flow
type ChatOutput struct {
	Response string `json:"response"`
}
