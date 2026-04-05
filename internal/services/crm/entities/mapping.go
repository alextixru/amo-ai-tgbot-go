package entities

import (
	"encoding/json"
	"fmt"
	"time"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// unixToISO конвертирует Unix timestamp в ISO-8601 строку. Возвращает "" если 0.
func unixToISO(ts int64) string {
	if ts == 0 {
		return ""
	}
	return time.Unix(ts, 0).UTC().Format(time.RFC3339)
}

// parseISO конвертирует ISO-8601 строку в Unix timestamp. Возвращает 0 при ошибке.
func parseISO(s string) int {
	if s == "" {
		return 0
	}
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return int(t.Unix())
		}
	}
	return 0
}

// leadToResult конвертирует SDK Lead в EntityResult.
func (s *service) leadToResult(lead *amomodels.Lead) *EntityResult {
	if lead == nil {
		return nil
	}
	r := &EntityResult{
		ID:                  lead.ID,
		Name:                lead.Name,
		Price:               lead.Price,
		PipelineName:        s.lookupPipelineName(lead.PipelineID),
		StatusName:          s.lookupStatusName(lead.StatusID),
		ResponsibleUserName: s.lookupUserName(lead.ResponsibleUserID),
		CreatedByName:       s.lookupUserName(lead.CreatedBy),
		UpdatedByName:       s.lookupUserName(lead.UpdatedBy),
		CreatedAt:           unixToISO(lead.CreatedAt),
		UpdatedAt:           unixToISO(lead.UpdatedAt),
	}

	if lead.LossReasonID != nil {
		r.LossReason = s.lookupLossReasonName(*lead.LossReasonID)
	}
	if lead.ClosedAt != nil {
		r.ClosedAt = unixToISO(*lead.ClosedAt)
	}

	// Embedded
	if lead.Embedded != nil {
		// Loss reason name из embedded (если запрашивали with=loss_reason)
		if len(lead.Embedded.LossReason) > 0 && r.LossReason == "" {
			r.LossReason = lead.Embedded.LossReason[0].Name
		}
		// Source name из embedded (если запрашивали with=source)
		if lead.Embedded.Source != nil && lead.Embedded.Source.Name != "" {
			r.SourceName = lead.Embedded.Source.Name
		}
		// Tags
		for _, t := range lead.Embedded.Tags {
			r.Tags = append(r.Tags, t.Name)
		}
		// Contacts
		for _, c := range lead.Embedded.Contacts {
			if c != nil {
				r.Contacts = append(r.Contacts, EntityRef{ID: c.ID, Name: c.Name})
			}
		}
		// Companies
		for _, c := range lead.Embedded.Companies {
			if c != nil {
				r.Companies = append(r.Companies, EntityRef{ID: c.ID, Name: c.Name})
			}
		}
	}

	// Custom fields
	r.CustomFieldsValues = mapCFVToEntries(lead.CustomFieldsValues)

	return r
}

// contactToResult конвертирует SDK Contact в EntityResult.
func (s *service) contactToResult(contact *amomodels.Contact) *EntityResult {
	if contact == nil {
		return nil
	}
	r := &EntityResult{
		ID:                  contact.ID,
		Name:                contact.Name,
		FirstName:           contact.FirstName,
		LastName:            contact.LastName,
		ResponsibleUserName: s.lookupUserName(contact.ResponsibleUserID),
		CreatedByName:       s.lookupUserName(contact.CreatedBy),
		UpdatedByName:       s.lookupUserName(contact.UpdatedBy),
		CreatedAt:           unixToISO(contact.CreatedAt),
		UpdatedAt:           unixToISO(contact.UpdatedAt),
	}

	if contact.Embedded != nil {
		for _, t := range contact.Embedded.Tags {
			r.Tags = append(r.Tags, t.Name)
		}
		for _, lead := range contact.Embedded.Leads {
			if lead != nil {
				r.Leads = append(r.Leads, EntityRef{ID: lead.ID, Name: lead.Name})
			}
		}
		if contact.Embedded.Company != nil {
			r.Companies = append(r.Companies, EntityRef{ID: contact.Embedded.Company.ID, Name: contact.Embedded.Company.Name})
		}
		for _, c := range contact.Embedded.Companies {
			if c != nil {
				r.Companies = append(r.Companies, EntityRef{ID: c.ID, Name: c.Name})
			}
		}
	}

	r.CustomFieldsValues = mapCFVToEntries(contact.CustomFieldsValues)
	return r
}

// companyToResult конвертирует SDK Company в EntityResult.
func (s *service) companyToResult(company *amomodels.Company) *EntityResult {
	if company == nil {
		return nil
	}
	r := &EntityResult{
		ID:                  company.ID,
		Name:                company.Name,
		ResponsibleUserName: s.lookupUserName(company.ResponsibleUserID),
		CreatedByName:       s.lookupUserName(company.CreatedBy),
		UpdatedByName:       s.lookupUserName(company.UpdatedBy),
		CreatedAt:           unixToISO(company.CreatedAt),
		UpdatedAt:           unixToISO(company.UpdatedAt),
	}

	if company.Embedded != nil {
		for _, t := range company.Embedded.Tags {
			r.Tags = append(r.Tags, t.Name)
		}
		for _, lead := range company.Embedded.Leads {
			if lead != nil {
				r.Leads = append(r.Leads, EntityRef{ID: lead.ID, Name: lead.Name})
			}
		}
		for _, c := range company.Embedded.Contacts {
			if c != nil {
				r.Contacts = append(r.Contacts, EntityRef{ID: c.ID, Name: c.Name})
			}
		}
	}

	r.CustomFieldsValues = mapCFVToEntries(company.CustomFieldsValues)
	return r
}

// mapCFVToEntries конвертирует []CustomFieldValue в []CustomFieldEntry для read-модели.
func mapCFVToEntries(cfv []amomodels.CustomFieldValue) []CustomFieldEntry {
	if len(cfv) == 0 {
		return nil
	}
	out := make([]CustomFieldEntry, 0, len(cfv))
	for _, cf := range cfv {
		entry := CustomFieldEntry{
			FieldCode: cf.FieldCode,
			FieldName: cf.FieldName,
		}
		for _, v := range cf.Values {
			entry.Values = append(entry.Values, v.Value)
		}
		out = append(out, entry)
	}
	return out
}

// mapCustomFieldsValues конвертирует map[string]any в []CustomFieldValue для SDK.
// Ключ — field_code. Значение — строка/массив/{value, enum_code}.
func mapCustomFieldsValues(cfv map[string]any) []amomodels.CustomFieldValue {
	if cfv == nil {
		return nil
	}

	var result []amomodels.CustomFieldValue

	for fieldCode, rawVal := range cfv {
		cfValue := amomodels.CustomFieldValue{
			FieldCode: fieldCode,
		}

		switch v := rawVal.(type) {
		case string:
			// Простое строковое значение
			cfValue.Values = []amomodels.FieldValueElement{{Value: v}}

		case []any:
			// Массив значений
			for _, item := range v {
				elem := amomodels.FieldValueElement{}
				switch iv := item.(type) {
				case string:
					elem.Value = iv
				case map[string]any:
					if val, ok := iv["value"]; ok {
						elem.Value = val
					}
					if enumCode, ok := iv["enum_code"].(string); ok {
						elem.EnumCode = enumCode
					}
				default:
					elem.Value = fmt.Sprintf("%v", iv)
				}
				cfValue.Values = append(cfValue.Values, elem)
			}

		case map[string]any:
			// Одно сложное значение
			elem := amomodels.FieldValueElement{}
			if val, ok := v["value"]; ok {
				elem.Value = val
			}
			if enumCode, ok := v["enum_code"].(string); ok {
				elem.EnumCode = enumCode
			}
			cfValue.Values = []amomodels.FieldValueElement{elem}

		default:
			// Прочие скалярные типы — через JSON roundtrip как fallback
			data, err := json.Marshal(rawVal)
			if err == nil {
				cfValue.Values = []amomodels.FieldValueElement{{Value: string(data)}}
			}
		}

		if len(cfValue.Values) > 0 {
			result = append(result, cfValue)
		}
	}

	return result
}

// mapTags конвертирует []EntityTag в []amomodels.Tag.
func mapTags(tags []gkitmodels.EntityTag) []amomodels.Tag {
	if len(tags) == 0 {
		return nil
	}
	result := make([]amomodels.Tag, len(tags))
	for i, t := range tags {
		result[i] = amomodels.Tag{Name: t.Name}
		if t.ID > 0 {
			result[i].ID = t.ID
		}
	}
	return result
}

// buildCustomFieldsFilter конвертирует []CustomFieldFilter в map[int]interface{} для SDK фильтра.
// customFieldsByCode — маппинг code→id для соответствующего типа сущности.
func buildCustomFieldsFilter(filters []gkitmodels.CustomFieldFilter, customFieldsByCode map[string]int) map[int]interface{} {
	if len(filters) == 0 {
		return nil
	}
	result := make(map[int]interface{})
	for _, f := range filters {
		fieldID, ok := customFieldsByCode[f.FieldCode]
		if !ok || fieldID == 0 {
			// код не найден — пропускаем
			continue
		}
		if len(f.Values) == 1 {
			result[fieldID] = fmt.Sprintf("%v", f.Values[0])
		} else if len(f.Values) > 1 {
			strs := make([]string, 0, len(f.Values))
			for _, v := range f.Values {
				strs = append(strs, fmt.Sprintf("%v", v))
			}
			result[fieldID] = strs
		}
	}
	return result
}
