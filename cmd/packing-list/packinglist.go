package main

import (
	"flag"
	"fmt"

	"github.com/sheeley/tools/packinglist"
)

func main() {
	in := &packinglist.Input{}

	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.Parse()

	out, err := packinglist.PackingList(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
