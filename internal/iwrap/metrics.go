package iwrap

// MetricsTmplt used to wrap an interface with metrics generators
const MetricsTmplt = `

{{$target := .InterfaceName}}

package {{ .PackageName }}
import (
	"context"

	"go.opencensus.io/stats"

	{{range .CustomImports}}
	"{{.}}"
	{{end}}
)

// force context to be used
var _ context.Context

// {{ lowerCase $target }}WithMetrics wraps another {{ $target }} and sends metrics to Prometheus.
type {{ lowerCase $target }}WithMetrics struct {
	wrapped{{$target}}    {{ $target }}
	observer *{{ lowerCase $target }}Observer
}

// {{ $target }}WithMetrics creates a new {{ $target }} with metrics
func {{ $target }}WithMetrics(toWrap {{ $target }}, logger *log.Logger) {{ $target }} {
	return &{{ lowerCase $target }}WithMetrics{wrapped{{$target}}: toWrap, observer: new{{ $target }}Observer(logger)}
}

var _ {{ $target }} = (*{{ lowerCase $target }}WithMetrics)(nil)

{{range .Methods}}
{{template "doc" . -}}
func (s *{{ lowerCase $target }}WithMetrics) {{.Name}}({{template "list" .Params}}) ({{template "list" .Returns}}) {
	done := s.observer.Observe(ctx, "{{.Name}}")
	defer done()
	{{template "call" .Returns}} = s.wrapped{{$target}}.{{.Name}}({{template "call" .Params}})
	{{if isLastReturnError .Returns }}
		if {{ lastReturnName .Returns }} != nil {
			stats.Record(ctx, {{ lowerCase $target }}CallErrorCount.M(1)) // Counter to track a wrapped{{$target}} call errors
		}
	{{end}}

	return {{template "call" .Returns}}
}
{{end}}

// Healthy calls Healthy on the wrapped wrapped{{$target}}.
func (s *{{ lowerCase $target }}WithMetrics) Healthy() error {
	return s.wrapped{{$target}}.Healthy()
}

// Ready calls ready on the wrapped wrapped{{$target}}.
func (s *{{ lowerCase $target }}WithMetrics) Ready() (bool, error) {
	return s.wrapped{{$target}}.Ready()
}

// Close calls Close on the wrapped wrapped{{$target}}.
func (s *{{ lowerCase $target }}WithMetrics) Close() error {
	return s.wrapped{{$target}}.Close()
}

{{define "list"}}{{range $index, $element := .}}{{if $index}}, {{end}}{{if $element.Name}}{{$element.Name}}{{end}} {{$element.Type}}{{end}}{{end}}
{{define "call"}}{{range $index, $element := .}}{{if $index}}, {{end}}{{if $element.Name}}{{$element.Name}}{{end}}{{end}}{{end}}
{{define "doc"}}
{{range .Doc}}
{{.}}
{{- else}}
// {{.Name}} .
{{- end}}
{{end}}
`
