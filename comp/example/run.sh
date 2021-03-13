#! /usr/bin/env bash
set -eu # x

go build -buildmode=plugin -o compdata.so main.go
comp -p compdata.so
