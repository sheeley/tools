package sievegen

import (
	"fmt"
	"io"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
)

type Input struct {
	Verbose bool
}

type Output struct {
	Sieve string
}

type RuleSet struct {
	Rules    []string
	Action   string
	Continue bool
}

func (rs *RuleSet) Write(w io.Writer) (int, error) {
	template := `
if %s {
	%s;%s
}
`

	stop := "\n\tstop;"
	if rs.Continue {
		stop = ""
	}

	condition := rs.Rules[0]
	if len(rs.Rules) > 1 {
		condTemp := `anyof(
	%s
)`
		condition = fmt.Sprintf(condTemp, strings.Join(rs.Rules, ",\n\t"))
	}

	return fmt.Fprintf(w, template, condition, rs.Action, stop)
}

type Sieve struct {
	RuleSets []*RuleSet
}

func (s *Sieve) Write(w io.Writer) error {
	_, err := fmt.Fprintf(w, `# https://www.fastmail.com/cgi-bin/sievetest.pl\n# require ["fileinto", "reject", "vacation", "notify", "envelope", "body", "relational", "regex", "subaddress", "copy", "mailbox", "mboxmetadata", "servermetadata", "date", "index", "comparator-i;ascii-numeric", "variables", "imap4flags", "editheader", "duplicate", "vacation-seconds"];`)
	if err != nil {
		return errs.Wrap(err)
	}

	for _, rs := range s.RuleSets {
		_, err = rs.Write(w)
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}
