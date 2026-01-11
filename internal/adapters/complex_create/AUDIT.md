# Аудит сервисов Complex Create

Этот файл содержит результаты последовательного аудита папки `adapters/complex_create/` на соответствие `tools_schema.md` и возможностям SDK.

---

## complex.go
**Layer:** complex
**Schema actions:** create
**SDK service:** LeadsService (`core/adapters/leads.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| AddComplex | ✅ | `CreateComplexBatch` | Батч-создание |
| AddOneComplex | ✅ | `CreateComplex` | Одиночное создание |

**Genkit Tool Handler:**
- ✅ **Полный маппинг данных**: Инструмент `complex_create.go` поддерживает все основные поля сделки, контактов и компании.
- ✅ **Custom Fields сделки**: Поддержка кастомных полей через `custom_fields_values`.
- ✅ **Теги сделки**: Поддержка тегов через поле `tags`.
- ✅ **Phone/Email контактов**: Телефон и email маппятся как кастомные поля PHONE/EMAIL.
- ✅ **Батч-режим**: Добавлен инструмент `complex_create_batch` для пакетного создания (до 50 сделок).

**Статус:** ✅ Полностью реализовано
**Обновлено:** 2026-01-11

### Capabilities Coverage

**Data Richness:**
- ✅ SDK: `ComplexLeadResult` возвращает `ContactID` и `CompanyID`.
- ✅ Bot: Сервис возвращает результат целиком, AI может увидеть созданные ID.
- ✅ Bot: AI может полноценно использовать комплексное создание с передачей всех данных.

**Поддерживаемые поля:**
| Сущность | Поля |
|----------|------|
| Lead | `name`, `price`, `pipeline_id`, `status_id`, `responsible_user_id`, `custom_fields_values`, `tags` |
| Contact | `name`, `first_name`, `last_name`, `phone`, `email`, `is_main`, `responsible_user_id`, `custom_fields_values` |
| Company | `name`, `responsible_user_id`, `custom_fields_values` |

---
