package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sheeley/tools/budget"
	"github.com/sheeley/tools/human"
)

var (
	headers           = []string{"name", "annual", "monthly", "daily"}
	blankRow          = []string{}
	separatorRow      = []string{}
	remaining         float64
	showDetails       = true
	showSectionTotals = false
)

func main() {
	tw := tablewriter.NewWriter(os.Stdout)
	if showDetails {
		headers = append(headers, "remaining")
	}
	tw.SetHeader(headers)
	for range headers {
		blankRow = append(blankRow, "")
		separatorRow = append(separatorRow, "====")
	}

	for _, section := range budget.Sections {
		if tw.NumLines() > 0 {
			tw.Append(blankRow)
		}

		for _, entry := range section.Entries {
			remaining -= entry.AnnualCost
			if section.Title != "Income" || showDetails {
				row := entry.Row()
				if showDetails {
					row = append(row, human.Float(remaining))
				}
				tw.Append(row)
			}
		}

		if showSectionTotals {
			row := append([]string{"Section Total"}, section.Row()...)
			if showDetails {
				row = append(row, human.Float(remaining))
			}
			tw.Append(row)
		}
	}

	tw.Append(blankRow)
	tw.Append(separatorRow)
	totalRow := append([]string{"Total"}, budget.CreateColumns(budget.Salary-remaining)...)
	if showDetails {
		totalRow = append(totalRow, "")
		tw.Append(totalRow)
		remaining := append([]string{"Remaining"}, budget.CreateColumns(remaining)...)
		tw.Append(append(remaining, ""))
	} else {
		tw.Append(totalRow)
	}

	alignments := []int{tablewriter.ALIGN_LEFT}
	for len(alignments) < len(headers) {
		alignments = append(alignments, tablewriter.ALIGN_RIGHT)
	}
	tw.SetColumnAlignment(alignments)

	// tw.SetAutoMergeCells(true)
	tw.Render()
}
