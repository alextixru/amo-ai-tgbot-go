// Package models — корневой пакет для структур данных Genkit.
//
// Структура:
//   - tools/ — Input DTOs для SDK-инструментов (полные схемы)
//   - flows/ — Input DTOs для Flow (упрощённые схемы для Main Agent)
//
// Принципы:
//   - Input структуры содержат jsonschema_description для LLM
//   - Output — модели из amocrm-sdk-go напрямую
package models
