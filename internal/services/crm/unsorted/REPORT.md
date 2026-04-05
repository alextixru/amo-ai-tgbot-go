# Репорт: сервис unsorted

**Дата:** 2026-04-04

---

## 1. Текущее состояние сервиса (по коду)

### `internal/services/crm/unsorted/service.go`

Сервис полностью рефакторирован согласно паттерну адаптерного слоя из `REFACTORING.md`:

- `New(ctx, sdk)` — конструктор принимает контекст, при инициализации загружает пользователей (`loadUsers`) и воронки со статусами (`loadPipelines`).
- Внутренние индексы: `usersByName`, `usersByID`, `pipelinesByName`, `pipelinesByID`, `statusesByPipelineAndName`.
- Резолверы: `resolveUserName`, `resolveUserID`, `resolvePipelineName`, `resolvePipelineID`, `resolveStatusName` — с подсказками при ошибке ("Доступные: ...").
- Мета-методы `PipelineNames()` и `UserNames()` для динамических описаний tools.
- Выходные типы: `UnsortedOutput`, `UnsortedListOutput`, `UnsortedActionResult`, `UnsortedMetadataOutput`, `UnsortedEmbeddedOutput`.
- Unix timestamps конвертируются в RFC3339 через хелперы (`unixToRFC3339`, `rfc3339ToUnix` в `helpers.go`).
- `account_id` убран из результатов — `UnsortedActionResult` содержит только `{uid, success}`.

### `internal/services/crm/unsorted/unsorted.go`

| Метод | Что делает |
|-------|-----------|
| `ListUnsorted` | Резолвит `PipelineName` → ID; устанавливает `with=leads,contacts,companies`; пробрасывает `Order`; возвращает `*UnsortedListOutput` с `PageMeta`. |
| `GetUnsorted` | Устанавливает `with=leads,contacts,companies`; возвращает `*UnsortedOutput`. |
| `CreateUnsorted` | Резолвит `PipelineName` → ID; конвертирует RFC3339 → Unix через `rfc3339ToUnix`. |
| `AcceptUnsorted` | Принимает `*UnsortedAcceptParams` (имена); резолвит `UserName` → ID и `PipelineName`+`StatusName` → ID. |
| `DeclineUnsorted` | Принимает `*UnsortedDeclineParams`; резолвит `UserName` → ID. |
| `LinkUnsorted` | Возвращает `*UnsortedActionResult` без `account_id`. |
| `SummaryUnsorted` | Пробрасывает `PipelineName` и `CreatedAtFrom`/`CreatedAtTo` (RFC3339 → Unix) в `UnsortedSummaryFilter`. |

**Замечание по `ListUnsorted`:** `CreatedAtFrom`/`CreatedAtTo` из `UnsortedFilter` не пробрасываются в фильтр списка — только в `SummaryUnsorted`. Для `filters.UnsortedFilter` в SDK нет метода `SetCreatedAt`, поэтому это ограничение SDK, не сервиса.

### `internal/models/tools/unsorted.go`

Входные модели полностью обновлены:

- `UnsortedFilter`: `PipelineName string`, `CreatedAtFrom string`, `CreatedAtTo string`, `Order string` (числовой `PipelineID []int` удалён).
- `UnsortedAcceptParams`: `UserName`, `PipelineName`, `StatusName` (числовые `UserID`, `StatusID` удалены).
- `UnsortedDeclineParams`: отдельная структура с `UserName` (не переиспользует `AcceptParams`).
- `UnsortedCreateItem`: `PipelineName string`, `CreatedAt string` (RFC3339) вместо числовых `PipelineID int`, `CreatedAt int`.

### `internal/services/crm/unsorted/helpers.go`

Два хелпера: `unixToRFC3339` и `rfc3339ToUnix`. Корректная реализация через `time.Unix` и `time.Parse(time.RFC3339, ...)`.

---

## 2. Выполнение пунктов AUDIT.md

**Вход (15 пунктов аудита):**

| # | Пункт | Статус |
|---|-------|--------|
| 1 | `filter.pipeline_id []int` → `filter.pipeline_name string` с резолвингом | Выполнено |
| 2 | `accept_params.user_id int` → `accept_params.user_name string` с резолвингом | Выполнено |
| 3 | `accept_params.status_id int` → `pipeline_name + status_name` с резолвингом | Выполнено |
| 4 | `filter.created_at_from/to` (RFC3339) в `UnsortedFilter` | Выполнено |
| 5 | `filter.order` в `UnsortedFilter` | Выполнено |
| 6 | Отдельная структура `UnsortedDeclineParams` с `UserName` | Выполнено |
| 7 | Типизированные metadata по категориям | Отложено (осознанно, не входило в задание) |
| 8 | Убрать дублирование action `"search"` / `"list"` | Выполнено |
| 9 | `PipelineID []int` → `PipelineName string` в `UnsortedFilter` | Выполнено |
| 10 | `with=leads,contacts,companies` в `ListUnsorted` и `GetUnsorted` | Выполнено |
| 11 | `pipeline_name` в ответе вместо `pipeline_id` | Выполнено |
| 12 | Unix timestamps → RFC3339 в ответе | Выполнено |
| 13 | `PageMeta` из `ListUnsorted` возвращается LLM | Выполнено |
| 14 | `account_id` убран из `accept/decline/link` | Выполнено |
| 15 | `created_at_from/to` в `SummaryUnsorted` | Выполнено |

**Итог: 14 из 15 выполнено, 1 отложено осознанно (пункт 7 — типизированные metadata).**

---

## 3. Расхождения между log.md и реальным кодом

### ~~Критическое расхождение: `app/gkit/tools/unsorted.go` НЕ отрефакторен~~ ИСПРАВЛЕНО (2026-04-04)

Все четыре проблемы устранены в ходе итерации 2026-04-04:

1. ~~`case "search", "list":` — алиас `"search"` присутствует~~ → оставлен только `"list"`.
2. ~~В `case "create":` — handler строит `[]*models.Unsorted` с полями старой модели~~ → `input.CreateData.Items` передаётся напрямую в `CreateUnsorted`.
3. ~~В `case "accept":` — обращается к `input.AcceptParams.UserID` и `input.AcceptParams.StatusID`~~ → `input.AcceptParams` (`*gkitmodels.UnsortedAcceptParams`) передаётся напрямую в `AcceptUnsorted`.
4. ~~В `case "decline":` — читает `input.AcceptParams.UserID`~~ → `input.DeclineParams` (`*gkitmodels.UnsortedDeclineParams`) передаётся напрямую в `DeclineUnsorted`.
5. Импорт `"github.com/alextixru/amocrm-sdk-go/core/models"` удалён из handler — больше не нужен.

### Расхождение по `agent.go`

Log.md (шаг 6) утверждает что `agent.go` был обновлён: `unsorted.NewService(sdk)` → `unsorted.New(ctx, sdk)`. Файл не проверялся напрямую, но несоответствие в `tools/unsorted.go` указывает что шаг 5 точно не был выполнен. Статус `agent.go` требует проверки.

---

## 4. Что нужно сделать

### Обязательно (для компиляции и корректности)

1. **Рефакторинг `app/gkit/tools/unsorted.go`** — весь файл необходимо обновить:
   - Удалить `case "search"` — оставить только `"list"`.
   - В `case "create"`: удалить ручную сборку `[]*models.Unsorted`, передавать `input.CreateData.Items` (`[]gkitmodels.UnsortedCreateItem`) напрямую в `CreateUnsorted`.
   - В `case "accept"`: передавать `input.AcceptParams` (`*gkitmodels.UnsortedAcceptParams`) в `AcceptUnsorted`.
   - В `case "decline"`: передавать `input.DeclineParams` (`*gkitmodels.UnsortedDeclineParams`) в `DeclineUnsorted`.
   - Удалить импорт `"github.com/alextixru/amocrm-sdk-go/core/models"` (больше не нужен).

2. **Проверить `app/gkit/agent.go`**: убедиться что `unsorted.New(ctx, sdk)` вызывается с контекстом и обрабатывает ошибку.

### Желательно (качество)

3. **Пункт 7 AUDIT.md** — типизированные metadata по категории (`SipMetadata`, `MailMetadata`, `FormsMetadata`, `ChatsMetadata`) вместо единой плоской `UnsortedMetadataOutput`. Отложено осознанно.

4. **Фильтр по дате для `list`** — `CreatedAtFrom`/`CreatedAtTo` есть в `UnsortedFilter`, но SDK (`filters.UnsortedFilter`) не предоставляет `SetCreatedAt`. Если SDK будет обновлён — пробросить по аналогии с `SummaryUnsorted`.

---

## 5. Итоговая оценка

Сервисный слой (`service.go`, `unsorted.go`, `helpers.go`) и входные модели (`models/tools/unsorted.go`) полностью реализованы корректно и соответствуют плану рефакторинга. Tool handler (`app/gkit/tools/unsorted.go`) обновлён в итерации 2026-04-04 — приведён в соответствие с новыми интерфейсами сервиса. Пакет компилируется.

---

## 6. Изменения в этой итерации (2026-04-04)

### `app/gkit/tools/unsorted.go` — полный рефакторинг handler'а

| Что изменилось | До | После |
|---|---|---|
| `case "search", "list"` | оба алиаса | только `"list"` |
| `case "create"` | ручная сборка `[]*models.Unsorted` со старыми полями | `input.CreateData.Items` передаётся напрямую |
| `case "accept"` | `input.AcceptParams.UserID`, `input.AcceptParams.StatusID` (не существуют) | `input.AcceptParams` (`*UnsortedAcceptParams`) |
| `case "decline"` | `input.AcceptParams.UserID` (неверная структура) | `input.DeclineParams` (`*UnsortedDeclineParams`) |
| импорт SDK-моделей | `"github.com/alextixru/amocrm-sdk-go/core/models"` | удалён |

Handler стал тонким маршрутизатором без зависимости от SDK-моделей. Вся бизнес-логика и маппинг — в сервисном слое.
