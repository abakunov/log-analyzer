package application

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// processFile processes a single file.
func (a *LogAnalyzer) processFile(filePath string, from, to time.Time, filterField, filterValue string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()
	return a.processLogs(file, from, to, filterField, filterValue)
}

// processURL processes logs directly from a URL without loading into memory.
func (a *LogAnalyzer) processURL(url string, from, to time.Time, filterField, filterValue string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch file from URL %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status for URL %s: %s", url, resp.Status)
	}
	return a.processLogs(resp.Body, from, to, filterField, filterValue)
}
