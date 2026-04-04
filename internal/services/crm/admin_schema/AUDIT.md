# Аудит: admin_schema

## Вход (LLM → tool → сервис)

### Поля, которые сейчас числовые ID, но должны быть именами

**`RequiredStatus.PipelineID int` / `StatusID int`** в `CustomField`
Чтобы сделать поле обязательным в статусе "Переговоры", LLM должна сначала запросить ID через `admin_pipelines`.
Сейчас: `"required_statuses": [{"pipeline_id": 1234567, "status_id": 9876543}]`
Нужно: `"required_statuses": [{"pipeline_name": "Основная", "status_name": "Переговоры"}]`

**`Source.PipelineID int`**
Источник привязывается к воронке по числовому ID. LLM обязана знать ID заранее.

**`CustomField.GroupID string`**
Строковый, но непрозрачный (значения типа `"statistic"`, `"general"` — нестандартные ID из CRM). Без предварительного `search` по `field_groups` LLM не знает допустимых значений.

### Поля, которые отсутствуют, но нужны LLM

- **Фильтр по имени** — нет поиска кастомного поля по `name`. Вынуждена получать полный список и искать вручную. (API не поддерживает — нужна client-side фильтрация.)
- **Фильтр по `code`** — символьный код (`"PHONE"`, `"SOURCE_ID"`) более стабилен чем имя, но фильтровать по нему нельзя.
- **`Order` в `SchemaFilter`** — `BaseFilter` SDK поддерживает `Order map[string]string`, не пробрасывается.
- **`With` параметры** — не экспонированы в `SchemaFilter`.
- **Поиск `loss_reasons` по имени** — `ListLossReasons` транслирует только `limit`/`page`.

### Неудобные структуры/типы

- **`Data map[string]any`** — полностью неструктурированный. Genkit не генерирует JSON Schema для вложенных объектов. LLM обязана знать что в `data.fields`, `data.groups` и т.д. — это задокументировано только в тексте description.
- **`Layer + Action` без валидации** — `layer=loss_reasons, action=update` допустимо по JSON Schema, упадёт в runtime. Ограничение только в текстовом описании.
- **`GetCustomField`/`GetFieldGroup`** — передают `url.Values{}` хардкодом, `With` недоступен даже теоретически.
- **`search` и `list` как алиасы** — дублирование без ценности.

---

## Выход (сервис → tool → LLM)

### Что сейчас возвращается

| Метод | Тип возврата |
|-------|-------------|
| `ListCustomFields` / `GetCustomField` | `[]*models.CustomField` / `*models.CustomField` |
| `CreateCustomFields` / `UpdateCustomFields` | `[]*models.CustomField` |
| `ListFieldGroups` / `GetFieldGroup` | `[]models.CustomFieldGroup` / `*models.CustomFieldGroup` |
| `CreateFieldGroups` / `UpdateFieldGroups` | `[]models.CustomFieldGroup` |
| `ListLossReasons` / `GetLossReason` | `[]*models.LossReason` / `*models.LossReason` |
| `CreateLossReasons` | `[]*models.LossReason` |
| `ListSources` / `GetSource` | `[]*models.Source` / `*models.Source` |
| `CreateSources` / `UpdateSources` | `[]*models.Source` |
| `delete` (все) | `nil` — нет подтверждения! |

`PageMeta` (второй параметр SDK-методов) **игнорируется везде**.

### Какие SDK With-параметры не используются

**Критически: `with=fields` для `ListFieldGroups`**
`CustomFieldGroup.Fields []CustomField` будет пустым массивом при list-запросе без `with=fields`. LLM видит только имена групп без содержимого — запрос бесполезен.

Для `custom_fields` — `With` технически пробрасывается через `buildCustomFieldsFilter`, но в `GetCustomField` передаётся `url.Values{}`.

### Числовые ID в ответе, которые LLM не может интерпретировать

- `CustomField.AccountID int` — всегда одинаковый, бесполезный шум
- `RequiredStatus.PipelineID int` / `StatusID int` — без имён воронки/статуса
- `Source.PipelineID int` — без имени воронки
- `CustomField.GroupID string` — непрозрачная строка без имени группы
- `LossReason.CreatedAt int` / `UpdatedAt int` — Unix timestamps
- `CustomFieldEnum.ID int` — при обновлении enum LLM должна оперировать этими числами

### Что теряется по сравнению с тем, что SDK может вернуть

1. **`PageMeta`** — полностью отбрасывается. При >50 полях LLM получит неполный список без предупреждения
2. **`with=fields` в `ListFieldGroups`** — `Fields` всегда пуст при list
3. **`CustomField.IsDeletable bool`, `IsPredefined bool`** — полезны LLM (можно ли удалить), но не описаны в input
4. **Подтверждение delete** — возвращается `nil`, необратимые операции без подтверждения
5. **`CustomField.Links *Links`** — HATEOAS-мусор в ответе

---

## Итого

**Приоритет рефакторинга: средний**

Функционально корректен, критических ошибок нет. Рефакторинг нужен для повышения качества работы LLM, а не исправления поломок. Один критический пропуск SDK: `with=fields` для групп полей.

### Список конкретных изменений

1. **[ВЫСОКИЙ] Добавить `with=fields` в `ListFieldGroups`** — передавать `v.Set("with", "fields")` в фильтре чтобы `Fields` приходил заполненным
2. **[ВЫСОКИЙ] Прокинуть `PageMeta` в ответ** — `{items, total_count, has_more}` для всех list-методов
3. **[ВЫСОКИЙ] Типизировать `Data map[string]any`** — заменить на `CustomFieldData`, `FieldGroupData`, `LossReasonData`, `SourceData`
4. **[СРЕДНИЙ] Добавить `name` как фильтр в `SchemaFilter`** — client-side фильтрация по имени
5. **[СРЕДНИЙ] Добавить `pipeline_name` в ответе `Source`** — резолвить рядом с `pipeline_id` из контекста сессии
6. **[СРЕДНИЙ] Добавить `order` в `SchemaFilter`** — пробросить в `buildCustomFieldsFilter`
7. **[СРЕДНИЙ] Возвращать `{success: true, deleted_id: N}` из delete-операций**
8. **[НИЗКИЙ] Добавить фильтр по `code` в `SchemaFilter`** — client-side
9. **[НИЗКИЙ] Убрать алиасы `search`/`list`** — оставить один
10. **[НИЗКИЙ] Убрать `AccountID` и `Links` из ответа** — шум для LLM
11. **[НИЗКИЙ] Добавить `required_statuses` с `pipeline_name`/`status_name`** — резолвинг в сервисе
