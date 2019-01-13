package budget

import (
	"github.com/sheeley/tools/human"
)

func CreateColumns(d float64) []string {
	return []string{
		human.Float(d),
		human.Float(d / 12),
		human.Float(d / 365),
	}
}

type Section struct {
	Title   string
	Entries []*Entry
}

func (s *Section) Row() []string {
	return CreateColumns(s.Total())
}

func (s *Section) Total() float64 {
	t := 0.0
	for _, e := range s.Entries {
		t -= e.AnnualCost
	}
	return t
}

type Entry struct {
	Name       string
	AnnualCost float64
}

func (m *Entry) Row() []string {
	return append([]string{m.Name}, CreateColumns(-1*m.AnnualCost)...)
}
