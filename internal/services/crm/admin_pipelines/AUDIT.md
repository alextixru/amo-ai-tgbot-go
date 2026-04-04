# Аудит: admin_pipelines

## Вход (LLM → tool → сервис)

### Поля, которые сейчас числовые ID, но должны быть именами

**`pipeline_id int`**
LLM не знает числовые ID воронок без предварительного `list`. Для каждой операции с конкретной воронкой (get, update, delete, list_statuses и т.д.) нужен лишний round-trip.
Нужно: `pipeline_name string` как альтернатива, резолвинг в handler.

**`status_id int`**
Для `get_status`, `update_status`, `delete_status` LLM должна знать числовой ID статуса.
Нужно: `status_name string` как альтернатива.

**`data.type int` для статуса (0/1/2)**
Тип статуса передаётся числом, семантика нигде не задокументирована в схеме. Допустимые значения: `0` (regular), `1` (won), `2` (lost) — LLM обязана угадывать.

### Поля, которые отсутствуют, но нужны LLM

- **`with_statuses bool`** — флаг для запроса воронки вместе со статусами. SDK поддерживает `?with=statuses`, но `ListPipelines` и `GetPipeline` всегда передают `url.Values{}`. LLM вынуждена делать два вызова.
- **`filter` для `list`** — нет пагинации и фильтрации при `list`. SDK принимает `url.Values` (page, limit, with).
- **Документация `data.color`** — 21 допустимый hex-цвет (`StatusColors`) нигде не перечислен. LLM угадывает.

### Неудобные структуры/типы

- **`data map[string]any`** — нетипизированный. Genkit генерирует JSON Schema как `object` без свойств. LLM угадывает ключи по тексту в description.
- **Три паттерна в одном поле `data`**: для `create` — `data.pipelines: [...]`, для `create_status` batch — `data.statuses: [...]`, для одиночного — поля прямо в `data`. Непоследовательно.
- **Расхождение action-алиасов**: `tools_schema.md` использует `search`/`get_statuses`, handler принимает `list`/`list_statuses`, description перечисляет оба варианта. Три источника истины.

---

## Выход (сервис → tool → LLM)

### Что сейчас возвращается

**`ListPipelines` / `GetPipeline`** → `[]*models.Pipeline` / `*models.Pipeline`:
- `ID`, `Name`, `Sort`, `IsMain`, `IsUnsortedOn`, `IsArchive`
- `AccountID int` — бесполезен для LLM
- `Embedded.Statuses` — **никогда не заполнен** (`with=statuses` не передаётся)
- `Links *Links` — HATEOAS, бесполезны для LLM

**`ListStatuses` / `GetStatus`** → `[]*models.Status` / `*models.Status`:
- `ID`, `Name`, `Sort`, `Color`, `IsEditable`
- `AccountID int` — бесполезен
- `PipelineID int` — число без имени воронки
- `Type int` — `0`/`1`/`2` без семантики
- `Embedded.Descriptions` — **никогда не запрашиваются**

### Какие SDK With-параметры не используются

- **`with=statuses` для воронок** — критичный пропуск. `ListPipelines` и `GetPipeline` всегда `url.Values{}`. За один вызов получить воронку со статусами невозможно.
- **`GetOne` для воронки** — SDK принимает `WithRelations("statuses")`, сервис вызывает `GetOne(ctx, id)` без опций.
- **`with=descriptions` для статуса** — `GetStatus` вызывает `GetOne(ctx, statusID, url.Values{})`, описания никогда не возвращаются.

### Числовые ID в ответе, которые LLM не может интерпретировать

| Поле | Модель | Проблема |
|------|--------|----------|
| `account_id` | Pipeline, Status | Всегда одинаковый, бесполезен |
| `pipeline_id` | Status | Число без имени воронки |
| `type` | Status | `0`, `1`, `2` без семантики |
| `created_by`, `updated_by` | StatusDescription | ID без имён |
| `created_at`, `updated_at` | StatusDescription | Unix timestamp (int) |

SDK предоставляет методы `IsWon()`, `IsLost()`, `IsClosed()` на модели `Status`, но они не используются при формировании ответа.

### Что теряется по сравнению с тем, что SDK может вернуть

1. **`Pipeline.Embedded.Statuses`** — список статусов воронки в одном запросе. Никогда не возвращается.
2. **`Status.Embedded.Descriptions`** — подсказки для менеджеров по уровням. Никогда не запрашиваются.
3. **Семантические метки типа статуса** — `IsWon()`, `IsLost()`, `IsClosed()` вычисляются на модели, в ответе только сырое `"type": 0`.
4. **Пагинация** — `*PageMeta` отбрасывается, LLM не знает есть ли следующая страница.
5. **`Pipeline.IsArchive`** — поле есть в модели, но в `ToAPI()` не сериализуется.

---

## Итого

**Приоритет рефакторинга: высокий**

Функционально рабочий, но без адаптации для LLM: вход требует числовых ID, выход содержит нечитаемые числовые поля, `with=statuses` не используется — главная SDK-возможность пропущена.

### Список конкретных изменений

**Вход:**
1. Добавить `pipeline_name string` в input — lookup по имени через `ListPipelines` если `pipeline_id == 0`
2. Добавить `status_name string` — аналогичный lookup для операций со статусами
3. Заменить `data map[string]any` на типизированные `PipelineData` и `StatusData` с документированными полями включая цвета и семантику `type`
4. Добавить `with_statuses bool` для запроса воронок со статусами в одном вызове
5. Унифицировать батч-режим: всегда `data.items: [...]`
6. Синхронизировать action-алиасы между description, `tools_schema.md` и handler

**Выход:**
7. В `ListPipelines` добавить `?with=statuses` по умолчанию
8. В `GetPipeline` передавать `WithRelations("statuses")`
9. В `GetStatus` запрашивать `with=descriptions`
10. Добавить в ответ по статусу: `is_won bool`, `is_lost bool`, `is_closed bool`, `type_label string` ("regular"/"won"/"lost")
11. Пробрасывать `*PageMeta` в ответ
12. Убрать `account_id` и `_links` из ответа LLM
