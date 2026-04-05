package unsorted

import (
	"time"
)

// unixToRFC3339 конвертирует Unix timestamp в строку RFC3339.
func unixToRFC3339(ts int64) string {
	if ts == 0 {
		return ""
	}
	return time.Unix(ts, 0).UTC().Format(time.RFC3339)
}

// rfc3339ToUnix конвертирует строку RFC3339 в Unix timestamp.
// Возвращает 0 если строка пустая или не парсится.
func rfc3339ToUnix(s string) int64 {
	if s == "" {
		return 0
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return 0
	}
	return t.Unix()
}
