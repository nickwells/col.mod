package colfmt

import (
	"bytes"
	"strings"

	"github.com/nickwells/col.mod/v6/col"
	"github.com/nickwells/twrap.mod/twrap"
)

// WrappedString records the values needed for the formatting of a string
// value.
//
// See [NilHdlr] and [DupHdlr] for the settings that can be given through
// those types.
type WrappedString struct {
	// W gives the width of the block that the string should fit within. This
	// must be set to some value greater than zero.
	W int
	// DupIndicator is the value to show if the value to be shown is the same
	// as the value shown on the previous line. Setting this value without
	// also setting the DupHdlr.SkipDups flag will have no effect. Note that
	// if the DupIndicator is too long to fit in the column it will be
	// truncated according to the setting of the W value.
	DupIndicator string

	NilHdlr
	DupHdlr
}

// Formatted returns the value formatted as a string. The string is wrapped
// to a maximum length of WrappedString.W and any trailing newlines are
// trimmed
func (f *WrappedString) Formatted(v any) string {
	if f.SkipNil(v) {
		return ""
	}

	if f.SkipDup(v) {
		return f.DupIndicator
	}

	width := f.W
	if width <= 0 {
		width = 1
	}

	b := bytes.Buffer{}
	twc := twrap.NewTWConfOrPanic(
		twrap.SetTargetLineLen(width),
		twrap.SetMinChars(width),
		twrap.SetWriter(&b))

	twc.Wrap(v.(string), 0)

	return strings.TrimRight(b.String(), "\n")
}

// Width returns the intended width of the value
func (f WrappedString) Width() int {
	if f.W <= 0 {
		return 1
	}

	return f.W
}

// Just returns the justification of the value
func (f WrappedString) Just() col.Justification {
	return col.Left
}

// Check returns a nil error
func (f WrappedString) Check() error {
	return nil
}
