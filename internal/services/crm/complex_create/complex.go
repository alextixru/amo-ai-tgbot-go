package complex_create

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) CreateComplex(ctx context.Context, input *gkitmodels.ComplexCreateInput) (*ComplexCreateResult, error) {
	lead, err := s.buildSDKLead(input)
	if err != nil {
		return nil, err
	}

	result, err := s.sdk.Leads().AddOneComplex(ctx, lead)
	if err != nil {
		return nil, err
	}

	return s.enrichResult(result.ID, result.ContactID, result.CompanyID, result.Merged, lead, input), nil
}

func (s *service) CreateComplexBatch(ctx context.Context, inputs []gkitmodels.ComplexCreateInput) ([]ComplexCreateResult, error) {
	leads := make([]*amomodels.Lead, 0, len(inputs))
	for i := range inputs {
		lead, err := s.buildSDKLead(&inputs[i])
		if err != nil {
			return nil, fmt.Errorf("элемент %d: %w", i, err)
		}
		leads = append(leads, lead)
	}

	results, err := s.sdk.Leads().AddComplex(ctx, leads)
	if err != nil {
		return nil, err
	}

	if len(results) != len(inputs) {
		return nil, fmt.Errorf("SDK вернул %d результатов для %d сделок", len(results), len(inputs))
	}

	out := make([]ComplexCreateResult, 0, len(results))
	for i, r := range results {
		out = append(out, *s.enrichResult(r.ID, r.ContactID, r.CompanyID, r.Merged, leads[i], &inputs[i]))
	}
	return out, nil
}

// buildSDKLead конвертирует ComplexCreateInput в SDK Lead, резолвя имена → ID.
func (s *service) buildSDKLead(input *gkitmodels.ComplexCreateInput) (*amomodels.Lead, error) {
	// Резолвим поля сделки
	pipelineID, err := s.resolvePipelineID(input.Lead.PipelineName)
	if err != nil {
		return nil, err
	}

	statusID, err := s.resolveStatusID(input.Lead.PipelineName, pipelineID, input.Lead.StatusName)
	if err != nil {
		return nil, err
	}

	responsibleID, err := s.resolveUserID(input.Lead.ResponsibleUserName)
	if err != nil {
		return nil, err
	}

	lead := &amomodels.Lead{
		Name:       input.Lead.Name,
		Price:      input.Lead.Price,
		PipelineID: pipelineID,
		StatusID:   statusID,
		Embedded:   &amomodels.LeadEmbedded{},
	}
	lead.ResponsibleUserID = responsibleID

	// Кастомные поля сделки
	if len(input.Lead.CustomFieldsValues) > 0 {
		lead.CustomFieldsValues = mapCustomFieldsValues(input.Lead.CustomFieldsValues)
	}

	// Теги
	if len(input.Lead.Tags) > 0 {
		lead.Embedded.Tags = make([]amomodels.Tag, 0, len(input.Lead.Tags))
		for _, tagName := range input.Lead.Tags {
			lead.Embedded.Tags = append(lead.Embedded.Tags, amomodels.Tag{Name: tagName})
		}
	}

	// Контакты
	if len(input.Contacts) > 0 {
		lead.Embedded.Contacts = make([]*amomodels.Contact, 0, len(input.Contacts))
		for _, c := range input.Contacts {
			contact, err := s.buildSDKContact(c)
			if err != nil {
				return nil, fmt.Errorf("контакт %q: %w", c.Name, err)
			}
			lead.Embedded.Contacts = append(lead.Embedded.Contacts, contact)
		}
	}

	// Компания
	if input.Company != nil {
		company, err := s.buildSDKCompany(input.Company)
		if err != nil {
			return nil, fmt.Errorf("компания %q: %w", input.Company.Name, err)
		}
		lead.Embedded.Companies = []*amomodels.Company{company}
	}

	return lead, nil
}

// buildSDKContact конвертирует ContactData в SDK Contact.
func (s *service) buildSDKContact(c gkitmodels.ContactData) (*amomodels.Contact, error) {
	responsibleID, err := s.resolveUserID(c.ResponsibleUserName)
	if err != nil {
		return nil, err
	}

	contact := &amomodels.Contact{
		Name:      c.Name,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		IsMain:    c.IsMain,
	}
	contact.ResponsibleUserID = responsibleID

	cfv := buildContactCustomFields(c)
	if len(cfv) > 0 {
		contact.CustomFieldsValues = cfv
	}

	return contact, nil
}

// buildSDKCompany конвертирует CompanyData в SDK Company.
func (s *service) buildSDKCompany(c *gkitmodels.CompanyData) (*amomodels.Company, error) {
	responsibleID, err := s.resolveUserID(c.ResponsibleUserName)
	if err != nil {
		return nil, err
	}

	company := &amomodels.Company{
		Name: c.Name,
	}
	company.ResponsibleUserID = responsibleID

	if len(c.CustomFieldsValues) > 0 {
		company.CustomFieldsValues = mapCustomFieldsValues(c.CustomFieldsValues)
	}

	return company, nil
}

// enrichResult строит обогащённый ответ для LLM, подставляя имена вместо числовых ID.
func (s *service) enrichResult(
	leadID int,
	contactID *int,
	companyID *int,
	merged bool,
	lead *amomodels.Lead,
	input *gkitmodels.ComplexCreateInput,
) *ComplexCreateResult {
	result := &ComplexCreateResult{
		Merged: merged,
		Lead: CreatedLeadView{
			ID:                  leadID,
			Name:                lead.Name,
			Price:               lead.Price,
			PipelineName:        s.lookupPipelineName(lead.PipelineID),
			StatusName:          s.lookupStatusName(lead.PipelineID, lead.StatusID),
			ResponsibleUserName: s.lookupUserName(lead.ResponsibleUserID),
			CreatedAt:           time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Контакты
	if contactID != nil && *contactID != 0 {
		views := make([]CreatedContactView, 0, len(input.Contacts))
		// Основной контакт из результата
		mainView := CreatedContactView{ID: *contactID}
		if len(input.Contacts) > 0 {
			mainView.Name = input.Contacts[0].Name
			mainView.ResponsibleUserName = input.Contacts[0].ResponsibleUserName
		}
		views = append(views, mainView)

		// Дополнительные контакты (без ID из ответа SDK, но сохраняем имена)
		if lead.Embedded != nil {
			for i, c := range lead.Embedded.Contacts {
				if i == 0 {
					continue // уже добавлен как основной
				}
				if c == nil {
					continue
				}
				view := CreatedContactView{
					Name: c.Name,
				}
				if i < len(input.Contacts) {
					view.ResponsibleUserName = input.Contacts[i].ResponsibleUserName
				}
				views = append(views, view)
			}
		}
		result.Contacts = views
	}

	// Компания
	if companyID != nil && *companyID != 0 {
		view := &CreatedCompanyView{ID: *companyID}
		if input.Company != nil {
			view.Name = input.Company.Name
			view.ResponsibleUserName = input.Company.ResponsibleUserName
		}
		result.Company = view
	}

	return result
}

// buildContactCustomFields собирает кастомные поля контакта из Phone, Email и CustomFieldsValues.
func buildContactCustomFields(c gkitmodels.ContactData) []amomodels.CustomFieldValue {
	var result []amomodels.CustomFieldValue

	if c.Phone != "" {
		result = append(result, amomodels.CustomFieldValue{
			FieldCode: "PHONE",
			Values: []amomodels.FieldValueElement{
				{Value: c.Phone, EnumCode: "WORK"},
			},
		})
	}

	if c.Email != "" {
		result = append(result, amomodels.CustomFieldValue{
			FieldCode: "EMAIL",
			Values: []amomodels.FieldValueElement{
				{Value: c.Email, EnumCode: "WORK"},
			},
		})
	}

	if len(c.CustomFieldsValues) > 0 {
		result = append(result, mapCustomFieldsValues(c.CustomFieldsValues)...)
	}

	return result
}

// mapCustomFieldsValues конвертирует map[string]any в []CustomFieldValue.
func mapCustomFieldsValues(cfv map[string]any) []amomodels.CustomFieldValue {
	data, err := json.Marshal(cfv)
	if err != nil {
		return nil
	}

	var result []amomodels.CustomFieldValue
	if err := json.Unmarshal(data, &result); err != nil {
		// Fallback: пробуем как map field_code -> values
		for fieldCode, values := range cfv {
			cfValue := amomodels.CustomFieldValue{
				FieldCode: fieldCode,
			}
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
