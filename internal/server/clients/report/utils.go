package report

import (
	"bytes"
	"text/template"
)

func FormatRecordInCSV(csvTemplate *template.Template, record Record) ([]byte, error) {
	var buf bytes.Buffer

	err := csvTemplate.Execute(&buf, record)

	return buf.Bytes(), err
}
