package rptmaker_test

import (
	"testing"

	"github.com/nickwells/col.mod/v6/col"
	"github.com/nickwells/col.mod/v6/rptmaker"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestColInfo(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		desc     string
		headings []string
		mkCol    rptmaker.ColMkFunc[P]
		colVal   rptmaker.ColValFunc[T]
		cmpVals  rptmaker.ColCmpFunc[T]

		expFullDesc   string
		expReportable bool
		expSortable   bool
	}{
		{
			ID:          testhelper.MkID("all nil"),
			expFullDesc: "This column is unheaded.",
		},
		{
			ID:          testhelper.MkID("has cmpVals"),
			cmpVals:     rptmaker.CmpVals,
			expFullDesc: "This column is unheaded.",
			expSortable: true,
		},
		{
			ID:          testhelper.MkID("has mkCol"),
			mkCol:       func(_ P, _ []string) *col.Col { return nil },
			expFullDesc: "This column is unheaded.",
		},
		{
			ID:          testhelper.MkID("has colVal"),
			colVal:      func(_ T) any { return 0 },
			expFullDesc: "This column is unheaded.",
		},
		{
			ID:            testhelper.MkID("has mkCol and colVal"),
			mkCol:         rptmaker.MkCol,
			colVal:        rptmaker.ColVal,
			expFullDesc:   "This column is unheaded.",
			expReportable: true,
		},
		{
			ID:          testhelper.MkID("has headings"),
			headings:    []string{"Hello", "World"},
			expFullDesc: `This column is headed: "Hello"/"World"`,
		},
		{
			ID:          testhelper.MkID("has desc"),
			desc:        "Description.",
			expFullDesc: "Description. This column is unheaded.",
		},
		{
			ID:       testhelper.MkID("has desc and headings"),
			desc:     "Description.",
			headings: []string{"Hello", "World"},
			expFullDesc: "Description. " +
				`This column is headed: "Hello"/"World"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ci := rptmaker.NewColInfo(
				tc.desc,
				tc.headings,
				tc.mkCol,
				tc.colVal,
				tc.cmpVals)

			testhelper.DiffBool(t,
				tc.IDStr(), "IsSortable",
				ci.IsSortable(), tc.expSortable)
			testhelper.DiffBool(t,
				tc.IDStr(), "IsReportable",
				ci.IsReportable(), tc.expReportable)
			testhelper.DiffString(t,
				tc.IDStr(), "FullDesc",
				ci.FullDesc(), tc.expFullDesc)
			testhelper.DiffString(t,
				tc.IDStr(), "Desc",
				ci.Desc(), tc.desc)
			testhelper.DiffStringSlice(t,
				tc.IDStr(), "Headings",
				ci.Headings(), tc.headings)
		})
	}
}
