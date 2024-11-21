package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) (string, error) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()

	if err != nil {
		return "", err
	}

	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	outputChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outputChan <- buf.String()
	}()

	f()

	_ = w.Close()

	return <-outputChan, nil
}

func TestRunAnalyzer_Success(t *testing.T) {
	// Создайте тестовые данные.
	err := os.MkdirAll("testdata", os.ModePerm)
	assert.NoError(t, err, "Failed to create testdata directory.")
	defer os.RemoveAll("testdata")

	logData := `127.0.0.1 - - [12/Dec/2021:15:04:05 +0000] "GET /index.html HTTP/1.1" 200 1024 "http://example.com" "Mozilla/5.0"`
	err = os.WriteFile("testdata/logfile.log", []byte(logData), 0o600)
	assert.NoError(t, err, "Failed to create logfile.log.")

	// Mock глобальные переменные.
	globPattern = "testdata/logfile.log"
	from = ""
	to = ""
	format = ""
	filterField = ""
	filterValue = ""

	// Capture the output of the analyzer.
	output, err := captureOutput(func() { runAnalyzer() })
	assert.NoError(t, err, "Failed to capture output.")
	assert.Contains(t, output, "Total Requests", "Output should contain 'Total Requests'.")
}
