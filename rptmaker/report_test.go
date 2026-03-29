package rptmaker_test

import (
	"errors"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/nickwells/col.mod/v6/col"
	"github.com/nickwells/col.mod/v6/colfmt"
	"github.com/nickwells/col.mod/v6/rptmaker"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

var (
	t1a = T{A: 1, B: "a"}
	t2a = T{A: 2, B: "a"}
	t3a = T{A: 3, B: "a"}
	t3b = T{A: 3, B: "b"}

	ts1 = []T{t2a, t1a, t3a}
	ts2 = []T{t2a, t1a, t3b, t3a}

	ciaName   = rptmaker.ColID("column a")
	cibName   = rptmaker.ColID("column b")
	badCIName = rptmaker.ColID("nonesuch")

	ciaNotSortable = rptmaker.NewColInfo(rptmaker.CIDesc, nil,
		func(_ P, headings []string) *col.Col {
			return col.New(&colfmt.Int{}, headings...)
		},
		func(t T) any { return t.A },
		nil,
	)

	ciaNotReportableNoMkCol = rptmaker.NewColInfo[P, T](rptmaker.CIDesc, nil,
		nil,
		func(t T) any { return t.A },
		func(a, b T) int { return a.A - b.A },
	)

	ciaNotReportableNoColVal = rptmaker.NewColInfo(rptmaker.CIDesc, nil,
		func(_ P, h []string) *col.Col { return col.New(&colfmt.Int{}, h...) },
		nil,
		func(a, b T) int { return a.A - b.A },
	)

	ciaNotReportableNilCol = rptmaker.NewColInfo(rptmaker.CIDesc, nil,
		func(_ P, _ []string) *col.Col { return nil },
		func(t T) any { return t.A },
		func(a, b T) int { return a.A - b.A },
	)

	cia = rptmaker.NewColInfo(rptmaker.CIDesc, []string{"column", "A"},
		func(_ P, h []string) *col.Col { return col.New(&colfmt.Int{}, h...) },
		func(t T) any { return t.A },
		func(a, b T) int { return a.A - b.A },
	)

	cib = rptmaker.NewColInfo(rptmaker.CIDesc, []string{"column", "B"},
		func(_ P, h []string) *col.Col {
			return col.New(&colfmt.String{}, h...)
		},
		func(t T) any { return t.B },
		func(a, b T) int { return strings.Compare(a.B, b.B) },
	)
)

func TestReport_MakeReport(t *testing.T) {
	const errIntro = "cannot create the Report:"

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		colsToAdd []ColsAddInfo
		repCols   []rptmaker.ColID
		hOpts     []col.HdrOptionFunc
	}{
		{
			ID:     testhelper.MkID("no columns"),
			ExpErr: testhelper.MkExpErr(errIntro + " no columns were given"),
		},
		{
			ID: testhelper.MkID("one column, bad header options"),
			ExpErr: testhelper.MkExpErr(errIntro +
				" cannot create the report header: bad header"),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: cia},
				{CID: cibName, CI: cib},
			},
			repCols: []rptmaker.ColID{ciaName},
			hOpts: []col.HdrOptionFunc{
				func(_ *col.Header) error {
					return errors.New("bad header")
				},
			},
		},
		{
			ID: testhelper.MkID("one column, not reportable, missing mkCol"),
			ExpErr: testhelper.MkExpErr(errIntro +
				` cannot GetReportableColInfo:` +
				` column: "column a": has no "mkCol" function` +
				` (it is not reportable)`),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: ciaNotReportableNoMkCol},
				{CID: cibName, CI: cib},
			},
			repCols: []rptmaker.ColID{ciaName},
		},
		{
			ID: testhelper.MkID("one column, not reportable, missing colVal"),
			ExpErr: testhelper.MkExpErr(errIntro +
				` cannot GetReportableColInfo:` +
				` column: "column a": has no "colVal" function` +
				` (it is not reportable)`),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: ciaNotReportableNoColVal},
				{CID: cibName, CI: cib},
			},
			repCols: []rptmaker.ColID{ciaName},
		},
		{
			ID: testhelper.MkID("one column, not reportable, nil col"),
			ExpErr: testhelper.MkExpErr(errIntro +
				` the mkCol function for "column a" returned a nil column`),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: ciaNotReportableNilCol},
				{CID: cibName, CI: cib},
			},
			repCols: []rptmaker.ColID{ciaName},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			b, err := rptmaker.MakeTestCols(tc.colsToAdd)
			if err != nil {
				t.Log(tc.IDStr())
				t.Fatal("\t: unexpected error making Cols: ", err)
			}

			r, err := b.MakeReport(P{},
				os.Stdout,
				tc.repCols,
				tc.hOpts...)

			testhelper.CheckExpErr(t, err, tc)

			if err == nil {
				if r == nil {
					t.Log(tc.IDStr())
					t.Error("\t: no error but no report was returned")
				}
			}
		})
	}
}

func TestReport_MkCmpFunc(t *testing.T) {
	const errIntro = "cannot make the comparison function: "

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		colsToAdd     []ColsAddInfo
		repCols       []rptmaker.ColID
		sortCols      []rptmaker.SortColumn
		data          []T
		expSortedData []T
	}{
		{
			ID: testhelper.MkID("bad sort cols, not found"),
			ExpErr: testhelper.MkExpErr(
				errIntro + `column: "nonesuch": not found`),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: cia},
				{CID: cibName, CI: cib},
			},
			repCols: []rptmaker.ColID{ciaName},
			sortCols: []rptmaker.SortColumn{
				{ID: badCIName},
			},
			data:          ts1,
			expSortedData: ts1,
		},
		{
			ID: testhelper.MkID("bad sort cols, nil CmpVals"),
			ExpErr: testhelper.MkExpErr(
				errIntro +
					`column: "column a": has no "cmpVals" function` +
					" (it is not sortable)"),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: ciaNotSortable},
				{CID: cibName, CI: cib},
			},
			repCols:       []rptmaker.ColID{ciaName},
			sortCols:      []rptmaker.SortColumn{{ID: ciaName}},
			data:          ts1,
			expSortedData: ts1,
		},
		{
			ID: testhelper.MkID("no sort cols exp no change"),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: cia},
				{CID: cibName, CI: cib},
			},
			repCols:       []rptmaker.ColID{ciaName},
			data:          ts1,
			expSortedData: ts1,
		},
		{
			ID: testhelper.MkID("sort by a, ascending"),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: cia},
				{CID: cibName, CI: cib},
			},
			repCols:       []rptmaker.ColID{ciaName},
			sortCols:      []rptmaker.SortColumn{{ID: ciaName}},
			data:          ts1,
			expSortedData: []T{t1a, t2a, t3a},
		},
		{
			ID: testhelper.MkID("sort by a, descending"),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: cia},
				{CID: cibName, CI: cib},
			},
			repCols:       []rptmaker.ColID{ciaName},
			sortCols:      []rptmaker.SortColumn{{ID: ciaName, Backwards: true}},
			data:          ts1,
			expSortedData: []T{t3a, t2a, t1a},
		},
		{
			ID: testhelper.MkID("sort by a, ascending and b, ascending"),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: cia},
				{CID: cibName, CI: cib},
			},
			repCols: []rptmaker.ColID{ciaName},
			sortCols: []rptmaker.SortColumn{
				{ID: ciaName},
				{ID: cibName},
			},
			data:          ts2,
			expSortedData: []T{t1a, t2a, t3a, t3b},
		},
		{
			ID: testhelper.MkID("sort by a, ascending and b, descending"),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: cia},
				{CID: cibName, CI: cib},
			},
			repCols: []rptmaker.ColID{ciaName},
			sortCols: []rptmaker.SortColumn{
				{ID: ciaName},
				{ID: cibName, Backwards: true},
			},
			data:          ts2,
			expSortedData: []T{t1a, t2a, t3b, t3a},
		},
		{
			ID: testhelper.MkID("sort by a, descending and b, ascending"),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: cia},
				{CID: cibName, CI: cib},
			},
			repCols: []rptmaker.ColID{ciaName},
			sortCols: []rptmaker.SortColumn{
				{ID: ciaName, Backwards: true},
				{ID: cibName},
			},
			data:          ts2,
			expSortedData: []T{t3a, t3b, t2a, t1a},
		},
		{
			ID: testhelper.MkID("sort by a, descending and b, descending"),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: cia},
				{CID: cibName, CI: cib},
			},
			repCols: []rptmaker.ColID{ciaName},
			sortCols: []rptmaker.SortColumn{
				{ID: ciaName, Backwards: true},
				{ID: cibName, Backwards: true},
			},
			data:          ts2,
			expSortedData: []T{t3b, t3a, t2a, t1a},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			b, err := rptmaker.MakeTestCols(tc.colsToAdd)
			if err != nil {
				t.Log(tc.IDStr())
				t.Fatal("\t: unexpected error making Cols: ", err)
			}

			r, err := (b).MakeReport(P{}, os.Stdout, tc.repCols)
			if err != nil {
				t.Log(tc.IDStr())
				t.Fatal("\t: unexpected error making Report: ", err)
			}

			cf, err := r.MkCmpFunc(tc.sortCols)
			testhelper.CheckExpErr(t, err, tc)

			if err == nil {
				slices.SortFunc(tc.data, cf)
				testhelper.DiffValsReport(t,
					tc.IDStr(), "sorted data",
					tc.data, tc.expSortedData)
			}
		})
	}
}

func TestReport_Print(t *testing.T) {
	const cfErrIntro = "cannot make the comparison function: "

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		colsToAdd []ColsAddInfo
		repCols   []rptmaker.ColID
		sortCols  []rptmaker.SortColumn
		data      []T
		expReport string
	}{
		{
			ID: testhelper.MkID("bad sort cols, not found"),
			ExpErr: testhelper.MkExpErr(
				cfErrIntro + `column: "nonesuch": not found`),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: cia},
				{CID: cibName, CI: cib},
			},
			repCols: []rptmaker.ColID{ciaName},
			sortCols: []rptmaker.SortColumn{
				{ID: badCIName},
			},
		},
		{
			ID: testhelper.MkID("no data, no expected report"),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: cia},
				{CID: cibName, CI: cib},
			},
			repCols: []rptmaker.ColID{ciaName},
			sortCols: []rptmaker.SortColumn{
				{ID: ciaName},
			},
		},
		{
			ID: testhelper.MkID("3 rows of data, 2 columns"),
			colsToAdd: []ColsAddInfo{
				{CID: ciaName, CI: cia},
				{CID: cibName, CI: cib},
			},
			repCols: []rptmaker.ColID{ciaName, cibName},
			sortCols: []rptmaker.SortColumn{
				{ID: ciaName},
			},
			data: ts1,
			expReport: `-column-
   A B  
   = =  
   1 a  
   2 a  
   3 a  
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			b, err := rptmaker.MakeTestCols(tc.colsToAdd)
			if err != nil {
				t.Log(tc.IDStr())
				t.Fatal("\t: unexpected error making Cols: ", err)
			}

			var rptOut strings.Builder

			r, err := (b).MakeReport(P{}, &rptOut, tc.repCols)
			if err != nil {
				t.Log(tc.IDStr())
				t.Fatal("\t: unexpected error making Report: ", err)
			}

			err = r.Print(tc.data, tc.sortCols)
			testhelper.CheckExpErr(t, err, tc)

			if err == nil {
				testhelper.DiffString(t,
					tc.IDStr(), "report",
					rptOut.String(), tc.expReport)
			}
		})
	}
}
