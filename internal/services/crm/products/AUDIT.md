# Аудит: products

## Вход (LLM → tool → сервис)

### Поля, которые сейчас числовые ID, но должны быть именами

**`product.price_id int`**
LLM не знает числовые ID ценовых полей каталога без отдельного запроса. Должно быть опциональным: если не передан — брать первое ценовое поле по умолчанию.
Сейчас: `{"product": {"id": 123, "quantity": 2, "price_id": 789}}`

**`ProductData.custom_fields_values map[string]any`**
На деле маршалится в `[]CustomFieldValue` где `field_id int` — числовой. LLM не знает `field_id` для SKU, описания и т.д. без запроса схемы. Нет поддержки `field_code` в input-модели.

**`entity.id int`**
Сам по себе нормален, но `get_by_entity` возвращает `EntityLink` с числовыми `ToEntityID` без имён товаров — цикл замыкается в числах.

### Поля, которые отсутствуют, но нужны LLM

- **`ProductFilter.With` не прокидывается в `SearchProducts`** — `f` строится без `with`, параметр полностью игнорируется при поиске
- **Нет `sort`/`order` в `ProductFilter`** — нет сортировки по имени или дате
- **`update_quantity` неотличим от `link`** — обе операции вызывают `LinkProduct`, ответ `nil` не говорит что это обновление
- **`ProductData` не содержит `id`** — batch update через `items []ProductData` нерабочий: нет способа передать ID обновляемого товара

### Неудобные структуры/типы

- **`With []string` в `ProductFilter`** — для action `get` LLM вынуждена передавать `with` внутри `filter`. Семантически `with` — не фильтр, а параметр представления.
- **`IDs []int` дублируется** — есть в `ProductsInput` (для delete) и `ProductFilter` (для search). Два разных места для похожей сущности.
- **`sync.Once` глобальный стейт для `productsCatalogID`** — глобальные переменные делают сервис непригодным для тестирования и нескольких SDK-инстансов.

---

## Выход (сервис → tool → LLM)

### Что сейчас возвращается

| Action | Тип ответа |
|--------|-----------|
| `search` | `[]*models.CatalogElement` |
| `get` | `*models.CatalogElement` |
| `create` | `[]*models.CatalogElement` |
| `update` | `[]*models.CatalogElement` |
| `delete` | `nil` |
| `get_by_entity` | `[]models.EntityLink` (только ID, без данных товара!) |
| `link` / `unlink` / `update_quantity` | `nil` |

### Какие SDK With-параметры не используются

`CatalogElement.AvailableWith()` возвращает: `["invoice_link", "supplier_field_values"]`

- **`supplier_field_values`** — поля поставщика (артикул, себестоимость) — никогда не запрашиваются при поиске
- **`invoice_link`** — доступен только при явной передаче в `get`, для search не работает вообще

В `SearchProducts` `f` строится из `ProductFilter` без `With` — параметр не прокидывается в `url.Values`.

### Числовые ID в ответе, которые LLM не может интерпретировать

- `CatalogID int` — числовой ID каталога
- `CurrencyID int` — LLM не знает что за валюта
- `CreatedBy int` / `UpdatedBy int` — без имён пользователей
- `AccountID int` — избыточно
- `PriceID int` — числовой ID ценового поля
- `CustomFieldsValues[].FieldID int` — `FieldName` заполнен только если API вернул, не гарантировано
- `CustomFieldsValues[].Values[].EnumID int` — без текстового значения
- **`get_by_entity` → `EntityLink.ToEntityID int`** — только ID товара, без имени, цены, SKU

### Что теряется по сравнению с тем, что SDK может вернуть

1. **`get_by_entity`** — `[]EntityLink` без обогащения: LLM получает список ID товаров без имён и цен. Сервис мог бы загрузить детали через `GetProduct` и вернуть `[]ProductWithLink{Link, Product}`
2. **`supplier_field_values`** — артикул поставщика, себестоимость теряются при bulk-поиске
3. **`PageMeta`** — отбрасывается везде, LLM не знает есть ли следующая страница
4. **`delete` и `link`/`unlink`** — возвращают `nil`, нет подтверждения операции
5. **`InvoiceWarning`** — никогда не обрабатывается в ответе

---

## Итого

**Приоритет рефакторинга: высокий**

Pass-through к SDK без адаптации. `get_by_entity` возвращает бесполезные числовые ID вместо данных товаров. `With` при поиске не работает. Batch update нерабочий.

### Список конкретных изменений

1. **[КРИТИЧНО] Обогатить `get_by_entity`** — после `[]EntityLink` загрузить детали каждого товара и вернуть `[]ProductWithLink{Link, Product}` с именами, полями и ценой
2. **[КРИТИЧНО] Добавить `field_code string` для custom fields** — заменить `map[string]any` на `[]ProductFieldInput{FieldCode string; Value any}`, маппинг в `field_id` делает сервис
3. **[ВЫСОКИЙ] Прокинуть `With` в `SearchProducts`** — добавить `with` в `url.Values` при построении запроса
4. **[ВЫСОКИЙ] Резолвить `currency_id` → `currency_code`** — через CurrenciesService (уже в контексте)
5. **[ВЫСОКИЙ] Резолвить `created_by`/`updated_by` в имена** — через UsersService
6. **[СРЕДНИЙ] Добавить `id` в `ProductData`** — чтобы batch update работал корректно
7. **[СРЕДНИЙ] Вынести `with []string` на верхний уровень `ProductsInput`**
8. **[СРЕДНИЙ] Вернуть `PageMeta` в ответ `search`** — `{items, total, page, has_more}`
9. **[СРЕДНИЙ] Сделать `price_id` опциональным при `link`** — по умолчанию первое ценовое поле
10. **[НИЗКИЙ] Возвращать подтверждение для `link`/`unlink`/`update_quantity`/`delete`** — `{ok: true, action, product_id, entity_id}`
11. **[НИЗКИЙ] Исправить typo**: `"привзяка"` → `"привязка"` в описании tool
12. **[НИЗКИЙ] Убрать глобальный `sync.Once`** — перенести `catalogIDOnce` и `productsCatalogID` в поля `service` struct
