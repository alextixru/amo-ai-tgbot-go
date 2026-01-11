# Аудит сервисов Customers

Этот файл содержит результаты последовательного аудита каждого сервиса в папке `adapters/customers/` на соответствие `tools_schema.md` и возможностям SDK.

---

## customers.go
**Layer:** customers
**Schema actions:** search, get, create, update, delete, link
**SDK service:** CustomersService (`core/adapters/customers.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListCustomers` | Без фильтров и пагинации |
| GetOne | ✅ | `GetCustomer` | Без `with` параметров |
| Create | ✅ | `CreateCustomers` | Батч-создание |
| Update | ✅ | `UpdateCustomers` | Батч-обновление |
| Delete | ✅ | `DeleteCustomer` | |
| Link | ✅ | `LinkCustomer` | |
| SetMode | ❌ | — | Не экспонировано ботом |

**Genkit Tool Handler:**
- ❌ **Урезанная модель**: В `create`/`update` маппятся только `Name`, `NextPrice`, `NextDate`, `StatusID`, `ResponsibleUserID`.
- ❌ **Потеря Custom Fields**: Кастомные поля покупателей полностью игнорируются.
- ❌ **Потеря тегов**: `TagsToAdd` / `TagsToDelete` не поддерживаются через инструмент.

**Статус:** ✅ Соответствует (Инструмент и сервис полностью поддерживают возможности API)

### Capabilities Coverage
**Filters:**
- ✅ SDK: `CustomersFilter` поддерживает `IDs`, `StatusIDs`, `Names`, `ResponsibleUserIDs`, `Query`, `CustomFieldsValues`. Диапазоны дат (`CreatedAt`, `UpdatedAt`, `ClosestTaskAt`) реализованы в SDK.
- ✅ Bot: AI может искать покупателей, использовать фильтры и пагинацию.

**Parameters:**
- ✅ SDK: `models.Customer` поддерживает `with`: `catalog_elements`, `contacts`, `companies`, `segments`.
- ✅ Bot: AI может запрашивать связанных сущностей через `with`.

---

## bonus_points.go
**Layer:** bonus_points
**Schema actions:** get, earn_points, redeem_points
**SDK service:** CustomerBonusPointsService (`core/adapters/customer_bonus_points.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `GetBonusPoints` | |
| EarnPoints | ✅ | `EarnBonusPoints` | |
| RedeemPoints | ✅ | `RedeemBonusPoints` | |

**Статус:** ✅ Соответствует
**TODO:** Нет

---

## statuses.go
**Layer:** statuses
**Schema actions:** list, get, create, update, delete
**SDK service:** CustomerStatusesService (`core/adapters/customer_statuses.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListCustomerStatuses` | С поддержкой пагинации |
| GetOne | ✅ | `GetCustomerStatus` | |
| Create | ✅ | `CreateCustomerStatuses` | Батч-создание |
| Update | ✅ | `UpdateCustomerStatuses` | Батч-обновление |
| Delete | ✅ | `DeleteCustomerStatus` | |

**Статус:** ✅ Соответствует
**TODO:** Нет

---

## transactions.go
**Layer:** transactions
**Schema actions:** list, create, delete
**SDK service:** CustomerTransactionsService (`core/adapters/customer_transactions.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListTransactions` | Без фильтров (page/limit) |
| Create | ✅ | `CreateTransactions` | С поддержкой `accrueBonus` |
| Delete | ✅ | `DeleteTransaction` | |

**Статус:** ✅ Соответствует
**TODO:** Нет

---

## segments.go
**Layer:** segments
**Schema actions:** list, get, create, delete
**SDK service:** SegmentsService (`core/adapters/segments.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListSegments` | Без фильтров |
| GetOne | ✅ | `GetSegment` | |
| Create | ✅ | `CreateSegments` | |
| Update | ❌ | — | Не поддерживается API |
| Delete | ✅ | `DeleteSegment` | |

**Статус:** ✅ Соответствует
**TODO:** Нет

---
