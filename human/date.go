package human

import (
	"strconv"
	"time"

	"github.com/richardwilkes/toolbox/errs"
)

const (
	Layout    = "20060102"
	LayoutOut = "2006/01/02"
)

func Itot(date int) (time.Time, error) {
	if date > 20500101 {
		return time.Time{}, errs.Newf("%d too large", date)
	}
	if date < 19840101 {
		return time.Time{}, errs.Newf("%d too small", date)
	}
	return time.Parse(Layout, strconv.Itoa(date))
}

func MustItot(date int) time.Time {
	d, err := Itot(date)
	if err != nil {
		panic(err)
	}
	return d
}

func Date(t time.Time) string {
	return t.Format(LayoutOut)
}

func Ttoi(t time.Time) (int, error) {
	return strconv.Atoi(t.Format(Layout))
}

func MustTtoi(t time.Time) int {
	i, err := Ttoi(t)
	if err != nil {
		panic(err)
	}
	return i
}
