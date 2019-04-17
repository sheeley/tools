package sievegen

import (
	"fmt"
	"strings"
)

func AllOf(rules ...string) string {
	return fmt.Sprintf("allof(\n\t\t%s\n\t)", strings.Join(rules, ",\n\t\t"))
}

func AnyOf(rules ...string) string {
	return fmt.Sprintf("anyof(\n\t\t%s\n\t)", strings.Join(rules, ",\n\t\t"))
}

func From(address string) string {
	if strings.Contains(address, "*") {
		return fmt.Sprintf(`address :matches "From" "%s"`, address)
	}
	return fmt.Sprintf(`address :is "From" "%s"`, address)
}

func FromContains(address string) string {
	return fmt.Sprintf(`address :contains "From" "%s"`, address)
}

func Contains(s string) string {
	return fmt.Sprintf(`body :text :contains "%s"`, s)
}

func To(address string) string {
	return Header("To", address)
}

func Subject(s string) string {
	return Header("Subject", s)
}

func Header(header, s string) string {
	if strings.Contains(s, "*") {
		return fmt.Sprintf(`header :matches "%s" "%s"`, header, s)
	}
	return fmt.Sprintf(`header :contains "%s" "%s"`, header, s)
}
