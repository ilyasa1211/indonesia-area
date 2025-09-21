package writer

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
)

type CSVWriter struct {
	headers   []string
	delimiter rune
	extension string
}

func NewCSVWriter(delimiter rune, headers []string) *CSVWriter {
	return &CSVWriter{
		headers:   headers,
		delimiter: delimiter,
		extension: ".csv",
	}
}

func (w *CSVWriter) SetHeaders(headers []string) {
	w.headers = headers
}

func (w *CSVWriter) Write(file *os.File, data any) error {
	wr := csv.NewWriter(file)

	defer wr.Flush()

	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice, got %T", data)
	}

	wr.Comma = w.delimiter

	if w.headers != nil {
		err := wr.Write(w.headers)
		if err != nil {
			return err
		}
	}

	v := reflect.ValueOf(data)

	for i := 0; i < v.Len(); i++ {
		var row []string
		elem := v.Index(i)

		for i := 0; i < elem.NumField(); i++ {
			row = append(row, fmt.Sprintf("%v", elem.Field(i).Interface()))
		}

		if err := wr.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}

	}

	return nil
}

func (w *CSVWriter) GetExtension() string {
	return w.extension
}
