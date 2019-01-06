package comp

import (
	"fmt"
	"io"
	"math"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/sheeley/tools/comp/data"
	"github.com/sheeley/tools/human"
)

type Input struct {
	Verbose bool
	Plugin  string

	Stocks []*data.Stock
	Cash   []*data.Cash
	Date   time.Time

	Out io.Writer
}

type Output struct{}

func Comp(in *Input) (*Output, error) {
	if in.Date.IsZero() {
		in.Date = time.Now()
	}

	if len(in.Stocks) > 0 {
		stocks(in)
	}

	if len(in.Cash) > 0 {
		cash(in)
	}

	return &Output{}, nil
}

func stocks(in *Input) {
	total := 0
	totalValue := 0.0
	vested := 0
	vestedValue := 0.0
	tw := tablewriter.NewWriter(in.Out)
	tw.SetHeader([]string{"Vesting Started", "Vesting Period (Years)", "Shares", "Strike", "Value", "Vested Shares", "Vested Value"})

	for _, g := range in.Stocks {
		if g.VestingPeriod == 0 {
			g.VestingPeriod = 4 * data.Year
		}

		total += g.ShareCount
		totalValue += float64(g.ShareCount) * g.StrikePrice
		vestedShares := float64(g.ShareCount)

		fullYears := math.Floor(in.Date.Sub(g.VestingStartDate).Hours() / float64(data.Year/time.Hour))

		vp := float64(g.VestingPeriod / data.Year)
		if fullYears < vp {
			vestedPercent := fullYears / vp
			vestedShares = vestedPercent * vestedShares
		}

		vested += int(vestedShares)
		vestedValue += vestedShares * g.StrikePrice
		tw.Append([]string{
			human.Date(g.VestingStartDate),
			fmt.Sprintf("%.0f", vp),
			human.Int(g.ShareCount),
			human.Float(g.StrikePrice),
			human.Float(float64(g.ShareCount) * g.StrikePrice),
			human.Int(int(vestedShares)),
			human.Float(vestedShares * g.StrikePrice),
		})
	}
	tw.Append([]string{"Total", "", human.Int(total), "", human.Float(totalValue), human.Int(vested), human.Float(vestedValue)})

	tw.Render()
}

func cash(in *Input) {
	tw := tablewriter.NewWriter(in.Out)
	tw.SetHeader([]string{"Granted", "Clawback Period (Years)", "Value", "Payback Amount"})

	paybackAmount := 0.0
	cashTotal := 0.0
	for _, cg := range in.Cash {
		if cg.ClawbackPeriod == 0 {
			cg.ClawbackPeriod = data.Year
		}
		pb := cg.Payback(in.Date)
		cashTotal += cg.Value
		tw.Append([]string{
			human.Date(cg.GrantDate),
			fmt.Sprintf("%.0f", float64(cg.ClawbackPeriod/data.Year)),
			human.Float(cg.Value),
			human.Float(pb),
		})
		paybackAmount += pb
	}
	tw.Append([]string{"Total", "", human.Float(cashTotal), human.Float(paybackAmount)})
	tw.Render()
}
