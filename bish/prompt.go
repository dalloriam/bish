package bish

import (
	"bytes"
	"strings"
	"text/template"
)

type promptContext struct {
	Username string
	Hostname string
	Cwd      string
}

func (p *promptContext) Render(temp string) ([]string, error) {
	ren, err := template.New("hello").Parse(temp)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := ren.Execute(&buf, p); err != nil {
		return nil, err
	}
	return strings.Split(buf.String(), "\n"), nil
}
