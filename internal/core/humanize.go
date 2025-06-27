// SPDX-FileCopyrightText: 2005
// 2008 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"math"
	"slices"
	"strconv"

	"github.com/sapcc/go-api-declarations/limes"
)

// suitableUnit takes a value and its current unit, and returns a suitable
// human friendly unit that this value can be converted to.
//
// E.g. A value 82854982 with unit limes.UnitBytes will return
// limes.UnitMebibytes.
func suitableUnit(v uint64, u limes.Unit) limes.Unit {
	if v < 10 {
		return u
	}

	var sizes []limes.Unit
	for i, unit := range LimesUnits {
		if u == unit {
			sizes = LimesUnits[i:]
			break
		}
	}

	base := 1024.0
	logn := math.Log(float64(v)) / math.Log(base)
	e := math.Floor(logn)
	return sizes[int(e)]
}

// convertValue converts a value from its old unit to the new unit.
// The returned value is rounded off.
//
// Use suitableUnit() to find an appropriate new unit for a specific value.
func convertValue(v uint64, source, target limes.Unit) string {
	if v == 0 {
		return "0"
	}
	_, sourceMultiple := source.Base()
	_, targetMultiple := target.Base()
	valueInBase := v * sourceMultiple
	newV := float64(valueInBase) / float64(targetMultiple)
	newV = math.Round(newV*100) / 100 // round to second decimal place
	return strconv.FormatFloat(newV, 'f', -1, 64)
}

// valToStrFunc is a function that can be used to convert a resource value to a
// string.
type valToStrFunc func(v uint64) string

var defaultValToStrFunc valToStrFunc = func(v uint64) string { return strconv.FormatUint(v, 10) }

// getValToStrFunc finds a human friendly unit for the given vals and returns a
// function that can be used to convert resource values to this new unit.
func getValToStrFunc(humanize bool, old limes.Unit, vals []uint64) (f valToStrFunc, newUnit limes.Unit) {
	f = defaultValToStrFunc
	newUnit = old
	if humanize && old != limes.UnitNone {
		// Find a new human friendly unit based on the smallest value.
		if sml := smallestValue(vals); sml > 0 {
			if newUnit = suitableUnit(sml, old); newUnit != old {
				f = func(v uint64) string {
					return convertValue(v, old, newUnit)
				}
			}
		}
	}
	return f, newUnit
}

func smallestValue(vals []uint64) uint64 {
	nonZero := make([]uint64, 0, len(vals))
	for _, v := range vals {
		if v > 0 {
			nonZero = append(nonZero, v)
		}
	}
	if len(nonZero) == 0 {
		return 0
	}
	slices.Sort(nonZero)
	return nonZero[0]
}
