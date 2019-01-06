#! /usr/bin/env bash
set -eu # x

# to run using the example data:
# go build -buildmode=plugin -o example/comp.so example/example.go
# go run ../cmd/comp/main.go -p example/comp.so

# to run using your data
# Fill in `compdata/main.go` with your data - see `example/example.go` for an _you guessed it_ example!
go build -buildmode=plugin -o compdata/compdata.so compdata/main.go
comp -p compdata/compdata.so