package utils

import "encoding/json"

// ToMap конвертирует любую структуру в map[string]any.
// Использует JSON сериализацию — все поля с json тегами включаются автоматически.
// При изменении структуры — map обновляется автоматически.
func ToMap(v any) map[string]any {
	if v == nil {
		return nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	return result
}
