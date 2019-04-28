package financetracker

import "bytes"

type Input struct {
	Verbose, Detail, Summary, ByTag, Payment bool
}

type Output struct {
	Detail, Summary, ByTag, Payment string
}

func FinanceTracker(in *Input, b *Book) (*Output, error) {
	var buf bytes.Buffer
	o := &Output{}

	if in.Detail {
		b.Detail(&buf)
		o.Detail = buf.String()
	}

	if in.Summary {
		buf.Reset()
		b.Summary(&buf)
		o.Summary = buf.String()
	}

	if in.Payment {
		buf.Reset()
		b.Payments(&buf)
		o.Payment = buf.String()
	}

	if in.ByTag {
		buf.Reset()
		b.ByTag(&buf)
		o.ByTag = buf.String()
	}

	return o, nil
}
