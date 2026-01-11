# Audit: Admin Schema Service

Этот файл содержит результаты аудита папки `adapters/admin_schema/` на соответствие `tools_schema.md` и возможностям SDK.

---

## custom_fields.go

**Layer:** custom_fields

**Schema actions:** search, list, get, create, update, delete

**SDK service:** `CustomFieldsService` (`core/adapters/custom_fields.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListCustomFields` | Вызывается без фильтров. SDK поддерживает `IDs`, `Types`. |
| GetOne | ✅ | `GetCustomField` | |
| Create | ✅ | `CreateCustomFields` | Батч-создание. |
| Update | ✅ | `UpdateCustomFields` | Батч-обновление. |
| Delete | ✅ | `DeleteCustomField` | |

**Gaps:**
- ✅ **Фильтрация**: AI теперь может фильтровать поля по типу (IDs, Types) через `SchemaFilter`.

---

## field_groups.go

**Layer:** field_groups

**Schema actions:** search, list, get, create, update, delete

**SDK service:** `CustomFieldGroupsService` (`core/adapters/custom_field_groups.go`)

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

**Schema actions:** search, list, get, create, delete

**SDK service:** `LossReasonsService` (`core/adapters/loss_reasons.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListLossReasons` | Поддерживает пагинацию. |
| GetOne | ✅ | `GetLossReason` | |
| Create | ✅ | `CreateLossReasons` | |
| Update | ❌ | `UpdateLossReasons` | **REMOVED**. API v4 не поддерживает обновление. |
| Delete | ✅ | `DeleteLossReason` | |

**Gaps:**
- ✅ **Update Action**: Удален нерабочий метод.

---

## sources.go

**Layer:** sources

**Schema actions:** search, list, get, create, update, delete

**SDK service:** `SourcesService` (`core/adapters/sources.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListSources` | Вызывается без фильтров. SDK поддерживает `ExternalIDs`. |
| GetOne | ✅ | `GetSource` | |
| Create | ✅ | `CreateSources` | |
| Update | ✅ | `UpdateSources` | |
| Delete | ✅ | `DeleteSource` | |

**Gaps:**
- ✅ **Фильтрация**: AI теперь может искать источники по `external_id`.

---

## Genkit Tool Handler (`app/gkit/tools/admin_schema.go`)

**Findings:**
- ✅ **Data Mapping**: Используется `json.Unmarshal` прямо в модели SDK.
- ✅ **Loss Reasons Update**: Удален нерабочий `update`.
- ✅ **Filtering**: Добавлены мапперы фильтров для всех слоев.

**Статус:** ✅ Готов
