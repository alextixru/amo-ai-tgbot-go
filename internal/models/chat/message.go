// Package chat defines framework-agnostic message types for AI conversations.
// These types replace ai.Message from Genkit and will map to ADK types later.
package chat

// Role constants for message roles.
const (
	RoleSystem = "system"
	RoleUser   = "user"
	RoleModel  = "model"
	RoleTool   = "tool"
)

// Message represents a single message in a conversation.
type Message struct {
	Role    string `json:"role"` // system, user, model, tool
	Content []Part `json:"content"`
}

// Part is a single content part of a message.
type Part struct {
	Text     string        `json:"text,omitempty"`
	ToolCall *ToolCall     `json:"tool_call,omitempty"`
	ToolResp *ToolResponse `json:"tool_response,omitempty"`
}

// ToolCall represents an LLM request to call a tool.
type ToolCall struct {
	Name  string         `json:"name"`
	Input map[string]any `json:"input"`
}

// ToolResponse represents the result of a tool call.
type ToolResponse struct {
	Name   string `json:"name"`
	Output any    `json:"output"`
}

// NewSystemMessage creates a system message with text content.
func NewSystemMessage(text string) *Message {
	return &Message{
		Role:    RoleSystem,
		Content: []Part{{Text: text}},
	}
}

// NewUserMessage creates a user message with text content.
func NewUserMessage(text string) *Message {
	return &Message{
		Role:    RoleUser,
		Content: []Part{{Text: text}},
	}
}

// NewModelMessage creates a model message with text content.
func NewModelMessage(text string) *Message {
	return &Message{
		Role:    RoleModel,
		Content: []Part{{Text: text}},
	}
}

// Text returns the concatenated text content of the message.
func (m *Message) Text() string {
	var s string
	for _, p := range m.Content {
		s += p.Text
	}
	return s
}
