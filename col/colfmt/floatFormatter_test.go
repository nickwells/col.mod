package colfmt_test

import (
	"testing"

	"github.com/nickwells/col.mod/v3/col/colfmt"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestFloatFormatter(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		ff     colfmt.Float
		val    interface{}
		expStr string
	}{
		{
			ID:     testhelper.MkID("basic"),
			val:    1.23,
			expStr: "1",
		},
		{
			ID:     testhelper.MkID("basic, pass nil"),
			expStr: "%!f(<nil>)",
		},
		{
			ID:     testhelper.MkID("ignore nil, pass a value"),
			ff:     colfmt.Float{IgnoreNil: true},
			val:    1.23,
			expStr: "1",
		},
		{
			ID:     testhelper.MkID("ignore nil, pass nil"),
			ff:     colfmt.Float{IgnoreNil: true},
			expStr: "",
		},
		{
			ID:     testhelper.MkID("with precision"),
			ff:     colfmt.Float{Prec: 2},
			val:    1.2345,
			expStr: "1.23",
		},
		{
			ID:     testhelper.MkID("with bad precision"),
			ff:     colfmt.Float{Prec: -1},
			val:    1.2345,
			expStr: "1",
		},
		{
			ID: testhelper.MkID("with zero handling, large (just) value"),
			ff: colfmt.Float{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
			},
			val:    0.50000001,
			expStr: "1",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, borderline value, zero precision"),
			ff: colfmt.Float{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
			},
			val:    0.5,
			expStr: "a",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, large (just) value, non-zero precision"),
			ff: colfmt.Float{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
				Prec: 1,
			},
			val:    0.05,
			expStr: "0.1",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, small value, non-zero precision"),
			ff: colfmt.Float{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
				Prec: 1,
			},
			val:    0.04999,
			expStr: "abc",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, small -ve value, non-zero precision"),
			ff: colfmt.Float{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
				Prec: 1,
			},
			val:    -0.04999,
			expStr: "abc",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, small value, non-zero precision & width"),
			ff: colfmt.Float{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
				Prec: 1,
				W:    6,
			},
			val:    0.04999,
			expStr: "abcd",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, zero value, as float64"),
			ff: colfmt.Float{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "",
				},
			},
			val:    float64(0.0),
			expStr: "",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, zero value, as float32"),
			ff: colfmt.Float{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "",
				},
			},
			val:    float32(0.0),
			expStr: "",
		},
	}

	for _, tc := range testCases {
		s := tc.ff.Formatted(tc.val)
		testhelper.DiffString(t, tc.IDStr(), "formatted value", s, tc.expStr)
	}
}

func TestFloatWidth(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		ff       colfmt.Float
		expWidth int
	}{
		{
			ID:       testhelper.MkID("default"),
			expWidth: 1,
		},
		{
			ID:       testhelper.MkID("0 prec, non-0 width"),
			ff:       colfmt.Float{W: 3},
			expWidth: 3,
		},
		{
			ID:       testhelper.MkID("non-0 prec, non-0 width, too narrow"),
			ff:       colfmt.Float{W: 3, Prec: 2},
			expWidth: 4,
		},
		{
			ID:       testhelper.MkID("non-0 prec, non-0 width, wide enough"),
			ff:       colfmt.Float{W: 5, Prec: 2},
			expWidth: 5,
		},
	}

	for _, tc := range testCases {
		testhelper.DiffInt(t, tc.IDStr(), "width", tc.ff.Width(), tc.expWidth)
	}
}
