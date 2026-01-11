package utils

import (
	"strings"
	"testing"
	"time"
)

func TestParseHumanDeadline(t *testing.T) {
	now := time.Now()
	location := now.Location()

	tests := []struct {
		input    string
		expected int64
		wantErr  bool
	}{
		{
			input:    "today",
			expected: time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, location).Unix(),
		},
		{
			input:    "tomorrow",
			expected: time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, location).AddDate(0, 0, 1).Unix(),
		},
		{
			input:    "in 2 hours",
			expected: now.Add(2 * time.Hour).Unix(),
		},
		{
			input:    "in 1 day",
			expected: now.AddDate(0, 0, 1).Unix(),
		},
		{
			input:    "2025-12-31",
			expected: time.Date(2025, 12, 31, 12, 0, 0, 0, location).Unix(),
		},
		{
			input:    "2025-12-31 15:30",
			expected: time.Date(2025, 12, 31, 15, 30, 0, 0, location).Unix(),
		},
		{
			input:    "invalid",
			expected: 0,
			wantErr:  true,
		},
		{
			input:    "",
			expected: 0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseHumanDeadline(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHumanDeadline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Для относительного времени даем погрешность в 2 секунды
			if strings.HasPrefix(tt.input, "in ") || tt.input == "" {
				if got < tt.expected-2 || got > tt.expected+2 {
					t.Errorf("ParseHumanDeadline() = %v, want %v (with 2s tolerance)", got, tt.expected)
				}
			} else {
				if got != tt.expected {
					t.Errorf("ParseHumanDeadline() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}
