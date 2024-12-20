package infrastructure

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/abakunov/log-analyzer/internal/application"
)

// ParseTimeBounds parses the 'from' and 'to' time bounds from strings into time.Time format.
func ParseTimeBounds(fromStr, toStr string) (fromTime, toTime time.Time, err error) {
	// List of supported time formats.
	formats := []string{
		time.RFC3339, // Full ISO8601 format with time (e.g., "2015-05-18T00:00:00Z").
		"2006-01-02", // Date only (e.g., "2015-05-18").
	}

	// Parse the "from" parameter.
	if fromStr != "" {
		fromTime, err = parseTimeWithFormats(fromStr, formats)
		if err != nil {
			err = fmt.Errorf("invalid from time: %w", err)
			return
		}
	}

	// Parse the "to" parameter.
	if toStr != "" {
		toTime, err = parseTimeWithFormats(toStr, formats)
		if err != nil {
			err = fmt.Errorf("invalid to time: %w", err)
			return
		}
	}

	return
}

// parseTimeWithFormats tries to parse a time string with multiple formats.
func parseTimeWithFormats(input string, formats []string) (time.Time, error) {
	for _, format := range formats {
		parsedTime, err := time.Parse(format, input)
		if err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("could not parse time: %s", input)
}

// ParseFiles parses the file path or URL pattern into a list of paths.
func ParseFiles(pattern string) ([]string, error) {
	// Check if the path is a URL.
	if application.IsURL(pattern) {
		return []string{pattern}, nil
	}

	// Use filepath.Glob for local files.
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("error finding files: %v", err)
	}

	// Check if files were found.
	if len(files) == 0 {
		return nil, fmt.Errorf("no files found matching the pattern: %s", pattern)
	}

	return files, nil
}
