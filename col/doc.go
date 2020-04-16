/*
Package col helps you to print a nicely formatted report to the terminal.

You can specify columns with multi-line headers and it will automatically align
the headers correctly for the column and format the data to fit neatly under
the headers.

You start by creating a header object, this can be used to set various
different behaviours - see the HdrOpt... functions for details of what options
are available. This cannot be used for anything except passing to the
constructor for the report object

Then you create a report object passing it the newly created header object and
a list of one or more columns.

Then to make use of all your work above you call the PrintRow... methods on the
report object that you just created.

*/
package col
