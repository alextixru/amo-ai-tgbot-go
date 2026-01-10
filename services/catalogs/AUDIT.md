#Audit: Catalogs Service

Этот файл содержит результаты аудита папки `services/catalogs/` на соответствие `tools_schema.md` и возможностям SDK.

---

##catalogs.go

**Layer:** catalogs

**Schema actions:** list, get, create, update

**SDK service:**`CatalogsService` (`core/services/catalogs.go`)

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

**Schema actions:** list_elements, get_element, create_element, update_element

**SDK service:**`CatalogElementsService` (`core/services/catalog_elements.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |

|-----------|----------------------|----------------|-------------|

| Get | ✅ |`ListElements`| Бот вызывает без фильтров (`nil`). |

| GetOne | ✅ |`GetElement`| Бот вызывает с пустыми `url.Values{}`, игнорируя `with`. |

| Create | ✅ |`CreateElements`||

| Update | ✅ |`UpdateElements`||

| Link | ✅ |`LinkElement`| Реализовано в сервисе, но нет в Genkit tool. |

| Unlink | ✅ |`UnlinkElement`| Реализовано в сервисе, но нет в Genkit tool. |

**Gaps:**

-**Фильтры**: SDK поддерживает `filters.CatalogElementsFilter` (Query, IDs). Бот их не использует, что делает поиск конкретного товара через AI невозможным (только полный листинг).

-**With параметры**: SDK поддерживает `invoice_link`, `supplier_field_values`. Бот их игнорирует.

-**Связи**: В Genkit tool отсутствуют действия для связывания элементов с лидами/покупателями.

---

##Genkit Tool Handler (`app/gkit/tools/catalogs.go`)

**Findings:**

- ✅ Базовый CRUD для каталогов и элементов присутствует.
- ❌ **Link/Unlink**: Действия `link_element` и `unlink_element` не зарегистрированы в инструменте.
- ❌ **Поиск**: Нет возможности передать `query` для поиска элементов. AI вынужден выкачивать весь список.
- ⚠️ **Data mapping**: Используется `json.Unmarshal` в `amomodels.CatalogElement`. Это лучше, чем ручное маппирование, но требует точности от AI в именовании полей (особенно `custom_fields_values`).

**Статус:** ⚠️ Частично (основной CRUD есть, функционал связей и поиска отсутствует)
