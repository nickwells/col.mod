package col_test

import (
	"bytes"
	"testing"

	"github.com/nickwells/col.mod/v3/col"
	"github.com/nickwells/col.mod/v3/col/colfmt"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestPrintRow(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		data        []interface{}
		hdrOpts     []col.HdrOptionFunc
		c           *col.Col
		cols        []*col.Col
		expectedVal string
	}{
		{
			ID:      testhelper.MkID("more data than columns"),
			data:    []interface{}{int64(5), float64(1.2), "test"},
			hdrOpts: []col.HdrOptionFunc{},
			c:       col.New(colfmt.Int{W: 3}, "an int"),
			ExpErr: testhelper.MkExpErr("Error printing row ",
				"wrong number of values",
				" Expected: ",
				" Received"),
			expectedVal: "",
		},
		{
			ID:      testhelper.MkID("3 columns - no header"),
			data:    []interface{}{int64(5), float64(1.2), "test"},
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontPrint},
			c:       col.New(&colfmt.Int{W: 3}, "an int"),
			cols: []*col.Col{
				col.New(&colfmt.Float{W: 3}, "a float"),
				col.New(colfmt.String{W: 3}, "string"),
			},
			expectedVal: `  5   1 test
`,
		},
		{
			ID:      testhelper.MkID("3 col - no underline"),
			data:    []interface{}{int64(5), float64(1.2), "test"},
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
				"3 col, 2 header lines, 1 span (narrow) - no underline"),
			data:    []interface{}{int64(5), float64(1.2), "test"},
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontUnderline},
			c: col.New(
				&colfmt.Int{W: 3},
				"first line",
				"an int"),

			cols: []*col.Col{
				col.New(
					&colfmt.Float{W: 3},
					"first line",
					"a float"),
				col.New(
					colfmt.String{W: 3},
					"a string"),
			},
			expectedVal: `--first line--         
an int a float a string
     5       1 test    
`,
		},
		{
			ID:      testhelper.MkID("5 col, 3 header lines - no underline"),
			data:    []interface{}{"c1", "c2", "c3", "c4", "c5"},
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontUnderline},
			c:       col.New(colfmt.String{W: 3}, "a", "b"),
			cols: []*col.Col{
				col.New(colfmt.String{W: 3}, "a", "c"),
				col.New(colfmt.String{W: 3}, "d"),
				col.New(colfmt.String{W: 3}, "e", "f", "g"),
				col.New(colfmt.String{W: 3}, "e", "f", "h"),
			},
			expectedVal: `            ---e---
---a---     ---f---
b   c   d   g   h  
c1  c2  c3  c4  c5 
`,
		},
		{
			ID:      testhelper.MkID("5 col, 3 header lines - default"),
			data:    []interface{}{"c1", "c2", "c3", "c4", "c5"},
			hdrOpts: []col.HdrOptionFunc{},
			c:       col.New(colfmt.String{W: 3}, "a", "b"),
			cols: []*col.Col{
				col.New(colfmt.String{W: 3}, "a", "c"),
				col.New(colfmt.String{W: 3}, "d"),
				col.New(colfmt.String{W: 3}, "e", "f", "d"),
				col.New(colfmt.String{W: 3}, "e", "f", "h"),
			},
			expectedVal: `            ---e---
---a---     ---f---
b   c   d   d   h  
=   =   =   =   =  
c1  c2  c3  c4  c5 
`,
		},
		{
			ID:      testhelper.MkID("5 col, 3 header lines - don't span dups"),
			data:    []interface{}{"c1", "c2", "c3", "c4", "c5"},
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontSpanDups},
			c:       col.New(colfmt.String{W: 3}, "a", "b"),
			cols: []*col.Col{
				col.New(colfmt.String{W: 3}, "a", "c"),
				col.New(colfmt.String{W: 3}, "d"),
				col.New(colfmt.String{W: 3}, "e", "f", "d"),
				col.New(colfmt.String{W: 3}, "e", "f", "h"),
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
		rpt := col.NewReport(h, &b, tc.c, tc.cols...)
		err = rpt.PrintRow(tc.data...)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			testhelper.DiffString(t, tc.IDStr(), "row",
				(&b).String(), tc.expectedVal)
		}
	}
}

func TestPrintRowSkipCols(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		data        []interface{}
		skip        uint
		hdrOpts     []col.HdrOptionFunc
		c           *col.Col
		cols        []*col.Col
		expectedVal string
	}{
		{
			ID:      testhelper.MkID("more data than col"),
			data:    []interface{}{int64(5), float64(1.2), "test"},
			skip:    1,
			hdrOpts: []col.HdrOptionFunc{},
			c:       col.New(colfmt.Int{W: 3}, "an int"),
			ExpErr: testhelper.MkExpErr("Error printing row ",
				"too many columns to skip"),
			expectedVal: "",
		},
		{
			ID:      testhelper.MkID("5 col, no header - skip 1"),
			data:    []interface{}{"c2", "c3", "c4", "c5"},
			skip:    1,
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontPrint},
			c:       col.New(colfmt.String{W: 3}),
			cols: []*col.Col{
				col.New(colfmt.String{W: 3}),
				col.New(colfmt.String{W: 3}),
				col.New(colfmt.String{W: 3}),
				col.New(colfmt.String{W: 3}),
			},
			expectedVal: `    c2  c3  c4  c5 
`,
		},
		{
			ID:      testhelper.MkID("5 col, no header - skip 0"),
			data:    []interface{}{"c1", "c2", "c3", "c4", "c5"},
			skip:    0,
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontPrint},
			c:       col.New(colfmt.String{W: 3}),
			cols: []*col.Col{
				col.New(colfmt.String{W: 3}),
				col.New(colfmt.String{W: 3}),
				col.New(colfmt.String{W: 3}),
				col.New(colfmt.String{W: 3}),
			},
			expectedVal: `c1  c2  c3  c4  c5 
`,
		},
		{
			ID:      testhelper.MkID("5 col, no header - skip all"),
			data:    []interface{}{},
			skip:    5,
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontPrint},
			c:       col.New(colfmt.String{W: 3}),
			cols: []*col.Col{
				col.New(colfmt.String{W: 3}),
				col.New(colfmt.String{W: 3}),
				col.New(colfmt.String{W: 3}),
				col.New(colfmt.String{W: 3}),
			},
			ExpErr: testhelper.MkExpErr("Error printing row ",
				"too many columns to skip:"),
		},
		{
			ID:      testhelper.MkID("5 col, no header - skip too many"),
			data:    []interface{}{},
			skip:    6,
			hdrOpts: []col.HdrOptionFunc{col.HdrOptDontPrint},
			c:       col.New(colfmt.String{W: 3}),
			cols: []*col.Col{
				col.New(colfmt.String{W: 3}),
				col.New(colfmt.String{W: 3}),
				col.New(colfmt.String{W: 3}),
				col.New(colfmt.String{W: 3}),
			},
			ExpErr: testhelper.MkExpErr("Error printing row ",
				"too many columns to skip:"),
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
		rpt := col.NewReport(h, &b, tc.c, tc.cols...)
		err = rpt.PrintRowSkipCols(tc.skip, tc.data...)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			testhelper.DiffString(t, tc.IDStr(), "row",
				(&b).String(), tc.expectedVal)
		}
	}
}
