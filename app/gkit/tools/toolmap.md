# SDK Methods Map

Анализ публичных методов amoCRM Go SDK для создания Genkit Tools.

---

## AccountService (`core/adapters/account.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `GetCurrent` | `ctx context.Context, opts ...GetOneOption` | `*models.Account, error` | get |

---

## CallsService (`core/adapters/calls.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Create` | `ctx context.Context, calls []*models.Call` | `[]*models.Call, *PageMeta, error` | create |
| `CreateOne` | `ctx context.Context, call *models.Call` | `*models.Call, error` | create |

---

## CatalogElementsService (`core/adapters/catalog_elements.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *CatalogElementsFilter` | `[]*models.CatalogElement, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, elementID int` | `*models.CatalogElement, error` | get |
| `Create` | `ctx context.Context, elements []*models.CatalogElement` | `[]*models.CatalogElement, *PageMeta, error` | create |
| `Update` | `ctx context.Context, elements []*models.CatalogElement` | `[]*models.CatalogElement, *PageMeta, error` | update |
| `Link` | `ctx context.Context, elementID int, entityType string, entityID int, metadata map[string]interface{}` | `error` | link |
| `Unlink` | `ctx context.Context, elementID int, entityType string, entityID int` | `error` | unlink |

---

## CatalogsService (`core/adapters/catalogs.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *CatalogsFilter` | `[]*models.Catalog, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Catalog, error` | get |
| `Create` | `ctx context.Context, catalogs []*models.Catalog` | `[]*models.Catalog, *PageMeta, error` | create |
| `Update` | `ctx context.Context, catalogs []*models.Catalog` | `[]*models.Catalog, *PageMeta, error` | update |

---

## ChatTemplatesService (`core/adapters/chat_templates.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *ChatTemplatesFilter` | `[]*ChatTemplate, *PageMeta, error` | list |
| `Create` | `ctx context.Context, templates []*ChatTemplate` | `[]*ChatTemplate, *PageMeta, error` | create |
| `Update` | `ctx context.Context, templates []*ChatTemplate` | `[]*ChatTemplate, *PageMeta, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |
| `DeleteMany` | `ctx context.Context, ids []int` | `error` | delete |
| `GetOne` | `ctx context.Context, id int, opts ...GetOneOption` | `*ChatTemplate, error` | get |
| `SendOnReview` | `ctx context.Context, templateID int` | `[]ChatTemplateReview, error` | update |
| `UpdateReviewStatus` | `ctx context.Context, templateID, reviewID int, status string` | `*ChatTemplateReview, error` | update |

---

## CompaniesService (`core/adapters/companies.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *CompaniesFilter` | `[]*models.Company, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int, opts ...GetOneOption` | `*models.Company, error` | get |
| `Create` | `ctx context.Context, companies []*models.Company` | `[]*models.Company, *PageMeta, error` | create |
| `Update` | `ctx context.Context, companies []*models.Company` | `[]*models.Company, *PageMeta, error` | update |
| `SyncOne` | `ctx context.Context, company *models.Company, with []string` | `*models.Company, error` | create/update |
| `Link` | `ctx context.Context, companyID int, entityType string, entityID int, metadata map[string]interface{}` | `error` | link |
| `Unlink` | `ctx context.Context, companyID int, entityType string, entityID int` | `error` | unlink |
| `GetLinks` | `ctx context.Context, companyID int` | `[]models.EntityLink, error` | list |

---

## ContactsService (`core/adapters/contacts.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *ContactsFilter` | `[]*models.Contact, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int, opts ...GetOneOption` | `*models.Contact, error` | get |
| `Create` | `ctx context.Context, contacts []*models.Contact` | `[]*models.Contact, *PageMeta, error` | create |
| `Update` | `ctx context.Context, contacts []*models.Contact` | `[]*models.Contact, *PageMeta, error` | update |
| `SyncOne` | `ctx context.Context, contact *models.Contact, with []string` | `*models.Contact, error` | create/update |
| `Link` | `ctx context.Context, contactID int, entityType string, entityID int, metadata map[string]interface{}` | `error` | link |
| `Unlink` | `ctx context.Context, contactID int, entityType string, entityID int` | `error` | unlink |
| `GetLinks` | `ctx context.Context, contactID int` | `[]models.EntityLink, error` | list |
| `GetChats` | `ctx context.Context, contactID int` | `[]models.ChatLink, error` | list |
| `LinkChats` | `ctx context.Context, links []models.ChatLink` | `[]models.ChatLink, error` | link |

---

## CurrenciesService (`core/adapters/currencies.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *CurrenciesFilter` | `[]*models.Currency, *PageMeta, error` | list |

---

## CustomFieldGroupsService (`core/adapters/custom_field_groups.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *CustomFieldGroupsFilter` | `[]*models.CustomFieldGroup, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id string` | `*models.CustomFieldGroup, error` | get |
| `Create` | `ctx context.Context, groups []*models.CustomFieldGroup` | `[]*models.CustomFieldGroup, *PageMeta, error` | create |
| `Update` | `ctx context.Context, groups []*models.CustomFieldGroup` | `[]*models.CustomFieldGroup, *PageMeta, error` | update |
| `Delete` | `ctx context.Context, id string` | `error` | delete |

---

## CustomFieldsService (`core/adapters/custom_fields.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, entityType string, filter *CustomFieldsFilter` | `[]*models.CustomField, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, entityType string, id int` | `*models.CustomField, error` | get |
| `Create` | `ctx context.Context, entityType string, fields []*models.CustomField` | `[]*models.CustomField, *PageMeta, error` | create |
| `Update` | `ctx context.Context, entityType string, fields []*models.CustomField` | `[]*models.CustomField, *PageMeta, error` | update |
| `Delete` | `ctx context.Context, entityType string, id int` | `error` | delete |

---

## CustomerBonusPointsService (`core/adapters/customer_bonus_points.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context` | `*BonusPoints, error` | get |
| `EarnPoints` | `ctx context.Context, points int` | `int, error` | update |
| `RedeemPoints` | `ctx context.Context, points int` | `int, error` | update |

---

## CustomerStatusesService (`core/adapters/customer_statuses.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, page, limit int` | `[]*models.Status, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Status, error` | get |
| `Create` | `ctx context.Context, statuses []*models.Status` | `[]*models.Status, *PageMeta, error` | create |
| `Update` | `ctx context.Context, statuses []*models.Status` | `[]*models.Status, *PageMeta, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |

---

## CustomerTransactionsService (`core/adapters/customer_transactions.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *TransactionsFilter` | `[]*Transaction, *PageMeta, error` | list |
| `Create` | `ctx context.Context, transactions []*Transaction, accrueBonus bool` | `[]*Transaction, *PageMeta, error` | create |
| `Delete` | `ctx context.Context, transactionID int` | `error` | delete |

---

## CustomersService (`core/adapters/customers.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *CustomersFilter` | `[]*models.Customer, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int, opts ...GetOneOption` | `*models.Customer, error` | get |
| `Create` | `ctx context.Context, customers []*models.Customer` | `[]*models.Customer, *PageMeta, error` | create |
| `Update` | `ctx context.Context, customers []*models.Customer` | `[]*models.Customer, *PageMeta, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |
| `Link` | `ctx context.Context, customerID int, links []models.EntityLink` | `[]models.EntityLink, error` | link |

---

## EntityFilesService (`core/adapters/entity_files.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, page, limit int` | `[]*models.FileLink, *PageMeta, error` | list |
| `Link` | `ctx context.Context, fileUUIDs []string` | `[]models.FileLink, error` | link |
| `Unlink` | `ctx context.Context, fileUUID string` | `error` | unlink |

---

## EntitySubscriptionsService (`core/adapters/entity_subscriptions.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, page, limit int` | `[]*models.Subscription, *PageMeta, error` | list |
| `Subscribe` | `ctx context.Context, userIDs []int` | `[]models.Subscription, error` | link |
| `Unsubscribe` | `ctx context.Context, userID int` | `error` | unlink |

---

## EventTypesService (`core/adapters/event_types.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *EventTypesFilter` | `[]*models.EventType, *PageMeta, error` | list |

---

## EventsService (`core/adapters/events.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *EventsFilter` | `[]*models.Event, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Event, error` | get |

---

## FilesService (`core/adapters/files.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *FilesFilter` | `[]*models.File, *PageMeta, error` | list |
| `GetOneByUUID` | `ctx context.Context, uuid string` | `*models.File, error` | get |
| `Add` | `ctx context.Context, files []*models.File` | `[]*models.File, *PageMeta, error` | create |
| `Update` | `ctx context.Context, files []*models.File` | `[]*models.File, *PageMeta, error` | update |
| `UpdateOne` | `ctx context.Context, file *models.File` | `*models.File, error` | update |
| `Delete` | `ctx context.Context, collection *FilesCollection` | `error` | delete |
| `DeleteOne` | `ctx context.Context, uuid string` | `error` | delete |
| `UploadOne` | `ctx context.Context, params FileUploadParams` | `*models.File, error` | create |

---

## LeadsService (`core/adapters/leads.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *LeadsFilter` | `[]*models.Lead, *PageMeta, error` | list |
| `GetWithPagination` | `ctx context.Context, filter *LeadsFilter` | `*models.LeadsResponse, error` | list |
| `GetOne` | `ctx context.Context, id int, opts ...GetOneOption` | `*models.Lead, error` | get |
| `Create` | `ctx context.Context, leads []*models.Lead` | `[]*models.Lead, *PageMeta, error` | create |
| `CreateOne` | `ctx context.Context, lead *models.Lead` | `*models.Lead, error` | create |
| `Update` | `ctx context.Context, leads []*models.Lead` | `[]*models.Lead, *PageMeta, error` | update |
| `UpdateOne` | `ctx context.Context, lead *models.Lead` | `*models.Lead, error` | update |
| `SyncOne` | `ctx context.Context, lead *models.Lead, with []string` | `*models.Lead, error` | create/update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |
| `Link` | `ctx context.Context, leadID int, links []models.EntityLink` | `[]models.EntityLink, error` | link |
| `Unlink` | `ctx context.Context, leadID int, links []models.EntityLink` | `error` | unlink |
| `AddComplex` | `ctx context.Context, leads []*models.Lead` | `[]ComplexLeadResult, error` | create |
| `AddOneComplex` | `ctx context.Context, lead *models.Lead` | `*ComplexLeadResult, error` | create |
| `AddComplexInPlace` | `ctx context.Context, leads []*models.Lead` | `error` | create |

---

## LinksService (`core/adapters/links.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, entityType string, entityID int, filter *LinksFilter` | `[]*models.EntityLink, *PageMeta, error` | list |
| `Link` | `ctx context.Context, entityType string, entityID int, links []models.EntityLink` | `[]models.EntityLink, error` | link |
| `Unlink` | `ctx context.Context, entityType string, entityID int, links []models.EntityLink` | `error` | unlink |

---

## LossReasonsService (`core/adapters/loss_reasons.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *LossReasonsFilter` | `[]*models.LossReason, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.LossReason, error` | get |
| `Create` | `ctx context.Context, reasons []*models.LossReason` | `[]*models.LossReason, *PageMeta, error` | create |
| `Update` | `ctx context.Context, reasons []*models.LossReason` | `[]*models.LossReason, *PageMeta, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |

---

## NotesService (`core/adapters/notes.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, entityType string, filter *NotesFilter` | `[]*models.Note, *PageMeta, error` | list |
| `GetByParent` | `ctx context.Context, entityType string, entityID int, filter *NotesFilter` | `[]*models.Note, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, entityType string, id int` | `*models.Note, error` | get |
| `Create` | `ctx context.Context, entityType string, notes []*models.Note` | `[]*models.Note, *PageMeta, error` | create |
| `Update` | `ctx context.Context, entityType string, notes []*models.Note` | `[]*models.Note, *PageMeta, error` | update |

---

## PipelinesService (`core/adapters/pipelines.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context` | `[]*models.Pipeline, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Pipeline, error` | get |
| `Create` | `ctx context.Context, pipelines []*models.Pipeline` | `[]*models.Pipeline, *PageMeta, error` | create |
| `Update` | `ctx context.Context, pipelines []*models.Pipeline` | `[]*models.Pipeline, *PageMeta, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |

**Примечание:** Для работы со статусами воронки используйте `sdk.Statuses(pipelineID)`.

---

## StatusesService (`core/adapters/statuses.go`)

Сервис для управления статусами воронок. Доступен через `sdk.Statuses(pipelineID)`.

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, params url.Values` | `[]*models.Status, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int, params url.Values` | `*models.Status, error` | get |
| `Create` | `ctx context.Context, statuses []*models.Status` | `[]*models.Status, *PageMeta, error` | create |
| `CreateOne` | `ctx context.Context, status *models.Status` | `*models.Status, error` | create |
| `UpdateOne` | `ctx context.Context, status *models.Status` | `*models.Status, error` | update |
| `DeleteOne` | `ctx context.Context, id int` | `error` | delete |

**Примечания:**
- `Update()` (batch) — возвращает `ErrNotAvailableForAction`
- `Delete()` (batch) — возвращает `ErrNotAvailableForAction`

---

## ProductsService (`core/adapters/products.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *ProductsFilter` | `[]*models.CatalogElement, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.CatalogElement, error` | get |
| `Create` | `ctx context.Context, products []*models.CatalogElement` | `[]*models.CatalogElement, *PageMeta, error` | create |
| `Update` | `ctx context.Context, products []*models.CatalogElement` | `[]*models.CatalogElement, *PageMeta, error` | update |
| `Delete` | `ctx context.Context, productIDs []int` | `error` | delete |

**⚠️ Примечания:**
- `Get`, `Create`, `Update` — возвращают `ErrNotAvailableForAction`
- API не поддерживает прямую работу с товарами через ProductsService
- Используйте `CatalogsService` с catalog_id товарного каталога

---

## RolesService (`core/adapters/roles.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *RoleFilter` | `[]*models.Role, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int, opts ...GetOneOption` | `*models.Role, error` | get |
| `Create` | `ctx context.Context, roles []*models.Role` | `[]*models.Role, *PageMeta, error` | create |
| `Update` | `ctx context.Context, roles []*models.Role` | `[]*models.Role, *PageMeta, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |

---

## SegmentsService (`core/adapters/segments.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *SegmentsFilter` | `[]*models.Segment, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Segment, error` | get |
| `Create` | `ctx context.Context, segments []*models.Segment` | `[]*models.Segment, *PageMeta, error` | create |
| `Delete` | `ctx context.Context, id int` | `error` | delete |

---

## ShortLinksService (`core/adapters/short_links.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *ShortLinksFilter` | `[]*models.ShortLink, *PageMeta, error` | list |
| `Create` | `ctx context.Context, links []*models.ShortLink` | `[]*models.ShortLink, *PageMeta, error` | create |
| `Delete` | `ctx context.Context, id int` | `error` | delete |

---

## SourcesService (`core/adapters/sources.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *SourcesFilter` | `[]*models.Source, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Source, error` | get |
| `Create` | `ctx context.Context, sources []*models.Source` | `[]*models.Source, *PageMeta, error` | create |
| `Update` | `ctx context.Context, sources []*models.Source` | `[]*models.Source, *PageMeta, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |

---

## TagsService (`core/adapters/tags.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, entityType string, filter *TagsFilter` | `[]*models.Tag, *PageMeta, error` | list |
| `Create` | `ctx context.Context, entityType string, tags []*models.Tag` | `[]*models.Tag, *PageMeta, error` | create |
| `Delete` | `ctx context.Context, entityType string, tags []models.Tag` | `error` | delete |

---

## TalksService (`core/adapters/talks.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *TalksFilter` | `[]*models.Talk, *PageMeta, error` | list |
| `Close` | `ctx context.Context, talkID string` | `error` | update |

**⚠️ Примечание:** `Get` возвращает `ErrNotAvailableForAction` — API не поддерживает получение списка бесед

---

## TasksService (`core/adapters/tasks.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *TasksFilter` | `[]*models.Task, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Task, error` | get |
| `Create` | `ctx context.Context, tasks []*models.Task` | `[]*models.Task, *PageMeta, error` | create |
| `Update` | `ctx context.Context, tasks []*models.Task` | `[]*models.Task, *PageMeta, error` | update |
| `Complete` | `ctx context.Context, id int, resultText string` | `*models.Task, error` | update |

---

## UnsortedService (`core/adapters/unsorted.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *UnsortedFilter` | `[]*models.Unsorted, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, uid string` | `*models.Unsorted, error` | get |
| `Create` | `ctx context.Context, category string, unsorted []*models.Unsorted` | `[]*models.Unsorted, *PageMeta, error` | create |
| `Accept` | `ctx context.Context, uid string, params map[string]interface{}` | `*models.UnsortedAcceptResult, error` | update |
| `Decline` | `ctx context.Context, uid string, params map[string]interface{}` | `*models.UnsortedDeclineResult, error` | delete |
| `Link` | `ctx context.Context, uid string, linkData map[string]interface{}` | `*models.UnsortedLinkResult, error` | link |
| `Summary` | `ctx context.Context, filter *UnsortedFilter` | `*models.UnsortedSummary, error` | get |

---

## UsersService (`core/adapters/users.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *UsersFilter` | `[]*models.User, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.User, error` | get |
| `Create` | `ctx context.Context, users []*models.User` | `[]*models.User, *PageMeta, error` | create |

---

## WebhooksService (`core/adapters/webhooks.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *WebhooksFilter` | `[]*models.Webhook, *PageMeta, error` | list |
| `Subscribe` | `ctx context.Context, webhook *models.Webhook` | `*models.Webhook, error` | create |
| `Unsubscribe` | `ctx context.Context, webhook *models.Webhook` | `error` | delete |

---

## WebsiteButtonsService (`core/adapters/website_buttons.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *WebsiteButtonsFilter` | `[]*models.WebsiteButton, *PageMeta, error` | list |
| `GetOne` | `ctx context.Context, sourceID int, opts ...GetOneOption` | `*models.WebsiteButton, error` | get |
| `CreateAsync` | `ctx context.Context, request *models.WebsiteButtonCreateRequest` | `*models.WebsiteButtonCreateResponse, error` | create |
| `UpdateAsync` | `ctx context.Context, request *models.WebsiteButtonUpdateRequest` | `*models.WebsiteButton, error` | update |
| `AddOnlineChatAsync` | `ctx context.Context, sourceID int` | `error` | create |

---

## WidgetsService (`core/adapters/widgets.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *WidgetsFilter` | `[]*models.Widget, *PageMeta, error` | list |
| `GetByCode` | `ctx context.Context, code string` | `*models.Widget, error` | get |
| `Add` | `ctx context.Context, widgets []*models.Widget` | `[]*models.Widget, *PageMeta, error` | create |
| `Update` | `ctx context.Context, widgets []*models.Widget` | `[]*models.Widget, *PageMeta, error` | update |
| `Install` | `ctx context.Context, widget *models.Widget` | `*models.Widget, error` | create |
| `InstallByCode` | `ctx context.Context, code string` | `*models.Widget, error` | create |
| `Uninstall` | `ctx context.Context, code string` | `error` | delete |

---


# Архитектура Genkit Tools (Гибридная группировка)

В этом документе описана структура инструментов (Tools) для Genkit, созданная на основе **комбинации Варианта 1 (Generics)** и **Варианта 3 (Core vs Helpers)**.

**Ключевая идея:**
1.  **Core Entities (Сделки, Контакты, Компании)** — получают **индивидуальные инструменты**, так как они критически важны, имеют сложные уникальные фильтры и структуры данных (Custom Fields).
2.  **Reference Data (Справочники)** — объединяются в **единый полиморфный инструмент** `crm_get_reference`, так как логика их получения идентична (список или ID).
3.  **Actions (Действия)** — группируются по типу действия (Связи, Задачи), а не по сущности.

---
