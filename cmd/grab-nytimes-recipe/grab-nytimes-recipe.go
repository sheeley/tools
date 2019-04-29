package main

import (
	"flag"
	"fmt"

	"github.com/sheeley/tools/bear"
	"github.com/sheeley/tools/grabnytimesrecipe"
)

func main() {
	in := &grabnytimesrecipe.Input{}
	flag.BoolVar(&in.CreateBearEntry, "b", true, "create Bear entry")
	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.Parse()

	in.URLs = flag.Args()

	out, err := grabnytimesrecipe.GrabNytimesRecipe(in)
	if err != nil {
		panic(err)
	}

	for _, r := range out.Results {
		if !in.CreateBearEntry {
			fmt.Println(r.Title + "\n\n" + r.Body)
			continue
		}

		err = bear.Create(&bear.Entry{
			Title: r.Title,
			Body:  r.Body,
		})
		if err != nil {
			panic(err)
		}
	}
}
