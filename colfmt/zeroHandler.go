package colfmt

// FloatZeroHandler is a mixin for handling zero values for columns taking
// floats
type FloatZeroHandler struct {
	// Handle, if set to true will check if the value to be printed is zero
	// (or closer than the given precision can reveal) and if so it will
	// print the Replacement string instead. The string printed will not be
	// wider than the minimum space indicated by the given width
	Handle bool
	// Replace is the value to be printed for zero if Handle is true
	Replace string
	epsilon float64
}

// setEpsilon sets the epsilon value if it hasn't already been set
func (fzh *FloatZeroHandler) setEpsilon(prec int) {
	if fzh.epsilon == 0.0 {
		fzh.epsilon = calcEpsilon(prec)
	}
}

// GetZeroStr calculates the appropriate zero string and returns it with a
// boolean indicating whether it should be used or not (if the value passed
// was actually zero)
func (fzh *FloatZeroHandler) GetZeroStr(prec int, v any) (bool, string) {
	if fzh != nil && fzh.Handle {
		fzh.setEpsilon(prec)

		f64, ok := getValAsFloat64(v)
		if ok &&
			((prec > 0 && f64 < fzh.epsilon && f64 > (-1*fzh.epsilon)) ||
				(prec == 0 && f64 <= fzh.epsilon && f64 >= (-1*fzh.epsilon))) {
			return true, fzh.Replace
		}
	}

	return false, ""
}
