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

	"github.com/alecthomas/kingpin"
)

// Handle is a convenient wrapper around kingpin.Fatalf.
// It takes an error and an arbitrary number of arguments and displays an error
// in the format:
//	limesctl: error: args[0]: args[1]: ... : args[n]: err
func Handle(err error, args ...interface{}) {
	if err == nil {
		return
	}

	if len(args) == 0 {
		kingpin.Fatalf(err.Error())
	} else {
		var format string
		for i := 0; i < len(args); i++ {
			if i == 0 {
				format += "%v"
			} else {
				format += ": %v"
			}
		}
		kingpin.Fatalf("%s: %v", fmt.Sprintf(format, args...), err)
	}
}

// New is a wrapper around errors.New(). Use this to avoid nameclash with the
// local errors package and the one in the standard library.
func New(str string) error {
	return errors.New(str)
}
