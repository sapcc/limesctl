// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package util

import "fmt"

func WrapError(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func CastStringsTo[O ~string, I ~string](input []I) (output []O) {
	output = make([]O, len(input))
	for idx, val := range input {
		output[idx] = O(val)
	}
	return
}
