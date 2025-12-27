package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// LeadToolInput входные параметры для инструмента crm_leads
type LeadToolInput struct {
	// Action действие: search, get, create, update
	Action string `json:"action"`

	// Filter фильтр для поиска (используется при action=search)
	Filter *LeadFilterInput `json:"filter,omitempty"`

	// ID идентификатор сделки (используется при action=get, update)
	ID int `json:"id,omitempty"`

	// Data данные для создания или обновления (используется при action=create, update)
	Data *LeadDataInput `json:"data,omitempty"`
}

// ... (skipping structs for brevity if possible, but replace_file_content needs contiguous block)
// I will just replace the import and the function call in two chunks or one if close.
// Imports are at top, function is further down.
// Let's do imports first.

// LeadFilterInput параметры фильтрации сделок
type LeadFilterInput struct {
	// Query поисковый запрос (2-3 символа минимум)
	Query string `json:"query,omitempty"`
	// PipelineID ID воронки
	PipelineID int `json:"pipeline_id,omitempty"`
	// StatusID ID статуса
	StatusID int `json:"status_id,omitempty"`
	// ResponsibleUserID ID ответственного
	ResponsibleUserID int `json:"responsible_user_id,omitempty"`
	// Limit количество элементов (по умолчанию 50)
	Limit int `json:"limit,omitempty"`
}

// LeadDataInput данные сделки
type LeadDataInput struct {
	// Title название сделки (обязательно при создании)
	Title string `json:"title,omitempty"`
	// Price бюджет сделки
	Price int `json:"price,omitempty"`
	// PipelineID ID воронки
	PipelineID int `json:"pipeline_id,omitempty"`
	// StatusID ID статуса
	StatusID int `json:"status_id,omitempty"`
	// ResponsibleUserID ID ответственного
	ResponsibleUserID int `json:"responsible_user_id,omitempty"`
}

// registerLeadsTools регистрирует инструменты для работы со сделками
func (r *Registry) registerLeadsTools() {
	r.addTool(genkit.DefineTool[LeadToolInput, any](
		r.g,
		"crm_leads",
		"Управление сделками (Leads). Позволяет искать, получать, создавать и обновлять сделки в amoCRM.",
		func(ctx *ai.ToolContext, input LeadToolInput) (any, error) {
			switch input.Action {
			case "search":
				return r.searchLeads(ctx.Context, input.Filter)
			case "get":
				if input.ID == 0 {
					return nil, fmt.Errorf("id is required for 'get' action")
				}
				return r.getLead(ctx.Context, input.ID)
			case "create":
				if input.Data == nil {
					return nil, fmt.Errorf("data is required for 'create' action")
				}
				return r.createLead(ctx.Context, input.Data)
			case "update":
				if input.ID == 0 {
					return nil, fmt.Errorf("id is required for 'update' action")
				}
				if input.Data == nil {
					return nil, fmt.Errorf("data is required for 'update' action")
				}
				return r.updateLead(ctx.Context, input.ID, input.Data)
			default:
				return nil, fmt.Errorf("unknown action: %s", input.Action)
			}
		},
	))
}

// searchLeads выполняет поиск сделок
func (r *Registry) searchLeads(ctx context.Context, filter *LeadFilterInput) ([]models.Lead, error) {
	sdkFilter := &services.LeadsFilter{
		Limit: 50,
	}

	if filter != nil {
		if filter.Limit > 0 {
			sdkFilter.Limit = filter.Limit
		}
		if filter.Query != "" {
			sdkFilter.Query = filter.Query
		}
		if filter.PipelineID > 0 {
			sdkFilter.FilterByPipelineID = []int{filter.PipelineID}
		}
		if filter.StatusID > 0 {
			sdkFilter.FilterByStatusID = []int{filter.StatusID}
		}
		if filter.ResponsibleUserID > 0 {
			sdkFilter.FilterByResponsibleUserID = []int{filter.ResponsibleUserID}
		}
	}

	return r.sdk.Leads().Get(ctx, sdkFilter)
}

// getLead получает сделку по ID
func (r *Registry) getLead(ctx context.Context, id int) (*models.Lead, error) {
	// Добавляем with=contacts, чтобы сразу видеть контакты
	return r.sdk.Leads().GetOne(ctx, id, []string{"contacts"})
}

// createLead создает новую сделку
func (r *Registry) createLead(ctx context.Context, data *LeadDataInput) (*models.Lead, error) {
	lead := models.Lead{
		Name:       data.Title,
		Price:      data.Price,
		StatusID:   data.StatusID, // При создании можно сразу задать статус
		PipelineID: data.PipelineID,
	}

	if data.ResponsibleUserID > 0 {
		lead.ResponsibleUserID = data.ResponsibleUserID
	}

	// SDK Create принимает слайс и возвращает слайс
	createdLeads, err := r.sdk.Leads().Create(ctx, []models.Lead{lead})
	if err != nil {
		return nil, err
	}

	if len(createdLeads) == 0 {
		return nil, fmt.Errorf("failed to create lead: empty response")
	}

	return &createdLeads[0], nil
}

// updateLead обновляет существующую сделку
func (r *Registry) updateLead(ctx context.Context, id int, data *LeadDataInput) (*models.Lead, error) {
	lead := models.Lead{
		BaseModel: models.BaseModel{ID: id},
	}

	// Заполняем только те поля, которые были переданы (не пустые)
	// Для string проверяем на пустоту, для int на 0.
	// В реальном сценарии для сброса значения в 0 нужны были бы указатели в InputData.
	// Но для упрощения считаем, что 0 не обновляет поле.

	if data.Title != "" {
		lead.Name = data.Title
	}
	if data.Price != 0 {
		lead.Price = data.Price
	}
	if data.StatusID != 0 {
		lead.StatusID = data.StatusID
	}
	if data.PipelineID != 0 {
		lead.PipelineID = data.PipelineID
	}
	if data.ResponsibleUserID != 0 {
		lead.ResponsibleUserID = data.ResponsibleUserID
	}

	updatedLeads, err := r.sdk.Leads().Update(ctx, []models.Lead{lead})
	if err != nil {
		return nil, err
	}

	if len(updatedLeads) == 0 {
		return nil, fmt.Errorf("failed to update lead: empty response")
	}

	return &updatedLeads[0], nil
}
