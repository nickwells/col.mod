package colfmt

import (
	"fmt"
	"math"
	"strings"

	"github.com/nickwells/col.mod/v4/col"
)

// Float records the values needed for the formatting of a
// float(64/32) value
type Float struct {
	// W gives the minimum space to be taken by the formatted value
	W uint
	// Prec gives the precision with which to print the value when formatted
	// Negative values are treated as zero
	Prec uint
	// IgnoreNil, if set to true will make nil values print as the empty string
	IgnoreNil bool
	// Zeroes records any desired special handling for zero values
	Zeroes *FloatZeroHandler
	// Verb specifies the formatting verb. If left unset it will use
	// 'f'. There will be a panic if it is not one of 'eEfFgGxX'
	Verb rune
	// TrimTrailingZeroes removes any trailing zeroes after the decimal
	// point. It leaves a zero immediately after the point
	TrimTrailingZeroes bool
	// ReformatOutOfBoundValues will generate a new format to be used if the
	// passed value is too big or too small to be shown in the space
	// available
	ReformatOutOfBoundValues bool
}

// makeFormat returns a format string to be used to report the value. It uses
// the Verb to construct the format string. It also consults the magnitude of
// the value and the ReformatOutOfBoundValues flag to decide whether to use a
// different format.
func (f Float) makeFormat(v any) string {
	format := ""
	switch f.Verb {
	case 0, 'f', 'F':
		format = "%.*f"
		if f.ReformatOutOfBoundValues {
			f64, ok := getValAsFloat64(v)
			if ok &&
				(f64 < math.Pow(10, -float64(f.Prec)) ||
					f64 > math.Pow(10, float64(f.Width())-float64(f.Prec))) {
				format = "%.*g"
			}
		}
	case 'e', 'E', 'g', 'G', 'x', 'X':
		format = "%.*" + string(f.Verb)
	default:
		panic(fmt.Errorf("%T: bad Format verb: %q", f, f.Verb))
	}
	return format
}

// trimTrailingZeros removes any trailing zeros after the decimal point
// (except the one immediately following)
func (f Float) trimTrailingZeros(s string) string {
	if !f.TrimTrailingZeroes {
		return s
	}

	r := []rune(s)
	postPointIdx := strings.LastIndex(s, ".")
	for i := len(s) - 1; i > postPointIdx+1; i-- {
		if r[i] == '0' {
			r[i] = ' '
		} else {
			break
		}
	}
	return string(r)
}

// Formatted returns the value formatted as a float
func (f *Float) Formatted(v any) string {
	if f.IgnoreNil && v == nil {
		return ""
	}

	format := f.makeFormat(v)

	if ok, str := f.Zeroes.GetZeroStr(f.Prec, v); ok {
		return fmt.Sprintf("%.*s", f.Width(), str)
	}
	return f.trimTrailingZeros(fmt.Sprintf(format, f.Prec, v))
}

// Width returns the intended width of the value. An invalid width or one
// incompatible with the given precision is ignored
func (f Float) Width() uint {
	var minWidth uint = 1
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

// Check returns a non-nil error if the Formatter has an invalid Verb
func (f Float) Check() error {
	switch f.Verb {
	case 0, 'f', 'F', 'e', 'E', 'g', 'G', 'x', 'X':
	default:
		return fmt.Errorf("%T: bad Format verb: %q", f, f.Verb)
	}

	return nil
}
