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

// Metrics stores statistics from analyzed logs.
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
	UniqueIPs       map[string]struct{} // To track unique IPs
	RPS             float64             // Requests Per Second
}

// NewMetrics initializes a new Metrics instance.
func NewMetrics(fileNames []string) *Metrics {
	return &Metrics{
		FileNames:     fileNames,
		Resources:     make(map[string]int),
		StatusCodes:   make(map[int]int),
		ResponseSizes: make([]int, 0),
		UniqueIPs:     make(map[string]struct{}),
	}
}

type LogParser interface {
	ParseLogLine(line string) (LogRecord, error)
}

type StreamReader interface {
	ReadLine() (string, error)
}
