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
func (f *FloatZeroHandler) setEpsilon(prec int) {
	if f.epsilon == 0.0 {
		f.epsilon = calcEpsilon(prec)
	}
}

// GetZeroStr calculates the appropriate zero string and returns it with a
// boolean indicating whether it should be used or not (if the value passed
// was actually zero)
func (z *FloatZeroHandler) GetZeroStr(prec int, v interface{}) (bool, string) {
	if z != nil && z.Handle {
		z.setEpsilon(prec)
		f64, ok := getValAsFloat64(v)
		if ok &&
			((prec > 0 && f64 < z.epsilon && f64 > (-1*z.epsilon)) ||
				(prec == 0 && f64 <= z.epsilon && f64 >= (-1*z.epsilon))) {
			return true, z.Replace
		}
	}
	return false, ""
}
