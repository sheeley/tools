package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/sheeley/tools/bear"
)

type input struct {
	name string
	next bool
}

func main() {
	in := &input{}

	flag.BoolVar(&in.next, "n", false, "for next time")
	flag.Parse()

	in.name = flag.Arg(0)
	if in.name == "" {
		panic("name required")
	}

	next := ""
	if in.next {
		next = "/next"
	}

	err := bear.Open("create", map[string]string{
		"title": fmt.Sprintf("%s %s", in.name, time.Now().Format("01/02")),
		"tags":  fmt.Sprintf("_fanatics/people/%s%s", in.name, next),
	})

	if err != nil {
		panic(err)
	}
}
