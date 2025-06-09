package colfmt

import (
	"fmt"

	"github.com/nickwells/col.mod/v5/col"
)

// Bool records the values needed for the formatting of a bool value.
//
// See [NilHdlr] and [DupHdlr] for the settings that can be given through
// those types.
type Bool struct {
	// W gives the minimum width of the bool that should be printed
	W uint
	// StrJust gives the justification to be used
	StrJust col.Justification

	NilHdlr
	DupHdlr
}

// Formatted returns the value formatted as a bool
func (f *Bool) Formatted(v any) string {
	if f.SkipNil(v) {
		return ""
	}

	if f.SkipDup(v) {
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
