package tpb

import "strconv"

// OrderBy represent the different values the search can be ordered by
type OrderBy int

// List of the different Orders available
const (
	OrderByName OrderBy = iota
	OrderByDate
	OrderBySize
	OrderBySeeds
	OrderByLeeches
)

// SortOrder represents the sort order
type SortOrder int

// List of the different sort order
const (
	Desc SortOrder = 1 + iota
	Asc
)

// mapOrderBy takes the orderBy and sort order parameter and return the
// corresponding option to pass to the website
func mapOrderBy(orderBy OrderBy, order SortOrder) string {
	return strconv.Itoa(int(orderBy)*2 + int(order))
}
