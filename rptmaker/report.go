package rptmaker

import (
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/nickwells/col.mod/v6/col"
)

// Report holds the details needed to generate a report
type Report[P, T any] struct {
	rpt    *col.Report
	colIDs []ColID
	cols   Cols[P, T]
}

// MakeReport creates a report
func (c Cols[P, T]) MakeReport(
	p P,
	w io.Writer,
	colIDs []ColID,
	hOpts ...col.HdrOptionFunc,
) (
	*Report[P, T], error,
) {
	const errIntro = "cannot create the Report:"

	if len(colIDs) == 0 {
		return nil, errors.New(errIntro + " no columns were given")
	}

	h, err := col.NewHeader(hOpts...)
	if err != nil {
		return nil,
			fmt.Errorf("%s cannot create the report header: %w", errIntro, err)
	}

	cols := make([]*col.Col, 0, len(colIDs))
	for _, cid := range colIDs {
		ci, err := c.GetReportableColInfo(cid)
		if err != nil {
			return nil, fmt.Errorf("%s %w", errIntro, err)
		}

		col := ci.mkCol(p, ci.Headings())
		if col == nil {
			return nil,
				fmt.Errorf("%s the mkCol function for %q returned a nil column",
					errIntro, cid)
		}

		cols = append(cols, col)
	}

	var rpt *col.Report

	if len(cols) == 1 {
		rpt, err = col.NewReport(h, w, cols[0])
	} else {
		rpt, err = col.NewReport(h, w, cols[0], cols[1:]...)
	}

	if err != nil {
		return nil, fmt.Errorf("%s %w", errIntro, err)
	}

	return &Report[P, T]{
		rpt:    rpt,
		colIDs: colIDs,
		cols:   c,
	}, nil
}

// MkCmpFunc returns a comparison function suitable to pass to
// slices.SortFunc. It is composed from the individual per-column comparison
// functions according to the columns given in the slice of [SortColumn]
// values. This also includes details on whether a column should be sorted in
// the reverse order in which case a per-column comparison function is
// generated with the parameter order swapped.
func (r Report[P, T]) MkCmpFunc(
	sortCols []SortColumn,
) (
	func(a, b T) int, error,
) {
	var cmpFuncs []ColCmpFunc[T]

	const errIntro = "cannot make the comparison function:"

	for _, sc := range sortCols {
		ci, ok := r.cols.colMap[sc.ID]
		if !ok {
			return nil, fmt.Errorf("%s %w", errIntro, MkColNotFoundErr(sc.ID))
		}

		if ci.cmpVals == nil {
			return nil,
				fmt.Errorf("%s %w", errIntro, MkNoFuncErr(sc.ID, cmpValsFName))
		}

		cf := ci.cmpVals
		if sc.Backwards {
			cf = func(a, b T) int {
				return ci.cmpVals(b, a)
			}
		}

		cmpFuncs = append(cmpFuncs, cf)
	}

	return func(a, b T) int {
		var rval int

		for _, cf := range cmpFuncs {
			rval = cf(a, b)

			if rval != 0 {
				return rval
			}
		}

		return rval
	}, nil
}

// PrintLine gathers the values to be printed from the v supplied using the
// Report's value functions. It returns a non-nil error if any of the columns
// is not found in the Report's [Cols], if the [ColInfo] has no value
// function or if the row printing fails.
func (r Report[P, T]) PrintLine(v T) error {
	vals := make([]any, 0, len(r.colIDs))

	for _, cid := range r.colIDs {
		ci, ok := r.cols.colMap[cid]
		if !ok {
			return MkColNotFoundErr(cid)
		}

		valF := ci.colVal
		if valF == nil {
			return MkNoFuncErr(cid, colValFName)
		}

		vals = append(vals, valF(v))
	}

	return r.rpt.PrintRow(vals...)
}

// Print takes the slice of values, sorts them according to the supplied
// sortCols and then prints them line by line.
func (r Report[P, T]) Print(vals []T, sortCols []SortColumn) error {
	if len(sortCols) > 0 {
		cf, err := r.MkCmpFunc(sortCols)
		if err != nil {
			return err
		}

		slices.SortFunc(vals, cf)
	}

	for _, v := range vals {
		if err := r.PrintLine(v); err != nil {
			return err
		}
	}

	return nil
}
