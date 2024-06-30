package filepath

import "strings"

// Comparison operation determines how paths are compared. IgnoreCase or CaseSensitive
type Comparison interface {
	// Equal checks if the two strings are equal
	Equal(first, second string) bool
	// comparison implements the comparison interface
	comparison()
}

type comparison string

const (
	IgnoreCase    comparison = "ignore_case"
	CaseSensitive comparison = "case_sensitive"
)

func (cmp comparison) comparison() {}

func (cmp comparison) Equal(s, t string) bool {
	switch cmp {
	case IgnoreCase:
		return strings.EqualFold(s, t)
	case CaseSensitive:
		return s == t
	}
	return false
}
