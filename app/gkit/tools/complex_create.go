package tools

import (
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// ComplexCreateInput входные параметры для создания сделки с контактами/компанией
type ComplexCreateInput struct {
	// Lead данные сделки
	Lead LeadData `json:"lead" jsonschema_description:"Данные сделки (обязательно)"`

	// Contacts контакты для привязки
	Contacts []ContactData `json:"contacts,omitempty" jsonschema_description:"Контакты для создания и привязки"`

	// Company компания для привязки
	Company *CompanyData `json:"company,omitempty" jsonschema_description:"Компания для создания и привязки"`
}

// LeadData данные сделки
type LeadData struct {
	Name              string `json:"name" jsonschema_description:"Название сделки"`
	Price             int    `json:"price,omitempty" jsonschema_description:"Бюджет"`
	PipelineID        int    `json:"pipeline_id,omitempty" jsonschema_description:"ID воронки"`
	StatusID          int    `json:"status_id,omitempty" jsonschema_description:"ID статуса"`
	ResponsibleUserID int    `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственного"`
}

// ContactData данные контакта
type ContactData struct {
	Name   string `json:"name" jsonschema_description:"Имя контакта"`
	Phone  string `json:"phone,omitempty" jsonschema_description:"Телефон"`
	Email  string `json:"email,omitempty" jsonschema_description:"Email"`
	IsMain bool   `json:"is_main,omitempty" jsonschema_description:"Основной контакт"`
}

// CompanyData данные компании
type CompanyData struct {
	Name string `json:"name" jsonschema_description:"Название компании"`
}

// registerComplexCreateTool регистрирует инструмент комплексного создания
func (r *Registry) registerComplexCreateTool() {
	r.addTool(genkit.DefineTool[ComplexCreateInput, any](
		r.g,
		"complex_create",
		"Создание сделки с контактами и компанией одним запросом. "+
			"Атомарная операция: создаёт lead + contacts + company и связывает их автоматически.",
		func(ctx *ai.ToolContext, input ComplexCreateInput) (any, error) {
			return r.handleComplexCreate(ctx, input)
		},
	))
}

func (r *Registry) handleComplexCreate(ctx *ai.ToolContext, input ComplexCreateInput) (any, error) {
	if input.Lead.Name == "" {
		return nil, fmt.Errorf("lead.name is required")
	}

	// Создаём базовую структуру Lead
	lead := &models.Lead{
		Name:       input.Lead.Name,
		Price:      input.Lead.Price,
		PipelineID: input.Lead.PipelineID,
		StatusID:   input.Lead.StatusID,
	}
	if input.Lead.ResponsibleUserID > 0 {
		lead.ResponsibleUserID = input.Lead.ResponsibleUserID
	}

	// Инициализируем embedded
	lead.Embedded = &models.LeadEmbedded{}

	// Добавляем контакты
	if len(input.Contacts) > 0 {
		contacts := make([]*models.Contact, 0, len(input.Contacts))
		for _, c := range input.Contacts {
			contact := &models.Contact{
				Name:   c.Name,
				IsMain: c.IsMain,
			}
			// TODO: добавить поддержку custom_fields для phone/email
			// Пока что phone и email можно передать через custom_fields_values
			contacts = append(contacts, contact)
		}
		lead.Embedded.Contacts = contacts
	}

	// Добавляем компанию
	if input.Company != nil && input.Company.Name != "" {
		company := &models.Company{
			Name: input.Company.Name,
		}
		lead.Embedded.Companies = []*models.Company{company}
	}

	// Вызываем AddOneComplex
	return r.sdk.Leads().AddOneComplex(ctx.Context, lead)
}
