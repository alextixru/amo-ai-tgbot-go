# Репорт: entities — текущее состояние (2026-04-04)

## Что реализовано

Сервис полностью переписан и соответствует архитектуре из REFACTORING.md.

### Инициализация (`service.go`)

`New(ctx, sdk)` загружает справочники при старте:

- Пользователи — `loadUsers` → `usersByName / usersByID`
- Воронки + статусы — `loadPipelines` → `pipelinesByName/ID`, `statusesByPipelineAndName`, `statusesByID`, `statusPipelineByID`
- Кастомные поля — `loadCustomFields` → `customFieldsLeads / Contacts / Companies` (code → id), нефатально при ошибке
- Причины отказа — `loadLossReasons` → `lossReasonsByName/ID`, нефатально при ошибке

Методы резолвинга (имя → ID): `resolvePipelineID`, `resolveUserID`, `resolveUserIDs`, `resolveStatusID`, `resolveLossReasonID`.
Методы обратного резолвинга (ID → имя): `lookupUserName`, `lookupPipelineName`, `lookupStatusName`, `lookupLossReasonName`.
При неизвестном ID возвращается `[unknown:N]`, ответ не ломается.

Метаданные для описания tools: `PipelineNames() []string`, `UserNames() []string`.

### Интерфейс Service

Покрывает сделки, контакты, компании, включая `Sync*` и `Link/Unlink`. Дополнительно: `GetContactChats`, `LinkContactChats`.

```
SearchLeads / GetLead / CreateLead / CreateLeads / UpdateLead / UpdateLeads / SyncLead / LinkLead / UnlinkLead
SearchContacts / GetContact / CreateContact / CreateContacts / UpdateContact / UpdateContacts / SyncContact / LinkContact / UnlinkContact / GetContactChats / LinkContactChats
SearchCompanies / GetCompany / CreateCompany / CreateCompanies / UpdateCompany / UpdateCompanies / SyncCompany / LinkCompany / UnlinkCompany
```

### Входная модель (`internal/models/tools/entities.go`)

Числовые ID заменены на имена во всех полях:

| Поле | Тип |
|------|-----|
| `EntitiesFilter.ResponsibleUserNames` | `[]string` |
| `EntitiesFilter.CreatedByNames` | `[]string` |
| `EntitiesFilter.UpdatedByNames` | `[]string` |
| `EntitiesFilter.PipelineNames` | `[]string` |
| `EntitiesFilter.Statuses` | `[]StatusPair{PipelineName, StatusName}` |
| `EntitiesFilter.CustomFieldsValues` | `[]CustomFieldFilter{FieldCode, Values}` |
| `EntitiesFilter.Created/Updated/ClosedAt` | ISO-8601 string |
| `EntityData.PipelineName` | `string` |
| `EntityData.StatusName` | `string` |
| `EntityData.ResponsibleUserName` | `string` |
| `EntityData.LossReasonName` | `string` |
| `EntityData.SourceName` | `string` |
| `EntityData.FirstName / LastName` | `string` (contacts) |

### Маппинг фильтров (`leads.go`, `contacts.go`, `companies.go`)

Все ранее отсутствующие вызовы реализованы:

- `SetStatuses` — вызывается, пары `{PipelineName, StatusName}` резолвятся в `{PipelineID, StatusID}` и передаются в SDK
- `SetCreatedBy` / `SetUpdatedBy` — вызываются для leads, contacts, companies
- `SetCustomFieldsValues` — вызывается для всех трёх типов через `buildCustomFieldsFilter` (field_code → field_id)
- `SetPipelineIDs` — имена воронок резолвятся в ID
- `SetResponsibleUserIDs` — имена резолвятся в ID
- `SetPrice`, `SetCreatedAt`, `SetUpdatedAt`, `SetClosedAt` — реализованы

Даты принимаются в ISO-8601 и конвертируются в Unix через `parseISO` (поддерживает RFC3339, `2006-01-02T15:04:05`, `2006-01-02`).

### Маппинг входных данных (`mapToLead`, `mapToContact`, `mapToCompany`)

- `mapToLead` — резолвит PipelineName, StatusName, ResponsibleUserName, LossReasonName; заполняет EmbeddedContacts и EmbeddedCompanies
- `mapToContact` — резолвит ResponsibleUserName; заполняет FirstName/LastName; реализует EmbeddedCompanies (ранее был баг — только комментарий)
- `mapToCompany` — резолвит ResponsibleUserName

### Выходная модель (`response.go`)

`EntityResult` — read-friendly структура без числовых ID:

```go
type EntityResult struct {
    ID, Name, Price, PipelineName, StatusName, LossReason, SourceName, ClosedAt
    FirstName, LastName
    ResponsibleUserName, CreatedByName, UpdatedByName
    CreatedAt, UpdatedAt  // ISO-8601
    Tags                  []string
    CustomFieldsValues    []CustomFieldEntry{FieldCode, FieldName, Values}
    Contacts, Companies, Leads []EntityRef{ID, Name}
}
```

`SearchResult` — обёртка с `Items []*EntityResult`, `HasMore bool`, `Page int`.
`LinkResult` — `{Success bool, Message string}` — возвращается из всех link/unlink операций.

Технические поля `account_id`, `_links` убраны.

### Маппинг ответов (`mapping.go`)

- `leadToResult` — заполняет все имена; обрабатывает `Embedded.LossReason`, `Embedded.Source`, Tags, Contacts, Companies
- `contactToResult` — заполняет FirstName/LastName, связанные Leads, Companies (включая `Embedded.Company` и `Embedded.Companies`)
- `companyToResult` — заполняет Leads, Contacts
- `mapCFVToEntries` — конвертирует `[]CustomFieldValue` в `[]CustomFieldEntry` (FieldCode, FieldName, Values)
- `mapCustomFieldsValues` — конвертирует `map[string]any` в `[]CustomFieldValue` для SDK; поддерживает строку, `[]any`, `map[string]any{value, enum_code}`

### Дефолтные `with` для GetLead

`GetLead` автоматически добавляет `["loss_reason", "source"]` если их нет в переданном with — LLM получает обогащённый ответ без явного запроса.

### Tool-слой (`app/gkit/tools/entities.go`)

Тонкий диспетчер, маршрутизирует по `entity_type` и `action`. Все операции возвращают конкретные типы (`*EntityResult`, `*SearchResult`, `*LinkResult`), а не сырые SDK-модели. Link/Unlink возвращают `*LinkResult` вместо `nil`.

---

## Что было исправлено (баги из AUDIT.md)

| Баг | Статус |
|-----|--------|
| `SetStatuses` не вызывался — `StatusID` фильтр полностью игнорировался | Исправлен |
| `SetCreatedBy` / `SetUpdatedBy` не вызывались в leads.go | Исправлен |
| `SetCustomFieldsValues` не вызывался нигде | Исправлен |
| `mapToContact`: блок `EmbeddedCompanies` был только комментарием | Исправлен |
| `CreateContact` / `CreateCompany` возвращали слайс вместо одного объекта | Исправлен |
| `Link/Unlink` возвращали `nil` — LLM не получала подтверждения | Исправлен |
| `PageMeta` отбрасывалась — LLM не знала о неполноте результатов | Исправлен |

---

## Соответствие log.md реальному коду

Все 9 шагов из log.md выполнены и соответствуют коду:

| Шаг из log.md | Статус |
|---------------|--------|
| 1. Чтение файлов и анализ | Артефакт (документация) |
| 2. Определение стратегии | Артефакт (документация) |
| 3. Изменение models/tools/entities.go | Подтверждён кодом |
| 4. Добавление read-model в response.go | Подтверждён кодом |
| 5. Переписан service.go | Подтверждён кодом |
| 6. Переписан leads.go | Подтверждён кодом |
| 7. Переписаны contacts.go и companies.go | Подтверждён кодом |
| 8. Обновлён mapping.go | Подтверждён кодом |
| 9. Обновлён app/gkit/tools/entities.go | Подтверждён кодом |

---

## Расхождения между AUDIT.md/log.md и реальным кодом

Расхождений нет. Все пункты AUDIT.md с пометкой "Список конкретных изменений" реализованы за исключением намеренно не включённых:

- **`GetLeadLinks` / `GetContactLinks` / `GetCompanyLinks`** (п. 19 AUDIT.md) — в интерфейс не добавлены. SDK-методы `GetLinks` существуют, но tool-слой их не представляет.
- **`SourceName` в `EntityData`** — поле есть в структуре, но `mapToLead` его не использует (SDK `Lead.SourceID` существует, резолвинг не реализован, так как sources не загружаются при init).

---

## Что остаётся сделать

### Незначительное (entities/)

1. **Резолвинг `SourceName` на входе** — сейчас поле `EntityData.SourceName` объявлено, но не обрабатывается в `mapToLead`. Для полноты нужно добавить загрузку `Sources` при init и резолвинг `source_name → source_id`.

2. **`GetLeadLinks` / `GetContactLinks` / `GetCompanyLinks`** — SDK-методы `GetLinks` есть, в интерфейсе отсутствуют.

3. **Обновление справочников без рестарта** — текущий MVP загружает данные один раз. При добавлении новых воронок/пользователей нужен рестарт.

### Рефакторинг других сервисов (из REFACTORING.md)

Entities — единственный завершённый сервис. Остальные по-прежнему используют числовые ID:

| Сервис | Что нужно |
|--------|-----------|
| `complex_create/` | Загрузка воронок, пользователей; замена ID → Names в LeadData, ContactData, CompanyData |
| `activities/` | Загрузка пользователей; замена `ResponsibleUserID`, `UserIDs`, `UserID` на имена |
| `unsorted/` | Загрузка воронок, пользователей; замена в AcceptParams, DeclineParams, UnsortedFilter |
| `customers/` | Загрузка пользователей, статусов покупателей, сегментов; замена в CustomerData, CustomersFilter |
| `catalogs/` | Загрузка каталогов; замена CatalogID, ElementID, EntityID на имена |
| `products/` | Загрузка каталога продуктов; замена ProductID, EntityID, PriceID на имена |
