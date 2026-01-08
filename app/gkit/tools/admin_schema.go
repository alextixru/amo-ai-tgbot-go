package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// AdminSchemaInput входные параметры для инструмента admin_schema
type AdminSchemaInput struct {
	// Layer слой схемы: custom_fields | field_groups | loss_reasons | sources
	Layer string `json:"layer" jsonschema_description:"Слой схемы: custom_fields, field_groups, loss_reasons, sources"`

	// Action действие: search | get | create | update | delete
	Action string `json:"action" jsonschema_description:"Действие: search, get, create, update, delete"`

	// EntityType тип сущности (для custom_fields и field_groups): leads | contacts | companies | customers
	EntityType string `json:"entity_type,omitempty" jsonschema_description:"Тип сущности: leads, contacts, companies, customers (для custom_fields и field_groups)"`

	// ID идентификатор элемента (для get, update, delete)
	ID int `json:"id,omitempty" jsonschema_description:"ID элемента (для custom_fields, loss_reasons, sources)"`

	// GroupID идентификатор группы полей (string в API)
	GroupID string `json:"group_id,omitempty" jsonschema_description:"ID группы полей (для field_groups)"`

	// Filter фильтры для search
	Filter *SchemaFilter `json:"filter,omitempty" jsonschema_description:"Фильтры для поиска"`

	// Data данные для create/update
	Data map[string]any `json:"data,omitempty" jsonschema_description:"Данные для создания/обновления"`
}

// SchemaFilter фильтры для поиска в admin_schema
type SchemaFilter struct {
	Limit int `json:"limit,omitempty" jsonschema_description:"Лимит результатов (по умолчанию 50)"`
	Page  int `json:"page,omitempty" jsonschema_description:"Номер страницы (по умолчанию 1)"`
}

// registerAdminSchemaTool регистрирует инструмент для управления структурой данных
func (r *Registry) registerAdminSchemaTool() {
	r.addTool(genkit.DefineTool[AdminSchemaInput, any](
		r.g,
		"admin_schema",
		"Управление структурой данных amoCRM: кастомные поля, группы полей, причины отказа, источники. "+
			"Layers: custom_fields (поля сущностей), field_groups (группы полей), loss_reasons (причины отказа для сделок), sources (источники лидов). "+
			"Actions: search (список), get (по ID), create, update, delete. "+
			"Для custom_fields и field_groups требуется entity_type (leads/contacts/companies/customers).",
		func(ctx *ai.ToolContext, input AdminSchemaInput) (any, error) {
			return r.handleAdminSchema(ctx.Context, input)
		},
	))
}

func (r *Registry) handleAdminSchema(ctx context.Context, input AdminSchemaInput) (any, error) {
	switch input.Layer {
	case "custom_fields":
		return r.handleCustomFields(ctx, input)
	case "field_groups":
		return r.handleFieldGroups(ctx, input)
	case "loss_reasons":
		return r.handleLossReasons(ctx, input)
	case "sources":
		return r.handleSources(ctx, input)
	default:
		return nil, fmt.Errorf("unknown layer: %s (expected: custom_fields, field_groups, loss_reasons, sources)", input.Layer)
	}
}

// ============================================================================
// Custom Fields
// ============================================================================

func (r *Registry) handleCustomFields(ctx context.Context, input AdminSchemaInput) (any, error) {
	if input.EntityType == "" {
		return nil, fmt.Errorf("entity_type is required for custom_fields (leads, contacts, companies, customers)")
	}

	switch input.Action {
	case "search":
		return r.searchCustomFields(ctx, input.EntityType, input.Filter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.sdk.CustomFields().GetOne(ctx, input.EntityType, input.ID, nil)
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		return r.createCustomField(ctx, input.EntityType, input.Data)
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		return r.updateCustomField(ctx, input.EntityType, input.ID, input.Data)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'delete'")
		}
		err := r.sdk.CustomFields().Delete(ctx, input.EntityType, input.ID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true, "deleted_id": input.ID}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s for custom_fields", input.Action)
	}
}

func (r *Registry) searchCustomFields(ctx context.Context, entityType string, filter *SchemaFilter) ([]*models.CustomField, error) {
	f := filters.NewCustomFieldsFilter()
	f.SetLimit(50)
	f.SetPage(1)
	if filter != nil {
		if filter.Limit > 0 {
			f.SetLimit(filter.Limit)
		}
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		}
	}
	fields, _, err := r.sdk.CustomFields().Get(ctx, entityType, f)
	return fields, err
}

func (r *Registry) createCustomField(ctx context.Context, entityType string, data map[string]any) ([]*models.CustomField, error) {
	field := &models.CustomField{}

	if name, ok := data["name"].(string); ok {
		field.Name = name
	}
	if fieldType, ok := data["type"].(string); ok {
		field.Type = models.CustomFieldType(fieldType)
	}
	if code, ok := data["code"].(string); ok {
		field.Code = code
	}
	if sort, ok := data["sort"].(float64); ok {
		field.Sort = int(sort)
	}
	if groupID, ok := data["group_id"].(string); ok {
		field.GroupID = groupID
	}
	if isRequired, ok := data["is_required"].(bool); ok {
		field.IsRequired = isRequired
	}

	fields, _, err := r.sdk.CustomFields().Create(ctx, entityType, []*models.CustomField{field})
	return fields, err
}

func (r *Registry) updateCustomField(ctx context.Context, entityType string, id int, data map[string]any) ([]*models.CustomField, error) {
	field := &models.CustomField{ID: id}

	if name, ok := data["name"].(string); ok {
		field.Name = name
	}
	if sort, ok := data["sort"].(float64); ok {
		field.Sort = int(sort)
	}
	if groupID, ok := data["group_id"].(string); ok {
		field.GroupID = groupID
	}
	if isRequired, ok := data["is_required"].(bool); ok {
		field.IsRequired = isRequired
	}

	fields, _, err := r.sdk.CustomFields().Update(ctx, entityType, []*models.CustomField{field})
	return fields, err
}

// ============================================================================
// Field Groups
// ============================================================================

func (r *Registry) handleFieldGroups(ctx context.Context, input AdminSchemaInput) (any, error) {
	if input.EntityType == "" {
		return nil, fmt.Errorf("entity_type is required for field_groups (leads, contacts, companies, customers)")
	}

	svc := r.sdk.CustomFieldGroups(input.EntityType)

	switch input.Action {
	case "search":
		return r.searchFieldGroups(ctx, svc, input.Filter)
	case "get":
		// SDK баг: GetOne ожидает int ID, но CustomFieldGroup.ID — string
		// TODO: Fix SDK BaseEntityTypeService to support string IDs
		return nil, fmt.Errorf("field_groups.get is not supported (SDK incompatibility: string ID vs int ID)")
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		return r.createFieldGroup(ctx, svc, input.Data)
	case "update":
		if input.GroupID == "" {
			return nil, fmt.Errorf("group_id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		return r.updateFieldGroup(ctx, svc, input.GroupID, input.Data)
	case "delete":
		if input.GroupID == "" {
			return nil, fmt.Errorf("group_id is required for action 'delete'")
		}
		err := svc.Delete(ctx, input.GroupID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true, "deleted_group_id": input.GroupID}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s for field_groups", input.Action)
	}
}

func (r *Registry) searchFieldGroups(ctx context.Context, svc *services.CustomFieldGroupsService, filter *SchemaFilter) ([]models.CustomFieldGroup, error) {
	// CustomFieldGroupsService наследует BaseEntityTypeService — Get принимает url.Values
	groups, _, err := svc.Get(ctx, nil)
	return groups, err
}

func (r *Registry) createFieldGroup(ctx context.Context, svc *services.CustomFieldGroupsService, data map[string]any) ([]models.CustomFieldGroup, error) {
	group := models.CustomFieldGroup{}

	if name, ok := data["name"].(string); ok {
		group.Name = name
	}
	if sort, ok := data["sort"].(float64); ok {
		group.Sort = int(sort)
	}

	groups, _, err := svc.Create(ctx, []models.CustomFieldGroup{group})
	return groups, err
}

func (r *Registry) updateFieldGroup(ctx context.Context, svc *services.CustomFieldGroupsService, groupID string, data map[string]any) ([]models.CustomFieldGroup, error) {
	group := models.CustomFieldGroup{ID: groupID}

	if name, ok := data["name"].(string); ok {
		group.Name = name
	}
	if sort, ok := data["sort"].(float64); ok {
		group.Sort = int(sort)
	}

	groups, _, err := svc.Update(ctx, []models.CustomFieldGroup{group})
	return groups, err
}

// ============================================================================
// Loss Reasons
// ============================================================================

func (r *Registry) handleLossReasons(ctx context.Context, input AdminSchemaInput) (any, error) {
	switch input.Action {
	case "search":
		return r.searchLossReasons(ctx, input.Filter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.sdk.LossReasons().GetOne(ctx, input.ID)
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		return r.createLossReason(ctx, input.Data)
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		return r.updateLossReason(ctx, input.ID, input.Data)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'delete'")
		}
		err := r.sdk.LossReasons().DeleteOne(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true, "deleted_id": input.ID}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s for loss_reasons", input.Action)
	}
}

func (r *Registry) searchLossReasons(ctx context.Context, filter *SchemaFilter) ([]*models.LossReason, error) {
	// LossReasonsService наследует BaseEntityService — Get принимает url.Values
	reasons, _, err := r.sdk.LossReasons().Get(ctx, nil)
	return reasons, err
}

func (r *Registry) createLossReason(ctx context.Context, data map[string]any) ([]*models.LossReason, error) {
	reason := &models.LossReason{}

	if name, ok := data["name"].(string); ok {
		reason.Name = name
	}

	reasons, _, err := r.sdk.LossReasons().Create(ctx, []*models.LossReason{reason})
	return reasons, err
}

func (r *Registry) updateLossReason(ctx context.Context, id int, data map[string]any) ([]*models.LossReason, error) {
	// LossReasons Update не поддерживается в SDK (возвращает ErrNotAvailableForAction)
	return nil, fmt.Errorf("update not supported for loss_reasons")
}

// ============================================================================
// Sources
// ============================================================================

func (r *Registry) handleSources(ctx context.Context, input AdminSchemaInput) (any, error) {
	switch input.Action {
	case "search":
		return r.searchSources(ctx, input.Filter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.sdk.Sources().GetOne(ctx, input.ID)
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		return r.createSource(ctx, input.Data)
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		return r.updateSource(ctx, input.ID, input.Data)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'delete'")
		}
		err := r.sdk.Sources().Delete(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"success": true, "deleted_id": input.ID}, nil
	default:
		return nil, fmt.Errorf("unknown action: %s for sources", input.Action)
	}
}

func (r *Registry) searchSources(ctx context.Context, filter *SchemaFilter) ([]*models.Source, error) {
	// SourcesFilter не поддерживает pagination
	sources, _, err := r.sdk.Sources().Get(ctx, nil)
	return sources, err
}

func (r *Registry) createSource(ctx context.Context, data map[string]any) ([]*models.Source, error) {
	source := &models.Source{}

	if name, ok := data["name"].(string); ok {
		source.Name = name
	}
	if externalID, ok := data["external_id"].(string); ok {
		source.ExternalID = externalID
	}
	if pipelineID, ok := data["pipeline_id"].(float64); ok {
		source.PipelineID = int(pipelineID)
	}

	sources, _, err := r.sdk.Sources().Create(ctx, []*models.Source{source})
	return sources, err
}

func (r *Registry) updateSource(ctx context.Context, id int, data map[string]any) ([]*models.Source, error) {
	source := &models.Source{ID: id}

	if name, ok := data["name"].(string); ok {
		source.Name = name
	}
	if externalID, ok := data["external_id"].(string); ok {
		source.ExternalID = externalID
	}

	sources, _, err := r.sdk.Sources().Update(ctx, []*models.Source{source})
	return sources, err
}
