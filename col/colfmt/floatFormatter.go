package colfmt

import (
	"fmt"

	"github.com/nickwells/col.mod/v2/col"
)

// Float records the values needed for the formatting of a
// float(64/32) value
type Float struct {
	// W gives the minimum space to be taken by the formatted value
	W int
	// Prec gives the precision with which to print the value when formatted
	// Negative values are treated as zero
	Prec int
	// IgnoreNil, if set to true will make nil values print as the empty string
	IgnoreNil bool
	// Zeroes records any desired special handling for zero values
	Zeroes *FloatZeroHandler
}

// Formatted returns the value formatted as a float
func (f *Float) Formatted(v interface{}) string {
	if f.IgnoreNil && v == nil {
		return ""
	}

	if f.Prec < 0 {
		f.Prec = 0
	}

	if ok, str := f.Zeroes.GetZeroStr(f.Prec, v); ok {
		return fmt.Sprintf("%.*s", f.Width(), str)
	}
	return fmt.Sprintf("%.*f", f.Prec, v)
}

// Width returns the intended width of the value. An invalid width or one
// incompatible with the given precision is ignored
func (f Float) Width() int {
	minWidth := 1
	if f.Prec > 0 {
		minWidth++ // for the decimal place
		minWidth += f.Prec
	}

	if minWidth > f.W {
		return minWidth
	}

	return f.W
}

// Just returns the justification of the value
func (f Float) Just() col.Justification {
	return col.Right
}
