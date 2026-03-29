package rptmaker

import (
	"github.com/nickwells/col.mod/v6/col"
)

type P struct{}

type T struct {
	A int
	B string
}

type ColsAddInfo struct {
	CID             ColID
	CI              *ColInfo[P, T]
	CommonAlias     ColID
	ReportableAlias ColID
	SortableAlias   ColID
	AliasVals       []ColID
}

// MakeTestCols constructs a Cols from the given information
func MakeTestCols(addInfo []ColsAddInfo) (
	*Cols[P, T],
	error,
) {
	c := NewCols[P, T]()

	var err error

	for _, ai := range addInfo {
		if ai.CID != "" {
			err = c.Add(ai.CID, ai.CI)
			if err != nil {
				break
			}
		}

		if ai.CommonAlias != "" {
			err = c.AddAlias(ai.CommonAlias,
				ai.AliasVals[0], ai.AliasVals[1:]...)
			if err != nil {
				break
			}
		}

		if ai.SortableAlias != "" {
			err = c.AddSortableAlias(ai.SortableAlias,
				ai.AliasVals[0], ai.AliasVals[1:]...)
			if err != nil {
				break
			}
		}

		if ai.ReportableAlias != "" {
			err = c.AddReportableAlias(ai.ReportableAlias,
				ai.AliasVals[0], ai.AliasVals[1:]...)
			if err != nil {
				break
			}
		}
	}

	return c, err
}

const (
	ColName     = "colName"
	RepColName  = "reportableCol"
	SortColName = "sortableCol"
	SRColName   = "sortableAndReportable"

	CIDesc    = "column description."
	DescInMap = CIDesc + " This column is unheaded."

	Alias1           = "alias1"
	SortableAlias1   = "sortableAlias1"
	ReportableAlias1 = "reportableAlias1"
)

var (
	CmpVals = func(_, _ T) int { return 0 }
	MkCol   = func(_ P, _ []string) *col.Col { return nil }
	ColVal  = func(_ T) any { return 0 }

	SortableCI          = NewColInfo[P, T](CIDesc, nil, nil, nil, CmpVals)
	ReportableCI        = NewColInfo[P, T](CIDesc, nil, MkCol, ColVal, nil)
	BadCINoFunc         = NewColInfo[P, T](CIDesc, nil, nil, nil, nil)
	SAndRCI             = NewColInfo[P, T](CIDesc, nil, MkCol, ColVal, CmpVals)
	BadSAndRCINoMkCol   = NewColInfo[P, T](CIDesc, nil, nil, ColVal, CmpVals)
	BadSAndRCINoColVal  = NewColInfo[P, T](CIDesc, nil, MkCol, nil, CmpVals)
	BadSAndRCINoCmpVals = NewColInfo[P, T](CIDesc, nil, MkCol, ColVal, nil)
)
