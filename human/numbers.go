package human

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var p = message.NewPrinter(language.English)

func Int(i int) string {
	return p.Sprintf("%d", i)
}

func Float(f float64) string {
	return p.Sprintf("%f", f)
}

func Dollar(f float64) string {
	return p.Sprintf("$%.2f", f)
}
