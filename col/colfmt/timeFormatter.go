package colfmt

import (
	"fmt"
	"time"

	"github.com/nickwells/col.mod/v3/col"
)

// DfltTimeFormat is the format that will be used by the Time Formatter for
// printing out times if no other format is given. Note that the order of the
// date is Year/Month/Day not Year/Day/Month
const DfltTimeFormat = "2006/01/02 15:04:05.000"

// Time records the values needed for the formatting of a time value
type Time struct {
	W         int
	Format    string
	IgnoreNil bool
}

// Formatted returns the value formatted as a time. If the format string is
// not set then it is set to the DfltTimeFormat.
func (f *Time) Formatted(v any) string {
	if f.IgnoreNil && v == nil {
		return ""
	}

	if f.Format == "" {
		f.Format = DfltTimeFormat
	}

	if t, ok := v.(time.Time); ok {
		return t.Format(f.Format)
	}
	return fmt.Sprintf("Not a time: %v", v)
}

// Width returns the intended width of the value. If it is set to zero then
// the length of the format string is used as a reasonable (but imperfect)
// value. If the format string is not set then it is set to the
// DfltTimeFormat before the width is calculated.
func (f *Time) Width() int {
	if f.W == 0 {
		if f.Format == "" {
			f.Format = DfltTimeFormat
		}
		f.W = len(f.Format)
	}
	return f.W
}

// Just returns the justification of the value
func (f Time) Just() col.Justification {
	return col.Left
}
