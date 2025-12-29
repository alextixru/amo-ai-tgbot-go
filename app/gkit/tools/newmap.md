# Задача
Покрыть все методы SDK оптимизированными инструментами.

---

# План: Bottom-Up Approach

## Этап 1: Каталогизация
**Цель:** Полный список всех сервисов и методов SDK.

**Результат:** Таблица:
| Service | Method | Params | Operation Type |
|---------|--------|--------|----------------|
| LeadsService | Get | filter | list |
| LeadsService | GetOne | id | get |
| ... | ... | ... | ... |

**Источник:** `toolmap.md` (уже есть).

---

## Этап 2: Классификация операций
**Цель:** Определить уникальные типы операций (actions).

**Результат:** Список actions:
- `list` / `search` — получение списка
- `get` — получение по ID
- `create` — создание
- `update` — обновление
- `delete` — удаление
- `link` / `unlink` — связывание сущностей
- `execute` — специфические действия (complete, accept, close, install...)

---

## Этап 3: Группировка по доменам
**Цель:** Объединить сервисы в логические группы.

**Результат:** Группы:
| Группа | Сервисы | Описание |
|--------|---------|----------|
| **Core Entities** | Leads, Contacts, Companies, Customers | Основные бизнес-объекты |
| **Reference** | Users, Pipelines, Tags, Roles, CustomFields | Справочники (read-only) |
| **Activities** | Tasks, Notes, Calls, Events | Активности и история |
| **Communication** | Talks, ChatTemplates, Unsorted | Общение с клиентами |
| **Catalog** | Catalogs, Products, CatalogElements | Товары и каталоги |
| **Admin** | Webhooks, Widgets, Sources, Files | Настройки и интеграции |

---

## Этап 4: Проектирование инструментов
**Цель:** Для каждой группы создать unified tool с общей Input Schema.

**Результат:** Список инструментов:
1. `crm_entities` — Core Entities (leads, contacts, companies, customers)
2. `crm_reference` — Reference data (users, pipelines, tags...)
3. `crm_activities` — Tasks, Notes, Calls
4. `crm_communication` — Talks, Unsorted, ChatTemplates
5. `crm_catalog` — Products, Catalogs
6. `crm_admin` — Webhooks, Widgets, Files

---

## Этап 5: Матрица покрытия
**Цель:** Таблица соответствия Service.Method → Tool.Action.

**Результат:**
| Service | Method | Tool | Action | ✓ |
|---------|--------|------|--------|---|
| LeadsService | Get | crm_entities | search | ✓ |
| LeadsService | Link | crm_entities | link | ✓ |
| TasksService | Complete | crm_activities | complete | ✓ |
| ... | ... | ... | ... | |

---

## Этап 6: Реализация
**Цель:** Кодирование каждого инструмента.

**Порядок:**
1. `crm_entities` (самый важный, CRUD для Core)
2. `crm_reference` (уже частично есть)
3. `crm_activities` (Tasks, Notes)
4. Остальные по приоритету

---

# Следующий шаг
Начать с **Этапа 1-2**: создать полную таблицу сервисов и классифицировать операции.
