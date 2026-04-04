# Рефакторинг crm/ — полноценный адаптерный слой

## Проблема

CRM-сервисы сейчас — тонкие прокладки с маппингом 1:1:

- На входе принимают числовые ID от tools, пробрасывают в SDK
- На выходе возвращают сырые SDK-модели (`*models.Lead` с `pipeline_id: 7628434`)
- LLM не знает что значат числа, вынуждена сама вызывать справочники

## Цель

Каждый CRM-сервис — самодостаточный переводчик между языком LLM и языком SDK:

- **При инициализации** — сервис сам загружает из SDK те справочники, которые ему нужны
- **На входе** — принимает имена от tools, сам конвертит в ID для SDK
- **На выходе** — получает SDK-модели с числами, отдаёт в tools читаемые структуры с именами
- **Для описаний tools** — отдаёт списки доступных значений ("Воронки: Продажи, VIP")

---

## Поток данных после рефакторинга

```
                    ИНИЦИАЛИЗАЦИЯ (один раз при старте)

agent.go → New(ctx, sdk) → сервис сам загружает нужные справочники из SDK
                           (воронки, статусы, пользователи, кастомные поля)
                           и хранит в своих внутренних мапах

                    ВХОД (каждый вызов tool)

Tool Input                        CRM-сервис                    SDK
{pipeline_name: "Продажи"}   →   ищет в своей мапе →7628434   →  sdk.Leads().Create()
{status_name: "Новая заявка"} →  ищет в своей мапе →status_id  →
{responsible_user: "Иван"}    →  ищет в своей мапе →user_id    →

                    ВЫХОД (каждый ответ)

SDK                              CRM-сервис                     Tool Output → LLM
Lead{PipelineID: 7628434}   →   ищет в своей мапе →"Продажи"  →  {pipeline: "Продажи"}
Lead{StatusID: 48953210}    →   ищет в своей мапе              →  {status: "Новая заявка"}
Lead{ResponsibleUserID: 9}  →   ищет в своей мапе              →  {responsible: "Иван Петров"}

                    ОПИСАНИЯ TOOLS (при регистрации)

CRM-сервис отдаёт метаданные из своих мап:
  PipelineNames() → ["Продажи", "VIP", "Партнёры"]
  UserNames()     → ["Иван Петров", "Мария Сидорова"]

Tool description формируется динамически:
  "Название воронки. Доступные: Продажи, VIP, Партнёры"
```

---

## Все сервисы

### Требуют рефакторинга (числовые ID → имена)

| Сервис | Что загружает при init | Вход: имя→ID | Выход: ID→имя |
|--------|----------------------|--------------|---------------|
| `entities/` | воронки, статусы, пользователи, кастомные поля | PipelineName, StatusName, ResponsibleUserName, CustomFields по имени | Lead/Contact/Company с именами вместо ID |
| `complex_create/` | воронки, статусы, пользователи | то же для LeadData, ContactData, CompanyData | Lead с именами |
| `activities/` | пользователи | ResponsibleUserName в TaskData | Task с именем ответственного |
| `unsorted/` | воронки, статусы, пользователи | PipelineName, StatusName, UserName в Accept/Decline/Create | Unsorted с именами воронки, статуса, пользователя |
| `customers/` | пользователи, статусы покупателей, сегменты | ResponsibleUserName, StatusName, SegmentName | Customer с именами статуса, сегмента, ответственного |
| `catalogs/` | каталоги (id→name) | CatalogName, ElementName, EntityName в Link/Unlink | Каталоги и элементы с именами |
| `products/` | каталог продуктов | ProductName, EntityName в Link/Unlink | Продукты с именами |

### Не требуют рефакторинга (уже работают с именами)

| Сервис | Почему |
|--------|--------|
| `admin_pipelines/` | возвращает уже с именами (Pipeline.Name, Status.Name) |
| `admin_users/` | User.Name уже строка |
| `admin_schema/` | CustomField.Name, Source.Name уже строки |
| `admin_integrations/` | работает с кодами и строками |
| `files/` | UUID-based, нет числовых ID справочников |

---

## Входные модели — замена ID на Name

### models/tools/entities.go — EntityData

| Сейчас                                                | После                           |
| ----------------------------------------------------------- | ------------------------------------ |
| `PipelineID int`                                          | `PipelineName string`              |
| `StatusID int`                                            | `StatusName string`                |
| `ResponsibleUserID int`                                   | `ResponsibleUserName string`       |
| `CustomFieldsValues map[string]any` (ключ = field_id) | ключ = название поля |

### models/tools/entities.go — EntitiesFilter

| Сейчас                         | После                        |
| ------------------------------------ | --------------------------------- |
| `PipelineID []int`                 | `PipelineNames []string`        |
| `StatusID []int`                   | `StatusNames []string`          |
| `ResponsibleUserID []int`          | `ResponsibleUserNames []string` |
| `CreatedBy []int`                  | `CreatedByNames []string`       |
| `UpdatedBy []int`                  | `UpdatedByNames []string`       |
| `CustomFieldsValues[].FieldID int` | `FieldName string`              |

### models/tools/activities.go

| Сейчас                            | После                        |
| --------------------------------------- | --------------------------------- |
| `TaskData.ResponsibleUserID int`      | `ResponsibleUserName string`    |
| `TasksFilter.ResponsibleUserID []int` | `ResponsibleUserNames []string` |
| `ActivitiesInput.UserIDs []int`       | `UserNames []string`            |
| `ActivitiesInput.UserID int`          | `UserName string`               |

### models/tools/complex_create.go

| Сейчас                          | После                     |
| ------------------------------------- | ------------------------------ |
| `LeadData.PipelineID int`           | `PipelineName string`        |
| `LeadData.StatusID int`             | `StatusName string`          |
| `LeadData.ResponsibleUserID int`    | `ResponsibleUserName string` |
| `ContactData.ResponsibleUserID int` | `ResponsibleUserName string` |
| `CompanyData.ResponsibleUserID int` | `ResponsibleUserName string` |

### models/tools/unsorted.go

| Сейчас | После |
|--------|-------|
| `AcceptParams.UserID int` | `UserName string` |
| `AcceptParams.StatusID int` | `StatusName string` |
| `DeclineParams.UserID int` | `UserName string` |
| `UnsortedFilter.PipelineID int` | `PipelineName string` |
| `CreateData.Items[].PipelineID int` | `PipelineName string` |

### models/tools/customers.go

| Сейчас | После |
|--------|-------|
| `CustomerData.ResponsibleUserID int` | `ResponsibleUserName string` |
| `CustomerData.StatusID int` | `StatusName string` |
| `CustomersFilter.ResponsibleUserIDs []int` | `ResponsibleUserNames []string` |
| `CustomersFilter.StatusIDs []int` | `StatusNames []string` |

### models/tools/catalogs.go

| Сейчас | После |
|--------|-------|
| `CatalogID int` | `CatalogName string` |
| `ElementID int` | `ElementName string` |
| `LinkData.EntityID int` | `EntityName string` |

### models/tools/products.go

| Сейчас | После |
|--------|-------|
| `ProductID int` | `ProductName string` |
| `LinkData.EntityID int` | `EntityName string` |
| `PriceID int` | `PriceName string` |

---

## Edge cases

**Имя не найдено** — сервис возвращает ошибку с подсказкой: "Воронка 'Продаж' не найдена. Доступные: Продажи, VIP, Партнёры"

**Дубли имён** — ошибка: "Найдено 2 пользователя 'Иван Петров', уточните"

**Обратный резолвинг: ID не найден** — вернуть `"[unknown:7628434]"`, не ломать ответ

**Устаревание** — MVP без автообновления. Рестарт = перезагрузка.

---

## Источники для рефакторинга

| Что                               | Путь                                 | Зачем                               |
| ------------------------------------ | ---------------------------------------- | ---------------------------------------- |
| Этот план                    | `internal/services/crm/REFACTORING.md` | Общий контекст              |
| Сервис (текущий код) | `internal/services/crm/{service}/`     | Что рефакторим              |
| Входная модель          | `internal/models/tools/{file}.go`      | Input от LLM                           |
| Tool (тонкий слой)         | `app/gkit/tools/{file}.go`             | Как вызывается сервис |
| SDK-модели                     | `amocrm-sdk-go/core/models/`           | Что возвращает SDK          |
| SDK-сервисы                   | `amocrm-sdk-go/core/services/`         | Какие методы доступны |
| SDK-фильтры                   | `amocrm-sdk-go/core/filters/`          | Параметры запросов      |
