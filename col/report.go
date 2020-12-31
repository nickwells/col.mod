package col

import (
	"fmt"
	"io"
	"strings"
)

//Report holds a collection of columns and header details
type Report struct {
	cols []*Col
	hdr  *Header
	w    io.Writer
}

// NewReport creates a new Report object. it will return a non-nil error if
// no columns have been given.
func NewReport(hdr *Header, w io.Writer, cols ...*Col) (*Report, error) {
	if len(cols) == 0 {
		return nil, fmt.Errorf("no columns have been given for this report")
	}

	hdr.initVals(cols)
	rpt := &Report{
		cols: cols,
		hdr:  hdr,
		w:    w,
	}

	return rpt, nil
}

// NewReportOrPanic creates a new Report object and panics if any error is
// reported
func NewReportOrPanic(hdr *Header, w io.Writer, cols ...*Col) *Report {
	rpt, err := NewReport(hdr, w, cols...)
	if err != nil {
		panic(err)
	}

	return rpt
}

// printFooter prints the footers under the numbered columns
func (rpt Report) printFooter(skip uint, vals ...interface{}) error {
	var pwe = printWithErr{w: rpt.w}

	sep := rpt.skipCols(&pwe, skip)

	for i, v := range vals {
		c := rpt.cols[i+int(skip)]

		pwe.print(sep)
		sep = c.sep

		text := ""
		if _, ok := v.(Skip); !ok {
			text = strings.Repeat(rpt.hdr.underlineCh, c.finalWidth)
		}
		pwe.print(c.stringInCol(text))
	}

	pwe.println()

	return pwe.error()
}

type printWithErr struct {
	w   io.Writer
	err error
}

// print uses fmt.Fprint to print the vals if no error has been found. It
// prints to the io.Writer of the printWithErr object.
func (pwe *printWithErr) print(vals ...interface{}) {
	if pwe.err == nil {
		_, pwe.err = fmt.Fprint(pwe.w, vals...)
	}
}

// println uses fmt.Fprintln to print the vals if no error has been found. It
// prints to the io.Writer of the printWithErr object.
func (pwe *printWithErr) println(vals ...interface{}) {
	if pwe.err == nil {
		_, pwe.err = fmt.Fprintln(pwe.w, vals...)
	}
}

// err returns any error found
func (pwe printWithErr) error() error {
	return pwe.err
}

// Skip is a type that can be passed as a column value that will print a
// blank value. It is an empty place-holder for a column.
type Skip struct{}

// PrintRow will print the values according to the specification of each
// corresponding column. It will also print the header as specified. It will
// return an error if there are not the same number of values as columns.
func (rpt *Report) PrintRow(vals ...interface{}) error {
	if len(vals) != len(rpt.cols) {
		return fmt.Errorf(
			"Error printing row %d: wrong number of values."+
				" Expected: %d,"+
				" Received: %d",
			rpt.hdr.dataRowsPrinted+1, len(rpt.cols), len(vals))
	}

	return rpt.printRowSkipping(0, vals...)
}

// PrintRowSkipCols will print the values according to the specification of
// each corresponding column. It will also print the header as specified. It
// will skip the first columns as specified. The most likely use for this is
// if you have several leading columns you want to skip. To skip individual
// columns you can use a col.Skip{}
func (rpt *Report) PrintRowSkipCols(skip uint, vals ...interface{}) error {
	if int(skip) >= len(rpt.cols) {
		return fmt.Errorf(
			"Error printing row %d: too many columns to skip: %d of %d",
			rpt.hdr.dataRowsPrinted+1, skip, len(rpt.cols))
	}

	if len(vals)+int(skip) != len(rpt.cols) {
		return fmt.Errorf(
			"Error printing row %d: wrong number of values."+
				" Skipped: %d,"+
				" Expected: %d,"+
				" Received: %d",
			rpt.hdr.dataRowsPrinted+1, skip, len(rpt.cols)-int(skip), len(vals))
	}

	return rpt.printRowSkipping(skip, vals...)
}

// printRowSkipping skips leading columns and prints the remainder. It prints
// the header as necessary and increments the number of rows printed
func (rpt *Report) printRowSkipping(skip uint, vals ...interface{}) error {
	defer rpt.hdr.incrDataRowsPrinted()

	rpt.hdr.printHeader(rpt.w, rpt.cols)

	return rpt.printValsSkipping(skip, vals...)
}

// printValsSkipping skips leading columns and prints the remainder. It does
// not print the header or increment the number of rows printed
func (rpt *Report) printValsSkipping(skip uint, vals ...interface{}) error {
	var pwe = printWithErr{w: rpt.w}

	// first collect all the strings to be printed (these may have embedded
	// new lines)

	var stringVals [][]string
	maxLines := 0
	for i, v := range vals {
		c := rpt.cols[i+int(skip)]
		str := ""
		if _, ok := v.(Skip); !ok {
			str = c.f.Formatted(v)
		}
		lines := strings.Split(str, "\n")
		if len(lines) > maxLines {
			maxLines = len(lines)
		}
		stringVals = append(stringVals, lines)
	}
	blanks := make([]string, maxLines)
	for i, lines := range stringVals {
		if len(lines) < maxLines {
			lines = append(lines, blanks[:maxLines-len(lines)]...)
			stringVals[i] = lines
		}
	}

	for j := 0; j < maxLines; j++ {
		sep := rpt.skipCols(&pwe, skip)
		for i, v := range stringVals {
			c := rpt.cols[i+int(skip)]

			pwe.print(sep)
			sep = c.sep

			pwe.print(c.stringInCol(v[j]))
		}
		pwe.println()
	}

	return pwe.error()
}

// skipCols skips leading columns
func (rpt *Report) skipCols(pwe *printWithErr, skip uint) string {
	sep := ""
	for i := uint(0); i < skip; i++ {
		c := rpt.cols[i]

		pwe.print(sep)
		sep = c.sep

		pwe.print(c.stringInCol(""))
	}
	return sep
}

// PrintFooterVals prints values for the footer. It does not print the header
// or increment the number of rows printed. It will print Header.underlineCh
// characters under the columns being printed
func (rpt Report) PrintFooterVals(skip uint, vals ...interface{}) error {
	if int(skip) >= len(rpt.cols) {
		return fmt.Errorf(
			"Error printing footer after row %d:"+
				" too many columns to skip: %d of %d",
			rpt.hdr.dataRowsPrinted, skip, len(rpt.cols))
	}

	if len(vals)+int(skip) != len(rpt.cols) {
		return fmt.Errorf(
			"Error printing footer after row %d: wrong number of values."+
				" Skipped: %d,"+
				" Expected: %d,"+
				" Received: %d",
			rpt.hdr.dataRowsPrinted, skip, len(rpt.cols)-int(skip), len(vals))
	}

	err := rpt.printFooter(skip, vals...)
	if err != nil {
		return err
	}

	return rpt.printValsSkipping(skip, vals...)
}
