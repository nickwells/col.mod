package colfmt_test

import (
	"testing"

	"github.com/nickwells/col.mod/v2/col/colfmt"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestIntFormatter(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		intF     colfmt.Int
		val      interface{}
		expected string
	}{
		{
			ID:       testhelper.MkID("basic"),
			val:      123,
			expected: "123",
		},
		{
			ID:       testhelper.MkID("basic, pass nil"),
			expected: "%!d(<nil>)",
		},
		{
			ID:       testhelper.MkID("ignore nil, pass a value"),
			intF:     colfmt.Int{IgnoreNil: true},
			val:      1,
			expected: "1",
		},
		{
			ID:       testhelper.MkID("ignore nil, pass nil"),
			intF:     colfmt.Int{IgnoreNil: true},
			expected: "",
		},
		{
			ID:       testhelper.MkID("with no zero handling, zero value"),
			val:      0,
			expected: "0",
		},
		{
			ID: testhelper.MkID("with zero handling, zero value"),
			intF: colfmt.Int{
				HandleZeroes:    true,
				ZeroReplacement: "abcd",
			},
			val:      0,
			expected: "a",
		},
		{
			ID: testhelper.MkID("with zero handling, non-zero value"),
			intF: colfmt.Int{
				HandleZeroes:    true,
				ZeroReplacement: "abcd",
			},
			val:      1,
			expected: "1",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, zero value, non-zero width"),
			intF: colfmt.Int{
				HandleZeroes:    true,
				ZeroReplacement: "abcd",
				W:               6,
			},
			val:      0,
			expected: "abcd",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, zero value, as int"),
			intF: colfmt.Int{
				HandleZeroes:    true,
				ZeroReplacement: "",
			},
			val:      int(0),
			expected: "",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, zero value, as int32"),
			intF: colfmt.Int{
				HandleZeroes:    true,
				ZeroReplacement: "",
			},
			val:      int32(0),
			expected: "",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, zero value, as int64"),
			intF: colfmt.Int{
				HandleZeroes:    true,
				ZeroReplacement: "",
			},
			val:      int64(0),
			expected: "",
		},
	}

	for _, tc := range testCases {
		s := tc.intF.Formatted(tc.val)
		if s != tc.expected {
			t.Log(tc.IDStr())
			t.Log("\t: expected:", tc.expected)
			t.Log("\t:      got:", s)
			t.Errorf("\t: badly formatted value\n")
		}
	}
}

func TestIntWidth(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		intF     colfmt.Int
		expected int
	}{
		{
			ID:       testhelper.MkID("zero width"),
			expected: 1,
		},
		{
			ID: testhelper.MkID("-ve width"),
			intF: colfmt.Int{
				W: -1,
			},
			expected: 1,
		},
		{
			ID: testhelper.MkID("width > 0"),
			intF: colfmt.Int{
				W: 9,
			},
			expected: 9,
		},
	}

	for _, tc := range testCases {
		w := tc.intF.Width()
		if w != tc.expected {
			t.Log(tc.IDStr())
			t.Log("\t: expected:", tc.expected)
			t.Log("\t:      got:", w)
			t.Errorf("\t: Width\n")
		}
	}
}
