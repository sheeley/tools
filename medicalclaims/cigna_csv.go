package medicalclaims

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/richardwilkes/toolbox/errs"
)

type CignaRow struct {
	ServiceDate      string `csv:"Service Date"`
	ProvidedBy       string `csv:"Provided By"`
	For              string `csv:"For"`
	Status           string `csv:"Status"`
	AmountBilled     string `csv:"Amount Billed"`
	WhatYourPlanPaid string `csv:"What Your Plan Paid"`
	MyAccountPaid    string `csv:"My Account Paid"`
	WhatIOwe         string `csv:"What I Owe"`
}

func ParseCSV(path string) ([]*CignaRow, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	var rows []*CignaRow
	if err = gocsv.UnmarshalFile(f, &rows); err != nil {
		return nil, errs.Wrap(err)
	}

	return rows, nil
}
