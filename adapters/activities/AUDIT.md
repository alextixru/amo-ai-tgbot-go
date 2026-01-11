# Аудит сервисов Activities

Этот файл содержит результаты последовательного аудита каждого сервиса в папке `adapters/activities/` на соответствие `tools_schema.md` и возможностям SDK.

---

## calls.go
**Layer:** calls
**Schema actions:** create (write-only)
**SDK service:** CallsService (`core/adapters/calls.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| CreateOne | ✅ | `CreateCall` | Соответствует схеме и SDK |
| Get | ❌ | - | Заблокировано в SDK (ErrNotAvailableForAction) |
| GetOne | ❌ | - | Заблокировано в SDK (ErrNotAvailableForAction) |
| Update | ❌ | - | Заблокировано в SDK (ErrNotAvailableForAction) |

**Статус:** ✅ Полностью соответствует (с учетом ограничений API)

---

## events.go
**Layer:** events
**Schema actions:** list, get (read-only)
**SDK service:** EventsService (`core/adapters/events.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListEvents` | С поддержкой `EventsFilter` |
| GetOne | ✅ | `GetEvent` | |

**Статус:** ✅ Полностью соответствует

### Capabilities Coverage

**Filters:**
| SDK Filter Method | Bot Field | Status | Comment |
|-------------------|-----------|--------|---------|
| SetLimit | filter.limit | ✅ | |
| SetPage | filter.page | ✅ | |
| SetTypes | filter.types | ✅ | |
| SetEntity | parent.type | ✅ | |
| SetEntityIDs | parent.id | ✅ | |
| SetCreatedBy | filter.created_by | ✅ | |

---

## files.go
**Layer:** files
**Schema actions:** list, link, unlink
**SDK service:** EntityFilesService (`core/adapters/entity_files.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListFiles` | С поддержкой `FilesFilter` |
| Link | ✅ | `LinkFiles` | |
| Unlink | ✅ | `UnlinkFile` | |

**Статус:** ✅ Полностью соответствует

### Capabilities Coverage

**Filters:**
| SDK Filter Method | Bot Field | Status | Comment |
|-------------------|-----------|--------|---------|
| SetLimit | filter.limit | ✅ | |
| SetPage | filter.page | ✅ | |

---

## links.go
**Layer:** links
**Schema actions:** list, link, unlink
**SDK service:** LinksService (`core/adapters/links.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListLinks` | С поддержкой `LinksFilter` |
| Link | ✅ | `LinkEntities` | Поддержка батч-линковки |
| Unlink | ✅ | `UnlinkEntity` | |

**Статус:** ✅ Полностью соответствует

### Capabilities Coverage

**Batch Operations:**
- ✅ SDK: `Link`/`Unlink` принимают массив.
- ✅ Bot: Поддерживает `LinksTo []LinkTarget` для массового связывания.

---

## notes.go
**Layer:** notes
**Schema actions:** list, get, get_by_parent, create, update
**SDK service:** NotesService (`core/adapters/notes.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| GetByParent| ✅ | `ListNotes` | С поддержкой `NotesFilter` и `With` |
| GetOne | ✅ | `GetNote` | |
| Create | ✅ | `CreateNotes` | Поддержка массового создания |
| Update | ✅ | `UpdateNote` | |

**Статус:** ✅ Полностью соответствует

### Capabilities Coverage

**Batch Operations:**
- ✅ SDK: `Create`/`Update` принимают массив.
- ✅ Bot: Добавлен метод `CreateNotes` для массового создания.

---

## subscriptions.go
**Layer:** subscriptions
**Schema actions:** list, subscribe, unsubscribe
**SDK service:** EntitySubscriptionsService (`core/adapters/entity_subscriptions.go`)

**Статус:** ✅ Полностью соответствует

---

## tags.go
**Layer:** tags
**Schema actions:** list, create, delete
**SDK service:** TagsService (`core/adapters/tags.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListTags` | С поддержкой `TagsFilter` |
| Create | ✅ | `CreateTags` | Поддержка массового создания |
| Update | ❌ | - | |
| Delete | ⚠️ | `DeleteTag` | Ограничение API v4 (использовать PATCH сущности) |

**Статус:** ✅ Полностью соответствует (с учетом ограничений API)

---

## talks.go
**Layer:** talks
**Schema actions:** close
**SDK service:** TalksService (`core/adapters/talks.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Close | ✅ | `CloseTalk` | |

**Статус:** ✅ Полностью соответствует

---

## tasks.go
**Layer:** tasks
**Schema actions:** list, get, create, update, complete
**SDK service:** TasksService (`core/adapters/tasks.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListTasks` | С полной фильтрацией и `With` |
| GetOne | ✅ | `GetTask` | С поддержкой `With` |
| Create | ✅ | `CreateTasks` | Массовое создание |
| Update | ✅ | `UpdateTask` | |
| Complete | ✅ | `CompleteTask` | |

**Статус:** ✅ Полностью соответствует

### Capabilities Coverage

**Filters:**
- ✅ `Page`, `Order`, `IDs`, `UpdatedAt` добавлены.
- ✅ `With` параметры поддерживаются в `ListTasks` и `GetTask`.
- ✅ Батч-операции: `CreateTasks` добавлен.
