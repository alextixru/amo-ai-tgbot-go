# Domain Layer

Чистая бизнес-логика. Не зависит от внешних систем.

## Структура

- `user/` — User, Mode, Permission
- `action/` — Action, ActionResult
- `ports/` — Интерфейсы (CRMRepository, AIProcessor)

## Правила

1. Никаких импортов из `infrastructure/` или `application/`
2. Только стандартная библиотека Go
3. Бизнес-логика без привязки к конкретным технологиям
