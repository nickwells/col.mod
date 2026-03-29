package rptmaker

import (
	"strings"

	"github.com/nickwells/col.mod/v6/col"
	"github.com/nickwells/english.mod/english"
)

// ColCmpFunc is the type of the column comparison function used to sort the
// data being reported on prior to printing. It should return an integer
// indicating whether the first value is less than, equal to or greater than
// the second value.
type ColCmpFunc[T any] func(a, b T) int

// ColMkFunc is the type of the column construction function used to create
// the col.Col. It should take a params argument giving details of how to
// construct the columns and a slice to be used as the column headings.
type ColMkFunc[P any] func(params P, headings []string) *col.Col

// ColValFunc is the type of a function used to extract a column specific
// value from a record of the data being reported on.
type ColValFunc[T any] func(r T) any

// ColInfo holds the information needed by the report cols for each
// column.
type ColInfo[P, T any] struct {
	desc     string
	headings []string
	mkCol    ColMkFunc[P]
	colVal   ColValFunc[T]
	cmpVals  ColCmpFunc[T]
}

// Headings returns the column headings.
func (ci ColInfo[P, T]) Headings() []string {
	return ci.headings
}

// headingDesc returns a description of the column headings.
func (ci ColInfo[P, T]) headingDesc() string {
	if len(ci.headings) == 0 {
		return "This column is unheaded."
	}

	return "This column is headed: " +
		english.JoinQuoted(ci.headings, "/", "/")
}

// Desc returns the column description.
func (ci ColInfo[P, T]) Desc() string {
	return ci.desc
}

// FullDesc returns the column description and a note giving the associated
// column headings.
func (ci ColInfo[P, T]) FullDesc() string {
	parts := []string{}

	if ci.desc != "" {
		parts = append(parts, ci.desc)
	}

	parts = append(parts, ci.headingDesc())

	return strings.Join(parts, " ")
}

// IsSortable returns true if the column can be used for sorting.
func (ci ColInfo[P, T]) IsSortable() bool {
	return ci.cmpVals != nil
}

// IsReportable returns true if the column can be used when printing the
// report.
func (ci ColInfo[P, T]) IsReportable() bool {
	return ci.mkCol != nil && ci.colVal != nil
}

// NewColInfo creates a [ColInfo]. Note that the mkCol, colVal and cmpVals
// functions can all be nil. Columns with a nil mkCol function or a nil
// colVal function cannot be used as a reported column but can still be used
// for sorting. Similarly columns with a nil cmpVals function cannot be
// sorted on. Columns that are neither sortable nor reportable will not be
// used at all.
func NewColInfo[P, T any](
	desc string,
	headings []string,
	mkCol ColMkFunc[P],
	colVal ColValFunc[T],
	cmpVals ColCmpFunc[T],
) *ColInfo[P, T] {
	return &ColInfo[P, T]{
		desc:     desc,
		headings: headings,
		mkCol:    mkCol,
		colVal:   colVal,
		cmpVals:  cmpVals,
	}
}
