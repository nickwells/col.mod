package colfmt

import "math"

// getValAsFloat64 converts the interface value into a float64 if
// possible. It will set the boolean return value to false if it is not
// possible
func getValAsFloat64(v interface{}) (float64, bool) {
	var f64 float64
	var f32 float32
	var ok bool
	f64, ok = v.(float64)
	if !ok {
		f32, ok = v.(float32)
		if ok {
			f64 = float64(f32)
		}
	}
	return f64, ok
}

// calcEpsilon calculates the appropriate epsilon value for the given precision
func calcEpsilon(prec int) float64 {
	scale := prec
	if scale < 0 {
		scale = 0
	}
	scale++
	scale *= -1
	return 5.0 * math.Pow10(scale)
}
