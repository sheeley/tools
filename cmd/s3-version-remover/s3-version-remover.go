package main

import (
	"flag"
	"fmt"

	"github.com/sheeley/tools/s3versionremover"
)

func main() {
	in := &s3versionremover.Input{}

	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.StringVar(&in.Bucket, "v", "", "s3 bucket name")
	flag.Parse()

	out, err := s3versionremover.S3VersionRemover(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
