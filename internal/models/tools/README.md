# Tools Models

Input DTOs для SDK-инструментов Genkit.

## Содержимое

| Файл | Описание |
|------|----------|
| `common.go` | Общие типы: `LinkTarget`, `ParentEntity` |
| `entities.go` | Сделки, контакты, компании |
| `activities.go` | Задачи, примечания, звонки, события, файлы, связи, теги |
| `customers.go` | Покупатели и транзакции |
| `products.go` | Товары |
| `catalogs.go` | Каталоги и элементы |
| `files.go` | Загрузка файлов |
| `unsorted.go` | Неразобранное |
| `complex_create.go` | Комплексное создание сущностей |
| `admin_*.go` | Административные операции |

## Использование

```go
import "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"

// В определении Tool
genkit.DefineTool[tools.EntitiesInput, any](g, "entities", "...", handler)
```

## Принципы

1. Каждое поле имеет `jsonschema_description` для LLM
2. Output — модели SDK напрямую (`github.com/alextixru/amocrm-sdk-go/core/models`)
3. Структуры максимально полные (30+ полей) для Sub-Agent
