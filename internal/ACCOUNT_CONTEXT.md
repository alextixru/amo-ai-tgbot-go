# Account Context — Init слой для AI агента

## Проблема

Нейронка получает инструменты с сырыми ID-полями (`pipeline_id int`, `status_id int`, `responsible_user_id int`).
Чтобы создать сделку, она вынуждена сначала вызывать `admin_pipelines → list`, потом резолвить ID вручную.
Каждый запрос = лишние API-вызовы + риск ошибки.

## Решение

Добавить init слой который:
1. При старте агента загружает метаданные аккаунта из amoCRM
2. Хранит кеш `имя → id` для воронок, статусов, пользователей, кастомных полей
3. Резолвер встраивается в адаптеры — нейронка работает с именами, конвертация происходит под капотом

## Архитектура

```
internal/
  services/
    account_context/
      context.go    // структура кеша (AccountContext)
      loader.go     // загрузка из API при старте (Load(ctx, sdk))
      resolver.go   // резолвинг имя→id (PipelineID, StatusID, UserID, FieldID)
```

## Что кешируется

| Данные | Источник | Ключ кеша |
|--------|----------|-----------|
| Воронки | `admin_pipelines → list` | name → id |
| Статусы воронок | `admin_pipelines → list_statuses` | pipeline+name → id |
| Пользователи | `admin_users → list` | name → id |
| Кастомные поля | `admin_schema → list` (по сущностям) | entity+name → id |
| Источники | `admin_schema → sources` | name → id |
| Причины отказа | `admin_schema → loss_reasons` | name → id |

## Изменения в моделях

### До (нейронка оперирует числами)

```go
type EntityData struct {
    PipelineID        int
    StatusID          int
    ResponsibleUserID int
    CustomFieldsValues map[string]any // ключ = field_id или field_code
}
```

### После (нейронка оперирует именами)

```go
type EntityData struct {
    PipelineName        string  // "Новые клиенты"
    StatusName          string  // "Первичный контакт"
    ResponsibleUserName string  // "Иван Петров"
    CustomFieldsValues  map[string]any // ключ = "Бюджет" (человеческое название)
}
```

## Точка встройки резолвера

`internal/adapters/entities/mapping.go` — функция `mapToLead` уже является точкой конвертации.
Резолвер передаётся в сервис через конструктор.

```go
// было
func (s *service) mapToLead(data *EntityData) *models.Lead {
    lead.PipelineID = data.PipelineID
}

// стало
func (s *service) mapToLead(data *EntityData) *models.Lead {
    lead.PipelineID        = s.resolver.PipelineID(data.PipelineName)
    lead.StatusID          = s.resolver.StatusID(data.PipelineName, data.StatusName)
    lead.ResponsibleUserID = s.resolver.UserID(data.ResponsibleUserName)
}
```

## Инициализация

В `app/gkit/agent.go` перед созданием сервисов:

```go
ac, err := account_context.Load(ctx, sdk)
// ac передаётся в адаптеры которым нужен резолвинг
entitiesSvc := entities.New(sdk, ac.Resolver())
```

## Порядок реализации

1. `internal/services/account_context/context.go` — структуры
2. `internal/services/account_context/loader.go` — загрузка из API
3. `internal/services/account_context/resolver.go` — резолвинг
4. `internal/models/tools/entities.go` — замена `*ID int` на `*Name string`
5. `internal/adapters/entities/mapping.go` — встройка резолвера
6. `app/gkit/agent.go` — init при старте
