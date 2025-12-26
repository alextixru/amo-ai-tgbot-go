# Application Layer

Use cases и оркестрация бизнес-логики.

## Структура

Каждый файл — один use case:
- `process_message.go` — ProcessMessageUseCase
- `get_leads.go` — GetLeadsUseCase
- `change_mode.go` — ChangeModeUseCase

## Правила

1. Зависит только от `domain/`
2. Использует интерфейсы из `domain/ports/`
3. Не знает о конкретных реализациях (amoCRM, Genkit)
