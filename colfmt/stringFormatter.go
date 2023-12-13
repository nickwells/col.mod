package colfmt

import (
	"fmt"

	"github.com/nickwells/col.mod/v4/col"
)

// String records the values needed for the formatting of a
// string value.
type String struct {
	// W gives the minimum width of the string that should be printed
	W uint
	// MaxW gives the maximum width of the string, if it is set to zero then
	// no limit is applied. If it is set to a negative value then the W
	// value is used. If it is a positive value then that is used
	MaxW int
	// StrJust gives the justification to be used
	StrJust col.Justification
	// IgnoreNil, if set to true will make nil values print as the empty string
	IgnoreNil bool
}

// Formatted returns the value formatted as a string
func (f String) Formatted(v any) string {
	if f.IgnoreNil && v == nil {
		return ""
	}

	switch {
	case f.MaxW == 0:
		return fmt.Sprintf("%s", v)
	case f.MaxW < 0:
		return fmt.Sprintf("%.*s", f.W, v)
	default:
		return fmt.Sprintf("%.*s", f.MaxW, v)
	}
}

// Width returns the intended width of the value
func (f String) Width() uint {
	return f.W
}

// Just returns the justification of the value
func (f String) Just() col.Justification {
	return f.StrJust
}

// Check returns a nil error
func (f String) Check() error {
	return nil
}
