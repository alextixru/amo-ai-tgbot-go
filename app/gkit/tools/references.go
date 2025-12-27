package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// ReferenceToolInput входные параметры для инструмента crm_get_reference
type ReferenceToolInput struct {
	// Action действие: users, pipelines, tags, roles, custom_fields, events
	Action string `json:"action"`

	// EntityType тип сущности (для tags, custom_fields): leads, contacts, companies, customers
	EntityType string `json:"entity_type,omitempty"`

	// ID идентификатор (если нужно получить конкретный элемент)
	ID int `json:"id,omitempty"`

	// Limit количество элементов (по умолчанию 50)
	Limit int `json:"limit,omitempty"`
}

// registerReferencesTool регистрирует инструмент для получения справочной информации
func (r *Registry) registerReferencesTool() {
	r.addTool(genkit.DefineTool[ReferenceToolInput, any](
		r.g,
		"crm_get_reference",
		"Получение справочной информации (пользователи, воронки, статусы, теги, поля). Используйте этот инструмент для поиска ID статусов, пользователей и другой мета-информации.",
		func(ctx *ai.ToolContext, input ReferenceToolInput) (any, error) {
			switch input.Action {
			case "users":
				return r.getUsers(ctx.Context, input)
			case "pipelines":
				return r.getPipelines(ctx.Context, input)
			case "tags":
				return r.getTags(ctx.Context, input)
			case "roles":
				return r.getRoles(ctx.Context, input)
			case "custom_fields":
				return r.getCustomFields(ctx.Context, input)
			case "events":
				return r.getEvents(ctx.Context, input)
			default:
				return nil, fmt.Errorf("unknown action: %s", input.Action)
			}
		},
	))
}

func (r *Registry) getUsers(ctx context.Context, input ReferenceToolInput) (any, error) {
	if input.ID > 0 {
		return r.sdk.Users().GetOne(ctx, input.ID)
	}

	filter := &services.UsersFilter{
		Limit: 50,
	}
	if input.Limit > 0 {
		filter.Limit = input.Limit
	}

	return r.sdk.Users().Get(ctx, filter)
}

func (r *Registry) getPipelines(ctx context.Context, input ReferenceToolInput) (any, error) {
	if input.ID > 0 {
		return r.sdk.Pipelines().GetOne(ctx, input.ID)
	}
	// Pipelines.Get не принимает фильтр в текущей версии SDK, возвращает всё
	return r.sdk.Pipelines().Get(ctx)
}

func (r *Registry) getTags(ctx context.Context, input ReferenceToolInput) (any, error) {
	if input.EntityType == "" {
		return nil, fmt.Errorf("entity_type is required for tags reference (leads, contacts, companies, customers)")
	}

	filter := &services.TagsFilter{}
	// В SDK TagsService.Get не принимает Limit в фильтре (на данный момент), но проверим если добавим
	// Пока игнорируем Limit, так как в TagsFilter его может не быть (надо проверить)
	// Проверил toolmap: Get(ctx, entityType string, filter *TagsFilter)

	return r.sdk.Tags().Get(ctx, input.EntityType, filter)
}

func (r *Registry) getRoles(ctx context.Context, input ReferenceToolInput) (any, error) {
	if input.ID > 0 {
		return r.sdk.Roles().GetOne(ctx, input.ID, nil)
	}

	filter := &services.RoleFilter{ // В SDK называется RoleFilter
		Limit: 50,
	}
	if input.Limit > 0 {
		filter.Limit = input.Limit
	}

	return r.sdk.Roles().Get(ctx, filter)
}

func (r *Registry) getCustomFields(ctx context.Context, input ReferenceToolInput) (any, error) {
	if input.EntityType == "" {
		return nil, fmt.Errorf("entity_type is required for fields reference (leads, contacts, companies, customers)")
	}

	if input.ID > 0 {
		return r.sdk.CustomFields().GetOne(ctx, input.EntityType, input.ID)
	}

	filter := &services.CustomFieldsFilter{
		Limit: 50,
	}
	if input.Limit > 0 {
		filter.Limit = input.Limit
	}

	return r.sdk.CustomFields().Get(ctx, input.EntityType, filter)
}

func (r *Registry) getEvents(ctx context.Context, input ReferenceToolInput) (any, error) {
	if input.ID > 0 {
		return r.sdk.Events().GetOne(ctx, input.ID)
	}

	filter := &services.EventsFilter{
		Limit: 50,
	}
	if input.Limit > 0 {
		filter.Limit = input.Limit
	}

	return r.sdk.Events().Get(ctx, filter)
}
