package rptmaker

import (
	"slices"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestCol_colNames(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		colsToAdd   []ColsAddInfo
		checkFunc   func(ColInfo[P, T]) bool
		expColNames []string
	}{
		{
			ID:        testhelper.MkID("no columns"),
			checkFunc: func(_ ColInfo[P, T]) bool { return true },
		},
		{
			ID: testhelper.MkID("3 columns, return all columns"),
			colsToAdd: []ColsAddInfo{
				{CID: RepColName, CI: ReportableCI},
				{CID: SortColName, CI: SortableCI},
				{CID: ColName, CI: SAndRCI},
			},
			checkFunc: func(_ ColInfo[P, T]) bool { return true },
			expColNames: []string{
				string(RepColName),
				string(SortColName),
				string(ColName),
			},
		},
		{
			ID: testhelper.MkID("3 columns, return no columns"),
			colsToAdd: []ColsAddInfo{
				{CID: RepColName, CI: ReportableCI},
				{CID: SortColName, CI: SortableCI},
				{CID: ColName, CI: SAndRCI},
			},
			checkFunc: func(_ ColInfo[P, T]) bool { return false },
		},
		{
			ID: testhelper.MkID("3 columns, return reportable columns"),
			colsToAdd: []ColsAddInfo{
				{CID: RepColName, CI: ReportableCI},
				{CID: SortColName, CI: SortableCI},
				{CID: ColName, CI: SAndRCI},
			},
			checkFunc: ColInfo[P, T].IsReportable,
			expColNames: []string{
				string(RepColName),
				string(ColName),
			},
		},
		{
			ID: testhelper.MkID("3 columns, return sortable columns"),
			colsToAdd: []ColsAddInfo{
				{CID: RepColName, CI: ReportableCI},
				{CID: SortColName, CI: SortableCI},
				{CID: ColName, CI: SAndRCI},
			},
			checkFunc: ColInfo[P, T].IsSortable,
			expColNames: []string{
				string(SortColName),
				string(ColName),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			c, err := MakeTestCols(tc.colsToAdd)
			if err != nil {
				t.Log(tc.IDStr())
				t.Fatal("\t: unexpected error: ", err)
			}

			colNames := c.colNames(tc.checkFunc)

			slices.Sort(colNames)
			slices.Sort(tc.expColNames)

			testhelper.DiffStringSlice(t,
				tc.IDStr(), "column names",
				colNames, tc.expColNames)
		})
	}
}
