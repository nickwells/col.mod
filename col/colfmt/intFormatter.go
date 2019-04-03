package colfmt

import (
	"fmt"

	"github.com/nickwells/col.mod/col"
)

// Int records the values needed for the formatting of an int value
type Int struct {
	W         int
	IgnoreNil bool
}

// Formatted returns the value formatted as an int
func (f Int) Formatted(v interface{}) string {
	if f.IgnoreNil && v == nil {
		return ""
	}

	return fmt.Sprintf("%d", v)
}

// Width returns the intended width of the value
func (f Int) Width() int {
	return f.W
}

// Just returns the justification of the value
func (f Int) Just() col.Justification {
	return col.Right
}
