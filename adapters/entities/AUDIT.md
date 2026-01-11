# Audit: Entities Service (Leads, Contacts, Companies)

**Status:** ✅ Refactored (v4 alignment)
**Goal:** Align with amoCRM SDK v4, support batch operations, rich data, and full filtering.

## Summary of Improvements

- **Rich Data:** DTOs now include `custom_fields_values`, `tags`, and `_embedded` data.
- **Full Filtering:** implemented comprehensive filtering for all entities, including IDs, Names, date ranges, and custom fields.
- **Batch Operations:** added `Create*s` and `Update*s` methods for batch API calls.
- **With-parameters:** added support for `with` parameters in all search and get operations.
- **Clean Interface:** removed non-functional `DeleteLead` method.

## Leads

| Action | Status | Method | Notes |
| :--- | :---: | :--- | :--- |
| Search | ✅ | `SearchLeads` | Supports all SDK filters and `with` params. |
| Get | ✅ | `GetLead` | Supports `with` params. |
| Create | ✅ | `CreateLead` / `CreateLeads` | Supports batch creation and rich data (CFV, Tags, Embedded). |
| Update | ✅ | `UpdateLead` / `UpdateLeads` | Supports batch update and rich data. |
| Sync | ✅ | `SyncLead` | Full sync with rich data support. |
| Delete | ⚠️ | - | Removed. API limitations handled at AI prompt level. |

## Contacts

| Action | Status | Method | Notes |
| :--- | :---: | :--- | :--- |
| Search | ✅ | `SearchContacts` | Supports all SDK filters and `with` params. |
| Get | ✅ | `GetContact` | Supports `with` params. |
| Create | ✅ | `CreateContact` / `CreateContacts` | Supports batch creation and rich data. |
| Update | ✅ | `UpdateContact` / `UpdateContacts` | Supports batch update and rich data. |
| Sync | ✅ | `SyncContact` | Full sync with rich data support. |
| Get Chats | ✅ | `GetContactChats` | Standard implementation. |

## Companies

| Action | Status | Method | Notes |
| :--- | :---: | :--- | :--- |
| Search | ✅ | `SearchCompanies` | Supports all SDK filters and `with` params. |
| Get | ✅ | `GetCompany` | Supports `with` params. |
| Create | ✅ | `CreateCompany` / `CreateCompanies` | Supports batch creation and rich data. |
| Update | ✅ | `UpdateCompany` / `UpdateCompanies` | Supports batch update and rich data. |
| Sync | ✅ | `SyncCompany` | Full sync with rich data support. |
