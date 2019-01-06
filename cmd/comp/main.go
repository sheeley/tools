package main

import (
	"flag"
	"os"
	"plugin"

	"github.com/sheeley/tools/comp"
	"github.com/sheeley/tools/comp/data"
	"github.com/sheeley/tools/human"
)

func main() {
	in := &comp.Input{}

	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.StringVar(&in.Plugin, "p", "comp.so", "Plugin path, required")
	flag.Parse()

	p, err := plugin.Open(in.Plugin)
	if err != nil {
		panic(err)
	}

	v, err := p.Lookup("Data")
	if err != nil {
		panic(err)
	}

	prov, ok := v.(data.Provider)
	if !ok {
		panic("plugin is not data.Provider")
	}

	in.Cash = prov.Cash()
	in.Stocks = prov.Stocks()
	in.Out = os.Stdout
	in.Date = human.MustItot(20120101)

	_, err = comp.Comp(in)
	if err != nil {
		panic(err)
	}
}
