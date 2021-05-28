package grpcwithgw

import (
	"embed"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/gobuffalo/flect"
	log "github.com/sirupsen/logrus"

	"github.com/cnative/servicebuilder/internal/builder"
)

type grpcServiceTemplateProvider struct {
	options   *builder.Options
	templates map[string]*template.Template
}

//go:embed tmplt/**
var assets embed.FS

var funcs = template.FuncMap{
	"TitleCase": strings.Title,
	"LowerCase": strings.ToLower,
	"UpperCase": strings.ToUpper,
	"Trim":      strings.Trim,
	"Pluralize": flect.Pluralize,
	"LCPluralize": func(s string) string {
		return flect.Pluralize(strings.ToLower(s))
	},
}

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

func walkAssets(fn func(path string, content []byte) error) error {
	rootDir := "tmplt"
	dirs, dir := []string{rootDir}, ""
	for len(dirs) > 0 {
		dir, dirs = dirs[0], dirs[1:]
		de, err := assets.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, d := range de {
			fqn := path.Join(dir, d.Name())
			if d.IsDir() {
				dirs = append(dirs, fqn)
			} else {
				content, err := assets.ReadFile(fqn)
				if err != nil {
					return err
				}
				if err := fn(strings.TrimPrefix(fqn, fmt.Sprintf("%s/", rootDir)), content); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (g *grpcServiceTemplateProvider) initialize() error {

	g.templates = make(map[string]*template.Template)

	collector := func(pth string, content []byte) error {

		fb := new(strings.Builder)
		t := template.Must(template.New("file name").Funcs(funcs).Parse(strings.TrimSuffix(pth, ".tmplt")))
		if err := t.Execute(fb, g.options); err != nil {
			return err
		}
		fpath := fb.String()
		tmplt := template.New(fpath).Funcs(funcs)
		if strings.HasPrefix(fpath, "helm") {
			if g.options.DeploymentType != builder.HemlChart {
				return nil //ignore required deployment is not helm
			}
			fpath = fmt.Sprintf("deployments%s", strings.TrimPrefix(fpath, "helm"))
			tmplt = template.New(fpath).Funcs(funcs).Delims("[[", "]]")
		} else if strings.HasPrefix(fpath, "kustomize") {
			if g.options.DeploymentType != builder.K8SManifest {
				return nil //ignore required deployment is not k8s
			}
			fpath = fmt.Sprintf("deployments%s", strings.TrimPrefix(fpath, "kustomize"))
		}

		pt, err := tmplt.Parse(string(content))
		if err != nil {
			return err
		}
		g.templates[fpath] = pt

		return nil
	}

	if err := walkAssets(collector); err != nil {
		return err
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
