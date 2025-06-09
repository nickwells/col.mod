package colfmt

// DupHdlr encapsulates all the parts needed to support the suppression of the
// printing of duplicate values
type DupHdlr struct {
	// SkipDups, if set to true will make values that are the same as the
	// previously printed value print as the empty string (or some other
	// value chosen by the formatter using this).
	SkipDups bool

	previousValue         any
	previousValueRecorded bool
}

// SkipDup returns true if the value, v, is the same as the previous value
// and should be skipped or false otherwise. If duplicates are being skipped
// then the value is recorded for comparison against the next value.
//
// This should be called by the column formatter's Formatted method to decide
// whether or not to replace the value with a new value (typically the empty
// string)
func (dh *DupHdlr) SkipDup(v any) bool {
	if !dh.SkipDups {
		return false
	}

	if dh.previousValueRecorded && v == dh.previousValue {
		return true
	}

	dh.previousValue = v
	dh.previousValueRecorded = true

	return false
}
