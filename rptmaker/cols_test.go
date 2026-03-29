package rptmaker_test

import (
	"testing"

	"github.com/nickwells/col.mod/v6/rptmaker"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestCols_AddCol(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		colsToAdd []ColsAddInfo
	}{
		{
			ID: testhelper.MkID("no cols, no error"),
		},
		{
			ID: testhelper.MkID("one col, bad ColInfo"),
			ExpErr: testhelper.MkExpErr(
				`column: "hello":` +
					` is neither sortable nor reportable`),
			colsToAdd: []ColsAddInfo{
				{CID: "hello", CI: rptmaker.BadCINoFunc},
			},
		},
		{
			ID: testhelper.MkID("one col, nil ColInfo"),
			ExpErr: testhelper.MkExpErr(
				`cannot Add: column: "hello":` +
					` no ColInfo has been supplied`),
			colsToAdd: []ColsAddInfo{
				{CID: "hello"},
			},
		},
		{
			ID: testhelper.MkID("two cols, col name == alias name"),
			ExpErr: testhelper.MkExpErr(
				`cannot Add: column: "world":` +
					` the name is already used as an alias`),
			colsToAdd: []ColsAddInfo{
				{CID: "hello", CI: rptmaker.SortableCI},
				{
					CommonAlias: "world",
					AliasVals: []rptmaker.ColID{
						"hello",
					},
				},
				{CID: "world", CI: rptmaker.ReportableCI},
			},
		},
		{
			ID: testhelper.MkID("two cols, dup name"),
			ExpErr: testhelper.MkExpErr(
				`cannot Add: column: "hello":` +
					` duplicate: there is already a column with that name`),
			colsToAdd: []ColsAddInfo{
				{CID: "hello", CI: rptmaker.SortableCI},
				{CID: "hello", CI: rptmaker.ReportableCI},
			},
		},
		{
			ID: testhelper.MkID("two cols, bad alias - dup with column"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddAlias: alias: "hello":` +
					` duplicate: there is already a column with that name`),
			colsToAdd: []ColsAddInfo{
				{CID: "hello", CI: rptmaker.SortableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					CommonAlias: "hello",
					AliasVals: []rptmaker.ColID{
						"world",
					},
				},
			},
		},
		{
			ID: testhelper.MkID("two cols, bad alias - dup with alias"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddAlias: alias: "greetings":` +
					` duplicate:` +
					` there is already a common alias with that name`),
			colsToAdd: []ColsAddInfo{
				{CID: "hello", CI: rptmaker.SortableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					CommonAlias: "greetings",
					AliasVals: []rptmaker.ColID{
						"world",
					},
				},
				{
					CommonAlias: "greetings",
					AliasVals: []rptmaker.ColID{
						"world",
					},
				},
			},
		},
		{
			ID: testhelper.MkID(
				"two cols, bad alias - dup with reportable alias"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddAlias: alias: "greetings":` +
					` duplicate:` +
					` there is already a reportable alias with that name`),
			colsToAdd: []ColsAddInfo{
				{CID: "hello", CI: rptmaker.SortableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					ReportableAlias: "greetings",
					AliasVals: []rptmaker.ColID{
						"world",
					},
				},
				{
					CommonAlias: "greetings",
					AliasVals: []rptmaker.ColID{
						"world",
					},
				},
			},
		},
		{
			ID: testhelper.MkID(
				"two cols, bad alias - dup with sortable alias"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddAlias: alias: "greetings":` +
					` duplicate:` +
					` there is already a sortable alias with that name`),
			colsToAdd: []ColsAddInfo{
				{CID: "hello", CI: rptmaker.SortableCI},
				{CID: "world", CI: rptmaker.SortableCI},
				{
					SortableAlias: "greetings",
					AliasVals: []rptmaker.ColID{
						"world",
					},
				},
				{
					CommonAlias: "greetings",
					AliasVals: []rptmaker.ColID{
						"world",
					},
				},
			},
		},
		{
			ID: testhelper.MkID("two cols, bad alias - no matching column"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddAlias: alias: "greetings":` +
					` there is no column called "universe"`),
			colsToAdd: []ColsAddInfo{
				{CID: "hello", CI: rptmaker.SortableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					CommonAlias: "greetings",
					AliasVals: []rptmaker.ColID{
						"universe",
					},
				},
			},
		},
		{
			ID: testhelper.MkID(
				"two cols, bad alias - no matching column, with suggestion"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddAlias: alias: "greetings":` +
					` there is no column called "helloo all"` +
					`, did you mean "hello all"`),
			colsToAdd: []ColsAddInfo{
				{CID: "hello all", CI: rptmaker.SortableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					CommonAlias: "greetings",
					AliasVals: []rptmaker.ColID{
						"helloo all",
					},
				},
			},
		},
		{
			ID: testhelper.MkID(
				"three cols, bad sortable alias" +
					" - no matching column, with suggestion"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddSortableAlias: alias: "greetings":` +
					` there is no column called "helloo all"` +
					`, did you mean "hello all"`),
			colsToAdd: []ColsAddInfo{
				{CID: "hello all", CI: rptmaker.SortableCI},
				{CID: "heloo all", CI: rptmaker.ReportableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					SortableAlias: "greetings",
					AliasVals: []rptmaker.ColID{
						"helloo all",
					},
				},
			},
		},
		{
			ID: testhelper.MkID(
				"three cols, bad sortable alias - column not sortable"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddSortableAlias: alias: "greetings":` +
					` column "` + rptmaker.RepColName + `" is not sortable`),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					SortableAlias: "greetings",
					AliasVals: []rptmaker.ColID{
						rptmaker.RepColName,
					},
				},
			},
		},
		{
			ID: testhelper.MkID(
				"three cols, bad sortable alias - alias name = column name"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddSortableAlias: alias: "` + rptmaker.SortColName + `":` +
					` duplicate: there is already a column with that name`),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					SortableAlias: rptmaker.SortColName,
					AliasVals: []rptmaker.ColID{
						rptmaker.SortColName,
					},
				},
			},
		},
		{
			ID: testhelper.MkID(
				"three cols, bad sortable alias - alias name = common alias"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddSortableAlias: alias: "` + rptmaker.Alias1 + `":` +
					` duplicate:` +
					` there is already a common alias with that name`),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					CommonAlias: rptmaker.Alias1,
					AliasVals: []rptmaker.ColID{
						rptmaker.SortColName,
					},
				},
				{
					SortableAlias: rptmaker.Alias1,
					AliasVals: []rptmaker.ColID{
						rptmaker.SortColName,
					},
				},
			},
		},
		{
			ID: testhelper.MkID(
				"three cols, bad sortable alias - alias name = sortable alias"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddSortableAlias: alias: "` + rptmaker.Alias1 + `":` +
					` duplicate:` +
					` there is already a sortable alias with that name`),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					SortableAlias: rptmaker.Alias1,
					AliasVals: []rptmaker.ColID{
						rptmaker.SortColName,
					},
				},
				{
					SortableAlias: rptmaker.Alias1,
					AliasVals: []rptmaker.ColID{
						rptmaker.SortColName,
					},
				},
			},
		},
		{
			ID: testhelper.MkID(
				"three cols, bad reportable alias" +
					" - no matching column, with suggestion"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddReportableAlias: alias: "greetings":` +
					` there is no column called "helloo all"` +
					`, did you mean "heloo all"`),
			colsToAdd: []ColsAddInfo{
				{CID: "hello all", CI: rptmaker.SortableCI},
				{CID: "heloo all", CI: rptmaker.ReportableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					ReportableAlias: "greetings",
					AliasVals: []rptmaker.ColID{
						"helloo all",
					},
				},
			},
		},
		{
			ID: testhelper.MkID(
				"three cols, bad reportable alias - column not reportable"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddReportableAlias: alias: "greetings":` +
					` column "` + rptmaker.SortColName + `" is not reportable`),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					ReportableAlias: "greetings",
					AliasVals: []rptmaker.ColID{
						rptmaker.SortColName,
					},
				},
			},
		},
		{
			ID: testhelper.MkID(
				"three cols, bad reportable alias - alias name = column name"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddReportableAlias: alias: "` + rptmaker.RepColName + `":` +
					` duplicate: there is already a column with that name`),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					ReportableAlias: rptmaker.RepColName,
					AliasVals: []rptmaker.ColID{
						rptmaker.RepColName,
					},
				},
			},
		},
		{
			ID: testhelper.MkID(
				"three cols, bad reportable alias - alias name = common alias"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddReportableAlias: alias: "` + rptmaker.Alias1 + `":` +
					` duplicate:` +
					` there is already a common alias with that name`),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					CommonAlias: rptmaker.Alias1,
					AliasVals: []rptmaker.ColID{
						rptmaker.RepColName,
					},
				},
				{
					ReportableAlias: rptmaker.Alias1,
					AliasVals: []rptmaker.ColID{
						rptmaker.RepColName,
					},
				},
			},
		},
		{
			ID: testhelper.MkID(
				"three cols, bad reportable alias" +
					" - alias name = reportable alias"),
			ExpErr: testhelper.MkExpErr(
				`cannot AddReportableAlias: alias: "` + rptmaker.Alias1 + `":` +
					` duplicate:` +
					` there is already a reportable alias with that name`),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
				{CID: "world", CI: rptmaker.ReportableCI},
				{
					ReportableAlias: rptmaker.Alias1,
					AliasVals: []rptmaker.ColID{
						rptmaker.RepColName,
					},
				},
				{
					ReportableAlias: rptmaker.Alias1,
					AliasVals: []rptmaker.ColID{
						rptmaker.RepColName,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := rptmaker.MakeTestCols(tc.colsToAdd)
			testhelper.CheckExpErr(t, err, tc)
		})
	}
}

func TestCols_GetColsAndAliases(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		colsToAdd            []ColsAddInfo
		expSortableCols      map[rptmaker.ColID]string
		expReportableCols    map[rptmaker.ColID]string
		expSortableAliases   map[rptmaker.ColID][]rptmaker.ColID
		expReportableAliases map[rptmaker.ColID][]rptmaker.ColID
	}{
		{
			ID: testhelper.MkID("empty cols"),
		},
		{
			ID: testhelper.MkID("just one reportable column"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
			},
			expReportableCols: map[rptmaker.ColID]string{
				rptmaker.RepColName: rptmaker.DescInMap,
			},
		},
		{
			ID: testhelper.MkID("just one sortable column"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
			},
			expSortableCols: map[rptmaker.ColID]string{
				rptmaker.SortColName: rptmaker.DescInMap,
			},
		},
		{
			ID: testhelper.MkID("1 sortable and 1 reportable column"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
			},
			expSortableCols: map[rptmaker.ColID]string{
				rptmaker.SortColName: rptmaker.DescInMap,
			},
			expReportableCols: map[rptmaker.ColID]string{
				rptmaker.RepColName: rptmaker.DescInMap,
			},
		},
		{
			ID: testhelper.MkID("1 sortable, 1 reportable and 1 both column"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
				{CID: rptmaker.SRColName, CI: rptmaker.SAndRCI},
			},
			expSortableCols: map[rptmaker.ColID]string{
				rptmaker.SortColName: rptmaker.DescInMap,
				rptmaker.SRColName:   rptmaker.DescInMap,
			},
			expReportableCols: map[rptmaker.ColID]string{
				rptmaker.RepColName: rptmaker.DescInMap,
				rptmaker.SRColName:  rptmaker.DescInMap,
			},
		},
		{
			ID: testhelper.MkID("1 S'able, 1 R'able and 1 both column aliased"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
				{CID: rptmaker.SRColName, CI: rptmaker.SAndRCI},
				{
					CommonAlias: rptmaker.Alias1,
					AliasVals: []rptmaker.ColID{
						rptmaker.SortColName,
						rptmaker.RepColName,
						rptmaker.SRColName,
					},
				},
			},
			expSortableCols: map[rptmaker.ColID]string{
				rptmaker.SortColName: rptmaker.DescInMap,
				rptmaker.SRColName:   rptmaker.DescInMap,
			},
			expSortableAliases: map[rptmaker.ColID][]rptmaker.ColID{
				rptmaker.Alias1: {rptmaker.SortColName, rptmaker.SRColName},
			},
			expReportableCols: map[rptmaker.ColID]string{
				rptmaker.RepColName: rptmaker.DescInMap,
				rptmaker.SRColName:  rptmaker.DescInMap,
			},
			expReportableAliases: map[rptmaker.ColID][]rptmaker.ColID{
				rptmaker.Alias1: {rptmaker.RepColName, rptmaker.SRColName},
			},
		},
		{
			ID: testhelper.MkID(
				"mixed columns, 1 S'able alias, 1 R'able alias"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
				{CID: rptmaker.SRColName, CI: rptmaker.SAndRCI},
				{
					CommonAlias: rptmaker.Alias1,
					AliasVals: []rptmaker.ColID{
						rptmaker.SortColName,
						rptmaker.RepColName,
						rptmaker.SRColName,
					},
				},
				{
					SortableAlias: rptmaker.SortableAlias1,
					AliasVals: []rptmaker.ColID{
						rptmaker.SortColName,
						rptmaker.SRColName,
					},
				},
				{
					ReportableAlias: rptmaker.ReportableAlias1,
					AliasVals: []rptmaker.ColID{
						rptmaker.RepColName,
						rptmaker.SRColName,
					},
				},
			},
			expSortableCols: map[rptmaker.ColID]string{
				rptmaker.SortColName: rptmaker.DescInMap,
				rptmaker.SRColName:   rptmaker.DescInMap,
			},
			expSortableAliases: map[rptmaker.ColID][]rptmaker.ColID{
				rptmaker.Alias1:         {rptmaker.SortColName, rptmaker.SRColName},
				rptmaker.SortableAlias1: {rptmaker.SortColName, rptmaker.SRColName},
			},
			expReportableCols: map[rptmaker.ColID]string{
				rptmaker.RepColName: rptmaker.DescInMap,
				rptmaker.SRColName:  rptmaker.DescInMap,
			},
			expReportableAliases: map[rptmaker.ColID][]rptmaker.ColID{
				rptmaker.Alias1:           {rptmaker.RepColName, rptmaker.SRColName},
				rptmaker.ReportableAlias1: {rptmaker.RepColName, rptmaker.SRColName},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			c, err := rptmaker.MakeTestCols(tc.colsToAdd)
			if err != nil {
				t.Log(tc.IDStr())
				t.Fatal("\t: unexpected error: ", err)
			}

			sc := c.Sortable()
			testhelper.DiffValsReport(t,
				tc.IDStr(), "sortable columns",
				sc, tc.expSortableCols)

			sa := c.SortableAliases()
			testhelper.DiffValsReport(t,
				tc.IDStr(), "sortable aliases",
				sa, tc.expSortableAliases)

			rc := c.Reportable()
			testhelper.DiffValsReport(t,
				tc.IDStr(), "reportable columns",
				rc, tc.expReportableCols)

			ra := c.ReportableAliases()
			testhelper.DiffValsReport(t,
				tc.IDStr(), "reportable aliases",
				ra, tc.expReportableAliases)
		})
	}
}

func TestGetColInfo(t *testing.T) {
	const (
		getReportable = iota
		getSortable
		getAny
	)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		colsToAdd []ColsAddInfo
		getWhich  int
		cid       rptmaker.ColID
	}{
		{
			ID: testhelper.MkID("no columns - getAny"),
			ExpErr: testhelper.MkExpErr(
				`cannot GetColInfo: column: "any": not found`),
			getWhich: getAny,
			cid:      "any",
		},
		{
			ID: testhelper.MkID("no columns - getReportable"),
			ExpErr: testhelper.MkExpErr(
				`cannot GetReportableColInfo: column: "any": not found`),
			getWhich: getReportable,
			cid:      "any",
		},
		{
			ID: testhelper.MkID("no columns - getSortable"),
			ExpErr: testhelper.MkExpErr(
				`cannot GetSortableColInfo: column: "any": not found`),
			getWhich: getSortable,
			cid:      "any",
		},

		{
			ID: testhelper.MkID("1 column - getAny, not found"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
			},
			ExpErr: testhelper.MkExpErr(
				`cannot GetColInfo: column: "any": not found`),
			getWhich: getAny,
			cid:      "any",
		},
		{
			ID: testhelper.MkID("1 column - getReportable, not found"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.RepColName, CI: rptmaker.ReportableCI},
			},
			ExpErr: testhelper.MkExpErr(
				`cannot GetReportableColInfo: column: "any": not found`),
			getWhich: getReportable,
			cid:      "any",
		},
		{
			ID: testhelper.MkID("1 column - getReportable, missing mkCol"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.ColName, CI: rptmaker.BadSAndRCINoMkCol},
			},
			ExpErr: testhelper.MkExpErr(
				`cannot GetReportableColInfo: column: "colName":` +
					` has no "mkCol" function (it is not reportable)`),
			getWhich: getReportable,
			cid:      rptmaker.ColName,
		},
		{
			ID: testhelper.MkID("1 column - getReportable, missing colVal"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.ColName, CI: rptmaker.BadSAndRCINoColVal},
			},
			ExpErr: testhelper.MkExpErr(
				`cannot GetReportableColInfo: column: "colName":` +
					` has no "colVal" function (it is not reportable)`),
			getWhich: getReportable,
			cid:      rptmaker.ColName,
		},
		{
			ID: testhelper.MkID("1 column - getSortable, not found"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.SortColName, CI: rptmaker.SortableCI},
			},
			ExpErr: testhelper.MkExpErr(
				`cannot GetSortableColInfo: column: "any": not found`),
			getWhich: getSortable,
			cid:      "any",
		},
		{
			ID: testhelper.MkID("1 column - getSortable, missing cmpVals"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.ColName, CI: rptmaker.BadSAndRCINoCmpVals},
			},
			ExpErr: testhelper.MkExpErr(
				`cannot GetSortableColInfo: column: "colName":` +
					` has no "cmpVals" function (it is not sortable)`),
			getWhich: getSortable,
			cid:      rptmaker.ColName,
		},
		{
			ID: testhelper.MkID("1 column - getAny, OK"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.ColName, CI: rptmaker.SAndRCI},
			},
			getWhich: getAny,
			cid:      rptmaker.ColName,
		},
		{
			ID: testhelper.MkID("1 column - getReportable, OK"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.ColName, CI: rptmaker.SAndRCI},
			},
			getWhich: getReportable,
			cid:      rptmaker.ColName,
		},
		{
			ID: testhelper.MkID("1 column - getSortable, OK"),
			colsToAdd: []ColsAddInfo{
				{CID: rptmaker.ColName, CI: rptmaker.SAndRCI},
			},
			getWhich: getSortable,
			cid:      rptmaker.ColName,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			c, err := rptmaker.MakeTestCols(tc.colsToAdd)
			if err != nil {
				t.Log(tc.IDStr())
				t.Fatal("\t: unexpected error:", err)
			}

			switch tc.getWhich {
			case getAny:
				_, err = c.GetColInfo(tc.cid)
			case getReportable:
				_, err = c.GetReportableColInfo(tc.cid)
			case getSortable:
				_, err = c.GetSortableColInfo(tc.cid)
			}

			testhelper.CheckExpErr(t, err, tc)
		})
	}
}
