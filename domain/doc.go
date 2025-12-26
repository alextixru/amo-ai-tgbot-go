/*
Package domain определяет бизнес-сущности приложения.

Доменная область ПОЛНОСТЬЮ делегируется amoCRM SDK:

	Domain Layer (сущности и бизнес-правила):
	  - Модели: github.com/alextixru/amocrm-sdk-go/core/models
	  - Query-объекты: github.com/alextixru/amocrm-sdk-go/core/filters
	  - Агрегаты (batch): github.com/alextixru/amocrm-sdk-go/core/collections

	Application Layer (use cases):
	  - Сервисы: github.com/alextixru/amocrm-sdk-go/core/services

	Infrastructure Layer:
	  - OAuth: github.com/alextixru/amocrm-sdk-go/core/oauth
	  - AI контекст: github.com/alextixru/amocrm-sdk-go/ai

Бот не определяет собственные модели данных, права доступа или режимы работы.
Идентификация пользователя = User.ID из amoCRM.
Права доступа = роли и права из amoCRM Account.

Архитектурное решение:
SDK v0.1.0 содержит полную доменную модель amoCRM,
включая сущности, их связи и бизнес-правила.
Дублирование этих структур в боте не имеет смысла.
*/
package domain
