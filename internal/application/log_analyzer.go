package application

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/abakunov/log-analyzer/internal/domain"
)

type LogAnalyzer struct {
	Paths   []string
	Metrics *domain.Metrics
}

// NewLogAnalyzer creates a new LogAnalyzer.
func NewLogAnalyzer(paths []string) *LogAnalyzer {
	return &LogAnalyzer{
		Paths:   paths,
		Metrics: domain.NewMetrics(paths),
	}
}

// AnalyzeLogs processes all log files or URLs based on the provided paths.
func (a *LogAnalyzer) AnalyzeLogs(from, to time.Time, filterField, filterValue string) error {
	for _, path := range a.Paths {
		err := a.processPath(path, from, to, filterField, filterValue)
		if err != nil {
			fmt.Printf("Error processing path %s: %v\n", path, err)
		}
	}

	a.calculateRPS()

	return nil
}

// processPath determines whether the path is a URL or local file and processes it.
func (a *LogAnalyzer) processPath(path string, from, to time.Time, filterField, filterValue string) error {
	if isURL(path) {
		return a.processURLPath(path, from, to, filterField, filterValue)
	}

	return a.processLocalPath(path, from, to, filterField, filterValue)
}

// processURLPath processes a URL path.
func (a *LogAnalyzer) processURLPath(path string, from, to time.Time, filterField, filterValue string) error {
	fmt.Printf("Processing URL: %s\n", path)

	err := a.processURL(path, from, to, filterField, filterValue)
	if err != nil {
		return fmt.Errorf("error processing URL %s: %w", path, err)
	}

	return nil
}

// processLocalPath processes a local path using filepath.Walk.
func (a *LogAnalyzer) processLocalPath(path string, from, to time.Time, filterField, filterValue string) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories.
		if info.IsDir() {
			return nil
		}

		fmt.Printf("Processing file: %s\n", filePath)

		err = a.processFile(filePath, from, to, filterField, filterValue)
		if err != nil {
			return fmt.Errorf("error processing file %s: %w", filePath, err)
		}

		return nil
	})
}

// calculateRPS calculates Requests Per Second (RPS).
func (a *LogAnalyzer) calculateRPS() {
	duration := a.Metrics.EndDate.Sub(a.Metrics.StartDate).Seconds()
	if duration > 0 {
		a.Metrics.RPS = float64(a.Metrics.TotalRequests) / duration
	}
}
