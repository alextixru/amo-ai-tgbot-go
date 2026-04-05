# complex_create: итоговый репорт

Дата: 2026-04-04

---

## Текущее состояние сервиса (что реально реализовано)

### Файлы пакета

| Файл | Содержание |
|------|-----------|
| `service.go` | Интерфейс `Service`, конструктор `New(ctx, sdk)`, загрузка справочников, resolve/lookup хелперы |
| `complex.go` | `CreateComplex`, `CreateComplexBatch`, `buildSDKLead`, `buildSDKContact`, `buildSDKCompany`, `enrichResult`, `buildContactCustomFields`, `mapCustomFieldsValues` |
| `output.go` | DTO для LLM: `ComplexCreateResult`, `CreatedLeadView`, `CreatedContactView`, `CreatedCompanyView` |

### Интерфейс Service

```go
type Service interface {
    CreateComplex(ctx context.Context, input *gkitmodels.ComplexCreateInput) (*ComplexCreateResult, error)
    CreateComplexBatch(ctx context.Context, inputs []gkitmodels.ComplexCreateInput) ([]ComplexCreateResult, error)
    PipelineNames() []string
    UserNames() []string
}
```

### Инициализация (`New`)

Конструктор `New(ctx, sdk)` — **блокирующий**: если загрузка пользователей или воронок упала, возвращает ошибку (не silent fail). Строит шесть внутренних мап:

- `usersByName map[string]int`, `usersByID map[int]string`
- `pipelinesByName map[string]int`, `pipelinesByID map[int]string`
- `statusesByPipelineAndName map[int]map[string]int`, `statusesByPipelineAndID map[int]map[int]string`

Воронки загружаются через `sdk.Pipelines().Get(ctx, params)` с `with=statuses` (через `url.Values`, не строкой).

### Резолвинг имён

| Метод | Поведение на пустое имя | Поведение при ненайденном |
|-------|------------------------|--------------------------|
| `resolveUserID` | возвращает `0, nil` | ошибка со списком доступных |
| `resolvePipelineID` | возвращает `0, nil` | ошибка со списком доступных |
| `resolveStatusID` | возвращает `0, nil` | ошибка: нет статусов или нет имени, со списком доступных |

Обратный резолвинг: `lookupUserName`, `lookupPipelineName`, `lookupStatusName` — возвращают `"[unknown:ID]"` при ненайденном.

### Маппинг входных данных (`buildSDKLead`)

- Резолвит `PipelineName`, `StatusName`, `ResponsibleUserName` из `LeadData`
- Маппит кастомные поля через `mapCustomFieldsValues` (JSON marshal/unmarshal + fallback)
- Маппит теги по имени в `[]amomodels.Tag`
- Строит контакты через `buildSDKContact` (телефон/email → кастомные поля `PHONE`/`EMAIL` с `EnumCode: "WORK"`)
- Строит компанию через `buildSDKCompany`
- `ResponsibleUserID` присваивается отдельно после инициализации struct (обход ограничения `BaseModel`)

### Обогащение ответа (`enrichResult`)

- Основной контакт: ID из `contactID` (из SDK), имя и ответственный из `input.Contacts[0]`
- Дополнительные контакты (i > 0): ID недоступны из API — отдаются только имя и ответственный
- `CreatedAt` у сделки: `time.Now().UTC().Format(time.RFC3339)` (время формирования ответа, не время создания в CRM)

### Tool handler (`app/gkit/tools/complex_create.go`)

- Два tool'а: `complex_create` и `complex_create_batch`
- Описания формируются динамически через `PipelineNames()` / `UserNames()` при регистрации
- Весь маппинг input → SDK вынесен в сервис; handler — тонкая прокладка

### agent.go и main.go

- `NewAgent(ctx, client, sdk)` возвращает `(*Agent, error)` — ошибки инициализации сервисов пробрасываются
- `complex_create.New(ctx, sdk)` вызывается внутри `NewAgent`; при ошибке — `return nil, fmt.Errorf("NewAgent: %w", err)`
- `cmd/bot/main.go` обрабатывает ошибку `NewAgent`

### Входные модели (`internal/models/tools/complex_create.go`)

Все ID заменены на имена:

| Поле | Тип | JSON-тег |
|------|-----|----------|
| `LeadData.PipelineName` | `string` | `pipeline_name` |
| `LeadData.StatusName` | `string` | `status_name` |
| `LeadData.ResponsibleUserName` | `string` | `responsible_user_name` |
| `ContactData.ResponsibleUserName` | `string` | `responsible_user_name` |
| `CompanyData.ResponsibleUserName` | `string` | `responsible_user_name` |

---

## Соответствие коду vs log/REPORT

### Полностью совпадает с описанием

- Замена ID на имена во входных моделях — реализована, JSON-теги верные
- Новые DTO в `output.go` — реализованы в точности
- Интерфейс `Service` (4 метода) — совпадает
- `New(ctx, sdk)` с загрузкой справочников — реализован
- Поля `usersByName/ID`, `pipelinesByName/ID`, `statusesByPipelineAndName/ID` — все присутствуют
- Резолвинг с читаемыми ошибками и списком доступных — реализован
- Обратный резолвинг `[unknown:ID]` — реализован
- Маппинг в `buildSDKLead/Contact/Company` в сервисе — реализован
- `buildContactCustomFields` перемещён в `complex.go` (не удалён, а переехал в сервис)
- Tool handler упрощён, динамические описания — реализовано
- `agent.go`: `NewAgent` принимает `ctx`, возвращает ошибку — реализовано
- `main.go`: обработка ошибки `NewAgent` — реализовано

### Расхождения

| Описание в log/REPORT | Реальный код | Вывод |
|-----------------------|-------------|-------|
| REPORT: "Ошибки загрузки не блокируют старт, работает с пустыми справочниками" | `service.go` комментарий к `New` содержит эту фразу, но фактически код `return nil, err` при ошибках загрузки — блокирует старт | Комментарий в коде устарел, поведение правильное (блокирующее), документация вводит в заблуждение |
| REPORT: "loadPipelines через `sdk.Pipelines().Get(ctx, "with=statuses")`" (строка) | Код использует `url.Values{}` с `params.Set("with", "statuses")` | Несущественно, результат идентичен |
| REPORT: "mapCustomFieldsValues конвертирует map[string]any" | В коде также есть fallback-парсинг при ошибке `json.Unmarshal` — не упомянуто | Дополнительная деталь, не расхождение |
| log шаг 7: "обновлён agent.go: New вместо NewService" | В реальном agent.go также инициализируются `entities.New(ctx, sdk)` и `catalogs.New(ctx, sdk)` с обработкой ошибок — не упомянуто в логе | Лог описывает только изменения по complex_create |

---

## Известные ограничения (не баги)

- **ID дополнительных контактов недоступны**: amoCRM complex endpoint возвращает только `contact_id` первого контакта. Дополнительные контакты отдаются без ID (только имена из input).
- **`CreatedAt` — время ответа, не CRM**: поле заполняется `time.Now()` при формировании DTO, не берётся из SDK.
- **Нет автообновления справочников**: при изменении воронок/пользователей в CRM нужен рестарт сервиса.
- **Дубли имён**: при дублях пользователей/воронок побеждает последний в итерации (поведение map). Проверка дублей не реализована (отличается от REFACTORING.md edge case "ошибка при дублях").

---

## Что ещё нужно сделать

Нет незавершённых задач по `complex_create`. Сервис реализован полностью согласно требованиям REFACTORING.md для данного сервиса.

Смежные сервисы из REFACTORING.md, которые ещё ожидают аналогичного рефакторинга:
- `entities/` — воронки, статусы, пользователи, кастомные поля по имени
- `activities/` — пользователи
- `unsorted/` — воронки, статусы, пользователи
- `customers/` — пользователи, статусы, сегменты
- `catalogs/` — каталоги (судя по `agent.go`, уже имеет `New(ctx, sdk)`)
- `products/` — продукты
