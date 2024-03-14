package filepath

import "strings"

// Comparison operation determines how paths are compared. IgnoreCase or CaseSensitive
type Comparison interface {
	Equal(first, second string) bool
	comparison()
}

type IgnoreCase struct{}
type CaseSensitive struct{}

func (cmp IgnoreCase) comparison() {}
func (cmp IgnoreCase) Equal(s, t string) bool {
	return strings.EqualFold(s, t)
}

func (cmp CaseSensitive) comparison() {}
func (cmp CaseSensitive) Equal(s, t string) bool {
	return s == t
}
