# Аудит сервиса Files

Этот файл содержит результаты последовательного аудита папки `adapters/files/` на соответствие `tools_schema.md` и возможностям SDK.

---

## files.go
**Layer:** files
**Schema actions:** search, get, upload, delete
**SDK service:** FilesService (`core/adapters/files.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListFiles` | Поддерживает пагинацию, UUID и **With** (deleted, unbilled) |
| GetOneByUUID | ✅ | `GetFile` | |
| UploadOne | ✅ | `UploadFile` | Поддерживает **FileUUID** для новых версий |
| UpdateOne | ✅ | `UpdateFile` | Реализовано (переименование файла) |
| Delete | ✅ | `DeleteFiles` | Поддерживает массив UUID для батч-удаления |

**Genkit Tool Handler:**
- ✅ Инструмент `files` правильно распределяет вызовы.
- ✅ **Расширенная загрузка**: `upload_params` теперь позволяют передать `FileUUID` для загрузки новой версии файла.
- ✅ **Batch & Update**: Добавлены действия `update` и `delete` (batch).

**Статус:** ✅ Выполнено (полная поддержка SDK и Genkit)

### Capabilities Coverage

**Versioning:**
- ✅ SDK: `FileUploadParams` поддерживает `FileUUID` для обновления версии существующего файла.
- ✅ Bot: AI может передать `file_uuid` для создания новой версии.

**Parameters:**
- ✅ SDK: `Get()` поддерживает `with`: `deleted`, `unbilled`.
- ✅ Bot: Поддерживается через поле `With` в фильтре.
- ⚠️ SDK: `models.File` содержит `Previews`, `DownloadLink`, `Metadata` (mime, extension).
- ✅ Bot: Сервис возвращает модель целиком, AI видит ссылки на скачивание.

**Batch Operations:**
- ✅ SDK: Метод `Delete` принимает список файлов для массового удаления.
- ✅ Bot: Поддерживается через массив `UUIDs` в `FilesInput`.

---
