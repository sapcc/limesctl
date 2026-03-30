// SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"strconv"

	"github.com/sapcc/go-api-declarations/limes"
)

// ValueFormatter is the interface for a function that renders integer values into their output representation in CSV.
// This may involve a conversion between units, based on the unit that was chosen for the respective CSV row during humanization.
type ValueFormatter func(uint64) string

// DefaultValueFormatter is a ValueFormatter that does not perform unit conversion.
func DefaultValueFormatter(value uint64) string {
	return strconv.FormatUint(value, 10)
}

var (
	allCleanUnitsBasedOnNone = []limes.Unit{
		limes.UnitNone,
	}
	allCleanUnitsBasedOnBytes = []limes.Unit{
		// This is sorted in descending order of size, so the first viable choice
		// will be the "best" choice (as in: producing the smallest number values).
		// e.g. 2 GiB is better than 2048 MiB is better than 2097152 KiB
		limes.UnitExbibytes,
		limes.UnitPebibytes,
		limes.UnitTebibytes,
		limes.UnitGibibytes,
		limes.UnitMebibytes,
		limes.UnitKibibytes,
		limes.UnitBytes,
	}
)

// PickHumanizedValueFormatter takes several values that are to be interpreted in the given unit,
// and chooses the best unit to format them in for human-readable values.
//
// For example, `PickHumanizedValueFormatter([2048 4096], UnitBytes)` returns
// UnitKibibytes and a formatter that renders 2048 as 2 and 4096 as 4.
func PickHumanizedValueFormatter(unit limes.Unit, values []uint64) (limes.Unit, ValueFormatter) {
	// possible choices of result units
	baseUnit, multiplierToBase := unit.Base()
	var possibleUnits []limes.Unit
	switch baseUnit {
	case limes.UnitNone:
		possibleUnits = allCleanUnitsBasedOnNone
	case limes.UnitBytes:
		possibleUnits = allCleanUnitsBasedOnBytes
	default:
		// if Limes adds new base units that we don't know about, not humanizing is a safe fallback
		return unit, DefaultValueFormatter
	}

	// pick the first candidate unit that produces clean integers for all presented values
UNIT:
	for _, targetUnit := range possibleUnits {
		targetBaseUnit, multiplierToTarget := targetUnit.Base()
		if targetBaseUnit != baseUnit {
			// defense in depth: the lists of `possibleUnits` above should be prepared such that all the base units match
			continue
		}
		for _, value := range values {
			rawValue := value * multiplierToBase
			if rawValue/multiplierToBase != value {
				// if conversion to the base unit runs into an integer overflow, not humanizing is a safe fallback
				return unit, DefaultValueFormatter
			}
			// the value in the target unit would be `rawValue / multiplierToTarget`, does this work within uint64?
			if rawValue%multiplierToTarget != 0 {
				continue UNIT
			}
		}
		// all `values` could be converted cleanly into `targetUnit`
		return targetUnit, func(value uint64) string {
			convertedValue := (value * multiplierToBase) / multiplierToTarget
			return strconv.FormatUint(convertedValue, 10)
		}
	}

	// defense in depth: if (somehow!) no candidate was viable, not humanizing is a safe fallback
	return unit, DefaultValueFormatter
}
