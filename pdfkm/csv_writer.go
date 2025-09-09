package pdfkm

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type CustomCSVWriter struct {
	Writer    *csv.Writer
	Delimiter string
}

func NewCustomCSVWriter(w *csv.Writer, delimiter string) *CustomCSVWriter {
	return &CustomCSVWriter{
		Writer:    w,
		Delimiter: delimiter,
	}
}

func (c *CustomCSVWriter) Write(record []string) error {
	// Join the fields using the custom delimiter
	// joinedRecord := fmt.Sprintf("%s%s%s", record.code, c.Delimiter, record.serial)
	joinedRecord := strings.Join(record, c.Delimiter)

	// Write the joined record to the CSV writer
	return c.Writer.Write([]string{joinedRecord})
}

func NewWriter(w io.Writer) (writer *csv.Writer) {
	writer = csv.NewWriter(w)
	writer.Comma = '\t'
	return
}

func saveCsvCustom(name string, data [][]string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '\t'
	writer.WriteAll(data) // calls Flush internally
	return writer.Error()
}

func saveTxtCustom(name string, data [][]string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, line := range data {
		if len(line) != 2 {
			return fmt.Errorf("строка для записи в файл должна содержать две строки ")
		}
		_, err := file.WriteString(line[0] + "\t" + line[1] + "\n")
		if err != nil {
			return fmt.Errorf("Error writing to file: %v", err)
		}
	}
	return nil
}
