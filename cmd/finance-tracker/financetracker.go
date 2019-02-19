package main

import (
	"flag"
	"fmt"

	"github.com/sheeley/tools/financetracker"
	"github.com/sheeley/tools/plugins"
)

func main() {
	in := &financetracker.Input{}
	l := &plugins.Loader{
		Name: "financetracker",
	}

	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.BoolVar(&in.Detail, "d", false, "show detail")
	flag.BoolVar(&in.ByTag, "t", true, "show tag summary")
	flag.BoolVar(&in.Summary, "s", false, "show account summary")
	flag.Parse()

	p, err := l.CompileAndLoad()
	if err != nil {
		panic(err)
	}

	v, err := p.Lookup("Book")
	if err != nil {
		panic(err)
	}

	// little ceaser pointer pointer due to plugin
	prov, ok := v.(**financetracker.Book)
	if !ok {
		panic("plugin.Data is not *financetracker.Book")
	}
	book := *prov

	out, err := financetracker.FinanceTracker(in, book)
	if err != nil {
		panic(err)
	}

	if in.Detail {
		fmt.Println(out.Detail)
	}

	if in.Summary {
		fmt.Println(out.Summary)
	}

	if in.ByTag {
		fmt.Println(out.ByTag)
	}
}
