# Аудит сервисов Admin Integrations

Этот файл содержит результаты последовательного аудита каждого сервиса в папке `adapters/admin_integrations/` на соответствие `tools_schema.md` и возможностям SDK.

---

## webhooks.go
**Layer:** webhooks
**Schema actions:** list, subscribe, unsubscribe
**SDK service:** WebhooksService (`core/adapters/webhooks.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListWebhooks` | Поддерживает фильтрацию по `destination` |
| Subscribe | ✅ | `SubscribeWebhook` | |
| Unsubscribe| ✅ | `UnsubscribeWebhook`| |

**Genkit Tool Handler:**
- 🛠 `handleWebhooks` поддерживает `list`, `subscribe`, `unsubscribe`. 
- 🛠 Поддерживает фильтр `destination`.

**Статус:** ✅ Полностью соответствует
**TODO:** Нет

### Capabilities Coverage

**Filters:**
| SDK Filter Method | Bot Field | Status | Comment |
|-------------------|-----------|--------|---------|
| SetDestination | ✅ | ✅ | Можно найти вебхук по конкретному URL. |

**Parameters:**
- ✅ SDK: `models.Webhook` имеет поле `Sort`.
- ℹ️ API вебхуков работает строго по одному URL.

**Batch Operations:**
- ℹ️ API вебхуков работает строго по одному URL.

**With/Relations:**
- ℹ️ Не поддерживается для вебхуков.

---

## widgets.go
**Layer:** widgets
**Schema actions:** search, get, install, uninstall
**SDK service:** WidgetsService (`core/adapters/widgets.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListWidgets` | Поддерживает пагинацию |
| GetByCode | ✅ | `GetWidget` | |
| Install | ✅ | `InstallWidget` | Поддерживает `settings` |
| Uninstall | ✅ | `UninstallWidget`| |
| Add | ℹ️ | - | Удалено (не поддерживается API/SDK) |
| Update | ℹ️ | - | Удалено (не поддерживается API/SDK) |

**Genkit Tool Handler:**
- 🛠 `handleWidgets` поддерживает `search`, `get`, `install`, `uninstall`.
- ✅ Экшен `install` поддерживает передачу `settings`.

**Статус:** ✅ Полностью соответствует
**TODO:** Нет

### Capabilities Coverage

**Filters:**
| SDK Filter Method | Bot Field | Status | Comment |
|-------------------|-----------|--------|---------|
| SetLimit | ✅ | ✅ | |
| SetPage | ✅ | ✅ | |

**Parameters:**
- ✅ SDK: `Install` принимает объект `models.Widget` с мапой `Settings`.
- ✅ Bot: `InstallWidget` принимает `code` и `settings`.

**Batch Operations:**
- ℹ️ Не поддерживаются для этого типа сущности.

**With/Relations:**
- ℹ️ Не используются для виджетов.

---

## website_buttons.go
**Layer:** website_buttons
**Schema actions:** search, get, create, update, add_chat
**SDK service:** WebsiteButtonsService (`core/adapters/website_buttons.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListWebsiteButtons` | Поддерживает пагинацию и `with` |
| GetOne | ✅ | `GetWebsiteButton` | Поддерживает `with` |
| CreateAsync | ✅ | `CreateWebsiteButton` | |
| UpdateAsync | ✅ | `UpdateWebsiteButton` | |
| AddOnlineChat | ✅ | `AddOnlineChat` | Реализовано |

**Genkit Tool Handler:**
- 🛠 `handleWebsiteButtons` поддерживает `list`, `get`, `create`, `update`, `add_chat`.
- ✅ Поддерживает `with=scripts` для получения кода кнопки.

**Статус:** ✅ Полностью соответствует
**TODO:** Нет

### Capabilities Coverage

**Filters:**
| SDK Filter Method | Bot Field | Status | Comment |
|-------------------|-----------|--------|---------|
| Page | ✅ | ✅ | |
| Limit | ✅ | ✅ | |

**Parameters:**
- ✅ SDK: `Get` и `GetOne` поддерживают `with[scripts, deleted]`.
- ✅ Bot: Запрашивает `scripts` по запросу AI.

**Batch Operations:**
- ℹ️ Не поддерживаются.

**With/Relations:**
- ✅ Реализована поддержка `with` параметров.

---

## chat_templates.go
**Layer:** chat_templates
**Schema actions:** search, list, delete, delete_many, send_review, update_review
**SDK service:** ChatTemplatesService (`core/adapters/chat_templates.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListChatTemplates` | Поддерживает пагинацию и фильтр по `external_id` |
| Delete | ✅ | `DeleteChatTemplate` | |
| DeleteMany | ✅ | `DeleteChatTemplates` | Батч-удаление |
| SendOnReview | ✅ | `SendChatTemplateOnReview` | |
| UpdateReviewStatus | ✅ | `UpdateChatTemplateReviewStatus` | |

**Genkit Tool Handler:**
- 🛠 `handleChatTemplates` поддерживает все экшены, включая `delete_many`.
- ✅ Поддерживает фильтрацию по внешним ID.

**Статус:** ✅ Полностью соответствует
**TODO:** Нет

### Capabilities Coverage

**Filters:**
| SDK Filter Method | Bot Field | Status | Comment |
|-------------------|-----------|--------|---------|
| SetExternalIDs | ✅ | ✅ | Поддерживается фильтрация. |
| Page | ✅ | ✅ | |
| Limit | ✅ | ✅ | |

**Batch Operations:**
- ✅ `DeleteMany` позволяет удалять массив шаблонов.

---

## short_links.go
**Layer:** short_links
**Schema actions:** search, list, create, delete
**SDK service:** ShortLinksService (`core/adapters/short_links.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListShortLinks` | Поддерживает пагинацию |
| Create | ✅ | `CreateShortLinks` | Батч-создание |
| CreateOne | ✅ | `CreateShortLink` | Одиночное создание |
| Delete | ✅ | `DeleteShortLink` | |

**Genkit Tool Handler:**
- 🛠 `handleShortLinks` поддерживает массив `urls` при создании.

**Статус:** ✅ Полностью соответствует
**TODO:** Нет

### Capabilities Coverage

**Filters:**
| SDK Filter Method | Bot Field | Status | Comment |
|-------------------|-----------|--------|---------|
| Page | ✅ | ✅ | |
| Limit | ✅ | ✅ | |

**Batch Operations:**
- ✅ Массовое создание коротких ссылок поддержано.
