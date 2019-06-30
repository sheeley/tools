package sievegen_test

import (
	"strings"
	"testing"

	"github.com/sheeley/tools/human"
	sg "github.com/sheeley/tools/sievegen"
)

func TestSieve(t *testing.T) {
	b := strings.Builder{}
	s := &sg.Sieve{
		RuleSets: []*sg.RuleSet{
			sg.Discard(sg.From("*mail.cybercoders.com")),
			sg.KeepInbox(sg.From("keep@in-inbox.com")),

			sg.MoveTo("gmail",
				sg.To("my-address@gmail.com")),

			sg.RedirectTo("mywork@email.com",
				sg.From("redirect-from@email.com")),

			sg.MoveTo("travel",
				sg.From("some-travel@site.com"),
				sg.From("some-other-travel@site.com"),
			).AndForward("some@external.site", "from@from.org", sg.Not(sg.From("*@external.site"))).Root(),

			sg.MoveTo("test", sg.Contains("huh!*").And(sg.Subject("testtest*"))).AndContinue(),

			sg.MoveTo("allof", sg.AllOf(sg.Contains("all"), sg.Contains("of"))),

			sg.MoveTo("anyof", sg.AnyOf(sg.Contains("any"), sg.Contains("of"))),
		},
	}

	s.Write(&b, 0)
	human.AssertOrDiff(t, strings.TrimSpace(expected), strings.TrimSpace(b.String()))
}

var expected = `
# https://www.fastmail.com/cgi-bin/sievetest.pl
# require ["fileinto", "reject", "vacation", "notify", "envelope", "body", "relational", "regex", "subaddress", "copy", "mailbox", "mboxmetadata", "servermetadata", "date", "index", "comparator-i;ascii-numeric", "variables", "imap4flags", "editheader", "duplicate", "vacation-seconds"];
if address :matches "From" "*mail.cybercoders.com" {
	discard;
	stop;
}

if address :is "From" "keep@in-inbox.com" {
	fileinto "INBOX";
	stop;
}

if header :contains "To" "my-address@gmail.com" {
	fileinto "INBOX.gmail";
	stop;
}

if address :is "From" "redirect-from@email.com" {
	redirect :copy "mywork@email.com";
	stop;
}

if anyof(
	address :is "From" "some-travel@site.com",
	address :is "From" "some-other-travel@site.com"
) {
	fileinto "INBOX.travel";
	if not address :matches "From" "*@external.site" {
		deleteheader "From";
		addheader "From" "from@from.org";
		redirect :copy "some@external.site";

	}

	stop;
}

if allof(
		body :text :contains "huh!*",
		header :matches "Subject" "testtest*"
	) {
	fileinto "INBOX.test";

}

if allof(
		body :text :contains "all",
		body :text :contains "of"
	) {
	fileinto "INBOX.allof";
	stop;
}

if anyof(
		body :text :contains "any",
		body :text :contains "of"
	) {
	fileinto "INBOX.anyof";
	stop;
}
`
