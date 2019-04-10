package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/cnative/servicebuilder/internal/builder"
	"github.com/cnative/servicebuilder/internal/templates/grpcwithgw"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "creates a new service",
	Long: `Generates necessary scafolding to get started with building the service.
For example:

$ servicebuilder new --module-name github.com/kustomers/contacts
or
$ servicebuilder new --module-name github.com/kustomers/contacts --path $GOPATH/src/github.com/kustomers/
or
$ servicebuilder new --module-name github.com/kustomers/contacts --image-name gcr.io/mycompany/customer_contants

This application is a tool to generate the needed files
to quickly create a cloud native micro service.`,
	Run: scafoldNewService,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("module-name", "m", "", `module name of the service
a typical value is of form <gitserver>/<gitorg>/<projectname>
an example module name is mycompany.com/kustomer/accounts
in this example 'accounts' is the service name`)

	newCmd.Flags().StringP("description", "", "", "a short description of the service")
	newCmd.Flags().StringP("image-name", "i", "", "container image name")
	newCmd.Flags().StringP("protoc-version", "", "v3.7.0", "protocol buffer version to use")
	newCmd.Flags().StringP("http-route-prefix", "", "/api/v1", "http route prefix")
	newCmd.Flags().StringP("deployment-type", "", "k8s", "deployment artifact to generate. Possible values [helm, k8s]")
	newCmd.Flags().StringP("domain-name", "", "localhost", "domain name")
	newCmd.Flags().StringP("resource", "r", "", "resource name")
	newCmd.Flags().StringP("path", "p", ".", "directory path where the project will be generated")
}

func parseAndValidateArgs(c *cobra.Command) (*builder.Options, error) {
	mname, err := c.Flags().GetString("module-name")
	if err != nil {
		return nil, err
	}

	mname = strings.TrimSpace(mname)
	mname = strings.Trim(mname, "/")

	if mname == "" {
		return nil, errors.New("module-name cannot be empty")
	}
	mparts := strings.Split(mname, "/")
	sz := len(mparts)

	// service name
	name := mparts[sz-1]
	name = strings.TrimSpace(name)

	description, err := c.Flags().GetString("description")
	if err != nil {
		return nil, err
	}

	protocVersion, err := c.Flags().GetString("protoc-version")
	if err != nil {
		return nil, err
	}

	p, err := c.Flags().GetString("path")
	if err != nil {
		return nil, err
	}

	if p == "." {
		cdir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		p = cdir
	}

	da, err := c.Flags().GetString("deployment-type")
	if err != nil {
		return nil, err
	}
	dtype, err := builder.ValueOf(da)
	if err != nil {
		return nil, err
	}

	routePrefix, err := c.Flags().GetString("http-route-prefix")
	if err != nil {
		return nil, err
	}

	domainName, err := c.Flags().GetString("domain-name")
	if err != nil {
		return nil, err
	}

	resName, err := c.Flags().GetString("resource")
	if err != nil {
		return nil, err
	}

	if resName == "" {
		resName = name
	}
	resName = strings.Title(resName)

	imgn, err := c.Flags().GetString("image-name")
	if err != nil {
		return nil, err
	}

	if imgn == "" {
		imgr := strings.Builder{}
		if sz > 2 {
			imgr.WriteString(strings.Trim(mparts[sz-2], " "))
			imgr.WriteString("/")
		}
		imgr.WriteString(name)
		imgn = imgr.String()
	}

	dir := path.Join(p, mname)
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil, errors.Errorf("directory %s already exists", dir)
	}

	return &builder.Options{
		Name:                  name,
		ModuleName:            mname,
		ResourceName:          resName,
		ImageName:             imgn,
		Description:           description,
		DstDir:                p,
		DeploymentType:        dtype,
		DomainName:            domainName,
		HTTPRoutePrefix:       routePrefix,
		ServiceBuilderVersion: getServiceBuilderVersion(),
		ProtocVersion:         protocVersion,
	}, nil
}

func getServiceBuilderVersion() string {

	if version == "dev" || version == "unknown" {
		return "latest"
	}

	return fmt.Sprintf("v%s", version)
}

func scafoldNewService(c *cobra.Command, args []string) {

	o, err := parseAndValidateArgs(c)
	if err != nil {
		log.WithError(err).Fatal("invalid args")
		os.Exit(1)
	}

	fmt.Println()
	log.WithField("version", version).Infof("sevicebuilder")
	log.WithFields(log.Fields{
		"name":            o.Name,
		"module-name":     o.ModuleName,
		"image-name":      o.ImageName,
		"destination-dir": o.DstDir,
		"protoc-version":  o.ProtocVersion,
	}).Info("parse and argument validation success")

	templateProvider, err := grpcwithgw.New(o)
	if err != nil {
		log.WithError(err).Fatal("error while creating template provider")
		os.Exit(1)
	}

	sb, err := builder.New(templateProvider)
	if err != nil {
		log.WithError(err).Fatal("error while creating service builder")
		os.Exit(1)
	}

	if err := sb.Generate(); err != nil {
		log.WithError(err).Fatal("error while generating project structure")
		os.Exit(1)
	}
}
