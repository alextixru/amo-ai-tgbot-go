# Рефакторинг crm/ — полноценный адаптерный слой

## Проблема

CRM-сервисы сейчас — тонкие прокладки с маппингом 1:1.

- **На входе:** LLM передаёт числовые ID (pipeline_id, status_id) — она их не знает, вынуждена делать цепочку вспомогательных вызовов
- **На выходе:** сервисы возвращают сырые SDK-модели (`*models.Lead`) с числовыми ID — LLM не понимает что значит `pipeline_id: 7628434`
- **При инициализации:** описания tools статичные, LLM не знает какие воронки/поля/пользователи есть в аккаунте

## Цель

CRM-сервисы — полноценный двусторонний адаптер между языком нейросети и языком SDK:

```
          ВХОД (Tools → SDK)
Tools → EntityData{PipelineName: "Продажи"} → crm/ → sdk.Lead{PipelineID: 123}

          ВЫХОД (SDK → Tools/LLM)
sdk.Lead{PipelineID: 123} → crm/ → читаемый ответ{Воронка: "Продажи", Этап: "Новая заявка"}

          ИНИЦИАЛИЗАЦИЯ (SDK → Tool descriptions)
sdk → crm/ загружает справочники → tools получают описания с реальными значениями аккаунта
```

---

## Три направления рефакторинга

### Направление 1: Вход — имена вместо ID

LLM передаёт строковые имена, сервисы резолвят в ID для SDK.

### Направление 2: Выход — читаемые ответы

Сервисы возвращают не сырые SDK-модели, а обогащённые DTO где ID заменены/дополнены именами.

### Направление 3: Инициализация — метаданные для tools

Resolver отдаёт списки доступных значений (воронки, этапы, пользователи, поля) для формирования динамических описаний tools.

---

## Направление 1: Вход

### Инвентаризация: где LLM передаёт числовые ID

#### models/tools/entities.go — EntityData

| Поле | Сейчас | После |
|------|--------|-------|
| `PipelineID` | int | `PipelineName` string |
| `StatusID` | int | `StatusName` string |
| `ResponsibleUserID` | int | `ResponsibleUserName` string |
| `CustomFieldsValues` | map[string]any (ключ = field_id) | ключ = название поля |

#### models/tools/entities.go — EntitiesFilter

| Поле | Сейчас | После |
|------|--------|-------|
| `PipelineID` | []int | `PipelineNames` []string |
| `StatusID` | []int | `StatusNames` []string |
| `ResponsibleUserID` | []int | `ResponsibleUserNames` []string |
| `CreatedBy` | []int | `CreatedByNames` []string |
| `UpdatedBy` | []int | `UpdatedByNames` []string |
| `CustomFieldsValues[].FieldID` | int | `FieldName` string |

#### models/tools/activities.go — TaskData

| Поле | Сейчас | После |
|------|--------|-------|
| `ResponsibleUserID` | int | `ResponsibleUserName` string |
| `TaskTypeID` | int | оставить (1=звонок, 2=встреча — фиксированные) |

#### models/tools/activities.go — TasksFilter

| Поле | Сейчас | После |
|------|--------|-------|
| `ResponsibleUserID` | []int | `ResponsibleUserNames` []string |

#### models/tools/activities.go — ActivitiesInput

| Поле | Сейчас | После |
|------|--------|-------|
| `UserIDs` | []int | `UserNames` []string |
| `UserID` | int | `UserName` string |

#### models/tools/complex_create.go — LeadData, ContactData, CompanyData

| Поле | Сейчас | После |
|------|--------|-------|
| `PipelineID` | int | `PipelineName` string |
| `StatusID` | int | `StatusName` string |
| `ResponsibleUserID` | int | `ResponsibleUserName` string |

---

## Направление 2: Выход

### Проблема сейчас

Сервисы возвращают `*models.Lead`, `[]*models.Contact` — сырые SDK-модели. Genkit сериализует через `json.Marshal`, LLM получает:

```json
{
  "id": 123,
  "name": "Сделка",
  "pipeline_id": 7628434,
  "status_id": 48953210,
  "responsible_user_id": 9182736,
  "custom_fields_values": [{"field_id": 284756, "values": [{"value": "500000"}]}]
}
```

LLM не знает что `7628434` — это "Продажи".

### Решение

Сервисы возвращают свои DTO с человекочитаемыми данными. Resolver обогащает ответ: `id → name`.

```json
{
  "id": 123,
  "name": "Сделка",
  "pipeline": "Продажи",
  "status": "Новая заявка",
  "responsible_user": "Иван Петров",
  "custom_fields": {"Бюджет": "500000", "Источник": "Сайт"}
}
```

### Что нужно

Выходные DTO в `models/tools/` (или рядом) — обогащённые структуры для ответов LLM.

Маппинг SDK → DTO в сервисах: `mapLeadToResponse(lead *models.Lead) *LeadResponse`.

### Какие сервисы затронуты

| Сервис | Возвращает SDK-модель | Нужна выходная DTO |
|--------|----------------------|-------------------|
| `entities/` | Lead, Contact, Company | да — основные сущности с воронками, статусами, полями |
| `activities/` | Task, Note, Call, Event, ... | частично — Task имеет ResponsibleUserID |
| `complex_create/` | ComplexLeadResult | да — содержит Lead |
| `admin_pipelines/` | Pipeline, Status | нет — возвращает уже с именами |
| `admin_users/` | User, Role | нет — уже читаемые |
| `admin_schema/` | CustomField, Source, LossReason | нет — справочники уже читаемые |
| `catalogs/` | Catalog, CatalogElement | возможно позже |
| `customers/` | Customer | возможно позже |
| `files/` | File | нет — UUID-based |
| `unsorted/` | Unsorted | возможно позже |
| `products/` | CatalogElement | возможно позже |
| `admin_integrations/` | Webhook, Widget, ... | нет |

**Приоритет:** entities, activities (tasks), complex_create. Остальное — позже.

---

## Направление 3: Инициализация

### Что нужно

Resolver хранит справочники и отдаёт их в двух форматах:

1. **Карты имя → ID** — для резолвинга (направление 1)
2. **Списки имён** — для описаний tools

### Пример использования в tool descriptions

Сейчас (статичное описание):
```
"ID воронки (только для leads)"
```

После (динамическое):
```
"Название воронки. Доступные: Продажи, VIP, Партнёры"
```

### API Resolver для инициализации

```go
func (r *Resolver) PipelineNames() []string
func (r *Resolver) StatusNames(pipelineName string) []string
func (r *Resolver) UserNames() []string
func (r *Resolver) CustomFieldNames(entityType string) []string
func (r *Resolver) SourceNames() []string
func (r *Resolver) LossReasonNames() []string
```

---

## Resolver — общий компонент

Файл: `internal/services/crm/resolver.go`

```go
type Resolver struct {
    pipelines    map[string]int                  // name → id
    statuses     map[string]map[string]int       // pipeline_name → status_name → id
    users        map[string]int                  // name → id
    customFields map[string]map[string]int       // entity_type → field_name → id
    sources      map[string]int                  // name → id
    lossReasons  map[string]int                  // name → id
}

// Резолвинг (вход)
func (r *Resolver) PipelineID(name string) (int, error)
func (r *Resolver) StatusID(pipelineName, statusName string) (int, error)
func (r *Resolver) UserID(name string) (int, error)
func (r *Resolver) CustomFieldID(entityType, fieldName string) (int, error)
func (r *Resolver) SourceID(name string) (int, error)
func (r *Resolver) LossReasonID(name string) (int, error)

// Обратный резолвинг (выход)
func (r *Resolver) PipelineName(id int) string
func (r *Resolver) StatusName(pipelineID, statusID int) string
func (r *Resolver) UserName(id int) string
func (r *Resolver) CustomFieldName(entityType string, fieldID int) string

// Списки для описаний tools (инициализация)
func (r *Resolver) PipelineNames() []string
func (r *Resolver) StatusNames(pipelineName string) []string
func (r *Resolver) UserNames() []string
func (r *Resolver) CustomFieldNames(entityType string) []string
```

Файл: `internal/services/crm/loader.go`

```go
func LoadResolver(ctx context.Context, sdk *amocrm.SDK) (*Resolver, error)
```

Загружает из SDK:

| Данные | Вызов SDK | Карта |
|--------|-----------|-------|
| Воронки | `sdk.Pipelines().Get()` | name ↔ id |
| Статусы | `pipeline.Embedded.Statuses` | (pipeline, status) ↔ id |
| Пользователи | `sdk.Users().Get()` | name ↔ id |
| Кастомные поля | `sdk.CustomFields().Get(entityType)` × 3 | (entity_type, field_name) ↔ id |
| Источники | `sdk.Sources().Get()` | name ↔ id |
| Причины отказа | `sdk.LossReasons().Get()` | name ↔ id |

---

## Порядок работы

### Фаза 1: Resolver
1. `crm/resolver.go` — структура, методы резолвинга (вход + выход), списки для описаний
2. `crm/loader.go` — загрузка справочников из SDK

### Фаза 2: Входные модели
3. `models/tools/entities.go` — EntityData, EntitiesFilter: ID → Name
4. `models/tools/activities.go` — TaskData, ActivitiesInput: ID → Name
5. `models/tools/complex_create.go` — LeadData, ContactData, CompanyData: ID → Name

### Фаза 3: Выходные модели
6. Выходные DTO для entities (LeadResponse, ContactResponse, CompanyResponse)
7. Выходная DTO для tasks (TaskResponse)

### Фаза 4: CRM-сервисы
8. `crm/entities/` — конструктор с resolver, маппинг вход (имя→ID) + выход (ID→имя)
9. `crm/activities/` — маппинг задач
10. `crm/complex_create/` — маппинг

### Фаза 5: Интеграция
11. `app/gkit/agent.go` — инициализация resolver, передача в сервисы
12. `app/gkit/tools/` — динамические описания из resolver (опционально, можно отдельным этапом)
13. Компиляция и проверка

---

## Edge cases

**Дубли имён пользователей** — ошибка с подсказкой: "Найдено 2 пользователя 'Иван Петров', уточните"

**Имя не найдено** — ошибка: "Воронка 'Продаж' не найдена. Доступные: Продажи, VIP, Партнёры"

**Обратный резолвинг: ID не найден** — возвращаем `"[unknown:7628434]"` (не ломаем ответ)

**Устаревание кеша** — MVP без автообновления. Рестарт = перезагрузка.
