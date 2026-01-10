# Audit: Unsorted Service

Этот файл содержит результаты аудита папки `services/unsorted/` на соответствие `tools_schema.md` и возможностям SDK.

---

## unsorted.go

**Layer:** unsorted

**Schema actions:** search, get, accept, decline, link, summary

**SDK service:** `UnsortedService` (`core/services/unsorted.go`)

| Метод SDK | Реализован в сервисе | Метод сервиса | Комментарий |
|-----------|----------------------|----------------|-------------|
| Get | ✅ | `ListUnsorted` | Использует фильтры `Category`, `PipelineID`. |
| GetOne | ✅ | `GetUnsorted` | Через `Get` с фильтром по `UID`. |
| Accept | ✅ | `AcceptUnsorted` | Принимает `user_id` и `status_id`. |
| Decline | ✅ | `DeclineUnsorted` | Принимает `user_id`. |
| Link | ✅ | `LinkUnsorted` | Принимает `lead_id`. |
| Summary | ✅ | `SummaryUnsorted` | Статистика по воронке. |
| Create | ❌ | — | Добавление заявок (SIP/Forms) не реализовано. |

**Особенности:**
- ✅ **Summary**: Позволяет AI быстро получить обзор текущего состояния очереди неразобранного.
- ✅ **Accept/Decline/Link**: Все основные способы трансформации заявки присутствуют.

---

## Genkit Tool Handler (`app/gkit/tools/unsorted.go`)

**Findings:**
- ✅ Хорошее покрытие действий: search, get, accept, decline, link, summary.
- ❌ **Create Action**: В инструменте нет действия `create`. SDK позволяет программно добавлять записи в Неразобранное (например, для источников, не имеющих прямой интеграции).
- ⚠️ **Data mapping**: Всё достаточно прямолинейно, так как Неразобранное больше про действия (UID + params), чем про сложную структуру данных.

**Статус:** ✅ Хорошо (основной рабочий процесс по разбору заявок полностью закрыт).

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
