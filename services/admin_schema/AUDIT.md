# Audit: Admin Schema Service

Этот файл содержит результаты аудита папки `services/admin_schema/` на соответствие `tools_schema.md` и возможностям SDK.

---

## custom_fields.go

**Layer:** custom_fields

**Schema actions:** search, list, get, create, update, delete

**SDK service:** `CustomFieldsService` (`core/services/custom_fields.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListCustomFields` | Вызывается без фильтров. SDK поддерживает `IDs`, `Types`. |
| GetOne | ✅ | `GetCustomField` | |
| Create | ✅ | `CreateCustomFields` | Батч-создание. |
| Update | ✅ | `UpdateCustomFields` | Батч-обновление. |
| Delete | ✅ | `DeleteCustomField` | |

**Gaps:**
- **Фильтрация**: Бот не позволяет AI фильтровать поля по типу (например, найти все поля типа `multiselect`). Это было бы полезно для упрощения настройки CRM.

---

## field_groups.go

**Layer:** field_groups

**Schema actions:** search, list, get, create, update, delete

**SDK service:** `CustomFieldGroupsService` (`core/services/custom_field_groups.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListFieldGroups` | |
| GetOne | ✅ | `GetFieldGroup` | |
| Create | ✅ | `CreateFieldGroups` | |
| Update | ✅ | `UpdateFieldGroups` | |
| Delete | ✅ | `DeleteFieldGroup` | Использует `string` ID (верно для групп). |

---

## loss_reasons.go

**Layer:** loss_reasons

**Schema actions:** search, list, get, create, update, delete

**SDK service:** `LossReasonsService` (`core/services/loss_reasons.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListLossReasons` | |
| GetOne | ✅ | `GetLossReason` | |
| Create | ✅ | `CreateLossReasons` | |
| Update | ❌ | `UpdateLossReasons` | **SDK возвращает ErrNotAvailableForAction**. В боте метод есть, но он не будет работать. |
| Delete | ✅ | `DeleteLossReason` | |

**Gaps:**
- ❌ **Update Action**: API v4 не поддерживает обновление причин отказа. Бот предоставляет это действие в Genkit, что приведет к ошибке при вызове.

---

## sources.go

**Layer:** sources

**Schema actions:** search, list, get, create, update, delete

**SDK service:** `SourcesService` (`core/services/sources.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListSources` | Вызывается без фильтров. SDK поддерживает `ExternalIDs`. |
| GetOne | ✅ | `GetSource` | |
| Create | ✅ | `CreateSources` | |
| Update | ✅ | `UpdateSources` | |
| Delete | ✅ | `DeleteSource` | |

**Gaps:**
- **Фильтрация**: Нет возможности искать источники по `external_id`.

---

## Genkit Tool Handler (`app/gkit/tools/admin_schema.go`)

**Findings:**
- ✅ **Data Mapping**: Используется `json.Unmarshal` прямо в модели SDK (`amomodels.CustomField`, `amomodels.Source`). Это позволяет AI передавать всю сложность структур (enums, settings), если он знает их структуру.
- ❌ **Loss Reasons Update**: Инструмент позволяет вызывать `update` для `loss_reasons`, что гарантированно упадет с ошибкой от SDK.
- ⚠️ **EntityType**: Для полей и групп `entity_type` обязателен. Бот это проверяет.

**Статус:** ⚠️ Частично (CRUD есть, но фильтры отсутствуют, а в `loss_reasons` есть нерабочее действие)
