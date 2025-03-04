package col

import (
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

// rowsAreEqual compares h with h2 and returns true if all the fields
// are the same and false otherwise
func (h Header) rowsAreEqual(h2 Header) bool {
	if len(h.headerRows) != len(h2.headerRows) {
		return false
	}

	for i, r := range h.headerRows {
		if r != h2.headerRows[i] {
			return false
		}
	}

	return true
}

// isEqual compares h with h2 and returns true if all the fields
// are the same and false otherwise
func (h Header) isEqual(h2 Header) bool {
	if h.underlineCh != h2.underlineCh {
		return false
	}

	if len(h.headerRows) != len(h2.headerRows) {
		return false
	}

	if h.dataRowsPrinted != h2.dataRowsPrinted {
		return false
	}

	if h.repeatHdrInterval != h2.repeatHdrInterval {
		return false
	}

	if h.headerRowCount != h2.headerRowCount {
		return false
	}

	if h.printHdr != h2.printHdr {
		return false
	}

	if h.hdrPrinted != h2.hdrPrinted {
		return false
	}

	if h.underlineHdr != h2.underlineHdr {
		return false
	}

	if !h.rowsAreEqual(h2) {
		return false
	}

	return true
}

func TestHdrCreate(t *testing.T) {
	dfltHdr := Header{
		underlineCh:       "=",
		headerRows:        nil,
		dataRowsPrinted:   0,
		repeatHdrInterval: 0,
		headerRowCount:    0,
		printHdr:          true,
		hdrPrinted:        false,
		underlineHdr:      true,
	}

	dontPrintHdr := dfltHdr
	dontPrintHdr.printHdr = false

	dontUnderlineHdr := dfltHdr
	dontUnderlineHdr.underlineHdr = false

	underlineHdrWith := dfltHdr
	underlineHdrWith.underlineCh = "X"

	repeatHdr1 := dfltHdr
	repeatHdr1.repeatHdrInterval = 1

	repeatHdr3 := dfltHdr
	repeatHdr3.repeatHdrInterval = 3

	const (
		badHdrRepeat    = "the header repeat count (0) must be >= 1"
		badHdrUnderline = "the header underline rune (U+0000) must be printable"
	)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		testhelper.ExpPanic
		hdrOpts         []HdrOptionFunc
		expectedHdr     Header
		expHeaderOutput string
	}{
		{
			ID:          testhelper.MkID("default"),
			hdrOpts:     []HdrOptionFunc{},
			expectedHdr: dfltHdr,
		},
		{
			ID:          testhelper.MkID("dont print"),
			hdrOpts:     []HdrOptionFunc{HdrOptDontPrint},
			expectedHdr: dontPrintHdr,
		},
		{
			ID:          testhelper.MkID("dont underline"),
			hdrOpts:     []HdrOptionFunc{HdrOptDontUnderline},
			expectedHdr: dontUnderlineHdr,
		},
		{
			ID:          testhelper.MkID("underline with"),
			hdrOpts:     []HdrOptionFunc{HdrOptUnderlineWith('X')},
			expectedHdr: underlineHdrWith,
		},
		{
			ID:          testhelper.MkID("good repeat header: 1"),
			hdrOpts:     []HdrOptionFunc{HdrOptRepeat(1)},
			expectedHdr: repeatHdr1,
		},
		{
			ID:          testhelper.MkID("good repeat header: 3"),
			hdrOpts:     []HdrOptionFunc{HdrOptRepeat(3)},
			expectedHdr: repeatHdr3,
		},
		{
			ID:       testhelper.MkID("bad repeat header: zero"),
			hdrOpts:  []HdrOptionFunc{HdrOptRepeat(0)},
			ExpErr:   testhelper.MkExpErr(badHdrRepeat),
			ExpPanic: testhelper.MkExpPanic(badHdrRepeat),
		},
		{
			ID:       testhelper.MkID("bad underline rune"),
			hdrOpts:  []HdrOptionFunc{HdrOptUnderlineWith(rune(0))},
			ExpErr:   testhelper.MkExpErr(badHdrUnderline),
			ExpPanic: testhelper.MkExpPanic(badHdrUnderline),
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			_ = NewHeaderOrPanic(tc.hdrOpts...)
		})
		testhelper.CheckExpPanicError(t, panicked, panicVal, tc)

		h, err := NewHeader(tc.hdrOpts...)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			if !h.isEqual(tc.expectedHdr) {
				t.Log(tc.IDStr())
				t.Logf("\t: expected header: %v\n", tc.expectedHdr)
				t.Logf("\t:   actual header: %v\n", h)
				t.Error("\t: header is incorrect\n")

				continue
			}
		}
	}
}
