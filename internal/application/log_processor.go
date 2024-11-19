package application

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

// processLogs processes logs from an io.Reader line by line.
func (a *LogAnalyzer) processLogs(reader io.Reader, from, to time.Time, filterField, filterValue string) error {
	scanner := bufio.NewScanner(reader)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		logRecord, err := ParseLogLine(line)

		if err != nil {
			fmt.Printf("Error parsing line: %s, Error: %v\n", line, err)
			continue
		}

		// Time range filter.
		if !from.IsZero() && logRecord.Timestamp.Before(from) {
			continue
		}

		if !to.IsZero() && logRecord.Timestamp.After(to) {
			continue
		}

		// Apply additional filters.
		if filterField != "" && filterValue != "" {
			if !matchesFilter(&logRecord, filterField, filterValue) {
				continue
			}
		}

		a.updateMetrics(&logRecord)
	}

	fmt.Printf("Processed %d lines from log source\n", lineCount)
	fmt.Println()

	return scanner.Err()
}
