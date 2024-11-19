package infrastructure

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/abakunov/log-analyzer/internal/domain"
)

// ReportFormatter is responsible for generating text reports.
type ReportFormatter struct {
	Metrics *domain.Metrics
}

// sortMapByValue sorts a map by value in descending order and returns a slice of key-value pairs.
func sortMapByValue(data map[string]int) []struct {
	Key   string
	Value int
} {
	pairs := make([]struct {
		Key   string
		Value int
	}, 0, len(data))

	for key, value := range data {
		pairs = append(pairs, struct {
			Key   string
			Value int
		}{key, value})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	return pairs
}

// sortIntMapByValue sorts a map with integer keys by value in descending order.
func sortIntMapByValue(data map[int]int) []struct {
	Key   int
	Value int
} {
	pairs := make([]struct {
		Key   int
		Value int
	}, 0, len(data))

	for key, value := range data {
		pairs = append(pairs, struct {
			Key   int
			Value int
		}{key, value})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	return pairs
}

// Render generates the report in the specified format.
func (rf *ReportFormatter) Render(format string) string {
	var sb strings.Builder

	// Add report creation timestamp.
	addHeader(&sb, format, fmt.Sprintf("Report created: %s", time.Now().Format("02.01.2006 15:04:05")))

	// Add general information section.
	addTable(&sb, format, "General Information", [][]string{
		{"Files", strings.Join(rf.Metrics.FileNames, ", ")},
		{"Start Date", rf.Metrics.StartDate.Format("02.01.2006")},
		{"End Date", rf.Metrics.EndDate.Format("02.01.2006")},
		{"Total Requests", fmt.Sprintf("%d", rf.Metrics.TotalRequests)},
		{"Unique IPs Count", fmt.Sprintf("%d", len(rf.Metrics.UniqueIPs))},
		{"RPS (Requests/sec)", fmt.Sprintf("%.2f", rf.Metrics.RPS)},
		{"Average Response Size", fmt.Sprintf("%db", int(math.Round(rf.Metrics.AverageRespSize)))},
		{"95th Percentile Size", fmt.Sprintf("%db", rf.Metrics.Percentile95)},
	})

	// Add resources section.
	sortedResources := sortMapByValue(rf.Metrics.Resources)

	resourcesTable := [][]string{{"Resource", "Count"}}
	for _, res := range sortedResources {
		resourcesTable = append(resourcesTable, []string{res.Key, fmt.Sprintf("%d", res.Value)})
	}

	addTable(&sb, format, "Requested Resources", resourcesTable)

	// Add status codes section.
	sortedStatusCodes := sortIntMapByValue(rf.Metrics.StatusCodes)

	statusTable := [][]string{{"Code", "Count"}}
	for _, code := range sortedStatusCodes {
		statusTable = append(statusTable, []string{fmt.Sprintf("%d", code.Key), fmt.Sprintf("%d", code.Value)})
	}

	addTable(&sb, format, "Response Codes", statusTable)

	return sb.String()
}

// addHeader adds a section header to the report in the specified format.
func addHeader(sb *strings.Builder, format, header string) {
	switch format {
	case "markdown":
		fmt.Fprintf(sb, "#### %s\n\n", header)
	case "adoc":
		fmt.Fprintf(sb, "= %s\n\n", header)
	default: // plain text
		fmt.Fprintf(sb, "=== %s ===\n\n", header)
	}
}

// addTable adds a table to the report in the specified format.
func addTable(sb *strings.Builder, format, title string, rows [][]string) {
	switch format {
	case "markdown":
		fmt.Fprintf(sb, "#### %s\n\n", title)

		for i, row := range rows {
			if i == 0 {
				fmt.Fprintf(sb, "| %s |\n", strings.Join(row, " | "))
				fmt.Fprintf(sb, "|%s|\n", strings.Repeat(":---|", len(row)))
			} else {
				fmt.Fprintf(sb, "| %s |\n", strings.Join(row, " | "))
			}
		}

	case "adoc":
		fmt.Fprintf(sb, "== %s\n\n", title)
		fmt.Fprintf(sb, "[cols=\"2,1\", options=\"header\"]\n|===\n")

		for _, row := range rows {
			fmt.Fprintf(sb, "| %s\n", strings.Join(row, " | "))
		}

		fmt.Fprint(sb, "|===\n")
	default: // plain text
		fmt.Fprintf(sb, "%s:\n", title)

		for _, row := range rows {
			fmt.Fprintf(sb, " %-25s %-15s\n", row[0], row[1])
		}
	}

	fmt.Fprintln(sb)
}
