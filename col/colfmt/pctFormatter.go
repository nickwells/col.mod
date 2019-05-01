package colfmt

import (
	"fmt"

	"github.com/nickwells/col.mod/col"
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
}

// Formatted returns the value formatted as a percentage
func (f Percent) Formatted(v interface{}) string {
	if f.IgnoreNil && v == nil {
		return ""
	}

	if f.Prec < 0 {
		f.Prec = 0
	}

	var pct float64
	switch f := v.(type) {
	case float64:
		pct = f * 100
	case float32:
		pct = float64(f) * 100
	default:
		return "not-a-float"
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
	pctWidth := 1
	if f.SuppressPct {
		pctWidth = 0
	}
	if f.Prec <= 0 {
		if f.W <= 0 {
			return 1 + pctWidth
		}
		return f.W + pctWidth
	}

	if (2 + f.Prec + pctWidth) > f.W {
		return 2 + f.Prec + pctWidth
	}

	return f.W
}

// Just returns the justification of the value
func (f Percent) Just() col.Justification {
	return col.Right
}
