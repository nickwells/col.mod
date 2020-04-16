package colfmt_test

import (
	"testing"

	"github.com/nickwells/col.mod/v2/col/colfmt"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestPctFormatter(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		pf       colfmt.Percent
		val      interface{}
		expected string
	}{
		{
			ID:       testhelper.MkID("basic"),
			val:      0.123,
			expected: "12%",
		},
		{
			ID:       testhelper.MkID("basic, pass nil"),
			expected: "nil",
		},
		{
			ID:       testhelper.MkID("basic, pass float64"),
			val:      float64(0.123),
			expected: "12%",
		},
		{
			ID:       testhelper.MkID("basic, pass float32"),
			val:      float32(0.123),
			expected: "12%",
		},
		{
			ID:       testhelper.MkID("basic, pass int"),
			val:      int(123),
			expected: "%!f(int=123)",
		},
		{
			ID:       testhelper.MkID("basic, pass nil"),
			expected: "nil",
		},
		{
			ID:       testhelper.MkID("ignore nil, pass a value"),
			pf:       colfmt.Percent{IgnoreNil: true},
			val:      1.23,
			expected: "123%",
		},
		{
			ID:       testhelper.MkID("ignore nil, pass nil"),
			pf:       colfmt.Percent{IgnoreNil: true},
			expected: "",
		},
		{
			ID:       testhelper.MkID("with precision"),
			pf:       colfmt.Percent{Prec: 2},
			val:      1.2345,
			expected: "123.45%",
		},
		{
			ID:       testhelper.MkID("with bad precision"),
			pf:       colfmt.Percent{Prec: -1},
			val:      1.2345,
			expected: "123%",
		},
		{
			ID: testhelper.MkID("with zero handling, large (just) value"),
			pf: colfmt.Percent{
				Zeroes: &colfmt.FloatZeroHandler{
					Handle:  true,
					Replace: "abcd",
				},
			},
			val:      0.0050000001,
			expected: "1%",
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
			val:      0.005,
			expected: "ab",
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
			val:      0.0,
			expected: "abcd",
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
			val:      0.0005,
			expected: "0.1%",
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
			val:      0.0004999,
			expected: "abcd",
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
			val:      -0.0004999,
			expected: "abcd",
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
			val:      0.0004999,
			expected: "abcd",
		},
		{
			ID: testhelper.MkID("basic, no % sign"),
			pf: colfmt.Percent{
				SuppressPct: true,
			},
			val:      0.123,
			expected: "12",
		},
	}

	for _, tc := range testCases {
		s := tc.pf.Formatted(tc.val)
		if s != tc.expected {
			t.Log(tc.IDStr())
			t.Log("\t: expected:", tc.expected)
			t.Log("\t:      got:", s)
			t.Errorf("\t: badly formatted value\n")
		}
	}
}

func TestPercentWidth(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		pf       colfmt.Percent
		expected int
	}{
		{
			ID:       testhelper.MkID("default"),
			expected: 2,
		},
		{
			ID:       testhelper.MkID("no % sign"),
			pf:       colfmt.Percent{SuppressPct: true},
			expected: 1,
		},
		{
			ID:       testhelper.MkID("zero precision, small, non-zero width"),
			pf:       colfmt.Percent{W: 1},
			expected: 2,
		},
		{
			ID:       testhelper.MkID("zero precision, large, non-zero width"),
			pf:       colfmt.Percent{W: 6},
			expected: 6,
		},
		{
			ID: testhelper.MkID(
				"non-zero precision, small, non-zero width"),
			pf:       colfmt.Percent{W: 1, Prec: 1},
			expected: 4,
		},
		{
			ID: testhelper.MkID(
				"non-zero precision, large, non-zero width"),
			pf:       colfmt.Percent{W: 6, Prec: 1},
			expected: 6,
		},
	}

	for _, tc := range testCases {
		w := tc.pf.Width()
		if w != tc.expected {
			t.Log(tc.IDStr())
			t.Log("\t: expected:", tc.expected)
			t.Log("\t:      got:", w)
			t.Errorf("\t: Width\n")
		}
	}
}
