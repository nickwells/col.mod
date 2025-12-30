package col

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Report holds a collection of columns and header details
type Report struct {
	cols []*Col
	hdr  *Header
	w    io.Writer
}

// NewReport creates a new Report object. If the header is nil, it is
// replaced with a newly constructed default header. If the writer is nil,
// Stdout is used.
func NewReport(hdr *Header, w io.Writer, c *Col, cs ...*Col) (*Report, error) {
	cols := []*Col{c}

	cols = append(cols, cs...)

	if err := checkColumns(cols); err != nil {
		return nil, err
	}

	if hdr == nil {
		hdr = NewHeaderOrPanic()
	}

	if w == nil {
		w = os.Stdout
	}

	hdr.initVals(cols)

	return &Report{
		cols: cols,
		hdr:  hdr,
		w:    w,
	}, nil
}

// NewReportOrPanic returns a new Report object. If an error was returned
// when the Report was created then this will panic.
func NewReportOrPanic(hdr *Header, w io.Writer, c *Col, cs ...*Col) *Report {
	r, err := NewReport(hdr, w, c, cs...)
	if err != nil {
		panic(err)
	}

	return r
}

// StdRpt creates a new Report object. The header used is a default value and
// the writer used is the os.Stdout; use the NewReport function if you want
// to use other values. Any errors will cause this to panic.
func StdRpt(c *Col, cs ...*Col) *Report {
	return NewReportOrPanic(NewHeaderOrPanic(), os.Stdout, c, cs...)
}

// checkColumns runs the Formatter's Check function for each Col and returns
// an error if any of them returns a non-nil error.
func checkColumns(cols []*Col) error {
	var allErrs []error

	for i, c := range cols {
		if err := c.f.Check(); err != nil {
			allErrs = append(allErrs,
				fmt.Errorf("column[%d] (%q) has a bad Formatter: %w",
					i, c.headers, err))
		}
	}

	return errors.Join(allErrs...)
}

// printFooter prints the footers under the numbered columns
func (rpt Report) printFooter(skip int, vals ...any) error {
	pwe := printWithErr{w: rpt.w}

	sep := rpt.skipCols(&pwe, skip)

	for i, v := range vals {
		c := rpt.cols[i+skip]

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
func (pwe *printWithErr) print(vals ...any) {
	if pwe.err == nil {
		_, pwe.err = fmt.Fprint(pwe.w, vals...)
	}
}

// println uses fmt.Fprintln to print the vals if no error has been found. It
// prints to the io.Writer of the printWithErr object.
func (pwe *printWithErr) println(vals ...any) {
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
func (rpt *Report) PrintRow(vals ...any) error {
	if len(vals) != len(rpt.cols) {
		return fmt.Errorf(
			"PrintRow(called from: %s): printing row %d:"+
				" wrong number of values."+
				" Expected: %d,"+
				" Received: %d",
			caller(), rpt.hdr.dataRowsPrinted+1,
			len(rpt.cols), len(vals))
	}

	return rpt.printRowSkipping(0, vals...)
}

// PrintRowSkipCols will print the values according to the specification of
// each corresponding column. It will also print the header as specified. It
// will skip the first columns as specified. The most likely use for this is
// if you have several leading columns you want to skip. To skip individual
// columns you can use a col.Skip{}
func (rpt *Report) PrintRowSkipCols(skip int, vals ...any) error {
	if err := rpt.checkSkipVal(skip, len(vals)); err != nil {
		return fmt.Errorf(
			"PrintRowSkipCols(called from: %s): printing row %d: %s",
			caller(), rpt.hdr.dataRowsPrinted+1, err)
	}

	return rpt.printRowSkipping(skip, vals...)
}

// printRowSkipping skips leading columns and prints the remainder. It prints
// the header as necessary and increments the number of rows printed
func (rpt *Report) printRowSkipping(skip int, vals ...any) error {
	defer rpt.hdr.incrDataRowsPrinted()

	rpt.hdr.printHeader(rpt.w, rpt.cols)

	return rpt.printValsSkipping(skip, vals...)
}

// printValsSkipping skips leading columns and prints the remainder. It does
// not print the header or increment the number of rows printed
func (rpt *Report) printValsSkipping(skip int, vals ...any) error {
	pwe := printWithErr{w: rpt.w}

	// generate all the lines to be printed for this row (note that some
	// columns can be formatted into multiple lines of text)
	lineVals, maxLines := rpt.splitVals(skip, vals...)
	lineVals = addBlanks(lineVals, maxLines)

	for j := range maxLines {
		sep := rpt.skipCols(&pwe, skip)

		for i, v := range lineVals {
			c := rpt.cols[i+skip]

			pwe.print(sep)
			pwe.print(c.stringInCol(v[j]))

			sep = c.sep
		}

		pwe.println()
	}

	return pwe.error()
}

// skipCols skips leading columns
func (rpt *Report) skipCols(pwe *printWithErr, skip int) string {
	sep := ""

	for i := range skip {
		c := rpt.cols[i]

		pwe.print(sep)
		pwe.print(c.stringInCol(""))

		sep = c.sep
	}

	return sep
}

// PrintFooterVals prints values for the footer. It does not print the header
// or increment the number of rows printed. It will print Header.underlineCh
// characters under the columns being printed
func (rpt Report) PrintFooterVals(skip int, vals ...any) error {
	if err := rpt.checkSkipVal(skip, len(vals)); err != nil {
		return fmt.Errorf(
			"PrintFooterVals(called from: %s):"+
				" printing footer after row %d: %s",
			caller(), rpt.hdr.dataRowsPrinted, err)
	}

	err := rpt.printFooter(skip, vals...)
	if err != nil {
		return err
	}

	return rpt.printValsSkipping(skip, vals...)
}

// checkSkipVal returns an error if the skip value is invalid. This can mean
// it is less than zero, greater than the number of columns or that the
// number of values plus the number to skip does not equal the number of
// columns expected. Note that the skip parameter is strictly uneccessary as
// the Print... func's could all take the number of values to skip as
// implicit in the number of values to print but this API makes it clearer
// what the intention of the developer is and, hopefully, will make any
// errors easier to find.
func (rpt Report) checkSkipVal(skip, valCnt int) error {
	colCnt := len(rpt.cols)

	if skip >= colCnt {
		return fmt.Errorf("too many columns to skip: %d of %d", skip, colCnt)
	}

	if skip < 0 {
		return fmt.Errorf("the number of columns to skip (%d) must be > 0",
			skip)
	}

	if valCnt+skip != colCnt {
		return fmt.Errorf("wrong number of values."+
			" Skipped: %d,"+
			" Expected: %d,"+
			" Received: %d",
			skip,
			colCnt-skip,
			valCnt)
	}

	return nil
}

// splitVals takes each value, formats it, splits the formatted value around
// newlines and adds each generated slice of values into the slice of slices
// to be returned. It simultaneously keeps track of the maximum number of
// lines detected. Finally it returns the collection of lines and the maximum
// number of lines encountered.
func (rpt *Report) splitVals(skip int, vals ...any) ([][]string, int) {
	var lineVals [][]string

	maxLines := 0

	for i, v := range vals {
		c := rpt.cols[i+skip]
		str := ""

		if _, ok := v.(Skip); !ok {
			str = c.f.Formatted(v)
		}

		lines := strings.Split(str, "\n")

		maxLines = max(len(lines), maxLines)

		lineVals = append(lineVals, lines)
	}

	return lineVals, maxLines
}

// addBlanks takes a slice of slices (each of which may have different
// numbers of members) and ensures that each has the same number of entries
// by adding blank strings at the end.
func addBlanks(lineVals [][]string, maxLines int) [][]string {
	blanks := make([]string, maxLines)

	for i, lines := range lineVals {
		if len(lines) < maxLines {
			lines = append(lines, blanks[:maxLines-len(lines)]...)
			lineVals[i] = lines
		}
	}

	return lineVals
}
