# Agent Init Flow

## Порядок инициализации

```
NewAgent(client, sdk)
  │
  ├── 1. account_context.Load(ctx, sdk)
  │       └── GET pipelines + statuses
  │       └── GET custom fields (leads, contacts, companies)
  │       └── GET users
  │       └── GET loss reasons, sources
  │       └── резолвер заполнен
  │
  ├── 2. tools.NewRegistry(g, resolver, services...)
  │       └── resolver передаётся в registry
  │
  ├── 3. registry.RegisterAll()
  │       └── каждый tool формирует описание из resolver.PipelineNames() и т.д.
  │       └── DefineTool("entities", динамическое описание, handler)
  │       └── нейронка видит реальные названия аккаунта в описаниях
  │
  ├── 4. flows.DefineChatFlow(g, model, registry.AllTools(), store)
  │       └── нейронка инициализируется уже с готовыми tools
  │
  └── 5. Бот готов — нейронка знает воронки, поля, пользователей по именам
```

## Что даёт этот порядок

- Нейронка видит в описании tool реальные значения: "Воронки: Новые клиенты, VIP"
- Нейронка передаёт строки, резолвер конвертит в ID под капотом
- Один источник правды — резолвер. Два потребителя — описание tool + mapping.go

## Где реализуется каждый шаг

| Шаг | Файл |
|-----|------|
| Load() | `internal/services/account_context/loader.go` |
| Resolver | `internal/services/account_context/resolver.go` |
| Динамические описания | `app/gkit/tools/registry.go` + каждый `Register*Tool()` |
| Резолвинг имя→ID | `internal/adapters/entities/mapping.go` |
| Сборка агента | `app/gkit/agent.go` |
