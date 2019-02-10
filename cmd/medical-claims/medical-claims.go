package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/sheeley/tools/medicalclaims"
	"github.com/sheeley/tools/medicalclaims/data"
	"github.com/sheeley/tools/plugins"
)

func main() {
	in := &medicalclaims.Input{}
	l := &plugins.Loader{
		Name: "medical-claims",
	}

	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.IntVar(&in.Year, "y", time.Now().Year(), "filter claims by year")
	flag.StringVar(&in.Provider, "provider", "", "filter claims by provider")
	flag.BoolVar(&data.ShowUUID, "u", false, "show UUIDs for claims")
	flag.BoolVar(&data.ShowMatches, "m", false, "show potential matches for claims found in CSV without uuid")

	flag.Parse()

	p, err := l.CompileAndLoad()
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

	claims = medicalclaims.Filter(in, claims)
	if len(claims) == 0 {
		fmt.Println("0 claims after filtering")
		return
	}
	medicalclaims.PrintTables(claims)
}
