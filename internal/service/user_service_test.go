package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "birthday already passed this year",
			dob:      time.Date(now.Year()-30, now.Month()-1, now.Day(), 0, 0, 0, 0, time.UTC),
			expected: 30,
		},
		{
			name:     "birthday not yet this year",
			dob:      time.Date(now.Year()-25, now.Month()+1, now.Day(), 0, 0, 0, 0, time.UTC),
			expected: 24,
		},
		{
			name:     "birthday is today",
			dob:      time.Date(now.Year()-20, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			expected: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAge(tt.dob)
			if got != tt.expected {
				t.Errorf("CalculateAge() = %d, want %d", got, tt.expected)
			}
		})
	}
}
