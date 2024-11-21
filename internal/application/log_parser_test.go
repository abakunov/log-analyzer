package application_test

import (
	"testing"
	"time"

	"github.com/abakunov/log-analyzer/internal/application"
	"github.com/abakunov/log-analyzer/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestParseLogLine(t *testing.T) {
	tests := []struct {
		name      string
		logLine   string
		expected  domain.LogRecord
		expectErr bool
	}{
		{
			name:    "Valid log line with all fields",
			logLine: `127.0.0.1 - - [12/Dec/2021:19:01:02 +0000] "GET /index.html HTTP/1.1" 200 1024 "http://example.com" "Mozilla/5.0"`,
			expected: domain.LogRecord{
				IP:           "127.0.0.1",
				Timestamp:    time.Date(2021, time.December, 12, 19, 1, 2, 0, time.UTC),
				Method:       "GET",
				URL:          "/index.html",
				Protocol:     "HTTP/1.1",
				StatusCode:   200,
				ResponseSize: 1024,
				Referer:      "http://example.com",
				UserAgent:    "Mozilla/5.0",
			},
			expectErr: false,
		},
		{
			name:    "Valid log line with size as '-'",
			logLine: `127.0.0.1 - - [12/Dec/2021:19:01:02 +0000] "POST /submit HTTP/1.1" 404 - "-" "-"`,
			expected: domain.LogRecord{
				IP:           "127.0.0.1",
				Timestamp:    time.Date(2021, time.December, 12, 19, 1, 2, 0, time.UTC),
				Method:       "POST",
				URL:          "/submit",
				Protocol:     "HTTP/1.1",
				StatusCode:   404,
				ResponseSize: 0,
				Referer:      "-",
				UserAgent:    "-",
			},
			expectErr: false,
		},
		{
			name:      "Invalid log line format",
			logLine:   `Invalid log line format`,
			expected:  domain.LogRecord{},
			expectErr: true,
		},
		{
			name:      "Invalid timestamp format",
			logLine:   `127.0.0.1 - - [2021/Dec/12:19:01:02 +0000] "GET /index.html HTTP/1.1" 200 1024 "http://example.com" "Mozilla/5.0"`,
			expected:  domain.LogRecord{},
			expectErr: true,
		},
		{
			name:      "Invalid status code",
			logLine:   `127.0.0.1 - - [12/Dec/2021:19:01:02 +0000] "GET /index.html HTTP/1.1" abc 1024 "http://example.com" "Mozilla/5.0"`,
			expected:  domain.LogRecord{},
			expectErr: true,
		},
		{
			name:      "Invalid response size",
			logLine:   `127.0.0.1 - - [12/Dec/2021:19:01:02 +0000] "GET /index.html HTTP/1.1" 200 invalid_size "http://example.com" "Mozilla/5.0"`,
			expected:  domain.LogRecord{},
			expectErr: true,
		},
		{
			name:      "Missing required fields",
			logLine:   `127.0.0.1 - - [12/Dec/2021:19:01:02 +0000] "GET /index.html HTTP/1.1" 200`,
			expected:  domain.LogRecord{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record, err := application.ParseLogLine(tt.logLine)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, record)
			}
		})
	}
}
