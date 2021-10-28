// Package template provides template formating for the output
// of the status command.
package template

import (
	"bytes"
	"text/template"
)

const status = `
{{- range .Pipeline.Workflows }}
{{- .Name }}:
  {{- range .Jobs }}
  {{ .Number | printf "%-2d" }} {{ .Name | printf "%-30s" }} {{ .Status }}
  {{- end }}
{{ end -}}`

var tmpl *template.Template

func init() {
	tmpl, _ = template.New("status").Parse(status)
}

// Render will render the given data using the template.
func Render(data interface{}) string {
	var b bytes.Buffer
	err := tmpl.Execute(&b, data)
	if err != nil {
		panic(err)
	}
	return b.String()
}
