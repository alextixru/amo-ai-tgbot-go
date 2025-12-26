package gkit

// ChatInput — вход для Chat Flow
type ChatInput struct {
	Message string `json:"message"`
}

// ChatOutput — выход для Chat Flow
type ChatOutput struct {
	Response string `json:"response"`
}
