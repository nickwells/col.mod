package col

// span holds details of a column span. A column span is a set of columns in a
// row all having the same header text and all belonging to the same span in
// the row above. For instance given the following columns:
//
//	Column number: 1         2          3
//	        row 1: average   average    total
//	        row 2: height    weight     income
//
// columns 1 and 2 form a span in row 1. Column 3 in row 1 and all columns in
// row 2 are single-column spans where the start and end fields are the same.
// If the column in the row above is in a single-column span then all the rows
// below will have that same column in a single-column span too.
type span struct {
	start, end int
	row        int
	hdrText    string
	width      int
	sepWidth   int
}

// isMultiCol returns true if the span represents multiple columns
func (s span) isMultiCol() bool {
	return s.start != s.end
}

// minWidth returns the minimum width of the span
func (s span) minWidth() int {
	w := len(s.hdrText)
	if s.isMultiCol() && w > 0 {
		// at least one hyphen at each end of the multi-column
		// span (but not for nameless spans, where w == 0)
		w += 2
	}

	return w
}

// spanGrid represents the row-by-row set of spans in the header
type spanGrid struct {
	spans [][]span
	cols  []*Col
}

// newSpanGrid creates a new spanGrid
func newSpanGrid(h *Header, cols []*Col) spanGrid {
	spans := make([][]span, h.headerRowCount)
	for i := range spans {
		spans[i] = make([]span, 0, len(cols))
	}
	return spanGrid{
		spans: spans,
		cols:  cols,
	}
}

// totalWidth gets the total width of all the spans in the row between start
// and end. It ignores any span before start, stops when it gets to the first
// span after end and adds extra space for every gap between spans
func (sg spanGrid) totalWidth(row, start, end int) int {
	w := 0
	gapIncr := 0
	for _, span := range sg.spans[row] {
		if span.end < start {
			continue
		}
		if span.start > end {
			break
		}
		w += gapIncr + span.width
		gapIncr = len(sg.cols[span.end].sep)
	}
	return w
}

// setWidths works out the width of each span in each row
func (sg spanGrid) setWidths() {
	sg.setWidthsFromSpansBelow()
	sg.setWidthsFromSpansAbove()
}

// setWidthsFromSpansBelow works out the width of each span in each row
func (sg spanGrid) setWidthsFromSpansBelow() {
	for row := len(sg.spans) - 2; row >= 0; row-- {
		for i, span := range sg.spans[row] {
			w := span.minWidth()
			nextRowWidth := sg.totalWidth(row+1, span.start, span.end)
			if nextRowWidth > w {
				w = nextRowWidth
			}
			span.width = w
			sg.spans[row][i] = span
		}
	}
}

func (sg spanGrid) setWidthsFromSpansAbove() {
	for row := 1; row <= len(sg.spans)-1; row++ {
		for _, span := range sg.spans[row-1] {
			w := sg.totalWidth(row, span.start, span.end)
			if span.width > w {
				// adjust the span widths below to take account of any extra
				// space from the spanning columns above
				extraSpaces := span.width - w
				colCount := 1 + span.end - span.start
				perCol := extraSpaces / colCount
				oneExtraCount := extraSpaces - (perCol * colCount)

				for i, spanBelow := range sg.spans[row] {
					if spanBelow.end < span.start {
						continue
					}
					if spanBelow.start > span.end {
						break
					}

					count := 1 + spanBelow.end - spanBelow.start
					spanBelow.width += (count * perCol)
					if oneExtraCount > 0 {
						spanBelow.width += count
						oneExtraCount -= count
					}
					sg.spans[row][i] = spanBelow
				}
			}
		}
	}
}

// setLastRowOfHeader sets the initial value of the last row of the header.
// This is the (non-underline) row just above the data values in the printed
// report and we do not span these headings so there is one entry per column
func (sg spanGrid) setLastRowOfHeader() {
	row := len(sg.spans) - 1
	for i, c := range sg.cols {
		span := span{
			start:    i,
			end:      i,
			row:      row,
			hdrText:  c.hdrText(row, len(sg.spans)),
			width:    c.finalWidth,
			sepWidth: len(c.sep),
		}
		if len(span.hdrText) > span.width {
			span.width = len(span.hdrText)
		}
		sg.spans[row] = append(sg.spans[row], span)
	}
}

func (sg spanGrid) setColWidthFromLastRow() {
	row := len(sg.spans) - 1
	for i, c := range sg.cols {
		c.finalWidth = sg.spans[row][i].width
		sg.cols[i] = c
	}
}
