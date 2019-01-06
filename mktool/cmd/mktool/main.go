package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/sheeley/tools/mktool"
)

func main() {
	in := &mktool.Input{}

	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.Parse()

	in.ToolName = flag.Arg(0)

	out, err := mktool.Mktool(in)
	if err != nil {
		panic(err)
	}

	// if an $EDITOR env var is set, open the enclosing folder, then the cmd file
	ed, ok := os.LookupEnv("EDITOR")
	if ok && ed != "" {
		fmt.Println("editor: ", ed)
		err = exec.Command(ed, out.ToolDir).Run()
		if err != nil {
			fmt.Println(err)
		}

		err = exec.Command(ed, out.CmdFile).Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}
