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

package core

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/sapcc/limesctl/internal/errors"
)

// writeJSON is a helper function that writes the JSON data to os.Stdout.
func (d jsonData) write(w io.Writer) {
	fmt.Fprintln(w, string(d))
}

// Write writes the CSV data to w.
func (d csvData) write(w io.Writer) {
	csvW := csv.NewWriter(w)
	csvW.Comma = rune(';') // use semicolon as delimiter

	err := csvW.WriteAll(d)
	errors.Handle(err, "could not write CSV data")
}

// writeTable is a helper function that writes the CSV data to os.Stdout in an ASCII table format.
func (d csvData) writeTable(w io.Writer) {
	t := tablewriter.NewWriter(w)
	t.SetHeader(d[0])
	t.AppendBulk(d[1:])
	t.Render()
}
