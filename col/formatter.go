package col

// Formatter is an interface which describes the methods to be provided by a
// column formatter. Various instances of the Formatter interface are given
// in the colfmt package. These should cover many common requirements.
type Formatter interface {
	// Formatted should return the value as a string
	Formatted(any) string
	// Width should return the expected width of the string printed with the
	// format string. Note that the actual width of the string may be greater
	// than this depending on the width of the column header
	Width() int
	// Just should return whether the resultant string is left or right
	// justified. This information is needed when deciding how to print the
	// header
	Just() Justification
	// Check should return an error if a formatter has invalid
	// configuration. It is called on each of the supplied Formatters (in the
	// Cols) before a Report is returned.
	Check() error
}
