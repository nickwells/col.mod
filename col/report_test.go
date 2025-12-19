package col_test

import (
	"testing"

	"github.com/nickwells/col.mod/v6/col"
	"github.com/nickwells/col.mod/v6/colfmt"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestNewReport(t *testing.T) {
	const (
		badFormatVerb = " has a bad Formatter: colfmt.Int: bad Format verb: '😀'"
		badCol0       = `column[0] (["bad" "formatter"])` + badFormatVerb
	)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		testhelper.ExpPanic
		col1 *col.Col
		cols []*col.Col
	}{
		{
			ID:   testhelper.MkID("ok"),
			col1: col.New(&colfmt.Int{Verb: 'd'}, "good"),
		},
		{
			ID:       testhelper.MkID("bad column"),
			ExpErr:   testhelper.MkExpErr(badCol0),
			ExpPanic: testhelper.MkExpPanic(badCol0),
			col1:     col.New(&colfmt.Int{Verb: '😀'}, "bad", "formatter"),
		},
		{
			ID: testhelper.MkID("multiple bad columns"),
			ExpErr: testhelper.MkExpErr(
				"column[1]",
				"column[2]",
				badFormatVerb),
			ExpPanic: testhelper.MkExpPanic(
				"column[1]",
				"column[2]",
				badFormatVerb),
			col1: col.New(&colfmt.Int{}, "good", "formatter"),
			cols: []*col.Col{
				col.New(&colfmt.Int{Verb: '😀'}, "bad", "formatter"),
				col.New(&colfmt.Int{Verb: '😀'}, "another bad", "formatter"),
			},
		},
	}

	for _, tc := range testCases {
		newRpt, err := col.NewReport(nil, nil, tc.col1, tc.cols...)
		testhelper.CheckExpErr(t, err, tc)

		var stdRpt *col.Report

		panicked, panicVal := testhelper.PanicSafe(func() {
			stdRpt = col.StdRpt(tc.col1, tc.cols...)
		})
		testhelper.CheckExpPanicError(t, panicked, panicVal, tc)

		err = testhelper.DiffVals(newRpt, stdRpt)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: the StdRpt and NewReport values differ: %s", err)
		}
	}
}
