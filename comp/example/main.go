package main

import (
	"github.com/sheeley/tools/comp/data"
	"github.com/sheeley/tools/human"
)

func main() {}

var Data = &data.Comp{}

func init() {
	Data.Stock = []*data.Stock{
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

	Data.Cash = []*data.Cash{
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
