package colfmt

import (
	"fmt"

	"github.com/nickwells/col.mod/v4/col"
)

// Int records the values needed for the formatting of an int value.
//
// See [NilHdlr] and [DupHdlr] for the settings that can be given through
// those types.
type Int struct {
	// W gives the minimum space to be taken by the formatted value
	W uint
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

	format string

	NilHdlr
	DupHdlr
}

// makeFormat sets the format string to be used to format the value. It uses
// the Verb to construct the format string.
func (f *Int) makeFormat() {
	if f.format == "" {
		switch f.Verb {
		case 0:
			f.format = "%d"
		case 'b', 'c', 'd', 'o', 'O', 'q', 'x', 'X', 'U':
			f.format = "%" + string(f.Verb)
		default:
			panic(fmt.Errorf("%T: bad Format verb: %q", f, f.Verb))
		}
	}
}

// isZero tests the interface value to see if it is a zero integer
//
//nolint:cyclop
func isZero(v any) bool {
	// Sadly we need this long list of switch types 'cause otherwise the 'i'
	// value remains as an 'any' and the test fails
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
func (f *Int) Formatted(v any) string {
	if f.SkipNil(v) {
		return ""
	}

	if f.SkipDup(v) {
		return ""
	}

	if f.HandleZeroes {
		if isZero(v) {
			return fmt.Sprintf("%.*s", f.Width(), f.ZeroReplacement)
		}
	}

	f.makeFormat()

	return fmt.Sprintf(f.format, v)
}

// Width returns the intended width of the value
func (f Int) Width() uint {
	if f.W == 0 {
		return 1
	}

	return f.W
}

// Just returns the justification of the value
func (f Int) Just() col.Justification {
	return col.Right
}

// Check returns a non-nil error if the Verb is invalid
func (f Int) Check() error {
	switch f.Verb {
	case 0, 'b', 'c', 'd', 'o', 'O', 'q', 'x', 'X', 'U':
	default:
		return fmt.Errorf("%T: bad Format verb: %q", f, f.Verb)
	}

	return nil
}
