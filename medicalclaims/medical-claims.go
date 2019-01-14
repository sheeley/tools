package medicalclaims

import (
	"os"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/sheeley/tools/human"
	"github.com/sheeley/tools/medicalclaims/data"
)

type Input struct {
	Verbose  bool
	Year     int
	Provider string
}

type Output struct {
}

func PrintGrouped(groupedClaims data.GroupedClaims) {
	tw := tablewriter.NewWriter(os.Stdout)
	total := 0
	for k, cs := range groupedClaims {
		tw.Append([]string{k, strconv.Itoa(len(cs))})
		total += len(cs)
	}
	tw.Append([]string{"Total", strconv.Itoa(total)})
	tw.Render()
}

func Filter(in *Input, claims []*data.Claim) []*data.Claim {
	if in.Year == 0 && in.Provider == "" {
		return claims
	}
	var o []*data.Claim
	for _, c := range claims {
		if in.Provider != "" && c.Provider != in.Provider {
			continue
		}
		if in.Year > 0 && c.Dates[0].Year() != in.Year {
			continue
		}
		o = append(o, c)
	}
	return o
}

func PrintTables(claims []*data.Claim) {
	tw := tablewriter.NewWriter(os.Stdout)

	sort.Slice(claims, func(i, j int) bool {
		if claims[i].Dates[0].Month() == claims[j].Dates[0].Month() {
			if claims[i].Treatment != claims[j].Treatment {
				return claims[i].Treatment < claims[j].Treatment
			}
		}
		return claims[i].Dates[0].Before(claims[j].Dates[0])
	})

	// tw.SetAutoMergeCells(true)

	headers := []string{
		"Date",
		"Match",
		"Patient",
		"Provider",
		"Type",
		"Amount Billed",
		"My Account Paid",
		"What I Owe",
		"What Your Plan Paid",
	}
	tw.SetHeader(headers)

	total := &data.Claim{}
	blank := make([]string, len(headers))
	month := claims[0].Dates[0].Month()
	matchCount := 0
	for _, c := range claims {
		newMonth := c.Dates[0].Month()
		if newMonth != month {
			month = newMonth
			tw.Append(blank)
		}

		matched := ""
		if c.Match != nil {
			matched = "y"
			matchCount += 1
		}

		// fmt.Println(human.Date(c.Dates[0]),
		// 	matched,
		// 	c.Patient,
		// 	c.Provider,
		// 	c.Treatment,
		// 	human.Float(c.Total()),
		// 	human.Float(c.MyAccountPaid),
		// 	human.Float(c.WhatIOwe),
		// 	human.Float(c.WhatYourPlanPaid))
		tw.Append([]string{
			human.Date(c.Dates[0]),
			matched,
			c.Patient,
			c.Provider,
			c.Treatment,
			human.Float(c.Total()),
			human.Float(c.MyAccountPaid),
			human.Float(c.WhatIOwe),
			human.Float(c.WhatYourPlanPaid),
		})

		total.Amounts = append(total.Amounts, c.Amounts...)
		total.AmountBilled += c.AmountBilled
		total.MyAccountPaid += c.MyAccountPaid
		total.WhatIOwe += c.WhatIOwe
		total.WhatYourPlanPaid += c.WhatYourPlanPaid
	}

	tw.Append(blank)
	tw.Append([]string{
		"Total",
		strconv.Itoa(matchCount) + "/" + strconv.Itoa(len(claims)),
		"",
		"",
		"",
		human.Float(total.Total()),
		human.Float(total.MyAccountPaid),
		human.Float(total.WhatIOwe),
		human.Float(total.WhatYourPlanPaid),
	})

	tw.Render()
}
