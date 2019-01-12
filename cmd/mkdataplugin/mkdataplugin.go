package main

import (
	"flag"
	"fmt"

	"github.com/sheeley/tools/mkdataplugin"
)

func main() {
	in := &mkdataplugin.Input{}

	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.StringVar(&in.Package, "p", "", "package name - e.g. github.com/sheeley/tools/comp/data")
	flag.StringVar(&in.Symbol, "s", "", "symbol name - e.g. Comp")
	flag.StringVar(&in.ExportName, "e", "Data", "export name")
	flag.Parse()

	out, err := mkdataplugin.Mkdataplugin(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(out.FileText)
}
