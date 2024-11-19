package application

import (
	"fmt"
	"github.com/abakunov/log-analyzer/internal/domain"
	"os"
	"path/filepath"
	"time"
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
				if !info.IsDir() { // Process files only.
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

// calculateRPS calculates Requests Per Second (RPS).
func (a *LogAnalyzer) calculateRPS() {
	duration := a.Metrics.EndDate.Sub(a.Metrics.StartDate).Seconds()
	if duration > 0 {
		a.Metrics.RPS = float64(a.Metrics.TotalRequests) / duration
	}
}
