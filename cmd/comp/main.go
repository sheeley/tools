package main

import (
	"flag"
	"os"
	"time"

	"github.com/sheeley/tools/comp"
	"github.com/sheeley/tools/comp/data"
	"github.com/sheeley/tools/human"
	"github.com/sheeley/tools/plugins"
)

func main() {
	in := &comp.Input{}
	l := &plugins.Loader{}

	date := human.MustTtoi(time.Now())

	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.StringVar(&l.Path, "p", "comp.so", "Plugin path, required")
	flag.IntVar(&date, "d", date, "date override, YYYYMMDD")
	flag.Parse()

	p, err := l.Open()
	if err != nil {
		panic(err)
	}

	v, err := p.Lookup("Data")
	if err != nil {
		panic(err)
	}

	// little ceaser pointer pointer due to plugin
	prov, ok := v.(**data.Comp)
	if !ok {
		panic("plugin.Data is not *data.Comp")
	}
	cData := *prov

	in.Cash = cData.Cash
	in.Stocks = cData.Stock
	in.Out = os.Stdout
	in.Date = human.MustItot(date)

	_, err = comp.Comp(in)
	if err != nil {
		panic(err)
	}
}
