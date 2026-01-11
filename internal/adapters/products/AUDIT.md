# Audit: Products Service

Этот файл содержит результаты аудита папки `adapters/products/` на соответствие `tools_schema.md` и возможностям SDK.

---

## products.go

**Layer:** products

**Schema actions:** search, get, create, update, delete, get_by_entity, link, unlink, update_quantity

**SDK service:** `CatalogElementsService` (`core/adapters/catalog_elements.go`) + `LinksService`

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `SearchProducts` | Использует фильтр `Query`. |
| GetOne | ✅ | `GetProduct` | Вызывается без `with` параметров. |
| Create | ✅ | `CreateProducts` | Батч-операция. |
| Update | ✅ | `UpdateProducts` | Батч-операция. |
| Delete | ✅ | `DeleteProducts` | Цикл по одному (SDK не поддерживает батч-удаление элементов). |
| Link | ✅ | `LinkProduct` | Использует `LinksService`. |
| Unlink | ✅ | `UnlinkProduct` | Использует `LinksService`. |

**Gaps:**
- ✅ **Фильтрация**: `SearchProducts` теперь поддерживает фильтрацию по списку `IDs` (`CatalogElementsFilter.IDs`).
- ✅ **Связи**: `GetProductsByEntity` фильтрует только по `catalog_id`, но SDK `LinksService` может возвращать больше данных.

---

## Genkit Tool Handler (`app/gkit/tools/products.go`)

**Findings:**
- ✅ **Data Richness (Custom Fields)**: При действиях `create` и `update` данные маппятся через `json.Unmarshal` в модель SDK, что позволяет передавать любые поля, включая `custom_fields_values`.
- ✅ **Batch Create/Update**: Инструмент теперь поддерживает как одиночные операции (`data`), так и батч-операции (`items`).
- ✅ **Update Quantity**: Правильно реализовано через повторный `Link` (специфика v4 API).

**Статус:** ✅ Выполнено (архитектура приведена в соответствие с `catalogs`, обеспечена максимальная гибкость данных).

---

## Capabilities Coverage

**Filters:**
- ✅ SDK: `CatalogElementsFilter` поддерживает `IDs`, `Query`.
- ✅ Bot: Поддерживает `Query`, `IDs`, `Limit`, `Page`.

**With параметры:**
- ✅ SDK: Поддерживает `invoice_link`, `supplier_field_values`.
- ✅ Bot: Использует в методе `GetProduct`.

**Batch Operations:**
- ✅ SDK: Поддерживает Create/Update.
- ✅ Bot (Service): Реализовано.
- ✅ Bot (Tool): Реализовано.
