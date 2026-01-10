# Audit: Products Service

Этот файл содержит результаты аудита папки `services/products/` на соответствие `tools_schema.md` и возможностям SDK.

---

## products.go

**Layer:** products

**Schema actions:** search, get, create, update, delete, get_by_entity, link, unlink, update_quantity

**SDK service:** `CatalogElementsService` (`core/services/catalog_elements.go`) + `LinksService`

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
- **Фильтрация**: `SearchProducts` игнорирует возможность фильтрации по списку `IDs`, которую предоставляет SDK (`CatalogElementsFilter.IDs`).
- **Связи**: `GetProductsByEntity` фильтрует только по `catalog_id`, но SDK `LinksService` может возвращать больше данных.

---

## Genkit Tool Handler (`app/gkit/tools/products.go`)

**Findings:**
- ❌ **Data Loss (Custom Fields)**: При действиях `create` и `update` маппится только `Name`. Кастомные поля (в которых обычно лежат SKU, цены, описание товара) полностью игнорируются.
- ❌ **Batch Create/Update**: Несмотря на то, что сервис поддерживает массивы, инструмент обрабатывает только один элемент за раз.
- ✅ **Update Quantity**: Правильно реализовано через повторный `Link` (специфика v4 API).

**Статус:** ⚠️ Частично (функциональность связей хорошая, но работа с данными самого товара сильно урезана).

---

## Capabilities Coverage

**Filters:**
- ⚠️ SDK: `CatalogElementsFilter` поддерживает `IDs`, `Query`.
- ✅ Bot: Использует `Query`.
- ❌ Bot: Не использует `IDs`.

**With параметры:**
- ⚠️ SDK: Поддерживает `invoice_link`, `supplier_field_values`.
- ❌ Bot: Не использует.

**Batch Operations:**
- ✅ SDK: Поддерживает Create/Update.
- ✅ Bot (Service): Реализовано.
- ❌ Bot (Tool): Не используется.
