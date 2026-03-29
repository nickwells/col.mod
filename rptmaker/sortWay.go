package rptmaker

import "slices"

// SortWay is the type of a column tag used to specify the direction of
// a comparison function. a column selected for sorting can be sorted in
// reverse order by decorating the column with the appropriate tag.
type SortWay string

const (
	// Forwards indicates that the standard ordering should be used
	Forwards = SortWay("forwards")
	// Backwards indicates that reverse ordering should be used
	Backwards = SortWay("backwards")

	// Reverse is an alias for Backwards
	Reverse = SortWay("reverse")
	// Rev is an alias for Backwards
	Rev = SortWay("rev")
	// Back is an alias for Backwards
	Back = SortWay("back")

	// Fwd is an alias for Forwards
	Fwd = SortWay("fwd")
	// Forward is an alias for Forwards
	Forward = SortWay("forward")
)

// AllowedSortDirections returns a map of SortDirections to
// descriptions. This is suitable to be used as an allowed values constraint
// for a parameter setter.
func AllowedSortDirections() map[SortWay]string {
	return map[SortWay]string{
		Forwards:  "in ascending order",
		Backwards: "in descending order",
	}
}

// SortDirectionAliases returns a map of alias values to SortDirections. This is
// suitable to be used as an alias values constraint for a parameter setter.
func SortDirectionAliases() map[SortWay][]SortWay {
	return map[SortWay][]SortWay{
		Fwd:     {Forwards},
		Forward: {Forwards},

		Reverse: {Backwards},
		Rev:     {Backwards},
		Back:    {Backwards},
	}
}

// SortColumn records a column for sorting and flag indicating whether the
// comparison function should be reversed.
type SortColumn struct {
	ID        ColID
	Backwards bool
}

// MakeSortColumn constructs a [SortColumn] from a [ColID] and a slice of
// [SortWay] values. If Backwards appears anywhere in the tags slice then the
// sort order will be in descending order; the default is ascending order (as
// given by the comparison function in [ColInfo]).
func MakeSortColumn(id ColID, tags []SortWay) SortColumn {
	return SortColumn{
		ID:        id,
		Backwards: slices.Contains(tags, Backwards),
	}
}
