package sievegen_test

import (
	"strings"
	"testing"

	"github.com/sheeley/tools/human"
	"github.com/sheeley/tools/sievegen"
)

func TestSieve(t *testing.T) {
	b := strings.Builder{}
	s := &sievegen.Sieve{}
	s.Write(&b)
	human.AssertOrDiff(t, strings.TrimSpace(b.String()), strings.TrimSpace(defaultT))
}

var defaultT = `
# https://www.fastmail.com/cgi-bin/sievetest.pl
# require ["fileinto", "reject", "vacation", "notify", "envelope", "body", "relational", "regex", "subaddress", "copy", "mailbox", "mboxmetadata", "servermetadata", "date", "index", "comparator-i;ascii-numeric", "variables", "imap4flags", "editheader", "duplicate", "vacation-seconds"];
`
