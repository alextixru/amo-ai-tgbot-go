package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseHumanDeadline конвертирует человеко-понятную строку в Unix timestamp.
//
// Поддерживаемые форматы:
//
//	"today"                 → конец сегодняшнего дня (23:59:59)
//	"tomorrow"              → конец завтрашнего дня (23:59:59)
//	"in 2 hours"            → текущее время + 2 часа
//	"in 3 days"             → текущее время + 3 дня
//	"2024-01-15"            → указанная дата в 12:00
//	"2024-01-15T14:00"      → ISO дата и время
//	"2024-01-15 14:00"      → дата и время через пробел
func ParseHumanDeadline(input string) (int64, error) {
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return 0, nil
	}

	now := time.Now()

	// 1. Фиксированные слова
	switch input {
	case "today":
		return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).Unix(), nil
	case "tomorrow":
		tomorrow := now.AddDate(0, 0, 1)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 59, 0, now.Location()).Unix(), nil
	}

	// 2. Относительное время: "in 2 hours", "in 3 days"
	reIn := regexp.MustCompile(`^in\s+(\d+)\s+(hour|min|minute|day|week)s?$`)
	if matches := reIn.FindStringSubmatch(input); len(matches) == 3 {
		count, _ := strconv.Atoi(matches[1])
		unit := matches[2]

		switch unit {
		case "hour":
			return now.Add(time.Duration(count) * time.Hour).Unix(), nil
		case "min", "minute":
			return now.Add(time.Duration(count) * time.Minute).Unix(), nil
		case "day":
			return now.AddDate(0, 0, count).Unix(), nil
		case "week":
			return now.AddDate(0, 0, count*7).Unix(), nil
		}
	}

	// 3. Форматы дат
	layouts := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, input, now.Location()); err == nil {
			// Если ввели только дату, ставим середину дня (12:00)
			if !strings.Contains(layout, "15:04") {
				return time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, now.Location()).Unix(), nil
			}
			return t.Unix(), nil
		}
	}

	return 0, fmt.Errorf("unknown deadline format: %s", input)
}
