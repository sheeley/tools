#! /usr/bin/env bash
set -eu # x

go build -buildmode=plugin -o comp.so example/example.go
go run cmd/comp/main.go