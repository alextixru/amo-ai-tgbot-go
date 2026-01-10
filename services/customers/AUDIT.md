# Аудит сервисов Customers

Этот файл содержит результаты последовательного аудита каждого сервиса в папке `services/customers/` на соответствие `tools_schema.md` и возможностям SDK.

---

## customers.go
**Layer:** customers
**Schema actions:** search, get, create, update, delete, link
**SDK service:** CustomersService (`core/services/customers.go`)

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

**Статус:** ⚠️ Частично (сервис ок, инструмент урезан, фильтры отсутствуют)

### Capabilities Coverage
**Filters:**
- ❌ SDK: `CustomersFilter` поддерживает `IDs`, `StatusIDs`, `Names`, `ResponsibleUserIDs`, `Query`, `CustomFieldsValues`, диапазоны дат.
- ❌ Bot: AI не может искать покупателя по телефону (через CF) или имени.

**Parameters:**
- ⚠️ SDK: `models.Customer` поддерживает `with`: `catalog_elements`, `contacts`, `companies`, `segments`.
- ❌ Bot: AI получает покупателей без связанных контактов и сегментов.

---

## bonus_points.go
**Layer:** bonus_points
**Schema actions:** get, earn_points, redeem_points
**SDK service:** CustomerBonusPointsService (`core/services/customer_bonus_points.go`)

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
**SDK service:** CustomerStatusesService (`core/services/customer_statuses.go`)

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
**SDK service:** CustomerTransactionsService (`core/services/customer_transactions.go`)

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
**SDK service:** SegmentsService (`core/services/segments.go`)

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
