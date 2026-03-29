package rptmaker

import "fmt"

// ColumnErr records details of a problem with a column
type ColumnErr struct {
	Column  ColID
	Problem string
}

// Error returns the string form of the ColumnErr
func (err ColumnErr) Error() string {
	return fmt.Sprintf("column: %q: %s", err.Column, err.Problem)
}

// AliasErr records details of a problem with an alias
type AliasErr struct {
	Alias   ColID
	Problem string
}

// Error returns the string form of the AliasErr
func (err AliasErr) Error() string {
	return fmt.Sprintf("alias: %q: %s", err.Alias, err.Problem)
}

const (
	mkColFName   = "mkCol"
	colValFName  = "colVal"
	cmpValsFName = "cmpVals"
)

// funcNameMeaning returns a string giving the meaning and importance of the
// named function
func funcNameMeaning(fName string) string {
	switch fName {
	case mkColFName:
		return "(it is not reportable)"
	case colValFName:
		return "(it is not reportable)"
	case cmpValsFName:
		return "(it is not sortable)"
	}

	return "(error: unknown function)"
}

// MkNoFuncErr returns a ColumnErr recording a missing function
func MkNoFuncErr(cid ColID, fName string) ColumnErr {
	return ColumnErr{
		Column: cid,
		Problem: fmt.Sprintf("has no %q function %s",
			fName, funcNameMeaning(fName)),
	}
}

// MkUnusableErr returns a ColumnErr describing a column which is neither
// sortable nor reportable.
func MkUnusableErr(cid ColID) ColumnErr {
	return ColumnErr{
		Column:  cid,
		Problem: "is neither sortable nor reportable",
	}
}

// MkColNotFoundErr returns a ColumnErr recording an unknown column
func MkColNotFoundErr(cid ColID) ColumnErr {
	return ColumnErr{
		Column:  cid,
		Problem: "not found",
	}
}

// MkDupColErr returns a ColumnErr recording a duplicate column
func MkDupColErr(cid ColID) ColumnErr {
	return ColumnErr{
		Column:  cid,
		Problem: "duplicate: there is already a column with that name",
	}
}

// MkAliasNameErr returns a ColumnErr recording a column whose name has
// already been used as an alias name.
func MkAliasNameErr(cid ColID) ColumnErr {
	return ColumnErr{
		Column:  cid,
		Problem: "the name is already used as an alias",
	}
}

// MkNoColInfoErr returns a ColumnErr recording a column with no associated
// ColInfo.
func MkNoColInfoErr(cid ColID) ColumnErr {
	return ColumnErr{
		Column:  cid,
		Problem: "no ColInfo has been supplied",
	}
}

// MkAliasIsColErr returns an AliasErr recording an alias name the same as a
// column name.
func MkAliasIsColErr(cid ColID) AliasErr {
	return AliasErr{
		Alias:   cid,
		Problem: "duplicate: there is already a column with that name",
	}
}

// MkDupAliasErr returns a AliasErr recording a duplicate alias
func MkDupAliasErr(cid ColID, aType string) AliasErr {
	return AliasErr{
		Alias:   cid,
		Problem: "duplicate: there is already a " + aType + " alias with that name",
	}
}

// MkAliasErr returns an AliasErr with a bespoke problem.
func MkAliasErr(cid ColID, problem string) AliasErr {
	return AliasErr{
		Alias:   cid,
		Problem: problem,
	}
}
