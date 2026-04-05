# admin_schema — итоговый отчёт

_Дата: 2026-04-04_

---

## Текущее состояние сервиса (реальный код)

### Слои и операции

Сервис управляет четырьмя слоями схемы CRM через единый Genkit-tool `admin_schema`:

| Слой | list | get | create | update | delete |
|------|------|-----|--------|--------|--------|
| `custom_fields` | + | + | + | + | + |
| `field_groups` | + | + | + | + | + |
| `loss_reasons` | + | + | + | — (API ограничение) | + |
| `sources` | + | + | + | + | + |

### Архитектура

```
AdminSchemaInput (typed)
    └── app/gkit/tools/admin_schema.go   — tool handler, конвертеры, client-side фильтры
        └── internal/services/crm/admin_schema/service.go  — Service interface
            ├── custom_fields.go
            ├── field_groups.go
            ├── loss_reasons.go
            └── sources.go
                └── amocrm-sdk-go
```

Входная модель: `internal/models/tools/admin_schema.go`

### Интерфейс Service (service.go)

Все list/create/update-методы возвращают `*PagedResult[T]`, delete-методы — `*DeleteResult`.

```go
type PagedResult[T any] struct {
    Items      []T
    TotalItems int
    Page       int
    HasMore    bool
    Meta       *services.PageMeta  // не сериализуется в JSON
}

type DeleteResult struct {
    Success   bool
    DeletedID any
}
```

### Входная модель (models/tools/admin_schema.go)

`Data map[string]any` полностью удалён. Вместо него — типизированные структуры:

- `CustomFieldData` — id, name, type, code, sort, group_id, is_api_only, is_required, enums[]
- `CustomFieldEnumData` — id, value, sort, code
- `FieldGroupData` — id, name, sort
- `LossReasonData` — name, sort
- `SourceData` — id, name, external_id, origin_code, pipeline_id, default

`AdminSchemaInput` принимает одиночные (`custom_field`, `field_group`, `loss_reason`, `source`) и batch-варианты (`custom_fields[]`, `field_groups[]`, `loss_reasons[]`, `sources[]`).

`SchemaFilter` содержит: `limit`, `page`, `name` (client-side), `ids[]`, `types[]`, `order map[string]string`, `external_ids[]`.

### Tool handler (app/gkit/tools/admin_schema.go)

- `buildCustomFieldsFilter` — строит `*filters.CustomFieldsFilter` с order через `SetOrder()`
- `buildSourcesFilter` — строит `*filters.SourcesFilter` с `ExternalIDs`
- `buildBaseFilter` — строит `url.Values` с limit/page для field_groups и loss_reasons
- `filterByName[T]()` — generic client-side фильтрация по имени (частичное совпадение, case-insensitive)
- `customFieldDataToModel`, `fieldGroupDataToModel`, `lossReasonDataToModel`, `sourceDataToModel` — конвертеры Data → SDK-модели
- `collectCustomFields`, `collectFieldGroups`, `collectLossReasons`, `collectSources` — объединяют одиночные и batch-поля из input

---

## Что было сделано (из log.md) — сверка с кодом

| Пункт из лога | Статус в коде |
|---------------|---------------|
| `with=fields` в `ListFieldGroups` | **Реализовано.** `field_groups.go:15`: `filter.Set("with", "fields")` |
| `with=fields` в `GetFieldGroup` | **Реализовано.** `field_groups.go:25`: `params.Set("with", "fields")` |
| `PagedResult[T]` и `DeleteResult` в service.go | **Реализовано.** `service.go:14-26` |
| `newPagedResult()` helper | **Реализовано.** `service.go:70-84` |
| Все list/create/update → `*PagedResult[T]` | **Реализовано.** Все 4 файла |
| Все delete → `*DeleteResult` | **Реализовано.** Все 4 файла |
| Удалён `Data map[string]any`, добавлены типизированные структуры | **Реализовано.** `models/tools/admin_schema.go` |
| `Name` в `SchemaFilter` | **Реализовано.** `models/tools/admin_schema.go:54` |
| `Order` в `SchemaFilter` + `buildCustomFieldsFilter` | **Реализовано.** `models/tools/admin_schema.go:61`, `tools/admin_schema.go:57-59` |
| `filterByName[T]()` в tool handler | **Реализовано.** `tools/admin_schema.go:90-102` |
| Конвертеры Data → SDK-модели | **Реализовано.** `tools/admin_schema.go:104-154` |
| Алиас action `"search"` убран | **Реализовано.** Только `"list"` во всех switch-ветках |
| Убран `encoding/json` из tool handler | **Реализовано.** Импорт отсутствует |

**Все задокументированные в log.md изменения совпадают с реальным кодом. Расхождений нет.**

---

## Расхождения между AUDIT.md и реальным кодом

Следующие пункты из AUDIT.md **не реализованы**, что соответствует явному решению, задокументированному в REPORT.md и log.md:

| Пункт AUDIT.md | Причина нереализации |
|----------------|---------------------|
| `pipeline_name` в ответе Source (СРЕДНИЙ) | Резолвинг требует загрузки воронок при init. `admin_schema` согласно `REFACTORING.md` не входит в список сервисов, требующих рефакторинга (уже работает с именами). Намеренно отложено. |
| `RequiredStatus` с `pipeline_name`/`status_name` в ответе (НИЗКИЙ) | Та же причина. |
| Убрать `AccountID` и `Links` из ответа (НИЗКИЙ) | SDK-модели возвращаются напрямую, без выходного маппера. Не критично для LLM. |
| Фильтр по `code` в `SchemaFilter` (НИЗКИЙ) | Низкий приоритет, не реализован. Добавляется аналогично `Name`-фильтру. |
| `With` параметры через `SchemaFilter` (НИЗКИЙ из AUDIT) | Частично: `with=fields` захардкожен в `field_groups.go`. Для custom_fields `With` пробрасывается через `CustomFieldsFilter`. |

Один оставшийся баг из AUDIT, **не исправленный**:
- `GetCustomField` по-прежнему передаёт `url.Values{}` без `With`. Если LLM когда-либо потребуются `with`-параметры для одиночного поля — это придётся добавить отдельно. Сейчас `GetCustomField` возвращает корректные данные, так как SDK заполняет все базовые поля без `with`.

---

## Что ещё можно сделать (если применимо)

Все пункты имеют низкий приоритет и не блокируют работу LLM:

1. **Фильтр по `code` в `SchemaFilter`** — аналогично `name`, client-side, одна строка в `filterByName`.
2. **Убрать `AccountID` и `Links` из ответа** — потребует выходного DTO/маппера, сейчас возвращаются сырые SDK-модели.
3. **`pipeline_name` в ответе Source** — резолвинг через `admin_pipelines` при инициализации (out of scope согласно REFACTORING.md).
4. **`RequiredStatus` с именами воронки/статуса** — аналогично предыдущему.
5. **`With` параметры для `GetCustomField`** — если потребуется, добавить `with` в `url.Values` аналогично `GetFieldGroup`.

---

## Файлы сервиса

| Файл | Содержимое |
|------|-----------|
| `service.go` | Интерфейс `Service`, `PagedResult[T]`, `DeleteResult`, `newPagedResult()`, `NewService()` |
| `custom_fields.go` | 5 методов: List/Get/Create/Update/Delete для кастомных полей |
| `field_groups.go` | 5 методов: List/Get/Create/Update/Delete для групп полей, `with=fields` |
| `loss_reasons.go` | 4 метода: List/Get/Create/Delete (нет Update — API ограничение) |
| `sources.go` | 5 методов: List/Get/Create/Update/Delete для источников |
| `internal/models/tools/admin_schema.go` | Типизированные входные структуры, `SchemaFilter` |
| `app/gkit/tools/admin_schema.go` | Tool handler, конвертеры, фильтры, `filterByName[T]()` |
