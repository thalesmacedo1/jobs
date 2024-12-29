package utils

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
	}{
		{
			input:    "01/02/2023",
			expected: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    " 15/08/1947 ",
			expected: time.Date(1947, 8, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "12/2023",
			expected: time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "31/02/2023",
			expected: time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "",
			expected: time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "   ",
			expected: time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "29/02/2020",
			expected: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "31/04/2021",
			expected: time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "15/13/2021",
			expected: time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "00/12/2021",
			expected: time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		result := ParseDate(test.input)
		if !result.Equal(test.expected) {
			t.Errorf("ParseDate(%q) = %v; esperado %v", test.input, result, test.expected)
		}
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput int
	}{
		{
			name:           "Valid integer",
			input:          "123",
			expectedOutput: 123,
		},
		{
			name:           "Scientific notation",
			input:          "1.68E+07",
			expectedOutput: 16800000,
		},
		{
			name:           "Integer with spaces",
			input:          "   456   ",
			expectedOutput: 456,
		},
		{
			name:           "Empty string",
			input:          "",
			expectedOutput: 0,
		},
		{
			name:           "Non-numeric string",
			input:          "abcd",
			expectedOutput: 0,
		},
		{
			name:           "Scientific notation with negative exponent",
			input:          "1.23E-4",
			expectedOutput: 0, // Converts to 0 when cast to int
		},
		{
			name:           "Floating point number",
			input:          "789.56",
			expectedOutput: 789, // Truncated to 789
		},
		{
			name:           "Negative integer",
			input:          "-42",
			expectedOutput: -42,
		},
		{
			name:           "Zero",
			input:          "0",
			expectedOutput: 0,
		},
		{
			name:           "Large number",
			input:          "9.99E+08",
			expectedOutput: 999000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseInt(tt.input)
			if result != tt.expectedOutput {
				t.Errorf("ParseInt(%q) = %d; want %d", tt.input, result, tt.expectedOutput)
			}
		})
	}
}
