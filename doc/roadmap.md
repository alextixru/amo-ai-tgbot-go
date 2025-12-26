# Roadmap

## Фаза 1: Init ✅

Базовая инфраструктура проекта.

- [x] Структура проекта (cmd, internal, doc)
- [x] Telegram бот (go-telegram/bot)
- [x] AI агент (Genkit + Ollama)
- [x] amoCRM SDK интеграция
- [x] Debug режим через ENV

---

## Фаза 2: Domain Layer

Доменная модель и бизнес-логика.

- [ ] Определить доменные сущности (Lead, Contact, Company)
- [ ] Описать use cases (GetLeads, CreateContact, UpdateLead)
- [ ] Интерфейсы репозиториев

---

## Фаза 3: Service Layer

Сервисный слой с адаптерами.

- [ ] CRM сервис с адаптером к amoCRM SDK
- [ ] Genkit Tools как адаптеры к сервисам
- [ ] AI контекст через ai.GlobalState

---

## Фаза 4: Infrastructure (Production-ready)

Инфраструктурные компоненты для production.

- [ ] Структурированное логирование (slog или zerolog)
- [ ] Конфигурация с валидацией (viper/envconfig)
- [ ] Graceful shutdown с timeout
- [ ] Health checks (/health эндпоинт)
- [ ] Метрики (Prometheus, опционально)
- [ ] Custom error types
