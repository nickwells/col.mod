package colfmt

import (
	"bytes"
	"strings"

	"github.com/nickwells/col.mod/v2/col"
	"github.com/nickwells/twrap.mod/twrap"
)

// WrappedString records the values needed for the formatting of a
// string value.
type WrappedString struct {
	// W gives the width of the block that the string should fit within. This
	// must be set to some non-zero value.
	W int
	// IgnoreNil, if set to true will make nil values print as the empty string
	IgnoreNil bool
}

// Formatted returns the value formatted as a string. The string is wrapped to a maximum length of WrappedString.W and any trailing newlines are trimmed
func (f WrappedString) Formatted(v interface{}) string {
	if f.IgnoreNil && v == nil {
		return ""
	}

	var b bytes.Buffer
	twc := twrap.NewTWConfOrPanic(
		twrap.SetTargetLineLen(f.W),
		twrap.SetMinChars(f.W),
		twrap.SetWriter(&b))

	twc.Wrap(v.(string), 0)

	return strings.TrimRight(b.String(), "\n")
}

// Width returns the intended width of the value
func (f WrappedString) Width() int {
	return f.W
}

// Just returns the justification of the value
func (f WrappedString) Just() col.Justification {
	return col.Left
}
