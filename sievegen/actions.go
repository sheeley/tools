package sievegen

import "fmt"

func Discard(rules ...*Rule) *RuleSet {
	return &RuleSet{
		Action: "discard",
		Rule:   &Rule{Children: rules},
	}
}

// RedirectTo sends the _exact_ email along to the destination - no modification of anything
func RedirectTo(to string, rules ...*Rule) *RuleSet {
	return &RuleSet{
		Action: fmt.Sprintf(`redirect :copy "%s"`, to),
		Rule:   &Rule{Children: rules},
	}
}

// ForwardTo modifies the From header, making it actually look like you've forwarded it, rather than the raw RedirectTo
func ForwardTo(to, from string, rules ...*Rule) *RuleSet {
	t := `deleteheader "From";
addheader "From" "%s";
redirect :copy "%s"`
	return &RuleSet{
		Action: fmt.Sprintf(t, from, to),
		Rule:   &Rule{Children: rules},
	}
}

func MoveTo(to string, rules ...*Rule) *RuleSet {
	return &RuleSet{
		Action: fmt.Sprintf(`fileinto "INBOX.%s"`, to),
		Rule:   &Rule{Children: rules},
	}
}

func KeepInbox(rules ...*Rule) *RuleSet {
	return &RuleSet{
		Action: `fileinto "INBOX"`,
		Rule:   &Rule{Children: rules},
	}
}
