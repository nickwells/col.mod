package colfmt_test

import (
	"testing"

	"github.com/nickwells/col.mod/v2/col/colfmt"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestFloatFormatter(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		ff       colfmt.Float
		val      interface{}
		expected string
	}{
		{
			ID:       testhelper.MkID("basic"),
			val:      1.23,
			expected: "1",
		},
		{
			ID:       testhelper.MkID("basic, pass nil"),
			expected: "%!f(<nil>)",
		},
		{
			ID:       testhelper.MkID("ignore nil, pass a value"),
			ff:       colfmt.Float{IgnoreNil: true},
			val:      1.23,
			expected: "1",
		},
		{
			ID:       testhelper.MkID("ignore nil, pass nil"),
			ff:       colfmt.Float{IgnoreNil: true},
			expected: "",
		},
		{
			ID:       testhelper.MkID("with precision"),
			ff:       colfmt.Float{Prec: 2},
			val:      1.2345,
			expected: "1.23",
		},
		{
			ID:       testhelper.MkID("with bad precision"),
			ff:       colfmt.Float{Prec: -1},
			val:      1.2345,
			expected: "1",
		},
		{
			ID: testhelper.MkID("with zero handling, large (just) value"),
			ff: colfmt.Float{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
			},
			val:      0.50000001,
			expected: "1",
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
			val:      0.5,
			expected: "a",
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
			val:      0.05,
			expected: "0.1",
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
			val:      0.04999,
			expected: "abc",
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
			val:      -0.04999,
			expected: "abc",
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
			val:      0.04999,
			expected: "abcd",
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
			val:      float64(0.0),
			expected: "",
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
			val:      float32(0.0),
			expected: "",
		},
	}

	for _, tc := range testCases {
		s := tc.ff.Formatted(tc.val)
		if s != tc.expected {
			t.Log(tc.IDStr())
			t.Log("\t: expected:", tc.expected)
			t.Log("\t:      got:", s)
			t.Errorf("\t: badly formatted value\n")
		}
	}
}

func TestFloatWidth(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		ff       colfmt.Float
		expected int
	}{
		{
			ID:       testhelper.MkID("default"),
			expected: 1,
		},
		{
			ID:       testhelper.MkID("zero precision, non-zero width"),
			ff:       colfmt.Float{W: 3},
			expected: 3,
		},
		{
			ID:       testhelper.MkID("non-zero precision, non-zero width"),
			ff:       colfmt.Float{W: 3, Prec: 2},
			expected: 4,
		},
	}

	for _, tc := range testCases {
		w := tc.ff.Width()
		if w != tc.expected {
			t.Log(tc.IDStr())
			t.Log("\t: expected:", tc.expected)
			t.Log("\t:      got:", w)
			t.Errorf("\t: Width\n")
		}
	}
}
