package tools

import (
	"encoding/json"
	"fmt"
	"net/url"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func (r *Registry) RegisterAdminSchemaTool() {
	r.addTool(genkit.DefineTool[gkitmodels.AdminSchemaInput, any](
		r.g,
		"admin_schema",
		"Work with CRM schema (custom fields, field groups, loss reasons and sources)",
		func(ctx *ai.ToolContext, input gkitmodels.AdminSchemaInput) (any, error) {
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
				return nil, fmt.Errorf("unknown layer: %s", input.Layer)
			}
		},
	))
}

func buildCustomFieldsFilter(f *gkitmodels.SchemaFilter) *filters.CustomFieldsFilter {
	if f == nil {
		return nil
	}
	filter := &filters.CustomFieldsFilter{}
	if f.Limit > 0 {
		filter.Limit = f.Limit
	}
	if f.Page > 0 {
		filter.Page = f.Page
	}
	if len(f.IDs) > 0 {
		filter.IDs = f.IDs
	}
	if len(f.Types) > 0 {
		filter.Types = f.Types
	}
	return filter
}

func buildSourcesFilter(f *gkitmodels.SchemaFilter) *filters.SourcesFilter {
	if f == nil {
		return nil
	}
	filter := &filters.SourcesFilter{}
	if len(f.ExternalIDs) > 0 {
		filter.ExternalIDs = f.ExternalIDs
	}
	return filter
}

func buildBaseFilter(f *gkitmodels.SchemaFilter) url.Values {
	if f == nil {
		return url.Values{}
	}
	v := url.Values{}
	if f.Limit > 0 {
		v.Add("limit", fmt.Sprintf("%d", f.Limit))
	}
	if f.Page > 0 {
		v.Add("page", fmt.Sprintf("%d", f.Page))
	}
	return v
}

func (r *Registry) handleCustomFields(ctx *ai.ToolContext, input gkitmodels.AdminSchemaInput) (any, error) {
	if input.EntityType == "" {
		return nil, fmt.Errorf("entity_type is required for custom_fields")
	}
	switch input.Action {
	case "search", "list":
		filter := buildCustomFieldsFilter(input.Filter)
		return r.adminSchemaService.ListCustomFields(ctx, input.EntityType, filter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for get")
		}
		return r.adminSchemaService.GetCustomField(ctx, input.EntityType, input.ID)
	case "create":
		var fields []*amomodels.CustomField
		data, _ := json.Marshal(input.Data["fields"])
		if err := json.Unmarshal(data, &fields); err != nil {
			return nil, fmt.Errorf("failed to parse fields: %w", err)
		}
		return r.adminSchemaService.CreateCustomFields(ctx, input.EntityType, fields)
	case "update":
		var fields []*amomodels.CustomField
		data, _ := json.Marshal(input.Data["fields"])
		if err := json.Unmarshal(data, &fields); err != nil {
			return nil, fmt.Errorf("failed to parse fields: %w", err)
		}
		return r.adminSchemaService.UpdateCustomFields(ctx, input.EntityType, fields)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for delete")
		}
		return nil, r.adminSchemaService.DeleteCustomField(ctx, input.EntityType, input.ID)
	default:
		return nil, fmt.Errorf("unknown action for custom_fields: %s", input.Action)
	}
}

func (r *Registry) handleFieldGroups(ctx *ai.ToolContext, input gkitmodels.AdminSchemaInput) (any, error) {
	if input.EntityType == "" {
		return nil, fmt.Errorf("entity_type is required for field_groups")
	}
	switch input.Action {
	case "search", "list":
		filter := buildBaseFilter(input.Filter)
		return r.adminSchemaService.ListFieldGroups(ctx, input.EntityType, filter)
	case "get":
		if input.GroupID == "" {
			return nil, fmt.Errorf("group_id is required for get")
		}
		return r.adminSchemaService.GetFieldGroup(ctx, input.EntityType, input.GroupID)
	case "create":
		var groups []amomodels.CustomFieldGroup
		data, _ := json.Marshal(input.Data["groups"])
		if err := json.Unmarshal(data, &groups); err != nil {
			return nil, fmt.Errorf("failed to parse groups: %w", err)
		}
		return r.adminSchemaService.CreateFieldGroups(ctx, input.EntityType, groups)
	case "update":
		var groups []amomodels.CustomFieldGroup
		data, _ := json.Marshal(input.Data["groups"])
		if err := json.Unmarshal(data, &groups); err != nil {
			return nil, fmt.Errorf("failed to parse groups: %w", err)
		}
		return r.adminSchemaService.UpdateFieldGroups(ctx, input.EntityType, groups)
	case "delete":
		if input.GroupID == "" {
			return nil, fmt.Errorf("group_id is required for delete")
		}
		return nil, r.adminSchemaService.DeleteFieldGroup(ctx, input.EntityType, input.GroupID)
	default:
		return nil, fmt.Errorf("unknown action for field_groups: %s", input.Action)
	}
}

func (r *Registry) handleLossReasons(ctx *ai.ToolContext, input gkitmodels.AdminSchemaInput) (any, error) {
	switch input.Action {
	case "search", "list":
		filter := buildBaseFilter(input.Filter)
		return r.adminSchemaService.ListLossReasons(ctx, filter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for get")
		}
		return r.adminSchemaService.GetLossReason(ctx, input.ID)
	case "create":
		var reasons []*amomodels.LossReason
		data, _ := json.Marshal(input.Data["reasons"])
		if err := json.Unmarshal(data, &reasons); err != nil {
			return nil, fmt.Errorf("failed to parse reasons: %w", err)
		}
		return r.adminSchemaService.CreateLossReasons(ctx, reasons)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for delete")
		}
		return nil, r.adminSchemaService.DeleteLossReason(ctx, input.ID)
	default:
		return nil, fmt.Errorf("unknown action for loss_reasons: %s", input.Action)
	}
}

func (r *Registry) handleSources(ctx *ai.ToolContext, input gkitmodels.AdminSchemaInput) (any, error) {
	switch input.Action {
	case "search", "list":
		filter := buildSourcesFilter(input.Filter)
		return r.adminSchemaService.ListSources(ctx, filter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for get")
		}
		return r.adminSchemaService.GetSource(ctx, input.ID)
	case "create":
		var sources []*amomodels.Source
		data, _ := json.Marshal(input.Data["sources"])
		if err := json.Unmarshal(data, &sources); err != nil {
			return nil, fmt.Errorf("failed to parse sources: %w", err)
		}
		return r.adminSchemaService.CreateSources(ctx, sources)
	case "update":
		var sources []*amomodels.Source
		data, _ := json.Marshal(input.Data["sources"])
		if err := json.Unmarshal(data, &sources); err != nil {
			return nil, fmt.Errorf("failed to parse sources: %w", err)
		}
		return r.adminSchemaService.UpdateSources(ctx, sources)
	case "delete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for delete")
		}
		return nil, r.adminSchemaService.DeleteSource(ctx, input.ID)
	default:
		return nil, fmt.Errorf("unknown action for sources: %s", input.Action)
	}
}
