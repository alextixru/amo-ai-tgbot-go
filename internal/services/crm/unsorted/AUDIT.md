# Аудит: unsorted

## Вход (LLM → tool → сервис)

### Поля, которые сейчас числовые ID, но должны быть именами

**`filter.pipeline_id []int`**
LLM знает воронки по именам (они вшиты в контекст). Сейчас модель вынуждена держать в голове или угадывать числовой ID. Должно быть:
```json
{ "filter": { "pipeline_name": "Входящие заявки" } }
```

**`accept_params.user_id int`**
LLM знает пользователей по именам. Передавать числовой ID — антипаттерн. Должно быть:
```json
{ "accept_params": { "user_name": "Иван Петров", "status_name": "..." } }
```

**`accept_params.status_id int`**
Статус — часть воронки, LLM знает его по имени ("Первичный контакт", "Переговоры"). Должно быть:
```json
{ "accept_params": { "pipeline_name": "Основная", "status_name": "Первичный контакт" } }
```

**`link_data.lead_id int`**
LLM не знает числовые ID сделок наизусть. Минимум — явная документация что требуется числовой ID сделки.

**`create_data.items[].pipeline_id int`**
Аналогично filter: LLM должна передавать имя воронки, сервис резолвит.

### Поля, которые отсутствуют, но нужны LLM

- **Фильтр по дате для `list`** — `UnsortedFilter` не содержит `created_at_from / created_at_to`, хотя `UnsortedSummaryFilter` их поддерживает. LLM не может фильтровать по времени.
- **Сортировка для `list`** — `BaseFilter.Order` есть в SDK, но не пробрасывается. LLM не может попросить "последние заявки".
- **`decline` использует `AcceptParams`** — нет отдельной структуры `UnsortedDeclineParams`, неочевидно для LLM.
- **Отсутствует типизированный `metadata` в `UnsortedCreateItem`** — `Data map[string]any` принимает это неявно, JSON Schema не описывает структуру.

### Неудобные структуры/типы

- **`UnsortedFilter.PipelineID []int`** — SDK поддерживает только один `pipeline_id`, срез вводит в заблуждение (LLM может передать несколько ожидая OR-фильтрацию).
- **`create_data.items[].created_at int`** — Unix timestamp неудобен для LLM. Лучше принимать RFC3339 строку.
- **`action: "search"` и `action: "list"`** — оба обрабатываются одинаково, дублирование без чёткого правила путает LLM.

---

## Выход (сервис → tool → LLM)

### Что сейчас возвращается

| Action | Тип возврата | Структура |
|--------|-------------|-----------|
| list/search | `[]*models.Unsorted` | uid, category, pipeline_id, created_at, source_name, metadata, _embedded |
| get | `*models.Unsorted` | один объект той же структуры |
| accept | `*models.UnsortedAcceptResult` | `{uid, account_id, result bool}` |
| decline | `*models.UnsortedDeclineResult` | `{uid, account_id, result bool}` |
| link | `*models.UnsortedLinkResult` | `{uid, account_id, result bool}` |
| summary | `*models.UnsortedSummary` | `{total, accepted, declined, average_sort_time, categories map[string]int}` |
| create | `[]*models.Unsorted` | массив созданных объектов |

### Какие SDK With-параметры не используются

`UnsortedFilter` наследует `BaseFilter` с `SetWith()` — механизм есть, но в `ListUnsorted` **никогда не вызывается**. Доступные `with`:
- `leads` — вложенная сделка с полными данными
- `contacts` — связанные контакты
- `companies` — связанные компании

Сервис игнорирует `with` полностью.

### Числовые ID в ответе, которые LLM не может интерпретировать

- **`pipeline_id int`** — LLM видит число, не знает имя воронки
- **`account_id int`** в `Accept/Decline/LinkResult` — технический ID, бесполезен для LLM
- **`Metadata.CalledAt int64`, `Metadata.ReceivedAt int64`** — Unix timestamps
- **`UnsortedEmbedded.Leads[].ID`** — если бы `with=leads` запрашивался, вернулись бы сделки с нечитаемыми полями

### Что теряется по сравнению с тем, что SDK может вернуть

1. **Вложенные контакты и компании** (`with=contacts,companies`) — не запрашиваются
2. **Связанная сделка с деталями** (`with=leads`) — не запрашивается
3. **Пагинация** — `ListUnsorted` отбрасывает `*PageMeta`
4. **Сортировка** — `BaseFilter.Order` не пробрасывается
5. **Типизированные метаданные** — SDK имеет `SipMetadata`, `MailMetadata`, `FormsMetadata`, `ChatsMetadata`, но в ответе используется единая плоская `UnsortedMetadata`
6. **Фильтр по дате для `summary`** — `UnsortedSummaryFilter` поддерживает `created_at[from/to]`, но не пробрасывается из входной модели

---

## Итого

**Приоритет рефакторинга: высокий**

Сервис — тонкая обёртка без адаптации: числовые ID насквозь в обе стороны, `with`-параметры не используются вообще, пагинация теряется, метаданные не типизированы.

### Список конкретных изменений

**Вход:**
1. Заменить `filter.pipeline_id []int` → `filter.pipeline_name string`; резолвинг в сервисе
2. Заменить `accept_params.user_id int` → `accept_params.user_name string`; резолвинг через users
3. Заменить `accept_params.status_id int` → `accept_params.pipeline_name + status_name`
4. Добавить `filter.created_at_from` и `filter.created_at_to` (RFC3339) в `UnsortedFilter`
5. Добавить `filter.order` (`created_at asc/desc`) в `UnsortedFilter`
6. Создать отдельную структуру `UnsortedDeclineParams` с `user_name string`
7. Раскрыть `metadata` в `UnsortedCreateItem` как типизированные объекты по категории
8. Убрать дублирование `action: "search"` / `"list"` — один алиас
9. Исправить `PipelineID []int` → `PipelineID *string` (после замены на имя)

**Выход:**
10. Запрашивать `with=leads,contacts,companies` в `ListUnsorted` и `GetUnsorted`
11. Добавить `pipeline_name string` в ответ (резолвить из reference пайплайнов)
12. Конвертировать Unix timestamps в RFC3339 (`created_at`, `metadata.called_at`, `metadata.received_at`)
13. Пробросить `PageMeta` из `ListUnsorted` — вернуть LLM информацию о пагинации
14. Убрать `account_id` из результатов `accept/decline/link` — маппировать в `{ "uid": "...", "success": true }`
15. Добавить `created_at_from/to` в `SummaryUnsorted` через входной фильтр
