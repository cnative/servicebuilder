package iwrap

// TracingTmplt used to wrap an interface with opencensus tracing
const TracingTmplt = `
{{$target := .InterfaceName}}

package {{ .PackageName }}

import (
	"context"
	"fmt"
	"strings"

	"go.opencensus.io/trace"

	{{range .CustomImports}}
	"{{.}}"
	{{end}}
)

// force context to be used
var _ context.Context

// {{ lowerCase $target}}WithTrace wraps another {{$target}} and records trace information.
type {{ lowerCase $target}}WithTrace struct {
	wrapped{{$target}}     {{$target}}
	component string
}

// {{$target}}WithTrace creates a new {{$target}} with trace.
func {{$target}}WithTrace(toWrap  {{$target}}, logger *log.Logger) {{$target}} {
	component := strings.TrimPrefix(fmt.Sprintf("%T", toWrap), "*")
	logger.Debugf("store tracing enabled for %v", component)
	
	return &{{ lowerCase $target}}WithTrace{
		wrapped{{$target}} :  toWrap,
		component: component,
	}
}

var _ {{$target}} = (*{{ lowerCase $target}}WithTrace)(nil)

{{range .Methods}}
{{template "doc" . -}}
func (s *{{ lowerCase $target}}WithTrace) {{.Name}}({{template "list" .Params}}) ({{template "list" .Returns}}) {
	ctx, span := trace.StartSpan(ctx, "{{.Name}}")
	defer span.End()

	{{template "call" .Returns}} = s.wrapped{{$target}}.{{.Name}}({{template "call" .Params}})
	{{if isLastReturnError .Returns }}
		if {{ lastReturnName .Returns }} != nil {
			span.Annotate([]trace.Attribute{
				trace.StringAttribute("error", {{ lastReturnName .Returns }}.Error()),
			}, "{{.Name}}")
		}
	{{end}}

	return {{template "call" .Returns}}
}
{{end}}


// Healthy calls Healthy on the wrapped {{$target}}.
func (s *{{ lowerCase $target}}WithTrace) Healthy() error {
	return s.wrapped{{$target}}.Healthy()
}

// Ready calls Ready on the wrapped {{$target}}.
func (s *{{ lowerCase $target}}WithTrace) Ready() (bool, error) {
	return s.wrapped{{$target}}.Ready()
}

// Close calls Close on the wrapped {{$target}}.
func (s *{{ lowerCase $target}}WithTrace) Close() error {
	return s.wrapped{{$target}}.Close()
}

{{define "list"}}{{range $index, $element := .}}{{if $index}}, {{end}}{{if $element.Name}}{{$element.Name}}{{end}} {{$element.Type}}{{end}}{{end}}
{{define "call"}}{{range $index, $element := .}}{{if $index}}, {{end}}{{if $element.Name}}{{$element.Name}}{{end}}{{end}}{{end}}

{{define "error"}}{{end}}

{{define "doc"}}
{{range .Doc}}
{{.}}
{{- else}}
// {{.Name}} .
{{- end}}
{{end}}
`
