# Genkit Tools Guide

Как правильно создавать tools для Genkit.

## Сигнатура Tool

```go
genkit.DefineTool(g, "toolName", "description for LLM",
    func(ctx *ai.ToolContext, input InputType) (OutputType, error) {
        return result, nil
    },
)
```

## 4 компонента Tool

| Компонент | Что это | Требования |
|-----------|---------|------------|
| **Name** | Идентификатор | camelCase, уникальный, понятный LLM |
| **Description** | Описание для LLM | Чётко что делает, когда использовать |
| **Input** | Go struct | `jsonschema_description` теги обязательны |
| **Output** | Go struct/примитив | Сериализуется в JSON автоматически |

---

## Input — самое важное!

LLM видит **только описание и JSON Schema input**. Всё остальное — чёрный ящик.

```go
type SearchLeadsInput struct {
    Query string `json:"query" jsonschema_description:"Search query (by name, phone, email)"`
    Limit int    `json:"limit,omitempty" jsonschema_description:"Max results (default 10)"`
}
```

### Правила Input:

1. `jsonschema_description` — **обязательно** для каждого поля
2. `json:"name"` — имя поля в JSON Schema
3. `omitempty` — опциональное поле
4. Используй примитивы (string, int, bool), не nested structs

### Struct Embedding

Genkit использует библиотеку `github.com/invopop/jsonschema` для генерации JSON Schema.

**Embedded structs автоматически flatten'ятся** (как в стандартном `encoding/json`):

```go
// Базовые поля (переиспользуемые)
type EntityIDInput struct {
    LeadID int `json:"lead_id" jsonschema_description:"Lead ID"`
}

type SearchParams struct {
    Query string `json:"query" jsonschema_description:"Search query"`
    Limit int    `json:"limit" jsonschema_description:"Max results"`
}

// Композиция через embedding
type GetLeadInput struct {
    EntityIDInput  // → lead_id поле будет на верхнем уровне
}

type UpdateLeadInput struct {
    EntityIDInput  // → lead_id
    Name  string `json:"name,omitempty" jsonschema_description:"New name"`
    Price int    `json:"price,omitempty" jsonschema_description:"New price"`
}
```

**Результат JSON Schema (flat):**
```json
{
  "properties": {
    "lead_id": { "type": "integer", "description": "Lead ID" },
    "name": { "type": "string", "description": "New name" },
    "price": { "type": "integer", "description": "New price" }
  }
}
```

---

## Output

```go
// SDK модель
func(...) (*models.Lead, error)

// Слайс
func(...) ([]models.Lead, error)

// Примитив
func(...) (string, error)
```

Genkit сериализует через `json.Marshal` автоматически.

---

## Description

**Плохо:** `"Gets leads"` — LLM не понимает когда использовать

**Хорошо:** `"Search leads in amoCRM by query (name, phone, email). Returns matching leads."`

---

## Что считать одним Tool?

**1 Tool = 1 атомарная операция**

| ✅ Правильно | ❌ Неправильно |
|-------------|---------------|
| `searchLeads` | `manageLeads` (CRUD в одном) |
| `getLead` | `getLeadWithContactsAndTasks` |
| `createLead` | `createLeadIfNotExists` (логика) |

**Принцип:** Tool делает одну вещь. Логику выбора оставляем LLM.

---

## Пример полного Tool

```go
// types.go
type GetLeadInput struct {
    LeadID int `json:"lead_id" jsonschema_description:"ID of the lead to retrieve"`
}

type Lead = models.Lead  // alias на SDK модель

// leads.go
func (r *Registry) defineGetLead() ai.Tool {
    return genkit.DefineTool(r.g, "getLead",
        "Get detailed information about a specific lead by its ID.",
        func(ctx *ai.ToolContext, input GetLeadInput) (*Lead, error) {
            return r.sdk.Leads().GetOne(ctx.Context, input.LeadID, nil)
        },
    )
}
```

---

## Чеклист нового Tool

- [ ] Имя в camelCase
- [ ] Description понятно LLM
- [ ] Input struct в `internal/models/tools/` с `jsonschema_description`
- [ ] Output = SDK модель или alias
- [ ] Регистрация через Registry

## Структура моделей

```
internal/models/
├── tools/           # Input для SDK-инструментов (полные схемы)
│   ├── entities.go
│   ├── activities.go
│   └── ...
└── flows/           # Input для Flow (упрощённые схемы для Main Agent)
    ├── base.go
    ├── entities.go
    └── ...
```
