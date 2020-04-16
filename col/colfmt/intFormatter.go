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
}

// getValAsInt64 converts the interface value into an int64 if possible. It
// will set the boolean return value to false if it is not possible
func getValAsInt64(v interface{}) (int64, bool) {
	var i64 int64
	var i32 int32
	var i int
	var ok bool

	i64, ok = v.(int64)
	if !ok {
		i32, ok = v.(int32)
		if ok {
			i64 = int64(i32)
		} else {
			i, ok = v.(int)
			if ok {
				i64 = int64(i)
			}
		}
	}
	return i64, ok
}

// Formatted returns the value formatted as an int
func (f Int) Formatted(v interface{}) string {
	if f.IgnoreNil && v == nil {
		return ""
	}

	if f.HandleZeroes {
		i64, ok := getValAsInt64(v)
		if ok && i64 == 0 {
			return fmt.Sprintf("%.*s", f.Width(), f.ZeroReplacement)
		}
	}

	return fmt.Sprintf("%d", v)
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
