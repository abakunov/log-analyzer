package application

import (
	"github.com/abakunov/log-analyzer/internal/domain"
	"sort"
)

// updateMetrics updates the metrics based on the log record.
func (a *LogAnalyzer) updateMetrics(logRecord domain.LogRecord) {
	metrics := a.Metrics

	metrics.TotalRequests++

	// Update StartDate and EndDate.
	if metrics.StartDate.IsZero() || metrics.StartDate.After(logRecord.Timestamp) {
		metrics.StartDate = logRecord.Timestamp
	}
	if metrics.EndDate.IsZero() || metrics.EndDate.Before(logRecord.Timestamp) {
		metrics.EndDate = logRecord.Timestamp
	}

	metrics.TotalRespSize += logRecord.ResponseSize
	metrics.AverageRespSize = float64(metrics.TotalRespSize) / float64(metrics.TotalRequests)

	metrics.ResponseSizes = append(metrics.ResponseSizes, logRecord.ResponseSize)
	metrics.Percentile95 = a.CalculatePercentile(metrics.ResponseSizes, 95)

	metrics.Resources[logRecord.URL]++
	metrics.StatusCodes[logRecord.StatusCode]++
	metrics.UniqueIPs[logRecord.IP] = struct{}{}
}

// CalculatePercentile calculates the value of the specified percentile.
func (a *LogAnalyzer) CalculatePercentile(values []int, percentile float64) int {
	if len(values) == 0 {
		return 0
	}
	sort.Ints(values)
	index := int(float64(len(values)) * percentile / 100)
	if index >= len(values) {
		index = len(values) - 1
	}
	return values[index]
}
