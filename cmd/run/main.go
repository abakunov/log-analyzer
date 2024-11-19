package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abakunov/log-analyzer/internal/application"
	"github.com/abakunov/log-analyzer/internal/infrastructure"

	"github.com/spf13/cobra"
)

var (
	globPattern string
	from        string
	to          string
	format      string
	filterField string
	filterValue string
	rootCmd     *cobra.Command
)

// setupRootCmd initializes the root command and its flags.
func setupRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyzer",
		Short: "Analyze NGINX log files.",
		Run: func(_ *cobra.Command, _ []string) {
			runAnalyzer()
		},
	}

	cmd.Flags().StringVar(&globPattern, "path", "", "Path(s) to log files (required).")
	cmd.Flags().StringVar(&from, "from", "", "Start date in ISO8601 format (optional).")
	cmd.Flags().StringVar(&to, "to", "", "End date in ISO8601 format (optional).")
	cmd.Flags().StringVar(&format, "format", "", "Output format: markdown or adoc (optional).")
	cmd.Flags().StringVar(&filterField, "filter-field", "", "Field to filter logs by (optional).")
	cmd.Flags().StringVar(&filterValue, "filter-value", "", "Value to filter logs by (supports glob patterns, optional).")

	err := cmd.MarkFlagRequired("path")
	if err != nil {
		log.Fatalf("Error marking path flag as required: %v", err)
	}

	return cmd
}

// runAnalyzer handles the log analysis process by parsing inputs and generating reports.
func runAnalyzer() {
	fromTime, toTime, err := infrastructure.ParseTimeBounds(from, to)
	if err != nil {
		log.Fatalf("Error parsing time bounds: %v", err)
	}

	paths, err := infrastructure.ParseFiles(globPattern)

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

	switch format {
	case "":
		output.OutputToConsole(formatter.Render("plain"))
	case "markdown":
		report := formatter.Render("markdown")
		err := output.OutputToFile(report, "log_report.md")

		if err != nil {
			log.Fatalf("Error saving markdown report: %v", err)
		}

		fmt.Println("Report saved as log_report.md")
	case "adoc":
		report := formatter.Render("adoc")
		err := output.OutputToFile(report, "log_report.adoc")

		if err != nil {
			log.Fatalf("Error saving AsciiDoc report: %v", err)
		}

		fmt.Println("Report saved as log_report.adoc")
	default:
		log.Fatalf("Unsupported format: %s", format)
	}
}

// main is the entry point of the program.
func main() {
	rootCmd = setupRootCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
