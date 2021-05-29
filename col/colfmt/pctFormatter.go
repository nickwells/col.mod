package colfmt

import (
	"fmt"

	"github.com/nickwells/col.mod/v3/col"
)

// Percent records the values needed for the formatting of a proportion as a
// percentage value. The value is expected to be a proportion and so is
// multiplied by 100 to convert it into a percentage value and then a % sign
// is added to the end (unless SuppressPct is set to true)
type Percent struct {
	// W gives the minimum space to be taken by the formatted value
	W int
	// Prec gives the precision with which to print the value when formatted
	// Negative values are treated as zero
	Prec int
	// IgnoreNil, if set to true will make nil values print as the empty string
	IgnoreNil bool
	// SuppressPct, if set to true will cause the '%' sign not to be printed
	SuppressPct bool
	// Zeroes records any desired special handling for zero values
	Zeroes *FloatZeroHandler
}

// Formatted returns the value formatted as a percentage
func (f *Percent) Formatted(v interface{}) string {
	if v == nil {
		if f.IgnoreNil {
			return ""
		}
		return "nil"
	}

	if f.Prec < 0 {
		f.Prec = 0
	}

	var pct float64
	switch flt := v.(type) {
	case float64:
		pct = flt * 100
	case float32:
		pct = float64(flt) * 100
	default:
		return fmt.Sprintf("%.*f", f.Prec, v)
	}
	if ok, str := f.Zeroes.GetZeroStr(f.Prec, pct); ok {
		return fmt.Sprintf("%.*s", f.Width(), str)
	}

	pctSign := "%%"
	if f.SuppressPct {
		pctSign = ""
	}
	return fmt.Sprintf("%.*f"+pctSign, f.Prec, pct)
}

// Width returns the intended width of the value. An invalid width or one
// incompatible with the given precision is ignored
func (f Percent) Width() int {
	minWidth := 1
	if !f.SuppressPct {
		minWidth++ // for the % sign
	}

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
func (f Percent) Just() col.Justification {
	return col.Right
}
