package col

import (
	"fmt"
	"io"
	"strings"
)

// PreHdrFunc is the signature of a function to be called immediately before
// the header itself is printed. It is intended for printing a report preamble
// or, if the header is periodically repeated, it could be used, for instance,
// to print sub-totals. The int64 parameter passes the number of data rows
// printed, if this is zero then the header is being printed for the first time
//
type PreHdrFunc func(io.Writer, uint64)

// Header holds the parameters which control how and when the header is printed
type Header struct {
	underlineCh       string
	headerRows        []string
	dataRowsPrinted   uint64
	repeatHdrInterval uint64
	headerRowCount    int
	preHeaderFunc     PreHdrFunc
	spanDups          bool
	printHdr          bool
	hdrPrinted        bool
	underlineHdr      bool
}

// initVals sets the header row count and records, for each column, the
// final, with-header width. If the header is not to be printed then the
// default column width from the formatter is used as is
// reuse if the header is to be reprinted
func (h *Header) initVals(cols []*Col) {
	if !h.printHdr {
		return
	}

	for _, c := range cols {
		if len(c.headers) > h.headerRowCount {
			h.headerRowCount = len(c.headers)
		}
	}
	h.headerRows = make([]string, h.headerRowCount)
}

// setSpanningCols populates the spans slice with any columns in the
// range start to end which are spanning in the given row
func (h *Header) setSpanningCols(row, start, end int, sg spanGrid) {
	span := span{
		start:   start,
		end:     start,
		row:     row,
		hdrText: sg.cols[start].hdrText(row, h.headerRowCount),
	}
	for i := start + 1; i <= end; i++ {
		hdrText := sg.cols[i].hdrText(row, h.headerRowCount)
		if span.hdrText == hdrText && h.spanDups {
			span.end = i
		} else {
			sg.spans[row] = append(sg.spans[row], span)

			span.start = i
			span.end = i
			span.hdrText = hdrText
		}
	}
	sg.spans[row] = append(sg.spans[row], span)
}

// addUnderlines adds the row of underlines as the last row in the set of
// cached headerRows
func (h *Header) addUnderlines(cols []*Col) {
	if h.underlineHdr {
		underline := ""
		sep := ""
		for _, c := range cols {
			underline += sep
			sep = strings.Repeat(" ", len(c.sep))
			s := c.headers[len(c.headers)-1]
			underline += c.stringInCol(strings.Repeat(h.underlineCh, len(s)))
		}
		h.headerRows = append(h.headerRows, underline)
	}
}

// createHeader creates the header rows and caches them in the Header for
// reuse if the header is to be reprinted
func (h *Header) createHeader(cols []*Col) {
	var sg = newSpanGrid(h, cols)

	if h.headerRowCount > 1 {
		h.setSpanningCols(0, 0, len(cols)-1, sg)

		for row := 1; row < h.headerRowCount-1; row++ {
			for _, span := range sg.spans[row-1] {
				h.setSpanningCols(row, span.start, span.end, sg)
			}
		}
	}

	sg.setLastRowOfHeader()

	sg.setWidths()

	sg.setColWidthFromLastRow()

	h.createHeaderFromSpans(sg)

	h.addUnderlines(cols)
}

func (h *Header) createHeaderFromSpans(sg spanGrid) {
	for row := 0; row < h.headerRowCount; row++ {
		sep := ""
		for _, span := range sg.spans[row] {
			h.headerRows[row] += sep
			sep = strings.Repeat(" ", len(sg.cols[span.end].sep))

			if span.isMultiCol() {
				textWidth := len(span.hdrText)

				if textWidth == 0 {
					h.headerRows[row] += fmt.Sprintf("%*s", span.width, "")
				} else {
					dashCount := (span.width - textWidth) / 2

					h.headerRows[row] +=
						fmt.Sprintf("%s%s%s",
							strings.Repeat("-", dashCount),
							span.hdrText,
							strings.Repeat("-", span.width-textWidth-dashCount))
				}
			} else {
				c := sg.cols[span.start]

				h.headerRows[row] += c.stringInCol(span.hdrText)
			}
		}
	}
}

// printHeader prints the header lines if necessary
func (h *Header) printHeader(w io.Writer, cols []*Col) {
	if !h.printHdr {
		return
	}

	if h.hdrPrinted {
		if h.repeatHdrInterval == 0 {
			return
		}
		if h.dataRowsPrinted%h.repeatHdrInterval != 0 {
			return
		}
	} else {
		h.createHeader(cols)
	}

	if h.preHeaderFunc != nil {
		h.preHeaderFunc(w, h.dataRowsPrinted)
	}
	for _, hr := range h.headerRows {
		fmt.Fprintln(w, hr)
	}

	h.hdrPrinted = true
}

// HdrOptionFunc is the signature of the function that is passed to the
// NewHeader function to set the header options
type HdrOptionFunc func(*Header) error

// HdrOptPreHdrFunc returns a HdrOptionFunc that will set the pre-header
// function
func HdrOptPreHdrFunc(f PreHdrFunc) HdrOptionFunc {
	return func(h *Header) error {
		h.preHeaderFunc = f
		return nil
	}
}

// HdrOptDontPrint prevents the header from being printed
func HdrOptDontPrint(h *Header) error {
	h.printHdr = false
	return nil
}

// HdrOptDontSpanDups prevents the header from spanning common headers
func HdrOptDontSpanDups(h *Header) error {
	h.spanDups = false
	return nil
}

// HdrOptDontUnderline prevents the header from being underlined
func HdrOptDontUnderline(h *Header) error {
	h.underlineHdr = false
	return nil
}

// HdrOptUnderlineWith returns a HdrOptionFunc that will set the rune used to
// underline the final header line
func HdrOptUnderlineWith(r rune) HdrOptionFunc {
	return func(h *Header) error {
		h.underlineCh = string(r)
		return nil
	}
}

// HdrOptRepeat returns a HdrOptionFunc that will set the number of
// lines of data that should be printed before the header is printed again. If
// this value is not set then the header is only printed once
func HdrOptRepeat(n uint64) HdrOptionFunc {
	return func(h *Header) error {
		if n < 1 {
			return fmt.Errorf(
				"the value for repeating the header must be >= 1: %d",
				n)
		}
		h.repeatHdrInterval = n
		return nil
	}
}

// NewHeader creates a new Header object
func NewHeader(options ...HdrOptionFunc) (*Header, error) {
	h := &Header{
		spanDups:     true,
		printHdr:     true,
		underlineHdr: true,
		underlineCh:  "=",
	}

	for _, o := range options {
		err := o(h)
		if err != nil {
			return nil, err
		}
	}

	return h, nil
}

// incrDataRowsPrinted increments the dataRowsPrinted (for defer)
func (h *Header) incrDataRowsPrinted() {
	h.dataRowsPrinted++
}
