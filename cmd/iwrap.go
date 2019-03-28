package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/cnative/servicebuilder/internal/iwrap"
	"github.com/spf13/cobra"
)

type (
	parameters struct {
		file           string
		interfaceName  string
		packageName    string
		templatePath   string
		formatCode     bool
		outputDir      string
		templates      []string
		ignoredMethods []string
		customImports  []string
	}

	templateParams struct {
		PackageName   string
		InterfaceName string
		Methods       []*method
		CustomImports []string
	}

	method struct {
		Name    string
		Doc     []string
		Params  []*arg
		Returns []*arg
	}

	arg struct {
		Name string
		Type string
	}
)

var (
	fns = template.FuncMap{
		"last": func(x int, a interface{}) bool {
			return x == reflect.ValueOf(a).Len()-1
		},
		"isLastReturnError": func(returns []*arg) bool {
			l := len(returns)
			if l == 0 {
				return false
			}
			return returns[l-1].Type == "error"
		},
		"lastReturnName": func(returns []*arg) string {
			l := len(returns)
			if l == 0 {
				return ""
			}
			return returns[l-1].Name
		},
		"lowerCase": strings.ToLower,
	}

	iwrapCmd = &cobra.Command{
		Use:   "iwrap",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: execute,
	}
)

func init() {
	rootCmd.AddCommand(iwrapCmd)

	iwrapCmd.Flags().StringP("file", "f", "", "path to the file containing the interface")
	iwrapCmd.Flags().StringP("interface-name", "i", "", "name of the interface to use")
	iwrapCmd.Flags().StringP("package-name", "p", "", "package name to use")
	iwrapCmd.Flags().StringP("template-path", "", "", "path to the template")
	iwrapCmd.Flags().StringSliceP("templates", "t", []string{"tracing", "metrics"}, "name of the templates to use. If template-path is specified templates will be ignored. If both template-path and templates are not specified then 'metrics' & 'tracing' will be applied")
	iwrapCmd.Flags().BoolP("format", "z", true, "format output using gofmt")
	iwrapCmd.Flags().StringP("output-dir", "o", "-", "path to the output file (use - for stdout)")
	iwrapCmd.Flags().StringSliceP("ignore", "g", []string{}, "ignore the following methods (separate with commas)")
	iwrapCmd.Flags().StringSliceP("imports", "m", []string{}, "custom imports (separate with commas)")
}

func contains(needle string, haystack []string) bool {
	for _, hay := range haystack {
		if hay == needle {
			return true
		}
	}

	return false
}

func getType(n ast.Expr) string {
	switch x := n.(type) {
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", x.X, x.Sel)
	case *ast.Ident:
		return fmt.Sprintf("%s", x.Name)
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", getType(x.X))
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", getType(x.Elt))
	}

	return fmt.Sprintf("unable to process type: %T", n)
}

func loadTemplates(templatePath string, knownTemplates []string) ([]*template.Template, error) {

	templates := []*template.Template{}
	templateStrs := make(map[string]string)

	if templatePath != "" {
		f, err := os.Open(templatePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		basename := filepath.Base(templatePath)
		fileName := strings.TrimSuffix(basename, filepath.Ext(basename))
		fileName = strings.TrimPrefix(fileName, filepath.Dir(basename))
		key := strings.Replace(fileName, ".go", "", -1)

		templateStrs[key] = string(b)
	} else {
		for _, k := range knownTemplates {
			key := strings.ToLower(k)
			if t, ok := iwrap.KnownInterfaceTemplates[key]; ok {
				templateStrs[key] = t
			} else {
				return nil, errors.New("unknown template - " + k)
			}
		}
	}

	for key, ts := range templateStrs {
		t, err := template.New(key).Funcs(fns).Parse(ts)
		if err != nil {
			return nil, err
		}

		templates = append(templates, t)
	}

	return templates, nil
}

func parseCommandArgs(c *cobra.Command) (*parameters, error) {

	file, err := c.Flags().GetString("file")
	if err != nil {
		return nil, err
	}
	if file == "" {
		return nil, errors.New("source file containing interface not specified")
	}

	interfaceName, err := c.Flags().GetString("interface-name")
	if err != nil {
		return nil, err
	}
	if interfaceName == "" {
		return nil, errors.New("interface name not specified")
	}

	packageName, err := c.Flags().GetString("package-name")
	if err != nil {
		return nil, err
	}
	if packageName == "" {
		return nil, errors.New("package name not specified")
	}

	templatePath, err := c.Flags().GetString("template-path")
	if err != nil {
		return nil, err
	}
	templates, err := c.Flags().GetStringSlice("templates")
	if err != nil {
		return nil, err
	}

	formatCode, err := c.Flags().GetBool("format")
	if err != nil {
		return nil, err
	}

	outputDir, _ := c.Flags().GetString("output-dir")
	if err != nil {
		return nil, err
	}

	ignoredMethods, err := c.Flags().GetStringSlice("ignore")
	if err != nil {
		return nil, err
	}

	customImports, err := c.Flags().GetStringSlice("imports")
	if err != nil {
		return nil, err
	}

	return &parameters{
		file:           file,
		interfaceName:  interfaceName,
		packageName:    packageName,
		templatePath:   templatePath,
		templates:      templates,
		formatCode:     formatCode,
		outputDir:      outputDir,
		ignoredMethods: ignoredMethods,
		customImports:  customImports,
	}, nil
}

func getInterfaceMethods(file, interfaceName string) ([]*ast.Field, error) {

	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	methods := []*ast.Field{}

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name != interfaceName {
				break
			}

			i, ok := x.Type.(*ast.InterfaceType)
			if !ok {
				break
			}

			for _, m := range i.Methods.List {
				methods = append(methods, m)
			}
		}

		return true
	})

	return methods, nil
}

func asTemplateMethodsParam(methods []*ast.Field, ignoredMethods []string) []*method {

	m := []*method{}
	for _, met := range methods {

		if t, ok := met.Type.(*ast.FuncType); ok {
			name := met.Names[0].Name
			if contains(name, ignoredMethods) {
				continue
			}

			counter := 0
			params := []*arg{}
			for _, par := range t.Params.List {
				name := fmt.Sprintf("p%d", counter)
				if len(par.Names) > 0 {
					name = par.Names[0].Name
				}

				params = append(params, &arg{Name: name, Type: getType(par.Type)})
				counter++
			}

			counter = 0
			returns := []*arg{}
			for _, ret := range t.Results.List {
				name := fmt.Sprintf("r%d", counter)
				if len(ret.Names) > 0 {
					name = ret.Names[0].Name
				}

				returns = append(returns, &arg{Name: name, Type: getType(ret.Type)})
				counter++
			}

			doclines := []string{}
			if met.Doc != nil && len(met.Doc.List) > 0 {
				for _, line := range met.Doc.List {
					doclines = append(doclines, line.Text)
				}
			}
			m = append(m, &method{Name: name, Params: params, Returns: returns, Doc: doclines})
		}
	}

	return m
}

func execute(c *cobra.Command, args []string) error {

	params, err := parseCommandArgs(c)
	if err != nil {
		return err
	}

	methods, err := getInterfaceMethods(params.file, params.interfaceName)
	if err != nil {
		return err
	}

	tmplts, err := loadTemplates(params.templatePath, params.templates)
	if err != nil {
		return err
	}

	vm := &templateParams{
		InterfaceName: params.interfaceName,
		PackageName:   params.packageName,
		Methods:       asTemplateMethodsParam(methods, params.ignoredMethods),
		CustomImports: params.customImports,
	}

	for _, t := range tmplts {
		var sink bytes.Buffer
		err = t.Execute(&sink, vm)
		if err != nil {
			return err
		}

		b := sink.Bytes()
		if params.formatCode {
			b, err = format.Source(b)
			if err != nil {
				return err
			}
		}

		var out io.Writer = os.Stdout
		if params.outputDir != "-" {
			fn := fmt.Sprintf("%s%c%s_with_%s.go", params.outputDir, filepath.Separator, strings.ToLower(params.interfaceName), t.Name())
			f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}
			defer f.Close()
			out = f
		}

		_, err = out.Write(b)
	}

	return err
}
