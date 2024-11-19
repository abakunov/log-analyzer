package infrastructure

import (
	"fmt"
	"os"
)

// ReportOutput handles outputting text reports to the console or a file.
type ReportOutput struct{}

// OutputToConsole prints the report text to the console.
func (ro *ReportOutput) OutputToConsole(report string) {
	fmt.Println(report)
}

// OutputToFile writes the report text to a file.
func (ro *ReportOutput) OutputToFile(report, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("failed to close file: %v\n", err)
		}
	}(file)

	_, err = file.WriteString(report)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
