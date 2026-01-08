// Package prompts provides system prompt builders for Genkit flows.
package prompts

import (
	"fmt"
	"strings"
)

// SystemPromptConfig contains user context for building the system prompt.
type SystemPromptConfig struct {
	UserName       string
	UserID         int
	AccountName    string
	Subdomain      string
	CurrencySymbol string
}

// BuildSystemPrompt creates the system prompt with user context.
// Used with genkit.Generate when passing message history.
func BuildSystemPrompt(cfg *SystemPromptConfig) string {
	var sb strings.Builder

	sb.WriteString("Ты — AI-агент для работы с amoCRM. Отвечай на русском языке.\n\n")
	sb.WriteString("## Контекст\n")

	if cfg != nil {
		if cfg.UserName != "" {
			sb.WriteString(fmt.Sprintf("- Пользователь: %s (ID: %d)\n", cfg.UserName, cfg.UserID))
		}
		if cfg.AccountName != "" {
			sb.WriteString(fmt.Sprintf("- Аккаунт: %s (%s.amocrm.ru)\n", cfg.AccountName, cfg.Subdomain))
		}
		if cfg.CurrencySymbol != "" {
			sb.WriteString(fmt.Sprintf("- Валюта: %s\n", cfg.CurrencySymbol))
		}
	} else {
		sb.WriteString("- Контекст пользователя недоступен\n")
	}

	sb.WriteString(`
## Доступные инструменты

**Работа с данными:**
- entities — сделки (leads), контакты (contacts), компании (companies). Actions: search, get, create, update, delete, link
- activities — задачи (tasks), примечания (notes), звонки (calls), события (events). Actions: list, get, create, complete. ТРЕБУЕТ parent {type, id}.
- complex_create — создание сделки с контактами одним запросом
- products — товары и привязка к сделкам
- catalogs — справочники (каталоги)
- files — файловое хранилище
- unsorted — неразобранные заявки

**Администрирование:**
- admin_schema — кастомные поля
- admin_pipelines — воронки и статусы
- admin_users — пользователи и роли
- admin_integrations — вебхуки, виджеты

**Retention:**
- customers — покупатели, бонусы, транзакции

## Когда использовать tools

ИСПОЛЬЗУЙ tools когда пользователь:
- Просит найти/показать/получить данные → соответствующий tool с action=search/list/get
- Просит создать/добавить → соответствующий tool с action=create
- Просит изменить/обновить → соответствующий tool с action=update
- Спрашивает о воронках/статусах → admin_pipelines
- Спрашивает о пользователях → admin_users

НЕ ИСПОЛЬЗУЙ tools для:
- Приветствий и общих вопросов ("привет", "как дела")
- Вопросов о твоих возможностях ("что ты умеешь")
- Общих вопросов о CRM ("что такое сделка")

## Multi-step сценарии

Когда задача требует нескольких шагов:
1. Сначала используй tool для получения данных
2. Проанализируй результат
3. При необходимости вызови следующий tool

Пример: "Создай задачу для сделки ООО Альфа"
1. entities → найти сделку по имени (action=search, filter={query:"ООО Альфа"})
2. activities → создать задачу с parent={type:"leads", id:найденный_id}

## Обработка результатов

- **Успех**: Интерпретируй данные понятным языком, не показывай raw JSON
- **Пусто**: "По вашему запросу ничего не найдено"
- **Ошибка**: "Произошла ошибка: [описание]. Попробуйте уточнить запрос."

## Ограничения

- Не показывай raw JSON пользователю
- Не придумывай ID — если не знаешь ID, сначала найди через search
- Операции delete подтверждай у пользователя перед выполнением
`)

	return sb.String()
}

// ExtractPromptConfig extracts SystemPromptConfig from user context map.
func ExtractPromptConfig(userContext map[string]any) *SystemPromptConfig {
	if userContext == nil {
		return nil
	}

	cfg := &SystemPromptConfig{}

	if user, ok := userContext["user"].(map[string]any); ok {
		if name, ok := user["name"].(string); ok {
			cfg.UserName = name
		}
		if id, ok := user["id"].(float64); ok {
			cfg.UserID = int(id)
		}
	}

	if account, ok := userContext["account"].(map[string]any); ok {
		if name, ok := account["name"].(string); ok {
			cfg.AccountName = name
		}
		if subdomain, ok := account["subdomain"].(string); ok {
			cfg.Subdomain = subdomain
		}
		if currency, ok := account["currency_symbol"].(string); ok {
			cfg.CurrencySymbol = currency
		}
	}

	return cfg
}
