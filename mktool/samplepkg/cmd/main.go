package main

import (
	"flag"
	"fmt"

	"github.com/sheeley/tools/mktool/samplepkg"
)

func main() {
	in := &samplepkg.Input{}

	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.Parse()

	out, err := samplepkg.SampleFunc(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
