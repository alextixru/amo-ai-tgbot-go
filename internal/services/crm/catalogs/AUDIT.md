# Аудит: catalogs

## Вход (LLM → tool → сервис)

### Поля, которые сейчас числовые ID, но должны быть именами

**`catalog_id int`**
LLM не знает числовой ID каталога — она знает имя ("Товары", "Услуги", "Счета"). Вынуждена сначала вызывать `list`, запоминать `id`, затем передавать в следующем вызове.
Нужно: `catalog_name string` как альтернатива, резолвинг в сервисе.

**`element_id int`**
LLM знает имя элемента ("iPhone 15 Pro"), но не его числовой ID. Для `get_element`, `update_element`, `link_element`, `unlink_element` нужен предварительный `list_elements`.

**`link_data.entity_id int`**
При привязке элемента к сущности нужен числовой ID. Нет подсказки откуда его брать.

**`custom_fields_values map[string]any` в `CatalogElementData`**
LLM не знает какие ключи использовать (`field_id`? `field_code`?). Конвертация через `json.Marshal → json.Unmarshal` в `[]CustomFieldValue` требует знания внутренней структуры SDK.

### Поля, которые отсутствуют, но нужны LLM

- **Нет `catalog_name string`** — поиск каталога по имени без предварительного `list`
- **Нет фильтра для `ListCatalogs`** — SDK имеет `CatalogsFilter.Type` (regular/invoices/products), в `CatalogsInput` нет `filter` для `list`
- **Нет `with` для `list_elements`** — `With []string` в `CatalogFilter` работает только для `get_element`, не для `list_elements`
- **Нет `delete` и `delete_element`** — в интерфейсе и tool handler отсутствуют, нарушает CRUD
- **`CatalogID` не прокидывается при `update`** — `input.CatalogID` не попадает в `models.Catalog`. Баг.

### Неудобные структуры/типы

- **Баг marshal/unmarshal**: `json.Marshal(input.Data)` → `json.Unmarshal(data, &catalogs)` не работает когда `input.Data` — одиночный объект, а target — slice. `create` и `update` нерабочие.
- **`CatalogFilter` смешивает фильтрацию и `with`** — `IDs`, `Query`, `Page`, `Limit` и `With []string` в одной структуре. Разные семантики.
- **`LinkData.Metadata map[string]any`** — нет документации о допустимых ключах (`quantity`, `price_id`)

---

## Выход (сервис → tool → LLM)

### Что сейчас возвращается

| Метод | Тип возврата |
|-------|-------------|
| `ListCatalogs` | `[]*models.Catalog` |
| `GetCatalog` | `*models.Catalog` |
| `CreateCatalogs` / `UpdateCatalogs` | `[]*models.Catalog` |
| `ListElements` | `[]*models.CatalogElement` |
| `GetElement` | `*models.CatalogElement` |
| `CreateElements` / `UpdateElements` | `[]*models.CatalogElement` |
| `LinkElement` / `UnlinkElement` | `{"success": true}` |

Сырые SDK-модели без обёртки. Числовые ID, Unix timestamps.

### Какие SDK With-параметры не используются

`CatalogElement.AvailableWith()` → `["invoice_link", "supplier_field_values"]`

- **`supplier_field_values`** — поля поставщика, себестоимость — **не запрашиваются нигде по умолчанию**
- **`invoice_link`** — ссылка на счёт — только при явном `with` в `get_element`, для `list_elements` недоступно
- **`CatalogsService.Get`** принимает `*filters.CatalogsFilter`, в `ListCatalogs` передаётся `nil`

### Числовые ID в ответе, которые LLM не может интерпретировать

**`models.Catalog`:**
- `created_by int` / `updated_by int` — без имён пользователей
- `account_id int` — технический шум
- `created_at int64` / `updated_at int64` — Unix timestamps

**`models.CatalogElement`:**
- `catalog_id int` — LLM не знает имя каталога
- `currency_id int` — LLM не знает код валюты (RUB, USD)
- `created_by int` / `updated_by int` — без имён
- `created_at int64` / `updated_at int64` — Unix timestamps
- `price_id int` — LLM не понимает смысл
- `account_id int` — технический шум
- `CustomFieldsValues[].enum_id int` — без текстового значения (есть `enum_code`, но не всегда заполнен)

### Что теряется по сравнению с тем, что SDK может вернуть

1. **`supplier_field_values`** — не запрашиваются при `list_elements`
2. **`invoice_link`** — не запрашивается при листинге
3. **Пагинация** — `*PageMeta` отбрасывается везде
4. **`CatalogsFilter.Type`** — фильтрация по типу недоступна через tool

---

## Итого

**Приоритет рефакторинга: высокий**

Сервис транслирует сырые SDK-объекты без нормализации. Два критических бага: marshal/unmarshal ломает `create`/`update`, `CatalogID` не прокидывается при `update`.

### Список конкретных изменений

1. **[КРИТИЧНО] Исправить баг marshal/unmarshal** — обернуть явно: `[]*models.Catalog{mapToModel(input.Data)}` вместо двойного JSON roundtrip
2. **[КРИТИЧНО] Прокидывать `CatalogID` в объект при `update`** — явно `catalog.ID = input.CatalogID`
3. **[ВЫСОКИЙ] Добавить резолвинг `catalog_name → catalog_id`** — принимать `catalog_name string`, резолвить через `ListCatalogs`
4. **[ВЫСОКИЙ] Добавить пагинацию в ответ** — `{items, total, page, has_more}`
5. **[ВЫСОКИЙ] Нормализовать числовые поля в ответе** — timestamps → RFC3339, `created_by`/`updated_by` → имена, `currency_id` → код валюты
6. **[СРЕДНИЙ] Добавить `filter` для `ListCatalogs`** — `type string` (regular/invoices/products)
7. **[СРЕДНИЙ] Автоматически запрашивать `with`** — в `GetElement` всегда `with=invoice_link,supplier_field_values`; для `ListElements` добавить поддержку через `url.Values`
8. **[СРЕДНИЙ] Вынести `with []string` из `CatalogFilter` на верхний уровень `CatalogsInput`**
9. **[СРЕДНИЙ] Добавить `delete` и `delete_element`**
10. **[НИЗКИЙ] Заменить `custom_fields_values map[string]any`** на `[]CatalogFieldValue{FieldCode string; Value any}`
11. **[НИЗКИЙ] Документировать `link_data.metadata`** — явные ключи: `quantity float64`, `price_id int`
12. **[НИЗКИЙ] Убрать `account_id`, `_links` из ответа** — технический шум для LLM
