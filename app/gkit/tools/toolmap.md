# SDK Methods Map

Анализ публичных методов amoCRM Go SDK для создания Genkit Tools.

---

## AccountService (`core/services/account.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `GetCurrent` | `ctx context.Context, with []string` | `*models.Account, error` | get |
| `AvailableWith` | - | `[]string` | list |

---

## CallsService (`core/services/calls.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Add` | `ctx context.Context, calls []models.Call` | `[]models.Call, error` | create |
| `AddOne` | `ctx context.Context, call *models.Call` | `*models.Call, error` | create |

---

## CatalogElementsService (`core/services/catalog_elements.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *CatalogElementsFilter` | `[]models.CatalogElement, error` | list |
| `GetOne` | `ctx context.Context, elementID int` | `*models.CatalogElement, error` | get |
| `Create` | `ctx context.Context, elements []models.CatalogElement` | `[]models.CatalogElement, error` | create |
| `Update` | `ctx context.Context, elements []models.CatalogElement` | `[]models.CatalogElement, error` | update |

---

## CatalogsService (`core/services/catalogs.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *CatalogsFilter` | `[]models.Catalog, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Catalog, error` | get |
| `Create` | `ctx context.Context, catalogs []models.Catalog` | `[]models.Catalog, error` | create |
| `Update` | `ctx context.Context, catalogs []models.Catalog` | `[]models.Catalog, error` | update |

---

## ChatTemplatesService (`core/services/chat_templates.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *ChatTemplatesFilter` | `[]ChatTemplate, error` | list |
| `Create` | `ctx context.Context, templates []ChatTemplate` | `[]ChatTemplate, error` | create |
| `Update` | `ctx context.Context, templates []ChatTemplate` | `[]ChatTemplate, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |
| `DeleteMany` | `ctx context.Context, ids []int` | `error` | delete |
| `GetOne` | `ctx context.Context, id int, with []string` | `*ChatTemplate, error` | get |
| `SendOnReview` | `ctx context.Context, templateID int` | `[]ChatTemplateReview, error` | update |
| `UpdateReviewStatus` | `ctx context.Context, templateID, reviewID int, status string` | `*ChatTemplateReview, error` | update |

---

## CompaniesService (`core/services/companies.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *CompaniesFilter` | `[]models.Company, error` | list |
| `GetOne` | `ctx context.Context, id int, with []string` | `*models.Company, error` | get |
| `Create` | `ctx context.Context, companies []models.Company` | `[]models.Company, error` | create |
| `Update` | `ctx context.Context, companies []models.Company` | `[]models.Company, error` | update |
| `Link` | `ctx context.Context, companyID int, entityType string, entityID int, metadata map[string]interface{}` | `error` | link |
| `Unlink` | `ctx context.Context, companyID int, entityType string, entityID int` | `error` | unlink |

---

## ContactsService (`core/services/contacts.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *ContactsFilter` | `[]models.Contact, error` | list |
| `GetOne` | `ctx context.Context, id int, with []string` | `*models.Contact, error` | get |
| `Create` | `ctx context.Context, contacts []models.Contact` | `[]models.Contact, error` | create |
| `Update` | `ctx context.Context, contacts []models.Contact` | `[]models.Contact, error` | update |
| `Link` | `ctx context.Context, contactID int, entityType string, entityID int, metadata map[string]interface{}` | `error` | link |
| `Unlink` | `ctx context.Context, contactID int, entityType string, entityID int` | `error` | unlink |
| `GetChats` | `ctx context.Context, contactID int` | `[]models.ChatLink, error` | list |
| `LinkChats` | `ctx context.Context, links []models.ChatLink` | `[]models.ChatLink, error` | link |

---

## CurrenciesService (`core/services/currencies.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *CurrenciesFilter` | `[]models.Currency, error` | list |

---

## CustomFieldGroupsService (`core/services/custom_field_groups.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *CustomFieldGroupsFilter` | `[]models.CustomFieldGroup, error` | list |
| `GetOne` | `ctx context.Context, id string` | `*models.CustomFieldGroup, error` | get |
| `Create` | `ctx context.Context, groups []models.CustomFieldGroup` | `[]models.CustomFieldGroup, error` | create |
| `Update` | `ctx context.Context, groups []models.CustomFieldGroup` | `[]models.CustomFieldGroup, error` | update |
| `Delete` | `ctx context.Context, id string` | `error` | delete |

---

## CustomFieldsService (`core/services/custom_fields.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, entityType string, filter *CustomFieldsFilter` | `[]models.CustomField, error` | list |
| `GetOne` | `ctx context.Context, entityType string, id int` | `*models.CustomField, error` | get |
| `Create` | `ctx context.Context, entityType string, fields []models.CustomField` | `[]models.CustomField, error` | create |
| `Update` | `ctx context.Context, entityType string, fields []models.CustomField` | `[]models.CustomField, error` | update |
| `Delete` | `ctx context.Context, entityType string, id int` | `error` | delete |

---

## CustomerBonusPointsService (`core/services/customer_bonus_points.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context` | `*BonusPoints, error` | get |
| `EarnPoints` | `ctx context.Context, points int` | `int, error` | update |
| `RedeemPoints` | `ctx context.Context, points int` | `int, error` | update |

---

## CustomerStatusesService (`core/services/customer_statuses.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, page, limit int` | `[]models.Status, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Status, error` | get |
| `Create` | `ctx context.Context, statuses []models.Status` | `[]models.Status, error` | create |
| `Update` | `ctx context.Context, statuses []models.Status` | `[]models.Status, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |

---

## CustomerTransactionsService (`core/services/customer_transactions.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *TransactionsFilter` | `[]Transaction, error` | list |
| `Create` | `ctx context.Context, transactions []Transaction, accrueBonus bool` | `[]Transaction, error` | create |
| `Delete` | `ctx context.Context, transactionID int` | `error` | delete |

---

## CustomersService (`core/services/customers.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *CustomersFilter` | `[]models.Customer, error` | list |
| `GetOne` | `ctx context.Context, id int, with []string` | `*models.Customer, error` | get |
| `Create` | `ctx context.Context, customers []models.Customer` | `[]models.Customer, error` | create |
| `Update` | `ctx context.Context, customers []models.Customer` | `[]models.Customer, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |
| `Link` | `ctx context.Context, customerID int, links []models.EntityLink` | `[]models.EntityLink, error` | link |

---

## EntityFilesService (`core/services/entity_files.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, page, limit int` | `[]models.FileLink, error` | list |
| `Link` | `ctx context.Context, fileUUIDs []string` | `[]models.FileLink, error` | link |
| `Unlink` | `ctx context.Context, fileUUID string` | `error` | unlink |

---

## EntitySubscriptionsService (`core/services/entity_subscriptions.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, page, limit int` | `[]models.Subscription, error` | list |
| `Subscribe` | `ctx context.Context, userIDs []int` | `[]models.Subscription, error` | link |
| `Unsubscribe` | `ctx context.Context, userID int` | `error` | unlink |

---

## EventTypesService (`core/services/event_types.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *EventTypesFilter` | `[]models.EventType, error` | list |

---

## EventsService (`core/services/events.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *EventsFilter` | `[]models.Event, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Event, error` | get |

---

## FilesService (`core/services/files.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *FilesFilter` | `[]models.File, error` | list |
| `GetOne` | `ctx context.Context, uuid string` | `*models.File, error` | get |
| `Delete` | `ctx context.Context, uuid string` | `error` | delete |
| `UploadOne` | `ctx context.Context, params FileUploadParams` | `*models.File, error` | create |

---

## LeadsService (`core/services/leads.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *LeadsFilter` | `[]models.Lead, error` | list |
| `GetWithPagination` | `ctx context.Context, filter *LeadsFilter` | `*models.LeadsResponse, error` | list |
| `GetOne` | `ctx context.Context, id int, with []string` | `*models.Lead, error` | get |
| `Create` | `ctx context.Context, leads []models.Lead` | `[]models.Lead, error` | create |
| `CreateOne` | `ctx context.Context, lead *models.Lead` | `*models.Lead, error` | create |
| `Update` | `ctx context.Context, leads []models.Lead` | `[]models.Lead, error` | update |
| `UpdateOne` | `ctx context.Context, lead *models.Lead` | `*models.Lead, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |
| `Link` | `ctx context.Context, leadID int, links []models.EntityLink` | `[]models.EntityLink, error` | link |
| `Unlink` | `ctx context.Context, leadID int, links []models.EntityLink` | `error` | unlink |
| `AddComplex` | `ctx context.Context, leads []models.Lead` | `[]ComplexLeadResult, error` | create |
| `AddOneComplex` | `ctx context.Context, lead *models.Lead` | `*ComplexLeadResult, error` | create |

---

## LinksService (`core/services/links.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, entityType string, entityID int, filter *LinksFilter` | `[]models.EntityLink, error` | list |
| `Link` | `ctx context.Context, entityType string, entityID int, links []models.EntityLink` | `[]models.EntityLink, error` | link |
| `Unlink` | `ctx context.Context, entityType string, entityID int, links []models.EntityLink` | `error` | unlink |

---

## LossReasonsService (`core/services/loss_reasons.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *LossReasonsFilter` | `[]models.LossReason, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.LossReason, error` | get |
| `Create` | `ctx context.Context, reasons []models.LossReason` | `[]models.LossReason, error` | create |
| `Update` | `ctx context.Context, reasons []models.LossReason` | `[]models.LossReason, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |

---

## NotesService (`core/services/notes.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, entityType string, filter *NotesFilter` | `[]models.Note, error` | list |
| `GetByParent` | `ctx context.Context, entityType string, entityID int, filter *NotesFilter` | `[]models.Note, error` | list |
| `GetOne` | `ctx context.Context, entityType string, id int` | `*models.Note, error` | get |
| `Create` | `ctx context.Context, entityType string, notes []models.Note` | `[]models.Note, error` | create |
| `Update` | `ctx context.Context, entityType string, notes []models.Note` | `[]models.Note, error` | update |

---

## PipelinesService (`core/services/pipelines.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context` | `[]models.Pipeline, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Pipeline, error` | get |
| `Create` | `ctx context.Context, pipelines []models.Pipeline` | `[]models.Pipeline, error` | create |
| `Update` | `ctx context.Context, pipelines []models.Pipeline` | `[]models.Pipeline, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |
| `GetStatuses` | `ctx context.Context, pipelineID int` | `[]models.Status, error` | list |
| `GetStatus` | `ctx context.Context, pipelineID int, statusID int` | `*models.Status, error` | get |
| `CreateStatus` | `ctx context.Context, pipelineID int, status *models.Status` | `*models.Status, error` | create |
| `UpdateStatus` | `ctx context.Context, pipelineID int, status *models.Status` | `*models.Status, error` | update |
| `DeleteStatus` | `ctx context.Context, pipelineID int, statusID int` | `error` | delete |

---

## ProductsService (`core/services/products.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *ProductsFilter` | `[]models.CatalogElement, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.CatalogElement, error` | get |
| `Create` | `ctx context.Context, products []models.CatalogElement` | `[]models.CatalogElement, error` | create |
| `Update` | `ctx context.Context, products []models.CatalogElement` | `[]models.CatalogElement, error` | update |
| `Delete` | `ctx context.Context, productIDs []int` | `error` | delete |

---

## RolesService (`core/services/roles.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *RoleFilter` | `[]models.Role, error` | list |
| `GetOne` | `ctx context.Context, id int, with []string` | `*models.Role, error` | get |
| `Create` | `ctx context.Context, roles []models.Role` | `[]models.Role, error` | create |
| `Update` | `ctx context.Context, roles []models.Role` | `[]models.Role, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |

---

## SegmentsService (`core/services/segments.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *SegmentsFilter` | `[]models.Segment, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Segment, error` | get |
| `Create` | `ctx context.Context, segments []models.Segment` | `[]models.Segment, error` | create |
| `Delete` | `ctx context.Context, id int` | `error` | delete |

---

## ShortLinksService (`core/services/short_links.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *ShortLinksFilter` | `[]models.ShortLink, error` | list |
| `Create` | `ctx context.Context, links []models.ShortLink` | `[]models.ShortLink, error` | create |
| `Delete` | `ctx context.Context, id int` | `error` | delete |

---

## SourcesService (`core/services/sources.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *SourcesFilter` | `[]models.Source, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Source, error` | get |
| `Create` | `ctx context.Context, sources []models.Source` | `[]models.Source, error` | create |
| `Update` | `ctx context.Context, sources []models.Source` | `[]models.Source, error` | update |
| `Delete` | `ctx context.Context, id int` | `error` | delete |

---

## TagsService (`core/services/tags.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, entityType string, filter *TagsFilter` | `[]models.Tag, error` | list |
| `Create` | `ctx context.Context, entityType string, tags []models.Tag` | `[]models.Tag, error` | create |
| `Delete` | `ctx context.Context, entityType string, tags []models.Tag` | `error` | delete |

---

## TalksService (`core/services/talks.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *TalksFilter` | `[]models.Talk, error` | list |
| `Close` | `ctx context.Context, talkID string` | `error` | update |

---

## TasksService (`core/services/tasks.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *TasksFilter` | `[]models.Task, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.Task, error` | get |
| `Create` | `ctx context.Context, tasks []models.Task` | `[]models.Task, error` | create |
| `Update` | `ctx context.Context, tasks []models.Task` | `[]models.Task, error` | update |
| `Complete` | `ctx context.Context, id int, resultText string` | `*models.Task, error` | update |

---

## UnsortedService (`core/services/unsorted.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *UnsortedFilter` | `[]models.Unsorted, error` | list |
| `GetOne` | `ctx context.Context, uid string` | `*models.Unsorted, error` | get |
| `Create` | `ctx context.Context, category string, unsorted []models.Unsorted` | `[]models.Unsorted, error` | create |
| `Accept` | `ctx context.Context, uid string, params map[string]interface{}` | `*models.UnsortedAcceptResult, error` | update |
| `Decline` | `ctx context.Context, uid string, params map[string]interface{}` | `*models.UnsortedDeclineResult, error` | delete |
| `Link` | `ctx context.Context, uid string, linkData map[string]interface{}` | `*models.UnsortedLinkResult, error` | link |
| `Summary` | `ctx context.Context, filter *UnsortedFilter` | `*models.UnsortedSummary, error` | get |

---

## UsersService (`core/services/users.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *UsersFilter` | `[]models.User, error` | list |
| `GetOne` | `ctx context.Context, id int` | `*models.User, error` | get |
| `Create` | `ctx context.Context, users []models.User` | `[]models.User, error` | create |
| `AddToGroup` | `ctx context.Context, userID int, groupID int` | `error` | update |
| `GetRoles` | `ctx context.Context, filter *RolesFilter` | `[]models.Role, error` | list |
| `GetRole` | `ctx context.Context, id int` | `*models.Role, error` | get |

---

## WebhooksService (`core/services/webhooks.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *WebhooksFilter` | `[]models.Webhook, error` | list |
| `Subscribe` | `ctx context.Context, webhook *models.Webhook` | `*models.Webhook, error` | create |
| `Unsubscribe` | `ctx context.Context, webhook *models.Webhook` | `error` | delete |

---

## WebsiteButtonsService (`core/services/website_buttons.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *WebsiteButtonsFilter, with []string` | `[]models.WebsiteButton, error` | list |
| `GetOne` | `ctx context.Context, sourceID int, with []string` | `*models.WebsiteButton, error` | get |
| `CreateAsync` | `ctx context.Context, request *models.WebsiteButtonCreateRequest` | `*models.WebsiteButtonCreateResponse, error` | create |
| `UpdateAsync` | `ctx context.Context, request *models.WebsiteButtonUpdateRequest` | `*models.WebsiteButton, error` | update |
| `AddOnlineChatAsync` | `ctx context.Context, sourceID int` | `error` | create |

---

## WidgetsService (`core/services/widgets.go`)

| Method | Params | Returns | Operation |
|--------|--------|---------|-----------|
| `Get` | `ctx context.Context, filter *WidgetsFilter` | `[]models.Widget, error` | list |
| `GetOne` | `ctx context.Context, code string` | `*models.Widget, error` | get |
| `Install` | `ctx context.Context, code string` | `*models.Widget, error` | create |
| `Uninstall` | `ctx context.Context, code string` | `error` | delete |

---


# Архитектура Genkit Tools (Гибридная группировка)

В этом документе описана структура инструментов (Tools) для Genkit, созданная на основе **комбинации Варианта 1 (Generics)** и **Варианта 3 (Core vs Helpers)**.

**Ключевая идея:**
1.  **Core Entities (Сделки, Контакты, Компании)** — получают **индивидуальные инструменты**, так как они критически важны, имеют сложные уникальные фильтры и структуры данных (Custom Fields).
2.  **Reference Data (Справочники)** — объединяются в **единый полиморфный инструмент** `crm_get_reference`, так как логика их получения идентична (список или ID).
3.  **Actions (Действия)** — группируются по типу действия (Связи, Задачи), а не по сущности.

---
