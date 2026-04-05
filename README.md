# amoCRM AI Telegram Bot

Telegram-бот для управления amoCRM через естественный язык с использованием AI-агента.

## Концепция

Пользователь общается с ботом в Telegram на естественном языке (на русском). AI-агент понимает намерение, вызывает нужные инструменты и выполняет операции в amoCRM через SDK. Все взаимодействие — AI-first, без ручных форм и кнопок.

```
Пользователь → Telegram → AI Agent → CRM Tools → amoCRM SDK → API
```

## Текущий стек

- **Go 1.24+**
- **Firebase Genkit** — AI-фреймворк (оркестрация, tool calling loop, flows)
- **Ollama / Gemini CLI** — LLM-провайдеры (локальный и облачный)
- **go-telegram/bot** — Telegram Bot API
- **amocrm-sdk-go** — SDK для amoCRM (кастомный форк)
- **OAuth 2.0 + PKCE** — авторизация Google для Gemini

## Структура проекта

```
├── cmd/bot/                  # Точка входа
├── app/
│   ├── gkit/
│   │   ├── agent.go          # AI-агент (оркестрация Genkit)
│   │   ├── tools/            # 12 категорий Genkit-инструментов
│   │   └── prompts/          # System prompt
│   └── telegram/
│       └── handler.go        # Telegram обработчик
├── internal/
│   ├── infrastructure/       # Инфраструктура (Genkit клиент, CRM клиент, Telegram)
│   ├── services/crm/         # Доменные CRM-сервисы (entities, activities, products, ...)
│   ├── services/flows/       # Genkit flows (multi-turn chat)
│   ├── services/session/     # In-memory session store
│   ├── services/auth/        # Google OAuth
│   └── models/               # Типы request/response для tools
└── config/                   # Конфигурация (env vars)
```

## CRM-сервисы (доменный слой)

12 категорий инструментов, покрывающих весь API amoCRM:

| Сервис | Что делает |
|--------|-----------|
| entities | Сделки, контакты, компании (CRUD, поиск, линковка) |
| activities | Задачи, заметки, звонки, события |
| complex_create | Атомарное создание сделка + контакт |
| products | Каталог товаров, привязка к сделкам |
| catalogs | Пользовательские каталоги |
| files | Файловое хранилище |
| unsorted | Неразобранное (входящие лиды) |
| customers | CRM+ управление покупателями |
| admin_schema | Пользовательские поля, группы полей |
| admin_pipelines | Воронки и статусы |
| admin_users | Пользователи и роли |
| admin_integrations | Вебхуки, виджеты, интеграции |

---

## Миграция: Genkit → Google ADK

### Почему переезжаем

Firebase Genkit решает задачу оркестрации вызовов к LLM: собрать messages, передать в модель, обработать tool calls, повторить. Это работает, но на практике для полноценного AI-чата этого недостаточно. Вот что Genkit **не даёт** и что приходится писать руками:

| Проблема | Что сейчас | Что нужно |
|----------|-----------|-----------|
| **Сессии** | In-memory map (`session.MemoryStore`). Теряются при перезапуске. | Персистентные сессии с хранением в Redis/PostgreSQL |
| **История** | Хранится в памяти, лимит 20 сообщений, без персистентности | Хранение в БД, треды (несколько разговоров у одного пользователя) |
| **Долговременная память** | Нет. Агент не помнит ничего между сессиями | Working memory (факты о пользователе), semantic recall (поиск по прошлым разговорам) |
| **Суммаризация контекста** | Нет. При длинных разговорах упираемся в контекстное окно | Автоматическая суммаризация старых сообщений |
| **Сборка контекста** | Полностью ручная: сам грузишь историю, сам собираешь messages[] | Runner, который автоматически загружает сессию, обогащает контекст, вызывает модель |
| **Observability** | Dev UI (localhost:4000) — только для разработки | Langfuse или аналог для production-трейсинга |

Genkit — это **SDK для вызова моделей**, а не платформа для AI-агентов. Он не управляет жизненным циклом разговора.

### Почему Google ADK

**Google ADK (Agent Development Kit)** — фреймворк от Google Cloud / Vertex AI, open-source, Go SDK в GA (v1.0.0 с ноября 2025). Это **следующее поколение** после Genkit, созданное другой командой внутри Google, специально для построения агентных систем.

ADK и Genkit — два разных проекта, не зависят друг от друга, код не шарят. Оба используют `google.golang.org/genai` как низкоуровневый клиент, но архитектурно решают разные задачи:

| | Genkit (Firebase) | ADK (Google Cloud) |
|---|---|---|
| Философия | "Добавь AI в приложение" | "Построй систему агентов" |
| Runner | Нет — сам пишешь flow | Есть — автоматически загружает сессию, собирает контекст, вызывает модель, обрабатывает tools |
| SessionService | `session.Store` (in-memory only) | Интерфейс — подключаешь любой backend (Redis, PostgreSQL) |
| MemoryService | Нет | Есть — долговременная память между сессиями, `LoadMemoryTool` |
| Plugins | Нет | BeforeModel / AfterModel callbacks |
| Multi-agent | Нет | Есть — делегация, sequential/parallel workflows |
| Go SDK статус | Beta | **GA v1.0.0** |

### Почему замена будет относительно простой

Genkit в нашем проекте — это **инфраструктурный слой**, а не доменный. Вся бизнес-логика (CRM-сервисы, tools, amoCRM SDK) не зависит от Genkit. Genkit используется только в:

1. `app/gkit/agent.go` — инициализация агента, вызов `chatFlow`
2. `app/gkit/tools/*.go` — регистрация tools через `genkit.DefineTool`
3. `internal/services/flows/chat.go` — chat flow (загрузка истории → generate → сохранение)
4. `internal/services/session/session.go` — in-memory session store
5. `internal/infrastructure/genkit/` — инициализация Genkit-клиента и провайдеров

Доменные сервисы (`internal/services/crm/*`) и модели (`internal/models/tools/*`) — чистый Go, без импортов Genkit. Они останутся как есть.

При миграции:
- `genkit.DefineTool` → `adk tool definition`
- `genkit.Generate` + ручной flow → `ADK Runner` (автоматический)
- `session.MemoryStore` → `ADK SessionService` (Redis через adk-utils-go)
- Добавляется `MemoryService` (pgvector через adk-utils-go)
- Добавляется Context Guard (автосуммаризация через adk-utils-go)

### Целевая архитектура

```
┌─────────────────────────────────────────────────┐
│                   Клиенты                        │
│         Telegram    /    Web UI (будущее)         │
└──────────────────────┬──────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────┐
│              Go Backend (ADK)                    │
│                                                  │
│  ┌─────────────────────────────────────────┐     │
│  │           ADK Runner                     │     │
│  │  автоматически: загрузка сессии →        │     │
│  │  сборка контекста → вызов LLM →          │     │
│  │  tool calling loop → сохранение          │     │
│  └─────────┬───────────────────┬───────────┘     │
│            │                   │                  │
│  ┌─────────▼─────────┐ ┌──────▼──────────┐      │
│  │  SessionService   │ │  MemoryService   │      │
│  │  (Redis)          │ │  (pgvector)      │      │
│  │  adk-utils-go     │ │  adk-utils-go    │      │
│  └─────��─────────────┘ └─���───────────────┘      │
│                                                  │
│  ┌─────────────────────────────────────────┐     │
│  │         Plugins                          │     │
│  │  Context Guard (суммаризация)            │     │
│  │  Langfuse (observability)                │     │
│  │  adk-utils-go                            │     │
│  └────�����───────────────────────────────────┘     │
│                                                  │
│  ┌───────────────────────────��─────────────┐     │
│  │         CRM Tools (12 категорий)         │     │
│  │  entities, activities, products, ...     │     │
│  │  ↓                                       │     │
│  │  CRM-сервисы (internal/services/crm/)    │     │
│  │  ↓                                       │     │
│  │  amoCRM SDK → amoCRM API                 │     │
│  └─────────────────────────────────────────┘     │
└──────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────┐
│          memory-service (опционально)             │
│  Отдельный микросервис для:                       │
│  - Персистентная история с тредами и forking      │
│  - Access control (права пользователей)           │
│  - Semantic + fulltext search по истории           │
│  - React demo UI для чата                         │
│  - REST/gRPC API                                  │
└─────────────────────────────────��────────────────┘
```

### Компоненты целевого стека

| Компонент | Решение | Роль |
|-----------|---------|------|
| **Agent runtime** | Google ADK (`google.golang.org/adk`) | Runner, agent, tool calling loop |
| **Sessions** | adk-utils-go → Redis | Персистентные сессии с TTL |
| **Memory** | adk-utils-go → PostgreSQL + pgvector | Долговременная память, semantic search |
| **Context management** | adk-utils-go → Context Guard plugin | Автоматическая суммаризация при приближении к лимиту контекстного окна |
| **Observability** | adk-utils-go → Langfuse plugin | Трейсинг каждого LLM-вызова с промптами, ответами, токенами |
| **LLM-провайдеры** | adk-utils-go → OpenAI/Anthropic клиенты | Поддержка Ollama, OpenRouter, Claude, GPT и др. |
| **Chat infrastructure** | memory-service (chirino) — опционально | Треды, forking, access control, UI |
| **CRM** | amoCRM SDK (без изменений) | Вся бизнес-логика CRM |
| **Telegram** | go-telegram/bot (без изменений) | Telegram Bot API |

### Исследовательские материалы

В директории `_research/` находятся клонированные репозитории для изучения:

- `_research/adk-utils-go/` — утилиты для ADK (Redis sessions, pgvector memory, Context Guard, Langfuse)
- `_research/memory-service/` — standalone микросервис для хранения и управления историей чатов

### План миграции (высокоуровневый)

1. Изучить ADK Go API и adk-utils-go examples
2. Переписать tool definitions с Genkit на ADK формат
3. Заменить Genkit agent + chat flow на ADK Runner
4. Подключить SessionService (Redis) вместо in-memory store
5. Подк��ючить MemoryService (pgvector) для долговременной памяти
6. Подключить Context Guard для автосуммаризации
7. Удалить Genkit из зависимостей

---

## Запуск

### Переменные окружения

```bash
# Telegram
TELEGRAM_BOT_TOKEN=           # Telegram Bot API токен

# AI Provider
AI_PROVIDER=                   # "ollama" или "gemini-cli"
OLLAMA_URL=                    # http://localhost:11434 (по умолчанию)
OLLAMA_MODEL=                  # Название модели

# amoCRM
AMOCRM_AUTH_MODE=              # "token" или "oauth"
AMOCRM_BASE_URL=               # https://subdomain.amocrm.ru
AMOCRM_ACCESS_TOKEN=           # Для token-режима
```

### Команды

```bash
make dev    # Hot reload + Genkit Dev UI (localhost:4000)
make run    # Запуск без перезагрузки
make build  # Сборка бинарника
```

## Лицензия

Private
