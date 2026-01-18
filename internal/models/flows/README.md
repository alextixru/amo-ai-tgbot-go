# Flows Models

Упрощённые Input DTOs для Tool Flows.

## Идея

Main Agent видит **только эти структуры** — компактные схемы с 6-7 полями вместо 30+.

## Содержимое

| Файл | Описание |
|------|----------|
| `base.go` | `FlowMode` (direct/complex), `BaseFlowInput` |
| `entities.go` | `EntitiesFlowInput` — упрощённый input для entities flow |

## Использование

```go
import "github.com/tihn/amo-ai-tgbot-go/internal/models/flows"

type EntitiesFlowInput struct {
    flows.BaseFlowInput
    EntityType string `json:"entity_type"`
    Action     string `json:"action,omitempty"`
    ID         int    `json:"id,omitempty"`
    // ...
}
```

## Режимы

| Mode | Описание |
|------|----------|
| `direct` | Прямой вызов SDK операции |
| `complex` | Передача Sub-Agent для сложной логики |

## Связь с tools/

```
flows.EntitiesFlowInput (Main Agent)
         ↓
    Tool Flow
         ↓
tools.EntitiesInput (Sub-Agent / SDK Adapter)
```
