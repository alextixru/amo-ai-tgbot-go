# models/

Входные структуры (Input DTOs) для инструментов Genkit.

## Содержимое

| Файл | Структуры |
|------|-----------|
| `common.go` | `LinkTarget`, `ParentEntity` |
| `entities.go` | `EntitiesInput`, `EntitiesFilter`, `EntityData` |
| `activities.go` | `ActivitiesInput`, `ActivityData` |
| `customers.go` | `CustomersInput`, `CustomerFilter`, `CustomerData`, `CustomerTransactionData`, `CustomerLinkData` |
| `admin_schema.go` | `AdminSchemaInput`, `SchemaFilter` |
| `admin_pipelines.go` | `AdminPipelinesInput` |
| `admin_users.go` | `AdminUsersInput`, `AdminUsersFilter` |
| `admin_integrations.go` | `AdminIntegrationsInput`, `IntegrationsFilter` |
| `products.go` | `ProductsInput`, `ProductFilter`, `ProductData` |
| `catalogs.go` | `CatalogsInput`, `CatalogFilter`, `CatalogData`, `CatalogElementData` |
| `files.go` | `FilesInput`, `FileFilter`, `FileUploadParams` |
| `unsorted.go` | `UnsortedInput`, `UnsortedFilter`, `UnsortedAcceptParams`, `UnsortedLinkData` |
| `complex_create.go` | `ComplexCreateInput`, `LeadData`, `ContactData`, `CompanyData` |

## Принципы

- **Input** → структуры здесь (с `jsonschema_description`)
- **Output** → модели SDK напрямую (`github.com/alextixru/amocrm-sdk-go/core/models`)
