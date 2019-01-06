package data

import (
	"time"

	"github.com/sheeley/tools/human"
)

type GroupedClaims map[string][]*Claim

type Claim struct {
	Provider  string
	Treatment string
	Patient   string
	Amounts   []float64
	Dates     []time.Time
	Paid      bool
}

func (c *Claim) SetAmounts(a ...float64) {
	c.Amounts = a
}

func (c *Claim) SetDates(ds ...int) {
	c.Dates = nil
	for _, d := range ds {
		c.Dates = append(c.Dates, human.MustItot(d))
	}
}

func (c *Claim) Total() float64 {
	var t float64
	for _, a := range c.Amounts {
		t += a
	}
	return t
}

type ClaimConfigurer func(c *Claim)
