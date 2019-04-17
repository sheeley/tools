package sievegen

import "fmt"

func Discard(rules ...string) *RuleSet {
	return &RuleSet{
		Action: "discard",
		Rules:  rules,
	}
}

func ForwardTo(to string, rules ...string) *RuleSet {
	return &RuleSet{
		Action:   fmt.Sprintf(`redirect :copy "%s"`, to),
		Rules:    rules,
		Continue: true,
	}
}

func MoveTo(to string, rules ...string) *RuleSet {
	return &RuleSet{
		Action: fmt.Sprintf(`fileinto "INBOX.%s"`, to),
		Rules:  rules,
	}
}

func KeepInbox(rules ...string) *RuleSet {
	return &RuleSet{
		Action: `fileinto "INBOX"`,
		Rules:  rules,
	}
}
