package rptmaker

import (
	"fmt"
	"maps"

	"github.com/nickwells/strdist.mod/v2/strdist"
)

// ColID is the type of a value indexing the column information
type ColID string

// Cols holds the collection of columns for the report and aliases for
// column names. Note that the zero value is not usable - it should be
// constructed with the [NewCols] function.
//
// The P type is the type of the parameter object to be used when making a
// column and the T type is the type of the data to be reported on.
type Cols[P, T any] struct {
	colMap            map[ColID]*ColInfo[P, T]
	aliases           map[ColID][]ColID
	reportableAliases map[ColID][]ColID
	sortableAliases   map[ColID][]ColID
}

// GetReportableColInfo returns the column information associated with the
// given ColID and a nil error. If there is no column data associated with
// the given ColID or it lacks the functions needed to report the column, a
// nil ColInfo pointer is returned and a non-nil error.
func (c Cols[P, T]) GetReportableColInfo(cid ColID) (*ColInfo[P, T], error) {
	const errIntro = "cannot GetReportableColInfo:"

	ci, ok := c.colMap[cid]
	if !ok {
		return nil, fmt.Errorf("%s %w", errIntro, MkColNotFoundErr(cid))
	}

	if ci.mkCol == nil {
		return nil, fmt.Errorf("%s %w", errIntro, MkNoFuncErr(cid, mkColFName))
	}

	if ci.colVal == nil {
		return nil, fmt.Errorf("%s %w", errIntro, MkNoFuncErr(cid, colValFName))
	}

	return ci, nil
}

// GetSortableColInfo returns the column information associated with the
// given ColID and a nil error. If there is no column data associated with
// the given ColID or it lacks the comparison function needed to sort the
// column, a nil ColInfo pointer is returned and a non-nil error.
func (c Cols[P, T]) GetSortableColInfo(cid ColID) (*ColInfo[P, T], error) {
	const errIntro = "cannot GetSortableColInfo:"

	ci, ok := c.colMap[cid]
	if !ok {
		return nil, fmt.Errorf("%s %w", errIntro, MkColNotFoundErr(cid))
	}

	if ci.cmpVals == nil {
		return nil,
			fmt.Errorf("%s %w", errIntro, MkNoFuncErr(cid, cmpValsFName))
	}

	return ci, nil
}

// GetColInfo returns the column information associated with the given ColID
// and a nil error. If there is no column data associated with the given
// ColID a nil ColInfo pointer is returned and a non-nil error.
func (c Cols[P, T]) GetColInfo(cid ColID) (*ColInfo[P, T], error) {
	const errIntro = "cannot GetColInfo:"

	ci, ok := c.colMap[cid]
	if !ok {
		return nil, fmt.Errorf("%s %w", errIntro, MkColNotFoundErr(cid))
	}

	return ci, nil
}

// NewCols returns a properly constructed Cols ready to use.
func NewCols[P, T any]() *Cols[P, T] {
	return &Cols[P, T]{
		colMap:            map[ColID]*ColInfo[P, T]{},
		aliases:           map[ColID][]ColID{},
		reportableAliases: map[ColID][]ColID{},
		sortableAliases:   map[ColID][]ColID{},
	}
}

// Add adds the column to the collection. It returns an error if the
// column is already present or if the column information is invalid.
func (c *Cols[P, T]) Add(cid ColID, ci *ColInfo[P, T]) error {
	const errIntro = "cannot Add:"

	if _, exists := c.colMap[cid]; exists {
		return fmt.Errorf("%s %w", errIntro, MkDupColErr(cid))
	}

	if _, exists := c.aliases[cid]; exists {
		return fmt.Errorf("%s %w", errIntro, MkAliasNameErr(cid))
	}

	if ci == nil {
		return fmt.Errorf("%s %w", errIntro, MkNoColInfoErr(cid))
	}

	if !ci.IsReportable() && !ci.IsSortable() {
		return MkUnusableErr(cid)
	}

	c.colMap[cid] = ci

	return nil
}

// colNames returns a list of column names matching the criteria checked
// by the supplied checkFunc.
func (c Cols[P, T]) colNames(checkFunc func(ci ColInfo[P, T]) bool) []string {
	colNames := []string{}

	for name, ci := range c.colMap {
		if checkFunc(*ci) {
			colNames = append(colNames, string(name))
		}
	}

	return colNames
}

// AddAlias adds an alias for a column name or names. The colName and any
// extras must already be known to the cols so it is advised that all the
// columns are added using [Cols.Add] before the aliases (if any) are
// added.
//
// Note that the alias name must not already exist in either the common set
// of aliases, the set of reportable aliases or the set of sortable
// aliases. Aliases added using this function will appear in the results of
// both the [ReportableAliases] and [SortableAliases] functions.
func (c *Cols[P, T]) AddAlias(
	alias ColID, colName ColID, extras ...ColID,
) error {
	const errIntro = "cannot AddAlias:"

	if _, exists := c.colMap[alias]; exists {
		return fmt.Errorf("%s %w", errIntro, MkAliasIsColErr(alias))
	}

	if _, exists := c.aliases[alias]; exists {
		return fmt.Errorf("%s %w", errIntro, MkDupAliasErr(alias, "common"))
	}

	if _, exists := c.reportableAliases[alias]; exists {
		return fmt.Errorf("%s %w", errIntro, MkDupAliasErr(alias, "reportable"))
	}

	if _, exists := c.sortableAliases[alias]; exists {
		return fmt.Errorf("%s %w", errIntro, MkDupAliasErr(alias, "sortable"))
	}

	cols := []ColID{colName}
	cols = append(cols, extras...)

	for _, cid := range cols {
		if _, exists := c.colMap[cid]; !exists {
			colNames := c.colNames(func(_ ColInfo[P, T]) bool { return true })

			problem := fmt.Sprintf("there is no column called %q%s",
				cid,
				strdist.SuggestionString(
					strdist.SuggestedVals(string(cid), colNames)))

			return fmt.Errorf("%s %w", errIntro, MkAliasErr(alias, problem))
		}
	}

	c.aliases[alias] = cols

	return nil
}

// AddReportableAlias adds an alias for a column name or names. The colName
// and any extras must already be known to the cols so it is advised that all
// the columns are added using [Cols.Add] before the aliases (if any) are
// added. Also each column that the alias maps to must be reportable.
//
// Note that the alias name must not already exist in either the common set
// of aliases or the set of reportable aliases. There may, though, be a
// sortable alias with the same name. Aliases added using this function will
// only appear in the results of the [ReportableAliases] function.
func (c *Cols[P, T]) AddReportableAlias(
	alias ColID, colName ColID, extras ...ColID,
) error {
	const errIntro = "cannot AddReportableAlias:"

	if _, exists := c.colMap[alias]; exists {
		return fmt.Errorf("%s %w", errIntro, MkAliasIsColErr(alias))
	}

	if _, exists := c.aliases[alias]; exists {
		return fmt.Errorf("%s %w", errIntro, MkDupAliasErr(alias, "common"))
	}

	if _, exists := c.reportableAliases[alias]; exists {
		return fmt.Errorf("%s %w", errIntro, MkDupAliasErr(alias, "reportable"))
	}

	cols := []ColID{colName}
	cols = append(cols, extras...)

	for _, cid := range cols {
		ci, exists := c.colMap[cid]
		if !exists {
			colNames := c.colNames(ColInfo[P, T].IsReportable)

			problem := fmt.Sprintf("there is no column called %q%s",
				cid,
				strdist.SuggestionString(
					strdist.SuggestedVals(string(cid), colNames)))

			return fmt.Errorf("%s %w", errIntro, MkAliasErr(alias, problem))
		}

		if !ci.IsReportable() {
			problem := fmt.Sprintf("column %q is not reportable",
				cid)

			return fmt.Errorf("%s %w", errIntro, MkAliasErr(alias, problem))
		}
	}

	c.reportableAliases[alias] = cols

	return nil
}

// AddSortableAlias adds an alias for a column name or names. The colName
// and any extras must already be known to the cols so it is advised that all
// the columns are added using [Cols.Add] before the aliases (if any) are
// added. Also each column that the alias maps to must be sortable.
//
// Note that the alias name must not already exist in either the common set
// of aliases or the set of sortable aliases. There may, though, be a
// reportable alias with the same name. Aliases added using this function will
// only appear in the results of the [SortableAliases] function.
func (c *Cols[P, T]) AddSortableAlias(
	alias ColID, colName ColID, extras ...ColID,
) error {
	const errIntro = "cannot AddSortableAlias:"

	if _, exists := c.colMap[alias]; exists {
		return fmt.Errorf("%s %w", errIntro, MkAliasIsColErr(alias))
	}

	if _, exists := c.aliases[alias]; exists {
		return fmt.Errorf("%s %w", errIntro, MkDupAliasErr(alias, "common"))
	}

	if _, exists := c.sortableAliases[alias]; exists {
		return fmt.Errorf("%s %w", errIntro, MkDupAliasErr(alias, "sortable"))
	}

	cols := []ColID{colName}
	cols = append(cols, extras...)

	for _, cid := range cols {
		ci, exists := c.colMap[cid]
		if !exists {
			colNames := c.colNames(ColInfo[P, T].IsSortable)

			problem := fmt.Sprintf("there is no column called %q%s",
				cid,
				strdist.SuggestionString(
					strdist.SuggestedVals(string(cid), colNames)))

			return fmt.Errorf("%s %w", errIntro, MkAliasErr(alias, problem))
		}

		if !ci.IsSortable() {
			problem := fmt.Sprintf("column %q is not sortable",
				cid)

			return fmt.Errorf("%s %w", errIntro, MkAliasErr(alias, problem))
		}
	}

	c.sortableAliases[alias] = cols

	return nil
}

// Reportable returns a map of ColIDs to their associated descriptions for
// the reportable columns only.
func (c Cols[P, T]) Reportable() map[ColID]string {
	rptCols := map[ColID]string{}

	for cid, ci := range c.colMap {
		if ci.IsReportable() {
			rptCols[cid] = ci.FullDesc()
		}
	}

	return rptCols
}

// ReportableAliases returns a map of aliases to the associated reportable
// columns.
func (c Cols[P, T]) ReportableAliases() map[ColID][]ColID {
	aliases := map[ColID][]ColID{}
	maps.Copy(aliases, c.reportableAliases)

	for alias, cids := range c.aliases {
		cols := []ColID{}

		for _, cid := range cids {
			if ci, exists := c.colMap[cid]; exists {
				if ci.IsReportable() {
					cols = append(cols, cid)
				}
			}
		}

		if len(cols) > 0 {
			aliases[alias] = cols
		}
	}

	return aliases
}

// Sortable returns a map of ColIDs to their associated descriptions for the
// sortable columns only.
func (c Cols[P, T]) Sortable() map[ColID]string {
	rptCols := map[ColID]string{}

	for cid, ci := range c.colMap {
		if ci.IsSortable() {
			rptCols[cid] = ci.FullDesc()
		}
	}

	return rptCols
}

// SortableAliases returns a map of aliases to the associated sortable
// columns.
func (c Cols[P, T]) SortableAliases() map[ColID][]ColID {
	aliases := map[ColID][]ColID{}
	maps.Copy(aliases, c.sortableAliases)

	for alias, cids := range c.aliases {
		cols := []ColID{}

		for _, cid := range cids {
			if ci, exists := c.colMap[cid]; exists {
				if ci.IsSortable() {
					cols = append(cols, cid)
				}
			}
		}

		if len(cols) > 0 {
			aliases[alias] = cols
		}
	}

	return aliases
}
