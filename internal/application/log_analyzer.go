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
	"strconv"
	"strings"
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

// AnalyzeLogs processes all log files or URLs based on the provided paths
func (a *LogAnalyzer) AnalyzeLogs(from, to time.Time, filterField, filterValue string) error {
	for _, path := range a.Paths {
		if isURL(path) {
			fmt.Printf("Processing URL: %s\n", path)
			err := a.processURL(path, from, to, filterField, filterValue)
			if err != nil {
				fmt.Printf("Error processing URL %s: %v\n", path, err)
			}
		} else {
			err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// Проверяем, что это файл (а не директория)
				if !info.IsDir() {
					fmt.Printf("Processing file: %s\n", filePath)
					err := a.processFile(filePath, from, to, filterField, filterValue)
					if err != nil {
						fmt.Printf("Error processing file %s: %v\n", filePath, err)
					}
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("error walking the path %s: %w", path, err)
			}
		}
	}

	a.calculateRPS()
	return nil
}

// processFile processes a single file
func (a *LogAnalyzer) processFile(filePath string, from, to time.Time, filterField, filterValue string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	return a.processLogs(file, from, to, filterField, filterValue)
}

// processURL processes logs directly from a URL without loading into memory
func (a *LogAnalyzer) processURL(url string, from, to time.Time, filterField, filterValue string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch file from URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status for URL %s: %s", url, resp.Status)
	}

	// Process the logs directly from the response body
	return a.processLogs(resp.Body, from, to, filterField, filterValue)
}

// processLogs processes logs from an io.Reader line by line
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

		// Skip filtering by time range if 'from' and 'to' are zero values
		if !from.IsZero() && logRecord.Timestamp.Before(from) {
			continue
		}
		if !to.IsZero() && logRecord.Timestamp.After(to) {
			continue
		}

		// Apply filter based on field and value
		if filterField != "" && filterValue != "" {
			if !matchesFilter(logRecord, filterField, filterValue) {
				continue
			}
		}

		a.updateMetrics(logRecord)
	}

	fmt.Printf("Processed %d lines from log source\n", lineCount)
	return scanner.Err()
}

var unknownFieldWarned = false // глобальная переменная для отслеживания предупреждения

func matchesFilter(logRecord domain.LogRecord, field, value string) bool {
	// Убираем символ `*` из конца значения, если он есть
	isWildcard := strings.HasSuffix(value, "*")
	if isWildcard {
		value = strings.TrimSuffix(value, "*")
	}

	switch field {
	case "ip":
		if isWildcard {
			return strings.HasPrefix(logRecord.IP, value)
		}
		return logRecord.IP == value
	case "timestamp":
		filterTime, err := time.Parse(time.RFC3339, value)
		if err != nil {
			fmt.Printf("Invalid timestamp filter value: %v\n", err)
			return false
		}
		return logRecord.Timestamp.Equal(filterTime)
	case "method":
		if isWildcard {
			return strings.HasPrefix(strings.ToLower(logRecord.Method), strings.ToLower(value))
		}
		return strings.EqualFold(logRecord.Method, value)
	case "url":
		if isWildcard {
			return strings.HasPrefix(logRecord.URL, value)
		}
		return logRecord.URL == value
	case "protocol":
		if isWildcard {
			return strings.HasPrefix(strings.ToLower(logRecord.Protocol), strings.ToLower(value))
		}
		return strings.EqualFold(logRecord.Protocol, value)
	case "status":
		if isWildcard {
			return strings.HasPrefix(fmt.Sprintf("%d", logRecord.StatusCode), value)
		}
		return fmt.Sprintf("%d", logRecord.StatusCode) == value
	case "response_size":
		size, err := strconv.Atoi(value)
		if err != nil {
			fmt.Printf("Invalid response_size filter value: %v\n", err)
			return false
		}
		return logRecord.ResponseSize == size
	case "referer":
		if isWildcard {
			return strings.HasPrefix(logRecord.Referer, value)
		}
		return logRecord.Referer == value
	case "agent":
		if isWildcard {
			return strings.HasPrefix(logRecord.UserAgent, value)
		}
		return logRecord.UserAgent == value
	default:
		if !unknownFieldWarned {
			fmt.Printf("Unknown filter field: %s\n", field)
			unknownFieldWarned = true // предупреждение выведено
		}
		return false
	}
}

// updateMetrics updates the metrics based on the log record
func (a *LogAnalyzer) updateMetrics(logRecord domain.LogRecord) {
	metrics := a.Metrics

	metrics.TotalRequests++

	// Update StartDate and EndDate
	if metrics.StartDate.IsZero() || metrics.StartDate.After(logRecord.Timestamp) {
		metrics.StartDate = logRecord.Timestamp
	}
	if metrics.EndDate.IsZero() || metrics.EndDate.Before(logRecord.Timestamp) {
		metrics.EndDate = logRecord.Timestamp
	}

	metrics.TotalRespSize += logRecord.ResponseSize

	metrics.AverageRespSize = float64(metrics.TotalRespSize) / float64(metrics.TotalRequests)

	metrics.ResponseSizes = append(metrics.ResponseSizes, logRecord.ResponseSize)

	metrics.Percentile95 = a.CalculatePercentile(metrics.ResponseSizes, 95)

	metrics.Resources[logRecord.URL] += 1

	metrics.StatusCodes[logRecord.StatusCode] += 1

	// Add unique IP
	metrics.UniqueIPs[logRecord.IP] = struct{}{}
}

// calculateRPS calculates Requests Per Second (RPS)
func (a *LogAnalyzer) calculateRPS() {
	metrics := a.Metrics
	duration := metrics.EndDate.Sub(metrics.StartDate).Seconds()
	if duration > 0 {
		metrics.RPS = float64(metrics.TotalRequests) / duration
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

// isURL checks if a given path is a URL
func isURL(path string) bool {
	return len(path) > 4 && (path[:4] == "http" || path[:5] == "https")
}
