# Audit: Unsorted Service

Этот файл содержит результаты аудита папки `adapters/unsorted/` на соответствие `tools_schema.md` и возможностям SDK.

---

## unsorted.go

**Layer:** unsorted

**Schema actions:** search, get, accept, decline, link, summary

**SDK service:** `UnsortedService` (`core/adapters/unsorted.go`)

| Create | ✅ | `CreateUnsorted` | Добавление заявок (SIP/Forms/Chats) реализовано. |

**Особенности:**
- ✅ **Summary**: Позволяет AI быстро получить обзор текущего состояния очереди неразобранного.
- ✅ **Accept/Decline/Link**: Все основные способы трансформации заявки присутствуют.
- ✅ **Create**: Поддерживается программное добавление заявок через батч-операцию.

---

## Genkit Tool Handler (`app/gkit/tools/unsorted.go`)

**Findings:**
- ✅ Хорошее покрытие действий: search, get, create, accept, decline, link, summary.
- ✅ **Create Action**: Добавлена поддержка создания записей в Неразобранном.
- ✅ **Data mapping**: Всё достаточно прямолинейно, поддерживается передача основных полей.

**Статус:** ✅ Отлично (все возможности SDK покрыты).

---

## Capabilities Coverage

**Filters:**
- ✅ SDK: `UnsortedFilter` поддерживает `UIDs`, `Category`, `PipelineID`.
- ✅ Bot: Использует все три (UID в `GetUnsorted`, категорию и воронку в `ListUnsorted`).

**Summary Filters:**
- ✅ SDK: `UnsortedSummaryFilter` поддерживает `PipelineID`.
- ✅ Bot: `SummaryUnsorted` использует его.

**Batch Operations:**
- ❌ SDK: Методы `Accept`, `Decline`, `Link` в API v4 работают строго с одним UID. Пакетная обработка невозможна.
- ✅ SDK: `Create` поддерживает массив элементов. Но в боте не реализовано.

---

## Рекомендации

1. **Добавить Create**: Если AI должен уметь "закидывать" заявки на проверку менеджерам (например, из сомнительных логов чатов), стоит реализовать `CreateUnsorted`.
2. **UID Validation**: Добавить более явную проверку формата UID (обычно это 32-символьная строка).
