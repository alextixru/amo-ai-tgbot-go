# Промт для агента: подключение Google ADK

## Контекст проекта

Go Telegram-бот для управления amoCRM через AI-агента. Genkit был удалён, проект компилируется с stub-агентом. Нужно подключить Google ADK (Agent Development Kit) чтобы AI снова заработал.

Ветка: `migration/remove-genkit`

## Текущая структура (после удаления Genkit)

```
cmd/bot/main.go                          — точка входа, создаёт agent.NewAgent(ctx, sdk)
app/agent/
├── agent.go                             — Agent struct со stub Process(), держит registry + session store
├── tools/                               — 12 tool-файлов (НЕ ПОДКЛЮЧАТЬ СЕЙЧАС, это отдельный этап)
│   ├── registry.go                      — ToolDefinition{Name, Description, InputSchema, Handler} + Registry
│   ├── entities.go ... (12 файлов)
└── prompts/
    └── system_prompt.go                 — BuildSystemPrompt() string, framework-agnostic

internal/
├── models/chat/message.go              — собственные типы Message (пока не используются, на будущее)
├── services/agent/agent.go             — интерфейс Processor{Process(ctx, sessionID, msg) (string, error)}
├── services/session/session.go         — Store interface + MemoryStore (in-memory, []*chat.Message)
├── services/telegram/service.go        — принимает agent.Processor, вызывает agent.Process()
├── services/crm/                       — 12 CRM-сервисов (чистый Go, не трогать)
├── services/auth/                      — Google OAuth (не трогать)
├── infrastructure/crm/                 — amoCRM SDK client (не трогать)
├── infrastructure/google/oauth/        — OAuth helpers (не трогать)
└── infrastructure/telegram/            — Telegram client (не трогать)

config/config.go                        — конфигурация из env vars

_research/adk-utils-go/                 — клонированный репо с примерами и утилитами ADK
```

## Что нужно сделать

Подключить ADK **без tools** — только LLM-провайдер, агент, сессии и Runner. Бот должен уметь вести диалог через AI, но пока без вызова CRM-инструментов. Подключение tools — отдельный следующий этап.

### 1. Установить зависимости

```bash
go get google.golang.org/adk@latest
go get google.golang.org/genai@latest
go get github.com/achetronic/adk-utils-go@latest
```

### 2. Создать ADK LLM-адаптер

Файл: `internal/infrastructure/llm/provider.go`

Использовать `github.com/achetronic/adk-utils-go/genai/openai` (работает с Ollama через OpenAI-compatible API).

Прочитать конфигурацию из `config.Config`:
- `cfg.OllamaURL` — base URL (по умолчанию `http://localhost:11434/v1`)
- `cfg.OllamaModel` — название модели

```go
import genaiopenai "github.com/achetronic/adk-utils-go/genai/openai"

llmModel := genaiopenai.New(genaiopenai.Config{
    BaseURL:   cfg.OllamaURL + "/v1",
    ModelName: cfg.OllamaModel,
    APIKey:    "ollama", // Ollama не требует ключ, но поле обязательное
})
```

Проверь `_research/adk-utils-go/genai/openai/` и `_research/adk-utils-go/examples/` для точного API.

### 3. Заменить stub Agent на ADK Runner (без tools)

Файл: `app/agent/agent.go`

```go
import (
    "google.golang.org/adk/agent/llmagent"
    "google.golang.org/adk/runner"
    "google.golang.org/adk/session"
    "google.golang.org/genai"
)

type Agent struct {
    runner         *runner.Runner
    sessionService session.Service
    appName        string
}

func NewAgent(ctx context.Context, sdk *amocrm.SDK, llmModel model.LLM) (*Agent, error) {
    // Создать ADK agent БЕЗ tools:
    adkAgent, _ := llmagent.New(llmagent.Config{
        Name:        "crm-assistant",
        Model:       llmModel,
        Description: "amoCRM AI assistant",
        Instruction: prompts.BuildSystemPrompt(),
        // Toolsets: пусто — tools подключим позже
    })

    // In-memory session service:
    sessionService := session.InMemoryService()

    // Runner:
    runnr, _ := runner.New(runner.Config{
        AppName:        "amocrm-bot",
        Agent:          adkAgent,
        SessionService: sessionService,
    })

    return &Agent{runner: runnr, sessionService: sessionService, appName: "amocrm-bot"}, nil
}
```

Инициализацию CRM-сервисов и registry можно пока убрать из NewAgent (или оставить закомментированными) — они понадобятся при подключении tools.

### 4. Реализовать Process() через ADK Runner

```go
func (a *Agent) Process(ctx context.Context, sessionID, message string) (string, error) {
    // Ensure session exists
    _, err := a.sessionService.Get(ctx, &session.GetRequest{
        AppName: a.appName, UserID: sessionID, SessionID: sessionID,
    })
    if err != nil {
        // Create new session
        a.sessionService.Create(ctx, &session.CreateRequest{
            AppName: a.appName, UserID: sessionID, SessionID: sessionID,
        })
    }

    // Run agent
    userMsg := genai.NewContentFromText(message, genai.RoleUser)
    var result strings.Builder
    for event, err := range a.runner.Run(ctx, sessionID, sessionID, userMsg, agent.RunConfig{}) {
        if err != nil { return "", err }
        if event.Content != nil {
            for _, part := range event.Content.Parts {
                if part.Text != "" {
                    result.WriteString(part.Text)
                }
            }
        }
    }
    return result.String(), nil
}
```

### 5. Обновить main.go

```go
// Создать LLM provider
llmModel := llm.NewProvider(cfg)

// Передать в agent (без genkitClient)
aiAgent, err := agent.NewAgent(ctx, crmClient.SDK(), llmModel)
```

### 6. Обновить config.go

Проверь что используются поля:
- `AI_PROVIDER` — "ollama" (по умолчанию, пока единственный)
- `OLLAMA_URL` — URL Ollama сервера
- `OLLAMA_MODEL` — название модели

Эти поля уже есть в конфиге, убедись что они читаются и передаются в LLM-адаптер.

## Важные ограничения

1. **НЕ ПОДКЛЮЧАТЬ TOOLS** — агент работает без tool calling, просто ведёт диалог. Подключение tools к ADK — отдельный следующий этап
2. **НЕ ТРОГАТЬ** `app/agent/tools/*`, `internal/services/crm/*`, `internal/models/tools/*`, `internal/services/auth/*`, `internal/infrastructure/crm/*`, `app/telegram/*`
3. **НЕ СТАВИТЬ Redis/PostgreSQL** — это отдельный этап. Использовать in-memory session service
4. Проект должен компилироваться: `go build ./...`
5. Изучай `_research/adk-utils-go/` для точного API — не полагайся на знания о стабильности API

## Ссылки для изучения

- `_research/adk-utils-go/examples/simple/main.go` — минимальный пример (начни с него)
- `_research/adk-utils-go/examples/session-memory/main.go` — с сессиями
- `_research/adk-utils-go/genai/openai/` — OpenAI-compatible провайдер (для Ollama)
- `google.golang.org/adk/agent/llmagent` — создание агента
- `google.golang.org/adk/runner` — Runner
- `google.golang.org/adk/session` — Session Service

## Критерий готовности

1. `go build ./...` — компилируется
2. Бот запускается с Ollama
3. AI отвечает на сообщения через ADK Runner (простой диалог, без tool calling)
4. История диалога сохраняется между сообщениями (in-memory)
5. System prompt из `app/agent/prompts/system_prompt.go` используется
