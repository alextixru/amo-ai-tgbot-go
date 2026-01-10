// Package models содержит входные структуры (Input DTOs) для инструментов Genkit.
//
// Все структуры предназначены для транспортного слоя и содержат теги
// jsonschema_description, которые используются AI для понимания назначения полей.
//
// Принципы:
//   - Input структуры определяются здесь
//   - Output — модели из amocrm-sdk-go напрямую
//   - Общие типы (LinkTarget, ParentEntity) вынесены в common.go
package models
