# Аудит: files

## Вход (LLM → tool → сервис)

### Поля, которые отсутствуют, но нужны LLM

`FileFilter` содержит только `page`, `limit`, `uuids`, `with`. SDK `filters.FilesFilter` поддерживает полноценный набор фильтров, ни один из которых не прокинут:

| Поле SDK | Метод | Значение для LLM |
|----------|-------|-----------------|
| `Name` | `SetName` | поиск по имени файла |
| `Term` | `SetTerm` | полнотекстовый поиск |
| `Extensions` | `SetExtensions` | фильтр по типу (pdf, xlsx, jpg) |
| `Deleted` | `SetDeleted` | показать удалённые |
| `DateFrom/DateTo/DateType` | `SetDate` | поиск за период |
| `DatePreset` | `SetDatePreset` | пресеты (today, week, month) |
| `SizeFrom/SizeTo` | `SetSize` | фильтр по размеру |
| `CreatedBy` | `SetCreatedBy` | по автору (числовые ID — нужен резолвинг) |
| `UpdatedBy` | `SetUpdatedBy` | по обновившему (числовые ID — нужен резолвинг) |

Из 9 фильтров SDK прокинуто только 2 (`page`, `limit`) + `uuids`.

### Расхождение документации и реализации (критично)

`tools_schema.md` декларирует:
```json
"upload_params": { "file_name": "...", "file_content": "base64...", "file_url": "https://..." }
```
Реальная `FileUploadParams` содержит:
```go
LocalPath   string  // путь на диске сервера
FileName    string
WithPreview bool
FileUUID    string
```
Поля `file_content` (base64) и `file_url` задокументированы, но не реализованы. LLM может попытаться передать URL или base64 — сервис их молча проигнорирует.

Также `tools_schema.md` декларирует `action: search`, реализация обрабатывает `action: list`. Расхождение.

### Поля с числовыми ID

- `FileFilter.CreatedBy []int` / `UpdatedBy []int` — если добавить в tool-модель, LLM не сможет заполнить без справочника. Нужен `created_by_name string` с резолвингом в сервисе.

### Неудобные структуры/типы

- `With []string` — семантика `deleted`/`unbilled` как `with` неочевидна (это флаги фильтрации, а не embedded данные)
- `UUIDs` (batch) и `UUID` (single) — два разных поля, handler проверяет их по очереди. LLM может запутаться.

---

## Выход (сервис → tool → LLM)

### Что сейчас возвращается

| Action | Тип возврата |
|--------|-------------|
| `list` | `[]*models.File` |
| `get` | `*models.File` |
| `upload` | `*models.File` (частично: только UUID, Name, Extension, MimeType, ID, Size) |
| `update` | `*models.File` |
| `delete` / batch | `nil` (только ошибка) |

### Какие SDK With-параметры не используются

`FileFilter.With` прокидывается через `f.SetWith(...)` — работает для `list`. Но `GetFile` (`GetOneByUUID`) вызывается без `with` — нельзя получить удалённый файл по UUID даже передав `with=deleted`.

### Числовые ID в ответе, которые LLM не может интерпретировать

- `CreatedBy *UserLink` — `{id: int, type: string}`, нет имени
- `UpdatedBy *UserLink` — то же
- `DeletedBy *UserLink` — то же
- `SourceID *int` — нет названия источника
- `ID *int` — числовой ID файла (бесполезен рядом с UUID)

Ни один ID не резолвится в человекочитаемое значение.

### Что теряется по сравнению с тем, что SDK может вернуть

1. **Пагинация**: `ListFiles` отбрасывает `PageMeta` — LLM не знает есть ли следующая страница
2. **`upload` возвращает только 6 полей** из 20+: теряются `CreatedAt`, `HasMultipleVersions`, `DownloadLink`, `Previews`, `Type`
3. **`Previews []FilePreview`** — URL превью с размерами, никогда не используется
4. **`DownloadLink`** — присутствует в ответе `get`, но не гарантируется и не проверяется
5. **9 фильтров SDK** полностью недоступны для LLM

---

## Итого

**Приоритет рефакторинга: высокий**

Сервис-слой написан корректно. Проблемы в трёх местах: неполная tool-модель входа, потеря `PageMeta`, критическое расхождение документации и реализации по upload.

### Список конкретных изменений

1. **Исправить `tools_schema.md`** — привести схему в соответствие с реальными actions и полями `upload_params` (`local_path`, `file_name`, `with_preview`, `file_uuid`). Убрать `file_content` и `file_url`.
2. **Расширить `FileFilter`** — добавить `name`, `term`, `extensions []string`, `date_from`, `date_to`, `date_preset`
3. **Прокинуть новые фильтры в `ListFiles`** — добавить маппинг в SDK `FilesFilter`
4. **Возвращать `PageMeta` из `ListFiles`** — LLM должна знать о пагинации
5. **Резолвить `CreatedBy`/`UpdatedBy`/`DeletedBy` ID в имена** — через reference UsersService
6. **Добавить `with` в `GetFile`** — поддержка получения удалённых файлов по UUID
7. **Объединить `UUID`/`UUIDs` в delete** — единое поле `uuids []string`
8. **Прокинуть `CreatedByID`/`CreatedByType`** из tool-модели в `FileUploadParams`
