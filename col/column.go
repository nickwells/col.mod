package col

import "fmt"

// Justification represents how a column is justified
type Justification int

// The justification types:
//    Left means left-justified
//    Right means right-justified
const (
	Left Justification = iota
	Right
)

// Formatter is an interface describes the methods to be provided by a column
// formatter.
type Formatter interface {
	// Formatted should return the value as a string
	Formatted(interface{}) string
	// Width should return the expected width of the string printed with the
	// format string. Note that the actual width of the string may be greater
	// than this depending on the width of the column header
	Width() int
	// Just should return whether the resultant string is left or right
	// justified. This information is needed when deciding how to print the
	// header
	Just() Justification
}

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
		sep:     " ",
	}

	if len(c.headers) == 0 {
		c.headers = make([]string, 1)
	}

	c.finalWidth = f.Width()

	return c
}

// SetSep sets the separator for the column from the default value (" ")
// to the value passed
func (c *Col) SetSep(s string) *Col {
	c.sep = s
	return c
}

// hdrText returns the text of the header corresponding to the given row
// If the row is before the start of the headers for that column then the
// empty string is returned
func (c Col) hdrText(rowIdx, rowCount int) string {
	val := ""
	valIdx := len(c.headers) - rowCount + rowIdx
	if valIdx >= 0 && valIdx < len(c.headers) {
		val = c.headers[valIdx]
	}
	return val
}

// stringInCol returns the string s formatted to fit in the column
func (c Col) stringInCol(s string) string {
	if c.f.Just() == Left {
		return fmt.Sprintf("%-*s", c.finalWidth, s)
	}

	return fmt.Sprintf("%*s", c.finalWidth, s)
}
