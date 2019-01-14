package data

import (
	"fmt"
	"time"

	"github.com/sheeley/tools/human"
)

var ShowUUID = false
var ShowMatches = false

type GroupedClaims map[string][]*Claim

type Claim struct {
	UUID string

	Provider  string
	Treatment string
	Patient   string
	Amounts   []float64
	Dates     []time.Time
	Paid      bool

	AmountBilled     float64
	WhatYourPlanPaid float64
	MyAccountPaid    float64
	WhatIOwe         float64

	Source string
	Match  interface{}
}

func MatchString(provider, patient string, cost float64, date time.Time) string {
	// return fmt.Sprintf("%s %.2f %s", provider, cost, human.Date(date))
	return fmt.Sprintf("%.2f %s", cost, human.Date(date))
}

func (c *Claim) String() string {
	return MatchString(c.Provider, c.Patient, c.Total(), c.Dates[0])
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
