package data

import "time"

// Provider is used to separate comp data from the internals
type Provider interface {
	Stocks() []*Stock
	Cash() []*Cash
}

// Stock represents RSUs or Options with
type Stock struct {
	ShareCount  int
	StrikePrice float64

	VestingStartDate time.Time
	// VestingPeriod defaults to 4 years
	VestingPeriod time.Duration

	// Optional
	GrantDate time.Time
}

// Cash represents a cash bonus with a clawback period
type Cash struct {
	Value          float64
	GrantDate      time.Time
	ClawbackPeriod time.Duration
}

func (cg *Cash) Payback(t time.Time) float64 {
	if cg.GrantDate.After(t) {
		return 0.0
	}

	diff := t.Sub(cg.GrantDate).Hours() / Year
	if diff > 1 {
		return 0.0
	}

	return cg.Value * (1.0 - diff)
}

const Year = 365 * 24
