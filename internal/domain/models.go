package domain

import "time"

type LogRecord struct {
	IP           string
	Timestamp    time.Time
	Method       string
	URL          string
	Protocol     string
	StatusCode   int
	ResponseSize int
	Referer      string
	UserAgent    string
}

type Metrics struct {
	FileNames       []string
	StartDate       time.Time
	EndDate         time.Time
	TotalRequests   int
	TotalRespSize   int
	AverageRespSize float64
	Percentile95    int
	ResponseSizes   []int
	Resources       map[string]int
	StatusCodes     map[int]int
}

func NewMetrics(fileNames []string) *Metrics {
	return &Metrics{
		FileNames:     fileNames,
		StartDate:     time.Now(),
		EndDate:       time.Time{},
		ResponseSizes: make([]int, 0),
		Resources:     make(map[string]int),
		StatusCodes:   make(map[int]int),
	}
}

type LogParser interface {
	ParseLogLine(line string) (LogRecord, error)
}

type StreamReader interface {
	ReadLine() (string, error)
}
