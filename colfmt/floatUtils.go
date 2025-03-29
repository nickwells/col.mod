package colfmt

import (
	"fmt"
	"math"
)

// getValAsFloat64 converts the interface value into a float64 if
// possible. It will set the boolean return value to false if it is not
// possible
func getValAsFloat64(v any) (float64, bool) {
	if f64, ok := v.(float64); ok {
		return f64, true
	}

	if f32, ok := v.(float32); ok {
		return float64(f32), true
	}

	return 0.0, false
}

// calcEpsilon calculates the appropriate epsilon value for the given precision
func calcEpsilon(prec uint) float64 {
	if prec > math.MaxInt-1 {
		panic(fmt.Errorf(
			"the precision value (%d) is too big, the maximum value is %d",
			prec, math.MaxInt-1))
	}

	scale := int(prec) + 1 //nolint:gosec
	scale *= -1

	return 5.0 * math.Pow10(scale) //nolint:mnd
}
