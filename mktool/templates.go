package mktool

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"reflect"

	"github.com/richardwilkes/toolbox/errs"
)

var (
	templateMain *template.Template
	templateLib  *template.Template
	currPkg      = reflect.TypeOf(tool{}).PkgPath()
)

type tool struct {
	Package, FuncName string
}

func loadTemplate(name, templatePath string) (*template.Template, error) {
	b, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	// main.go: replace the import
	b = bytes.Replace(b, []byte(currPkg), []byte(path.Dir(currPkg)), -1)

	// main.go/lib.go: replace any package name / function name
	b = bytes.Replace(b, []byte("samplepkg"), []byte("{{.Package}}"), -1)
	b = bytes.Replace(b, []byte("SampleFunc"), []byte("{{.FuncName}}"), -1)

	return template.New("main").Parse(string(b))
}

func setTemplates() error {
	gp, ok := os.LookupEnv("GOPATH") // TODO: maybe? build.Default.GOPATH
	if !ok {
		return errors.New("$GOPATH not set")
	}

	var err error
	templateLib, err = loadTemplate("lib", path.Join(gp, "src", currPkg, "samplepkg/sample.go"))
	if err != nil {
		return errs.Wrap(err)
	}

	templateMain, err = loadTemplate("main", path.Join(gp, "src", currPkg, "samplepkg/cmd/main.go"))
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}
