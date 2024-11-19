package application

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/abakunov/log-analyzer/internal/domain"
)

var unknownFieldWarned = false // Tracks if the warning has been shown.

// matchesFilter checks if a log record matches the filter criteria.
func matchesFilter(logRecord *domain.LogRecord, field, value string) bool {
	isWildcard := strings.HasSuffix(value, "*")
	if isWildcard {
		value = strings.TrimSuffix(value, "*")
	}

	switch field {
	case "ip":
		return matchIP(logRecord.IP, value, isWildcard)
	case "timestamp":
		return matchTimestamp(logRecord.Timestamp, value)
	case "method":
		return matchStringField(logRecord.Method, value, isWildcard)
	case "url":
		return matchStringField(logRecord.URL, value, isWildcard)
	case "protocol":
		return matchStringField(strings.ToLower(logRecord.Protocol), strings.ToLower(value), isWildcard)
	case "status":
		return matchStatusCode(logRecord.StatusCode, value, isWildcard)
	case "response_size":
		return matchResponseSize(logRecord.ResponseSize, value)
	case "referer":
		return matchStringField(logRecord.Referer, value, isWildcard)
	case "agent":
		return matchStringField(logRecord.UserAgent, value, isWildcard)
	default:
		if !unknownFieldWarned {
			fmt.Printf("Unknown filter field: %s\n", field)

			unknownFieldWarned = true
		}

		return false
	}
}

// matchIP checks if an IP matches the given filter value.
func matchIP(ip, value string, isWildcard bool) bool {
	if isWildcard {
		return strings.HasPrefix(ip, value)
	}

	return ip == value
}

// matchTimestamp checks if a timestamp matches the given filter value.
func matchTimestamp(timestamp time.Time, value string) bool {
	filterTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		fmt.Printf("Invalid timestamp filter value: %v\n", err)
		return false
	}

	return timestamp.Equal(filterTime)
}

// matchStringField checks if a string field matches the given filter value.
func matchStringField(field, value string, isWildcard bool) bool {
	if isWildcard {
		return strings.HasPrefix(field, value)
	}

	return field == value
}

// matchStatusCode checks if a status code matches the given filter value.
func matchStatusCode(statusCode int, value string, isWildcard bool) bool {
	statusStr := fmt.Sprintf("%d", statusCode)
	if isWildcard {
		return strings.HasPrefix(statusStr, value)
	}

	return statusStr == value
}

// matchResponseSize checks if a response size matches the given filter value.
func matchResponseSize(responseSize int, value string) bool {
	size, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("Invalid response_size filter value: %v\n", err)
		return false
	}

	return responseSize == size
}
