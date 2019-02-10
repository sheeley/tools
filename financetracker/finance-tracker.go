package financetracker

import "bytes"

type Input struct {
	Verbose, Detail, Summary bool
}

type Output struct {
	Detail, Summary string
}

func FinanceTracker(in *Input, b *Book) (*Output, error) {
	var buf bytes.Buffer
	o := &Output{}

	b.Detail(&buf)
	o.Detail = buf.String()

	buf.Reset()
	b.Summary(&buf)
	o.Summary = buf.String()

	return o, nil
}
