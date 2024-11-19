package main

import (
	"fmt"
	"github.com/abakunov/log-analyzer/internal/application"
	"github.com/abakunov/log-analyzer/internal/infrastructure"
	"log"
	"os"

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

// runAnalyzer handles the log analysis process by parsing inputs and generating reports.
func runAnalyzer() {
	// Parse time bounds from user input.
	fromTime, toTime, err := infrastructure.ParseTimeBounds(from, to)
	if err != nil {
		log.Fatalf("Error parsing time bounds: %v", err)
	}

	// Parse file paths or URLs from the glob pattern.
	paths, err := infrastructure.ParseFiles(globPattern)
	if err != nil {
		log.Fatalf("Error parsing files: %v", err)
	}

	// Initialize the log analyzer and analyze logs.
	analyzer := application.NewLogAnalyzer(paths)
	err = analyzer.AnalyzeLogs(fromTime, toTime, filterField, filterValue)
	if err != nil {
		log.Fatalf("Error analyzing logs: %v", err)
	}

	// Create the metrics object and formatter.
	metrics := analyzer.Metrics
	formatter := infrastructure.ReportFormatter{Metrics: metrics}
	output := infrastructure.ReportOutput{}

	// Render the report based on the specified format.
	if format == "" {
		// Print the report to the console if no format is specified.
		output.OutputToConsole(formatter.Render("plain"))
	} else if format == "markdown" {
		report := formatter.Render("markdown")
		err := output.OutputToFile(report, "log_report.md")
		if err != nil {
			log.Fatalf("Error saving markdown report: %v", err)
		}
		fmt.Println("Report saved as log_report.md")
	} else if format == "adoc" {
		report := formatter.Render("adoc")
		err := output.OutputToFile(report, "log_report.adoc")
		if err != nil {
			log.Fatalf("Error saving AsciiDoc report: %v", err)
		}
		fmt.Println("Report saved as log_report.adoc")
	} else {
		log.Fatalf("Unsupported format: %s", format)
	}
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
