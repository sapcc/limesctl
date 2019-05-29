/*******************************************************************************
*
* Copyright 2018 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package errors

import (
	"errors"
	"fmt"
	"os"
)

// Handle takes an error and a string writes an error message to os.Stderr.
func Handle(err error, str string) {
	if err != nil {
		msg := "limesctl: error"
		if str != "" {
			msg += ": " + str
		}
		fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
		os.Exit(1)
	}
}

// New is a wrapper around errors.New(). Use this to avoid nameclash with the
// local errors package and the one in the standard library.
func New(str string) error {
	return errors.New(str)
}
