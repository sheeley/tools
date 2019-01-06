package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
)

func main() {
	filesWithR := 0

	err := filepath.Walk("/Users/sheeley/Music", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errs.Wrap(err)
		}

		if !strings.Contains(path, "\r") {
			return nil
		}

		filesWithR++

		err = os.Rename(path, strings.Replace(path, "\r", "", -1))
		if err != nil {
			fmt.Println(path, strings.Replace(path, "\r", "", -1))
			return errs.Wrap(err)
		}

		return filepath.SkipDir
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(filesWithR)
}
