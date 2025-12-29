# Tools Schema

> Группировка по смыслу (техническому), не по бизнес-контексту.
> Бизнес-логика реализуется на уровне Flows.

---

## Implementation Pattern

**Unified Input + Switch по Action**

Каждый tool использует единую input-структуру с опциональными полями (`omitempty`). 
Genkit автоматически генерирует JSON Schema из Go struct. LLM передаёт только нужные поля.

```go
type ToolInput struct {
    Action     string         `json:"action"`              // обязательное
    EntityType string         `json:"entity_type,omitempty"`
    ID         int            `json:"id,omitempty"`
    Filter     map[string]any `json:"filter,omitempty"`
    Data       map[string]any `json:"data,omitempty"`
}

// Handler:
switch input.Action {
case "get":
    if input.ID == 0 {
        return nil, fmt.Errorf("id required for 'get'")
    }
    // ...
}
```

> **Важно:** Genkit не делает conditional validation. Проверка обязательности полей в зависимости от action — в коде handler'а.

---

## Reference Context

> Справочники НЕ являются отдельными Tools.
> Данные из Reference автоматически вшиваются в контекст Tools.
> LLM получает сведённые данные — не нужно вручную запрашивать справочники.

| Service                  | Данные                           | Вшивается в   |
| ------------------------ | -------------------------------------- | ----------------------- |
| AccountService           | Аккаунт, домен, права | Session context         |
| UsersService             | Пользователи               | entities, activities    |
| PipelinesService         | Воронки + статусы        | entities                |
| RolesService             | Роли                               | Session context         |
| CustomFieldsService      | Кастомные поля            | entities                |
| CustomFieldGroupsService | Группы полей                | entities                |
| CurrenciesService        | Валюты                           | Session context         |
| LossReasonsService       | Причины отказа            | entities (при close) |
| EventTypesService        | Типы событий                | activities              |
| SourcesService           | Источники                     | entities                |

---

## Tool 1: `entities`

**Что:** Основные бизнес-объекты (standalone)

**Services:** Leads, Contacts, Companies

**Actions:** search, get, create, update, delete, link, unlink, get_chats, link_chats

```json
{
  "entity_type": "leads | contacts | companies",
  "action": "...",
  "id": 123,
  "filter": {},
  "data": {},
  "link_to": { "type": "...", "id": 456 }
}
```

---

## Tool 2: `activities`

**Что:** Данные, привязанные к сущностям (attached)

**Services:** Tasks, Notes, Calls, Events, EntityFiles, Links, Tags, EntitySubscriptions, Talks

**Actions:** list, get, create, update, delete, complete, link, unlink, subscribe, unsubscribe, close

> **Layer-specific actions:**
> - `tasks`: list, get, create, update, complete
> - `notes`: list, get, create, update
> - `calls`: create (write-only, нет чтения через SDK)
> - `events`: list, get (read-only)
> - `files`: list, link, unlink
> - `links`: list, link, unlink
> - `tags`: list, create, link, unlink (delete не поддерживается API)
> - `subscriptions`: list, subscribe, unsubscribe
> - `talks`: list, close

```json
{
  "parent": { "type": "leads | contacts | companies", "id": 123 },
  "layer": "tasks | notes | calls | events | files | links | tags | subscriptions | talks",
  "action": "...",
  "id": 789,
  "data": {},
  "user_ids": [456, 789],
  "user_id": 456,
  "talk_id": "abc123"
}
```

---

## Tool 3: `complex_create`

**Что:** Создание сделки с контактами/компанией одним запросом

**Services:** Leads (AddComplex, AddOneComplex)

**Actions:** create

```json
{
  "lead": {
    "name": "...",
    "price": 100000,
    "_embedded": {
      "contacts": [...],
      "companies": [...]
    }
  }
}
```

---

## Tool 4: `products`

**Что:** Работа с товарами (элементами каталога) + привязка к сущностям

**Services:** ProductsService + Lead/Contact/Company (catalog_elements)

**Actions:** search, get, create, update, delete, get_by_entity, link, unlink, update_quantity

```json
{
  "action": "...",
  "product_id": 123,
  "filter": { "name": "iPhone" },
  "entity": { "type": "leads | contacts | companies", "id": 456 },
  "product": { "id": 123, "quantity": 2, "price_id": 789 },
  "data": { "name": "...", "custom_fields_values": [] }
}
```

---

## Tool 5: `catalogs`

**Что:** Конфигурация справочников (каталогов) — admin

**Services:** CatalogsService, CatalogElementsService

**Actions:** list, get, create, update, list_elements, get_element, create_element, update_element

```json
{
  "action": "...",
  "catalog_id": 1,
  "element_id": 123,
  "filter": {},
  "data": {}
}
```

---

## Tool 6: `files`

**Что:** Файловое хранилище аккаунта (storage layer)

**Services:** FilesService

**Actions:** search, get, upload, delete

```json
{
  "action": "search | get | upload | delete",
  "uuid": "file-uuid-here",
  "filter": { "page": 1, "limit": 50 },
  "upload_params": {
    "file_name": "document.pdf",
    "file_content": "base64...",
    "file_url": "https://..."
  }
}
```

---

## Tool 7: `unsorted`

**Что:** Обработка входящих заявок (Неразобранное)

**Services:** UnsortedService

**Actions:** search, get, accept, decline, link, summary

```json
{
  "action": "search | get | accept | decline | link | summary",
  "uid": "unsorted-uid",
  "filter": {
    "category": ["sip", "mail", "forms", "chats"],
    "pipeline_id": 123
  },
  "accept_params": { "user_id": 456, "status_id": 789 },
  "link_data": { "lead_id": 123 }
}
```

---

## Tool 8: `admin_schema`

**Что:** Структура данных (кастомные поля, справочники)

**Services:** CustomFieldsService, CustomFieldGroupsService, LossReasonsService, SourcesService

**Actions:** search, get, create, update, delete

```json
{
  "layer": "custom_fields | field_groups | loss_reasons | sources",
  "action": "search | get | create | update | delete",
  "entity_type": "leads | contacts | companies",
  "id": 123,
  "data": {}
}
```

---

## Tool 9: `admin_pipelines`

**Что:** Воронки и этапы

**Services:** PipelinesService

**Actions:** search, get, create, update, delete, get_statuses, get_status, create_status, update_status, delete_status

```json
{
  "action": "search | get | create | update | delete | get_statuses | create_status",
  "pipeline_id": 123,
  "status_id": 456,
  "data": {}
}
```

---

## Tool 10: `admin_users`

**Что:** Пользователи и права

**Services:** UsersService, RolesService

**Actions:** search, get, create, update, delete, add_to_group

```json
{
  "layer": "users | roles",
  "action": "search | get | create | update | delete | add_to_group",
  "id": 123,
  "user_id": 456,
  "group_id": 789,
  "data": {}
}
```

---

## Tool 11: `admin_integrations`

**Что:** Интеграции

**Services:** WebhooksService, WidgetsService, WebsiteButtonsService, ChatTemplatesService, ShortLinksService

**Actions:** search, get, create, update, delete, subscribe, unsubscribe, install, uninstall

```json
{
  "layer": "webhooks | widgets | website_buttons | chat_templates | short_links",
  "action": "search | get | create | update | delete | subscribe | install",
  "id": 123,
  "code": "widget_code",
  "data": {}
}
```

---

## Tool 12: `customers`

**Что:** Retention (Покупатели)

**Services:** CustomersService, CustomerBonusPointsService, CustomerStatusesService, CustomerTransactionsService, SegmentsService

**Actions:** search, get, create, update, delete, link, earn_points, redeem_points

> **Layer-specific limitations:**
> - `segments`: только search, get, create, delete (update не поддерживается API)

```json
{
  "layer": "customers | bonus_points | statuses | transactions | segments",
  "action": "search | get | create | update | delete | link | earn_points | redeem_points",
  "id": 123,
  "customer_id": 456,
  "points": 100,
  "data": {}
}
```

---

## Итого: 12 инструментов

### По режимам работы:

| Mode | Tool | Сервисов | Описание |
|------|------|----------|----------|
| **Work** | entities | 3 | Сделки, контакты, компании |
| **Work** | activities | 9 | Задачи, примечания, звонки, события, файлы, связи, теги, подписки, чаты |
| **Work** | complex_create | 1 | Создание сделки с контактами |
| **Work** | products | 1+ | Товары и привязка к сущностям |
| **Work** | catalogs | 2 | Справочники (каталоги) |
| **Work** | files | 1 | Файловое хранилище |
| **Work** | unsorted | 1 | Неразобранное |
| **Admin** | admin_schema | 4 | Кастомные поля, причины отказа, источники |
| **Admin** | admin_pipelines | 1 | Воронки и этапы |
| **Admin** | admin_users | 2 | Пользователи и роли |
| **Admin** | admin_integrations | 5 | Вебхуки, виджеты, шаблоны |
| **Retention** | customers | 5 | Покупатели, бонусы, транзакции |

### Reference Context (auto-inject):

AccountService, CurrenciesService, EventTypesService — read-only, вшиваются автоматически.

### Editable Reference (через admin tools):

CustomFieldsService, PipelinesService, UsersService, RolesService, SourcesService, LossReasonsService — доступны через `admin_*` инструменты.





