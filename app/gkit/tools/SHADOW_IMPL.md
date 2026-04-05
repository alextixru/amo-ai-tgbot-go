# Shadow Tools — Итоги реализации

## Статус

Все 12 tools переведены на Shadow Tool паттерн.

## Паттерн по каждому tool

| Tool | Input type | available_values |
|------|-----------|-----------------|
| entities | `any` + `WithInputSchema` (минимальная схема: entity_type + action) | pipelines, statuses, users, loss_reasons, custom_field_codes |
| activities | `map[string]any` | users (для tasks/subscriptions), entity_types |
| complex_create | `any` | pipelines, statuses, users |
| products | `map[string]any` | product_names (первые 20 из SearchProducts) |
| catalogs | `map[string]any` | catalog_names |
| files | `map[string]any` | нет |
| unsorted | `map[string]any` | pipelines, statuses, users |
| customers | `map[string]any` | users, customer_statuses |
| admin_schema | `any` | нет |
| admin_pipelines | `map[string]any` | нет |
| admin_users | `map[string]any` | нет |
| admin_integrations | `any` | нет |

## Важные нюансы

**`isSchemaMode` — пакетная функция в unsorted.go**
Используется также в catalogs.go. Остальные tools реализуют свою логику локально (разные имена: `entitiesIsSchemaMode`, `adminPipelinesIsSchemaMode` и т.д.).

**`entities` — единственный с `WithInputSchema`**
Остальные используют `map[string]any` как тип → genkit генерирует permissive `{type: object}` схему без полей. Оба подхода дают одинаковый эффект для LLM.

**`search`/`list` — всегда Execute mode**
Во всех tools action без обязательных полей сразу уходит в execute, не возвращает схему.

**JSON roundtrip**
Execute mode во всех tools: `map[string]any → json.Marshal → json.Unmarshal → OriginalInputStruct`. Числа из JSON приходят как `float64` — проверки `== 0` и `!= float64(0)` разбросаны по коду.

**`complex_create` объединил два tool в один**
`complex_create_batch` удалён как отдельный tool, стал action `create_batch` внутри `complex_create`. Registry.RegisterAll() обновлён.

**`unsorted` — нормализация плоского формата**
В execute mode `executeUnsorted` нормализует плоский input в вложенный до roundtrip (поля `link_data`, `create_data`).

**`products` — dynamic available_values через API**
Единственный tool где schema response делает реальный API-вызов (`SearchProducts limit=20`).

**admin_* — без available_values**
Сервисы сами возвращают ошибку с перечислением доступных значений если имя не найдено.

## Новые методы в сервисах

| Сервис | Добавлено |
|--------|-----------|
| entities | `StatusesByPipeline()`, `LossReasonNames()`, `CustomFieldCodes(entityType)` |
| complex_create | `StatusesByPipeline()` |
| unsorted | `StatusNames() map[string][]string` |
