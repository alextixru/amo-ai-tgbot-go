# Репорт: admin_pipelines

_Дата: 2026-04-04_

---

## Текущее состояние сервиса

Сервис **полностью отрефакторен** на уровне `service.go`, `pipelines.go`, `statuses.go` и `internal/models/tools/admin_pipelines.go`. Слой `app/gkit/tools/admin_pipelines.go` **обновлён** в итерации от 2026-04-04.

---

## Что реально реализовано в коде

### internal/models/tools/admin_pipelines.go

- Типизированные структуры `PipelineData` и `StatusData` (вместо `map[string]any`)
- `StatusData.Type string` — принимает `"regular"`, `"won"`, `"lost"` (вместо `int`)
- `StatusData.Color` — задокументированы все 21 допустимый hex-цвет из `StatusColors`
- `AdminPipelinesInput` содержит: `PipelineID`, `PipelineName`, `StatusID`, `StatusName`, `WithStatuses bool`
- Батч-режим: `Pipeline *PipelineData` (один), `Status *StatusData` (один), `Items []map[string]any` (множество)
- Output-типы: `PipelineOutput`, `StatusOutput`, `ListPipelinesOutput`
- `PipelineOutput`: не содержит `account_id` и `_links`, имеет `IsArchive`, вложенные `Statuses`
- `StatusOutput`: содержит `TypeLabel string`, `IsWon bool`, `IsLost bool`, `IsClosed bool` — без `account_id`
- `ListPipelinesOutput`: содержит `Pipelines` и `PageMeta`

### internal/services/crm/admin_pipelines/service.go

- Интерфейс `Service` с полными сигнатурами (все методы принимают `name string` как альтернативу ID)
- `ListPipelines(ctx, withStatuses bool) (*ListPipelinesOutput, error)`
- `GetPipeline(ctx, id, name, withStatuses bool) (*PipelineOutput, error)`
- `CreatePipelines(ctx, []PipelineData) ([]*PipelineOutput, error)`
- `UpdatePipeline(ctx, id, name, PipelineData) (*PipelineOutput, error)`
- `DeletePipeline(ctx, id, name) error`
- `ListStatuses(ctx, pipelineID, pipelineName) ([]*StatusOutput, error)`
- `GetStatus(ctx, pipelineID, pipelineName, statusID, statusName) (*StatusOutput, error)`
- `CreateStatus/CreateStatuses/UpdateStatus/DeleteStatus` — аналогично
- Helpers: `resolvePipelineID`, `resolveStatusID` — lookup по имени с сообщением об ошибке, включающим список доступных значений

### internal/services/crm/admin_pipelines/pipelines.go

- `ListPipelines` — передаёт `?with=statuses` когда `withStatuses=true`, возвращает `*ListPipelinesOutput` с `PageMeta`
- `GetPipeline` — использует `sdkservices.WithRelations("statuses")` когда `withStatuses=true`
- `pipelineToOutput` — маппит в `PipelineOutput`, включая вложенные статусы, без `account_id`/`_links`

### internal/services/crm/admin_pipelines/statuses.go

- `GetStatus` — передаёт `?with=descriptions`
- `statusToOutput` — маппит в `StatusOutput` с `TypeLabel`, `IsWon`, `IsLost`, `IsClosed`
- `statusTypeToInt` — конвертирует `"regular"/"won"/"lost"` → числовой тип SDK
- `statusDataToModel` — конвертирует `StatusData` в SDK-модель

### app/gkit/tools/admin_pipelines.go ✓ ИСПРАВЛЕН

- Убран импорт `amomodels` — весь маппинг в сервисе
- Все вызовы методов обновлены под новые сигнатуры
- `list`: передаёт `input.WithStatuses`
- `get`: передаёт `(ctx, input.PipelineID, input.PipelineName, input.WithStatuses)`
- `create`: батч из `input.Items` → `[]PipelineData` или одиночный из `input.Pipeline`
- `update`: передаёт `(ctx, input.PipelineID, input.PipelineName, *input.Pipeline)`
- `delete`: передаёт `(ctx, input.PipelineID, input.PipelineName)`
- `list_statuses`: передаёт `(ctx, input.PipelineID, input.PipelineName)`
- `get_status`: передаёт все 4 параметра идентификации
- `create_status`: батч из `input.Items` → `[]StatusData` или одиночный из `input.Status`
- `update_status`: передаёт `(ctx, pipelineID, pipelineName, statusID, statusName, *input.Status)`
- `delete_status`: передаёт `(ctx, pipelineID, pipelineName, statusID, statusName)`
- Описание tool обновлено — убраны `data.pipelines`/`data.statuses`, описан новый формат

---

## Соответствие log.md и реального кода

| Шаг | Что описано | Соответствует коду |
|-----|-------------|-------------------|
| Шаг 2 | Замена `Data map[string]any` на `PipelineData`/`StatusData`, добавление `PipelineName`, `StatusName`, `WithStatuses`, унификация батч-режима через `Items` | Да — модели полностью соответствуют |
| Шаг 3 | Обновление интерфейса `Service` с новыми сигнатурами и helpers | Да — `service.go` полностью соответствует |
| Шаг 4 | `ListPipelines` с `with=statuses`, `GetPipeline` с `WithRelations`, `resolvePipelineID` | Да — `pipelines.go` полностью соответствует |
| Шаг 5 | `GetStatus` с `with=descriptions`, `resolveStatusID` | Да — `statuses.go` полностью соответствует |
| Шаг 6 | Обновление `app/gkit/tools/admin_pipelines.go` под новые типы | **Да — исправлено 2026-04-04** |
| Шаг 7 | Создание output-типов с семантическими полями | Да — модели полностью соответствуют |

---

## Изменения в этой итерации (2026-04-04)

### Исправлено: app/gkit/tools/admin_pipelines.go

**Было:** старый handler, не компилирующийся с новым интерфейсом сервиса.

**Стало:** handler полностью переписан под новые сигнатуры:

1. Убран импорт `amomodels` — маппинг SDK-моделей теперь исключительно в сервисе
2. `list`: добавлен `input.WithStatuses`
3. `get`: добавлены `input.PipelineName` и `input.WithStatuses`
4. `create`: читает `input.Items` (батч) или `input.Pipeline` (одиночный), маппит в `[]toolmodels.PipelineData` через JSON round-trip
5. `update`: передаёт `input.PipelineID`, `input.PipelineName`, `*input.Pipeline`
6. `delete`: передаёт `input.PipelineID`, `input.PipelineName`
7. `list_statuses`: передаёт оба поля идентификации воронки
8. `get_status`: передаёт все 4 параметра идентификации
9. `create_status`: батч из `input.Items` → `[]StatusData` через `CreateStatuses`; одиночный из `input.Status` через `CreateStatus`
10. `update_status`: передаёт все 5 параметров (pipelineID, pipelineName, statusID, statusName, data)
11. `delete_status`: передаёт все 4 параметра идентификации
12. Описание tool обновлено: убраны `data.pipelines`/`data.statuses`, добавлено описание нового формата

---

## Что нужно сделать

### Опционально (улучшения)

1. **Строгая типизация батча** — заменить `Items []map[string]any` на отдельные поля `Pipelines []PipelineData` и `Statuses []StatusData` в `AdminPipelinesInput`, устранив неоднозначность и JSON round-trip при десериализации

2. **`REFACTORING.md` пометить admin_pipelines как "готов"** — сервис и tool handler полностью отрефакторены
