# Unified Tools Architecture

> Идея: микросервис с группированными endpoints вместо 1:1 маппинга SDK методов.

## Концепция

Вместо создания отдельного gRPC endpoint для каждого метода SDK (100+ методов), группируем их в **12 логических инструментов** с unified input.

## Преимущества

| Аспект | Выгода |
|--------|--------|
| **AI-friendly** | LLM удобнее работать с 12 tools + switch, чем с 100+ отдельных методов |
| **Меньше endpoint'ов** | Проще поддерживать, меньше boilerplate |
| **Семантическая группировка** | `entities` / `activities` — понятная ментальная модель |
| **Единый интерфейс** | `{ action, entity_type, id, filter, data }` — универсально |
| **Reference Context** | Справочники вшиваются автоматически |

## Tools Map

| Mode | Tool | Сервисов | Описание |
|------|------|----------|----------|
| **Work** | entities | 3 | Сделки, контакты, компании |
| **Work** | activities | 9 | Задачи, примечания, звонки, события, файлы, связи, теги, подписки, чаты |
| **Work** | complex_create | 1 | Создание сделки с контактами |
| **Work** | products | 1+ | Товары и привязка к сущностям |
| **Work** | catalogs | 2 | Справочники (каталоги) |
| **Work** | files | 1 | Файловое хранилище |
| **Work** | unsorted | 1 | Неразобранное |
| **Admin** | admin_schema | 4 | Кастомные поля, причины отказа, источники |
| **Admin** | admin_pipelines | 1 | Воронки и этапы |
| **Admin** | admin_users | 2 | Пользователи и роли |
| **Admin** | admin_integrations | 5 | Вебхуки, виджеты, шаблоны |
| **Retention** | customers | 5 | Покупатели, бонусы, транзакции |

## Unified Input Pattern

```go
type ToolInput struct {
    Action     string         `json:"action"`               // обязательное
    EntityType string         `json:"entity_type,omitempty"`
    ID         int            `json:"id,omitempty"`
    Filter     map[string]any `json:"filter,omitempty"`
    Data       map[string]any `json:"data,omitempty"`
}
```

Пример JSON:
```json
{
  "entity_type": "leads",
  "action": "search",
  "filter": { "pipeline_id": 123 }
}
```

## Proto Example

```protobuf
service EntitiesService {
  rpc Execute(EntitiesRequest) returns (EntitiesResponse);
}

message EntitiesRequest {
  string entity_type = 1;  // leads | contacts | companies
  string action = 2;       // search | get | create | update | delete | link
  int64 id = 3;
  google.protobuf.Struct filter = 4;
  google.protobuf.Struct data = 5;
  LinkTarget link_to = 6;
}
```

## Reference Context (auto-inject)

Справочники **не являются отдельными Tools**. Данные автоматически вшиваются в контекст:

| Service | Данные | Вшивается в |
|---------|--------|-------------|
| AccountService | Аккаунт, домен, права | Session context |
| UsersService | Пользователи | entities, activities |
| PipelinesService | Воронки + статусы | entities |
| CustomFieldsService | Кастомные поля | entities |
| CurrenciesService | Валюты | Session context |

## Архитектура в crm-core

```
crm-core (rk-boot + go-kit)
       ↓
   12 gRPC endpoints (entities, activities, products, ...)
       ↓
   handler → switch by action → SDK call
       ↓
   amocrm-sdk-go
```

## Что учесть при реализации

1. **Валидация** — conditional validation в handler:
   ```go
   switch input.Action {
   case "get":
       if input.ID == 0 {
           return nil, errors.New("id required for 'get'")
       }
   }
   ```

2. **Типизация ответов** — использовать `oneof` или `google.protobuf.Struct` для полиморфных ответов

3. **Документация** — Swagger/gRPC-UI должны показывать какие `action` поддерживает какой tool

---

---

## Industry Patterns

Это устоявшиеся паттерны, которые используют крупные компании:

### GraphQL
- **Один endpoint** (`/graphql`) → внутри dispatch по query/mutation
- Facebook, GitHub, Shopify, Stripe
- То же самое, только другой протокол

### MCP (Model Context Protocol)
Anthropic ввели паттерны для работы с tools:
- **Toolhost Pattern** — консолидирует много tools за одним dispatcher
- **Tool Grouping Pattern** — группировка tools в логические единицы
- **Hierarchical MCP Pattern** — иерархия по domain

> *"Toolhost pattern consolidates many tools behind a single dispatcher, particularly useful when clients struggle with large tool lists"*

### API Gateway + Action Routing
AWS, Kong, Apigee:
- Single endpoint → routing по `body.action`
- Пример: `POST /api` с `{ "action": "createLead", ... }`

### JSON-RPC / gRPC с Dispatch
Ethereum, финтех:
```json
{ "method": "eth_getBalance", "params": [...] }
```

---

## Про дублирование операций

Схема **сознательно дублирует** некоторые операции:

```
entities.link()         ← связать сделку с контактом
activities.links.link() ← то же самое, но из слоя activities
```

### Почему это адекватно для AI?

| Причина | Объяснение |
|---------|------------|
| **Разные ментальные модели** | LLM может думать "я работаю со сделкой" → `entities` или "я добавляю связь" → `activities.links` |
| **Context window** | AI не нужно держать в голове весь граф сущностей |
| **Fault tolerance** | Если одна формулировка не сработала, LLM попробует другую |

### Когда это НЕ адекватно?

- **Программный API для разработчиков** — дублирование = путаница
- **Строгая типизация** — лучше 1:1 mapping

---

## Когда какой подход?

| Задача | Подходит unified tools? |
|--------|-------------------------|
| **AI-агент (TG бот, MCP)** | ✅ Идеально — меньше tools, context-aware |
| **Микросервис для других сервисов** | ⚠️ Частично — нужна документация |
| **Public API для разработчиков** | ❌ Лучше REST 1:1 с SDK |

---

## Рекомендация: Двухслойная архитектура

```
┌────────────────────────────────────────────────────┐
│                    crm-core                        │
├────────────────────────────────────────────────────┤
│                                                    │
│   AI clients ──► Unified Tools (12 gRPC endpoints) │
│                         │                          │
│   Dev clients ──► REST 1:1 wrapper (опционально)   │
│                         │                          │
│                         ▼                          │
│                  amocrm-sdk-go                     │
│                                                    │
└────────────────────────────────────────────────────┘
```

1. **Internal gRPC** — unified tools (12 endpoints) для AI-клиентов
2. **Thin REST wrapper** (опционально) — если понадобится программный API

---

## Источники

- [toolmap.md](file:///Users/tihn/amo-ai-tgbot-go/app/gkit/tools/toolmap.md) — маппинг SDK методов
- [tools_schema.md](file:///Users/tihn/amo-ai-tgbot-go/app/gkit/tools/tools_schema.md) — схема группировки tools
