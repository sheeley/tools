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

type RuleMode int

const (
	ModeAny RuleMode = iota
	ModeAll
)

func rule(format string, a ...interface{}) *Rule {
	return &Rule{s: fmt.Sprintf(format, a...)}
}

type Rule struct {
	s        string
	Not      bool
	Mode     RuleMode
	Children []*Rule
}

func (r *Rule) String(level int) string {
	primaryIndent := strings.Repeat("\t", level-1)
	secondaryIndent := primaryIndent + "\t"

	var o strings.Builder
	//o.WriteString(primaryIndent)
	if r.Not {
		o.WriteString("not ")
	}
	if r.s != "" {
		o.WriteString(r.s)
	}

	if len(r.Children) == 0 {
		return o.String()
	}

	if len(r.Children) == 1 {
		o.WriteString(r.Children[0].String(level + 1))
		return o.String()
	}

	if r.Mode == ModeAll {
		o.WriteString("allof(")
	} else {
		o.WriteString("anyof(")
	}

	l := len(r.Children) - 1
	for idx, r := range r.Children {
		o.WriteString("\n" + secondaryIndent + r.String(level+1))
		if idx < l {
			o.WriteString(",")
		}
	}

	o.WriteString("\n" + primaryIndent + ")")

	return o.String()
}

type RuleSet struct {
	Rule     *Rule
	Action   string
	Continue bool

	parent   *RuleSet
	children []*RuleSet
}

func (rs *RuleSet) Root() *RuleSet {
	if rs.parent != nil {
		return rs.parent.Root()
	}
	return rs
}

func (rs *RuleSet) AndContinue() *RuleSet {
	rs.Continue = true
	return rs
}

func (rs *RuleSet) AndForward(to, from string, rules ...*Rule) *RuleSet {
	child := ForwardTo(to, from, rules...)
	child.parent = rs
	child.Continue = true
	rs.children = append(rs.children, child)
	return child
}

func (rs *RuleSet) Write(w io.Writer, level int) error {
	tabs := strings.Repeat("\t", level-1)
	secondaryTabs := tabs + "\t"
	_, err := fmt.Fprintf(w, "%sif %s {\n", tabs, rs.Rule.String(level))
	if err != nil {
		return errs.Wrap(err)
	}
	_, err = fmt.Fprintf(w, "%s%s;\n", secondaryTabs, strings.Replace(rs.Action, "\n", "\n"+secondaryTabs, -1))
	if err != nil {
		return errs.Wrap(err)
	}

	if len(rs.children) > 0 {
		for _, c := range rs.children {
			err = c.Write(w, level+1)
			if err != nil {
				return errs.Wrap(err)
			}
		}
	}

	stop := fmt.Sprintf("%sstop;", secondaryTabs)
	if rs.Continue {
		stop = ""
	}
	_, err = fmt.Fprintf(w, "%s\n%s}\n\n", stop, tabs)
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}

type Sieve struct {
	RuleSets []*RuleSet
}

const header = `# https://www.fastmail.com/cgi-bin/sievetest.pl
# require ["fileinto", "reject", "vacation", "notify", "envelope", "body", "relational", "regex", "subaddress", "copy", "mailbox", "mboxmetadata", "servermetadata", "date", "index", "comparator-i;ascii-numeric", "variables", "imap4flags", "editheader", "duplicate", "vacation-seconds"];
`

func (s *Sieve) Write(w io.Writer, level int) error {
	_, err := fmt.Fprintf(w, header)
	if err != nil {
		return errs.Wrap(err)
	}

	for _, rs := range s.RuleSets {
		err = rs.Write(w, level+1)
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}
