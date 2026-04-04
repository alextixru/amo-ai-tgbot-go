# Аудит: activities

## Вход (LLM → tool → сервис)

### Поля, которые сейчас числовые ID, но должны быть именами

**`TaskData.responsible_user_id int`**
LLM знает пользователей по именам из вшитого контекста, но вынуждена самостоятельно резолвить имя → ID.
Нужно: `responsible_user_name string` + резолвинг в сервисе.

**`TasksFilter.responsible_user_id []int`**
Фильтр по ответственному — числовые ID.
Нужно: `responsible_user_names []string`.

**`TasksFilter.task_type_id int`**
Тип задачи (1=звонок, 2=встреча) передаётся числом.
Нужно: `task_type string` — `"follow_up" | "meeting"`, сервис конвертирует.

**`EventsFilter.created_by []int`**
Фильтр событий по создателям — числовые ID.
Нужно: `created_by_names []string`.

**`UserIDs []int` / `UserID int`** — для subscribe/unsubscribe
LLM знает пользователей по именам.
Нужно: `user_names []string` / `user_name string`.

**`TagID int`** — для tags.delete
LLM не знает ID тегов без предварительного `tags.list`. Нужно удаление по имени либо `tag_name` как альтернатива.

### Поля, которые отсутствуют, но нужны LLM

- **`TasksFilter.created_by`** — SDK `filters.TasksFilter` поддерживает `SetCreatedBy([]int)`, но поле не пробрасывается из входной модели
- **`TasksFilter.lead_statuses`** — SDK поддерживает `SetLeadStatuses()` (фильтр задач по воронке/этапу), полностью отсутствует во входной модели
- **`TasksFilter.updated_at_to`** — SDK принимает диапазон, входная модель предоставляет только `updated_at` (только from)
- **`EventsFilter.with`** — `Event.AvailableWith()` объявляет 7 полезных значений (`contact_name`, `lead_name`, `company_name`, `catalog_element_name`, `customer_name`, `catalog_name`, `note`), но поля `with` в `EventsFilter` нет
- **`NoteData.responsible_user_id`** — `models.Note` наследует `BaseModel.ResponsibleUserID`, установить при создании нельзя
- **`CallData.responsible_user_id`** — `models.Call` имеет `ResponsibleUserID` и `CallResponsible`, но `CallData` их не предоставляет

### Неудобные структуры/типы

- **Дублирование одиночных и батч полей** — `task_data` + `tasks_data`, `note_data` + `notes_data`, `link_to` + `links_to`, `tag_name` + `tag_names` = 8 полей вместо 4. LLM должна выбирать нужный вариант. Лучше: всегда массив, один элемент = одиночная операция.
- **`With []string` — свободная строка без подсказки** — LLM не знает допустимые значения для каждого layer. Нет enum-подсказки.
- **`talks` — нет `action: "get"`** — `GetTalk` реализован в интерфейсе сервиса, но не экспонируется через tool handler.
- **Валидация `parent` — запутана** — вложенные `if` с дублирующимися проверками. При `action != "list"` parent проверяется частично, потом повторно в каждом handler'е. Логика противоречива: `subscriptions.list` требует parent, но верхний guard пропускает `list` без parent.
- **`CallStatus` enum расходится с SDK** — описание в коде (`1=успех...6=неправильный номер`) не совпадает с константами SDK (`models.CallStatusSuccess = 7`).

---

## Выход (сервис → tool → LLM)

### Что сейчас возвращается

| Метод | Тип возврата |
|-------|-------------|
| `ListTasks` / `GetTask` | `[]*models.Task` |
| `CreateTask` / `UpdateTask` / `CompleteTask` | `*models.Task` |
| `ListNotes` / `GetNote` / `CreateNote` / `UpdateNote` | `*models.Note` |
| `CreateCall` | `*models.Call` |
| `ListEvents` / `GetEvent` | `[]*models.Event` |
| `ListFiles` / `LinkFiles` | `[]models.FileLink` (только UUID!) |
| `ListLinks` / `LinkEntities` | `[]*models.EntityLink` |
| `ListTags` / `CreateTag` | `[]*models.Tag` |
| `ListSubscriptions` / `Subscribe` | `[]models.Subscription` |

### Какие SDK With-параметры не используются

**Events** — `with` вообще не поддерживается на уровне входа. `Event.AvailableWith()` объявляет 7 полезных значений: `contact_name`, `lead_name`, `company_name`, `catalog_element_name`, `customer_name`, `catalog_name`, `note` — **все теряются**. Event без `with` — только `type` + `entity_id` (число) + `value_before/after` (interface{}). LLM не может интерпретировать.

**GetNote** — вызывается как `sdk.Notes().GetOne(ctx, entityType, id, nil)` — параметр `params` всегда `nil`.

**Tasks** — `with` пробрасывается, но LLM не знает допустимые значения. Task поддерживает `with=leads`, `with=contacts`.

### Числовые ID в ответе, которые LLM не может интерпретировать

- `Task.ResponsibleUserID int` — нет имени
- `Task.TaskTypeID int` — возвращается `1` или `2`, LLM должна помнить значения
- `Task.CompleteTill *int64` — Unix timestamp
- `Task.CreatedAt / UpdatedAt int64` — Unix timestamps
- `Note.CreatedBy / UpdatedBy int` — числовые ID
- `Event.CreatedBy int` — числовой ID
- `Call.CreatedBy / ResponsibleUserID int` — числовые ID
- `Subscription.SubscriberID int` — LLM не знает кто подписан
- `EntityLink.ToEntityID / EntityID int` — без имён

### Что теряется по сравнению с тем, что SDK может вернуть

1. **Events `with`** — весь контекст событий (имена контактов, сделок, компаний) недоступен
2. **`FileLink`** — только UUID файла, без `name`, `size`, `download_link`, `created_at`
3. **Пагинация** — `PageMeta` отбрасывается везде, LLM не знает есть ли следующая страница
4. **`Note.Text`** — данные приходят через `Params.Text`, а не плоское `text`. Неудобно.
5. **`Task.Duration`** — поле модели, не документируется для LLM

---

## Итого

**Приоритет рефакторинга: высокий**

Числовые ID проходят насквозь в обоих направлениях. SDK-возможности (особенно Events `with`) полностью не используются. LLM получает нечитаемые timestamps и ID вместо имён и дат.

### Список конкретных изменений

**Вход:**
1. Добавить `responsible_user_name string` в `TaskData`, `responsible_user_names []string` в `TasksFilter`; резолвинг через Users reference
2. Добавить `user_names []string` / `user_name string` для subscriptions вместо числовых ID
3. Заменить `TaskTypeID int` → `task_type string` (`"follow_up" | "meeting"`)
4. Добавить `with []string` в `EventsFilter`; по умолчанию запрашивать `["contact_name", "lead_name", "company_name"]`
5. Добавить `created_by_names []string` в `EventsFilter`
6. Добавить `created_by []string` в `TasksFilter` (пробросить `SetCreatedBy`)
7. Добавить `lead_statuses` в `TasksFilter` (пробросить `SetLeadStatuses`)
8. Добавить `updated_at_to` в `TasksFilter`
9. Объединить одиночные/батч поля: всегда массив (`tasks_data []TaskData`, `notes_data []NoteData` и т.д.)
10. Добавить `tag_name string` как альтернативу `tag_id` для delete
11. Добавить `action: "get"` в `handleTalks`
12. Исправить валидацию `parent` — централизовать в начале `handleActivities`, убрать дублирование
13. Исправить `CallStatus` enum — привести к актуальным константам SDK

**Выход:**
14. Конвертировать все timestamps в ISO 8601 на выходе (`CompleteTill`, `CreatedAt`, `UpdatedAt`)
15. Резолвить user ID в имена (`responsible_user_name`, `created_by_name`) в ответе
16. Возвращать `TaskType string` вместо `task_type_id int`
17. При `ListFiles` возвращать полные метаданные файла (через `FilesService`), не только UUID
18. Пробросить `PageMeta` из всех list-методов — `HasMore bool`, `Total int`
