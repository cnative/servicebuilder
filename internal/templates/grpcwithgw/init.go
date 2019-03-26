package grpcwithgw

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/cnative/servicebuilder/internal/builder"
	log "github.com/sirupsen/logrus"
)

//go:generate go-bindata -o ./grpc_sevice_with_gw.go -pkg grpcwithgw -nomemcopy -nometadata -prefix tmplt tmplt/...

type (
	grpcServiceTemplateProvider struct {
		options   *builder.Options
		templates map[string]*template.Template
	}
)

var (
	funcs = template.FuncMap{
		"TitleCase": strings.Title,
		"Trim":      strings.Trim,
	}
)

// New creates GRPC Service Builder with Gateway Builder
func New(o *builder.Options) (builder.TemplateProvider, error) {

	s := &grpcServiceTemplateProvider{
		options: o,
	}
	if err := s.initialize(); err != nil {
		return nil, err
	}

	return s, nil
}

func (g *grpcServiceTemplateProvider) initialize() error {

	g.templates = make(map[string]*template.Template)

	for k, v := range _bindata {
		t, err := v()
		if err != nil {
			return err
		}
		s := string(t.bytes)

		pt, err := template.New(k).Funcs(funcs).Parse(s)
		if err != nil {
			return err
		}
		f := strings.TrimSuffix(k, ".tmplt")
		if f == "proto" {
			f = fmt.Sprintf("%s.proto", g.options.Name)
		}
		g.templates[f] = pt
	}

	log.Info("GRPC service with Gateway template provider initialized")
	return nil
}

func (g *grpcServiceTemplateProvider) GetTemplates() map[string]*template.Template {
	return g.templates
}

func (g *grpcServiceTemplateProvider) GetOptions() *builder.Options {
	return g.options
}
