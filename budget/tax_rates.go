package budget

import "fmt"

var (
	StateTaxRates = []*TaxRate{
		{105224, 9.3, 0},
		{537500, 10.3, 0},
	}

	FederalTaxRates = []*TaxRate{
		{19050, 12, 1905},
		{77400, 22, 8907},
		{165000, 24, 28179},
		{315000, 32, 64179},
	}
)

type TaxRate struct {
	IncomeFloor, Rate, Base float64
}

func TaxAmount(rates []*TaxRate, incomes ...float64) float64 {
	income := 0.0
	for _, i := range incomes {
		income += i
	}

	for _, r := range rates {
		if income > r.IncomeFloor {
			taxRate := r.Rate / 100
			fmt.Println(taxRate, income)
			if r.Base == 0 {
				return taxRate * income
			}

			excess := income - r.IncomeFloor
			return (excess * taxRate) + r.Base
		}
	}

	return 0
}
