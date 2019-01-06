package main

import (
	"time"

	"github.com/sheeley/tools/bear"
)

func main() {
	err := bear.Open("create", map[string]string{
		"title": time.Now().Format("01/02"),
		"tags":  "journal",
	})
	if err != nil {
		panic(err)
	}
}
