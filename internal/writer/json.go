package writer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type JSONWriter struct {
	Indent    int
	extension string
}

func NewJSONWriter(indent int) *JSONWriter {
	if indent < 0 {
		fmt.Println("Indent cannot be negative, setting to 0")
		indent = 0
	}

	return &JSONWriter{
		Indent:    indent,
		extension: ".json",
	}
}

func (w *JSONWriter) Write(file *os.File, data any) error {
	enc := json.NewEncoder(file)
	enc.SetIndent("", strings.Repeat(" ", w.Indent))
	err := enc.Encode(data)

	if err != nil {
		return err
	}

	return nil
}

func (w *JSONWriter) GetExtension() string {
	return w.extension
}
