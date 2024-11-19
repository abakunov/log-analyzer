package application

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/abakunov/log-analyzer/internal/domain"
)

func ParseLogLine(line string) (domain.LogRecord, error) {
	var log domain.LogRecord
	pattern := `^(\S+) - - \[([^\]]+)\] "(\S+) (\S+) (\S+)" (\d+) (\d+|-) "([^"]*)" "([^"]*)"$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if len(matches) != 10 {
		return log, fmt.Errorf("failed to parse line: %s", line)
	}

	log.IP = matches[1]

	// Parse timestamp.
	timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
		return log, fmt.Errorf("failed to parse time: %v", err)
	}

	log.Timestamp = timestamp

	log.Method = matches[3]
	log.URL = matches[4]
	log.Protocol = matches[5]

	// Status code.
	statusCode, err := strconv.Atoi(matches[6])
	if err != nil {
		return log, fmt.Errorf("failed to parse status code: %v", err)
	}

	log.StatusCode = statusCode

	// Size of response.
	if matches[7] != "-" {
		responseSize, err := strconv.Atoi(matches[7])
		if err != nil {
			return log, fmt.Errorf("failed to parse response size: %v", err)
		}
		log.ResponseSize = responseSize
	} else {
		log.ResponseSize = 0
	}

	log.Referer = matches[8]
	log.UserAgent = matches[9]

	return log, nil
}
