# Аудит: admin_users

## Вход (LLM → tool → сервис)

### Поля, которые сейчас числовые ID, но должны быть именами

**`ID int` (users/roles)**
Единственный идентификатор для `get` и `delete` — числовой. LLM не знает ID заранее, работает с именами. Вынуждена сначала вызывать `list`, вытаскивать ID, потом делать нужный вызов.
Нужно: `name string` как альтернатива, резолвинг в сервисе.

**`UserRights.GroupID int` / `UserRights.RoleID int` в ответе**
В `User.Rights` числа без расшифровки. LLM не может интерпретировать без отдельного `list roles`.

**`UserRights.StatusRights[].PipelineID int` / `.StatusID int`**
Права на статусы — числовые ID воронок и статусов без названий.

**`UserRights.CatalogRights[].CatalogID int`**
Права на каталоги — числовой ID без названия каталога.

### Поля, которые отсутствуют, но нужны LLM

- **Нет фильтра по имени/email** — `AdminUsersFilter` содержит только `Limit`, `Page`, `With`. LLM вынуждена грузить всю страницу и фильтровать на стороне. (Ограничение API — нужно задокументировать явно.)
- **`add_to_group` не реализован** — `tools_schema.md` декларирует action, но handler падает в `default` с "unknown action". Нет ни реализации, ни явной ошибки "not implemented".
- **`user_id` и `group_id`** — поля из `tools_schema.md` для `add_to_group` отсутствуют в `AdminUsersInput`.
- **Нет `Order` в фильтре** — `BaseFilter` SDK поддерживает `Order map[string]string`, в `AdminUsersFilter` нет. Нельзя запросить "последних добавленных пользователей".

### Неудобные структуры/типы

- **`Data map[string]any` для create/update** — двойная сериализация хрупка, Genkit генерирует JSON Schema как `object` без свойств. LLM не знает допустимых полей.
- **`With []string` — детали HTTP API видны LLM** — LLM не должна знать про `with=role,uuid,group`. Сервис должен сам решать что подгружать.
- **Нет типизированного поля `Password`** — `User.Password` помечен `json:"-"`. LLM не знает как передать пароль при create.
- **`search` как алиас `list`** — нигде не задокументирован для LLM, только в `case`.

---

## Выход (сервис → tool → LLM)

### Что сейчас возвращается

| Метод | Тип возврата |
|-------|-------------|
| `ListUsers` | `[]*amomodels.User` |
| `GetUser` | `*amomodels.User` |
| `CreateUsers` | `[]*amomodels.User` |
| `ListRoles` | `[]*amomodels.Role` |
| `GetRole` | `*amomodels.Role` |
| `CreateRoles` / `UpdateRoles` | `[]*amomodels.Role` |
| `DeleteRole` | `error` → tool возвращает `nil, nil` (нет подтверждения!) |

Все методы — голые SDK-структуры без маппинга.

### Какие SDK With-параметры не используются

**`GetUser` игнорирует `with` полностью:**
```go
func (s *service) GetUser(ctx context.Context, id int, with []string) (*amomodels.User, error) {
    return s.sdk.Users().GetOne(ctx, id)  // with молча игнорируется
}
```
SDK поддерживает `WithRelations(...)`: `role`, `uuid`, `group`, `amojo_id`, `user_rank`, `phone_number` — все шесть параметров игнорируются.

**`GetRole` игнорирует `with` полностью:**
```go
func (s *service) GetRole(ctx context.Context, id int, with []string) (*amomodels.Role, error) {
    return s.sdk.Roles().GetOne(ctx, id)  // with молча игнорируется
}
```
SDK поддерживает `WithRelations("users")` — список пользователей роли. Игнорируется.

**`ListUsers`** — `With` присваивается в `sdkFilter`, но `UsersFilter` не наследует полноценный `BaseFilter.With` — работает корректно только для list, не для GetOne.

### Числовые ID в ответе, которые LLM не может интерпретировать

| Поле | Проблема |
|------|----------|
| `User.Rights.RoleID int` | ID роли без названия |
| `User.Rights.GroupID int` | ID группы без названия |
| `User.Rights.StatusRights[].PipelineID int` | ID воронки без названия |
| `User.Rights.StatusRights[].StatusID int` | ID статуса без названия |
| `User.Rights.CatalogRights[].CatalogID int` | ID каталога без названия |
| `UserEmbedded.Roles[].ID int` | без названия если `with=role` не запрошен |

Без `with=role,group` LLM видит пустые `_embedded` или не видит вообще.

### Что теряется по сравнению с тем, что SDK может вернуть

1. **`GetUser`/`GetRole` с `WithRelations`** — никогда не вызывается. Роль пользователя по имени, список пользователей роли — недоступны.
2. **`PageMeta`** — отбрасывается при `ListUsers`/`ListRoles`. LLM не знает что список неполный.
3. **`User.UUID`** — заполняется только при `with=uuid`, по умолчанию пусто.
4. **`User.Rank`** — только при `with=user_rank`.
5. **`User.PhoneNumber`** — только при `with=phone_number`.
6. **`DeleteRole` возвращает `nil, nil`** — LLM не получает подтверждения удаления.
7. **Нет валидации `create`** — обязательные поля (`name`, `email`) не проверяются до вызова API.

---

## Итого

**Приоритет рефакторинга: высокий**

Сервис работает, но как адаптерный слой плох: `with` игнорируется в `GetOne`, пагинация теряется, схема расходится с реализацией (`add_to_group`).

### Список конкретных изменений

**Критические:**
1. Исправить `GetUser` — передавать `with` в SDK: `GetOne(ctx, id, services.WithRelations(with...))`
2. Исправить `GetRole` — аналогично
3. Реализовать или явно запретить `add_to_group` с понятной ошибкой
4. `DeleteRole` — возвращать `{"success": true, "deleted_id": id}` вместо `nil`

**Высокий приоритет:**
5. Добавить `UserID` и `GroupID` в `AdminUsersInput` (поля из схемы которые отсутствуют в структуре)
6. Вернуть `PageMeta` в ответе list-методов — `{items, page, total, has_next}`
7. Убрать `With []string` из входной модели — сервис сам решает что подгружать. Дефолт для `list`: `role`, `group`.
8. Добавить `Name string` / `Email string` в `AdminUsersFilter` (client-side фильтрация с пометкой в документации)
9. Расшифровать числовые ID в ответе: `role_name` рядом с `role_id`, `group_name` рядом с `group_id` через DTO `UserView`

**Средний приоритет:**
10. Автоматически запрашивать `with=role,group` при `ListUsers` по умолчанию
11. Добавить `Order` в `AdminUsersFilter` — пробросить в SDK
12. Добавить валидацию `create` на обязательные поля (`name`, `email`) до вызова API
13. Ввести типизированный input для create: `Users []UserCreateInput`, `Roles []RoleCreateInput`
14. Задокументировать `search` как алиас `list` либо убрать из handler
