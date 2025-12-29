# amoCRM AI Bot: Layered Architecture

**Подход: Сверху вниз (Top-Down)**

```
┌─────────────────────────────────────────────────────────────────────┐
│ Layer 1: INTERFACE                                                  │
│ Telegram Bot | Custom UI | amoCRM Widget                            │
├─────────────────────────────────────────────────────────────────────┤
│ Layer 2: ROUTER + DEFAULT CHAT                                      │
│ Routing по режимам, fallback chat flow                              │
├─────────────────────────────────────────────────────────────────────┤
│ Layer 3: MODE FLOWS                                                 │
│ Reader Mode | Sales Mode | Manager Mode | Admin Mode                │
├─────────────────────────────────────────────────────────────────────┤
│ Layer 4: SPECIALIZED FLOWS                                          │
│ AnalyzeLead | CreateLead | MyTasks | TeamDashboard | ...            │
├─────────────────────────────────────────────────────────────────────┤
│ Layer 5: TOOLS                                                      │
│ Все инструменты из amoCRM SDK                                       │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Layer 1: Interface

**Цель:** Подключение к пользовательскому интерфейсу

### Варианты

| Interface | Описание | Сложность |
|-----------|----------|-----------|
| Telegram Bot | Текущая реализация | ✅ Есть |
| REST API | HTTP endpoints для flows | Низкая |
| amoCRM Widget | Встраивание в интерфейс amo | Средняя |
| Custom Web UI | Standalone web app | Высокая |

### Telegram Bot (текущий)

```
app/telegram/
├── handler.go            # Message handler
├── formatter.go          # Response formatting
└── keyboards.go          # Inline keyboards
```

### Чеклист Layer 1

- [ ] Обновить handler.go → использовать router flow
- [ ] Форматирование ответов для Telegram
- [ ] Inline keyboards для confirmation
- [ ] Error handling

---

## Layer 2: Router + Default Chat

**Цель:** Единая точка входа, routing по режимам, fallback chat

### Структура

```
app/gkit/
├── router.go             # Router Flow
├── chat.go               # Default Chat Flow (fallback)
└── agent.go              # Main Agent (entry point)
```

### Router Flow

```go
routerFlow := genkit.DefineFlow(g, "router",
    func(ctx context.Context, input Input) (Output, error) {
        // 1. Определить режим по правам
        mode := determineMode(input.UserContext.Rights)
        
        // 2. Попробовать intent → specialized flow
        if flow := matchFlow(input.Message, mode.AvailableFlows()); flow != nil {
            return flow.Run(ctx, input)
        }
        
        // 3. Fallback: chat с tools режима
        return chatFlow.Run(ctx, ChatInput{
            Message:      input.Message,
            SystemPrompt: mode.SystemPrompt(),
            Tools:        mode.AvailableTools(),
        })
    },
)
```

### Default Chat Flow

```go
chatFlow := genkit.DefineFlow(g, "chat",
    func(ctx context.Context, input ChatInput) (ChatOutput, error) {
        return genkit.Generate(ctx, g,
            ai.WithPrompt(input.SystemPrompt + "\n\n" + input.Message),
            ai.WithTools(input.Tools...),
        )
    },
)
```

### Чеклист Layer 2

- [ ] `router.go` — routing logic
- [ ] `chat.go` — default chat with tools
- [ ] `agent.go` — entry point
- [ ] `router.prompt` — intent classification (optional)

---

## Layer 3: Mode Flows

**Цель:** Режимы работы, управляющие flows и tools из нижних слоёв

### Структура

```
app/gkit/modes/
├── mode.go               # Interface Mode
├── reader.go             # Reader Mode
├── sales.go              # Sales Mode
├── manager.go            # Manager Mode
└── admin.go              # Admin Mode
```

### Mode Interface

```go
type Mode interface {
    Name() string
    SystemPrompt() string
    AvailableTools() []ai.Tool
    AvailableFlows() []FlowDef
    CanAccess(action string) bool
}
```

### Режимы

| Mode | Rights Check | Tools (Layer 5) | Flows (Layer 4) |
|------|--------------|-----------------|-----------------|
| **Reader** | view: A | search*, get*, getPipelines | AnalyzeLead (view), MyTasks |
| **Sales** | add/edit: D | Reader + create*, update*, link* | + CreateLeadWizard, QualifyLead |
| **Manager** | delete: D | Sales + delete*, getUsers | + TeamDashboard, PipelineReport, BulkUpdate |
| **Admin** | структура | Manager + createPipeline, createField | + все flows |

### Prompts для режимов

```
app/gkit/prompts/
├── reader_mode.prompt      # "Ты помощник для просмотра..."
├── sales_mode.prompt       # "Ты ассистент менеджера..."
├── manager_mode.prompt     # "Ты ассистент руководителя..."
└── admin_mode.prompt       # "Ты ассистент администратора..."
```

### Чеклист Layer 3

- [ ] `mode.go` — interface + determineMode()
- [ ] `reader.go` + reader_mode.prompt
- [ ] `sales.go` + sales_mode.prompt
- [ ] `manager.go` + manager_mode.prompt
- [ ] `admin.go` + admin_mode.prompt

---

## Layer 4: Specialized Flows

**Цель:** Бизнес-сценарии, использующие tools из Layer 5

### Структура

```
app/gkit/flows/
├── analyze_lead.go       # Sequential: полный анализ сделки
├── create_lead.go        # Iterative: wizard создания
├── qualify_lead.go       # Sequential: квалификация
├── my_tasks.go           # Tool Calling: задачи пользователя
├── team_tasks.go         # Parallel: задачи команды
├── pipeline_report.go    # Sequential: отчёт по воронке
├── search.go             # Routing: универсальный поиск
└── bulk_update.go        # Tool Calling: массовые операции
```

### Flows с паттернами

| Flow | Паттерн | Использует Tools | Output |
|------|---------|------------------|--------|
| `AnalyzeLead` | Sequential | getLead, getLinkedContacts, getNotes, getCalls | `LeadAnalysis` |
| `CreateLeadWizard` | Iterative | getPipelines, createLead, createContact, linkEntities | `Lead` |
| `QualifyLead` | Sequential | getLead, updateLead | `QualificationResult` |
| `MyTasks` | Tool Calling | getMyTasks, completeTask | `TaskList` |
| `TeamDashboard` | Parallel + Sequential | getUsers, getMyTasks (per user), aggregate | `Dashboard` |
| `PipelineReport` | Sequential | getPipelines, searchLeads, aggregate | `Report` |
| `UniversalSearch` | Routing | searchLeads, searchContacts, searchCompanies | `SearchResults` |
| `BulkUpdate` | Tool Calling | searchLeads, updateLead (batch) | `UpdateResult` |

### Prompts для flows

```
app/gkit/prompts/
├── analyze_lead.prompt     # System prompt для анализа
├── create_lead.prompt      # Wizard prompts
├── qualify_lead.prompt     # Квалификационные критерии
└── report.prompt           # Генерация отчётов
```

### Чеклист Layer 4

- [ ] `analyze_lead.go` + prompt
- [ ] `create_lead.go` + prompt
- [ ] `qualify_lead.go` + prompt
- [ ] `my_tasks.go`
- [ ] `team_dashboard.go`
- [ ] `pipeline_report.go` + prompt
- [ ] `search.go`
- [ ] `bulk_update.go`

---

## Layer 5: Tools

**Цель:** 12 unified tools покрывающих весь SDK

> Детальная схема: [tools/tools_schema.md](tools/tools_schema.md)

### Архитектура

```
app/gkit/tools/
├── tools.go              # Registry
├── entities.go           # leads, contacts, companies
├── activities.go         # tasks, notes, calls, events, files, links, tags, subscriptions, talks
├── complex_create.go     # создание сделки с контактами
├── products.go           # товары + привязка к сущностям
├── catalogs.go           # справочники (каталоги)
├── files.go              # файловое хранилище
├── unsorted.go           # неразобранное
├── admin_schema.go       # кастомные поля, источники
├── admin_pipelines.go    # воронки и статусы
├── admin_users.go        # пользователи и роли
├── admin_integrations.go # вебхуки, виджеты
└── customers.go          # retention (покупатели)
```

### Tools по режимам

| Mode | Tools | Описание |
|------|-------|----------|
| **Work** | entities, activities, complex_create, products, catalogs, files, unsorted | Основная работа |
| **Admin** | admin_schema, admin_pipelines, admin_users, admin_integrations | Администрирование |
| **Retention** | customers | Покупатели |

### Чеклист Layer 5

- [ ] `tools.go` — registry
- [ ] `entities.go` — leads, contacts, companies (3 сервиса)
- [ ] `activities.go` — tasks, notes, calls и др. (9 сервисов)
- [ ] `complex_create.go` — AddComplex
- [ ] `products.go` — ProductsService + привязка
- [ ] `catalogs.go` — CatalogsService, CatalogElementsService
- [ ] `files.go` — FilesService
- [ ] `unsorted.go` — UnsortedService
- [ ] `admin_schema.go` — CustomFields, LossReasons, Sources
- [ ] `admin_pipelines.go` — PipelinesService
- [ ] `admin_users.go` — UsersService, RolesService
- [ ] `admin_integrations.go` — Webhooks, Widgets, etc.
- [ ] `customers.go` — Customers, BonusPoints, Segments

---

## Зависимости между слоями

```
┌─────────────────────────────────────────────────────────────────────┐
│  [External] Layer 1: Interface (Telegram, REST API, Widget)         │
│             ↓ вызывает                                              │
└─────────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────────┐
│  Layer 2: Router                                                    │
│      ↓ управляет                                                    │
│  Layer 3: Mode Flows                                                │
│      ↓ оркестрирует                                                 │
│  Layer 4: Specialized Flows                                         │
│      ↓ использует                                                   │
│  Layer 5: Tools                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### Направление зависимостей (DDD-стиль)

```
| Верхний слой     | Управляет | Нижний слой       |
|------------------|-----------|-------------------|
| Router           | →         | Mode Flows        |
| Mode Flows       | →         | Specialized Flows |
| Specialized Flows| →         | Tools             |
```

> **Правило:** Верхние слои управляют нижними. Нижние слои ничего не знают о верхних.

### Interface (Layer 1)

Interface — это **внешний слой**, не относящийся к бэкенду:
- Он лишь **вызывает** Router (Layer 2)
- Знает только про вход/выход Router
- Может быть заменён без изменения бэкенда (Telegram → REST API → Widget)

---

## План реализации

### Этап 1: Tools — Layer 5 (3-5 дней)

- [ ] Структура `tools/`
- [ ] Core tools: leads, contacts, companies
- [ ] Tasks & Notes tools
- [ ] Structure tools: pipelines, users
- [ ] Registry и тесты

### Этап 2: Specialized Flows — Layer 4 (3-5 дней)

- [ ] AnalyzeLead flow
- [ ] CreateLeadWizard flow
- [ ] MyTasks flow
- [ ] Search flow
- [ ] Prompts для flows

### Этап 3: Mode Flows — Layer 3 (2-3 дня)

- [ ] Mode interface
- [ ] Reader + Sales modes
- [ ] Manager + Admin modes
- [ ] Mode prompts

### Этап 4: Router — Layer 2 (1-2 дня)

- [ ] Router flow
- [ ] Default chat flow
- [ ] Agent entry point
- [ ] Integration tests

### Этап 5: Interface — Layer 1 (1-2 дня)

- [ ] Update Telegram handler
- [ ] Formatting
- [ ] Error handling

---

**Общий срок: ~2-3 недели**