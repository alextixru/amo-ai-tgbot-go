package flows

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// ChatInput — вход для Chat Flow
type ChatInput struct {
	Message string `json:"message"`
}

// ChatOutput — выход для Chat Flow
type ChatOutput struct {
	Response string `json:"response"`
}

// DefineChatFlow регистрирует Chat Flow и возвращает функцию для его запуска
func DefineChatFlow(g *genkit.Genkit, model ai.Model) func(context.Context, ChatInput) (ChatOutput, error) {
	// Получаем prompt из Dotprompt файла
	chatPrompt := genkit.LookupPrompt(g, "chat")
	if chatPrompt == nil {
		panic("prompt 'chat' not found in prompts directory")
	}

	flow := genkit.DefineFlow(g, "chat",
		func(ctx context.Context, input ChatInput) (ChatOutput, error) {
			// Выполняем prompt с входными данными
			resp, err := chatPrompt.Execute(ctx,
				ai.WithModelName(model.Name()),
				ai.WithInput(map[string]any{"message": input.Message}),
			)
			if err != nil {
				return ChatOutput{}, fmt.Errorf("prompt execute: %w", err)
			}
			return ChatOutput{Response: resp.Text()}, nil
		},
	)
	return flow.Run
}
