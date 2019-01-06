package medicalclaims

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/sheeley/tools/human"
	"github.com/sheeley/tools/medicalclaims/data"
)

type Input struct {
	Verbose bool
}

type Output struct {
}

func MedicalClaims(in *Input) (*Output, error) {
	return &Output{}, nil
}

func PrintGrouped(groupedClaims data.GroupedClaims) {
	total := 0
	for k, cs := range groupedClaims {
		fmt.Println(k, len(cs))
		total += len(cs)
	}
	fmt.Println(total)
}

func PrintTable(claims []*data.Claim) {
	tw := tablewriter.NewWriter(os.Stdout)

	sort.Slice(claims, func(i, j int) bool {
		if claims[i].Dates[0].Month() == claims[j].Dates[0].Month() {
			if claims[i].Treatment != claims[j].Treatment {
				return claims[i].Treatment < claims[j].Treatment
			}
		}
		return claims[i].Dates[0].Before(claims[j].Dates[0])
	})

	tw.SetAutoMergeCells(true)
	headers := []string{"Date", "Patient", "Provider", "Type", "Total", "Paid"}
	blank := make([]string, len(headers))
	tw.SetHeader(headers)

	var total float64
	var paidTotal float64
	month := claims[0].Dates[0].Month()
	for _, c := range claims {
		newMonth := c.Dates[0].Month()
		if newMonth != month {
			month = newMonth
			tw.Append(blank)
		}
		paid := ""
		if c.Paid {
			paid = "X"
			paidTotal += c.Total()
		}
		total += c.Total()

		cells := []string{
			human.Date(c.Dates[0]),
			c.Patient,
			c.Provider,
			c.Treatment,
			human.Float(c.Total()),
			paid,
		}
		tw.Append(cells)
	}
	tw.Render()

	tw = tablewriter.NewWriter(os.Stdout)
	tw.SetHeader([]string{"Type", "Count", "Total $", "Re-paid $"})
	tw.Append([]string{"Claimed", strconv.Itoa(len(claims)), human.Float(total), human.Float(paidTotal)})
	// tw.Append([]string{"Received", strconv.Itoa(len(receivedClaims)), human.Float(receivedAmt), human.Float(receivedPaid)})
	tw.Render()
}
