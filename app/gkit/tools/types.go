// Package tools provides Genkit tool definitions for amoCRM SDK operations.
package tools

import (
	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
)

// === SDK Models Re-export (Output types) ===
// Tools use these aliases instead of importing SDK directly.
// If SDK models change, tools automatically get updated types.

// Lead is an alias for SDK Lead model.
type Lead = models.Lead

// LeadsFilter is an alias for SDK LeadsFilter.
type LeadsFilter = services.LeadsFilter

// Pipeline is an alias for SDK Pipeline model.
type Pipeline = models.Pipeline

// Status is an alias for SDK Status model.
type Status = models.Status

// === Input types (for LLM tool calling) ===
// These require jsonschema_description tags for Genkit.

// SearchLeadsInput input for searchLeads tool.
type SearchLeadsInput struct {
	Query string `json:"query,omitempty" jsonschema_description:"Search query for leads (by name, phone, email)"`
	Limit int    `json:"limit,omitempty" jsonschema_description:"Maximum number of results (default 10, max 50)"`
}

// GetLeadInput input for getLead tool.
type GetLeadInput struct {
	LeadID int `json:"lead_id" jsonschema_description:"ID of the lead to retrieve"`
}

// CreateLeadInput input for createLead tool.
type CreateLeadInput struct {
	Name       string `json:"name" jsonschema_description:"Lead name (required)"`
	Price      int    `json:"price,omitempty" jsonschema_description:"Lead budget/price"`
	PipelineID int    `json:"pipeline_id,omitempty" jsonschema_description:"Pipeline ID (uses default if not specified)"`
	StatusID   int    `json:"status_id,omitempty" jsonschema_description:"Status ID within pipeline"`
}

// UpdateLeadInput input for updateLead tool.
type UpdateLeadInput struct {
	LeadID   int    `json:"lead_id" jsonschema_description:"ID of the lead to update"`
	Name     string `json:"name,omitempty" jsonschema_description:"New lead name"`
	Price    int    `json:"price,omitempty" jsonschema_description:"New budget/price"`
	StatusID int    `json:"status_id,omitempty" jsonschema_description:"New status ID"`
}
