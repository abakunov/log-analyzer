package infrastructure_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/abakunov/log-analyzer/internal/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestParseTimeBounds(t *testing.T) {
	testCases := []struct {
		name         string
		fromStr      string
		toStr        string
		expectedFrom time.Time
		expectedTo   time.Time
		expectErr    bool
	}{
		{
			name:         "Valid ISO8601 format",
			fromStr:      "2023-11-20T15:04:05Z",
			toStr:        "2023-11-21T15:04:05Z",
			expectedFrom: time.Date(2023, time.November, 20, 15, 4, 5, 0, time.UTC),
			expectedTo:   time.Date(2023, time.November, 21, 15, 4, 5, 0, time.UTC),
			expectErr:    false,
		},
		{
			name:         "Valid date only format",
			fromStr:      "2023-11-20",
			toStr:        "2023-11-21",
			expectedFrom: time.Date(2023, time.November, 20, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2023, time.November, 21, 0, 0, 0, 0, time.UTC),
			expectErr:    false,
		},
		{
			name:      "Invalid from date",
			fromStr:   "invalid-date",
			toStr:     "2023-11-21",
			expectErr: true,
		},
		{
			name:      "Invalid to date",
			fromStr:   "2023-11-20",
			toStr:     "invalid-date",
			expectErr: true,
		},
		{
			name:      "Empty strings",
			fromStr:   "",
			toStr:     "",
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			from, to, err := infrastructure.ParseTimeBounds(tc.fromStr, tc.toStr)

			if tc.expectErr {
				assert.Error(t, err, "Expected an error, but got none.")
			} else {
				assert.NoError(t, err, "Did not expect an error, but got one.")
				assert.Equal(t, tc.expectedFrom, from, "From time mismatch.")
				assert.Equal(t, tc.expectedTo, to, "To time mismatch.")
			}
		})
	}
}

func TestParseFiles(t *testing.T) {
	testCases := []struct {
		name      string
		pattern   string
		mockFiles []string
		expectErr bool
		expected  []string
	}{
		{
			name:      "Valid local files",
			pattern:   "testdata/*.log",
			mockFiles: []string{"testdata/file1.log", "testdata/file2.log"},
			expected:  []string{"testdata/file1.log", "testdata/file2.log"},
			expectErr: false,
		},
		{
			name:      "No matching files",
			pattern:   "testdata/nonexistent/*.log",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "URL pattern",
			pattern:   "http://example.com/file.log",
			expected:  []string{"http://example.com/file.log"},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.mockFiles) > 0 {
				err := setupMockFiles(tc.mockFiles)
				assert.NoError(t, err, "Failed to set up mock files.")
				defer cleanupMockFiles(tc.mockFiles)
			}

			files, err := infrastructure.ParseFiles(tc.pattern)

			if tc.expectErr {
				assert.Error(t, err, "Expected an error, but got none.")
			} else {
				assert.NoError(t, err, "Did not expect an error, but got one.")
				assert.ElementsMatch(t, tc.expected, files, "File list mismatch.")
			}
		})
	}
}

// Helper functions for mock file setup and cleanup.
func setupMockFiles(files []string) error {
	for _, file := range files {
		dir := filepath.Dir(file)
		err := os.MkdirAll(dir, os.ModePerm)

		if err != nil {
			return err
		}

		_, err = os.Create(file)

		if err != nil {
			return err
		}
	}

	return nil
}

func cleanupMockFiles(files []string) {
	for _, file := range files {
		os.RemoveAll(filepath.Dir(file))
	}
}
