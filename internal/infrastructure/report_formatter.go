package infrastructure

import (
	"fmt"
	"github.com/abakunov/log-analyzer/internal/domain"
	"math"
	"sort"
	"strings"
	"time"
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
		sb.WriteString(fmt.Sprintf("#### %s\n\n", header))
	case "adoc":
		sb.WriteString(fmt.Sprintf("= %s\n\n", header))
	default: // plain text
		sb.WriteString(fmt.Sprintf("=== %s ===\n\n", header))
	}
}

// addTable adds a table to the report in the specified format.
func addTable(sb *strings.Builder, format, title string, rows [][]string) {
	switch format {
	case "markdown":
		sb.WriteString(fmt.Sprintf("#### %s\n\n", title))
		for i, row := range rows {
			if i == 0 {
				sb.WriteString("| " + strings.Join(row, " | ") + " |\n")
				sb.WriteString("|" + strings.Repeat(":---|", len(row)) + "\n")
			} else {
				sb.WriteString("| " + strings.Join(row, " | ") + " |\n")
			}
		}
		sb.WriteString("\n")
	case "adoc":
		sb.WriteString(fmt.Sprintf("== %s\n\n", title))
		sb.WriteString("[cols=\"2,1\", options=\"header\"]\n|===\n")
		for _, row := range rows {
			sb.WriteString("| " + strings.Join(row, " | ") + "\n")
		}
		sb.WriteString("|===\n\n")
	default: // plain text
		sb.WriteString(fmt.Sprintf("%s:\n", title))
		for _, row := range rows {
			sb.WriteString(fmt.Sprintf(" %-25s %-15s\n", row[0], row[1]))
		}
		sb.WriteString("\n")
	}
}
