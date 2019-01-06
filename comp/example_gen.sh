#! /usr/bin/env bash
set -eu # x

mkdir compdata
mkdataplugin -p github.com/sheeley/tools/comp/data -s Comp > compdata/main.go
