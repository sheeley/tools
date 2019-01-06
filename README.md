# Tools
A set of golang/bash tools that I use in daily life.

## Installation
```
go install github.com/sheeley/tools/cmd/...
```

Additional suggestion, or, simply an instruction, if you're me:

Add `$GOPATH/src/github.com/sheeley/tools/scripts/` to your path.

### Requirements:
Go 1.11
Bash
Fun

## mktool
A tool for creating additional tools with useful defaults.
```
mktool some-new-tool
```

## human
A package for easy input & output of dates / numbers.

* Convert `20180101` => `time.Time` using `human.Itot` or `human.MustItot`.
* Easily output comma-delimited numbers with `human.Int`, `human.Float`, and `human.Dollar`