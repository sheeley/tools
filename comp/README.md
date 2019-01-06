# comp
A tool for tracking vesting of stock grants & cash bonuses with clawback periods.

## usage
See `example/run.sh`

## generate your own data
```
mkdir compdata
mkdataplugin -p github.com/sheeley/tools/comp/data -s Comp > compdata/main.go
```