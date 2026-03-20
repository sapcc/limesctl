// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package util

import "fmt"

// WrapError contains a new error, constructed from a prefix string and the original error.
func WrapError(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// CastStringsTo is a generic function which casts all values of an array to a defined
// output type and returns them as array.
func CastStringsTo[O ~string, I ~string](input []I) (output []O) {
	output = make([]O, len(input))
	for idx, val := range input {
		output[idx] = O(val)
	}
	return
}
