# Infrastructure Layer

Адаптеры к внешним системам.

## Структура

- `ai/` — Genkit/Ollama адаптер
- `config/` — Конфигурация
- `crm/` — amoCRM SDK адаптер
- `telegram/` — Telegram обработчик

## Правила

1. Реализует интерфейсы из `domain/ports/`
2. Зависит от `domain/` и `application/`
3. Содержит всю работу с внешними API
