# Domain Layer

Доменная область бота **полностью делегируется** [amoCRM SDK](https://github.com/alextixru/amocrm-sdk-go).

## Архитектурное решение

Бот не определяет собственные модели. Вместо этого:

| Концепция | Источник |
|-----------|----------|
| Сущности (Lead, Contact, Company...) | `sdk/core/models` |
| AI контекст | `sdk/ai` |
| Авторизация и токены | `sdk/core/oauth` |
| Идентификация пользователя | `models.User.ID` из amoCRM |
| Права доступа | Роли пользователя в amoCRM Account |

## Почему так?

1. SDK уже содержит полную доменную модель amoCRM
2. Дублирование структур = двойная поддержка
3. AI получает контекст напрямую из SDK (`ai.GlobalState`)
4. Права пользователя определяются amoCRM, а не ботом

## Использование

```go
import (
    "github.com/alextixru/amocrm-sdk-go/core/models"
    "github.com/alextixru/amocrm-sdk-go/ai"
)

// Работаем с SDK моделями напрямую
func ProcessLead(lead *models.Lead) {
    // ...
}
```
