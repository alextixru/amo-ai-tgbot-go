# Services Layer (Business Logic)

Этот слой содержит бизнес-логику приложения и маппинг данных между Domain DTO и моделями SDK.

## Архитектура

Каждый сервис в этой папке соответствует логической группе функций (например, `entities`, `activities`, `admin`).

### Обязательные компоненты:

1. **README.md**: Описания бизнес-логики сервиса.
2. **service.go**: Определение интерфейса `Service` и его базовой реализации.
3. **Функциональные файлы**: Разделение логики по сущностям (например, `leads.go`, `contacts.go` внутри `entities`).

## Взаимодействие

```
Transport (Genkit Tool) → Service (Business Logic) → amoCRM SDK
```

Services знают о:
- `models/` (Transport DTOs как входные данные)
- `github.com/alextixru/amocrm-sdk-go/core/models` (SDK модели)
- `github.com/alextixru/amocrm-sdk-go/core/adapters` (SDK сервисы)
