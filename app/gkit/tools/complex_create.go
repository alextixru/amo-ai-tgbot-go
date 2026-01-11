package tools

import (
	"encoding/json"

	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// mapLeadToSDK конвертирует LeadData в SDK Lead
func (r *Registry) mapLeadToSDK(input gkitmodels.ComplexCreateInput) *amomodels.Lead {
	lead := &amomodels.Lead{
		Name:       input.Lead.Name,
		Price:      input.Lead.Price,
		PipelineID: input.Lead.PipelineID,
		StatusID:   input.Lead.StatusID,
	}
	lead.ResponsibleUserID = input.Lead.ResponsibleUserID

	// Инициализируем вложенные данные
	lead.Embedded = &amomodels.LeadEmbedded{}

	// Маппинг кастомных полей сделки
	if len(input.Lead.CustomFieldsValues) > 0 {
		lead.CustomFieldsValues = mapCustomFieldsValues(input.Lead.CustomFieldsValues)
	}

	// Маппинг тегов
	if len(input.Lead.Tags) > 0 {
		lead.Embedded.Tags = make([]amomodels.Tag, 0, len(input.Lead.Tags))
		for _, tagName := range input.Lead.Tags {
			lead.Embedded.Tags = append(lead.Embedded.Tags, amomodels.Tag{Name: tagName})
		}
	}

	// Добавляем контакты
	if len(input.Contacts) > 0 {
		lead.Embedded.Contacts = make([]*amomodels.Contact, 0, len(input.Contacts))
		for _, c := range input.Contacts {
			contact := &amomodels.Contact{
				Name: c.Name,
			}

			// Маппинг имени/фамилии отдельно
			if c.FirstName != "" {
				contact.FirstName = c.FirstName
			}
			if c.LastName != "" {
				contact.LastName = c.LastName
			}
			if c.ResponsibleUserID != 0 {
				contact.ResponsibleUserID = c.ResponsibleUserID
			}
			if c.IsMain {
				contact.IsMain = c.IsMain
			}

			// Собираем кастомные поля контакта (включая phone/email)
			contactCFV := buildContactCustomFields(c)
			if len(contactCFV) > 0 {
				contact.CustomFieldsValues = contactCFV
			}

			lead.Embedded.Contacts = append(lead.Embedded.Contacts, contact)
		}
	}

	// Добавляем компанию
	if input.Company != nil {
		company := &amomodels.Company{
			Name: input.Company.Name,
		}
		if input.Company.ResponsibleUserID != 0 {
			company.ResponsibleUserID = input.Company.ResponsibleUserID
		}
		if len(input.Company.CustomFieldsValues) > 0 {
			company.CustomFieldsValues = mapCustomFieldsValues(input.Company.CustomFieldsValues)
		}
		lead.Embedded.Companies = []*amomodels.Company{company}
	}

	return lead
}

// buildContactCustomFields собирает кастомные поля контакта из Phone, Email и CustomFieldsValues
func buildContactCustomFields(c gkitmodels.ContactData) []amomodels.CustomFieldValue {
	var result []amomodels.CustomFieldValue

	// Телефон — используем специальный тип поля
	if c.Phone != "" {
		result = append(result, amomodels.CustomFieldValue{
			FieldCode: "PHONE",
			Values: []amomodels.FieldValueElement{
				{Value: c.Phone, EnumCode: "WORK"},
			},
		})
	}

	// Email — используем специальный тип поля
	if c.Email != "" {
		result = append(result, amomodels.CustomFieldValue{
			FieldCode: "EMAIL",
			Values: []amomodels.FieldValueElement{
				{Value: c.Email, EnumCode: "WORK"},
			},
		})
	}

	// Прочие кастомные поля
	if len(c.CustomFieldsValues) > 0 {
		result = append(result, mapCustomFieldsValues(c.CustomFieldsValues)...)
	}

	return result
}

// mapCustomFieldsValues конвертирует map[string]any в []CustomFieldValue
func mapCustomFieldsValues(cfv map[string]any) []amomodels.CustomFieldValue {
	// Пробуем десериализовать через JSON для гибкости
	data, err := json.Marshal(cfv)
	if err != nil {
		return nil
	}

	var result []amomodels.CustomFieldValue
	if err := json.Unmarshal(data, &result); err != nil {
		// Fallback: пробуем как map field_id -> values
		for fieldID, values := range cfv {
			cfValue := amomodels.CustomFieldValue{
				FieldCode: fieldID,
			}
			// Пытаемся распарсить values как массив
			if valArr, ok := values.([]any); ok {
				for _, v := range valArr {
					if valMap, ok := v.(map[string]any); ok {
						item := amomodels.FieldValueElement{}
						if val, ok := valMap["value"]; ok {
							item.Value = val
						}
						if enum, ok := valMap["enum_code"].(string); ok {
							item.EnumCode = enum
						}
						cfValue.Values = append(cfValue.Values, item)
					}
				}
			}
			if len(cfValue.Values) > 0 {
				result = append(result, cfValue)
			}
		}
	}

	return result
}

func (r *Registry) RegisterComplexCreateTool() {
	r.addTool(genkit.DefineTool[gkitmodels.ComplexCreateInput, any](
		r.g,
		"complex_create",
		"Create a lead with contacts and company in one request. Supports custom fields, tags, and full contact/company data.",
		func(ctx *ai.ToolContext, input gkitmodels.ComplexCreateInput) (any, error) {
			lead := r.mapLeadToSDK(input)
			return r.complexCreateService.CreateComplex(ctx, lead)
		},
	))
}

func (r *Registry) RegisterComplexCreateBatchTool() {
	r.addTool(genkit.DefineTool[gkitmodels.ComplexCreateBatchInput, any](
		r.g,
		"complex_create_batch",
		"Create multiple leads with contacts and companies in one batch request (up to 50 items). Each item has lead, contacts, and company.",
		func(ctx *ai.ToolContext, input gkitmodels.ComplexCreateBatchInput) (any, error) {
			leads := make([]*amomodels.Lead, 0, len(input.Items))
			for _, item := range input.Items {
				leads = append(leads, r.mapLeadToSDK(item))
			}
			return r.complexCreateService.CreateComplexBatch(ctx, leads)
		},
	))
}
