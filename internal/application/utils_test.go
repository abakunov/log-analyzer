package application_test

import (
	"testing"

	"github.com/abakunov/log-analyzer/internal/application"

	"github.com/stretchr/testify/assert"
)

func TestIsURL(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid HTTP URL",
			input:    "http://example.com",
			expected: true,
		},
		{
			name:     "Valid HTTPS URL",
			input:    "https://example.com",
			expected: true,
		},
		{
			name:     "Invalid URL (short path)",
			input:    "ftp://example.com",
			expected: false,
		},
		{
			name:     "Invalid URL (missing protocol)",
			input:    "example.com",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "File path with forward slashes",
			input:    "/usr/local/bin",
			expected: false,
		},
		{
			name:     "File path with backslashes",
			input:    "C:\\Program Files",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := application.IsURL(tc.input)
			assert.Equal(t, tc.expected, actual, "Input: %s", tc.input)
		})
	}
}
