#Audit: Admin Users Service

Этот файл содержит результаты аудита папки `services/admin_users/` на соответствие `tools_schema.md` и возможностям SDK.

---

##users.go

**Layer:** users

**Schema actions:** search, list, get, create

**SDK service:**`UsersService` (`core/services/users.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |

|-----------|----------------------|----------------|-------------|

| Get | ✅ |`ListUsers`| Вызывается без `with` параметров. |

| GetOne | ✅ |`GetUser`| Вызывается без `with` параметров. |

| Create | ✅ |`CreateUsers`||

| Update | ❌ | — |**Не поддерживается API v4**. |

| UpdateOne | ❌ | — |**Не поддерживается API v4**. |

**Gaps:**

-**With параметры**: SDK поддерживает `role`, `uuid`, `group`, `amojo_id`, `user_rank`, `phone_number`. Бот их не запрашивает, что может ограничить AI в понимании прав пользователя или его чат-идентификатора.

---

##roles.go

**Layer:** roles

**Schema actions:** search, list, get, create, update, delete

**SDK service:**`RolesService` (`core/services/roles.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |

|-----------|----------------------|----------------|-------------|

| Get | ✅ |`ListRoles`||

| GetOne | ✅ |`GetRole`||

| Create | ✅ |`CreateRoles`||

| Update | ✅ |`UpdateRoles`||

| Delete | ✅ |`DeleteRole`| Один за раз. |

**Gaps:**

-**With параметры**: SDK поддерживает `users`. Бот не использует.

---

##Genkit Tool Handler (`app/gkit/tools/admin_users.go`)

**Findings:**

- ✅ Хорошая реализация разделения прав и действий.
- ✅ **Обработка ограничений API**: Явная блокировка `update`/`delete` для пользователей с сообщением об ошибке.
- ⚠️ **Data mapping**: Используется консистентный батч-подход (`input.Data["users"]`).

**Статус:** ✅ Хорошо (соответствует ограничениям платформы)
