package colfmt

import (
	"fmt"

	"github.com/nickwells/col.mod/v4/col"
)

// Bool records the values needed for the formatting of a
// bool value.
type Bool struct {
	// W gives the minimum width of the bool that should be printed
	W uint
	// StrJust gives the justification to be used
	StrJust col.Justification
	// IgnoreNil, if set to true will make nil values print as the empty bool
	IgnoreNil bool
}

// Formatted returns the value formatted as a bool
func (f Bool) Formatted(v any) string {
	if f.IgnoreNil && v == nil {
		return ""
	}
	return fmt.Sprintf("%t", v)
}

// Width returns the intended width of the value
func (f Bool) Width() uint {
	return f.W
}

// Just returns the justification of the value
func (f Bool) Just() col.Justification {
	return f.StrJust
}

// Check returns a nil error
func (f Bool) Check() error {
	return nil
}
