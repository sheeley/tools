package sievegen

import (
	"strings"
)

func Not(rule *Rule) *Rule {
	rule.Not = true
	return rule
}

func AllOf(rules ...*Rule) *Rule {
	return &Rule{
		Children: rules,
		Mode:     ModeAll,
	}
}

func AnyOf(rules ...*Rule) *Rule {
	return &Rule{
		Children: rules,
		Mode:     ModeAny,
	}
}

func From(address string) *Rule {
	if strings.Contains(address, "*") {
		return rule(`address :matches "From" "%s"`, address)
	}
	return rule(`address :is "From" "%s"`, address)
}

func Contains(s string) *Rule {
	return rule(`body :text :contains "%s"`, s)
}

func To(address string) *Rule {
	return Header("To", address)
}

func Subject(s string) *Rule {
	return Header("Subject", s)
}

func Header(header, s string) *Rule {
	if strings.Contains(s, "*") {
		return rule(`header :matches "%s" "%s"`, header, s)
	}
	return rule(`header :contains "%s" "%s"`, header, s)
}
