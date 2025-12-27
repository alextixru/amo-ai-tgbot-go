package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// CompanyToolInput входные параметры для инструмента crm_companies
type CompanyToolInput struct {
	// Action действие: search, get, create, update
	Action string `json:"action"`

	// Filter фильтр для поиска (используется при action=search)
	Filter *CompanyFilterInput `json:"filter,omitempty"`

	// ID идентификатор компании (используется при action=get, update)
	ID int `json:"id,omitempty"`

	// Data данные для создания или обновления (используется при action=create, update)
	Data *CompanyDataInput `json:"data,omitempty"`
}

// CompanyFilterInput параметры фильтрации компаний
type CompanyFilterInput struct {
	// Query поисковый запрос (название, email, телефон)
	Query string `json:"query,omitempty"`
	// ResponsibleUserID ID ответственного
	ResponsibleUserID int `json:"responsible_user_id,omitempty"`
	// Limit количество элементов (по умолчанию 50)
	Limit int `json:"limit,omitempty"`
}

// CompanyDataInput данные компании
type CompanyDataInput struct {
	// Name название компании (обязательно при создании)
	Name string `json:"name,omitempty"`
	// ResponsibleUserID ID ответственного
	ResponsibleUserID int `json:"responsible_user_id,omitempty"`
}

// registerCompaniesTools регистрирует инструменты для работы с компаниями
func (r *Registry) registerCompaniesTools() {
	r.addTool(genkit.DefineTool[CompanyToolInput, any](
		r.g,
		"crm_companies",
		"Управление компаниями (Companies). Позволяет искать, получать, создавать и обновлять компании в amoCRM.",
		func(ctx *ai.ToolContext, input CompanyToolInput) (any, error) {
			switch input.Action {
			case "search":
				return r.searchCompanies(ctx.Context, input.Filter)
			case "get":
				if input.ID == 0 {
					return nil, fmt.Errorf("id is required for 'get' action")
				}
				return r.getCompany(ctx.Context, input.ID)
			case "create":
				if input.Data == nil {
					return nil, fmt.Errorf("data is required for 'create' action")
				}
				return r.createCompany(ctx.Context, input.Data)
			case "update":
				if input.ID == 0 {
					return nil, fmt.Errorf("id is required for 'update' action")
				}
				if input.Data == nil {
					return nil, fmt.Errorf("data is required for 'update' action")
				}
				return r.updateCompany(ctx.Context, input.ID, input.Data)
			default:
				return nil, fmt.Errorf("unknown action: %s", input.Action)
			}
		},
	))
}

// searchCompanies выполняет поиск компаний
func (r *Registry) searchCompanies(ctx context.Context, filter *CompanyFilterInput) ([]models.Company, error) {
	sdkFilter := &services.CompaniesFilter{
		Limit: 50,
	}

	if filter != nil {
		if filter.Limit > 0 {
			sdkFilter.Limit = filter.Limit
		}
		if filter.Query != "" {
			sdkFilter.Query = filter.Query
		}
		if filter.ResponsibleUserID > 0 {
			sdkFilter.FilterByResponsibleUserID = []int{filter.ResponsibleUserID}
		}
	}

	return r.sdk.Companies().Get(ctx, sdkFilter)
}

// getCompany получает компанию по ID
func (r *Registry) getCompany(ctx context.Context, id int) (*models.Company, error) {
	// Добавляем with=contacts,leads, чтобы видеть связи
	return r.sdk.Companies().GetOne(ctx, id, []string{"contacts", "leads"})
}

// createCompany создает новую компанию
func (r *Registry) createCompany(ctx context.Context, data *CompanyDataInput) (*models.Company, error) {
	company := models.Company{
		Name: data.Name,
	}

	if data.ResponsibleUserID > 0 {
		company.ResponsibleUserID = data.ResponsibleUserID
	}

	createdCompanies, err := r.sdk.Companies().Create(ctx, []models.Company{company})
	if err != nil {
		return nil, err
	}

	if len(createdCompanies) == 0 {
		return nil, fmt.Errorf("failed to create company: empty response")
	}

	return &createdCompanies[0], nil
}

// updateCompany обновляет существующую компанию
func (r *Registry) updateCompany(ctx context.Context, id int, data *CompanyDataInput) (*models.Company, error) {
	company := models.Company{
		BaseModel: models.BaseModel{ID: id},
	}

	if data.Name != "" {
		company.Name = data.Name
	}
	if data.ResponsibleUserID != 0 {
		company.ResponsibleUserID = data.ResponsibleUserID
	}

	updatedCompanies, err := r.sdk.Companies().Update(ctx, []models.Company{company})
	if err != nil {
		return nil, err
	}

	if len(updatedCompanies) == 0 {
		return nil, fmt.Errorf("failed to update company: empty response")
	}

	return &updatedCompanies[0], nil
}
