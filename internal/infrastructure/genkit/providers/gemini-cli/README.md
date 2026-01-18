# Спецификация интеграции Gemini CLI (Code Assist)

Этот документ описывает реверс-инжиниринг спецификаций **Google Gemini Code Assist API** (Internal), который используется официальным `gemini-cli` в режиме авторизации `AuthType.LOGIN_WITH_GOOGLE`.

> **⚠️ ПРЕДУПРЕЖДЕНИЕ**: Это **Внутренний API (`v1internal`)**. Он не задокументирован для публичного использования. Google может сломать, изменить или закрыть доступ к нему в любой момент. Используйте на свой страх и риск.

## 1. Аутентификация (OAuth 2.0)

Чтобы представиться как официальный CLI, мы должны использовать его публичный Client ID и Secret.

*   **Тип**: OAuth 2.0 (User Redirect / Copy-Paste Code)
*   **Token Endpoint**: `https://oauth2.googleapis.com/token`
*   **Scopes (Права)**:
    *   `https://www.googleapis.com/auth/cloud-platform`
    *   `https://www.googleapis.com/auth/userinfo.email`
    *   `https://www.googleapis.com/auth/userinfo.profile`

### Учетные данные (из `packages/core/src/code_assist/oauth2.ts`)
*   **Client ID**: `681255809395-oo8ft2oprdrnp9e3aqf6av3hmdib135j.apps.googleusercontent.com`
*   **Client Secret**: `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`

> **Примечание**: В комментариях исходного кода сказано: *"It's ok to save this in git because this is an installed application... the client secret is obviously not treated as a secret."*

## 2. API Эндпоинт

*   **Базовый URL**: `https://cloudcode-pa.googleapis.com/v1internal`
*   **Методы**:
    *   `POST /:generateContent` (Обычный запрос)
    *   `POST /:streamGenerateContent?alt=sse` (Стриминг)
    *   `POST /:countTokens`

## 3. Кастомная структура запроса

В отличие от стандартного Vertex AI/Gemini API, этот внутренний API требует специфическую JSON-обертку. Стандартные SDK (Genkit GoogleAI, VertexAI) **не будут работать** "из коробки", так как они отправляют чистый JSON без обертки.

**Требуемая структура JSON:**

```json
{
  "model": "models/gemini-2.0-flash",  // Имя модели
  "project": "Your-Project-ID",       // ID GCP проекта (из OAuth токена или User Info)
  "user_prompt_id": "optional-uuid",
  "request": {
     // --- Сюда вкладывается стандартный Payload от Vertex AI ---
     "contents": [
       {
         "role": "user",
         "parts": [{ "text": "Hello world" }]
       }
     ],
     "generationConfig": {
       "temperature": 0.5,
       "maxOutputTokens": 1024
     },
     "tools": [ ... ]
  }
}
```

### Обработка ответа
*   **Streaming**: Ответ приходит в формате Server-Sent Events (SSE).
*   **Данные**: Именование полей соответствует внутренним protobuf определениям (см. код `gemini-cli` для `CaGenerateContentResponse`).

## 4. Руководство по реализации (Custom Provider)

Чтобы использовать "Killer Features" Code Assist в этом проекте:

1.  **НЕ ИСПОЛЬЗУЙТЕ** напрямую `github.com/firebase/genkit/go/plugins/vertexai`.
2.  **Создайте Custom Provider**, который:
    *   Реализует интерфейс Genkit Model.
    *   Использует библиотеку `golang.org/x/oauth2` с указанными выше ID/Secret.
    *   Вручную формирует JSON-пейлоуд с оберткой `request`.
    *   Отправляет HTTP POST запросы на `cloudcode-pa.googleapis.com`.

### В чем "Киллер-фича"?
Этот API используется плагинами Google для IDE (Cloud Code, IDX). Потенциальные преимущества:
*   Другие лимиты/квоты (привязка к пользователю, а не только к проекту).
*   Доступ к специфическим внутренним моделям или логике, оптимизированной для кодинга.
*   Использование Free Tier для Code Assist, привязанного к аккаунту разработчика, вместо требований биллинга GCP на каждый запрос.
