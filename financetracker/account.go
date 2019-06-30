package financetracker

import (
	"sort"
	"time"
)

type Account struct {
	Source      *Account
	Destination *Account

	Name      string
	Entries   map[time.Time]*Entry
	Dates     []time.Time
	Invisible bool // this will hide a specific account without removing it from the total
	Tags      []string

	PaymentDay    int
	PaymentAmount float64
}

type AccountConfig func(*Account)

func AccountInvisible(a *Account) {
	a.Invisible = true
}

func NewAccount(name string, entries []*Entry, configs ...AccountConfig) *Account {
	e := map[time.Time]*Entry{}
	var dates []time.Time
	if len(entries) > 0 {
		for _, entry := range entries {
			e[entry.Date] = entry
			dates = append(dates, entry.Date)
		}
		sort.Slice(dates, func(i, j int) bool {
			return dates[i].Before(dates[j])
		})
	}
	a := &Account{
		Name:    name,
		Entries: e,
		Dates:   dates,
	}

	for _, cfg := range configs {
		cfg(a)
	}

	return a
}

func (a *Account) CurrentValue() int {
	if len(a.Dates) == 0 {
		return 0
	}

	if e, ok := a.Entries[a.Dates[len(a.Dates)-1]]; ok {
		return e.Value
	}

	return 0
}

func (a *Account) ValueAt(d time.Time) int {
	if len(a.Entries) == 0 {
		return 0
	}

	if e, ok := a.Entries[d]; ok {
		return e.Value
	}

	first := a.Dates[0]
	if d.Before(first) {
		return 0 // a.Entries[first].Value
	}

	last := a.Dates[len(a.Dates)-1]
	if d.After(last) {
		return a.Entries[last].Value
	}

	var previous time.Time
	var next time.Time

	for _, ad := range a.Dates {
		if ad.Before(d) {
			previous = ad
		}
		if ad.After(d) {
			next = ad
			break
		}
	}

	before := float64(a.Entries[previous].Value)
	after := float64(a.Entries[next].Value)
	beginningToEnd := next.Sub(previous).Hours()
	slope := (after - before) / beginningToEnd
	beginningToDate := d.Sub(previous).Hours()
	value := (slope * beginningToDate) + before

	return int(value)
}
