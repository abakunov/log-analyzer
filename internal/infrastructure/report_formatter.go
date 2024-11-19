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

// RenderMarkdown generates a markdown representation of the report.
func (rf *ReportFormatter) RenderMarkdown() string {
	var sb strings.Builder

	// Report creation timestamp.
	sb.WriteString(fmt.Sprintf("#### Report created: %s\n\n", time.Now().Format("02.01.2006 15:04:05")))

	// General information.
	sb.WriteString("#### General Information\n\n")
	sb.WriteString("| Metric             | Value          |\n")
	sb.WriteString("|:-------------------|:----------------|\n")                                     // Left-aligned table.
	sb.WriteString(fmt.Sprintf("| File(s)           | `%s`           |\n", rf.Metrics.FileNames[0])) // First file.
	for _, file := range rf.Metrics.FileNames[1:] {
		sb.WriteString(fmt.Sprintf("|                   | `%s`           |\n", file)) // Remaining files.
	}
	sb.WriteString(fmt.Sprintf("| Start Date        | %s             |\n", rf.Metrics.StartDate.Format("02.01.2006")))
	sb.WriteString(fmt.Sprintf("| End Date          | %s             |\n", rf.Metrics.EndDate.Format("02.01.2006")))
	sb.WriteString(fmt.Sprintf("| Total Requests    | %d             |\n", rf.Metrics.TotalRequests))
	sb.WriteString(fmt.Sprintf("| Unique IPs Count  | %d             |\n", len(rf.Metrics.UniqueIPs)))
	sb.WriteString(fmt.Sprintf("| RPS (Requests/sec)| %.2f           |\n", rf.Metrics.RPS))
	sb.WriteString(fmt.Sprintf("| Average Resp. Size| %db            |\n", int(math.Round(rf.Metrics.AverageRespSize))))
	sb.WriteString(fmt.Sprintf("| 95th Percentile   | %db            |\n", rf.Metrics.Percentile95))

	// Resources section.
	sb.WriteString("\n#### Requested Resources\n\n")
	sb.WriteString("| Resource          | Count          |\n")
	sb.WriteString("|:------------------|---------------:|\n")
	sortedResources := sortMapByValue(rf.Metrics.Resources)
	for _, pair := range sortedResources {
		sb.WriteString(fmt.Sprintf("| `%s`             | %d             |\n", pair.Key, pair.Value))
	}

	// Status codes section.
	sb.WriteString("\n#### Response Codes\n\n")
	sb.WriteString("| Code              | Count          |\n")
	sb.WriteString("|:------------------|---------------:|\n")
	sortedStatusCodes := sortIntMapByValue(rf.Metrics.StatusCodes)
	for _, pair := range sortedStatusCodes {
		sb.WriteString(fmt.Sprintf("| %d               | %d             |\n", pair.Key, pair.Value))
	}

	return sb.String()
}

// RenderAdoc generates an AsciiDoc representation of the report.
func (rf *ReportFormatter) RenderAdoc() string {
	var sb strings.Builder

	// Report creation timestamp.
	sb.WriteString(fmt.Sprintf("= Report created: %s\n\n", time.Now().Format("02.01.2006 15:04:05")))

	// General information.
	sb.WriteString("== General Information\n\n")
	sb.WriteString("[cols=\"2,1\", options=\"header\"]\n|===\n")
	sb.WriteString("| Metric | Value\n")
	sb.WriteString(fmt.Sprintf("| File(s) | `%s`\n", rf.Metrics.FileNames[0])) // First file.
	for _, file := range rf.Metrics.FileNames[1:] {
		sb.WriteString(fmt.Sprintf("|         | `%s`\n", file)) // Remaining files.
	}
	sb.WriteString(fmt.Sprintf("| Start Date | %s\n", rf.Metrics.StartDate.Format("02.01.2006")))
	sb.WriteString(fmt.Sprintf("| End Date   | %s\n", rf.Metrics.EndDate.Format("02.01.2006")))
	sb.WriteString(fmt.Sprintf("| Total Requests | %d\n", rf.Metrics.TotalRequests))
	sb.WriteString(fmt.Sprintf("| Unique IPs Count | %d\n", len(rf.Metrics.UniqueIPs)))
	sb.WriteString(fmt.Sprintf("| RPS (Requests/sec) | %.2f\n", rf.Metrics.RPS))
	sb.WriteString(fmt.Sprintf("| Average Resp. Size | %db\n", int(math.Round(rf.Metrics.AverageRespSize))))
	sb.WriteString(fmt.Sprintf("| 95th Percentile | %db\n", rf.Metrics.Percentile95))
	sb.WriteString("|===\n\n")

	// Resources section.
	sb.WriteString("== Requested Resources\n\n")
	sb.WriteString("[cols=\"2,1\", options=\"header\"]\n|===\n")
	sb.WriteString("| Resource | Count\n")
	sortedResources := sortMapByValue(rf.Metrics.Resources)
	for _, pair := range sortedResources {
		sb.WriteString(fmt.Sprintf("| `%s` | %d\n", pair.Key, pair.Value))
	}
	sb.WriteString("|===\n\n")

	// Status codes section.
	sb.WriteString("== Response Codes\n\n")
	sb.WriteString("[cols=\"2,1\", options=\"header\"]\n|===\n")
	sb.WriteString("| Code | Count\n")
	sortedStatusCodes := sortIntMapByValue(rf.Metrics.StatusCodes)
	for _, pair := range sortedStatusCodes {
		sb.WriteString(fmt.Sprintf("| %d | %d\n", pair.Key, pair.Value))
	}
	sb.WriteString("|===\n\n")

	return sb.String()
}

// RenderConsole generates a plain text representation of the report for console output.
func (rf *ReportFormatter) RenderConsole() string {
	var sb strings.Builder

	// Report creation timestamp.
	sb.WriteString(fmt.Sprintf("=== Log Analysis Report (Created: %s) ===\n\n", time.Now().Format("02.01.2006 15:04:05")))

	// General information section.
	sb.WriteString("General Information:\n")
	sb.WriteString(fmt.Sprintf(" %-25s %-15s\n", "Metric", "Value"))
	sb.WriteString(fmt.Sprintf(" %-25s %-15s\n", "-------------------------", "---------------"))
	sb.WriteString(fmt.Sprintf(" %-25s %-15s\n", "Files", rf.Metrics.FileNames[0])) // First file.
	for _, file := range rf.Metrics.FileNames[1:] {
		sb.WriteString(fmt.Sprintf(" %-25s %-15s\n", "", file)) // Remaining files.
	}
	sb.WriteString(fmt.Sprintf(" %-25s %-15s\n", "Start Date", rf.Metrics.StartDate.Format("02.01.2006")))
	sb.WriteString(fmt.Sprintf(" %-25s %-15s\n", "End Date", rf.Metrics.EndDate.Format("02.01.2006")))
	sb.WriteString(fmt.Sprintf(" %-25s %-15d\n", "Total Requests", rf.Metrics.TotalRequests))
	sb.WriteString(fmt.Sprintf(" %-25s %-15d\n", "Unique IPs", len(rf.Metrics.UniqueIPs)))
	sb.WriteString(fmt.Sprintf(" %-25s %-15.2f\n", "RPS (Requests/sec)", rf.Metrics.RPS))
	sb.WriteString(fmt.Sprintf(" %-25s %-15d\n", "Average Resp. Size", int(math.Round(rf.Metrics.AverageRespSize))))
	sb.WriteString(fmt.Sprintf(" %-25s %-15d\n", "95th Percentile Size", rf.Metrics.Percentile95))

	// Resources section.
	sb.WriteString("\nRequested Resources:\n")
	sb.WriteString(fmt.Sprintf(" %-40s %-10s\n", "Resource", "Count"))
	sb.WriteString(fmt.Sprintf(" %-40s %-10s\n", "----------------------------------------", "----------"))
	sortedResources := sortMapByValue(rf.Metrics.Resources)
	for _, res := range sortedResources {
		sb.WriteString(fmt.Sprintf(" %-40s %-10d\n", res.Key, res.Value))
	}

	// Status codes section.
	sb.WriteString("\nResponse Codes:\n")
	sb.WriteString(fmt.Sprintf(" %-10s %-10s\n", "Code", "Count"))
	sb.WriteString(fmt.Sprintf(" %-10s %-10s\n", "----------", "----------"))
	sortedStatusCodes := sortIntMapByValue(rf.Metrics.StatusCodes)
	for _, code := range sortedStatusCodes {
		sb.WriteString(fmt.Sprintf(" %-10d %-10d\n", code.Key, code.Value))
	}

	return sb.String()
}
