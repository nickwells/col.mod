package colfmt

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/nickwells/col.mod/v5/col"
	"github.com/nickwells/twrap.mod/twrap"
)

// WrappedString records the values needed for the formatting of a string
// value.
//
// See [NilHdlr] and [DupHdlr] for the settings that can be given through
// those types.
type WrappedString struct {
	// W gives the width of the block that the string should fit within. This
	// must be set to some non-zero value.
	W uint
	// DupIndicator is the value to show if the value to be shown is the same
	// as the value shown on the previous line. Setting this value without
	// also setting the DupHdlr.SkipDups flag will have no effect. Note that
	// if the DupIndicator is too long to fit in the column it will be
	// truncated according to the settings of the W and MaxW values.
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

	if f.W > math.MaxInt {
		panic(fmt.Errorf(
			"the width (%d) is too big, the maximum value is %d",
			f.W, math.MaxInt))
	}

	width := int(f.W) //nolint:gosec
	if width == 0 {
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
func (f WrappedString) Width() uint {
	if f.W == 0 {
		return 1
	}

	return f.W
}

// Just returns the justification of the value
func (f WrappedString) Just() col.Justification {
	return col.Left
}

// Check returns a non-nil error if the parameters are invalid
func (f WrappedString) Check() error {
	if f.W == 0 {
		return fmt.Errorf("the width (%d) must be > 0", f.W)
	}

	return nil
}
