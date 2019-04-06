package financetracker

import (
	"io"
	"sort"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/sheeley/tools/human"
)

type Book struct {
	Accounts []*Account
}

func (b *Book) CurrentValue() int {
	val := 0
	for _, acct := range b.Accounts {
		val += acct.CurrentValue()
	}
	return val
}

func (b *Book) Dates() []time.Time {
	dateSet := map[time.Time]bool{}
	for _, a := range b.Accounts {
		for _, d := range a.Dates {
			dateSet[d] = true
		}
	}

	if len(dateSet) == 0 {
		return nil
	}

	var dates []time.Time
	for d := range dateSet {
		dates = append(dates, d)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	return dates
}

func (b *Book) Summary(o io.Writer) {
	w := tablewriter.NewWriter(o)
	w.SetHeader([]string{"Account", "Last Date", "Last Value"})
	for _, a := range b.Accounts {
		date := ""
		value := 0
		if len(a.Dates) > 0 {
			date = human.Date(a.Dates[len(a.Dates)-1])
			value = a.CurrentValue()
		}
		w.Append([]string{a.Name, date, human.Int(value)})
	}
	w.Append([]string{"Total", "", human.Int(b.CurrentValue())})

	w.Render()
}

func (b *Book) Detail(o io.Writer) {
	w := tablewriter.NewWriter(o)

	last := 0
	headers := []string{"Date"}
	for _, a := range b.Accounts {
		if a.Invisible {
			continue
		}
		headers = append(headers, a.Name)
	}
	headers = append(headers, "Delta", "Net worth")
	w.SetHeader(headers)

	for _, d := range b.Dates() {
		row := []string{human.Date(d)}
		dVal := 0
		for _, a := range b.Accounts {
			val := a.ValueAt(d)
			dVal += val
			if a.Invisible {
				// this will hide a specific account without removing it from the total
				continue
			}
			row = append(row, human.Int(val))
		}

		delta := dVal - last
		if last == 0 {
			delta = 0
		}
		row = append(row, human.Int(delta), human.Int(dVal))
		last = dVal
		w.Append(row)
	}

	w.Render()
}

type tagSummary struct {
	current, lastMonth, lastYear float64
}

func (b *Book) ByTag(o io.Writer) {
	tv := map[string]*tagSummary{}
	lm := time.Now().Add(-24 * 30 * time.Hour)
	ly := time.Now().Add(-24 * 365 * time.Hour)

	for _, acct := range b.Accounts {
		for _, t := range acct.Tags {
			if _, ok := tv[t]; !ok {
				tv[t] = &tagSummary{}
			}
			tv[t].current += float64(acct.CurrentValue())
			tv[t].lastMonth += float64(acct.ValueAt(lm))
			tv[t].lastYear += float64(acct.ValueAt(ly))
		}
	}

	w := tablewriter.NewWriter(o)
	w.SetHeader([]string{"tag", "ly", "lm", "current"})
	for t, v := range tv {
		w.Append([]string{t, human.Dollar(v.lastYear), human.Dollar(v.lastMonth), human.Dollar(v.current)})
	}
	w.Render()
}
