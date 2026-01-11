#Audit: Catalogs Service

Этот файл содержит результаты аудита папки `adapters/catalogs/` на соответствие `tools_schema.md` и возможностям SDK.

---

##catalogs.go

**Layer:** catalogs

**Schema actions:** list, get, create, update

**SDK service:**`CatalogsService` (`core/adapters/catalogs.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |

|-----------|----------------------|----------------|-------------|

| Get | ✅ |`ListCatalogs`| Бот вызывает без фильтров. SDK поддерживает `filters.CatalogsFilter`|

| GetOne | ✅ |`GetCatalog`||

| Create | ✅ |`CreateCatalogs`||

| Update | ✅ |`UpdateCatalogs`||

**Gaps:**

- Бот не поддерживает фильтрацию каталогов (хотя их обычно мало).

---

##elements.go

**Layer:** catalogs (elements)

**Schema actions:** list_elements, get_element, create_element, update_element, link_element, unlink_element

**SDK service:**`CatalogElementsService` (`core/adapters/catalog_elements.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |

|-----------|----------------------|----------------|-------------|

| Get | ✅ |`ListElements`| ✅ Поддерживает `CatalogElementsFilter` (Query, IDs, Page, Limit) |

| GetOne | ✅ |`GetElement`| ✅ Поддерживает `with` параметры (invoice_link, supplier_field_values) |

| Create | ✅ |`CreateElements`||

| Update | ✅ |`UpdateElements`||

| Link | ✅ |`LinkElement`| ✅ Доступен в Genkit tool как `link_element` |

| Unlink | ✅ |`UnlinkElement`| ✅ Доступен в Genkit tool как `unlink_element` |

**Gaps:**

- ~~**Фильтры**: SDK поддерживает `filters.CatalogElementsFilter` (Query, IDs). Бот их не использует.~~ ✅ Исправлено
- ~~**With параметры**: SDK поддерживает `invoice_link`, `supplier_field_values`. Бот их игнорирует.~~ ✅ Исправлено
- ~~**Связи**: В Genkit tool отсутствуют действия для связывания элементов с лидами/покупателями.~~ ✅ Исправлено

---

##Genkit Tool Handler (`app/gkit/tools/catalogs.go`)

**Findings:**

- ✅ Базовый CRUD для каталогов и элементов присутствует.
- ✅ **Link/Unlink**: Действия `link_element` и `unlink_element` зарегистрированы в инструменте.
- ✅ **Поиск**: AI может передать `query` и `ids` для поиска элементов.
- ✅ **With-параметры**: AI может запросить `invoice_link`, `supplier_field_values` через `filter.with`.
- ⚠️ **Data mapping**: Используется `json.Unmarshal` в `amomodels.CatalogElement`. Это лучше, чем ручное маппирование, но требует точности от AI в именовании полей (особенно `custom_fields_values`).

**Статус:** ✅ Полнофункционален

