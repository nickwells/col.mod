package colfmt_test

import (
	"testing"

	"github.com/nickwells/col.mod/v4/colfmt"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestIntFormatter(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		intF   colfmt.Int
		val    any
		expStr string
	}{
		{
			ID:     testhelper.MkID("basic"),
			val:    123,
			expStr: "123",
		},
		{
			ID:     testhelper.MkID("basic, pass nil"),
			expStr: "%!d(<nil>)",
		},
		{
			ID:     testhelper.MkID("ignore nil, pass a value"),
			intF:   colfmt.Int{NilHdlr: colfmt.NilHdlr{IgnoreNil: true}},
			val:    1,
			expStr: "1",
		},
		{
			ID:     testhelper.MkID("ignore nil, pass nil"),
			intF:   colfmt.Int{NilHdlr: colfmt.NilHdlr{IgnoreNil: true}},
			expStr: "",
		},
		{
			ID:     testhelper.MkID("with no zero handling, zero value"),
			val:    0,
			expStr: "0",
		},
		{
			ID: testhelper.MkID("with zero handling, zero value"),
			intF: colfmt.Int{
				HandleZeroes:    true,
				ZeroReplacement: "abcd",
			},
			val:    0,
			expStr: "a",
		},
		{
			ID: testhelper.MkID("with zero handling, non-zero value"),
			intF: colfmt.Int{
				HandleZeroes:    true,
				ZeroReplacement: "abcd",
			},
			val:    1,
			expStr: "1",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, zero value, non-zero width"),
			intF: colfmt.Int{
				HandleZeroes:    true,
				ZeroReplacement: "abcd",
				W:               6,
			},
			val:    0,
			expStr: "abcd",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, zero value, as int"),
			intF: colfmt.Int{
				HandleZeroes:    true,
				ZeroReplacement: "",
			},
			val:    int(0),
			expStr: "",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, zero value, as int32"),
			intF: colfmt.Int{
				HandleZeroes:    true,
				ZeroReplacement: "",
			},
			val:    int32(0),
			expStr: "",
		},
		{
			ID: testhelper.MkID(
				"with zero handling, zero value, as int64"),
			intF: colfmt.Int{
				HandleZeroes:    true,
				ZeroReplacement: "",
			},
			val:    int64(0),
			expStr: "",
		},
	}

	for _, tc := range testCases {
		s := tc.intF.Formatted(tc.val)
		testhelper.DiffString(t, tc.IDStr(), "formatted value", s, tc.expStr)
	}
}

func TestIntWidth(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		intF     colfmt.Int
		expWidth uint
	}{
		{
			ID:       testhelper.MkID("zero width"),
			expWidth: 1,
		},
		{
			ID: testhelper.MkID("width > 0"),
			intF: colfmt.Int{
				W: 9,
			},
			expWidth: 9,
		},
	}

	for _, tc := range testCases {
		testhelper.DiffInt(t, tc.IDStr(), "width", tc.intF.Width(), tc.expWidth)
	}
}
