package exporter

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	FilenameTemplate = "flag-evaluation-{{ .Timestamp }}.{{ .Format }}"
	CsvTemplate      = "{{ .CreationDate }};{{ .Key }};{{ .Value }};{{ .Variant }};{{ .Reason }};{{ .ErrorCode }};{{ .ErrorMessage }}\n"
)

type Parser struct {
	FilenameTemplate *template.Template
	CsvTemplate      *template.Template
}

func (p *Parser) ParseTemplate(name, temp string) *template.Template {
	t, _ := template.New(name).Parse(temp)
	return t
}

func (p *Parser) ParseFilename(format string) (string, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	format = strings.ToLower(format)

	var buf bytes.Buffer

	err := p.FilenameTemplate.Execute(&buf, struct {
		Timestamp string
		Format    string
	}{
		Timestamp: timestamp,
		Format:    format,
	})

	return buf.String(), err
}
