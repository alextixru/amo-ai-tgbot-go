# Аудит сервиса Files

Этот файл содержит результаты последовательного аудита папки `services/files/` на соответствие `tools_schema.md` и возможностям SDK.

---

## files.go
**Layer:** files
**Schema actions:** search, get, upload, delete
**SDK service:** FilesService (`core/services/files.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListFiles` | Поддерживает пагинацию и фильтр по UUID |
| GetOneByUUID | ✅ | `GetFile` | |
| UploadOne | ✅ | `UploadFile` | |
| UpdateOne | ❌ | — | Не реализовано в боте (переименование файла) |
| Delete | ⚠️ | `DeleteFile` | Реализовано как `DeleteOne`, батч-удаление (SDK) не используется |

**Genkit Tool Handler:**
- ✅ Инструмент `files` правильно распределяет вызовы.
- ❌ **Урезанная загрузка**: `upload_params` не позволяют передать `FileUUID` для загрузки новой версии файла или `CreatedBy` информацию.

**Статус:** ⚠️ Частично (основные действия есть, но нет обновления и расширенной загрузки)

### Capabilities Coverage

**Versioning:**
- ❌ SDK: `FileUploadParams` поддерживает `FileUUID` для обновления версии существующего файла.
- ❌ Bot: AI всегда создает новый файл, не может обновить версию.

**Parameters:**
- ⚠️ SDK: `Get()` поддерживает `with`: `deleted`, `unbilled`.
- ❌ Bot: Бот не предоставляет доступ к удаленным или неоплаченным файлам.
- ⚠️ SDK: `models.File` содержит `Previews`, `DownloadLink`, `Metadata` (mime, extension).
- ✅ Bot: Сервис возвращает модель целиком, AI видит ссылки на скачивание.

**Batch Operations:**
- ⚠️ SDK: Метод `Delete` принимает список файлов для массового удаления.
- ❌ Bot: Бот удаляет только по одному файлу.

---
