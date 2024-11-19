package application

import (
	"fmt"
	"github.com/abakunov/log-analyzer/internal/domain"
	"strconv"
	"strings"
	"time"
)

var unknownFieldWarned = false // Tracks if the warning has been shown.

// matchesFilter checks if a log record matches the filter criteria.
func matchesFilter(logRecord domain.LogRecord, field, value string) bool {
	isWildcard := strings.HasSuffix(value, "*")
	if isWildcard {
		value = strings.TrimSuffix(value, "*")
	}

	switch field {
	case "ip":
		if isWildcard {
			return strings.HasPrefix(logRecord.IP, value)
		}
		return logRecord.IP == value
	case "timestamp":
		filterTime, err := time.Parse(time.RFC3339, value)
		if err != nil {
			fmt.Printf("Invalid timestamp filter value: %v\n", err)
			return false
		}
		return logRecord.Timestamp.Equal(filterTime)
	case "method":
		if isWildcard {
			return strings.HasPrefix(strings.ToLower(logRecord.Method), strings.ToLower(value))
		}
		return strings.EqualFold(logRecord.Method, value)
	case "url":
		if isWildcard {
			return strings.HasPrefix(logRecord.URL, value)
		}
		return logRecord.URL == value
	case "protocol":
		if isWildcard {
			return strings.HasPrefix(strings.ToLower(logRecord.Protocol), strings.ToLower(value))
		}
		return strings.EqualFold(logRecord.Protocol, value)
	case "status":
		if isWildcard {
			return strings.HasPrefix(fmt.Sprintf("%d", logRecord.StatusCode), value)
		}
		return fmt.Sprintf("%d", logRecord.StatusCode) == value
	case "response_size":
		size, err := strconv.Atoi(value)
		if err != nil {
			fmt.Printf("Invalid response_size filter value: %v\n", err)
			return false
		}
		return logRecord.ResponseSize == size
	case "referer":
		if isWildcard {
			return strings.HasPrefix(logRecord.Referer, value)
		}
		return logRecord.Referer == value
	case "agent":
		if isWildcard {
			return strings.HasPrefix(logRecord.UserAgent, value)
		}
		return logRecord.UserAgent == value
	default:
		if !unknownFieldWarned {
			fmt.Printf("Unknown filter field: %s\n", field)
			unknownFieldWarned = true
		}
		return false
	}
}
