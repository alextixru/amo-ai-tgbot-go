package flows

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/tihn/amo-ai-tgbot-go/app/gkit/prompts"
	"github.com/tihn/amo-ai-tgbot-go/app/gkit/session"
)

// ChatInput — вход для Chat Flow
type ChatInput struct {
	SessionID   string         `json:"session_id"` // Telegram chat ID or unique session identifier
	Message     string         `json:"message"`
	UserContext map[string]any `json:"user_context,omitempty"`
}

// ChatOutput — выход для Chat Flow
type ChatOutput struct {
	Response string `json:"response"`
}

// DefineChatFlow регистрирует Chat Flow с поддержкой истории и возвращает функцию для его запуска
func DefineChatFlow(
	g *genkit.Genkit,
	model ai.Model,
	tools []ai.ToolRef,
	store session.Store,
) func(context.Context, ChatInput) (ChatOutput, error) {

	flow := genkit.DefineFlow(g, "chat",
		func(ctx context.Context, input ChatInput) (ChatOutput, error) {
			// 1. Загружаем историю сессии
			history := store.Load(input.SessionID)

			// 2. Если история пустая — добавляем системный prompt
			if len(history) == 0 {
				cfg := prompts.ExtractPromptConfig(input.UserContext)
				systemPrompt := prompts.BuildSystemPrompt(cfg)
				history = append(history, ai.NewSystemMessage(ai.NewTextPart(systemPrompt)))
			}

			// 3. Добавляем новое сообщение пользователя
			history = append(history, ai.NewUserMessage(ai.NewTextPart(input.Message)))

			// 4. Генерируем ответ с полной историей
			resp, err := genkit.Generate(ctx, g,
				ai.WithModelName(model.Name()),
				ai.WithMessages(history...),
				ai.WithTools(tools...),
			)
			if err != nil {
				return ChatOutput{}, fmt.Errorf("generate: %w", err)
			}

			// 5. Сохраняем обновлённую историю (включает tool calls и responses)
			store.Save(input.SessionID, resp.History())

			return ChatOutput{Response: resp.Text()}, nil
		},
	)
	return flow.Run
}
