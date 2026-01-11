#Audit: Admin Users Service

Этот файл содержит результаты аудита папки `adapters/admin_users/` на соответствие `tools_schema.md` и возможностям SDK.

---

##users.go

**Layer:** users

**Schema actions:** search, list, get, create

**SDK service:**`UsersService` (`core/adapters/users.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |

|-----------|----------------------|----------------|-------------|

| Get | ✅ |`ListUsers`| Поддерживает `with` и фильтрацию. |
| GetOne | ✅ |`GetUser`| Поддерживает `with` параметры. |

| Create | ✅ |`CreateUsers`||

| Update | ❌ | — |**Не поддерживается API v4**. |

| UpdateOne | ❌ | — |**Не поддерживается API v4**. |

**Gaps:**
- ✅ **With параметры**: Реализована поддержка `role`, `uuid`, `group`, `amojo_id`, `user_rank`, `phone_number`.
- ✅ **Фильтрация**: Реализована поддержка `limit` и `page`.

---

##roles.go

**Layer:** roles

**Schema actions:** search, list, get, create, update, delete

**SDK service:**`RolesService` (`core/adapters/roles.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |

|-----------|----------------------|----------------|-------------|

| Get | ✅ |`ListRoles`||

| GetOne | ✅ |`GetRole`||

| Create | ✅ |`CreateRoles`||

| Update | ✅ |`UpdateRoles`||

| Delete | ✅ |`DeleteRole`| Один за раз. |

**Gaps:**
- ✅ **With параметры**: Реализована поддержка `users`.

---

##Genkit Tool Handler (`app/gkit/tools/admin_users.go`)

**Findings:**

- ✅ Хорошая реализация разделения прав и действий.
- ✅ **Обработка ограничений API**: Явная блокировка `update`/`delete` для пользователей с сообщением об ошибке.
- ⚠️ **Data mapping**: Используется консистентный батч-подход (`input.Data["users"]`).

**Статус:** ✅ Хорошо (соответствует ограничениям платформы)
