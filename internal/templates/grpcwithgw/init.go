package grpcwithgw

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/cnative/servicebuilder/internal/builder"
	log "github.com/sirupsen/logrus"
)

//go:generate go-bindata -o ./grpc_service_with_gw.go -pkg grpcwithgw -nomemcopy -nometadata -prefix tmplt tmplt/...

type (
	grpcServiceTemplateProvider struct {
		options   *builder.Options
		templates map[string]*template.Template
	}
)

var (
	funcs = template.FuncMap{
		"TitleCase": strings.Title,
		"LowerCase": strings.ToLower,
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
		f := strings.TrimSuffix(k, ".tmplt")
		tmplt := template.New(f).Funcs(funcs)

		if strings.HasPrefix(k, "helm") {
			if g.options.DeploymentType != builder.HemlChart {
				continue //ignore required deployment is not helm
			}
			f = fmt.Sprintf("deployment%s", strings.TrimPrefix(k, "helm"))
			tmplt = template.New(f).Funcs(funcs).Delims("[[", "]]")
		} else if strings.HasPrefix(k, "kustomize") {
			if g.options.DeploymentType != builder.K8SManifest {
				continue //ignore required deployment is not k8s
			}
			f = fmt.Sprintf("deployment%s", strings.TrimPrefix(k, "kustomize"))
		}

		pt, err := tmplt.Parse(s)
		if err != nil {
			return err
		}
		if f == "proto" {
			f = fmt.Sprintf("%s.proto", g.options.Name)
		}
		g.templates[f] = pt
	}

	log.Info("gRPC service with Gateway template provider initialized")
	return nil
}

func (g *grpcServiceTemplateProvider) GetTemplates() map[string]*template.Template {
	return g.templates
}

func (g *grpcServiceTemplateProvider) GetOptions() *builder.Options {
	return g.options
}
