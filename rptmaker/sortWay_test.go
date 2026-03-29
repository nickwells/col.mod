package rptmaker_test

import (
	"testing"

	"github.com/nickwells/col.mod/v6/rptmaker"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestMakeSortColumn(t *testing.T) {
	testColID := rptmaker.ColID("test")

	testCases := []struct {
		testhelper.ID
		tags          []rptmaker.SortWay
		expSortColumn rptmaker.SortColumn
	}{
		{
			ID: testhelper.MkID("empty tags"),
			expSortColumn: rptmaker.SortColumn{
				ID:        testColID,
				Backwards: false,
			},
		},
		{
			ID:   testhelper.MkID("one tag - Forwards"),
			tags: []rptmaker.SortWay{rptmaker.Forwards},
			expSortColumn: rptmaker.SortColumn{
				ID:        testColID,
				Backwards: false,
			},
		},
		{
			ID:   testhelper.MkID("one tag - Backwards"),
			tags: []rptmaker.SortWay{rptmaker.Backwards},
			expSortColumn: rptmaker.SortColumn{
				ID:        testColID,
				Backwards: true,
			},
		},
		{
			ID:   testhelper.MkID("two tags - Forwards and Backwards"),
			tags: []rptmaker.SortWay{rptmaker.Backwards, rptmaker.Forwards},
			expSortColumn: rptmaker.SortColumn{
				ID:        testColID,
				Backwards: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actSC := rptmaker.MakeSortColumn(testColID, tc.tags)
			testhelper.DiffValsReport(t,
				tc.IDStr(), "sort column",
				actSC, tc.expSortColumn)
		})
	}
}
