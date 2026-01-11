package entities

import (
	"encoding/json"

	amomodels "github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models"
)

// mapCustomFieldsValues конвертирует map[string]any в []CustomFieldValue
func mapCustomFieldsValues(cfv map[string]any) []amomodels.CustomFieldValue {
	if cfv == nil {
		return nil
	}
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

func mapTags(tags []gkitmodels.EntityTag) []amomodels.Tag {
	if len(tags) == 0 {
		return nil
	}
	result := make([]amomodels.Tag, len(tags))
	for i, t := range tags {
		result[i] = amomodels.Tag{}
		if t.ID > 0 {
			result[i].ID = t.ID
		}
		result[i].Name = t.Name
	}
	return result
}
