# Аудит: entities

## Вход (LLM → tool → сервис)

### Поля, которые сейчас числовые ID, но должны быть именами

**`EntitiesFilter.responsible_user_id []int`**
LLM знает пользователей по именам из Reference-контекста, но вынуждена передавать числовые ID.
Сейчас: `"responsible_user_id": [12345]`
Должно быть: `"responsible_user_name": ["Иван Петров"]`

**`EntitiesFilter.created_by []int` / `updated_by []int`**
Та же проблема — имена пользователей из контекста не резолвятся в сервисе.

**`EntitiesFilter.pipeline_id []int` (только leads)**
LLM знает воронку по имени, нет резолвинга имя → ID.
Сейчас: `"pipeline_id": [100500]`
Должно быть: `"pipeline_name": ["Основная воронка"]`

**`EntitiesFilter.status_id []int` (только leads)**
SDK-фильтр по статусам требует пары `{status_id, pipeline_id}`, но в адаптере `SetStatuses` **не вызывается вообще** — поле объявлено в модели, но полностью потеряно при маппинге. Баг.

**`EntityData.status_id int` / `pipeline_id int` (create/update leads)**
При создании/обновлении сделки LLM обязана знать числовые ID. Сервис не резолвит имена.

**`EntityData.responsible_user_id int`**
Имя пользователя есть в контексте, числовой ID — нет. Нет резолвинга.

**`EntitiesFilter.CustomFieldsValues[].field_id int`**
Фильтр по кастомным полям принимает числовой `field_id`. LLM знает поля по кодам/именам из Reference. Нет поддержки `field_code` как альтернативы.

**`EntityData.CustomFieldsValues map[string]any`**
Ключи — `field_code` (строковые, что хорошо), но маппинг через двойной JSON-roundtrip (`mapCustomFieldsValues`) нестабилен: сначала пробует `[]CustomFieldValue`, при неудаче — `map[string]any`. Неочевидный контракт, провоцирует молчаливые ошибки.

### Поля, которые отсутствуют, но нужны LLM

- **Фильтр по статусам через имена**: нет `status_name []string` или `[{"pipeline": "...", "status": "..."}]` пар. SDK ожидает `{StatusID, PipelineID}` — вызов `SetStatuses` отсутствует в адаптере.
- **`EntitiesFilter.closest_task_at`**: SDK поддерживает `SetClosestTaskAt(from, to)`, но поля нет в модели.
- **`EntitiesFilter.created_by` / `updated_by` для leads**: поля есть в модели, `SetCreatedBy`/`SetUpdatedBy` есть в SDK, но вызовы в `leads.go` отсутствуют — баг.
- **`EntityData.loss_reason_id` / `loss_reason_name`**: при закрытии сделки как проигранной нужна причина отказа. SDK `Lead.LossReasonID *int` есть, в `EntityData` — нет.
- **`EntityData.source_id` / `source_name`**: SDK `Lead.SourceID *int` есть, в `EntityData` — нет.
- **`EntityData` для contacts: `first_name` / `last_name`**: SDK `Contact` поддерживает раздельные поля, `EntityData` имеет только `Name`.
- **`mapToContact` EmbeddedCompanies не реализован**: блок `if len(data.EmbeddedCompanies) > 0` содержит только комментарий — привязка компаний к контакту не работает. Баг.

### Неудобные структуры/типы

- **Временные метки как Unix timestamp `int64`**: `created_at_from/to`, `updated_at_from/to`, `closed_at_from/to` — LLM сама конвертирует. Нет поддержки ISO-8601.
- **`filter.status_id []int` + `filter.pipeline_id []int` — разорванная пара**: SDK требует `{status_id, pipeline_id}`, а в модели два независимых массива без гарантии соответствия.
- **`CustomFieldFilter.values []any`**: слишком широкий тип, нет подсказки LLM какие типы допустимы.
- **`input.With` + `filter.With`**: два места для `with`, сливаются через `append` — неочевидно для LLM.

---

## Выход (сервис → tool → LLM)

### Что сейчас возвращается

Сервис возвращает сырые SDK-модели без какого-либо преобразования:

| Метод | Тип возврата |
|-------|-------------|
| `SearchLeads` | `[]*models.Lead` |
| `GetLead` | `*models.Lead` |
| `CreateLead` | `*models.Lead` |
| `CreateLeads` | `[]*models.Lead` |
| `UpdateLead` / `UpdateLeads` | `*models.Lead` / `[]*models.Lead` |
| `SearchContacts` | `[]*models.Contact` |
| `GetContact` | `*models.Contact` |
| `CreateContact` / `UpdateContact` | `[]*models.Contact` (слайс, а не один!) |
| `SearchCompanies` | `[]*models.Company` |
| `GetCompany` | `*models.Company` |
| `CreateCompany` / `UpdateCompany` | `[]*models.Company` |
| `LinkLead` | `[]models.EntityLink` |
| `LinkContact` / `UnlinkContact` / `LinkCompany` / `UnlinkCompany` | `error` (nil → LLM не получает ничего) |

### Какие SDK With-параметры не используются

По умолчанию `with` передаётся только если LLM явно его указала. Нет дефолтных значений.

Для **leads** доступны но не запрашиваются по умолчанию:
- `contacts` — связанные контакты
- `companies` — связанные компании
- `catalog_elements` — товары
- `loss_reason` — причина отказа (содержит имя!)
- `source` — источник (содержит имя!)

Для **contacts**:
- `leads`, `companies`, `catalog_elements`, `social_profiles`

Для **companies**:
- `leads`, `contacts`, `catalog_elements`

SDK также предоставляет `ContactsService.GetLinks` и `CompaniesService.GetLinks` — не представлены в интерфейсе сервиса вообще.

### Числовые ID в ответе, которые LLM не может интерпретировать

- `responsible_user_id` — нет имени пользователя
- `created_by` / `updated_by` — нет имён
- `pipeline_id` — нет названия воронки
- `status_id` — нет названия статуса
- `loss_reason_id` — нет названия причины отказа
- `source_id` — нет названия источника
- `group_id` — нет названия группы
- `account_id` — избыточное поле (всегда одинаково)
- `CustomFieldsValues[].field_id` — числовые ID кастомных полей (хотя `field_code` тоже есть)
- `ClosestTaskAt` — Unix timestamp `*int64`, нечитаемо

### Что теряется по сравнению с тем, что SDK может вернуть

1. **Пагинация**: `sdk.Leads().Get()` возвращает `([]*Lead, *PageMeta, error)` — сервис отбрасывает `*PageMeta`. LLM не знает есть ли следующая страница.
2. **`LossReason` при `with=loss_reason`**: содержит читаемое название причины — не запрашивается.
3. **`Source` при `with=source`**: содержит `Name string` источника — не запрашивается.
4. **`SocialProfile` у контактов при `with=social_profiles`** — не запрашивается.
5. **`GetLinks`** у contacts и companies — не представлен в интерфейсе.
6. **Несогласованность типов**: `CreateLead` → `*Lead` (один), `CreateContact` → `[]*Contact` (слайс). Разные форматы для одинаковых операций.
7. **`link`/`unlink` возвращает `nil`**: LLM не получает подтверждения успеха.

---

## Итого

**Приоритет рефакторинга: высокий**

Сервис — тонкая обёртка без адаптации. Кроме проблем с именами/ID есть реальные баги: `status_id` фильтр не работает, `created_by`/`updated_by` теряются, `EmbeddedCompanies` не реализован, кастомные поля в фильтре не работают.

### Список конкретных изменений

**Вход — баги:**
1. Добавить вызов `f.SetStatuses()` в `leads.go` — сейчас `filter.StatusID` игнорируется полностью
2. Добавить вызовы `f.SetCreatedBy` и `f.SetUpdatedBy` в `leads.go` — поля потеряны
3. Добавить маппинг `filter.CustomFieldsValues` во всех трёх файлах — `SetCustomFieldsValues` не вызывается нигде
4. Исправить `mapToContact`: реализовать блок `EmbeddedCompanies`

**Вход — резолвинг имён:**
5. Заменить `responsible_user_id int/[]int` → `responsible_user_name string/[]string` в `EntityData` и `EntitiesFilter`
6. Заменить `pipeline_id int/[]int` → `pipeline_name string/[]string` в `EntityData` и `EntitiesFilter`
7. Заменить `status_id int/[]int` → структурированные пары `[{"pipeline": "...", "status": "..."}]` в фильтре, `status_name string` в `EntityData`
8. Добавить `field_code string` как альтернативу `field_id` в `CustomFieldFilter`
9. Добавить `closest_task_at_from/to string` (ISO-8601) в `EntitiesFilter`
10. Принимать даты в ISO-8601 вместо Unix timestamp

**Вход — недостающие поля:**
11. Добавить `loss_reason_name string` и `source_name string` в `EntityData`
12. Добавить `first_name string` / `last_name string` в `EntityData` для контактов

**Выход — критические:**
13. Возвращать `PageMeta` из поиска — LLM не знает о неполноте результатов
14. Унифицировать типы возврата: `CreateContact`/`CreateCompany` должны возвращать один объект как `CreateLead`
15. Возвращать читаемый ответ из `link`/`unlink` вместо `nil`

**Выход — обогащение:**
16. Добавить дефолтные `with=["loss_reason", "source"]` для `GetLead`
17. Создать read-model для LLM с человеко-читаемыми полями: `responsible_user_name`, `pipeline_name`, `status_name`, `created_by_name`, `created_at` в ISO-8601
18. Убрать технические поля из ответа: `account_id`, `_links`
19. Добавить в интерфейс `Service` методы `GetLeadLinks` / `GetContactLinks` / `GetCompanyLinks`
