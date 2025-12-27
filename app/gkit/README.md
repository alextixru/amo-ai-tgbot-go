# Genkit AI Agent

AI агент для управления amoCRM через естественный язык.

## Архитектура

```
User Message → Router Flow → Specialized Flow → CRM Tools → amoCRM SDK
```

### Router Flow
Классифицирует intent пользователя и направляет в нужный flow:
- `leads` — работа со сделками
- `contacts` — работа с контактами  
- `tasks` — работа с задачами
- `general` — общие вопросы (без CRM операций)

### Specialized Flows
Каждый flow получает свой набор tools:

| Flow | Tools |
|------|-------|
| `leads_flow` | getLeads, createLead, updateLead |
| `contacts_flow` | getContacts, createContact |
| `tasks_flow` | getTasks, createTask |

## План реализации

1. [x] Базовый Chat Flow (работает)
2. [ ] Router Flow — классификация intent
3. [ ] Подключение Router → Chat Flow
4. [ ] Leads Flow + tools
5. [ ] Contacts Flow + tools

## Принципы

- **Минимум tools на flow** — LLM видит только релевантные инструменты
- **Изоляция** — каждый flow тестируется независимо
- **Наблюдаемость** — все flows видны в Genkit UI

## Концепции архитектуры

### Определения

| Понятие | Роль |
|---------|------|
| **Flow** | Бизнес-сценарий, оркестрация |
| **Agent** | Роль внутри flow (не объект, а ответственность) |
| **Tool** | Единственный способ воздействовать на мир |

### Иерархия

```
┌────────────┐
│ Flow       │  ← orchestration
├────────────┤
│ Agent A    │  ← prompt + tools
│ Agent B    │
├────────────┤
│ Tools      │  ← side effects
└────────────┘
```

**Агент не живёт сам по себе — он всегда подчинён flow.**

### Поток данных

```
User input
   ↓
Router agent (intent, task type)
   ↓
Specialized flow / agent
```

### В Genkit

- **Tool = граница доверия** — LLM не знает реализацию
- **Flow знает последствия** — контролирует side effects

### ❌ Антипаттерны

- "Главный агент", который сам решает, что делать
- Бесконечный loop: reasoning → tool → reasoning

---

## Genkit Core Abstractions

В Developer UI каждый action (Flow, Prompt, Generate, Model) имеет 4 секции:

| Секция | Что это | Пример |
|--------|---------|--------|
| **Input** | Входные данные для action | `{ message, user_context }` |
| **Context** | Auth/session контекст (заполняется при Firebase деплое) | `auth`, `app`, `instanceIdToken` |
| **Output** | Результат выполнения | `{ response }` |
| **Attributes** | OpenTelemetry span атрибуты | `trace_id`, `span_id`, `genkit.name` |

### Иерархия вызовов

```
chat (Flow) ─── 1.93s
│   Input:  ChatInput { message, user_context }
│   Output: ChatOutput { response }
│
└── user_chat (Prompt) ─── 3ms
    │   Input:  { query, user_context }
    │   Output: rendered prompt
    │
    └── generate (Util) ─── 1.88s
        │   Input:  prompt + config
        │   Output: model response
        │
        └── ollama/model (Model) ─── 1.87s
                Input:  messages
                Output: text
```

### Context vs UserContext

- **Context** (в UI) — auth контекст при Firebase деплое. Пустой локально!
- **UserContext** (наш) — данные пользователя из amoCRM, передаются в Input

### Когда Context заполняется?

| Сценарий | Context |
|----------|---------|
| Локальный запуск с Ollama | ❌ Пустой |
| Firebase Cloud Functions (`onCallGenkit`) | ✅ `auth`, `app` |
| HTTP server с `ContextProvider` | ✅ Custom auth |

---

## Genkit Agentic Patterns (официальная документация)

Шкала от надёжных Workflow до гибких Agents:

```
Надёжность ←─────────────────────────────────────→ Гибкость

WORKFLOW                  HYBRID                    AGENT
├── Sequential           ├── Tool Calling          └── Autonomous
├── Routing              └── Iterative                 Operation
└── Parallel                 Refinement
```

---

### 1. Sequential Processing (Workflow)

Фиксированная последовательность LLM вызовов. Предсказуемо, легко отлаживать.

```go
// Пример: research → draft → review
researchResult, _ := genkit.Generate(ctx, g,
    ai.WithPrompt("Research the topic: "+topic),
)

draftResult, _ := genkit.Generate(ctx, g,
    ai.WithPrompt("Write a draft based on: "+researchResult.Text()),
)

finalResult, _ := genkit.Generate(ctx, g,
    ai.WithPrompt("Review and polish: "+draftResult.Text()),
)
```

**Когда использовать:** Pipeline обработки данных, генерация контента с review.

---

### 2. Conditional Routing (Workflow)

Ветвление на основе классификации LLM.

```go
// Классифицируем intent
classifyResult, _ := genkit.Generate(ctx, g,
    ai.WithPrompt("Classify intent: leads, contacts, tasks, or general. Query: "+query),
    ai.WithOutputFormat(ai.OutputFormatText),
)

intent := strings.TrimSpace(classifyResult.Text())

// Роутим
switch intent {
case "leads":
    return leadsFlow.Run(ctx, input)
case "contacts":
    return contactsFlow.Run(ctx, input)
case "tasks":
    return tasksFlow.Run(ctx, input)
default:
    return chatFlow.Run(ctx, input)
}
```

**Когда использовать:** Router агент, выбор специализации.

---

### 3. Parallel Execution (Workflow)

Несколько LLM вызовов параллельно для скорости или разных точек зрения.

```go
var wg sync.WaitGroup
results := make([]string, 3)

prompts := []string{
    "Analyze as sales manager...",
    "Analyze as support agent...",
    "Analyze as product manager...",
}

for i, prompt := range prompts {
    wg.Add(1)
    go func(i int, p string) {
        defer wg.Done()
        resp, _ := genkit.Generate(ctx, g, ai.WithPrompt(p+content))
        results[i] = resp.Text()
    }(i, prompt)
}
wg.Wait()

// Агрегируем результаты
finalResult, _ := genkit.Generate(ctx, g,
    ai.WithPrompt("Synthesize these perspectives: "+strings.Join(results, "\n\n")),
)
```

**Когда использовать:** Мультиперспективный анализ, ускорение обработки.

---

### 4. Tool Calling (Hybrid)

LLM сам решает какие tools вызвать. Genkit автоматически выполняет tools.

```go
// Определяем tool
searchLeadsTool := genkit.DefineTool(g, "searchLeads",
    "Search leads in amoCRM by query",
    func(ctx *ai.ToolContext, input struct {
        Query string `json:"query" jsonschema_description:"Search query"`
        Limit int    `json:"limit,omitempty"`
    }) ([]Lead, error) {
        return sdk.Leads().Search(ctx, input.Query, input.Limit)
    },
)

// Используем в генерации
response, _ := genkit.Generate(ctx, g,
    ai.WithPrompt("Find all leads related to: "+userQuery),
    ai.WithTools(searchLeadsTool, createLeadTool, updateLeadTool),
)
// Genkit автоматически вызывает tools и возвращает финальный ответ
```

**Когда использовать:** Доступ к внешним данным, CRUD операции.

---

### 5. Iterative Refinement (Hybrid)

Цикл самоулучшения: генерация → критика → улучшение.

```go
const maxIterations = 3
var draft string

for i := 0; i < maxIterations; i++ {
    // Генерируем/улучшаем
    if i == 0 {
        draftResp, _ := genkit.Generate(ctx, g,
            ai.WithPrompt("Write initial draft for: "+topic),
        )
        draft = draftResp.Text()
    }
    
    // Критикуем
    critiqueResp, _ := genkit.Generate(ctx, g,
        ai.WithPrompt("Critique this draft. List specific improvements:\n\n"+draft),
    )
    
    // Проверяем качество
    if strings.Contains(critiqueResp.Text(), "no improvements needed") {
        break
    }
    
    // Улучшаем
    improveResp, _ := genkit.Generate(ctx, g,
        ai.WithPrompt("Improve draft based on feedback:\n\nDraft: "+draft+"\n\nFeedback: "+critiqueResp.Text()),
    )
    draft = improveResp.Text()
}
```

**Когда использовать:** Генерация качественного контента, code review.

---

### 6. Autonomous Operation (Agent)

Агент сам планирует и выполняет, пока не достигнет цели. Максимальная гибкость.

```go
var history []*ai.Message
history = append(history, ai.NewUserMessage(ai.NewTextPart(userGoal)))

for {
    response, err := genkit.Generate(ctx, g,
        ai.WithMessages(history...),
        ai.WithTools(searchTool, createTool, updateTool, completeTool),
    )
    if err != nil {
        return "", err
    }
    
    // Обновляем историю
    history = response.History()
    
    // Проверяем завершение: нет tool requests и stop
    if response.FinishReason() == "stop" && len(response.ToolRequests()) == 0 {
        return response.Text(), nil // Агент закончил
    }
    
    // Защита от бесконечного цикла
    if len(history) > 20 {
        return "", errors.New("max iterations exceeded")
    }
}
```

**Когда использовать:** Сложные многошаговые задачи, исследование.

**⚠️ Осторожно:** Может зациклиться, требует ограничений.

---

### 7. Stateful Interactions (Bonus)

Сохранение истории между вызовами для диалога.

```go
// Хранилище истории (в production — Redis/DB)
var historyStore = make(map[string][]*ai.Message)

func loadHistory(sessionID string) []*ai.Message {
    return historyStore[sessionID]
}

func saveHistory(sessionID string, history []*ai.Message) {
    historyStore[sessionID] = history
}

// Flow с состоянием
statefulChatFlow := genkit.DefineFlow(g, "statefulChat",
    func(ctx context.Context, req ChatRequest) (string, error) {
        // 1. Загрузить историю
        history := loadHistory(req.SessionID)
        
        // 2. Добавить новое сообщение
        history = append(history, ai.NewUserMessage(ai.NewTextPart(req.Message)))
        
        // 3. Генерировать с историей
        response, err := genkit.Generate(ctx, g,
            ai.WithMessages(history...),
        )
        if err != nil {
            return "", err
        }
        
        // 4. Сохранить обновлённую историю
        saveHistory(req.SessionID, response.History())
        
        return response.Text(), nil
    },
)
```

**Когда использовать:** Чат-боты, многоходовые диалоги.

---

## Применение паттернов для amoCRM бота

| Паттерн | Применение |
|---------|------------|
| **Routing** | Router Flow → выбор режима по ролям |
| **Sequential** | AnalyzeLead: getLead → getContacts → getNotes → summary |
| **Tool Calling** | CRUD операции со сделками/контактами |
| **Stateful** | История диалога в Telegram сессии |
| **Iterative** | Улучшение формулировок задач |

**GitHub примеры:** [genkit-ai/samples/agentic-patterns](https://github.com/genkit-ai/samples/tree/main/agentic-patterns)

---

## 🔮 AI SDK (на будущее)

**Идея:** Вынести переиспользуемую AI-логику в отдельный SDK.

### Что точно в SDK:

- **Tools (Layer 5)** — обёртки Genkit Tools над amoCRM SDK методами
  - `searchLeads`, `getLead`, `createLead`, `updateLead`
  - `getContacts`, `createContact`, `getTasks`, `createTask`
  - Универсальны для любого amoCRM проекта

### Под вопросом:

- **Flows (Layer 4)** — уровень абстракции пока не ясен
  - ❓ Насколько универсальны `AnalyzeLead`, `CreateLeadWizard`?
  - ❓ Или flows слишком привязаны к конкретной бизнес-логике?
  - ❓ Может, только "примитивные" flows в SDK, а сложные — в приложении?

### Что остаётся в приложении:

- **Modes (Layer 3)** — специфичны под роли конкретной компании
- **Router (Layer 2)** — бизнес-логика конкретного чат-бота
- **Interface (Layer 1)** — Telegram/REST/Widget
