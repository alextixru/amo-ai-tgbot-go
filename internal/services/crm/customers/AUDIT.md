# Аудит: customers

## Вход (LLM → tool → сервис)

### Поля, которые сейчас числовые ID, но должны быть именами

**`CustomerData.responsible_user_id int`**
LLM знает пользователей по именам (UsersService в Reference Context), но сервис не делает резолвинг.
Нужно: `responsible_user_name string` + резолвинг в сервисе.

**`CustomerData.status_id int`**
Статусы покупателей отсутствуют в Reference Context — LLM не знает ни имён, ни ID. Для смены статуса нужен двойной round-trip: сначала `customers/statuses/list`, потом update.
Нужно: `status_name string` + резолвинг, либо добавить CustomerStatuses в Reference Context.

**`CustomerFilter.responsible_user_ids []int`**
Тот же случай — фильтр по ответственным принимает числовые ID.
Нужно: `responsible_user_names []string`.

**`CustomerFilter.status_ids []int`**
Фильтр по статусам — числовые ID которые LLM не знает без предварительного list.
Нужно: `status_names []string`.

**`CustomerLinkData.entity_id int`**
При привязке покупателя к контакту/компании нужен числовой entity_id. Нет возможности указать контакт по имени.

### Поля, которые отсутствуют, но нужны LLM

- **`CustomerData.periodicity int`** — SDK `Customer.Periodicity` есть, в `CustomerData` и `mapCustomerData` полностью отсутствует. LLM не может задать периодичность покупок.
- **`CustomerFilter.next_date_from / next_date_to`** — поля декларируются в схеме, но handler содержит комментарий `// NextDate filter is NOT supported`. Молчаливое игнорирование вводит LLM в заблуждение — нет ошибки, но фильтр не работает.
- **`with []string` на верхнем уровне input** — сейчас `with` вложен в `Filter`. Для `get`/`create` операций LLM вынуждена передавать `filter: { "with": [...] }` — неестественная структура.
- **Batch-операции для `statuses`** — SDK поддерживает батч, tool handler принимает только одиночный объект.

### Неудобные структуры/типы

- **`CustomFieldsValues map[string]any`** — двойной JSON marshal/unmarshal в `mapCustomerCustomFieldsValues` ненадёжен: первая попытка десериализовать как `[]CustomFieldValue` почти всегда падает в fallback. В JSON Schema генерируется как `object` без структуры.
- **`NextDate int64` (Unix timestamp)** — LLM не работает с unix timestamp нативно. Нет конвертации из ISO 8601.
- **`CustomerData` переиспользуется для разных слоёв** — при create `statuses`/`segments` handler берёт только `input.Data.Name`. Остальные поля (`responsible_user_id`, `next_date`, `status_id`) бессмысленны для этих слоёв, но присутствуют в схеме — LLM вводится в заблуждение.

---

## Выход (сервис → tool → LLM)

### Что сейчас возвращается

| Action | Тип возврата |
|--------|-------------|
| `list/get customers` | `[]models.Customer` / `models.Customer` |
| `get bonus_points` | `*models.BonusPoints` |
| `earn_points / redeem_points` | `int` (баланс без контекста) |
| `list/get statuses` | `[]models.Status` / `models.Status` |
| `list transactions` | `[]models.Transaction` |
| `list/get segments` | `[]*models.CustomerSegment` / `*models.CustomerSegment` |

### Какие SDK With-параметры не используются

`Customer.AvailableWith()` возвращает: `catalog_elements`, `contacts`, `companies`, `segments`, `group_id`.

`with` передаётся из `input.Filter.With` — работает для `list`/`get`. Но:
- `group_id` не упомянут в схеме для LLM
- При `CreateCustomers`/`UpdateCustomers` `with` не применяется — SDK возвращает объект без embedded данных, повторный get с `with` не делается. LLM не получает теги, сегменты, контакты после create.
- По умолчанию (без явного `with`) ответ — голый объект без связей.

### Числовые ID в ответе, которые LLM не может интерпретировать

- `responsible_user_id int` — нет имени пользователя
- `status_id int` — нет названия статуса (статусы не в Reference Context!)
- `created_by int` / `updated_by int` — нет имён
- `group_id int` — нет названия группы
- В `CustomerEmbedded.Contacts`/`Companies` — объекты только при `with=contacts,companies`, иначе пусто

### Что теряется по сравнению с тем, что SDK может вернуть

1. **`Customer.Ltv`** (Lifetime Value) — полезная метрика, в ответе есть, но не документирована в схеме
2. **`Customer.PurchasesCount`** и **`Customer.AverageCheck`** — SDK возвращает, LLM не знает что они будут в ответе
3. **`Customer.Embedded.Tags`** — поведение неочевидно: при `list` могут не возвращаться
4. **`ListSegments`** — передаёт `nil` в SDK вместо фильтра. Нет пагинации, нет фильтрации.
5. **`ListTransactions`** — передаёт `nil` в SDK вместо фильтра. Нет пагинации.
6. **`earn`/`redeem`** — возвращают только `int` (баланс?), без подтверждения операции и timestamp.

---

## Итого

**Приоритет рефакторинга: высокий**

Тонкая обёртка без адаптации. LLM работает с числовыми ID, не получает читаемых имён, не может корректно использовать ряд операций без дополнительных round-trip запросов. Есть молчаливые баги (next_date фильтр).

### Список конкретных изменений

1. Добавить `responsible_user_name string` в `CustomerData` и `responsible_user_names []string` в `CustomerFilter`; резолвинг через Users reference
2. Добавить `status_name string` в `CustomerData` и `status_names []string` в `CustomerFilter`; резолвинг через CustomerStatuses
3. Добавить `CustomerStatuses` в Reference Context (`tools_schema.md`) — по аналогии с PipelinesService
4. Убрать или явно возвращать ошибку для `next_date_from/to` — текущее молчаливое игнорирование недопустимо
5. Добавить `periodicity int` в `CustomerData` и `mapCustomerData`
6. Перенести `with []string` на верхний уровень `CustomersInput`
7. Обогатить ответ именами: `responsible_user_name`, `status_name` вместо сырых ID
8. Запрашивать `with=contacts,companies,segments` по умолчанию при `get`
9. Исправить `mapCustomerCustomFieldsValues` — заменить двойной JSON roundtrip на типизированную схему
10. Добавить пагинацию и фильтрацию для `ListSegments` и `ListTransactions` (сейчас `nil`)
11. Добавить поддержку batch для `statuses` в tool handler
12. Принимать даты в ISO 8601 (`next_date`, `next_date_from/to`), конвертировать в Unix timestamp в сервисе
