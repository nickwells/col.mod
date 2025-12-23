package col_test

import (
	"fmt"
	"io"
	"os"

	"github.com/nickwells/col.mod/v6/col"
	"github.com/nickwells/col.mod/v6/colfmt"
)

// Example_report demonstrates how the col package might be used to generate
// a report. Notice how we have multi-line column headings and how common
// header text is shared between adjacent columns (this makes for more
// concise reports). This uses the StdRpt func to create the report; this is
// suitable when no special header handling is needed and the report should
// be written to os.Stdout. Note that StdRpt will panic if any columns are
// incorrectly created.
func Example_report() {
	// Now create the report. List the columns we want giving the format
	// types and the column headings
	rpt := col.StdRpt(col.New(&colfmt.Int{}, "Date"),
		col.New(&colfmt.Int{W: 2}, "Number of", "Boys"),
		col.New(&colfmt.Int{W: 2}, "Number of", "Girls"),
		col.New(&colfmt.Int{W: 3}, "Class", "Size"),
		col.New(&colfmt.Float{Prec: 2}, "Ratio", "Boys-Girls"),
	)

	type rowStruct struct {
		y     int
		boys  int
		girls int
	}

	for _, v := range []rowStruct{
		{y: 2011, boys: 14, girls: 13},
		{y: 2012, boys: 12, girls: 16},
		{y: 2013, boys: 13, girls: 13},
	} {
		err := rpt.PrintRow(v.y,
			v.boys, v.girls,
			v.boys+v.girls,
			float64(v.boys)/float64(v.girls))
		if err != nil {
			fmt.Println("Unexpected error found while printing a row:", err)
			break
		}
	}
	// Output:
	//      -Number of- Class      Ratio
	// Date  Boys Girls  Size Boys-Girls
	// ====  ==== =====  ==== ==========
	// 2011    14    13    27       1.08
	// 2012    12    16    28       0.75
	// 2013    13    13    26       1.00
}

// Example_report2 demonstrates how the col package might be used to generate
// a report. This is a more sophisticated report demonstrating how you can
// customise the header, skip columns and print column totals.
func Example_report2() {
	// First, create the header
	h, err := col.NewHeader(
		// On the first page only, print a report description
		col.HdrOptPreHdrFunc(func(w io.Writer, n int64) {
			if n == 0 {
				fmt.Fprintln(w,
					"A report on the variation in class sizes over time")
			}
		}),
		// Use '-' to underline the column headings
		col.HdrOptUnderlineWith('-'),
	)
	if err != nil {
		fmt.Println("Unexpected error found while building the Header:", err)
		return
	}

	// Now create the report. List the columns we want giving the format
	// types and the column headings
	rpt, err := col.NewReport(h, os.Stdout,
		col.New(&colfmt.Int{}, "Academic", "Year"),
		col.New(&colfmt.Int{}, "Date"),
		col.New(&colfmt.Int{W: 2}, "Number of", "Boys"),
		col.New(&colfmt.Int{W: 2}, "Number of", "Girls"),
		col.New(&colfmt.Float{Prec: 2}, "Ratio", "Boys-Girls"),
		col.New(&colfmt.Int{W: 3}, "Class", "Size"),
	)
	if err != nil {
		fmt.Println("Unexpected error found while building the Report:", err)
		return
	}

	type rowStruct struct {
		year  int
		date  int
		boys  int
		girls int
	}

	lastYear := 0
	totBoys := 0
	totGirls := 0
	count := 0

	for _, v := range []rowStruct{
		{year: 4, date: 2011, boys: 14, girls: 13},
		{year: 4, date: 2012, boys: 12, girls: 16},
		{year: 4, date: 2013, boys: 13, girls: 13},
		{year: 5, date: 2011, boys: 14, girls: 13},
		{year: 5, date: 2012, boys: 12, girls: 16},
		{year: 5, date: 2013, boys: 13, girls: 13},
		{year: 6, date: 2011, boys: 13, girls: 13},
	} {
		count++
		totBoys += v.boys
		totGirls += v.girls
		tot := v.boys + v.girls
		ratio := float64(v.boys) / float64(v.girls)

		var ratioVal any

		if ratio >= 1.005 || ratio <= 0.995 {
			ratioVal = any(ratio)
		} else {
			ratioVal = any(col.Skip{})
		}

		if v.year == lastYear {
			// This illustrates the use of the PrintRowSkipCols func. Note
			// that this could equally have been done by passing col.Skip{}
			// as the first value
			err = rpt.PrintRowSkipCols(1,
				v.date, v.boys, v.girls, ratioVal, tot)
		} else {
			err = rpt.PrintRow(v.year, v.date, v.boys, v.girls, ratioVal, tot)
		}

		if err != nil {
			fmt.Println("Unexpected error found while printing a row:",
				err)

			break
		}

		lastYear = v.year
	}

	// now print the column totals using PrintFooterVals
	ratio := float64(totBoys) / float64(totGirls)
	avgClassSize := (totBoys + totGirls) / count

	err = rpt.PrintFooterVals(2, totBoys, totGirls, ratio, avgClassSize)
	if err != nil {
		fmt.Println("Unexpected error found while printing the report footer:",
			err)
	}
	// Output:
	// A report on the variation in class sizes over time
	// Academic      -Number of-      Ratio Class
	//     Year Date  Boys Girls Boys-Girls  Size
	//     ---- ----  ---- ----- ----------  ----
	//        4 2011    14    13       1.08    27
	//          2012    12    16       0.75    28
	//          2013    13    13               26
	//        5 2011    14    13       1.08    27
	//          2012    12    16       0.75    28
	//          2013    13    13               26
	//        6 2011    13    13               26
	//               ----- ----- ---------- -----
	//                  91    97       0.94    26
}
