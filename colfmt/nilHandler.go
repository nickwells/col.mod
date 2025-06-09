package colfmt

// NilHdlr encapsulates all the parts needed to support the suppression of the
// printing of duplicate values
type NilHdlr struct {
	// IgnoreNil, if set to true will make nil values print as the empty
	// string (or some other value chosen by the formatter using this).
	IgnoreNil bool
}

// SkipNil returns true if the value, v, is nil and the IgnoreNil flag is
// set.
//
// This should be called by the column formatter's Formatted method to decide
// whether or not to replace the value with a new value (typically the empty
// string)
func (nh *NilHdlr) SkipNil(v any) bool {
	return nh.IgnoreNil && v == nil
}
