# Аудит сервисов Entities (Leads, Contacts, Companies)

Этот файл содержит результаты последовательного аудита каждого сервиса в папке `services/entities/` на соответствие `tools_schema.md` и возможностям SDK.

---

## leads.go
**Layer:** leads
**Schema actions:** search, get, create, update, sync, delete, link, unlink
**SDK service:** LeadsService (`core/services/leads.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `SearchLeads` | Реализована лишь часть фильтров |
| GetOne | ✅ | `GetLead` | Без `with` параметров |
| Create | ✅ | `CreateLead` | Одиночное создание (в SDK есть батч) |
| Update | ✅ | `UpdateLead` | Одиночное обновление (в SDK есть батч) |
| SyncOne | ✅ | `SyncLead` | |
| Delete | ⚠️ | `DeleteLead` | SDK возвращает ошибку (v4 не поддерживает прямое удаление) |
| Link | ✅ | `LinkLead` | |
| Unlink | ✅ | `UnlinkLead` | |
| AddComplex | ✅ | — | Используется в `complex_create` |

**Genkit Tool Handler:**
- ✅ Инструмент `entities` правильно распределяет вызовы для сделок.
- ❌ **Урезанная модель**: При создании/обновлении не передаются кастомные поля, теги или вложенные сущности (кроме sync).

**Статус:** ⚠️ Частично (основные действия есть, но фильтры и данные урезаны)

### Capabilities Coverage
**Filters:**
- ❌ SDK: `LeadsFilter` поддерживает `Price`, `Statuses`, `CreatedBy`, `ClosedAt`, `CustomFieldsValues` и др.
- ❌ Bot: AI не может найти сделки конкретного статуса или сделки с бюджетом > 100k.

**Parameters:**
- ⚠️ SDK: `models.Lead` поддерживает `with`: `catalog_elements`, `contacts`, `companies`, `loss_reason`, `source`.
- ❌ Bot: AI получает «голые» сделки без связанных контактов.

---

## contacts.go
**Layer:** contacts
**Schema actions:** search, get, create, update, sync, link, unlink, get_chats, link_chats
**SDK service:** ContactsService (`core/services/contacts.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `SearchContacts` | Ограниченные фильтры |
| GetOne | ✅ | `GetContact` | Без `with` параметров |
| Create | ✅ | `CreateContact` | |
| Update | ✅ | `UpdateContact` | |
| SyncOne | ✅ | `SyncContact` | |
| Link | ✅ | `LinkContact` | |
| Unlink | ✅ | `UnlinkContact` | |
| GetChats | ✅ | `GetContactChats` | |
| LinkChats | ✅ | `LinkContactChats` | |

**Genkit Tool Handler:**
- ✅ Инструмент `entities` правильно обрабатывает контакты и чаты.
- ❌ **Урезанная модель**: Кастомные поля и теги игнорируются при создании/обновлении.

**Статус:** ⚠️ Частично (фильтры и данные урезаны)

### Capabilities Coverage
**Filters:**
- ❌ SDK: `ContactsFilter` поддерживает `IDs`, `Names`, `CreatedBy`, `CreatedAt`, `CustomFieldsValues`.
- ❌ Bot: AI не может найти контакт по телефону (через CF) или по дате создания.

**Parameters:**
- ⚠️ SDK: `models.Contact` поддерживает `with`: `leads`, `company`, `customers`, `catalog_elements`, `social_profiles`.
- ❌ Bot: AI получает контакты без информации о связанных сделках или компании.

---

## companies.go
**Layer:** companies
**Schema actions:** search, get, create, update, sync, link, unlink
**SDK service:** CompaniesService (`core/services/companies.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `SearchCompanies` | Ограниченные фильтры |
| GetOne | ✅ | `GetCompany` | Без `with` параметров |
| Create | ✅ | `CreateCompany` | |
| Update | ✅ | `UpdateCompany` | |
| SyncOne | ✅ | `SyncCompany` | |
| Link | ✅ | `LinkCompany` | |
| Unlink | ✅ | `UnlinkCompany` | |
| GetLinks | ✅ | `GetLinks` | |

**Genkit Tool Handler:**
- ✅ Инструмент `entities` правильно обрабатывает компании.
- ❌ **Урезанная модель**: Кастомные поля и теги игнорируются при создании/обновлении.

**Статус:** ⚠️ Частично (фильтры и данные урезаны)

### Capabilities Coverage
**Filters:**
- ❌ SDK: `CompaniesFilter` поддерживает `IDs`, `Names`, `CreatedBy`, `CreatedAt`, `CustomFieldsValues`.
- ❌ Bot: AI не может искать компании по кастомным полям.

**Parameters:**
- ⚠️ SDK: `models.Company` поддерживает `with`: `leads`, `contacts`, `customers`, `catalog_elements`.
- ❌ Bot: AI получает компании без информации о связанных сделках.

---
