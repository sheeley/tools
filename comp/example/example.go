package main

import (
	"github.com/sheeley/tools/comp/data"
	"github.com/sheeley/tools/human"
)

func main() {}

var Data = example{}

type example struct{}

func (d *example) Stocks() []*data.Stock {
	return []*data.Stock{
		&data.Stock{
			GrantDate:        human.MustItot(20100101),
			VestingStartDate: human.MustItot(20100101),
			ShareCount:       100,
			StrikePrice:      10,
		},
		&data.Stock{
			GrantDate:        human.MustItot(20110101),
			VestingStartDate: human.MustItot(20110101),
			ShareCount:       100,
			StrikePrice:      10,
		},
		&data.Stock{
			GrantDate:        human.MustItot(20120101),
			VestingStartDate: human.MustItot(20120101),
			ShareCount:       100,
			StrikePrice:      10,
		},
	}
}

func (d *example) Cash() []*data.Cash {
	return []*data.Cash{
		&data.Cash{
			GrantDate: human.MustItot(20100101),
			Value:     10000,
		},
		&data.Cash{
			GrantDate: human.MustItot(20110101),
			Value:     20000,
		},
		&data.Cash{
			GrantDate: human.MustItot(20120101),
			Value:     30000,
		},
	}
}
