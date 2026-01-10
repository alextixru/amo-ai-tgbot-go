package tools

import (
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func (r *Registry) RegisterComplexCreateTool() {
	r.addTool(genkit.DefineTool[gkitmodels.ComplexCreateInput, any](
		r.g,
		"complex_create",
		"Create leads with contacts and companies in one request",
		func(ctx *ai.ToolContext, input gkitmodels.ComplexCreateInput) (any, error) {
			// Маппинг входных данных в модель SDK
			lead := &amomodels.Lead{
				Name:       input.Lead.Name,
				Price:      input.Lead.Price,
				PipelineID: input.Lead.PipelineID,
				StatusID:   input.Lead.StatusID,
			}
			lead.ResponsibleUserID = input.Lead.ResponsibleUserID

			// Инициализируем вложенные данные
			lead.Embedded = &amomodels.LeadEmbedded{}

			// Добавляем контакты
			if len(input.Contacts) > 0 {
				lead.Embedded.Contacts = make([]*amomodels.Contact, 0, len(input.Contacts))
				for _, c := range input.Contacts {
					contact := &amomodels.Contact{
						Name: c.Name,
					}
					// Для простоты пока маппим только имя.
					lead.Embedded.Contacts = append(lead.Embedded.Contacts, contact)
				}
			}

			// Добавляем компанию
			if input.Company != nil {
				lead.Embedded.Companies = []*amomodels.Company{
					{
						Name: input.Company.Name,
					},
				}
			}

			return r.complexCreateService.CreateComplex(ctx, lead)
		},
	))
}
