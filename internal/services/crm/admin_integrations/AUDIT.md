# Аудит: admin_integrations

## Вход (LLM → tool → сервис)

### Поля, которые сейчас числовые ID, но должны быть именами

**`ID int` в website_buttons / short_links / chat_templates**
LLM не знает ID без предварительного `list`. Для `website_buttons` поле `id` = `source_id` — неочевидно из названия.

**`data["review_id"]` кастится как `float64`**
```go
reviewID, _ := input.Data["review_id"].(float64)
```
Числовой ID ревью в свободном `map[string]any`. Хрупко и неочевидно для LLM.

**`PipelineID *int` в WebsiteButtonCreateRequest / UpdateRequest**
LLM должна передать числовой ID воронки в `data`. Поле не задокументировано в `AdminIntegrationsInput`. Воронки инжектируются в контекст, но маппинг "имя → ID" нигде в схеме инструмента не закреплён.

### Поля, которые отсутствуют, но нужны LLM

- **`destination` и `settings` для webhooks** — скрытый контракт через `data`. Явных полей `Destination string` и `WebhookSettings []string` нет в `AdminIntegrationsInput`.
- **`website_buttons` create/update: структура `data` не задокументирована** — поля `name`, `pipeline_id`, `trusted_websites`, `is_duplication_control_enabled` нигде не описаны в JSON Schema.
- **`short_links` create: `entity_id`/`entity_type`** — SDK поддерживает привязку к сущности через metadata, но `CreateShortLink(url)` создаёт `ShortLink{URL: url}`, теряя entity-привязку.
- **`chat_templates`: нет фильтра по `type` (amocrm/waba)** — SDK-модель имеет `Type ChatTemplateType`, фильтра нет.
- **`chat_templates`: нет create/update** — `Service` интерфейс не предоставляет этих методов несмотря на `ToAPI()` в SDK-модели.
- **Перечисление webhook-событий** — 20+ констант (`add_lead`, `update_contact` и др.) нигде не описаны в schema. LLM должна угадывать.

### Неудобные структуры/типы

- **`Settings map[string]any` для webhook** — `settings` это `[]string` (список событий), но передаётся через `data["settings"].([]any)` с ручным кастингом. Если LLM передаст строку вместо массива — молча вернётся пустой список.
- **`IDs []int` только для `chat_templates: delete_many`** — поле видно в schema для всех layers, шум.
- **`with` внутри `filter`** — `With []string` — параметр обогащения ответа, не фильтрации. Для `get` LLM вынуждена создавать `filter` только ради `with`.
- **`Widget.Version interface{}`** — строка или число. LLM не может предсказать тип.
- **`ChatTemplateReview` без `ID`** — при `update_review` нужен `review_id`, но `send_review` возвращает `[]ChatTemplateReview` без ID. Цепочка `send_review → update_review` неработоспособна.

---

## Выход (сервис → tool → LLM)

### Что сейчас возвращается

| Layer / Action | Тип возврата |
|----------------|-------------|
| `webhooks / list` | `[]models.Webhook` |
| `webhooks / subscribe` | `*models.Webhook` |
| `widgets / list` | `[]*models.Widget` |
| `widgets / get, install` | `*models.Widget` |
| `website_buttons / list, get` | `[]*models.WebsiteButton` / `*models.WebsiteButton` |
| `website_buttons / create` | `*models.WebsiteButtonCreateResponse` (только source_id, creation_status) |
| `website_buttons / update` | `*models.WebsiteButton` |
| `chat_templates / list` | `[]*models.ChatTemplate` |
| `chat_templates / send_review` | `[]models.ChatTemplateReview` |
| `chat_templates / update_review` | `*models.ChatTemplateReview` |
| `short_links / list, create` | `[]models.ShortLink` / `models.ShortLink` |

### Какие SDK With-параметры не используются

**`website_buttons`: `WebsiteButtonWithDeleted = "deleted"`**
Константа объявлена в SDK, но не упомянута в `jsonschema_description`. LLM не знает что можно получить удалённые кнопки.

**`widgets`** — `WidgetsService.Get` не принимает `with`. `widgets_template` с описанием настроек возвращается, но нигде не документируется.

**`webhooks`** — `WebhooksService.Get` не принимает `with`. Нормально для API, но не задокументировано.

### Числовые ID в ответе, которые LLM не может интерпретировать

- `Webhook.AccountID int` / `Webhook.CreatedBy int` — без имени пользователя
- `WebsiteButton.AccountID int`, `WebsiteButton.PipelineID *int` — без имени воронки
- `WebsiteButton.SourceID int` / `ButtonID int` — два разных ID одного объекта, различие нигде не объяснено
- `Widget.PipelineID int` — без имени воронки
- `ShortLink.EntityID int` — числовой ID сущности
- `ChatTemplate.AccountID int` — технический шум
- `ChatTemplate.CreatedAt int64` / `ChatTemplateReview.CreatedAt int64` — Unix timestamps

### Что теряется по сравнению с тем, что SDK может вернуть

1. **`short_links` create** — `entity_id`/`entity_type` теряются: `CreateShortLink(url)` создаёт `ShortLink{URL: url}` без metadata
2. **`website_buttons` update — баг**: `SourceID` помечен `json:"-"`, PATCH всегда уходит на `/website_buttons/0`
3. **`chat_templates`** — нет create/update несмотря на полный `ToAPI()` в SDK
4. **`PageMeta`** для webhooks отбрасывается
5. **`WebsiteButtonWithDeleted`** — полностью недокументирован

---

## Итого

**Приоритет рефакторинга: высокий**

Два функциональных бага, несколько скрытых контрактов, систематически потерянные SDK-возможности для WABA и website_buttons сценариев.

### Список конкретных изменений

**Баги:**
1. **[BUG] Исправить `UpdateWebsiteButton`** — `SourceID` не заполняется из-за `json:"-"`. Передавать из `input.ID` напрямую: `req.SourceID = input.ID`
2. **[BUG] `ChatTemplateReview` без `ID`** — добавить `ID int` в модель либо документировать источник `review_id`. Цепочка `send_review → update_review` сейчас неработоспособна.

**Вход:**
3. Вынести `Destination` и `EventTypes []string` на верхний уровень `AdminIntegrationsInput` для webhooks вместо `data`
4. Вынести `With []string` из `filter` на верхний уровень input
5. Задокументировать `WebsiteButtonWithDeleted = "deleted"` в description поля `with`
6. Добавить `EntityID int` и `EntityType string` в input для `short_links: create`
7. Добавить перечисление webhook-событий в `jsonschema_description`
8. Добавить `TemplateType string` в `IntegrationsFilter` для фильтрации шаблонов по типу
9. Заменить `data map[string]any` для `website_buttons` на типизированный `WebsiteButtonData`
10. Убрать `IDs []int` из общего input, заменить на `data["ids"]` только для `chat_templates: delete_many`

**Выход:**
11. Добавить `CreateChatTemplate` и `UpdateChatTemplate` в `Service` интерфейс и реализацию
12. Форматировать Unix timestamps в RFC3339 перед возвратом LLM
13. Убрать `AccountID` из ответных структур на уровне tool-wrapper
