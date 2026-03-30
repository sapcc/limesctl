// SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"testing"

	"github.com/sapcc/go-api-declarations/limes"
	"github.com/sapcc/go-bits/assert"
	"github.com/sapcc/go-bits/must"
)

func TestValueFormatter(t *testing.T) {
	// DefaultValueFormatter just formats everything as-is
	assert.Equal(t, DefaultValueFormatter(0), "0")
	assert.Equal(t, DefaultValueFormatter(42), "42")
	assert.Equal(t, DefaultValueFormatter(1<<20), "1048576")

	// PickHumanizedValueFormatter should pick optimal units
	u, f := PickHumanizedValueFormatter(limes.UnitMebibytes, []uint64{2, 1024})
	assert.Equal(t, u.String(), "MiB")
	assert.Equal(t, f(2), "2")
	assert.Equal(t, f(1024), "1024")

	u, f = PickHumanizedValueFormatter(limes.UnitMebibytes, []uint64{2048, 1024})
	assert.Equal(t, u.String(), "GiB")
	assert.Equal(t, f(2048), "2")
	assert.Equal(t, f(1024), "1")

	// PickHumanizedValueFormatter should work with non-standard units and turn them into "clean" units
	weirdUnit := must.Return(limes.UnitMebibytes.MultiplyBy(2032))
	u, f = PickHumanizedValueFormatter(weirdUnit, []uint64{2, 4})
	assert.Equal(t, u.String(), "MiB")
	assert.Equal(t, f(2), "4064")
	assert.Equal(t, f(4), "8128")

	u, f = PickHumanizedValueFormatter(weirdUnit, []uint64{64, 128})
	assert.Equal(t, u.String(), "GiB")
	assert.Equal(t, f(64), "127")
	assert.Equal(t, f(128), "254")
}
