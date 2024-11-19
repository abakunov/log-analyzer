package application

import (
	"fmt"
	"net/http"
	"net/url"
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
func (a *LogAnalyzer) processURL(rawURL string, from, to time.Time, filterField, filterValue string) error {
	// Parse and validate the URL.
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL %s: %w", rawURL, err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme for %s", rawURL)
	}

	// Perform the HTTP GET request.
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return fmt.Errorf("failed to fetch file from URL %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status for URL %s: %s", rawURL, resp.Status)
	}

	return a.processLogs(resp.Body, from, to, filterField, filterValue)
}
