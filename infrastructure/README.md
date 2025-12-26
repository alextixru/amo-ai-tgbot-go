# Infrastructure Layer

Клиенты и адаптеры к внешним системам.

## Пакеты

| Пакет | Файл | Назначение |
|-------|------|------------|
| `genkit/` | `client.go` | Genkit + Ollama клиент |
| `telegram/` | `bot.go` | Telegram Bot API клиент |
| `crm/` | `client.go` | amoCRM SDK обёртка |
| `config/` | `config.go` | Конфигурация из ENV |

## Принцип

Infrastructure содержит **клиенты** — инициализацию и подключение.

Бизнес-логика (AI Agent, Telegram обработчики) находится в `app/`.

## Зависимости

```
cmd/bot/main.go
    ├── infrastructure/genkit   → создаёт Client
    ├── infrastructure/telegram → создаёт Bot  
    ├── infrastructure/crm      → создаёт Client
    └── infrastructure/config   → загружает Config
```
