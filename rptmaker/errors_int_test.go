package rptmaker

import (
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestFuncNameMeaning(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		fName      string
		expMeaning string
	}{
		{
			ID:         testhelper.MkID("not reportable - mkColFName"),
			fName:      mkColFName,
			expMeaning: "(it is not reportable)",
		},
		{
			ID:         testhelper.MkID("not reportable - colValFName"),
			fName:      colValFName,
			expMeaning: "(it is not reportable)",
		},
		{
			ID:         testhelper.MkID("not sortable - cmpValsFName"),
			fName:      cmpValsFName,
			expMeaning: "(it is not sortable)",
		},
		{
			ID:         testhelper.MkID("unknown func name"),
			fName:      "nonesuch",
			expMeaning: "(error: unknown function)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			testhelper.DiffString(t,
				tc.IDStr(), "meaning",
				funcNameMeaning(tc.fName), tc.expMeaning)
		})
	}
}
