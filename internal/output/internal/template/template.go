// Package template provides template formating for the output
// of the output command.
package template

import (
	"bytes"
	"html/template"
)

const output = `
{{- range .Steps }}
-- {{ .Name | printf "%-50s" }} -------------------------- 

{{ .Output }}
{{ end -}}
`

var tmpl *template.Template

func init() {
	tmpl, _ = template.New("output").Parse(output)
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
