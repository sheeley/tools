package mkdataplugin

import (
	"path"
	"strings"

	"github.com/alecthomas/template"
	"github.com/richardwilkes/toolbox/errs"
)

type Input struct {
	Verbose                     bool
	Package, Symbol, ExportName string
}

type Output struct {
	FileText string
}

type templateData struct {
	Symbol, Package, ShortPackage, ExportName string
}

const mainTemplate = `package main

import (
	"{{.Package}}"
)

func main() {}

var {{.ExportName}} = &{{.ShortPackage}}.{{.Symbol}}{}

func init() {
	// TODO: fill {{.ExportName}} in here.
}
`

func Mkdataplugin(in *Input) (*Output, error) {
	t, err := template.New("").Parse(mainTemplate)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	data := &templateData{
		ExportName:   strings.TrimSpace(in.ExportName),
		Package:      strings.TrimSpace(in.Package),
		Symbol:       strings.TrimSpace(in.Symbol),
		ShortPackage: path.Base(in.Package),
	}

	if data.ExportName == "" || data.Package == "" || data.Symbol == "" {
		return nil, errs.Newf("must provide ExportName (%s), Package (%s), Symbol (%s)", data.ExportName, data.Package, data.Symbol)
	}

	sb := &strings.Builder{}
	err = t.Execute(sb, data)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return &Output{
		FileText: sb.String(),
	}, nil
}
