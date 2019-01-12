package main

import (
	"flag"
	"fmt"

	"github.com/sheeley/tools/grabnytimesrecipe"
)

func main() {
	in := &grabnytimesrecipe.Input{}

	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.Parse()

	in.URLs = flag.Args()

	out, err := grabnytimesrecipe.GrabNytimesRecipe(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
