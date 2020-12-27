package colfmt

import (
	"fmt"

	"github.com/nickwells/col.mod/v2/col"
)

// Int records the values needed for the formatting of an int value
type Int struct {
	// W gives the minimum space to be taken by the formatted value
	W int
	// IgnoreNil, if set to true will make nil values print as the empty string
	IgnoreNil bool
	// HandleZeroes, if set to true will check if the value to be printed is
	// zero and if so it will print the ZeroReplacement string instead. The
	// string printed will not be wider than the minimum space given by the W
	// value
	HandleZeroes bool
	// ZeroReplacement is the value to be printed for zero if HandleZeroes is
	// true
	ZeroReplacement string
	// Verb specifies the formatting verb. If left unset it will use
	// 'd'. There will be a panic if it is not one of 'bcdoOqxXU'
	Verb rune
}

// makeFormat returns a format string to be used to format the value. It uses
// the Verb to construct the format string.
func (f Int) makeFormat() string {
	switch f.Verb {
	case 0:
		return "%d"
	case 'b', 'c', 'd', 'o', 'O', 'q', 'x', 'X', 'U':
		return "%" + string(f.Verb)
	default:
		panic(fmt.Errorf("%T: bad Format verb: %q", f, f.Verb))
	}
}

// isZero tests the interface value to see if it is a zero integer
func isZero(v interface{}) bool { // nolint: gocyclo
	switch i := v.(type) {
	case int64:
		return i == 0
	case int32:
		return i == 0
	case int16:
		return i == 0
	case int8:
		return i == 0
	case int:
		return i == 0
	case uint64:
		return i == 0
	case uint32:
		return i == 0
	case uint16:
		return i == 0
	case uint8:
		return i == 0
	default:
		return false
	}
}

// Formatted returns the value formatted as an int
func (f Int) Formatted(v interface{}) string {
	if f.IgnoreNil && v == nil {
		return ""
	}

	if f.HandleZeroes {
		if isZero(v) {
			return fmt.Sprintf("%.*s", f.Width(), f.ZeroReplacement)
		}
	}

	format := f.makeFormat()
	return fmt.Sprintf(format, v)
}

// Width returns the intended width of the value
func (f Int) Width() int {
	if f.W <= 0 {
		return 1
	}
	return f.W
}

// Just returns the justification of the value
func (f Int) Just() col.Justification {
	return col.Right
}
