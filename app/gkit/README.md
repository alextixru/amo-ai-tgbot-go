# Genkit AI Agent

AI агент для управления amoCRM через естественный язык.

## Архитектура

```
User Message → Router Flow → Specialized Flow → CRM Tools → amoCRM SDK
```

### Router Flow
Классифицирует intent пользователя и направляет в нужный flow:
- `leads` — работа со сделками
- `contacts` — работа с контактами  
- `tasks` — работа с задачами
- `general` — общие вопросы (без CRM операций)

### Specialized Flows
Каждый flow получает свой набор tools:

| Flow | Tools |
|------|-------|
| `leads_flow` | getLeads, createLead, updateLead |
| `contacts_flow` | getContacts, createContact |
| `tasks_flow` | getTasks, createTask |

## Структура

```
gkit/
├── agent.go          # Точка входа, оркестратор flows
├── flows/
│   ├── router.go     # Router Flow (классификация intent)
│   ├── chat.go       # Chat Flow (текущий)
│   ├── leads.go      # Leads Flow + tools
│   └── contacts.go   # Contacts Flow + tools
├── prompts/
│   ├── router.prompt # Промпт для роутера
│   └── chat.prompt   # Общий чат промпт
├── tools.go          # Общие утилиты для tools
└── types.go          # Типы данных
```

## План реализации

1. [x] Базовый Chat Flow (работает)
2. [ ] Router Flow — классификация intent
3. [ ] Подключение Router → Chat Flow
4. [ ] Leads Flow + tools
5. [ ] Contacts Flow + tools

## Принципы

- **Минимум tools на flow** — LLM видит только релевантные инструменты
- **Изоляция** — каждый flow тестируется независимо
- **Наблюдаемость** — все flows видны в Genkit UI
