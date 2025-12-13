package colfmt

import (
	"fmt"

	"github.com/nickwells/col.mod/v5/col"
	"github.com/nickwells/mathutil.mod/v2/mathutil"
)

// Percent records the values needed for the formatting of a proportion as a
// percentage value. The value is expected to be a proportion and so is
// multiplied by 100 to convert it into a percentage value and then a % sign
// is added to the end (unless SuppressPct is set to true)
type Percent struct {
	// W gives the minimum space to be taken by the formatted value
	W uint
	// Prec gives the precision with which to print the value when formatted
	Prec uint
	// IgnoreNil, if set to true will make nil values print as the empty string
	IgnoreNil bool
	// SuppressPct, if set to true will cause the '%' sign not to be printed
	SuppressPct bool
	// Zeroes records any desired special handling for zero values
	Zeroes *FloatZeroHandler
}

// Formatted returns the value formatted as a percentage. That is it is taken
// to be a proportion and is converted into a percentage value. So passing it
// a value of 1.25 will return a value of 125% (or 125 depending on the
// setting of SuppressPct)
//
//nolint:cyclop
func (f *Percent) Formatted(v any) string {
	if v == nil {
		if f.IgnoreNil {
			return ""
		}

		return "nil"
	}

	pctSign := "%%"
	if f.SuppressPct {
		pctSign = ""
	}

	var pct float64

	switch flt := v.(type) {
	case float64:
		pct = mathutil.ToPercent(flt)
	case float32:
		pct = float64(mathutil.ToPercent(flt))
	case int64:
		pct = mathutil.ToPercent(float64(flt))
	case int32:
		pct = mathutil.ToPercent(float64(flt))
	case int16:
		pct = mathutil.ToPercent(float64(flt))
	case int8:
		pct = mathutil.ToPercent(float64(flt))
	case int:
		pct = mathutil.ToPercent(float64(flt))
	case uint64:
		pct = mathutil.ToPercent(float64(flt))
	case uint32:
		pct = mathutil.ToPercent(float64(flt))
	case uint16:
		pct = mathutil.ToPercent(float64(flt))
	case uint8:
		pct = mathutil.ToPercent(float64(flt))
	default:
		return fmt.Sprintf("Numeric value expected (got: %T): %v", v, v)
	}

	if ok, str := f.Zeroes.GetZeroStr(f.Prec, pct); ok {
		return fmt.Sprintf("%.*s", f.Width(), str)
	}

	return fmt.Sprintf("%.*f"+pctSign, f.Prec, pct)
}

// Width returns the intended width of the value. An invalid width or one
// incompatible with the given precision is ignored
func (f Percent) Width() uint {
	var minWidth uint = 1
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

// Check returns a nil error
func (f Percent) Check() error {
	return nil
}
