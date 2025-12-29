package flows

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// ChatInput — вход для Chat Flow
type ChatInput struct {
	Message     string         `json:"message"`
	UserContext map[string]any `json:"user_context,omitempty"`
}

// ChatOutput — выход для Chat Flow
type ChatOutput struct {
	Response string `json:"response"`
}

// DefineChatFlow регистрирует Chat Flow и возвращает функцию для его запуска
func DefineChatFlow(g *genkit.Genkit, model ai.Model, tools []ai.ToolRef) func(context.Context, ChatInput) (ChatOutput, error) {
	// Получаем prompt из Dotprompt файла
	chatPrompt := genkit.LookupPrompt(g, "user_chat")
	if chatPrompt == nil {
		panic("prompt 'user_chat' not found in prompts directory")
	}

	flow := genkit.DefineFlow(g, "chat",
		func(ctx context.Context, input ChatInput) (ChatOutput, error) {
			// Передаём всё в prompt как есть
			promptInput := map[string]any{
				"query":        input.Message,
				"user_context": input.UserContext,
			}

			resp, err := chatPrompt.Execute(ctx,
				ai.WithModelName(model.Name()),
				ai.WithInput(promptInput),
				ai.WithTools(tools...),
			)
			if err != nil {
				return ChatOutput{}, fmt.Errorf("prompt execute: %w", err)
			}
			return ChatOutput{Response: resp.Text()}, nil
		},
	)
	return flow.Run
}
