package builder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	// K8SManifest used to generate K8S Manifest that is compatible with kustomize
	K8SManifest DeploymentType = iota - 1

	// HemlChart Deployment artifact
	HemlChart

	// UnknownDeployemntType is invalid deployment type
	UnknownDeployemntType
)

type (

	// DeploymentType indicates artifacts to use for deployment
	DeploymentType int8

	// Options used for Service builder
	Options struct {
		Name            string
		ModuleName      string
		ResourceName    string
		ImageName       string
		Description     string
		DstDir          string
		HTTPRoutePrefix string
		DeploymentType  DeploymentType
		DomainName      string

		ProtocVersion         string
		ServiceBuilderVersion string
	}

	// ServiceBuilder that register templates can generates a service
	ServiceBuilder interface {
		Generate() error
	}

	// TemplateProvider for service builder
	TemplateProvider interface {
		GetOptions() *Options
		GetTemplates() map[string]*template.Template
	}

	serviceBuilder struct {
		tmpDirPath       string
		templateProvider TemplateProvider
	}
)

//New ServiceBuilder with a given template provider
func New(templateProvider TemplateProvider) (ServiceBuilder, error) {
	tempDirPath, err := ioutil.TempDir("", "servicebuilder")
	if err != nil {
		return nil, err
	}
	log.WithField("dir", tempDirPath).Debugf("temp folder created")

	return &serviceBuilder{
		tmpDirPath:       tempDirPath,
		templateProvider: templateProvider,
	}, nil
}

func (g *serviceBuilder) Generate() error {

	if g.templateProvider == nil {
		return errors.New("builder not initialized")
	}

	tmplts := g.templateProvider.GetTemplates()
	options := g.templateProvider.GetOptions()

	for k, v := range tmplts {
		p := path.Join(g.tmpDirPath, path.Dir(k))
		if err := os.MkdirAll(p, os.ModePerm); err != nil {
			return err
		}

		var sink bytes.Buffer
		if err := v.Execute(&sink, options); err != nil {
			return err
		}
		f, err := os.Create(path.Join(g.tmpDirPath, k))
		if err != nil {
			log.WithError(err).Fatal("error while creating file")
			return err
		}
		defer f.Close()
		if _, err := f.Write(sink.Bytes()); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(options.DstDir, os.ModePerm); err != nil {
		return err
	}

	dir := path.Join(options.DstDir, options.Name)
	if err := os.Rename(g.tmpDirPath, dir); err != nil {
		return err
	}
	log.Info("generation done")

	f := color.GreenString(`
	cd %s
	
	# Download dependent tools. dep, protoc and protoc plugins
	make install-deptools

	# build 
	make clean build

	./bin/%s`)
	fmt.Printf(f, dir, options.Name)
	fmt.Println()

	return nil
}

func (d DeploymentType) String() string {

	switch d {
	case K8SManifest:
		return "k8s"
	case HemlChart:
		return "helm"
	default:
		return "unknown"
	}
}

// ValueOf returns Typed DeploymentType
func ValueOf(d string) (DeploymentType, error) {

	ud := strings.ToLower(d)
	switch ud {
	case "k8s":
		return K8SManifest, nil
	case "helm":
		return HemlChart, nil
	default:
		return UnknownDeployemntType, errors.New("unknown deployment type")
	}
}
