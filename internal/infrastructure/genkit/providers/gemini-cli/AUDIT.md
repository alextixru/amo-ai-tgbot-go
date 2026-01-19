# Аудит Провайдера Gemini CLI: Обнаруженные Баги и Проблемы

**Дата аудита:** 2026-01-19  
**Аудитор:** Expert Golang Developer & Security Auditor

---

## Сводка Проблем

- **Критические баги:** 3 (2 исправлено)
- **Высокий приоритет:** 7
- **Средний приоритет:** 6

---

## 1. OAuth и Безопасность

### 1.1 ✅ **ИСПРАВЛЕНО: Персистентность Токенов**

Реализовано: Keyring Storage + retry с exponential backoff.

---

### 1.2 ✅ **ИСПРАВЛЕНО: User Management в OAuth Flow**

Реализовано: `UserAccountManager`, `fetchAndCacheUserInfo`, Headless OAuth API.

---

## 2. API Совместимость

### 2.1 ✅ **ИСПРАВЛЕНО: Поля GenerationConfig**

Реализовано: Добавлены все 21 поле в структуру `GenerationConfig` и маппинг в `Generate()`/`GenerateStream()`.

---

### 2.2 ✅ **ИСПРАВЛЕНО: Jitter в Retry**

Реализовано: Добавлен jitter (0-50% от delay) в `postWithRetry` для предотвращения thundering herd.

---

## 3. Конвертеры Типов

### 3.1 ✅ **ИСПРАВЛЕНО: Разрешающий Fallback в Схемах**

**Расположение:** [gemini_genai.go](file:///Users/tihn/amo-ai-tgbot-go/internal/infrastructure/genkit/providers/gemini-cli/gemini_genai.go)

**Реализовано:**
- Убран permissive fallback на `TypeString`
- Добавлена поддержка `oneOf` (паритет с `anyOf`)
- Добавлено разрешение `$ref` внутри `anyOf`/`oneOf`
- Добавлен вывод типа из структуры (`properties`/`additionalProperties` → object)
- Написаны тесты для 12 моделей инструментов

---

## 4. Телеметрия

### ⏭️ **НЕ ПРИМЕНИМО для данного use-case**

Телеметрия (`RecordConversationOffered`, `RecordConversationInteraction`) — это внутренняя аналитика Google для IDE-плагинов. Для Telegram-бота не требуется.

---

## 5. Отсутствующие Функции

### 5.1 ⚠️ **СРЕДНИЙ: Отсутствует Experiments API**

**Расположение:** [client.go](file:///Users/tihn/amo-ai-tgbot-go/internal/infrastructure/genkit/providers/gemini-cli/client.go)

**Проблема:** TypeScript поддерживает experiments API, Go — нет.

**TypeScript-эталон:**
```typescript
async listExperiments(metadata: ClientMetadata): Promise<ListExperimentsResponse> {
    if (!this.projectID) {
        throw new Error('projectId is not defined for CodeAssistServer.');
    }
    const req: ListExperimentsRequest = {
        project: this.projectID,
        metadata: { ...metadata, duetProject: this.projectID },
    };
    return this.requestPost<ListExperimentsResponse>('listExperiments', req);
}
```

**Последствия:** Невозможность использования экспериментальных функций API

---

### 5.2 ⚠️ **НИЗКИЙ: Отсутствует EmbedContent**

**Расположение:** [client.go](file:///Users/tihn/amo-ai-tgbot-go/internal/infrastructure/genkit/providers/gemini-cli/client.go)

**Проблема:** Отсутствует метод `embedContent` (хотя в TypeScript он тоже не реализован).

**TypeScript-эталон:**
```typescript
async embedContent(_req: EmbedContentParameters): Promise<EmbedContentResponse> {
    throw Error();
}
```

**Последствия:** Невозможность использования embeddings через этот провайдер

---

## Сводная Таблица Проблем

| № | Проблема | Серьёзность | Файл | Строки |
|---|----------|-------------|------|--------|
| 1 | Персистентность токенов | ✅ ИСПРАВЛЕНО | oauth.go | 91-111 |
| 2 | User Management в OAuth flow | ✅ ИСПРАВЛЕНО | oauth.go | - |
| 3 | Отсутствуют 6 полей GenerationConfig | 🔴 КРИТИЧНО | client.go | 136-152 |
| 4 | Нет jitter в retry | 🟢 СРЕДНИЙ | client.go | 405-428 |
| 5 | Разрешающий fallback схем | ✅ ИСПРАВЛЕНО | gemini_genai.go | - |
| 6 | RecordConversationOffered | ⏭️ Н/П | - | - |
| 7 | Отсутствуют методы телеметрии | ⏭️ Н/П | - | - |
| 8 | Отсутствует Experiments API | 🟢 СРЕДНИЙ | client.go | - |
| 9 | Отсутствует EmbedContent | 🟢 НИЗКИЙ | client.go | - |

---

## Приоритет Исправлений

### Исправлено ✅

1. Персистентность токенов
2. User Management в OAuth flow
3. Поля GenerationConfig
4. Jitter в retry
5. Разрешающий fallback схем

### Не применимо ⏭️

- Телеметрия (для IDE-плагинов)

### Средний Приоритет

6. **Реализовать Experiments API** - паритет функций
7. **Реализовать EmbedContent** - для embeddings

---

**Дата:** 2026-01-19  
**Статус:** Требуется исправление критических багов перед использованием
