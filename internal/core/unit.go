package core

import "github.com/sapcc/limes"

// LimesUnits is a slice of units that Limes understands.
// Note: the units in this slice **must** be in ascending order.
var LimesUnits = []limes.Unit{
	limes.UnitBytes,
	limes.UnitKibibytes,
	limes.UnitMebibytes,
	limes.UnitGibibytes,
	limes.UnitTebibytes,
	limes.UnitPebibytes,
	limes.UnitExbibytes,
}

// ResourceUnits is a map of service name to resource name to the resource's
// default unit in Limes.
type ResourceUnits map[string]map[string]limes.Unit
