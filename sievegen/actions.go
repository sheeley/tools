package sievegen

import "fmt"

func Discard(from ...string) *RuleSet {
	var rules []string
	for _, f := range from {
		rules = append(rules, From(f))
	}
	return &RuleSet{
		Action: "discard",
		Rules:  rules,
	}
}

func ForwardTo(to string, from ...string) *RuleSet {
	return &RuleSet{
		Action:   fmt.Sprintf(`redirect :copy "%s"`, to),
		Rules:    from,
		Continue: true,
	}
}

func MoveTo(to string, from ...string) *RuleSet {
	return &RuleSet{
		Action: fmt.Sprintf(`fileinto "INBOX.%s"`, to),
		Rules:  from,
	}
}

func KeepInbox(from ...string) *RuleSet {
	return &RuleSet{
		Action: `fileinto "INBOX"`,
		Rules:  from,
	}
}
