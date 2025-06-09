package col_test

import (
	"bytes"
	"testing"

	"github.com/nickwells/col.mod/v4/col"
	"github.com/nickwells/col.mod/v4/colfmt"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestPrintRow(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		data        []any
		extraRows   int
		hdrOpts     []col.HdrOptionFunc
		c           *col.Col
		cols        []*col.Col
		expectedVal string
	}{
		{
			ID:      testhelper.MkID("more data than columns"),
			data:    []any{int64(5), float64(1.2), "test"},
			hdrOpts: []col.HdrOptionFunc{},
			c:       col.New(&colfmt.Int{W: 3}, "an int"),
			ExpErr: testhelper.MkExpErr("Error printing row 1:" +
				" wrong number of values." +
				" Expected: 1," +
				" Received: 3"),
			expectedVal: "",
		},
		{
			ID:      testhelper.MkID("more columns than data"),
			data:    []any{},
			hdrOpts: []col.HdrOptionFunc{},
			c:       col.New(&colfmt.Int{W: 3}, "an int"),
			ExpErr: testhelper.MkExpErr("Error printing row 1:" +
				" wrong number of values." +
				" Expected: 1," +
				" Received: 0"),
			expectedVal: "",
		},
		{
			ID:      testhelper.MkID("3 columns - no header"),
			data:    []any{int64(5), float64(1.2), "test"},
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontPrint},
			c:       col.New(&colfmt.Int{W: 3}, "an int"),
			cols: []*col.Col{
				col.New(&colfmt.Float{W: 3}, "a float"),
				col.New(&colfmt.String{W: 3}, "string"),
			},
			expectedVal: `  5   1 test
`,
		},
		{
			ID:      testhelper.MkID("3 col - no underline"),
			data:    []any{int64(5), float64(1.2), "test"},
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontUnderline},
			c:       col.New(&colfmt.Int{W: 3}, "an int"),
			cols: []*col.Col{
				col.New(&colfmt.Float{W: 3}, "a float"),
				col.New(&colfmt.String{W: 3}, "string"),
			},
			expectedVal: `an int a float string
     5       1 test  
`,
		},
		{
			ID: testhelper.MkID(
				"3 col, 2 hdr lines, 1 span (narrow) - no underline"),
			data:    []any{int64(5), float64(1.2), "test"},
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontUnderline},
			c:       col.New(&colfmt.Int{W: 3}, "first line", "an int"),
			cols: []*col.Col{
				col.New(&colfmt.Float{W: 3}, "first line", "a float"),
				col.New(&colfmt.String{W: 3}, "a string"),
			},
			expectedVal: `--first line--         
an int a float a string
     5       1 test    
`,
		},
		{
			ID:      testhelper.MkID("5 col, 3 hdr lines - no underline"),
			data:    []any{"c1", "c2", "c3", "c4", "c5"},
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontUnderline},
			c:       col.New(&colfmt.String{W: 3}, "a", "b"),
			cols: []*col.Col{
				col.New(&colfmt.String{W: 3}, "a", "c"),
				col.New(&colfmt.String{W: 3}, "d"),
				col.New(&colfmt.String{W: 3}, "e", "f", "g"),
				col.New(&colfmt.String{W: 3}, "e", "f", "h"),
			},
			expectedVal: `            ---e---
---a---     ---f---
b   c   d   g   h  
c1  c2  c3  c4  c5 
`,
		},
		{
			ID:      testhelper.MkID("5 col, 3 hdr lines - default"),
			data:    []any{"c1", "c2", "c3", "c4", "c5"},
			hdrOpts: []col.HdrOptionFunc{},
			c:       col.New(&colfmt.String{W: 3}, "a", "b"),
			cols: []*col.Col{
				col.New(&colfmt.String{W: 3}, "a", "c"),
				col.New(&colfmt.String{W: 3}, "d"),
				col.New(&colfmt.String{W: 3}, "e", "f", "d"),
				col.New(&colfmt.String{W: 3}, "e", "f", "h"),
			},
			expectedVal: `            ---e---
---a---     ---f---
b   c   d   d   h  
=   =   =   =   =  
c1  c2  c3  c4  c5 
`,
		},
		{
			ID:      testhelper.MkID("5 col, 3 hdr lines - with col Sep"),
			data:    []any{"c1", "c2", "c3", "c4", "c5"},
			hdrOpts: []col.HdrOptionFunc{},
			c:       col.New(&colfmt.String{W: 3}, "a", "b").SetSep("=|= "),
			cols: []*col.Col{
				col.New(&colfmt.String{W: 3}, "a", "c"),
				col.New(&colfmt.String{W: 3}, "d"),
				col.New(&colfmt.String{W: 3}, "e", "f", "d"),
				col.New(&colfmt.String{W: 3}, "e", "f", "h"),
			},
			expectedVal: `               ---e---
----a-----     ---f---
b      c   d   d   h  
=      =   =   =   =  
c1 =|= c2  c3  c4  c5 
`,
		},
		{
			ID:        testhelper.MkID("5 col, 3 hdr lines - 2 rows"),
			data:      []any{"c1", "c2", "c3", "c4", "c5"},
			extraRows: 1,
			hdrOpts:   []col.HdrOptionFunc{},
			c:         col.New(&colfmt.String{W: 3}, "a", "b"),
			cols: []*col.Col{
				col.New(&colfmt.String{W: 3}, "a", "c"),
				col.New(&colfmt.String{W: 3}, "d"),
				col.New(&colfmt.String{W: 3}, "e", "f", "d"),
				col.New(&colfmt.String{W: 3}, "e", "f", "h"),
			},
			expectedVal: `            ---e---
---a---     ---f---
b   c   d   d   h  
=   =   =   =   =  
c1  c2  c3  c4  c5 
c1  c2  c3  c4  c5 
`,
		},
		{
			ID:        testhelper.MkID("5 col, 3 hdr lines - 3 rows, rpt hdr"),
			data:      []any{"c1", "c2", "c3", "c4", "c5"},
			extraRows: 2,
			hdrOpts: []col.HdrOptionFunc{
				col.HdrOptRepeat(2),
			},
			c: col.New(&colfmt.String{W: 3}, "a", "b"),
			cols: []*col.Col{
				col.New(&colfmt.String{W: 3}, "a", "c"),
				col.New(&colfmt.String{W: 3}, "d"),
				col.New(&colfmt.String{W: 3}, "e", "f", "d"),
				col.New(&colfmt.String{W: 3}, "e", "f", "h"),
			},
			expectedVal: `            ---e---
---a---     ---f---
b   c   d   d   h  
=   =   =   =   =  
c1  c2  c3  c4  c5 
c1  c2  c3  c4  c5 
            ---e---
---a---     ---f---
b   c   d   d   h  
=   =   =   =   =  
c1  c2  c3  c4  c5 
`,
		},
		{
			ID:      testhelper.MkID("5 col, 3 hdr lines - don't span dups"),
			data:    []any{"c1", "c2", "c3", "c4", "c5"},
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontSpanDups},
			c:       col.New(&colfmt.String{W: 3}, "a", "b"),
			cols: []*col.Col{
				col.New(&colfmt.String{W: 3}, "a", "c"),
				col.New(&colfmt.String{W: 3}, "d"),
				col.New(&colfmt.String{W: 3}, "e", "f", "d"),
				col.New(&colfmt.String{W: 3}, "e", "f", "h"),
			},
			expectedVal: `            e   e  
a   a       f   f  
b   c   d   d   h  
=   =   =   =   =  
c1  c2  c3  c4  c5 
`,
		},
	}

	for _, tc := range testCases {
		h, err := col.NewHeader(tc.hdrOpts...)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: making the Header returned an error: %s", err)

			continue
		}

		var b bytes.Buffer

		rpt, err := col.NewReport(h, &b, tc.c, tc.cols...)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: making the report returned an error: %s", err)

			continue
		}

		err = rpt.PrintRow(tc.data...)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			for i := range tc.extraRows {
				err = rpt.PrintRow(tc.data...)
				if err != nil {
					t.Log(tc.IDStr())
					t.Errorf("\t: unexpected error printing row %d: %s",
						i+2, err)

					break
				}
			}

			testhelper.DiffString(t, tc.IDStr(), "row",
				(&b).String(), tc.expectedVal)
		}
	}
}

func TestPrintRowSkipCols(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		data        []any
		skip        uint
		hdrOpts     []col.HdrOptionFunc
		c           *col.Col
		cols        []*col.Col
		expectedVal string
	}{
		{
			ID:      testhelper.MkID("more data than col"),
			data:    []any{int64(5), float64(1.2), "test"},
			skip:    1,
			hdrOpts: []col.HdrOptionFunc{},
			c:       col.New(&colfmt.Int{W: 3}, "an int"),
			ExpErr: testhelper.MkExpErr("Error printing row ",
				"too many columns to skip"),
			expectedVal: "",
		},
		{
			ID:      testhelper.MkID("5 col, no header - skip 1"),
			data:    []any{"c2", "c3", "c4", "c5"},
			skip:    1,
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontPrint},
			c:       col.New(&colfmt.String{W: 3}),
			cols: []*col.Col{
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
			},
			expectedVal: `    c2  c3  c4  c5 
`,
		},
		{
			ID:      testhelper.MkID("5 col, no header - skip 0"),
			data:    []any{"c1", "c2", "c3", "c4", "c5"},
			skip:    0,
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontPrint},
			c:       col.New(&colfmt.String{W: 3}),
			cols: []*col.Col{
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
			},
			expectedVal: `c1  c2  c3  c4  c5 
`,
		},
		{
			ID:      testhelper.MkID("5 col, no header - skip all"),
			data:    []any{},
			skip:    5,
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontPrint},
			c:       col.New(&colfmt.String{W: 3}),
			cols: []*col.Col{
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
			},
			ExpErr: testhelper.MkExpErr("Error printing row ",
				"too many columns to skip:"),
		},
		{
			ID:      testhelper.MkID("5 col, no header - skip too many"),
			data:    []any{},
			skip:    6,
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontPrint},
			c:       col.New(&colfmt.String{W: 3}),
			cols: []*col.Col{
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
			},
			ExpErr: testhelper.MkExpErr("PrintRowSkipCols(called from: ", "): ",
				"Error printing row 1:"+
					" too many columns to skip: 6 of 5"),
		},
		{
			ID:      testhelper.MkID("5 col, no header - skip too few"),
			data:    []any{},
			skip:    4,
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontPrint},
			c:       col.New(&colfmt.String{W: 3}),
			cols: []*col.Col{
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
				col.New(&colfmt.String{W: 3}),
			},
			ExpErr: testhelper.MkExpErr("PrintRowSkipCols(called from: ", "): ",
				"Error printing row 1:"+
					" wrong number of values."+
					" Skipped: 4, Expected: 1, Received: 0"),
		},
	}

	for _, tc := range testCases {
		h, err := col.NewHeader(tc.hdrOpts...)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: making the Header returned an error: %s", err)

			continue
		}

		var b bytes.Buffer

		rpt, err := col.NewReport(h, &b, tc.c, tc.cols...)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: making the Report returned an error: %s", err)

			continue
		}

		err = rpt.PrintRowSkipCols(tc.skip, tc.data...)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			testhelper.DiffString(t, tc.IDStr(), "row",
				(&b).String(), tc.expectedVal)
		}
	}
}
