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
	if strings.Contains(address, "*") {
		return fmt.Sprintf(`header :matches "To" "%s"`, address)
	}
	return fmt.Sprintf(`header :contains "To" "%s"`, address)
}

func Subject(s string) string {
	if strings.Contains(s, "*") {
		return fmt.Sprintf(`header :matches "Subject" "%s"`, s)
	}
	return fmt.Sprintf(`header :contains "Subject" "%s"`, s)
}
