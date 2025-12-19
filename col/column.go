package col

import "fmt"

// Justification represents how a column is justified
type Justification int

// The justification types:
//
//	Left means left-justified
//	Right means right-justified
const (
	Left Justification = iota
	Right
)

// DfltColSep is the default column separator
const DfltColSep = " "

// Col holds the values needed in order to represent a column
type Col struct {
	headers    []string
	f          Formatter
	finalWidth int
	sep        string
}

// New creates a new Col object
func New(f Formatter, colHead ...string) *Col {
	c := &Col{
		headers: colHead,
		f:       f,
		sep:     DfltColSep,
	}

	if len(c.headers) == 0 {
		c.headers = make([]string, 1)
	}

	c.finalWidth = int(f.Width())

	return c
}

// SetSep sets the separator for the column from the default value (see
// DfltColSep) to the value passed
func (c *Col) SetSep(s string) *Col {
	c.sep = s
	return c
}

// hdrText returns the text of the header corresponding to the given row
// If the row is before the start of the headers for that column then the
// empty string is returned
func (c Col) hdrText(rowIdx, rowCount int) string {
	valIdx := len(c.headers) - rowCount + rowIdx
	if valIdx < 0 || valIdx >= len(c.headers) {
		return ""
	}

	return c.headers[valIdx]
}

// stringInCol returns the string s formatted to fit in the column
func (c Col) stringInCol(s string) string {
	if c.f.Just() == Left {
		return fmt.Sprintf("%-*s", c.finalWidth, s)
	}

	return fmt.Sprintf("%*s", c.finalWidth, s)
}
