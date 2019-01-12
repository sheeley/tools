package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/sheeley/tools/medicalclaims"
	"github.com/sheeley/tools/medicalclaims/data"
	"github.com/sheeley/tools/plugins"
)

func main() {
	in := &medicalclaims.Input{}
	l := &plugins.Loader{}

	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.StringVar(&l.Path, "p", "claims.so", "Plugin path, required")
	flag.Parse()

	p, err := l.Open()
	if err != nil {
		panic(err)
	}

	showGrouped := true
	gc, err := p.Lookup("GroupedClaims")
	if err != nil {
		if !strings.Contains(err.Error(), "symbol GroupedClaims not found") {
			panic(err)
		}
		showGrouped = false
	}

	if showGrouped {
		pgroupedClaims, ok := gc.(*data.GroupedClaims)
		if !ok {
			panic("plugin.GroupedClaims is not *data.GroupedClaims")
		}
		groupedClaims := *pgroupedClaims
		medicalclaims.PrintGrouped(groupedClaims)
	}

	c, err := p.Lookup("Claims")
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%T\n", c)
	pClaims, ok := c.(*[]*data.Claim)
	if !ok {
		panic("plugin.GroupedClaims is not *data.GroupedClaims")
	}
	claims := *pClaims
	medicalclaims.PrintTable(claims)

	out, err := medicalclaims.MedicalClaims(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
