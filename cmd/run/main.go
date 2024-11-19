package main

import (
	"fmt"
	"github.com/abakunov/log-analyzer/internal/application"
	"github.com/abakunov/log-analyzer/internal/infrastructure"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var (
	globPattern string
	from        string
	to          string
	format      string
	filterField string
	filterValue string
)

// parseTimeBounds parses the 'from' and 'to' time bounds from strings into time.Time format.
func parseTimeBounds(fromStr, toStr string) (time.Time, time.Time, error) {
	var fromTime, toTime time.Time
	var err error

	// List of supported time formats.
	formats := []string{
		time.RFC3339, // Full ISO8601 format with time (e.g., "2015-05-18T00:00:00Z").
		"2006-01-02", // Date only (e.g., "2015-05-18").
	}

	// Parse the "from" parameter.
	if fromStr != "" {
		fromTime, err = parseTimeWithFormats(fromStr, formats)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid from time: %w", err)
		}
	}

	// Parse the "to" parameter.
	if toStr != "" {
		toTime, err = parseTimeWithFormats(toStr, formats)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid to time: %w", err)
		}
	}

	return fromTime, toTime, nil
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

// parseFiles parses the file path or URL pattern into a list of paths.
func parseFiles(pattern string) ([]string, error) {
	// Check if the path is a URL.
	if isURL(pattern) {
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

// isURL checks if a given path is a URL.
func isURL(path string) bool {
	return len(path) > 4 && (path[:4] == "http" || path[:5] == "https")
}

// runAnalyzer handles the log analysis process by parsing inputs and generating reports.
func runAnalyzer() {
	fromTime, toTime, err := parseTimeBounds(from, to)
	if err != nil {
		log.Fatalf("Error parsing time bounds: %v", err)
	}

	paths, err := parseFiles(globPattern)
	if err != nil {
		log.Fatalf("Error parsing files: %v", err)
	}

	analyzer := application.NewLogAnalyzer(paths)
	err = analyzer.AnalyzeLogs(fromTime, toTime, filterField, filterValue)
	if err != nil {
		log.Fatalf("Error analyzing logs: %v", err)
	}

	metrics := analyzer.Metrics
	formatter := infrastructure.ReportFormatter{Metrics: metrics}
	output := infrastructure.ReportOutput{}

	// If the format is not specified, print the report to the console.
	if format == "" {
		output.OutputToConsole(formatter.RenderConsole())
		return
	}

	// Generate the report in the specified format.
	var report string
	outputFile := "log_report.md"
	if format == "markdown" {
		report = formatter.RenderMarkdown()
	} else if format == "adoc" {
		report = formatter.RenderAdoc()
		outputFile = "log_report.adoc"
	} else {
		log.Fatalf("Unsupported format: %s", format)
	}

	err = output.OutputToFile(report, outputFile)
	if err != nil {
		log.Fatalf("Error saving report: %v", err)
	}

	fmt.Printf("Report saved as %s\n", outputFile)
}

var rootCmd = &cobra.Command{
	Use:   "analyzer",
	Short: "Analyze NGINX log files.",
	Run: func(cmd *cobra.Command, args []string) {
		runAnalyzer()
	},
}

// init initializes command-line flags and their descriptions.
func init() {
	rootCmd.Flags().StringVar(&globPattern, "path", "", "Path(s) to log files (required).")
	rootCmd.Flags().StringVar(&from, "from", "", "Start date in ISO8601 format (optional).")
	rootCmd.Flags().StringVar(&to, "to", "", "End date in ISO8601 format (optional).")
	rootCmd.Flags().StringVar(&format, "format", "", "Output format: markdown or adoc (optional).")
	rootCmd.Flags().StringVar(&filterField, "filter-field", "", "Field to filter logs by (optional).")
	rootCmd.Flags().StringVar(&filterValue, "filter-value", "", "Value to filter logs by (supports glob patterns, optional).")

	err := rootCmd.MarkFlagRequired("path")
	if err != nil {
		return
	}
}

// main is the entry point of the program.
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
