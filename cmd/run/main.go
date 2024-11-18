package main

import (
	"fmt"
	"github.com/abakunov/log-analyzer/internal/application"
	"github.com/abakunov/log-analyzer/internal/infrastructure"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	paths  []string
	from   string
	to     string
	format string
)

func parseTimeBounds(fromStr, toStr string) (time.Time, time.Time, error) {
	var fromTime, toTime time.Time
	var err error

	if fromStr != "" {
		fromTime, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid from time: %w", err)
		}
	}

	if toStr != "" {
		toTime, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid to time: %w", err)
		}
	}

	return fromTime, toTime, nil
}

func runAnalyzer() {
	fromTime, toTime, err := parseTimeBounds(from, to)
	if err != nil {
		log.Fatalf("Error parsing time bounds: %v", err)
	}

	analyzer := application.NewLogAnalyzer(paths)
	err = analyzer.AnalyzeLogs(fromTime, toTime)
	if err != nil {
		log.Fatalf("Error analyzing logs: %v", err)
	}

	metrics := analyzer.Metrics
	formatter := infrastructure.ReportFormatter{Metrics: metrics}

	// Generate the report in the specified format
	var report string
	outputFile := "log_report.md"
	if format == "markdown" || format == "" {
		report = formatter.RenderMarkdown()
	} else if format == "adoc" {
		report = formatter.RenderAdoc()
		outputFile = "log_report.adoc"
	} else {
		log.Fatalf("Unsupported format: %s", format)
	}

	output := infrastructure.ReportOutput{}
	err = output.OutputToFile(report, outputFile)
	if err != nil {
		log.Fatalf("Error saving report: %v", err)
	}

	fmt.Printf("Report saved as %s\n", outputFile)
}

var rootCmd = &cobra.Command{
	Use:   "analyzer",
	Short: "Analyze NGINX log files",
	Run: func(cmd *cobra.Command, args []string) {
		runAnalyzer()
	},
}

func init() {
	rootCmd.Flags().StringSliceVar(&paths, "path", nil, "Path(s) to log files (required)")
	rootCmd.Flags().StringVar(&from, "from", "", "Start date in ISO8601 format (optional)")
	rootCmd.Flags().StringVar(&to, "to", "", "End date in ISO8601 format (optional)")
	rootCmd.Flags().StringVar(&format, "format", "", "Output format: markdown or adoc (optional)")
	rootCmd.MarkFlagRequired("path")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
