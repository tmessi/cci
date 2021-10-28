// Package template provides template formating for the output
// of the status command.
package template

import (
	"bytes"
	"fmt"
	"math"
	"text/template"
	"time"
)

func since(t *time.Time) time.Duration {
	return time.Since(*t)
}

func duration(duration time.Duration) string {
	puralize := func(a int64, singular string) string {
		if a > 1 {
			return fmt.Sprintf("%d %ss", a, singular)
		}
		return fmt.Sprintf("%d %s", a, singular)
	}

	if days := int64(duration.Hours() / 24); days > 0 {
		return puralize(days, "day")
	}

	if hours := int64(math.Mod(duration.Hours(), 24)); hours > 0 {
		return puralize(hours, "hour")
	}

	if minutes := int64(math.Mod(duration.Minutes(), 60)); minutes > 0 {
		return puralize(minutes, "minute")
	}

	if seconds := int64(math.Mod(duration.Seconds(), 60)); seconds > 0 {
		return puralize(seconds, "second")
	}

	return ""
}

const status = `
{{- range .Pipelines }}
{{ .Number }} pipeline: {{ .Updated | since | duration }} ago
  {{- range .Workflows }}
  {{ .Name }}:
    {{- range .Jobs }}
    {{ .Number | printf "%-2d" }} {{ .Name | printf "%-30s" }} {{ .Status }}
    {{- end }}
  {{- end -}}
{{- end -}}`

var tmpl *template.Template

func init() {
	funcMap := template.FuncMap{
		"since":    since,
		"duration": duration,
	}
	tmpl, _ = template.New("status").Funcs(funcMap).Parse(status)
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
