# Аудит toolmap.md

## AccountService ⚠️

### Отсутствует в SDK:
- `AvailableWith` — метод упомянут в toolmap, но отсутствует в `account.go`

### Соответствует:
- `GetCurrent(ctx, with)`

---

## CallsService ⚠️

### Неточности в названиях:
- в toolmap: `Add`, `AddOne`
- в SDK: `Create` (унаследовано от BaseEntityService), `CreateOne`

### Отсутствует в toolmap:
- Метод `Create` (batch) доступен через наследование, но в toolmap назван `Add`.

---

## CatalogElementsService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`
- `GetOne`: принимает дополнительный `params url.Values`
- `Create`, `Update`: возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`

### Отсутствует в toolmap:
- `Link(ctx, elementID, entityType, entityID, metadata)`
- `Unlink(ctx, elementID, entityType, entityID)`
- Helper methods: `SetCatalogId`, `GetCatalogId`, `ValidateEntityId`

---

## CatalogsService ⚠️

### ~~Отсутствует в SDK:~~
- ~~`GetOne` — метод есть в toolmap, но отсутствует в `catalogs.go`~~
- **ИСПРАВЛЕНО:** `GetOne` наследуется от `BaseEntityService` — метод доступен

### Неточности в сигнатурах:
- `Get`, `Create`, `Update`: возвращают `*PageMeta`, в toolmap этого нет.

---

## ChatTemplatesService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `GetOne`: в toolmap указан аргумент `with []string`, но метод наследуется от `BaseEntityService` и принимает `params url.Values`.
- `Create`, `Update`: наследуются, возвращают `*PageMeta`, в toolmap пропущен.

### Соответствует:
- `Delete`, `DeleteMany`, `SendOnReview`, `UpdateReviewStatus`

---

## CompaniesService ⚠️

### Отсутствует в toolmap:
- `CreateOne(ctx, company)`
- `UpdateOne(ctx, company)`
- `SyncOne(ctx, company, with)`
- `GetLinks(ctx, companyID)`

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `Create`, `Update`: возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `GetOne`: в toolmap указан, но в `companies.go` он явно не определен (используется наследуемый от `BaseEntityService`, который принимает `params url.Values` вместо `with []string`).

---

## ContactsService ⚠️

### Отсутствует в toolmap:
- `CreateOne(ctx, contact)`
- `UpdateOne(ctx, contact)`
- `SyncOne(ctx, contact, with)`
- `GetLinks(ctx, contactID)`
- `GetChats(ctx, contactID)`
- `LinkChats(ctx, links)`

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `Create`, `Update`: возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `GetOne`: не переопределен в `contacts.go`, используется из `BaseEntityService`. Соответственно, аргумент `with []string` неверен, должно быть `params url.Values`.

### Соответствует:
- `Link`, `Unlink`

---

## CurrenciesService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.

---

## CustomFieldGroupsService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `Create`, `Update`: наследуются от `BaseEntityTypeService`, возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `GetOne`: не определен явно, наследуется. Принимает `params url.Values` (как `BaseEntityTypeService`), а не просто `id string`.

### Соответствует:
- `Delete`

---

## CustomFieldsService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `Create`, `Update`: возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `GetOne`: принимает дополнительный `params url.Values`, которого нет в toolmap.

### Соответствует:
- `Delete`

---

## CustomerBonusPointsService ✅
Соответствует

---

## CustomerStatusesService ⚠️

### Неточности в сигнатурах:
- `Get`: в toolmap указаны аргументы `page, limit int`, но метод наследуется от `BaseEntityService` и принимает `params url.Values`. Также возвращает `*PageMeta`, что пропущено в toolmap.
- `Create`, `Update`: наследуются, возвращают `*PageMeta`, в toolmap пропущен.
- `GetOne`: не определён явно (наследуется), принимает дополнительный `params url.Values` вместо `id int`.

### Соответствует:
- `Delete`

---

## CustomerTransactionsService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.

### Соответствует:
- `Create`, `Delete`

---

## CustomersService ⚠️

### Отсутствует в toolmap:
- `SetMode(ctx, mode, isEnabled)`

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `Create`, `Update`: возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `GetOne`: в toolmap указан `with []string`, но метод наследуется и принимает `params url.Values`.

### Соответствует:
- `Delete`, `Link`

---

## EntityFilesService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `Link`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.

### Соответствует:
- `Unlink`

---

## EntitySubscriptionsService ✅
Соответствует

---

## EventTypesService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.

---

## EventsService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `GetOne`: наследуется от BaseEntityService, принимает `params url.Values`, а не `id int`.

---

## FilesService ⚠️

### Неточности в названиях:
- в toolmap: `GetOne` (принимает uuid), в SDK: `GetOneByUUID`.
- в toolmap: `Delete` (принимает uuid), в SDK: `DeleteOne`. Метод `Delete` в SDK принимает `FilesCollection`.

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.

### Отсутствует в toolmap:
- `UpdateOne`

### Соответствует:
- `UploadOne`

---

## LeadsService ⚠️

### Отсутствует в SDK:
- `GetWithPagination`: функционал пагинации встроен в `Get`, который возвращает `PageMeta`.

### Отсутствует в toolmap:
- `SyncOne(ctx, lead, with)`
- `AddComplexInPlace(ctx, leads)`

### Неточности в сигнатурах:
- `Get`, `Create`, `Update`: возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `GetOne`: в toolmap указан `with []string`, но метод наследуется и принимает `params url.Values`.

### Соответствует:
- `CreateOne`, `UpdateOne`, `Delete`, `Link`, `Unlink`, `AddComplex`, `AddOneComplex`

---

## NotesService ⚠️

### Неточности в сигнатурах:
- `Get`, `GetByParent`, `Create`, `Update`: возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `GetOne`: принимает `params url.Values`, а не просто `id int`.

---

## PipelinesService ⚠️

### ~~Отсутствует в SDK:~~
- ~~`GetStatuses`, `GetStatus`, `CreateStatus`, `UpdateStatus`, `DeleteStatus`: методы для работы со статусами внутри воронки не реализованы в этом сервисе.~~
- **ИСПРАВЛЕНО:** `StatusesService` существует как отдельный сервис. Accessor `sdk.Statuses(pipelineID)` добавлен.
- `Update` (batch), `Delete` (batch): методы существуют, но явно возвращают ошибку `ErrNotAvailableForAction`. В toolmap они указаны как рабочие.

### Отсутствует в toolmap:
- `UpdateOne`, `DeleteOne`

### Неточности в сигнатурах:
- `Get`, `Create`: возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `GetOne`: принимает `params url.Values`, а не просто `id int`.

---

## ProductsService ⚠️

### Отсутствует в SDK:
- `Get`, `GetOne`, `Create`, `Update`, `Delete`: эти методы для работы с товарами (элементами каталога) не реализованы в `ProductsService`. В Go SDK (и PHP SDK) для этого используется `CatalogElementsService` с ID каталога товаров.

### Отсутствует в toolmap:
- `Settings(ctx)`: для получения настроек товаров.
- `UpdateSettings(ctx, settings)`: для обновления настроек товаров.

---

## RolesService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `Create`, `Update`: наследуются, возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `GetOne`: наследуется, принимает `params url.Values`, а не `with []string`.

---

## SegmentsService ⚠️

### Отсутствует в SDK:
- `GetOne`: наследуется от `BaseEntityService`, но в `SegmentsService` он не переопределен. Однако `BaseEntityService.GetOne` принимает `params url.Values`, а не `id`.
- `Update`: существует, но возвращает `ErrNotAvailableForAction`. В toolmap указан как рабочий метод.

### Неточности в сигнатурах:
- `Get`, `Create`: возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.

---

## ShortLinksService ⚠️

### Отсутствует в toolmap:
- `CreateOne`: метод реализован в SDK, но нет в карта.
- `Update`:  существует в SDK, но возвращает ошибку `ErrNotAvailableForAction`. В toolmap не упомянут (что верно, так как не работает), но стоит отметить.

### Неточности в сигнатурах:
- `Get`, `Create`: возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.

---

## SourcesService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `Create`, `Update`: возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `GetOne`: присутствует в интерфейсе и наследуется от `BaseEntityService`, принимает `params url.Values`, но в toolmap указан `id int`.

---

## TagsService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, SDK-метод принимает `filter *filters.TagsFilter`, в toolmap пропущен `PageMeta`.
- `Create`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.

### Неточности в поведении:
- `Update`: присутствует в SDK, но работает *только* для leads. В toolmap метод отсутствует.
- `Delete`: присутствует в SDK, но возвращает ошибку `ErrTagsDeleteNotSupported`, так как API v4 не поддерживает прямое удаление тегов. В toolmap метод указан как рабочий.

---

## TalksService ⚠️

### Неточности в сигнатурах:
- `Close`: принимает `(ctx, talkID, opts *TalkCloseOptions)`, а в toolmap только `talkID`.
- `GetOne`: присутствует в SDK (наследуется), но отсутствует в toolmap.
- `Get`: присутствует в toolmap как рабочий метод, но в SDK возвращает `ErrNotAvailableForAction`.

### Отсутствует в SDK:
- `Get`: заблокирован.

---

## TasksService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `Create`, `Update`: возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `GetOne`: наследуется от `BaseEntityService`, принимает `params url.Values`, но в toolmap указан `id int`.

### Отсутствует в toolmap:
- `CreateOne`: существует в SDK.
- `UpdateOne`: существует в SDK.

---

## UnsortedService ⚠️

### Неточности в сигнатурах:
- `Get`, `Create`: возвращают `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `Create`: в SDK принимает `(ctx, category, unsorted)`, в toolmap так же, но возвращает `PageMeta`.
- `Summary`: в SDK принимает `filter *filters.UnsortedSummaryFilter`, в toolmap указан `*UnsortedFilter`.

### Отсутствует в toolmap:
- `GetOne`: присутствует в SDK (наследуются от BaseEntityService), но отсутствует в toolmap.

---

## UsersService ⚠️

### Отсутствует в SDK:
- `AddToGroup`: метод упомянут в toolmap, но отсутствует в `UsersService`.
- `GetRoles`, `GetRole`: методы упомянуты в toolmap, но отсутствуют в `UsersService`.

### Отсутствует в SDK (заблокировано):
- `Update`, `UpdateOne`: методы существуют, но возвращают ошибку.

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.

---

## WebhooksService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.

---

## WebsiteButtonsService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.

### Отсутствует в toolmap:
- `Create`, `Update`, `AddOnlineChat`: синхронные методы-обертки (существуют в SDK). В toolmap указаны только Async версии.

---

## WidgetsService ⚠️

### Неточности в сигнатурах:
- `Get`: возвращает `(..., *PageMeta, error)`, в toolmap пропущен `PageMeta`.
- `Install`: в toolmap принимает `code string`. В SDK метод `Install` принимает `*models.Widget`. Метод, принимающий `code string`, называется `InstallByCode`.

### Неточности в названиях:
- `GetOne`: в toolmap называется `GetOne` (принимает code). В SDK метод называется `GetByCode`.
