# Models

Пакет содержит структуры данных для взаимодействия с LLM через Genkit.

## Структура

| Папка | Описание |
|-------|----------|
| `tools/` | Входные данные для прямых инструментов (SDK). Содержат полный набор полей. |
| `flows/` | Входные данные для Tool Flows. Упрощены для Main Agent. |

## Принципы

1. **Input** → структуры здесь (с `jsonschema_description`).
2. **Output** → модели SDK напрямую (`github.com/alextixru/amocrm-sdk-go/core/models`).
3. **Разделение** → Main Agent видит только модели из `flows/`.

## Использование

### Для инструментов (Tools)
```go
import "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
```

### Для флоу (Flows)
```go
import "github.com/tihn/amo-ai-tgbot-go/internal/models/flows"
```
