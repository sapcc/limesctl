// Copyright 2022 SAP SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
