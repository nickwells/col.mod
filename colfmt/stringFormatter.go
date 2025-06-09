package colfmt

import (
	"fmt"

	"github.com/nickwells/col.mod/v4/col"
)

// String records the values needed for the formatting of a string value.
//
// See [NilHdlr] and [DupHdlr] for the settings that can be given through
// those types.
type String struct {
	// W gives the minimum width of the string that should be printed
	W uint
	// MaxW gives the maximum width of the string, if it is set to zero then
	// no limit is applied. If it is set to a negative value then the W
	// value is used. If it is a positive value then that is used
	MaxW int
	// StrJust gives the justification to be used
	StrJust col.Justification
	// DuplicateIndicator is the value to show if the value to be shown is
	// the same as the value shown on the previous line. Setting this value
	// without also setting the SkipDuplicates flag will have no effect. Note
	// that if the DuplicateIndicator is too long to fit in the column it
	// will be truncated according to the settings of the W and MaxW values.
	DuplicateIndicator string

	format string

	NilHdlr
	DupHdlr
}

// makeFormat sets the format string to be used to format the value. It uses
// the Verb to construct the format string.
func (f *String) makeFormat() {
	if f.format == "" {
		switch {
		case f.MaxW == 0:
			f.format = "%s"
		case f.MaxW < 0:
			f.format = fmt.Sprintf("%%.%ds", f.W)
		default:
			f.format = fmt.Sprintf("%%.%ds", f.MaxW)
		}
	}
}

// Formatted returns the value formatted as a string
func (f *String) Formatted(v any) string {
	if f.SkipNil(v) {
		return ""
	}

	if f.SkipDup(v) {
		v = f.DuplicateIndicator
	}

	f.makeFormat()

	return fmt.Sprintf(f.format, v)
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
