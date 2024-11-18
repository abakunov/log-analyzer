package application

import (
	"bufio"
	"fmt"
	"github.com/abakunov/log-analyzer/internal/domain"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type LogAnalyzer struct {
	Paths   []string
	Metrics *domain.Metrics
}

// NewLogAnalyzer creates a new LogAnalyzer
func NewLogAnalyzer(paths []string) *LogAnalyzer {
	return &LogAnalyzer{
		Paths:   paths,
		Metrics: domain.NewMetrics(paths),
	}
}

// CalculatePercentile counts the value of the specified percentile
func (a *LogAnalyzer) CalculatePercentile(values []int, percentile float64) int {
	if len(values) == 0 {
		return 0
	}
	sort.Ints(values)
	index := int(float64(len(values)) * percentile / 100)
	if index >= len(values) {
		index = len(values) - 1
	}
	return values[index]
}

// AnalyzeLogs analyzes logs from multiple sources (local files or URLs) in the specified time range
func (a *LogAnalyzer) AnalyzeLogs(from, to time.Time) error {
	for _, path := range a.Paths {
		if isURL(path) {
			// Analyze from URL
			if err := a.analyzeFromURL(path, from, to); err != nil {
				return err
			}
		} else {
			// Analyze from local files with glob support
			matches, err := filepath.Glob(path)
			if err != nil {
				return fmt.Errorf("invalid glob pattern %s: %w", path, err)
			}
			for _, file := range matches {
				if err := a.analyzeFromFile(file, from, to); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// analyzeFromURL processes logs from a URL
func (a *LogAnalyzer) analyzeFromURL(url string, from, to time.Time) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch logs from URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching logs from URL %s: %s", url, resp.Status)
	}

	return a.processLogs(resp.Body, from, to)
}

// analyzeFromFile processes logs from a local file
func (a *LogAnalyzer) analyzeFromFile(filePath string, from, to time.Time) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	return a.processLogs(file, from, to)
}

// processLogs processes logs from an io.Reader
func (a *LogAnalyzer) processLogs(reader io.Reader, from, to time.Time) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		logRecord, err := ParseLogLine(scanner.Text())
		if err != nil {
			continue
		}

		// Skip filtering by time range if 'from' and 'to' are zero values
		if !from.IsZero() && logRecord.Timestamp.Before(from) {
			continue
		}
		if !to.IsZero() && logRecord.Timestamp.After(to) {
			continue
		}

		a.updateMetrics(logRecord)
	}

	return scanner.Err()
}

// updateMetrics updates the metrics based on the log record
func (a *LogAnalyzer) updateMetrics(logRecord domain.LogRecord) {
	metrics := a.Metrics

	metrics.TotalRequests++

	if metrics.StartDate.After(logRecord.Timestamp) {
		metrics.StartDate = logRecord.Timestamp
	}
	if metrics.EndDate.Before(logRecord.Timestamp) {
		metrics.EndDate = logRecord.Timestamp
	}

	metrics.TotalRespSize += logRecord.ResponseSize

	metrics.AverageRespSize = float64(metrics.TotalRespSize) / float64(metrics.TotalRequests)

	metrics.ResponseSizes = append(metrics.ResponseSizes, logRecord.ResponseSize)

	metrics.Percentile95 = a.CalculatePercentile(metrics.ResponseSizes, 95)

	metrics.Resources[logRecord.URL]++

	metrics.StatusCodes[logRecord.StatusCode]++
}

// isURL checks if a path is a URL
func isURL(path string) bool {
	return len(path) > 6 && (path[:7] == "http://" || path[:8] == "https://")
}
