package financetracker

import (
	"time"

	"github.com/sheeley/tools/human"
)

type Entry struct {
	Date  time.Time
	Value int
}

func NewEntry(date int, value int) *Entry {
	return &Entry{
		Date:  human.MustItot(date),
		Value: value,
	}
}
