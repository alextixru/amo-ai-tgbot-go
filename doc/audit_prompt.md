# Промпт: Аудит toolmap.md vs SDK

## Задача

Последовательно проверить каждый раздел `toolmap.md` на соответствие реальным методам в `amocrm-sdk-go/core/adapters/`.

## Инструкция

Ты — аудитор SDK. Работай **строго последовательно**:

1. **Один раздел за раз** — берёшь следующий сервис из списка ниже
2. **Читаешь toolmap.md** — находишь соответствующий раздел
3. **Читаешь реальный .go файл** — смотришь outline и методы
4. **Записываешь различия** в файл `toolmap_audit.md`
5. **Спрашиваешь** — "Раздел X проверен. Продолжить?"

## Порядок проверки (32 сервиса)

```
1.  AccountService          → account.go
2.  CallsService            → calls.go
3.  CatalogElementsService  → catalog_elements.go
4.  CatalogsService         → catalogs.go
5.  ChatTemplatesService    → chat_templates.go
6.  CompaniesService        → companies.go
7.  ContactsService         → contacts.go
8.  CurrenciesService       → currencies.go
9.  CustomFieldGroupsService → custom_field_groups.go
10. CustomFieldsService     → custom_fields.go
11. CustomerBonusPointsService → customer_bonus_points.go
12. CustomerStatusesService → customer_statuses.go
13. CustomerTransactionsService → customer_transactions.go
14. CustomersService        → customers.go
15. EntityFilesService      → entity_files.go
16. EntitySubscriptionsService → entity_subscriptions.go
17. EventTypesService       → event_types.go
18. EventsService           → events.go
19. FilesService            → files.go
20. LeadsService            → leads.go
21. LinksService            → links.go
22. LossReasonsService      → loss_reasons.go
23. NotesService            → notes.go
24. PipelinesService        → pipelines.go
25. ProductsService         → products.go
26. RolesService            → roles.go
27. SegmentsService         → segments.go
28. ShortLinksService       → short_links.go
29. SourcesService          → sources.go
30. TagsService             → tags.go
31. TalksService            → talks.go
32. TasksService            → tasks.go
33. UnsortedService         → unsorted.go
34. UsersService            → users.go
35. WebhooksService         → webhooks.go
36. WebsiteButtonsService   → website_buttons.go
37. WidgetsService          → widgets.go
```

## Формат отчёта (toolmap_audit.md)

```markdown
# Аудит toolmap.md

## AccountService ✅
Соответствует

---

## LeadsService ⚠️

### Отсутствует в toolmap:
- `CreateOne(ctx, lead) (*Lead, error)` — convenience для одной записи
- `SyncOne(ctx, lead, with) (*Lead, error)` — create or update

### Неточности:
- `GetWithPagination` — в SDK это просто `Get()` с возвратом `PageMeta`
- `Delete` — метод есть, но возвращает ошибку (API не поддерживает)

---
```

## Пути файлов

- toolmap: `/Users/tihn/amo-ai-tgbot-go/app/gkit/tools/toolmap.md`
- SDK services: `/Users/tihn/amocrm-sdk-go/core/adapters/`
- Отчёт: `/Users/tihn/amo-ai-tgbot-go/doc/toolmap_audit.md`

## Начало

Скажи: **"Начинаю аудит. Первый раздел: AccountService. Продолжить?"**

После подтверждения:

1. `view_file_outline` на `account.go`
2. Сравни с разделом AccountService в toolmap.md
3. Запиши результат
4. Спроси про следующий раздел
