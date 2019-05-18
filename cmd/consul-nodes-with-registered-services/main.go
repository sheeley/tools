package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sheeley/counter"

	"github.com/hashicorp/consul/api"
)

func main() {
	addr := ""
	flag.StringVar(&addr, "addr", "", "consul http addr")
	flag.Parse()
	cl, err := api.NewClient(&api.Config{
		Address: addr,
	})
	if err != nil {
		panic(err)
	}

	nodes := map[string]bool{}
	svcNodeCount := counter.String{}

	svcs, _, err := cl.Catalog().Services(nil)
	if err != nil {
		panic(err)
	}
	for s := range svcs {
		// fmt.Println(s)
		svc, _, err := cl.Catalog().Service(s, "", nil)
		if err != nil {
			panic(err)
		}
		for _, sv := range svc {
			nodes[sv.Node] = true
			svcNodeCount.Add(s)
			// fmt.Println(sv.Node)
		}
	}

	fmt.Printf("\n%d\ttotal\n", len(nodes))
	svcNodeCount.Print(os.Stdout)
}
