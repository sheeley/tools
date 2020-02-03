package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/sheeley/tools/bearexport"
	"github.com/sheeley/tools/input"
)

func main() {
	in := &bearexport.Input{}

	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.BoolVar(&in.IncludeTrashed, "include-trashed", false, "include trashed notes")
	flag.BoolVar(&in.SeparateFiles, "separate-files", false, "write notes into separate json files")
	flag.StringVar(&in.Outdir, "o", ".", "output directory")
	flag.Parse()

	running, err := bearexport.BearRunning(in)
	if err != nil {
		panic(err)
	}

	if running {
		fmt.Println("Bear is running - this can cause problems. It's suggested that you close Bear. Continue? [y/N]")
		char, err := input.ReadChar()
		if err != nil {
			panic(err)
		}

		if char != 'y' {
			os.Exit(0)
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if !path.IsAbs(in.Outdir) {
		in.Outdir = path.Join(wd, in.Outdir)
	}

	out, err := bearexport.BearExport(in)
	if err != nil {
		panic(err)
	}

	err = bearexport.WriteNotes(in, out.Notes)
	if err != nil {
		panic(err)
	}
}
