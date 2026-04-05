# Репорт: admin_integrations

**Дата:** 2026-04-04
**Статус:** Рефакторинг завершён. Код соответствует плану.

---

## Текущее состояние сервиса

### Архитектура

Сервис состоит из трёх слоёв:

| Слой | Путь | Роль |
|------|------|------|
| Модель входа | `internal/models/tools/admin_integrations.go` | Типизированный input для LLM |
| Сервис | `internal/services/crm/admin_integrations/` | Бизнес-логика, вызовы SDK |
| Tool-handler | `app/gkit/tools/admin_integrations.go` | Маршрутизация, DTO-выход, форматирование |

### Поддерживаемые операции

**Webhooks (`layer: "webhooks"`)**
- `list` / `search` — список вебхуков с фильтром по `destination`; ответ через DTO без `AccountID`/`CreatedBy`, timestamps в RFC3339
- `subscribe` — создание вебхука; читает `input.Destination` и `input.EventTypes` (типизированный массив, 20 событий перечислены в jsonschema_description)
- `unsubscribe` — удаление вебхука; те же поля

**Widgets (`layer: "widgets"`)**
- `list` / `search` — список виджетов с фильтром по limit/page; возвращает сырые `[]*models.Widget`
- `get` — по `code`
- `install` — по `code` + `settings map[string]any`
- `uninstall` — по `code`

**Website Buttons (`layer: "website_buttons"`)**
- `list` / `search` — фильтр limit/page, `with []string` на верхнем уровне; ответ через `websiteButtonOut` без `AccountID`
- `get` — по `id` + `with`
- `create` — из типизированного `WebsiteButtonData` (`Name`, `PipelineID *int`, `TrustedWebsites`, `IsDuplicationControlEnabled *bool`)
- `update` — **BUG FIX применён**: `req.SourceID = input.ID` явно в handler; ранее `SourceID` не заполнялся из-за `json:"-"` в SDK-модели
- `add_chat` — по `id` (source_id)

**Chat Templates (`layer: "chat_templates"`)**
- `list` / `search` — фильтр limit/page/external_ids; клиентская фильтрация по `TemplateType` (API не поддерживает); ответ через `chatTemplateOut` без `AccountID`, timestamps в RFC3339
- `create` — из `ChatTemplateData`; добавлено в рамках рефакторинга
- `update` — из `ChatTemplateData` + `id`; добавлено в рамках рефакторинга
- `delete` — по `id`
- `delete_many` — по `ids []int`
- `send_review` — по `id`; ответ через `chatTemplateReviewOut` с полем `id` ревью (BUG FIX в SDK)
- `update_review` — по `id` + `review_id` (верхний уровень input, не через `data`) + `review_status`

**Short Links (`layer: "short_links"`)**
- `list` / `search` — фильтр limit/page; возвращает сырые `[]models.ShortLink`
- `create` — одиночная: `url` + `entity_id` / `entity_type`; батч: `urls []string`; подписи изменены на `models.ShortLink` (EntityID/EntityType передаются в SDK)
- `delete` — по `id`

---

## Что было сделано в рефакторинге (и соответствует коду)

### Баги исправлены

| # | Проблема | Статус в коде |
|---|----------|---------------|
| BUG-1 | `UpdateWebsiteButton`: `SourceID` не заполнялся, PATCH уходил на `/website_buttons/0` | **Исправлено.** `req.SourceID = input.ID` — строка 184 в `app/gkit/tools/admin_integrations.go` |
| BUG-2 | `ChatTemplateReview` без `ID` — цепочка `send_review → update_review` была нерабочей | **Исправлено.** Поле `ID int` добавлено в `amocrm-sdk-go/core/models/chat_template.go` (строка 116). Поле присутствует. |

### Вход — модель `AdminIntegrationsInput`

| # | Изменение | Статус |
|---|-----------|--------|
| 3 | `Destination` и `EventTypes []string` вынесены на верхний уровень для webhooks | Реализовано. Поля присутствуют в модели. |
| 4 | `With []string` вынесен из `filter` на верхний уровень | Реализовано. |
| 5 | `WebsiteButtonWithDeleted = "deleted"` задокументирован в jsonschema_description поля `With` | Реализовано. |
| 6 | `EntityID int` и `EntityType string` добавлены для `short_links: create` | Реализовано. Передаются в `models.ShortLink` при вызове SDK. |
| 7 | Перечисление webhook-событий в jsonschema_description поля `EventTypes` | Реализовано. 20 событий перечислены. |
| 8 | `TemplateType string` в `IntegrationsFilter` для фильтрации шаблонов по типу | Реализовано. Клиентская фильтрация в handler. |
| 9 | `data map[string]any` для `website_buttons` заменён на `WebsiteButtonData` | Реализовано. |
| 10 | `IDs []int` оставлен на верхнем уровне (не перенесён в `data`) | Осознанное решение. Поле семантически чистое и используется только в `delete_many`. |
| 11 | `ReviewID int` и `ReviewStatus string` вынесены на верхний уровень для `update_review` | Реализовано. |
| 12 | `URL string` и `URLs []string` добавлены как явные поля для short_links | Реализовано. |
| — | `ChatTemplateData` — типизированная структура для create/update шаблонов | Реализовано. |

### Выход — tool-handler

| # | Изменение | Статус |
|---|-----------|--------|
| 13 | `CreateChatTemplate` и `UpdateChatTemplate` в `Service` интерфейсе и реализации | Реализовано. В `service.go` и `chat_templates.go`. |
| 14 | Unix timestamps форматируются в RFC3339 через `unixToRFC3339()` | Реализовано. Применяется к Webhook, ChatTemplate, ChatTemplateReview. |
| 15 | `AccountID` убран из ответов через DTO-структуры | Реализовано. `webhookOut`, `websiteButtonOut`, `chatTemplateOut` не содержат `AccountID`. |
| 16 | `CreateShortLink` принимает `models.ShortLink` вместо `string url` | Реализовано. Сигнатура изменена в сервисе и handler. |

### SDK и инфраструктура

| # | Изменение | Статус |
|---|-----------|--------|
| 17 | `replace` directive в `go.mod` для локального SDK | Подтверждено в log.md. Позволяет использовать локальные изменения SDK без публикации. |

---

## Расхождения между документами и кодом

Расхождений нет. Всё описанное в REPORT.md (предыдущая версия) и log.md реализовано в коде:

- `ChatTemplateReview.ID` — поле присутствует в SDK (`chat_template.go:116`)
- BUG-FIX `req.SourceID = input.ID` — строка 184 в `app/gkit/tools/admin_integrations.go`
- Типизированные структуры `WebsiteButtonData`, `ChatTemplateData` — присутствуют в `internal/models/tools/admin_integrations.go`
- `CreateChatTemplate`, `UpdateChatTemplate` — реализованы в `chat_templates.go`
- `CreateShortLink(ctx, models.ShortLink)` — новая сигнатура присутствует в `service.go` и `short_links.go`
- `With []string` на верхнем уровне input — присутствует
- `EventTypes`, `Destination`, `ReviewID`, `ReviewStatus`, `URL`, `URLs`, `EntityID`, `EntityType` — все поля присутствуют

---

## Что не было сделано (осознанные решения)

| Пункт из AUDIT.md | Решение |
|-------------------|---------|
| `Widget.Version interface{}` | Не исправлено. Это SDK-поле, изменение затронет десериализацию. Приемлемо для LLM. |
| `PipelineID *int` для WebsiteButton/Widget не резолвится в имя | По `REFACTORING.md`, `admin_integrations` не требует резолвинга справочников. PipelineID задокументирован в `WebsiteButtonData`. |
| `PageMeta` для webhooks не возвращается | API webhooks не возвращает пагинацию, отбрасывание `_` оправдано. |
| `WebsiteButtonWithDeleted` добавляется автоматически | Не автоматизировано. LLM должна явно передать `with: ["deleted"]`. Задокументировано в jsonschema_description. |
| `Webhook.AccountID`, `Webhook.CreatedBy` без имён | Убраны из ответного DTO целиком. LLM не получает числовые ID, которые не может интерпретировать. |

---

## Файлы, изменённые в рамках рефакторинга

| Файл | Что изменено |
|------|-------------|
| `/Users/tihn/ssm/amocrm-sdk-go/core/models/chat_template.go` | Добавлено `ID int` в `ChatTemplateReview` |
| `/Users/tihn/amo-ai-tgbot-go/go.mod` | Добавлен `replace` directive для локального SDK |
| `/Users/tihn/amo-ai-tgbot-go/internal/models/tools/admin_integrations.go` | Полная переработка: `Destination`, `EventTypes`, `With`, `WebsiteButtonData`, `ChatTemplateData`, `ReviewID`, `ReviewStatus`, `URL`, `URLs`, `EntityID`, `EntityType`, `TemplateType` |
| `/Users/tihn/amo-ai-tgbot-go/internal/services/crm/admin_integrations/service.go` | Добавлены `CreateChatTemplate`, `UpdateChatTemplate`; изменена сигнатура `CreateShortLink` |
| `/Users/tihn/amo-ai-tgbot-go/internal/services/crm/admin_integrations/chat_templates.go` | Реализации `CreateChatTemplate`, `UpdateChatTemplate` |
| `/Users/tihn/amo-ai-tgbot-go/internal/services/crm/admin_integrations/short_links.go` | `CreateShortLink` принимает `models.ShortLink`; `CreateShortLinks` принимает `[]models.ShortLink` |
| `/Users/tihn/amo-ai-tgbot-go/app/gkit/tools/admin_integrations.go` | Полная переработка handlers: типизированные входы, DTO-выходы, BUG-FIX SourceID, RFC3339 timestamps, убран AccountID |
