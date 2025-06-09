package colfmt_test

import (
	"testing"

	"github.com/nickwells/col.mod/v5/colfmt"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestPctFormatter(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		pf     colfmt.Percent
		val    any
		expStr string
	}{
		{
			ID:     testhelper.MkID("basic"),
			val:    0.123,
			expStr: "12%",
		},
		{
			ID:     testhelper.MkID("basic, pass nil"),
			expStr: "nil",
		},
		{
			ID:     testhelper.MkID("basic, pass float64"),
			val:    float64(0.123),
			expStr: "12%",
		},
		{
			ID:     testhelper.MkID("basic, pass float32"),
			val:    float32(0.123),
			expStr: "12%",
		},
		{
			ID:     testhelper.MkID("basic, pass int"),
			val:    int(123),
			expStr: "12300%",
		},
		{
			ID:     testhelper.MkID("basic, pass non-float / non-int"),
			val:    "not a number",
			expStr: "Numeric value expected (got: string): not a number",
		},
		{
			ID:     testhelper.MkID("ignore nil, pass a value"),
			pf:     colfmt.Percent{IgnoreNil: true},
			val:    1.23,
			expStr: "123%",
		},
		{
			ID:     testhelper.MkID("ignore nil, pass nil"),
			pf:     colfmt.Percent{IgnoreNil: true},
			expStr: "",
		},
		{
			ID:     testhelper.MkID("with precision"),
			pf:     colfmt.Percent{Prec: 2},
			val:    1.2345,
			expStr: "123.45%",
		},
		{
			ID: testhelper.MkID("with zero handling, large (just) value"),
			pf: colfmt.Percent{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
			},
			val:    0.0050000001,
			expStr: "1%",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, borderline value, zero precision"),
			pf: colfmt.Percent{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
			},
			val:    0.005,
			expStr: "ab",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, zero value, zero precision and width"),
			pf: colfmt.Percent{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
				W: 4,
			},
			val:    0.0,
			expStr: "abcd",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, large (just) value, non-zero precision"),
			pf: colfmt.Percent{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
				Prec: 1,
			},
			val:    0.0005,
			expStr: "0.1%",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, small value, non-zero precision"),
			pf: colfmt.Percent{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
				Prec: 1,
			},
			val:    0.0004999,
			expStr: "abcd",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, small -ve value, non-zero precision"),
			pf: colfmt.Percent{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
				Prec: 1,
			},
			val:    -0.0004999,
			expStr: "abcd",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, small value, non-zero precision & width"),
			pf: colfmt.Percent{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
				Prec: 1,
				W:    6,
			},
			val:    0.0004999,
			expStr: "abcd",
		},
		{
			ID: testhelper.MkID("basic, no % sign"),
			pf: colfmt.Percent{
				SuppressPct: true,
			},
			val:    0.123,
			expStr: "12",
		},
	}

	for _, tc := range testCases {
		s := tc.pf.Formatted(tc.val)
		testhelper.DiffString(t, tc.IDStr(), "formatted value", s, tc.expStr)
	}
}

func TestPercentWidth(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		pf       colfmt.Percent
		expWidth uint
	}{
		{
			ID:       testhelper.MkID("default"),
			expWidth: 2,
		},
		{
			ID:       testhelper.MkID("no % sign"),
			pf:       colfmt.Percent{SuppressPct: true},
			expWidth: 1,
		},
		{
			ID:       testhelper.MkID("zero precision, small, non-zero width"),
			pf:       colfmt.Percent{W: 1},
			expWidth: 2,
		},
		{
			ID:       testhelper.MkID("zero precision, large, non-zero width"),
			pf:       colfmt.Percent{W: 6},
			expWidth: 6,
		},
		{
			ID: testhelper.MkID(
				"non-zero precision, small, non-zero width"),
			pf:       colfmt.Percent{W: 1, Prec: 1},
			expWidth: 4,
		},
		{
			ID: testhelper.MkID(
				"non-zero precision, large, non-zero width"),
			pf:       colfmt.Percent{W: 6, Prec: 1},
			expWidth: 6,
		},
	}

	for _, tc := range testCases {
		testhelper.DiffInt(t, tc.IDStr(), "width", tc.pf.Width(), tc.expWidth)
	}
}
