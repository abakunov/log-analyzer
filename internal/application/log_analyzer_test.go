package application_test

import (
	"os"
	"testing"
	"time"

	"github.com/abakunov/log-analyzer/internal/application"
	"github.com/stretchr/testify/assert"
)

// MockLogFileGenerator отвечает за создание тестовых данных в виде лог-файлов.
type MockLogFileGenerator struct{}

func (m *MockLogFileGenerator) GenerateLogFile(filePath, content string) error {
	err := os.MkdirAll("testdata", os.ModePerm)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(content), 0o600)
}

func (m *MockLogFileGenerator) Cleanup() {
	os.RemoveAll("testdata")
}

func TestLogAnalyzer_AnalyzeLogs(t *testing.T) {
	mockGenerator := &MockLogFileGenerator{}

	// Данные для лог-файла.
	logData := `127.0.0.1 - - [12/Dec/2021:15:04:05 +0000] "GET /index.html HTTP/1.1" 200 1024 "http://example.com" "Mozilla/5.0"
127.0.0.1 - - [12/Dec/2021:16:04:05 +0000] "POST /submit HTTP/1.1" 404 - "-" "-"
192.168.1.1 - - [13/Dec/2021:15:04:05 +0000] "GET /home HTTP/1.1" 200 512 "-" "Mozilla/5.0"`

	// Генерация тестового лог-файла.
	err := mockGenerator.GenerateLogFile("testdata/logfile.log", logData)
	assert.NoError(t, err, "Failed to create logfile.log.")

	defer mockGenerator.Cleanup() // Удалить тестовые данные после тестов.

	// Массив с тестовыми сценариями.
	testCases := []struct {
		name         string
		paths        []string
		from         time.Time
		to           time.Time
		filterField  string
		filterValue  string
		expectedErr  bool
		expectedReqs int
	}{
		{
			name:         "Filter by IP",
			paths:        []string{"testdata/logfile.log"},
			from:         time.Time{},
			to:           time.Time{},
			filterField:  "ip",
			filterValue:  "127.0.0.1",
			expectedErr:  false,
			expectedReqs: 2,
		},
		{
			name:         "Filter with date range",
			paths:        []string{"testdata/logfile.log"},
			from:         time.Date(2021, time.December, 12, 0, 0, 0, 0, time.UTC),
			to:           time.Date(2021, time.December, 13, 0, 0, 0, 0, time.UTC),
			filterField:  "",
			filterValue:  "",
			expectedErr:  false,
			expectedReqs: 2,
		},
		{
			name:         "Filter with no matches",
			paths:        []string{"testdata/logfile.log"},
			from:         time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			to:           time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
			filterField:  "",
			filterValue:  "",
			expectedErr:  false,
			expectedReqs: 0,
		},
		{
			name:         "Filter by Method",
			paths:        []string{"testdata/logfile.log"},
			from:         time.Time{},
			to:           time.Time{},
			filterField:  "method",
			filterValue:  "GET",
			expectedErr:  false,
			expectedReqs: 2,
		},
		{
			name:         "Filter by Status Code",
			paths:        []string{"testdata/logfile.log"},
			from:         time.Time{},
			to:           time.Time{},
			filterField:  "status",
			filterValue:  "404",
			expectedErr:  false,
			expectedReqs: 1,
		},
		{
			name:         "Filter by URL",
			paths:        []string{"testdata/logfile.log"},
			from:         time.Time{},
			to:           time.Time{},
			filterField:  "url",
			filterValue:  "/home",
			expectedErr:  false,
			expectedReqs: 1,
		},
		{
			name:         "Filter by Protocol",
			paths:        []string{"testdata/logfile.log"},
			from:         time.Time{},
			to:           time.Time{},
			filterField:  "protocol",
			filterValue:  "HTTP/1.1",
			expectedErr:  false,
			expectedReqs: 3,
		},
		{
			name:         "Filter by User Agent",
			paths:        []string{"testdata/logfile.log"},
			from:         time.Time{},
			to:           time.Time{},
			filterField:  "agent",
			filterValue:  "Mozilla/5.0",
			expectedErr:  false,
			expectedReqs: 2,
		},
		{
			name:         "Filter by Referer",
			paths:        []string{"testdata/logfile.log"},
			from:         time.Time{},
			to:           time.Time{},
			filterField:  "referer",
			filterValue:  "http://example.com",
			expectedErr:  false,
			expectedReqs: 1,
		},
	}

	// Итерация через тестовые сценарии.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analyzer := application.NewLogAnalyzer(tc.paths)
			err := analyzer.AnalyzeLogs(tc.from, tc.to, tc.filterField, tc.filterValue)

			// Проверка на ошибку.
			if tc.expectedErr {
				assert.Error(t, err, "Expected an error, but got none.")
			} else {
				assert.NoError(t, err, "Expected no error, but got one.")
			}

			// Проверка количества обработанных запросов.
			assert.Equal(t, tc.expectedReqs, analyzer.Metrics.TotalRequests, "TotalRequests mismatch.")
		})
	}
}
