# Audit: Admin Pipelines Service

Этот файл содержит результаты аудита папки `services/admin_pipelines/` на соответствие `tools_schema.md` и возможностям SDK.

---

## pipelines.go

**Layer:** pipelines

**Schema actions:** list, search, get, create, update, delete

**SDK service:** `PipelinesService` (`core/services/pipelines.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListPipelines` | Без фильтров. |
| GetOne | ✅ | `GetPipeline` | |
| Create | ✅ | `CreatePipelines` | Батч-создание. |
| Update | ❌ | — | **SDK запрещает батч-обновление** (ErrNotAvailableForAction). |
| UpdateOne | ✅ | `UpdatePipeline` | |
| Delete | ❌ | — | **SDK запрещает батч-удаление**. |
| DeleteOne | ✅ | `DeletePipeline` | |

**Gaps:**
- Бот полностью следует архитектуре SDK, блокируя запрещенные API v4 батч-операции.
- В `ListPipelines` не передаются параметры (хотя для воронок их практически нет в API v4, кроме HATEOAS).

---

## statuses.go

**Layer:** statuses (в контексте воронки)

**Schema actions:** list_statuses, get_status, create_status, update_status, delete_status

**SDK service:** `StatusesService` (`core/services/statuses.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListStatuses` | |
| GetOne | ✅ | `GetStatus` | |
| Create | ✅ | `CreateStatus` | В сервисе по одному, но SDK поддерживает батч. |
| Update | ❌ | — | **SDK запрещает батч-обновление**. |
| UpdateOne | ✅ | `UpdateStatus` | |
| Delete | ❌ | — | **SDK запрещает батч-удаление**. |
| DeleteOne | ✅ | `DeleteStatus` | |

**Gaps:**
- **Батч-создание**: Сервис `StatusesService.Create` в боте принимает только один статус, хотя SDK метод `Create` работает со слайсом.

---

## Genkit Tool Handler (`app/gkit/tools/admin_pipelines.go`)

**Findings:**
- ✅ **Data Mapping**: Используется `json.Unmarshal` в модели `amomodels.Pipeline` и `amomodels.Status`. Это позволяет AI задавать цвета статусов (`color`), порядок сортировки (`sort`) и другие поля без ручного маппинга.
- ✅ **Обработка ID**: Корректная подстановка ID из `input.PipelineID` / `input.StatusID` в структуру данных, если они отсутствуют в `Data`.
- ✅ **Действия**: Поддерживаются все заявленные в схеме действия.

**Статус:** ✅ Хорошо (соответствует ограничениям платформы и предоставляет полный контроль над воронками)

---

## Capabilities Coverage

**Batch Operations:**
- ✅ SDK: `Pipelines.Create` поддерживает батч.
- ✅ Bot: `CreatePipelines` реализовано.
- ⚠️ SDK: `Statuses.Create` поддерживает батч.
- ❌ Bot: `CreateStatus` принимает один элемент.

**API Restrictions:**
- ✅ SDK: Запрещает `Update` и `Delete` (только `UpdateOne` / `DeleteOne`).
- ✅ Bot: Сервис не предоставляет эти методы, что предотвращает ошибки рантайма.
